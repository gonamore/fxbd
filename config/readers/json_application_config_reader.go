package readers

import (
	"encoding/json"
	"github.com/gonamore/fxbd/config/models"
	"io/ioutil"
	"os"
)

// Reads application configuration properties from specified JSON file
type JsonApplicationConfigReader struct {
	path string

	ApplicationConfigReader
}

func (rcv *JsonApplicationConfigReader) Read() (*models.ApplicationConfig, error) {
	file, err := os.Open(rcv.path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	config := &models.ApplicationConfig{}
	err = json.Unmarshal(bytes, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func NewJsonApplicationConfigReader(path string) *JsonApplicationConfigReader {
	return &JsonApplicationConfigReader{path: path}
}
