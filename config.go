package otelsampleprocessor

import (
	"go.opentelemetry.io/collector/config"
)

type GlooProcessorConfig struct {
	config.ProcessorSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
}
