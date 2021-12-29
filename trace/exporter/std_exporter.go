package exporter

import (
	"io"
	"os"

	"github.com/elonsolar/library/trace/config"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// stdoutExporter
func newStdExporter(attr config.ExportAttribute) (sdktrace.SpanExporter, error) {

	var cfg = &config.StdoutExporterConfig{}
	if err := attr.Decode(cfg); err != nil {
		return nil, err
	}

	var (
		writer io.Writer
		err    error
	)
	if cfg.FileName == "" {
		writer = os.Stdout
	} else {
		writer, err = os.Create(cfg.FileName)
		// writer, err = os.OpenFile(cfg.FileName, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return nil, err
		}
	}
	var options = []stdouttrace.Option{stdouttrace.WithWriter(writer)}

	if cfg.PrettyPrint {
		options = append(options, stdouttrace.WithPrettyPrint())
	}

	if !cfg.Timestamps {
		options = append(options, stdouttrace.WithoutTimestamps())
	}

	return stdouttrace.New(options...)
}
