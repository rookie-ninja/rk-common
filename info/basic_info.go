package rk_info

import (
	"github.com/rookie-ninja/rk-common/context"
	"go.uber.org/zap"
	"os"
	"time"
)

// Basic information for a running applilcation
type BasicInfo struct {
	StartTime   string `json:"start_time"`
	UpTime      string `json:"up_time"`
	Application string `json:"application"`
	Region      string `json:"region"`
	AZ          string `json:"az"`
	Realm       string `json:"realm"`
	Domain      string `json:"domain"`
}

func BasicInfoToStruct() *BasicInfo {
	return &BasicInfo{
		StartTime:   rk_ctx.GlobalAppCtx.GetStartTime().Format(time.RFC3339),
		UpTime:      rk_ctx.GlobalAppCtx.GetUpTime().String(),
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
func BasicInfoToFields() []zap.Field{
	return structToFields(BasicInfoToStruct())
}