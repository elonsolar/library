package trace

import (
	"context"
	"fmt"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/elonsolar/library/trace/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func TestConfig(t *testing.T) {

	var cfg = &config.Config{}
	_, err := toml.DecodeFile("/Users/chenxiangqian/try/library/trace/application.toml", cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)

}

func TestTraceRegister(t *testing.T) {

	var cfg = &config.Config{}
	_, err := toml.DecodeFile("/Users/chenxiangqian/try/library/trace/application.toml", cfg)
	if err != nil {
		panic(err)
	}
	err = RegisterTrace(cfg)
	if err != nil {
		panic(err)
	}
}

func TestTrace(t *testing.T) {

	var cfg = &config.Config{}
	_, err := toml.DecodeFile("/Users/chenxiangqian/try/library/trace/application.toml", cfg)
	if err != nil {
		panic(err)
	}
	err = RegisterTrace(cfg)
	if err != nil {
		panic(err)
	}

	ctx, span := otel.Tracer("xx").Start(context.Background(), "begin")
	// defer span.End()
	DoFunc(ctx)
	span.End()
}

func DoFunc(ctx context.Context) {
	ctx, span := otel.Tracer("xx").Start(ctx, "DoFunc")
	defer span.End()
	span.SetAttributes(attribute.String("a_key", "a_value"))
}
