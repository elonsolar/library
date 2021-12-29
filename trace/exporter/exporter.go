package exporter

import (
	"errors"

	"github.com/elonsolar/library/trace/config"
	"github.com/elonsolar/library/trace/constant"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	ErrNoExportSpecified = errors.New("没有指定trace exporter")
)

// 我们从配置文件 构建对象，
func NewExporter(cfg *config.ExportConfig) (sdktrace.SpanExporter, error) {

	switch cfg.Name {
	case constant.STDOUT:
		return newStdExporter(cfg.Attribute)
	case constant.JEAJER:
		return newJeajerExporter(cfg.Attribute)
	default:
		return nil, ErrNoExportSpecified
	}
}
