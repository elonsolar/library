package trace

import (
	"fmt"

	"github.com/elonsolar/library/trace/config"
	"github.com/elonsolar/library/trace/exporter"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// register a trace which configed by user  to otel
func RegisterTrace(cfg *config.Config) (func() error, error) {

	release, exp, err := exporter.NewExporter(cfg.ExportConfig)
	if err != nil {
		return release, fmt.Errorf("trace init exporter err:%w  ", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(newResource(cfg.ResourceConfig)),
	)
	otel.SetTracerProvider(tp)
	return release, nil
}
