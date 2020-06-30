package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/cxqlkk/library/discovery/consul/env"
	"github.com/cxqlkk/library/encode"
	"github.com/cxqlkk/library/log"
	consulapi "github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"os"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	_configTimeout           = 60 * time.Second
	_configCheckInterval     = 10 * time.Second
	_configInitRetryTime     = 3
	_configInitRetryInterval = 1 * time.Second
)

type Config struct {
	conf          interface{}
	consuleClient *consulapi.Client
	//data          atomic.Value
	lastIndex uint64
}

func InitConfig(conf interface{}) error {
	if env.Conf == "" {
		return RemoteConfig(conf)
	}
	return LocalConfig(conf)
}

func LocalConfig(conf interface{}) error {
	if _, err := toml.DecodeFile(env.Conf+"/application.toml", conf); err != nil {
		log.Logger.Error("初始化配置失败", zap.Error(err))
		return err
	}
	return nil
}

func RemoteConfig(conf interface{}) error {
	config := &Config{
		consuleClient: GetConsulClient(),
		conf:          conf,
	}
	if err := config.init(); err != nil {
		log.Logger.Error("初始化配置失败", zap.Error(err))
		return fmt.Errorf("%w", err)
	}
	go config.checkConfig()
	return nil
}

func GetConsulClient() *consulapi.Client {
	config := consulapi.DefaultConfig()
	config.Address = env.ConsulAddr
	if client, err := consulapi.NewClient(config); err != nil {
		panic("getConsulClient()->error:配置中心连接失败")
	} else {
		return client
	}
}

func (conf *Config) checkConfig() {
	//configTimer := time.NewTicker(_configCheckInterval)
	//for {
	//	select {
	//	case <-configTimer.C:
	//		if _, meta, err := conf.consuleClient.KV().List(env.KVPrefix, &consulapi.QueryOptions{
	//			WaitIndex:         conf.lastIndex,
	//			WaitTime:          time.Second*100,
	//		}); err != nil {
	//			log.Logger.Error("conf.checkConfig()", zap.Error(err))
	//		} else if conf.lastIndex != meta.LastIndex {
	//			syscall.Kill(os.Getpid(), syscall.SIGHUP)
	//		}
	//	}
	//}
	for {
		log.Logger.Info("开始检查配置更新", zap.String("now", time.Now().Format("2006-01-02 15:04:05")))
		if _, meta, err := conf.consuleClient.KV().List(env.KVPrefix, &consulapi.QueryOptions{
			WaitIndex: conf.lastIndex,
			WaitTime:  time.Second * 100,
		}); err != nil {
			log.Logger.Error("conf.checkConfig()", zap.Error(err))
		} else if conf.lastIndex != meta.LastIndex {
			log.Logger.Info("配置检查结束，存在更新", zap.String("now", time.Now().Format("2006-01-02 15:04:05")),zap.Int("conf.lastindex",int(conf.lastIndex)),zap.Int("meta.lastindex",int(meta.LastIndex)))
			syscall.Kill(os.Getpid(), syscall.SIGHUP)
		}
		log.Logger.Info("配置检查结束，无更新", zap.String("now", time.Now().Format("2006-01-02 15:04:05")))
		time.Sleep(_configCheckInterval)
	}
}

func (conf *Config) init() (err error) {
	var kvs consulapi.KVPairs
	var meta *consulapi.QueryMeta
	for i := 0; i < _configInitRetryTime; i++ {

		if kvs, meta, err = conf.consuleClient.KV().List(env.KVPrefix, nil); err == nil {
			conf.lastIndex = meta.LastIndex
			return conf.load(kvs)
		}
		log.Logger.Error("conf.init() error", zap.Error(encode.InterfaceError))
		time.Sleep(_configInitRetryInterval)
	}
	return
}

func (conf *Config) load(kvpairs consulapi.KVPairs) error {
	mp := map[string]string{}
	for _, kv := range kvpairs {
		mp[kv.Key] = string(kv.Value)
	}
	return conf.parse([]string{env.KVPrefix}, reflect.ValueOf(conf.conf), mp)
}

func (conf *Config) parse(prefix []string, rv reflect.Value, kvs map[string]string) error {

	rv = reflect.Indirect(rv)
	if rv.Kind() == reflect.Interface {
		rv = reflect.New(rv.Elem().Type())
	}
	rt := rv.Type() //获取tag
	for i := 0; i < rv.NumField(); i++ {
		fieldv := rv.Field(i)
		fmt.Println(rt.Field(i).Name)
		if tag := rt.Field(i).Tag.Get("dc"); tag != "" {
			if fieldv.Kind() == reflect.Ptr {
				newGuy := reflect.New(fieldv.Type().Elem())
				fieldv.Set(newGuy)
				fieldv = fieldv.Elem()
			}
			switch fieldv.Kind() {
			case reflect.Struct:
				part := append(prefix, tag)
				if err := conf.parse(part, fieldv, kvs); err != nil {
					return err
				}

			default:
				part := append(prefix, tag)
				value := kvs[strings.Join(part, "/")]
				fmt.Println(strings.Join(part, "/"), value)
				if fieldv.Kind() == reflect.Int {
					if value, err := strconv.Atoi(value); err != nil {
						fieldv.SetInt(int64(value))
					} else {
						log.Logger.Error("conf.parse()", zap.Error(fmt.Errorf("目标类型为int")))
						return encode.ServerErr
					}
				} else if fieldv.Kind() == reflect.String {
					fieldv.SetString(value)
				} else {
					log.Logger.Error("conf.parse() ", zap.Error(fmt.Errorf("不支持的解析类型%v", fieldv.Kind())))
					return encode.ServerErr
				}
			}

		}
	}
	return nil
}
