// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_info

import (
	"fmt"
	"github.com/rookie-ninja/rk-common/context"
	"go.uber.org/zap"
)

// Config information
type ConfigInfo struct {
	Name string `json:"name"`
	Raw  string `json:"raw"`
}

// As struct
func ViperConfigToStruct() []*ConfigInfo {
	res := make([]*ConfigInfo, 0)

	for k, v := range rk_ctx.GlobalAppCtx.ListViperConfigs() {
		res = append(res, &ConfigInfo{
			Name: k,
			Raw:  fmt.Sprintf("%v", v.AllSettings()),
		})
	}

	return res
}

// As json string
func ViperConfigToJSON() string {
	return structToJSON(ViperConfigToStruct())
}

// As pretty json string
func ViperConfigToJSONPretty() string {
	return structToJSONPretty(ViperConfigToStruct())
}

// As byte array
func ViperConfigToBytes() []byte {
	return structToBytes(ViperConfigToStruct())
}

// As map
func ViperConfigToMap() map[string]interface{} {
	res := make(map[string]interface{})
	configList := ViperConfigToStruct()
	for i := range configList {
		res[configList[i].Name] = configList[i].Raw
	}
	return res
}

// As zap.Field
func ViperConfigToFields() []zap.Field {
	res := make([]zap.Field, 0)

	configList := ViperConfigToStruct()
	for i := range configList {
		res = append(res, zap.Any(configList[i].Name, configList[i].Raw))
	}

	return res
}
