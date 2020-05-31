package readers

import (
	model2 "github.com/gonamore/fxbd/config/models"
)

type ApplicationConfigReader interface {
	Read() (*model2.ApplicationConfig, error)
}
