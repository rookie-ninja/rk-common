// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_ctx

import (
	"fmt"
	rk_config "github.com/rookie-ninja/rk-config"
	rk_logger "github.com/rookie-ninja/rk-logger"
	"github.com/rookie-ninja/rk-query"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	// Global application context
	GlobalAppCtx = &appContext{
		application:  "unknown-application",
		startTime:    time.Unix(0, 0),
		loggers:      make(map[string]*LoggerPair, 0),
		viperConfigs: make(map[string]*viper.Viper, 0),
		rkConfigs:    make(map[string]*rk_config.RkConfig, 0),
		eventFactory: rk_query.NewEventFactory(),
		shutdownSig:  make(chan os.Signal),
		customValues: make(map[string]interface{}),
	}
)

// Contains zap.logger and zap.config
// Make sure zap.logger was built from zap.config in order to
// dynamically change zap.logger attributes
type LoggerPair struct {
	Logger *zap.Logger
	Config *zap.Config
}

// Standard application context which contains bellow
// 1: application - name of running application
// 2: startTime - application start time
// 3: loggers - loggers with name as a key
// 4: viperConfigs - viper configs with name as a key
// 5: rkConfigs - rk style configs with name as a key
// 6: logger - default logger whose name is "default"
// 7: eventFactory - event data factory
// 8: shutdownSig - a channel receiving shutdown signals
// 9: customValues - custom k/v store
type appContext struct {
	application  string
	startTime    time.Time
	loggers      map[string]*LoggerPair
	viperConfigs map[string]*viper.Viper
	rkConfigs    map[string]*rk_config.RkConfig
	logger       *zap.Logger
	eventFactory *rk_query.EventFactory
	shutdownSig  chan os.Signal
	customValues map[string]interface{}
}

type appContextOption func(*appContext)

func WithApplication(app string) appContextOption {
	return func(ctx *appContext) {
		ctx.application = app
	}
}

func WithStartTime(ts time.Time) appContextOption {
	return func(ctx *appContext) {
		ctx.startTime = ts
	}
}

func WithLoggers(pairs map[string]*LoggerPair) appContextOption {
	return func(ctx *appContext) {
		ctx.loggers = pairs
	}
}

func WithViperConfigs(vipers map[string]*viper.Viper) appContextOption {
	return func(ctx *appContext) {
		ctx.viperConfigs = vipers
	}
}

func WithRkConfigs(rks map[string]*rk_config.RkConfig) appContextOption {
	return func(ctx *appContext) {
		ctx.rkConfigs = rks
	}
}

func WithEventFactory(fac *rk_query.EventFactory) appContextOption {
	return func(ctx *appContext) {
		ctx.eventFactory = fac
	}
}

func NewAppContext(opts ...appContextOption) *appContext {
	ctx := &appContext{
		application:  "unknown-application",
		startTime:    time.Unix(0, 0),
		loggers:      make(map[string]*LoggerPair, 0),
		viperConfigs: make(map[string]*viper.Viper, 0),
		rkConfigs:    make(map[string]*rk_config.RkConfig, 0),
		eventFactory: rk_query.NewEventFactory(),
		shutdownSig:  make(chan os.Signal),
		customValues: make(map[string]interface{}),
	}

	for i := range opts {
		opts[i](ctx)
	}

	// assign default logger
	if pair, ok := ctx.loggers["default"]; ok {
		ctx.logger = pair.Logger
	} else {
		ctx.loggers["default"] = &LoggerPair{
			Logger: rk_logger.StdoutLogger,
			Config: &rk_logger.StdoutLoggerConfig,
		}
		ctx.logger = rk_logger.StdoutLogger
	}

	// register signal
	signal.Notify(ctx.shutdownSig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	GlobalAppCtx = ctx

	return ctx
}

func (ctx *appContext) GetValue(key string) interface{} {
	return ctx.customValues[key]
}

func (ctx *appContext) ListValues() map[string]interface{} {
	return ctx.customValues
}

func (ctx *appContext) AddValue(key string, value interface{}) {
	ctx.customValues[key] = value
}

func (ctx *appContext) AddLoggerPair(name string, pair *LoggerPair) string {
	if len(name) < 1 {
		name = fmt.Sprintf("logger-%d", len(ctx.loggers)+1)
	}

	ctx.loggers[name] = pair
	return name
}

func (ctx *appContext) GetLoggerPair(name string) *LoggerPair {
	return ctx.loggers[name]
}

func (ctx *appContext) ListLoggerPairs() map[string]*LoggerPair {
	return ctx.loggers
}

func (ctx *appContext) GetLogger(name string) *zap.Logger {
	if val, ok := ctx.loggers[name]; ok {
		return val.Logger
	}

	return zap.NewNop()
}

func (ctx *appContext) GetDefaultLogger() *zap.Logger {
	return ctx.GetLogger("default")
}

func (ctx *appContext) GetLoggerConfig(name string) *zap.Config {
	if val, ok := ctx.loggers[name]; ok {
		return val.Config
	}
	
	return nil
}

func (ctx *appContext) GetApplication() string {
	return ctx.application
}

func (ctx *appContext) GetStartTime() time.Time {
	return ctx.startTime
}

func (ctx *appContext) GetUpTime() time.Duration {
	return time.Since(ctx.startTime)
}

func (ctx *appContext) GetViperConfig(name string) *viper.Viper {
	return ctx.viperConfigs[name]
}

func (ctx *appContext) ListViperConfigs() map[string]*viper.Viper {
	return ctx.viperConfigs
}

func (ctx *appContext) GetRkConfig(name string) *rk_config.RkConfig {
	return ctx.rkConfigs[name]
}

func (ctx *appContext) ListRkConfigs() map[string]*rk_config.RkConfig {
	return ctx.rkConfigs
}

func (ctx *appContext) GetShutdownSig() chan os.Signal {
	return ctx.shutdownSig
}
