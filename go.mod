module github.com/elonsolar/library

go 1.13

require (
	github.com/BurntSushi/toml v0.4.1
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/mitchellh/mapstructure v1.4.3
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.opentelemetry.io/otel v1.3.0
	go.opentelemetry.io/otel/exporters/jaeger v1.3.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.2.0
	go.opentelemetry.io/otel/sdk v1.3.0
	go.opentelemetry.io/otel/trace v1.3.0
	go.uber.org/zap v1.19.1
)
