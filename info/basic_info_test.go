package rk_info

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBasicInfoToStruct_HappyCase(t *testing.T) {
	os.Setenv("REALM", "unit-test-realm")
	os.Setenv("REGION", "unit-test-region")
	os.Setenv("AZ", "unit-test-az")
	os.Setenv("DOMAIN", "unit-test-domain")

	info := BasicInfoToStruct()
	assert.NotNil(t, info)

	assert.NotEmpty(t, info.Username)
	assert.NotEmpty(t, info.UID)
	assert.NotEmpty(t, info.GID)
	assert.NotEmpty(t, info.StartTime)
	assert.Zero(t, info.UpTimeSec)
	assert.NotEmpty(t, info.UpTimeStr)
	assert.NotEmpty(t, info.ApplicationName)
	assert.Equal(t, "unit-test-realm", info.Realm)
	assert.Equal(t, "unit-test-region", info.Region)
	assert.Equal(t, "unit-test-az", info.AZ)
	assert.Equal(t, "unit-test-domain", info.Domain)
}

func TestBasicInfoToJSON_HappyCase(t *testing.T) {
	assert.NotNil(t, BasicInfoToJSON())
}

func TestBasicInfoToJSONPretty_HappyCase(t *testing.T) {
	assert.NotNil(t, BasicInfoToJSONPretty())
}

func TestBasicInfoToBytes_HappyCase(t *testing.T) {
	assert.NotEmpty(t, BasicInfoToBytes())
}

func TestBasicInfoToMap_HappyCase(t *testing.T) {
	assert.True(t, len(BasicInfoToMap()) > 0)
}

func TestBasicInfoToZapFields_HappyCase(t *testing.T) {
	assert.NotEmpty(t, BasicInfoToFields())
}

func TestStructToBytes_HappyCase(t *testing.T) {
	assert.NotEmpty(t, structToBytes(BasicInfoToStruct()))
}

func TestStructToJSON_HappyCase(t *testing.T) {
	assert.NotEmpty(t, structToJSON(BasicInfoToStruct()))
}

func TestStructToJSONPretty_HappyCase(t *testing.T) {
	assert.NotEmpty(t, structToJSONPretty(BasicInfoToStruct()))
}

func TestStructToMap_HappyCase(t *testing.T) {
	assert.NotEmpty(t, structToMap(BasicInfoToStruct()))
}

func TestStructToFields_HappyCase(t *testing.T) {
	assert.NotEmpty(t, structToFields(BasicInfoToStruct()))
}
