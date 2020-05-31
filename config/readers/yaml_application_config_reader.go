package readers

import (
	"github.com/gonamore/fxbd/config/models"
	"gopkg.in/yaml.v2"
	"os"
)

// Reads application configuration properties from specified YAML file
type YamlApplicationConfigReader struct {
	path string

	ApplicationConfigReader
}

func (rcv *YamlApplicationConfigReader) Read() (*models.ApplicationConfig, error) {
	file, err := os.Open(rcv.path)
	if err != nil {
		return nil, err
	}

	decoder := yaml.NewDecoder(file)

	config := &models.ApplicationConfig{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func NewYamlApplicationConfigReader(path string) *YamlApplicationConfigReader {
	return &YamlApplicationConfigReader{path: path}
}
