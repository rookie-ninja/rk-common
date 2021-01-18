package rk_common_example

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk-common/context"
	"github.com/rookie-ninja/rk-common/entry"
	"github.com/rookie-ninja/rk-query"
	"go.uber.org/zap"
	"time"
)

// Register entry, must be in init() function since we need to register entry
// at beginning
func init() {
	rk_ctx.RegisterEntry(NewEntry)
}

// A struct which is for unmarshalled YAML
type bootConfig struct {
	Example struct {
		Enabled bool   `yaml:"enabled"`
		Key     string `yaml:"key"`
	} `yaml:"example"`
}

// An implementation of:
// type EntryRegFunc func(string, *rk_query.EventFactory, *zap.Logger) map[string]Entry
func NewEntry(configFilePath string, factory *rk_query.EventFactory, logger *zap.Logger) map[string]rk_entry.Entry {
	// 1: read config file
	bytes := rk_common.MustReadFile(configFilePath)
	config := &bootConfig{}
	// 2: unmarshal config to struct
	if err := yaml.Unmarshal(bytes, config); err != nil {
		rk_common.ShutdownWithError(err)
	}

	res := make(map[string]rk_entry.Entry)

	// 3: construct entry
	if config.Example.Enabled {
		entry := &EntryImpl{
			name:    "example",
			key:     config.Example.Key,
			logger:  logger,
			factory: factory,
		}
		res[entry.GetName()] = entry
	}

	return res
}

type EntryImpl struct {
	name    string
	key     string
	logger  *zap.Logger
	factory *rk_query.EventFactory
}

// Bootstrap will be called from boot.Bootstrap()
func (entry *EntryImpl) Bootstrap(event rk_query.Event) {
	// do your stuff
	event.AddPair("bootstrap", "true")
}

// Important:
// WaitForShutdownSig won't be called from boot.Shutdown
// User could call this function while start entry only without bootstrapper
// We recommend to call Shutdown in this function
func (entry *EntryImpl) WaitForShutdownSig(duration time.Duration) {}

// Shutdown will be called from boot.Shutdown()
func (entry *EntryImpl) Shutdown(event rk_query.Event) {
	// do your stuff
	event.AddPair("shutdown", "true")
}

// Return name of entry
func (entry *EntryImpl) GetName() string {
	return entry.name
}

// Return type of entry
func (entry *EntryImpl) GetType() string {
	return "example-entry"
}

// Modify as needed
func (entry *EntryImpl) String() string {
	m := map[string]string{
		"name": entry.GetName(),
		"type": entry.GetType(),
		"key":  entry.key,
	}

	bytes, _ := json.Marshal(m)

	return string(bytes)
}
