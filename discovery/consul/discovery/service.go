package discovery

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go-library/discovery/consul/config"
	errors "go-library/encode"
	"go-library/log"
	"go.uber.org/zap"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	_serviceInitRetryTime     = 3
	_serviceCheckInterval     = time.Second * 10
	_serviceInitRetryInterval = 1 * time.Second
)

type Service struct {
	value        atomic.Value
	consulClient *api.Client
	lastIndex    uint64
}

func NewService()*Service {
	var srv = &Service{
		value:        atomic.Value{},
		consulClient: config.GetConsulClient(),
	}
	if err := srv.initService(); err != nil {
		panic(fmt.Errorf("初始化service 错误:%v", err))
	}
	go srv.serviceCheck()
	return srv
}

func (srv *Service) GetService(serviceName string) (addr string, err error) {
	services := srv.value.Load().(map[string]string)
	if addr, ok := services[serviceName]; ok {
		return addr, nil
	}
	return "", errors.ServiceNotFound
}

func (srv *Service) initService() (err error) {
	var services = map[string][]string{}
	var meta *api.QueryMeta
	for i := 0; i < _serviceInitRetryTime; i++ {
		if services, meta, err = srv.consulClient.Catalog().Services(nil); err == nil {
			srv.lastIndex = meta.LastIndex
			return srv.loadServices(services)
		}
		log.Logger.Error("srv.initService()", zap.Error(err))
	}
	return err
}
func (srv *Service) loadServices(serves map[string][]string) (err error) {
	var serviceMap= map[string]string{}
	for k, _ := range serves {
		if entry, _, err := srv.consulClient.Health().Service(k, "", true, nil); err != nil {
			log.Logger.Error("srv.loadServices()", zap.Error(err))
			return err
		} else if len(entry)>0 {
			//todo fixme
			fmt.Println(entry[0].Service.Address,
				entry[0].Service.Port)
			serviceMap[k] = entry[0].Service.Address + ":" + strconv.Itoa(entry[0].Service.Port)
		}
	}
	srv.value.Store(serviceMap)
	return nil
}

func (srv *Service) serviceCheck() {
	configTimer := time.NewTicker(_serviceCheckInterval)
	for {
		select {
		case <-configTimer.C:
			if services, meta, err := srv.consulClient.Catalog().Services(nil); err != nil {
				log.Logger.Error("srv.serviceCheck()", zap.Error(err))
			} else if meta.LastIndex != srv.lastIndex {
				srv.loadServices(services)
			}
		}
	}

}

