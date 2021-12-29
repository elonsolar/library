package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"testing"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func NewExp() sdktrace.SpanExporter {
	l := log.New(os.Stdout, "", 0)
	// Write telemetry data to a file.
	f, err := os.Create("traces.txt")
	// f, err := os.OpenFile("traces.txt", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		l.Fatal(err)
	}

	// 比较简单 new( +options)
	exp, err := newExporter(f)
	if err != nil {
		l.Fatal(err)
	}
	return exp
}

func TestMain(t *testing.T) {
	l := log.New(os.Stdout, "", 0)

	// Write telemetry data to a file.
	//f, err := os.Create("traces.txt")
	// f, err := os.OpenFile("traces.txt", os.O_CREATE|os.O_RDWR, 0666)
	// if err != nil {
	// 	l.Fatal(err)
	// }
	// defer f.Close()

	// 比较简单 new( +options)
	// exp, err := newExporter(f)
	// if err != nil {
	// 	l.Fatal(err)
	// }

	//
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(NewExp()), //先创建span process  + 回上层创建provider
		sdktrace.WithResource(newResource()),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			l.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tp)

	// l := log.New(os.Stdout, "", 0)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	errCh := make(chan error)
	app := NewApp(os.Stdin, l)
	go func() {
		errCh <- app.Run(context.Background())
	}()

	select {
	case <-sigCh:
		l.Println("\ngoodbye")
		return
	case err := <-errCh:
		if err != nil {
			l.Fatal(err)
		}
	}
}
