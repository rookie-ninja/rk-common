package rk_ctx

import (
	"github.com/rookie-ninja/rk-common/entry"
	rk_logger "github.com/rookie-ninja/rk-logger"
	"github.com/rookie-ninja/rk-query"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestGlobalAppCtx_init(t *testing.T) {
	assert.NotNil(t, GlobalAppCtx)

	assert.NotEmpty(t, GlobalAppCtx.applicationName)
	assert.NotNil(t, GlobalAppCtx.startTime)
	assert.NotNil(t, GlobalAppCtx.loggers)
	assert.NotEmpty(t, GlobalAppCtx.loggers)
	assert.NotNil(t, GlobalAppCtx.viperConfigs)
	assert.Empty(t, GlobalAppCtx.viperConfigs)
	assert.NotNil(t, GlobalAppCtx.eventFactory)
	assert.NotNil(t, GlobalAppCtx.shutdownSig)
	assert.NotNil(t, GlobalAppCtx.customValues)
	assert.Empty(t, GlobalAppCtx.customValues)
	assert.NotNil(t, GlobalAppCtx.entries)
	assert.Empty(t, GlobalAppCtx.entries)
}

func TestRegisterEntry_WithNilInput(t *testing.T) {
	RegisterEntry(nil)
	assert.Empty(t, entryRegFuncList)
}

func TestRegisterEntry_HappyCase(t *testing.T) {
	regFunc := func(string, *rk_query.EventFactory, *zap.Logger) map[string]rk_entry.Entry {
		return make(map[string]rk_entry.Entry)
	}

	RegisterEntry(regFunc)
	assert.Len(t, ListEntryRegFunc(), 1)
	// clear reg functions
	entryRegFuncList = entryRegFuncList[:0]
}

func TestListEntryRegFunc_HappyCase(t *testing.T) {
	regFunc := func(string, *rk_query.EventFactory, *zap.Logger) map[string]rk_entry.Entry {
		return make(map[string]rk_entry.Entry)
	}

	RegisterEntry(regFunc)
	assert.Len(t, ListEntryRegFunc(), 1)
	// clear reg functions
	entryRegFuncList = entryRegFuncList[:0]
}

// value related
func TestAppContext_AddValue_WithEmptyKey(t *testing.T) {
	key := ""
	value := "value"
	GlobalAppCtx.AddValue(key, value)
	assert.Equal(t, value, GlobalAppCtx.GetValue(key).(string))
	GlobalAppCtx.ClearValues()
}

func TestAppContext_AddValue_WithEmptyValue(t *testing.T) {
	key := "key"
	value := ""
	GlobalAppCtx.AddValue(key, value)
	assert.Equal(t, value, GlobalAppCtx.GetValue(key).(string))
	GlobalAppCtx.ClearValues()
}

func TestAppContext_AddValue_HappyCase(t *testing.T) {
	key := "key"
	value := "value"
	GlobalAppCtx.AddValue(key, value)
	assert.Equal(t, value, GlobalAppCtx.GetValue(key).(string))
	GlobalAppCtx.ClearValues()
}

func TestAppContext_GetValue_WithEmptyKey(t *testing.T) {
	key := ""
	value := "value"
	GlobalAppCtx.AddValue(key, value)
	assert.Equal(t, value, GlobalAppCtx.GetValue(key).(string))
	GlobalAppCtx.ClearValues()
}

func TestAppContext_GetValue_WithEmptyValue(t *testing.T) {
	key := "key"
	value := ""
	GlobalAppCtx.AddValue(key, value)
	assert.Equal(t, value, GlobalAppCtx.GetValue(key).(string))
	GlobalAppCtx.ClearValues()
}

func TestAppContext_GetValue_HappyCase(t *testing.T) {
	key := "key"
	value := "value"
	GlobalAppCtx.AddValue(key, value)
	assert.Equal(t, value, GlobalAppCtx.GetValue(key).(string))
	GlobalAppCtx.ClearValues()
}

func TestAppContext_ListValue_WithEmptyKey(t *testing.T) {
	key := ""
	value := "value"
	GlobalAppCtx.AddValue(key, value)
	assert.True(t, len(GlobalAppCtx.ListValues()) == 1)
	assert.Equal(t, value, GlobalAppCtx.ListValues()[key])
	GlobalAppCtx.ClearValues()
}

func TestAppContext_ListValue_WithEmptyValue(t *testing.T) {
	key := "key"
	value := ""
	GlobalAppCtx.AddValue(key, value)
	assert.True(t, len(GlobalAppCtx.ListValues()) == 1)
	assert.Equal(t, value, GlobalAppCtx.ListValues()[key])
	GlobalAppCtx.ClearValues()
}

func TestAppContext_ListValue_HappyCase(t *testing.T) {
	key := "key"
	value := "value"
	GlobalAppCtx.AddValue(key, value)
	assert.True(t, len(GlobalAppCtx.ListValues()) == 1)
	assert.Equal(t, value, GlobalAppCtx.ListValues()[key])
	GlobalAppCtx.ClearValues()
}

// logger related
func TestAppContext_AddLoggerPair_WithEmptyName(t *testing.T) {
	pair := &LoggerPair{}
	assert.Equal(t, "logger-2", GlobalAppCtx.AddLoggerPair("", pair))
	assert.Equal(t, pair, GlobalAppCtx.GetLoggerPair("logger-2"))

	// clear logger pairs
	GlobalAppCtx.clearLoggerPairs()
}

func TestAppContext_AddLoggerPair_WithNilPair(t *testing.T) {
	assert.Equal(t, "", GlobalAppCtx.AddLoggerPair("logger", nil))
	assert.Nil(t, GlobalAppCtx.GetLoggerPair("logger"))

	// clear logger pairs
	GlobalAppCtx.clearLoggerPairs()
}

func TestAppContext_AddLoggerPair_HappyCase(t *testing.T) {
	pair := &LoggerPair{}
	loggerName := "logger-unit-test"
	assert.Equal(t, loggerName, GlobalAppCtx.AddLoggerPair(loggerName, pair))
	assert.Equal(t, pair, GlobalAppCtx.GetLoggerPair(loggerName))

	// clear logger pairs
	GlobalAppCtx.clearLoggerPairs()
}

func TestAppContext_GetLoggerPair_WithEmptyName(t *testing.T) {
	loggerName := ""
	assert.Nil(t, GlobalAppCtx.GetLoggerPair(loggerName))

	// clear logger pairs
	GlobalAppCtx.clearLoggerPairs()
}

func TestAppContext_GetLoggerPair_WithNonExist(t *testing.T) {
	loggerName := "non-exist"
	assert.Nil(t, GlobalAppCtx.GetLoggerPair(loggerName))

	// clear logger pairs
	GlobalAppCtx.clearLoggerPairs()
}

func TestAppContext_GetLoggerPair_HappyCase(t *testing.T) {
	pair := &LoggerPair{}
	loggerName := "logger-unit-test"
	assert.Equal(t, loggerName, GlobalAppCtx.AddLoggerPair(loggerName, pair))
	assert.Equal(t, pair, GlobalAppCtx.GetLoggerPair(loggerName))

	// clear logger pairs
	GlobalAppCtx.clearLoggerPairs()
}

func TestAppContext_ListLoggerPairs_HappyCase(t *testing.T) {
	pair := &LoggerPair{}
	loggerName := "logger-unit-test"
	assert.Equal(t, loggerName, GlobalAppCtx.AddLoggerPair(loggerName, pair))
	assert.Equal(t, 1, len(GlobalAppCtx.ListLoggerPairs()))

	// clear logger pairs
	GlobalAppCtx.clearLoggerPairs()
}

func TestAppContext_GetLogger_WithNonExist(t *testing.T) {
	loggerName := "non-exist"
	assert.NotNil(t, GlobalAppCtx.GetLogger(loggerName))
	// clear logger pairs
	GlobalAppCtx.clearLoggerPairs()
}

func TestAppContext_GetLogger_HappyCase(t *testing.T) {
	loggerName := "logger-unit-test"
	pair := &LoggerPair{
		Logger: rk_logger.StdoutLogger,
		Config: &rk_logger.StdoutLoggerConfig,
	}
	GlobalAppCtx.AddLoggerPair(loggerName, pair)
	assert.Equal(t, pair.Logger, GlobalAppCtx.GetLogger(loggerName))

	// clear logger pairs
	GlobalAppCtx.clearLoggerPairs()
}

func TestAppContext_GetDefaultLogger_HappyCase(t *testing.T) {
	assert.NotNil(t, GlobalAppCtx.GetDefaultLogger())

}

func TestAppContext_GetLoggerConfig_WithNonExist(t *testing.T) {
	loggerName := "non-exist"
	assert.Nil(t, GlobalAppCtx.GetLoggerConfig(loggerName))
	// clear logger pairs
	GlobalAppCtx.clearLoggerPairs()
}

func TestAppContext_GetLoggerConfig_HappyCase(t *testing.T) {
	loggerName := "logger-unit-test"
	pair := &LoggerPair{
		Logger: rk_logger.StdoutLogger,
		Config: &rk_logger.StdoutLoggerConfig,
	}
	GlobalAppCtx.AddLoggerPair(loggerName, pair)
	assert.Equal(t, pair.Config, GlobalAppCtx.GetLoggerConfig(loggerName))

	// clear logger pairs
	GlobalAppCtx.clearLoggerPairs()
}

// application related
func TestAppContext_SetApplicationName_HappyCase(t *testing.T) {
	GlobalAppCtx.SetApplicationName("unit-test")
	assert.Equal(t, "unit-test", GlobalAppCtx.GetApplicationName())
}

func TestAppContext_GetApplicationName_HappyCase(t *testing.T) {
	GlobalAppCtx.SetApplicationName("unit-test")
	assert.Equal(t, "unit-test", GlobalAppCtx.GetApplicationName())
}

// start time related
func TestAppContext_SetStartTime_HappyCase(t *testing.T) {
	now := time.Now()
	GlobalAppCtx.SetStartTime(now)
	assert.Equal(t, now, GlobalAppCtx.GetStartTime())
}

func TestAppContext_GetStartTime_HappyCase(t *testing.T) {
	now := time.Now()
	GlobalAppCtx.SetStartTime(now)
	assert.Equal(t, now, GlobalAppCtx.GetStartTime())
}

func TestAppContext_GetUpTime_HappyCase(t *testing.T) {
	now := time.Now()
	GlobalAppCtx.SetStartTime(now)
	time.Sleep(1 * time.Second)
	assert.True(t, GlobalAppCtx.GetUpTime() >= 1*time.Second)
}

// viper config related
func TestAppContext_AddViperConfig_WithEmptyName(t *testing.T) {
	name := ""
	vp := viper.New()
	GlobalAppCtx.AddViperConfig(name, vp)
	assert.True(t, len(GlobalAppCtx.viperConfigs) == 0)
}

func TestAppContext_AddViperConfig_WithNilViper(t *testing.T) {
	name := "viper-config"
	GlobalAppCtx.AddViperConfig(name, nil)
	assert.True(t, len(GlobalAppCtx.viperConfigs) == 0)
	// clear viper config
	GlobalAppCtx.clearViperConfigs()
}

func TestAppContext_AddViperConfig_HappyCase(t *testing.T) {
	name := "viper-config"
	vp := viper.New()
	GlobalAppCtx.AddViperConfig(name, vp)
	assert.True(t, len(GlobalAppCtx.viperConfigs) == 1)
	assert.Equal(t, vp, GlobalAppCtx.GetViperConfig(name))
	// clear viper config
	GlobalAppCtx.clearViperConfigs()
}

func TestAppContext_GetViperConfig_WithNonExist(t *testing.T) {
	name := "non-exist"
	assert.Nil(t, GlobalAppCtx.GetViperConfig(name))
	// clear viper config
	GlobalAppCtx.clearViperConfigs()
}

func TestAppContext_GetViperConfig_HappyCase(t *testing.T) {
	name := "viper-config"
	vp := viper.New()
	GlobalAppCtx.AddViperConfig(name, vp)
	assert.True(t, len(GlobalAppCtx.viperConfigs) == 1)
	assert.NotNil(t, GlobalAppCtx.GetViperConfig(name))
	// clear viper config
	GlobalAppCtx.clearViperConfigs()
}

func TestAppContext_ListViperConfigs_WithEmptyList(t *testing.T) {
	assert.True(t, len(GlobalAppCtx.ListViperConfigs()) == 0)
	// clear viper config
	GlobalAppCtx.clearViperConfigs()
}

func TestAppContext_ListViperConfigs_HappyCase(t *testing.T) {
	name := "viper-config"
	vp := viper.New()
	GlobalAppCtx.AddViperConfig(name, vp)
	assert.True(t, len(GlobalAppCtx.ListViperConfigs()) == 1)
	// clear viper config
	GlobalAppCtx.clearViperConfigs()
}

// event factory related
func TestAppContext_SetEventFactory_WithNilFactory(t *testing.T) {
	GlobalAppCtx.SetEventFactory(nil)
	assert.Nil(t, GlobalAppCtx.GetEventFactory())
}

func TestAppContext_SetEventFactory_HappyCase(t *testing.T) {
	fac := rk_query.NewEventFactory()
	GlobalAppCtx.SetEventFactory(fac)
	assert.Equal(t, fac, GlobalAppCtx.GetEventFactory())
}

func TestAppContext_GetEventFactory_HappyCase(t *testing.T) {
	fac := rk_query.NewEventFactory()
	GlobalAppCtx.SetEventFactory(fac)
	assert.Equal(t, fac, GlobalAppCtx.GetEventFactory())
}

// shutdown signal related
func TestAppContext_GetShutdownSig_HappyCase(t *testing.T) {
	assert.NotNil(t, GlobalAppCtx.GetShutdownSig())
}

// shutdown hook related
func TestAppContext_AddShutdownHook_WithEmptyName(t *testing.T) {
	name := ""
	f := func() {}
	GlobalAppCtx.AddShutdownHook(name, f)
	assert.Equal(t, 1, len(GlobalAppCtx.ListShutdownHooks()))
	assert.NotNil(t, GlobalAppCtx.GetShutdownHook(name))
	// clear shutdown hooks
	GlobalAppCtx.clearShutdownHooks()
}

func TestAppContext_AddShutdownHook_WithNilFunc(t *testing.T) {
	name := ""
	GlobalAppCtx.AddShutdownHook(name, nil)
	assert.Equal(t, 0, len(GlobalAppCtx.ListShutdownHooks()))
	assert.Nil(t, GlobalAppCtx.GetShutdownHook(name))
	// clear shutdown hooks
	GlobalAppCtx.clearShutdownHooks()
}

func TestAppContext_AddShutdownHook_HappyCase(t *testing.T) {
	name := "unit-test-hook"
	f := func() {}
	GlobalAppCtx.AddShutdownHook(name, f)
	assert.Equal(t, 1, len(GlobalAppCtx.ListShutdownHooks()))
	assert.NotNil(t, GlobalAppCtx.GetShutdownHook(name))
	// clear shutdown hooks
	GlobalAppCtx.clearShutdownHooks()
}

func TestAppContext_GetShutdownHook_WithNonExistHooks(t *testing.T) {
	name := "non-exist"
	assert.Nil(t, GlobalAppCtx.GetShutdownHook(name))
	// clear shutdown hooks
	GlobalAppCtx.clearShutdownHooks()
}

func TestAppContext_GetShutdownHook_HappyCase(t *testing.T) {
	name := "unit-test-hook"
	f := func() {}
	GlobalAppCtx.AddShutdownHook(name, f)
	assert.NotNil(t, GlobalAppCtx.GetShutdownHook(name))
	// clear shutdown hooks
	GlobalAppCtx.clearShutdownHooks()
}

func TestAppContext_ListShutdownHooks_WithEmptyHooks(t *testing.T) {
	assert.Equal(t, 0, len(GlobalAppCtx.ListShutdownHooks()))
	// clear shutdown hooks
	GlobalAppCtx.clearShutdownHooks()
}

func TestAppContext_ListShutdownHooks_HappyCase(t *testing.T) {
	name := "unit-test-hook"
	f := func() {}
	GlobalAppCtx.AddShutdownHook(name, f)
	assert.Equal(t, 1, len(GlobalAppCtx.ListShutdownHooks()))
	// clear shutdown hooks
	GlobalAppCtx.clearShutdownHooks()
}

// entry related
func TestAppContext_AddEntry_WithEmptyName(t *testing.T) {
	name := ""
	entry := &EntryMock{}
	GlobalAppCtx.AddEntry(name, entry)
	assert.Equal(t, 1, len(GlobalAppCtx.ListEntries()))
	assert.Equal(t, entry, GlobalAppCtx.GetEntry(name))

	// clear entries
	GlobalAppCtx.clearEntries()
}

func TestAppContext_AddEntry_WithNilEntry(t *testing.T) {
	name := ""
	GlobalAppCtx.AddEntry(name, nil)
	assert.Equal(t, 0, len(GlobalAppCtx.ListEntries()))
	assert.Nil(t, GlobalAppCtx.GetEntry(name))

	// clear entries
	GlobalAppCtx.clearEntries()
}

func TestAppContext_AddEntry_HappyCase(t *testing.T) {
	name := "unit-test-entry"
	entry := &EntryMock{}
	GlobalAppCtx.AddEntry(name, entry)
	assert.Equal(t, 1, len(GlobalAppCtx.ListEntries()))
	assert.Equal(t, entry, GlobalAppCtx.GetEntry(name))

	// clear entries
	GlobalAppCtx.clearEntries()
}

func TestAppContext_GetEntry_WithNonExistEntry(t *testing.T) {
	name := "non-exist"
	assert.Nil(t, GlobalAppCtx.GetEntry(name))

	// clear entries
	GlobalAppCtx.clearEntries()
}

func TestAppContext_GetEntry_HappyCase(t *testing.T) {
	name := "unit-test-entry"
	entry := &EntryMock{}
	GlobalAppCtx.AddEntry(name, entry)
	assert.Equal(t, entry, GlobalAppCtx.GetEntry(name))

	// clear entries
	GlobalAppCtx.clearEntries()
}

func TestAppContext_ListEntries_WithEmptyEntry(t *testing.T) {
	assert.Equal(t, 0, len(GlobalAppCtx.ListEntries()))

	// clear entries
	GlobalAppCtx.clearEntries()
}

func TestAppContext_ListEntries_HappyCase(t *testing.T) {
	name := "unit-test-entry"
	entry := &EntryMock{}
	GlobalAppCtx.AddEntry(name, entry)
	assert.Equal(t, 1, len(GlobalAppCtx.ListEntries()))
	// clear entries
	GlobalAppCtx.clearEntries()
}

func TestAppContext_MergeEntries_HappyCase(t *testing.T) {
	target := map[string]rk_entry.Entry{
		"target-entry": &EntryMock{},
	}

	name := "unit-test-entry"
	entry := &EntryMock{}

	GlobalAppCtx.AddEntry(name, entry)
	GlobalAppCtx.MergeEntries(target)
	assert.Equal(t, 2, len(GlobalAppCtx.ListEntries()))
	// clear entries
	GlobalAppCtx.clearEntries()
}

type EntryMock struct{}

// Bootstrap will be called from boot.Bootstrap()
func (entry *EntryMock) Bootstrap(rk_query.Event) {}

// Important:
// WaitForShutdownSig won't be called from boot.Shutdown
// User could call this function while start entry only without bootstrapper
// We recommend to call Shutdown in this function
func (entry *EntryMock) WaitForShutdownSig(time.Duration) {}

// Shutdown will be called from boot.Shutdown()
func (entry *EntryMock) Shutdown(rk_query.Event) {}

// Return name of entry
func (entry *EntryMock) GetName() string {
	return ""
}

// Return type of entry
func (entry *EntryMock) GetType() string {
	return ""
}

// Modify as needed
func (entry *EntryMock) String() string {
	return ""
}
