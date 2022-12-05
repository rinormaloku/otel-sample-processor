package otelsampleprocessor

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

const (
	// The value of "type" key in configuration.
	processorType = "glooMetricsProcessor"
	stability     = component.StabilityLevelDevelopment
)

var processorCapabilities = consumer.Capabilities{MutatesData: true}

func NewFactory() component.ProcessorFactory {
	return component.NewProcessorFactory(
		processorType,
		createDefaultConfig,
		component.WithMetricsProcessor(createGlooMetricsProcessor, stability),
	)
}

func createDefaultConfig() component.ProcessorConfig {
	return &GlooProcessorConfig{
		ProcessorSettings: config.NewProcessorSettings(component.NewID(processorType)),
	}
}

func createGlooMetricsProcessor(
	ctx context.Context,
	set component.ProcessorCreateSettings,
	cfg component.ProcessorConfig,
	nextConsumer consumer.Metrics,
) (component.MetricsProcessor, error) {
	oCfg := cfg.(*GlooProcessorConfig)
	glooMetricsProcessor, err := newGlooMetricProcessor(set.Logger, oCfg)

	if err != nil {
		return nil, err
	}
	return processorhelper.NewMetricsProcessor(
		ctx,
		set,
		cfg,
		nextConsumer,
		glooMetricsProcessor.processMetrics,
		processorhelper.WithCapabilities(processorCapabilities))
}

type GlooMetricsProcessor struct {
	consumer.Metrics
	Started bool
	Stopped bool
}

func (ep *GlooMetricsProcessor) Start(_ context.Context, _ component.Host) error {
	ep.Started = true
	return nil
}

func (ep *GlooMetricsProcessor) Shutdown(_ context.Context) error {
	ep.Stopped = true
	return nil
}

func (ep *GlooMetricsProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}
