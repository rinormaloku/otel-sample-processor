package otel_sample_processor

import (
	"context"
	"fmt"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"time"

	"go.uber.org/zap"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

type glooMetricProcessor struct {
	cfg    *GlooProcessorConfig
	logger *zap.Logger
}

func newGlooMetricProcessor(logger *zap.Logger, cfg *GlooProcessorConfig) (*glooMetricProcessor, error) {
	logger.Info("Metric timestamp processor configured")
	return &glooMetricProcessor{cfg: cfg, logger: logger}, nil
}

// processMetrics takes incoming metrics and adjusts their timestamps to
// the nearest time unit (specified by duration in the config)
func (fmp *glooMetricProcessor) processMetrics(_ context.Context, src pmetric.Metrics) (pmetric.Metrics, error) {
	// set the timestamps to the nearest time unit
	for i := 0; i < src.ResourceMetrics().Len(); i++ {
		rm := src.ResourceMetrics().At(i)
		for j := 0; j < rm.ScopeMetrics().Len(); j++ {
			ms := rm.ScopeMetrics().At(j)
			for k := 0; k < ms.Metrics().Len(); k++ {
				m := ms.Metrics().At(k)

				switch m.Type() {
				case pmetric.MetricTypeGauge:
					dataPoints := m.Gauge().DataPoints()
					for l := 0; l < dataPoints.Len(); l++ {
						gotDataPoint := dataPoints.At(l)
						snappedTimestamp := gotDataPoint.Timestamp().AsTime().Truncate(time.Hour)
						gotDataPoint.SetTimestamp(pcommon.NewTimestampFromTime(snappedTimestamp))
					}
				case pmetric.MetricTypeSum:
					dataPoints := m.Sum().DataPoints()
					for l := 0; l < dataPoints.Len(); l++ {
						gotDataPoint := dataPoints.At(l)
						snappedTimestamp := gotDataPoint.Timestamp().AsTime().Truncate(time.Hour)
						gotDataPoint.SetTimestamp(pcommon.NewTimestampFromTime(snappedTimestamp))
					}
				case pmetric.MetricTypeHistogram:
					dataPoints := m.Histogram().DataPoints()
					for l := 0; l < dataPoints.Len(); l++ {
						gotDataPoint := dataPoints.At(l)
						snappedTimestamp := gotDataPoint.Timestamp().AsTime().Truncate(time.Hour)
						gotDataPoint.SetTimestamp(pcommon.NewTimestampFromTime(snappedTimestamp))
					}
				case pmetric.MetricTypeSummary:
					dataPoints := m.Summary().DataPoints()
					for l := 0; l < dataPoints.Len(); l++ {
						gotDataPoint := dataPoints.At(l)
						snappedTimestamp := gotDataPoint.Timestamp().AsTime().Truncate(time.Hour)
						gotDataPoint.SetTimestamp(pcommon.NewTimestampFromTime(snappedTimestamp))
					}
				default:
					fmt.Printf("Unknown type")
					return src, fmt.Errorf("unknown type: %s", m.Type().String())
				}
			}
		}
	}
	return src, nil
}
