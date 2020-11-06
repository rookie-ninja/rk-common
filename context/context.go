// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_ctx

import (
	"fmt"
	"github.com/rookie-ninja/rk-config"
	"github.com/rookie-ninja/rk-query"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"time"
)

type LoggerPair struct {
	Logger *zap.Logger
	Config *zap.Config
}

type AppContext struct {
	application  string
	startTime    time.Time
	loggers      map[string]*LoggerPair
	eventFactory *rk_query.EventFactory
	viperConfigs map[string]*viper.Viper
	rkConfigs    map[string]*rk_config.RkConfig
	shutdownSig  chan os.Signal
}

type AppContextOption func(*AppContext)

func WithApplication(app string) AppContextOption {
	return func(ctx *AppContext) {
		ctx.application = app
	}
}

func WithStartTime(ts time.Time) AppContextOption {
	return func(ctx *AppContext) {
		ctx.startTime = ts
	}
}

func WithEventFactory(fac *rk_query.EventFactory) AppContextOption {
	return func(ctx *AppContext) {
		ctx.eventFactory = fac
	}
}

func WithShutdownSig(sig chan os.Signal) AppContextOption {
	return func(ctx *AppContext) {
		ctx.shutdownSig = sig
	}
}

func NewAppContext(opts ...AppContextOption) *AppContext {
	ctx := &AppContext{
		application:  "unknown-application",
		startTime:    time.Unix(0, 0),
		loggers:      make(map[string]*LoggerPair, 0),
		eventFactory: rk_query.NewEventFactory(),
		viperConfigs: make(map[string]*viper.Viper, 0),
		rkConfigs:    make(map[string]*rk_config.RkConfig, 0),
		shutdownSig:  make(chan os.Signal),
	}

	for i := range opts {
		opts[i](ctx)
	}

	return ctx
}

func (ctx *AppContext) AddLoggerPair(name string, pair *LoggerPair) string {
	if len(name) < 1 {
		name = fmt.Sprintf("logger-%d", len(ctx.loggers)+1)
	}

	ctx.loggers[name] = pair
	return name
}

func (ctx *AppContext) AddRkConfig(name string, config *rk_config.RkConfig) string {
	if len(name) < 1 {
		name = fmt.Sprintf("rkconfig-%d", len(ctx.rkConfigs)+1)
	}

	ctx.rkConfigs[name] = config
	return name
}

func (ctx *AppContext) AddViperConfig(name string, config *viper.Viper) string {
	if len(name) < 1 {
		name = fmt.Sprintf("viper-%d", len(ctx.viperConfigs)+1)
	}

	ctx.viperConfigs[name] = config
	return name
}

func (ctx *AppContext) GetApplication() string {
	return ctx.application
}

func (ctx *AppContext) GetStartTime() time.Time {
	return ctx.startTime
}

func (ctx *AppContext) GetUpTime() time.Duration {
	return time.Since(ctx.startTime)
}

func (ctx *AppContext) GetRkConfig(name string) *viper.Viper {
	val, ok := ctx.viperConfigs[name]
	if ok {
		return val
	}

	return viper.New()
}

func (ctx *AppContext) ListRkConfigs() map[string]*rk_config.RkConfig {
	return ctx.rkConfigs
}

func (ctx *AppContext) GetViperConfig(name string) *viper.Viper {
	res, _ := ctx.viperConfigs[name]
	return res
}

func (ctx *AppContext) ListViperConfigs() map[string]*viper.Viper {
	return ctx.viperConfigs
}

func (ctx *AppContext) GetDefaultLogger() *zap.Logger {
	return ctx.GetLogger("default")
}

func (ctx *AppContext) GetLogger(name string) *zap.Logger {
	if val, ok := ctx.loggers[name]; ok {
		return val.Logger
	}

	return zap.NewNop()
}

func (ctx *AppContext) GetLoggerConfig(name string) *zap.Config {
	if val, ok := ctx.loggers[name]; ok {
		return val.Config
	}

	return nil
}

func (ctx *AppContext) ListLoggers() map[string]*LoggerPair {
	return ctx.loggers
}

func (ctx *AppContext) GetShutdownSig() chan os.Signal {
	return ctx.shutdownSig
}
