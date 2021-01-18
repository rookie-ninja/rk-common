package rk_metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	opts = prometheus.SummaryOpts{
		Namespace:  "namespace",
		Subsystem:  "subSystem",
		Name:       "name",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001, 0.999: 0.0001},
	}

	labelKeys   = []string{"path", "res_code"}
	summaryVec  = prometheus.NewSummaryVec(opts, labelKeys)
	observer, _ = summaryVec.GetMetricWithLabelValues("GET", "200")
)

func TestGetRequestMetrics_HappyCase(t *testing.T) {
	observer.Observe(100)
	observer.Observe(200)
	metrics := GetRequestMetrics(summaryVec)
	assert.Len(t, metrics, 1)

	metric := metrics[0]

	assert.Equal(t, uint64(2), metric.Count)
	assert.Equal(t, "GET", metric.Path)
	assert.Len(t, metric.ResCode, 1)
	assert.Equal(t, "200", metric.ResCode[0].ResCode)
	assert.Equal(t, uint64(2), metric.ResCode[0].Count)
}
