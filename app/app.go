package app

import (
	"net/http"

	trace "github.com/elonsolar/library/trace"
	"github.com/fvbock/endless"
)

type Config struct {
	appName string
	handler http.Handler
}

type App struct {
	cfg Config
}

func NewApp(options ...Option) *App {

	var defalutCfg = &Config{
		appName: "default_app",
	}
	for _, op := range options {
		op.apply(defalutCfg)
	}

	return &App{
		cfg: *defalutCfg,
	}
}

func (app *App) Run() error {
	go func() {
		_, err := trace.RegisterTrace(nil)
		if err != nil {
			panic(err)
		}

	}()

	go func() {
		err := endless.ListenAndServe("localhost:4244", app.cfg.handler) // 底层 调用  NewServer().Server()
		if err != nil {
			panic(err)
		}
	}()
	return nil
}

type Option interface {
	apply(cfg *Config)
}

type appOptionFunc func(cfg *Config)

func (fn appOptionFunc) apply(cfg *Config) {
	fn(cfg)
}

// WithHandler with specific handler
func WithHandler(handler http.Handler) Option {
	return appOptionFunc(func(cfg *Config) {
		cfg.handler = handler
	})
}
