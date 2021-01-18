// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_ctx

import (
	"fmt"
	"github.com/rookie-ninja/rk-common/entry"
	"github.com/rookie-ninja/rk-logger"
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
	GlobalAppCtx = &appContext{}
	// list of entry registration function
	entryRegFuncList = make([]rk_entry.EntryRegFunc, 0)
)

type ShutdownHook func()

// init global app context
func init() {
	GlobalAppCtx = &appContext{
		applicationName: "rk-app",
		startTime:       time.Now(),
		loggers:         make(map[string]*LoggerPair),
		viperConfigs:    make(map[string]*viper.Viper),
		eventFactory:    rk_query.NewEventFactory(),
		shutdownSig:     make(chan os.Signal),
		shutdownHooks:   make(map[string]ShutdownHook),
		customValues:    make(map[string]interface{}),
		entries:         make(map[string]rk_entry.Entry),
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
// 1: applicationName - name of running application
// 2: startTime - application start time
// 3: loggers - loggers with name as a key
// 4: viperConfigs - viper configs with name as a key
// 5: eventFactory - event data factory
// 6: shutdownSig - a channel receiving shutdown signals
// 7: shutdownHooks - a list of shutdownHook function registered by user
// 8: customValues - custom k/v store
type appContext struct {
	applicationName string
	startTime       time.Time
	loggers         map[string]*LoggerPair
	viperConfigs    map[string]*viper.Viper
	eventFactory    *rk_query.EventFactory
	customValues    map[string]interface{}
	entries         map[string]rk_entry.Entry
	shutdownSig     chan os.Signal
	shutdownHooks   map[string]ShutdownHook
}

func RegisterEntry(regFunc rk_entry.EntryRegFunc) {
	if regFunc == nil {
		return
	}
	entryRegFuncList = append(entryRegFuncList, regFunc)
}

func ListEntryRegFunc() []rk_entry.EntryRegFunc {
	// make a copy of it
	res := make([]rk_entry.EntryRegFunc, 0)
	for i := range entryRegFuncList {
		res = append(res, entryRegFuncList[i])
	}

	return res
}

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

func (ctx *appContext) DeleteValue(key string) {
	delete(ctx.customValues, key)
}

func (ctx *appContext) ClearValues() {
	for k := range ctx.customValues {
		delete(ctx.customValues, k)
	}
}

// logger related
func (ctx *appContext) AddLoggerPair(name string, pair *LoggerPair) string {
	if pair == nil {
		return ""
	}

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

func (ctx *appContext) clearLoggerPairs() {
	for k := range ctx.loggers {
		delete(ctx.loggers, k)
	}
}

// application related
func (ctx *appContext) SetApplicationName(name string) {
	ctx.applicationName = name
}

func (ctx *appContext) GetApplicationName() string {
	return ctx.applicationName
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
	if vp == nil || len(name) < 1 {
		return
	}
	ctx.viperConfigs[name] = vp
}

func (ctx *appContext) GetViperConfig(name string) *viper.Viper {
	return ctx.viperConfigs[name]
}

func (ctx *appContext) ListViperConfigs() map[string]*viper.Viper {
	return ctx.viperConfigs
}

func (ctx *appContext) clearViperConfigs() {
	for k := range ctx.viperConfigs {
		delete(ctx.viperConfigs, k)
	}
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

// shutdown hook related
func (ctx *appContext) AddShutdownHook(name string, f ShutdownHook) {
	if f == nil {
		return
	}
	ctx.shutdownHooks[name] = f
}

func (ctx *appContext) GetShutdownHook(name string) ShutdownHook {
	return ctx.shutdownHooks[name]
}

func (ctx *appContext) ListShutdownHooks() map[string]ShutdownHook {
	return ctx.shutdownHooks
}

func (ctx *appContext) clearShutdownHooks() {
	for k := range ctx.shutdownHooks {
		delete(ctx.shutdownHooks, k)
	}
}

// entry related
func (ctx *appContext) AddEntry(name string, entry rk_entry.Entry) {
	if entry == nil {
		return
	}
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

func (ctx *appContext) clearEntries() {
	for k := range ctx.entries {
		delete(ctx.entries, k)
	}
}
