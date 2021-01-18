// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_info

import (
	"bytes"
	"encoding/json"
	"github.com/hako/durafmt"
	"github.com/rookie-ninja/rk-common/context"
	"go.uber.org/zap"
	"os"
	"os/user"
	"time"
)

// Basic information for a running application
// 1: UID - user id which runs process
// 2: GID - group id which runs process
// 3: Username - username which runs process
// 4: StartTime - application start time
// 5: UpTimeSec - application up time in seconds
// 6: UpTimeStr - application up time in string
// 7: ApplicationName - name of application
// 8: Region - region where process runs
// 9: AZ - availability zone where process runs
// 10: Realm - realm where process runs
// 11: Domain - domain where process runs
type BasicInfo struct {
	UID             string `json:"uid"`
	GID             string `json:"gid"`
	Username        string `json:"username"`
	StartTime       string `json:"start_time"`
	UpTimeSec       int64  `json:"up_time_sec"`
	UpTimeStr       string `json:"up_time_str"`
	ApplicationName string `json:"application_name"`
	Region          string `json:"region"`
	AZ              string `json:"az"`
	Realm           string `json:"realm"`
	Domain          string `json:"domain"`
}

func BasicInfoToStruct() *BasicInfo {
	u, err := user.Current()
	// assign unknown value to user in order to prevent panic
	if err != nil {
		u = &user.User{
			Name: "unknown",
			Uid:  "unknown",
			Gid:  "unknown",
		}
	}
	return &BasicInfo{
		Username:        u.Name,
		UID:             u.Uid,
		GID:             u.Gid,
		StartTime:       rk_ctx.GlobalAppCtx.GetStartTime().Format(time.RFC3339),
		UpTimeSec:       int64(rk_ctx.GlobalAppCtx.GetUpTime().Seconds()),
		UpTimeStr:       durafmt.ParseShort(rk_ctx.GlobalAppCtx.GetUpTime()).String(),
		ApplicationName: rk_ctx.GlobalAppCtx.GetApplicationName(),
		Realm:           os.Getenv("REALM"),
		Region:          os.Getenv("REGION"),
		AZ:              os.Getenv("AZ"),
		Domain:          os.Getenv("DOMAIN"),
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

// Marshal struct to json string
func structToJSON(src interface{}) string {
	return string(structToBytes(src))
}

// Marshal struct to pretty json string
func structToJSONPretty(src interface{}) string {
	mid := structToBytes(src)
	dest := &bytes.Buffer{}
	if err := json.Indent(dest, mid, "", "  "); err != nil {
		return "{}"
	}

	return dest.String()
}

// Marshal struct to bytes
func structToBytes(src interface{}) []byte {
	bytes, _ := json.Marshal(src)
	return bytes
}

// Convert struct to map
func structToMap(src interface{}) map[string]interface{} {
	bytes := structToBytes(src)
	res := make(map[string]interface{})

	// just catch the error
	if err := json.Unmarshal(bytes, &res); err != nil {
		return res
	}

	return res
}

// Convert struct to zap fields
func structToFields(src interface{}) []zap.Field {
	mid := structToMap(src)
	fields := make([]zap.Field, 0)

	for k, v := range mid {
		fields = append(fields, zap.Any(k, v))
	}

	return fields
}
