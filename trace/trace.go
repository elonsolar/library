package trace

import (
	"fmt"

	"github.com/elonsolar/library/trace/config"
	"github.com/elonsolar/library/trace/exporter"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// register a trace which configed by user  to otel
func RegisterTrace(cfg *config.Config) error {

	exp, err := exporter.NewExporter(cfg.ExportConfig)
	if err != nil {
		return fmt.Errorf("trace init exporter err:%w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(newResource(cfg.ResourceConfig)),
	)
	otel.SetTracerProvider(tp)
	return nil
}
