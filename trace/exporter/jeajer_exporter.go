package exporter

import (
	"github.com/elonsolar/library/trace/config"
	"go.opentelemetry.io/otel/exporters/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func newJeajerExporter(attr config.ExportAttribute) (sdktrace.SpanExporter, error) {

	var cfg = &config.JeagerExporterConfig{}

	if err := attr.Decode(cfg); err != nil {
		return nil, err
	}

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Url)))
	if err != nil {
		return nil, err
	}
	return exp, nil
}
