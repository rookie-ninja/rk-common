// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rkcommon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Generate random string.
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// Override lumberjack config.
// This function will override fields of non empty and non-nil.
func OverrideLumberjackConfig(origin *lumberjack.Logger, override *lumberjack.Logger) {
	if override == nil {
		return
	}
	origin.Compress = override.Compress
	origin.LocalTime = override.LocalTime
	if override.MaxAge > 0 {
		origin.MaxAge = override.MaxAge
	}

	if override.MaxBackups > 0 {
		origin.MaxBackups = override.MaxBackups
	}

	if override.MaxSize > 0 {
		origin.MaxSize = override.MaxSize
	}

	if len(override.Filename) > 0 {
		origin.Filename = override.Filename
	}
}

// Override zap config.
// This function will override fields of non empty and non-nil.
func OverrideZapConfig(origin *zap.Config, override *zap.Config) {
	if override == nil {
		return
	}

	// by default, these fields would be false
	// so just override it with new config
	origin.Development = override.Development
	origin.DisableCaller = override.DisableCaller
	origin.DisableStacktrace = override.DisableStacktrace

	if len(override.Encoding) > 0 {
		origin.Encoding = override.Encoding
	}

	if !reflect.ValueOf(override.Level).Field(0).IsNil() {
		origin.Level.SetLevel(override.Level.Level())
	}

	if len(override.InitialFields) > 0 {
		origin.InitialFields = override.InitialFields
	}

	if len(override.ErrorOutputPaths) > 0 {
		origin.ErrorOutputPaths = override.ErrorOutputPaths
	}

	if len(override.OutputPaths) > 0 {
		origin.OutputPaths = override.OutputPaths
	}

	if override.Sampling != nil {
		origin.Sampling = override.Sampling
	}

	// deal with encoder config
	if len(override.EncoderConfig.CallerKey) > 0 {
		origin.EncoderConfig.CallerKey = override.EncoderConfig.CallerKey
	}

	if len(override.EncoderConfig.ConsoleSeparator) > 0 {
		origin.EncoderConfig.ConsoleSeparator = override.EncoderConfig.ConsoleSeparator
	}

	if override.EncoderConfig.EncodeCaller != nil {
		origin.EncoderConfig.EncodeCaller = override.EncoderConfig.EncodeCaller
	}

	if override.EncoderConfig.EncodeDuration != nil {
		origin.EncoderConfig.EncodeDuration = override.EncoderConfig.EncodeDuration
	}

	if override.EncoderConfig.EncodeLevel != nil {
		origin.EncoderConfig.EncodeLevel = override.EncoderConfig.EncodeLevel
	}

	if override.EncoderConfig.EncodeName != nil {
		origin.EncoderConfig.EncodeName = override.EncoderConfig.EncodeName
	}

	if override.EncoderConfig.EncodeTime != nil {
		origin.EncoderConfig.EncodeTime = override.EncoderConfig.EncodeTime
	}

	if len(override.EncoderConfig.MessageKey) > 0 {
		origin.EncoderConfig.MessageKey = override.EncoderConfig.MessageKey
	}

	if len(override.EncoderConfig.LevelKey) > 0 {
		origin.EncoderConfig.LevelKey = override.EncoderConfig.LevelKey
	}

	if len(override.EncoderConfig.TimeKey) > 0 {
		origin.EncoderConfig.TimeKey = override.EncoderConfig.TimeKey
	}

	if len(override.EncoderConfig.NameKey) > 0 {
		origin.EncoderConfig.NameKey = override.EncoderConfig.NameKey
	}

	if len(override.EncoderConfig.CallerKey) > 0 {
		origin.EncoderConfig.CallerKey = override.EncoderConfig.CallerKey
	}

	if len(override.EncoderConfig.FunctionKey) > 0 {
		origin.EncoderConfig.FunctionKey = override.EncoderConfig.FunctionKey
	}

	if len(override.EncoderConfig.StacktraceKey) > 0 {
		origin.EncoderConfig.StacktraceKey = override.EncoderConfig.StacktraceKey
	}

	if len(override.EncoderConfig.LineEnding) > 0 {
		origin.EncoderConfig.LineEnding = override.EncoderConfig.LineEnding
	}
}

// Shutdown and panic.
func ShutdownWithError(err error) {
	if err == nil {
		err = errors.New("error is nil")
	}
	panic(err)
}

// Read files with provided path, use working directory if given path is relative path.
// Shutdown process if any error occurs, this should be used for MUST SUCCESS scenario like reading config files.
func MustReadFile(filePath string) []byte {
	if len(filePath) < 1 {
		ShutdownWithError(errors.New("empty file path"))
	}

	if !path.IsAbs(filePath) {
		wd, err := os.Getwd()

		if err != nil {
			ShutdownWithError(err)
		}
		filePath = path.Join(wd, filePath)
	}

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		ShutdownWithError(err)
	}

	return bytes
}

// Check File existence, file path should be full path.
func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Return default value if original string is empty.
func GetDefaultIfEmptyString(origin, def string) string {
	if len(origin) < 1 {
		return def
	}

	context.Background()

	return origin
}

// Return default value if environment variable is empty or not exist.
func GetEnvValueOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)

	if len(value) < 1 {
		return defaultValue
	}

	return value
}

// This is a tricky function.
// We will iterate through all the network interfaces，but will choose the first one since we are assuming that
// eth0 will be the default one to use in most of the case.
//
// Currently, we do not have any interfaces for selecting the network interface yet.
func GetLocalIP() string {
	localIP := "localhost"

	// skip the error since we don't want to break RPC calls because of it
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return localIP
	}

	for _, addr := range addresses {
		items := strings.Split(addr.String(), "/")
		if len(items) < 2 || items[0] == "127.0.0.1" {
			continue
		}

		if match, err := regexp.MatchString(`\d+\.\d+\.\d+\.\d+`, items[0]); err == nil && match {
			localIP = items[0]
		}
	}

	return localIP
}

// Get hostname of localhost, return "unknown" if error occurs or hostname is empty.
func GetLocalHostname() string {
	hostname, err := os.Hostname()
	if err != nil || len(hostname) < 1 {
		hostname = "unknown"
	}

	return hostname
}

// Generate request id based on google/uuid.
// UUIDs are based on RFC 4122 and DCE 1.1: Authentication and Security Services.
//
// A UUID is a 16 byte (128 bit) array. UUIDs may be used as keys to maps or compared directly.
func GenerateRequestId() string {
	// do not use uuid.New() since it would panic if any error occurs
	requestId, err := uuid.NewRandom()

	// currently, we will return empty string if error occurs
	if err != nil {
		return ""
	}

	return requestId.String()
}

// Generate request id based on google/uuid.
// UUIDs are based on RFC 4122 and DCE 1.1: Authentication and Security Services.
//
// A UUID is a 16 byte (128 bit) array. UUIDs may be used as keys to maps or compared directly.
func GenerateRequestIdWithPrefix(prefix string) string {
	// Do not use uuid.New() since it would panic if any error occurs
	requestId, err := uuid.NewRandom()

	// Currently, we will return empty string if error occurs
	if err != nil {
		return ""
	}

	if len(prefix) > 0 {
		return prefix + "-" + requestId.String()
	}

	return requestId.String()
}

// Override source map with new map items.
// It will iterate through all items in map and check map and slice types of item to recursively override values
//
// Mainly used for unmarshalling YAML to map.
func OverrideMap(src map[interface{}]interface{}, override map[interface{}]interface{}) {
	if src == nil || override == nil {
		return
	}

	for k, overrideItem := range override {
		originalItem, ok := src[k]
		if ok && reflect.TypeOf(originalItem) == reflect.TypeOf(overrideItem) {
			switch overrideItem.(type) {
			case []interface{}:
				OverrideSlice(originalItem.([]interface{}), overrideItem.([]interface{}))
			case map[interface{}]interface{}:
				OverrideMap(originalItem.(map[interface{}]interface{}), overrideItem.(map[interface{}]interface{}))
			default:
				src[k] = overrideItem
			}
		}
	}
}

// Override source slice with new slice items.
// It will iterate through all items in slice and check map and slice types of item to recursively override values
//
// Mainly used for unmarshalling YAML to map.
func OverrideSlice(src []interface{}, override []interface{}) {
	if src == nil || override == nil {
		return
	}

	for i := range override {
		if override[i] != nil && len(src)-1 >= i && reflect.TypeOf(override[i]) == reflect.TypeOf(src[i]) {
			overrideItem := override[i]
			originalItem := src[i]
			switch overrideItem.(type) {
			case []interface{}:
				OverrideSlice(originalItem.([]interface{}), overrideItem.([]interface{}))
			case map[interface{}]interface{}:
				OverrideMap(originalItem.(map[interface{}]interface{}), overrideItem.(map[interface{}]interface{}))
			default:
				src[i] = override[i]
			}
		}
	}
}

// Convert JSON style string to map[string]interface{}.
// Return empty map if length of input parameter is less than 2 which can not construct
// a valid JSON string.
func ConvertJSONToMap(str string) map[string]interface{} {
	res := make(map[string]interface{})
	if len(str) < 2 {
		return res
	}

	json.Unmarshal([]byte(str), &res)

	return res
}

// Marshal struct to json string.
// Return empty string if input parameter is nil.
func ConvertStructToJSON(src interface{}) string {
	if src == nil {
		return ""
	}

	return string(ConvertStructToBytes(src))
}

// Marshal struct to pretty json string.
// Return empty string if input parameter is nil.
func ConvertStructToJSONPretty(src interface{}) string {
	if src == nil {
		return ""
	}

	mid := ConvertStructToBytes(src)
	dest := &bytes.Buffer{}
	if err := json.Indent(dest, mid, "", "  "); err != nil {
		return "{}"
	}

	return dest.String()
}

// Marshal struct to bytes.
// Return empty byte slice if input parameter is nil.
func ConvertStructToBytes(src interface{}) []byte {
	if src == nil {
		return []byte{}
	}
	bytes, _ := json.Marshal(src)
	return bytes
}

// Convert struct to map.
// Return empty map if input parameter is nil.
func ConvertStructToMap(src interface{}) map[string]interface{} {
	res := make(map[string]interface{})

	if src == nil {
		return res
	}

	bytes := ConvertStructToBytes(src)

	// just catch the error
	if err := json.Unmarshal(bytes, &res); err != nil {
		return res
	}

	return res
}

// Convert struct to zap fields.
// Return empty zap.Field array if input parameter is nil.
func ConvertStructToZapFields(src interface{}) []zap.Field {
	fields := make([]zap.Field, 0)
	if src == nil {
		return fields
	}
	mid := ConvertStructToMap(src)

	for k, v := range mid {
		fields = append(fields, zap.Any(k, v))
	}

	return fields
}
