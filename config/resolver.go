package config

import (
	"github.com/gonamore/fxbd/config/models"
	"github.com/gonamore/fxbd/config/readers"
	"github.com/imdario/mergo"
	"log"
)

type Resolver struct {
}

// Returns conf object in order: YAML, ENV.
func (rcv *Resolver) GetConfig() (*models.ApplicationConfig, error) {
	jsonReader := readers.NewJsonApplicationConfigReader("application_config.json")
	yamlReader := readers.NewYamlApplicationConfigReader("application_config.yaml")
	envReader := readers.NewEnvApplicationConfigReader()

	return merge(jsonReader, yamlReader, envReader)
}

func merge(readers ...readers.ApplicationConfigReader) (*models.ApplicationConfig, error) {
	config := &models.ApplicationConfig{}

	for _, configReader := range readers {
		newConfig, err := configReader.Read()
		if err != nil {
			log.Println(err)
			continue
		}
		err = mergo.Merge(config, newConfig, mergo.WithOverride)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func NewResolver() *Resolver {
	return &Resolver{}
}
