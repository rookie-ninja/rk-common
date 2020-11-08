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

// Config
type ConfigInfo struct {
	Name  string        `json:"name"`
	Raw   string        `json:"raw"`
}

// As struct
func ViperConfigToStruct() []*ConfigInfo {
	res := make([]*ConfigInfo, 0)

	for k, v := range rk_ctx.GlobalAppCtx.ListViperConfigs() {
		res = append(res, &ConfigInfo{
			Name: k,
			Raw: fmt.Sprintf("%v", v.AllSettings()),
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
	return structToMap(ViperConfigToStruct())
}

// As zap.Field
func ViperConfigToFields() []zap.Field {
	return structToFields(ViperConfigToStruct())
}

func RkConfigToStruct() []*ConfigInfo {
	res := make([]*ConfigInfo, 0)

	for k, v := range rk_ctx.GlobalAppCtx.ListRkConfigs() {
		res = append(res, &ConfigInfo{
			Name: k,
			Raw: fmt.Sprintf("%v", v.GetViper().AllSettings()),
		})
	}

	return res
}

// As json string
func RkConfigToJSON() string {
	return structToJSON(RkConfigToStruct())
}

// As pretty json string
func RkConfigToJSONPretty() string {
	return structToJSONPretty(RkConfigToStruct())
}

// As byte array
func RkConfigToBytes() []byte {
	return structToBytes(RkConfigToStruct())
}

// As map
func RkConfigToMap() map[string]interface{} {
	return structToMap(RkConfigToStruct())
}

// As zap.Field
func RkConfigToFields() []zap.Field {
	return structToFields(RkConfigToStruct())
}
