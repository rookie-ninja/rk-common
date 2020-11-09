// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_info

import (
	"github.com/hako/durafmt"
	"github.com/rookie-ninja/rk-common/context"
	"go.uber.org/zap"
	"os"
	"os/user"
	"time"
)

// Basic information for a running applilcation
type BasicInfo struct {
	UID         string `json:"uid"`
	GID         string `json:"gid"`
	Username    string `json:"username"`
	StartTime   string `json:"start_time"`
	UpTimeSec   int64  `json:"up_time_sec"`
	UpTimeStr   string `json:"up_time_str"`
	Application string `json:"application"`
	Region      string `json:"region"`
	AZ          string `json:"az"`
	Realm       string `json:"realm"`
	Domain      string `json:"domain"`
}

func BasicInfoToStruct() *BasicInfo {
	u, _ := user.Current()
	return &BasicInfo{
		Username:    u.Name,
		UID:         u.Uid,
		GID:         u.Gid,
		StartTime:   rk_ctx.GlobalAppCtx.GetStartTime().Format(time.RFC3339),
		UpTimeSec:   int64(rk_ctx.GlobalAppCtx.GetUpTime().Seconds()),
		UpTimeStr:   durafmt.ParseShort(rk_ctx.GlobalAppCtx.GetUpTime()).String(),
		Application: rk_ctx.GlobalAppCtx.GetApplication(),
		Realm:       os.Getenv("REALM"),
		Region:      os.Getenv("REGION"),
		AZ:          os.Getenv("AZ"),
		Domain:      os.Getenv("DOMAIN"),
	}
}

// As json string
func BasicInfoToJSON() string {
	return structToJSON(BasicInfoToStruct())
}

// As pretty json string
func BasicInfoToJSONPretty() string {
	return structToJSONPretty(BasicInfoToStruct())
}

// As byte array
func BasicInfoToBytes() []byte {
	return structToBytes(BasicInfoToStruct())
}

// As a map
func BasicInfoToMap() map[string]interface{} {
	return structToMap(BasicInfoToStruct())
}

// As zap.Field
func BasicInfoToFields() []zap.Field {
	return structToFields(BasicInfoToStruct())
}
