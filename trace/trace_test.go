package trace

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/elonsolar/library/trace/config"
	gt "github.com/elonsolar/library/trace/plugin/gorm_trace"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestConfig(t *testing.T) {

	var cfg = &config.Config{}
	_, err := toml.DecodeFile("/Users/chenxiangqian/try/library/trace/application.toml", cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)

}

func TestTraceRegister(t *testing.T) {

	var cfg = &config.Config{}
	_, err := toml.DecodeFile("/Users/chenxiangqian/try/library/trace/application.toml", cfg)
	if err != nil {
		panic(err)
	}
	_, err = RegisterTrace(cfg)
	if err != nil {
		panic(err)
	}
}

func TestTrace(t *testing.T) {

	var cfg = &config.Config{}
	_, err := toml.DecodeFile("/Users/chenxiangqian/try/library/trace/application.toml", cfg)
	if err != nil {
		panic(err)
	}
	release, err := RegisterTrace(cfg)
	defer release()
	if err != nil {
		panic(err)
	}
	eg := gin.New()

	// eg.Use(gin)
	eg.Run(":9090")
	ctx, span := otel.Tracer("xx").Start(context.Background(), "begin")
	// defer span.End()
	// for i := 0; i < 100; i++ {

	DoFunc(ctx)
	// }
	span.End()

	time.Sleep(time.Second * 20)
}

func DoFunc(ctx context.Context) {
	ctx, span := otel.Tracer("xx").Start(ctx, "DoFunc")
	db := NewOrm()

	// db.Statement.Context = ctx
	var data = make([]map[string]interface{}, 0)
	if err := db.WithContext(ctx).Table("USER").Scan(data).Error; err != nil {
		panic(err)
	}
	defer span.End()
	// span.SetAttributes(attribute.String("a_key", "a_value"))
}

func NewOrm() *gorm.DB {
	db, err := gorm.Open(mysqlConfig(DefaultConfig), gormConfig(DefaultConfig))
	if err != nil {
		panic(err)
	}
	// demo.FuckSomething()
	_ = db.Use(&gt.TracePlugin{})
	return db
}

func mysqlConfig(config *Config) mysql.Dialector {

	return mysql.Dialector{Config: &mysql.Config{
		DSN:                       mysqlDsn(config), // DSN data source name
		DefaultStringSize:         191,              // string 类型字段的默认长度
		DisableDatetimePrecision:  true,             // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,             // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,             // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,            // 根据版本自动配置
	}}
}

func mysqlDsn(c *Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&&loc=Local", c.UserName, c.Password, c.Host, c.Port, c.DatabaseName)
}

func gormConfig(config *Config) *gorm.Config {

	return &gorm.Config{
		SkipDefaultTransaction:                   true,
		NamingStrategy:                           nil,
		FullSaveAssociations:                     false,
		Logger:                                   logger.Default.LogMode(logger.Info),
		NowFunc:                                  nil,
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		AllowGlobalUpdate:                        false,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  nil,
	}
}

type Config struct {
	UserName     string
	Password     string
	Host         string
	Port         int
	DatabaseName string
}

var DefaultConfig = &Config{
	UserName:     "root",
	Password:     "Maneng1234!@#$",
	Host:         "ntyj.codenai.com",
	Port:         3306,
	DatabaseName: "codenai",
}
