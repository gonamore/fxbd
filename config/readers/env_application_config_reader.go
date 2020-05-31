package readers

import (
	"github.com/gonamore/fxbd/config/models"
	"github.com/kelseyhightower/envconfig"
)

// Reads application configuration properties from OS Environment Variables
type EnvApplicationConfigReader struct {
	ApplicationConfigReader
}

func (rcv *EnvApplicationConfigReader) Read() (*models.ApplicationConfig, error) {
	config := &models.ApplicationConfig{}

	err := envconfig.Process("", config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func NewEnvApplicationConfigReader() *EnvApplicationConfigReader {
	return &EnvApplicationConfigReader{}
}
