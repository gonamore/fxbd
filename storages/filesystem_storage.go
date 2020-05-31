package storages

import (
	"encoding/json"
	cfgmodels "github.com/gonamore/fxbd/config/models"
	accmodels "github.com/gonamore/fxbd/models"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type FilesystemStorage struct {
	applicationConfig *cfgmodels.ApplicationConfig
	accountConfig     *accmodels.AccountConfig

	Storage
}

func NewFilesystemStorage(applicationConfig *cfgmodels.ApplicationConfig, accountConfig *accmodels.AccountConfig) *FilesystemStorage {
	return &FilesystemStorage{applicationConfig: applicationConfig, accountConfig: accountConfig}
}

func (rcv *FilesystemStorage) Save(accountStats accmodels.AccountStats) error {
	accountStatsJSON, err := json.Marshal(accountStats)
	if err != nil {
		return err
	}
	log.Println(string(accountStatsJSON))

	filepath := path.Join(rcv.applicationConfig.StatsDir, rcv.accountConfig.Name+".json")

	err = os.MkdirAll(rcv.applicationConfig.StatsDir, os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath, accountStatsJSON, 0644)
	if err != nil {
		return err
	}
	return nil
}
