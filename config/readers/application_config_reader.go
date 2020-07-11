package readers

import (
	model2 "github.com/gonamore/fxbd/config/models"
)

// Interface for reading application configuration
type ApplicationConfigReader interface {
	Read() (*model2.ApplicationConfig, error)
}
