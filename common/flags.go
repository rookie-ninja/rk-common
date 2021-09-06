// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkcommon

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path"
)

var (
	// GlobalFlags will read pflags passed while starting main entry
	GlobalFlags *pflag.FlagSet
)

const (
	BootConfigPathFlagKey = "rkboot"
	BootConfigOverrideKey = "rkset"
)

// pflag.FlagSet which contains rkboot and rkset as key.
//
// Usage of rkboot:
// Receives path of boot config file path, can be either absolute of relative path.
// If relative path was provided, then current working directory would be attached in front of provided path.
// example:
// ./your_compiled_binary --rkboot <your path to config file>
//
// Usage of rkset:
//
// Receives flattened boot config file(YAML) keys and override them in provided boot config.
// We follow the way of HELM does while overriding keys, refer to https://helm.sh/docs/intro/using_helm/
// example:
//
// Lets assuming we have boot config YAML file as bellow:
//
// example-boot.yaml:
// gin:
//   - port: 1949
//     commonService:
//       enabled: true
//
// We can override values in example-boot.yaml file as bellow:
//
// ./your_compiled_binary --rkboot example-boot.yaml --rkset "gin[0].port=2008,gin[0].commonService.enabled=false"
//
// Basic rules:
// 1: Using comma(,) to separate different k/v section.
// 2: Using [index] to access arrays in YAML file.
// 3: Using equal sign(=) to distinguish key and value.
// 4: Using dot(.) to access map in YAML file.
func init() {
	// GlobalFlags will continue with error
	GlobalFlags = pflag.NewFlagSet("rk", pflag.ContinueOnError)
	GlobalFlags.String(BootConfigPathFlagKey, "", "set config file path")
	GlobalFlags.String(BootConfigOverrideKey, "", "set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	GlobalFlags.Parse(os.Args[1:])
}

// GetBootConfigPath this function will do the following things.
// First, override config file path if --rkboot <config file path> was provided by user.
// Second, join path with current working directory if user provided path is relative path.
// Finally, validate file existence, shutdown process if file is missing.
func GetBootConfigPath(configFilePath string) string {
	// get config file path overrides from input args, shutdown apps if unexpected error occur
	if pathFromFlag, err := GlobalFlags.GetString(BootConfigPathFlagKey); err != nil {
		ShutdownWithError(err)
	} else if len(pathFromFlag) > 0 {
		configFilePath = pathFromFlag
	}

	// join the path with current working directory if user provided path is relative path
	if !path.IsAbs(configFilePath) {
		wd, err := os.Getwd()

		if err != nil {
			ShutdownWithError(err)
		}
		configFilePath = path.Join(wd, configFilePath)
	}

	// validate file existence，shutdown if config file does not exists
	if !FileExists(configFilePath) {
		ShutdownWithError(fmt.Errorf("config file does not exist with path:%s", configFilePath))
	}

	return configFilePath
}

// GetBootConfigOverrides this function will read user provided config content overrides and construct it into a map.
func GetBootConfigOverrides() map[interface{}]interface{} {
	bootConfigOverrides, err := GlobalFlags.GetString(BootConfigOverrideKey)

	if err != nil {
		ShutdownWithError(err)
	}

	res, err := ParseBootConfigOverrides(bootConfigOverrides)

	if err != nil {
		ShutdownWithError(err)
	}

	return res
}

// GetBootConfigOriginal read config file content and unmarshal into map.
func GetBootConfigOriginal(configFilePath string) map[interface{}]interface{} {
	configFilePath = GetBootConfigPath(configFilePath)

	vp := viper.New()
	vp.SetConfigFile(configFilePath)

	if err := vp.ReadInConfig(); err != nil {
		ShutdownWithError(err)
	}

	originalMap := make(map[interface{}]interface{})
	if err := vp.Unmarshal(&originalMap); err != nil {
		ShutdownWithError(err)
	}

	return originalMap
}

// UnmarshalBootConfig this function is combination of GetBootConfigPath, GetBootConfigOverrides and
// GetBootConfigOriginal.
// User who want to implement his/her own entry, may use this function to parse YAML config into struct.
// This function would also parse --rkset flags.
//
// This function would do the following:
// First, read config file and unmarshal content into a map (--rkboot flag would be read).
// Second, read --rkset flags and override values in map unmarshalled at above step.
// Finally, unmarshal map into user provided struct.
func UnmarshalBootConfig(configFilePath string, config interface{}) {
	// 1: unmarshal config file into map
	configMap := GetBootConfigOriginal(configFilePath)

	// 2: read command line flags and override original config map with flags
	OverrideMap(configMap, GetBootConfigOverrides())

	// 3: decode config map into boot config struct
	if err := mapstructure.Decode(configMap, config); err != nil {
		ShutdownWithError(err)
	}
}
