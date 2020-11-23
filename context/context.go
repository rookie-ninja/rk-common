// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_ctx

import (
	"fmt"
	rk_entry "github.com/rookie-ninja/rk-common/entry"
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
		application:  "rk-application",
		startTime:    time.Now(),
		loggers:      make(map[string]*LoggerPair, 0),
		viperConfigs: make(map[string]*viper.Viper, 0),
		rkConfigs:    make(map[string]*rk_config.RkConfig, 0),
		eventFactory: rk_query.NewEventFactory(),
		shutdownSig:  make(chan os.Signal),
		customValues: make(map[string]interface{}),
		entries:      make(map[string]rk_entry.Entry),
	}

	entryTypeList = make([]EntryInitializer, 0)
)

type EntryInitializer func(string, *rk_query.EventFactory, *zap.Logger) map[string]rk_entry.Entry

func RegisterEntryInitializer(init EntryInitializer) {
	entryTypeList = append(entryTypeList, init)
}

func ListEntryInitializer() []EntryInitializer {
	// make a copy of it
	res := make([]EntryInitializer, 0)
	for i := range entryTypeList {
		res = append(res, entryTypeList[i])
	}

	return res
}

// init global app context
func init() {
	GlobalAppCtx = &appContext{
		application:  "rk-application",
		startTime:    time.Now(),
		loggers:      make(map[string]*LoggerPair, 0),
		viperConfigs: make(map[string]*viper.Viper, 0),
		rkConfigs:    make(map[string]*rk_config.RkConfig, 0),
		eventFactory: rk_query.NewEventFactory(),
		shutdownSig:  make(chan os.Signal),
		customValues: make(map[string]interface{}),
		entries:      make(map[string]rk_entry.Entry),
	}

	// init logger
	// add stdout logger named with rk
	GlobalAppCtx.AddLoggerPair("default", &LoggerPair{
		Logger: rk_logger.StdoutLogger,
		Config: &rk_logger.StdoutLoggerConfig,
	})

	signal.Notify(GlobalAppCtx.shutdownSig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGQUIT)
}

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
// 6: logger - default logger whose name is "rk"
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
	entries      map[string]rk_entry.Entry
	quitters     map[string]QuitterFunc
}

type QuitterFunc func()

// value related
func (ctx *appContext) AddValue(key string, value interface{}) {
	ctx.customValues[key] = value
}

func (ctx *appContext) GetValue(key string) interface{} {
	return ctx.customValues[key]
}

func (ctx *appContext) ListValues() map[string]interface{} {
	return ctx.customValues
}

// logger related
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

// application related
func (ctx *appContext) SetApplication(application string) {
	ctx.application = application
}

func (ctx *appContext) GetApplication() string {
	return ctx.application
}

// start time related
func (ctx *appContext) SetStartTime(time time.Time) {
	ctx.startTime = time
}

func (ctx *appContext) GetStartTime() time.Time {
	return ctx.startTime
}

func (ctx *appContext) GetUpTime() time.Duration {
	return time.Since(ctx.startTime)
}

// viper config related
func (ctx *appContext) AddViperConfig(name string, vp *viper.Viper) {
	ctx.viperConfigs[name] = vp
}

func (ctx *appContext) GetViperConfig(name string) *viper.Viper {
	return ctx.viperConfigs[name]
}

func (ctx *appContext) ListViperConfigs() map[string]*viper.Viper {
	return ctx.viperConfigs
}

// rk config related
func (ctx *appContext) AddRkConfig(name string, rk *rk_config.RkConfig) {
	ctx.rkConfigs[name] = rk
}

func (ctx *appContext) GetRkConfig(name string) *rk_config.RkConfig {
	return ctx.rkConfigs[name]
}

func (ctx *appContext) ListRkConfigs() map[string]*rk_config.RkConfig {
	return ctx.rkConfigs
}

// event factory related
func (ctx *appContext) SetEventFactory(fac *rk_query.EventFactory) {
	ctx.eventFactory = fac
}

func (ctx *appContext) GetEventFactory() *rk_query.EventFactory {
	return ctx.eventFactory
}

// shutdown signal related
func (ctx *appContext) GetShutdownSig() chan os.Signal {
	return ctx.shutdownSig
}

// quitter related
func (ctx *appContext) AddQuitter(name string, f QuitterFunc) {
	ctx.quitters[name] = f
}

func (ctx *appContext) GetQuitter(name string) QuitterFunc {
	return ctx.quitters[name]
}

func (ctx *appContext) ListQuitters() map[string]QuitterFunc {
	return ctx.quitters
}

// entry related
func (ctx *appContext) AddEntry(name string, entry rk_entry.Entry) {
	ctx.entries[name] = entry
}

func (ctx *appContext) GetEntry(name string) rk_entry.Entry {
	return ctx.entries[name]
}

func (ctx *appContext) MergeEntries(entries map[string]rk_entry.Entry) {
	for k, v := range entries {
		ctx.entries[k] = v
	}
}

func (ctx *appContext) ListEntries() map[string]rk_entry.Entry {
	return ctx.entries
}
