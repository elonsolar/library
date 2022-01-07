package trace

import (
	"log"

	"github.com/elonsolar/library/trace/config"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// newResource returns a resource describing this application.
func newResource(cfg *config.ResourceConfig) *resource.Resource {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
		),
	)
	if err != nil {
		log.Fatalf("resource merge err%v", err)
	}
	return r
}
