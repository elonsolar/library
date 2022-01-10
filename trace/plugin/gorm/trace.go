package gorm_trace

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/gorm"
)

const (
	callBackBeforeName = "opentracing:before"
	callBackAfterName  = "opentracing:after"
)

type TracePlugin struct{}

func (op *TracePlugin) Name() string {
	return "tracePlugin"
}

func (op *TracePlugin) Initialize(db *gorm.DB) (err error) {
	// 开始前 - 并不是都用相同的方法，可以自己自定义
	db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before)
	db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before)
	db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, before)
	db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before)
	db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before)

	// 结束后 - 并不是都用相同的方法，可以自己自定义
	db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after("CREATE"))
	db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after("QUERY"))
	db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after("DELETE"))
	db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after("UPDATE"))
	db.Callback().Row().After("gorm:row").Register(callBackAfterName, after("ROW"))
	db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after("RAW"))
	return
}

var _ gorm.Plugin = &TracePlugin{}

// 包内静态变量
const gormSpanKey = "gorm"

func before(db *gorm.DB) {

	// db.Statement.Context, _ = otel.Tracer(gormSpanKey).Start(db.Statement.Context, "gorm_before")

}
func after(name string) func(db *gorm.DB) {
	return func(db *gorm.DB) {

		_, span := otel.Tracer(gormSpanKey).Start(db.Statement.Context, name)

		defer span.End()
		// sql --> 写法来源GORM V2的日志
		span.SetAttributes(attribute.String("sql", db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)))
	}
}
