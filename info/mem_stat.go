// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_info

import (
	"go.uber.org/zap"
	"runtime"
	"time"
)

// Memory stats of current running process
type MemStats struct {
	MemAllocByte    uint64  `json:"mem_alloc_byte"`
	SysAllocByte    uint64  `json:"sys_alloc_byte"`
	MemPercentage   float64 `json:"mem_usage_percentage"`
	LastGCTimestamp string  `json:"last_gc_timestamp"`
	GCCount         uint32  `json:"gc_count_total"`
	ForceGCCount    uint32  `json:"force_gc_count"`
}

func MemStatsToStruct() *MemStats {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	return &MemStats{
		MemAllocByte:    stats.Alloc,
		SysAllocByte:    stats.Sys,
		MemPercentage:   float64(stats.Alloc) / float64(stats.Sys),
		LastGCTimestamp: time.Unix(0, int64(stats.LastGC)).Format(time.RFC3339),
		GCCount:         stats.NumGC,
		ForceGCCount:    stats.NumForcedGC,
	}
}

func MemStatsToJSON() string {
	return structToJSON(MemStatsToStruct())
}

func MemStatsToJSONPretty() string {
	return structToJSONPretty(MemStatsToStruct())
}

func MemStatsToBytes() []byte {
	return structToBytes(MemStatsToStruct())
}

func MemStatsToMap() map[string]interface{} {
	return structToMap(MemStatsToStruct())
}

func MemStatsToFields() []zap.Field {
	return structToFields(MemStatsToStruct())
}
