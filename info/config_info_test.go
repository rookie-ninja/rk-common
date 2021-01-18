package rk_info

import (
	"github.com/rookie-ninja/rk-common/context"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestViperConfigToStruct_WithEmptyViperConfigs(t *testing.T) {
	assert.Empty(t, ViperConfigToStruct())
}

func TestViperConfigToStruct_HappyCase(t *testing.T) {
	rk_ctx.GlobalAppCtx.AddViperConfig("unit-test-config", viper.New())
	assert.Len(t, ViperConfigToStruct(), 1)
}

func TestViperConfigToJSON_HappyCase(t *testing.T) {
	rk_ctx.GlobalAppCtx.AddViperConfig("unit-test-config", viper.New())
	assert.NotEmpty(t, ViperConfigToJSON())
}

func TestViperConfigToJSONPretty_HappyCase(t *testing.T) {
	rk_ctx.GlobalAppCtx.AddViperConfig("unit-test-config", viper.New())
	assert.NotEmpty(t, ViperConfigToJSONPretty())
}

func TestViperConfigToBytes_HappyCase(t *testing.T) {
	rk_ctx.GlobalAppCtx.AddViperConfig("unit-test-config", viper.New())
	assert.NotEmpty(t, ViperConfigToBytes())
}

func TestViperConfigToMap_HappyCase(t *testing.T) {
	rk_ctx.GlobalAppCtx.AddViperConfig("unit-test-config", viper.New())
	assert.True(t, len(ViperConfigToMap()) == 1)
}

func TestViperConfigToFields_HappyCase(t *testing.T) {
	rk_ctx.GlobalAppCtx.AddViperConfig("unit-test-config", viper.New())
	assert.NotEmpty(t, ViperConfigToFields())
}
