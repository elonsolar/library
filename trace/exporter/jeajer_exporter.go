package exporter

import (
	"github.com/elonsolar/library/trace/config"
	"go.opentelemetry.io/otel/exporters/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// https://github.com/open-telemetry/opentelemetry-go/tree/main/exporters/jaeger
func newJeajerExporter(attr config.ExportAttribute) (func() error, sdktrace.SpanExporter, error) {

	var cfg = &config.JeagerExporterConfig{}

	if err := attr.Decode(cfg); err != nil {
		return nil, nil, err
	}

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Url)))
	if err != nil {
		return nil, nil, err
	}
	return func() error { return nil }, exp, nil
}
