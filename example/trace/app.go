package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// newExporter returns a console exporter.
func newExporter(w io.Writer) (sdktrace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

// newResource returns a resource describing this application.
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("fib"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}

// App is a Fibonacci computation application.
type App struct {
	r io.Reader
	l *log.Logger
}

// NewApp returns a new App.
func NewApp(r io.Reader, l *log.Logger) *App {
	return &App{r: r, l: l}
}

// Run starts polling users for Fibonacci number requests and writes results.
func (a *App) Run(ctx context.Context) error {
	for {
		var span trace.Span

		ctx, span = otel.Tracer("xx").Start(ctx, "Run")
		// n, err := a.Poll(ctx)
		// if err != nil {
		// 	return err
		// }

		span.SetAttributes(attribute.String("request.n", "hahah"))
		a.Write(ctx, uint(12))
		span.End()
		time.Sleep(time.Second * 3)
	}
	// return nil
}

// Poll asks a user for input and returns the request.
func (a *App) Poll(ctx context.Context) (uint, error) {

	_, span := otel.Tracer("xx").Start(ctx, "Poll")
	defer span.End()
	a.l.Print("What Fibonacci number would you like to know: ")

	var n uint
	_, err := fmt.Fscanf(a.r, "%d\n", &n)
	nstr := strconv.FormatUint(uint64(n), 10)
	span.SetAttributes(attribute.String("request.n", nstr))
	return n, err
}

// Write writes the n-th Fibonacci number back to the user.
func (a *App) Write(ctx context.Context, n uint) {
	var span trace.Span
	ctx, span = otel.Tracer("xx").Start(ctx, "Write")
	defer span.End()
	f, err := func(ctx context.Context) (uint64, error) {
		_, span := otel.Tracer("xx").Start(ctx, "Fibonacci")
		defer span.End()
		return Fibonacci(n)
	}(ctx)

	if err != nil {
		a.l.Printf("Fibonacci(%d): %v\n", n, err)
	} else {
		a.l.Printf("Fibonacci(%d) = %d\n", n, f)
	}
}

// Fibonacci returns the n-th fibonacci number.
func Fibonacci(n uint) (uint64, error) {
	if n <= 1 {
		return uint64(n), nil
	}

	var n2, n1 uint64 = 0, 1
	for i := uint(2); i < n; i++ {
		n2, n1 = n1, n1+n2
	}

	return n2 + n1, nil
}
