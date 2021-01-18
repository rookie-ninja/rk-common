<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [rk-common](#rk-common)
  - [Installation](#installation)
  - [Quick Start](#quick-start)
    - [Context](#context)
    - [Entry](#entry)
    - [Basic Info](#basic-info)
    - [Config info](#config-info)
    - [Memory Stats](#memory-stats)
    - [Request Metrics](#request-metrics)
  - [Contributing](#contributing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# rk-common
The common library mainly used by rk-boot

## Installation
`go get -u github.com/rookie-ninja/rk-common`

## Quick Start

### Context
A struct called AppContext witch contains RK style application metadata.

**How to use AppContext?** 
Access it via GlobalAppCtx variable 
```go
    rk_ctx.GlobalAppCtx
```

| Element | Description | Default |
| ------ | ------ | ------ |
| application | name of running application | empty |
| startTime | application start time | 0001-01-01 00:00:00 +0000 UTC |
| loggers | loggers with a name as a key | empty map |
| viperConfigs | viper configs with a name as a key | empty map |
| eventFactory | event data factory | standard event factory which logs to stdout |
| shutdownSig | a channel receiving shutdown signals | empty channel |
| shutdownHooks | a list of shutdownHook function registered by user | empty list |
| customValues | custom k/v store | empty map |

### Entry
rk_entry.Entry is an interface for rk_boot.Bootstrapper to bootstrap entry.

Users could implement rk_entry.Entry interface and bootstrap any service/process with rk_boot.Bootstrapper

**How to create a new custom entry? Please see example/ for details**
- **Step 1:**
Construct your own entry YAML struct as needed
example:
```yaml
---
myEntry:
  enabled: true
  key: value
```

- **Step 2:**
Create a struct which implements Entry interface. 
```go
// A struct which is for unmarshalled YAML
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
```

- **Step 3:**
Implements EntryRegFunc and define a struct which could be marshaled from YAML config
```go
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
```
- **Step 4:**
Register your reg function in init() in order to register your entry while application starts
```go
// Register entry, must be in init() function since we need to register entry
// at beginning
func init() {
	rk_ctx.RegisterEntry(NewEntry)
}
```

How entry interact with rk-boot.Bootstrapper?**
**
1: Entry will be created and registered into rk_ctx.GlobalAppCtx

2: Bootstrap will be called from Bootstrapper.Bootstrap() function

3: Application will wait for shutdown signal

4: Shutdown will be called from Bootstrapper.Shutdown() function

### Basic Info
Basic information for a running application

**How to get basic information?** 
Access it via BasicInfoToStruct() function 
```go
    rk_info.BasicInfoToStruct()
```

| Element | Description | Default |
| ------ | ------ | ------ |
| UID | user id which runs process | based on system |
| GID | group id which runs process | based on system |
| Username | username which runs process | based on system |
| StartTime | application start time | time.now() |
| UpTimeSec | application up time in seconds | zero at the beginning |
| UpTimeStr | application up time in string | zero as string at the beginning |
| ApplicationName | name of application | empty |
| Region | region where process runs | based on environment variable REGION |
| AZ | availability zone where process runs | based on environment variable AZ |
| Realm | realm where process runs | based on environment variable REALM |
| Domain | domain where process runs | based on environment variable DOMAIN |

### Config info
Config information stored in GlobalAppCtx

### Memory Stats
Memory stats of current running process

**How to get memory stats?** 
Access it via MemStatsToStruct() function 
```go
    rk_info.MemStatsToStruct()
```

| Element | Description | Default |
| ------ | ------ | ------ |
| MemAllocByte | bytes of allocated heap objects, from runtime.MemStats.Alloc | based on system |
| SysAllocByte | total bytes of memory obtained from the OS, from runtime.MemStats.Sys | based on system |
| MemPercentage | float64(stats.Alloc) / float64(stats.Sys) | based on system |
| LastGCTimestamp | the time the last garbage collection finished as RFC3339 | based on system |
| GCCount | the number of completed GC cycles | zero at the beginning |
| ForceGCCount | the number of GC cycles that were forced by the application calling the GC function | zero as string at the beginning |

### Request Metrics
Request metrics to struct from prometheus summary collector

**How to get request metrics?** 
Access it via MemStatsToStruct() function 
```go
    rk_metrics.GetRequestMetrics()
```

| Element | Description | Default |
| ------ | ------ | ------ |
| Path | API path | based on system |
| ElapsedNanoP50 | quantile of p50 with time elapsed | based on prometheus collector |
| ElapsedNanoP90 | quantile of p90 with time elapsed | based on prometheus collector |
| ElapsedNanoP99 | quantile of p99 with time elapsed | based on prometheus collector |
| ElapsedNanoP999 | quantile of p999 with time elapsed | based on prometheus collector |
| Count | total number of requests | based on prometheus collector |
| ResCode | response code labels | based on prometheus collector |



## Contributing
We encourage and support an active, healthy community of contributors &mdash;
including you! Details are in the [contribution guide](CONTRIBUTING.md) and
the [code of conduct](CODE_OF_CONDUCT.md). The rk maintainers keep an eye on
issues and pull requests, but you can also report any negative conduct to
dongxuny@gmail.com. That email list is a private, safe space; even the zap
maintainers don't have access, so don't hesitate to hold us to a high
standard.

<hr>

Released under the [MIT License](LICENSE).