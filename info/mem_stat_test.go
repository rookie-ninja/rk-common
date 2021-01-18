package rk_info

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemStatsToStruct_HappyCase(t *testing.T) {
	stats := MemStatsToStruct()
	assert.NotNil(t, stats)
	assert.True(t, stats.MemAllocByte > 0)
	assert.True(t, stats.SysAllocByte > 0)
	assert.True(t, stats.MemPercentage > 0)
	assert.NotEmpty(t, stats.LastGCTimestamp)
}

func TestMemStatsToJSON_HappyCase(t *testing.T) {
	assert.NotEmpty(t, MemStatsToJSON())
}

func TestMemStatsToJSONPretty_HappyCase(t *testing.T) {
	assert.NotEmpty(t, MemStatsToJSONPretty())
}

func TestMemStatsToBytes_HappyCase(t *testing.T) {
	assert.NotEmpty(t, MemStatsToBytes())
}

func TestMemStatsToMap_HappyCase(t *testing.T) {
	assert.True(t, len(MemStatsToMap()) > 0)
}

func TestMemStatsToFields_HappyCase(t *testing.T) {
	assert.NotEmpty(t, MemStatsToFields())
}
