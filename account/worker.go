package account

import (
	"github.com/gonamore/fxbd/account/models"
	cfgmodels "github.com/gonamore/fxbd/config/models"
	"github.com/gonamore/fxbd/providers"
	"github.com/gonamore/fxbd/storages"
	"log"
	"time"
)

// Makes the important job to fetch and save account stats
type Worker struct {
	applicationConfig *cfgmodels.ApplicationConfig
}

func NewWorker(applicationConfig *cfgmodels.ApplicationConfig) *Worker {
	return &Worker{applicationConfig: applicationConfig}
}

func (rcv *Worker) Start(accountConfig models.AccountConfig) {
	for true {
		log.Println("Worker cycle started for account", accountConfig.Name)

		myfxbookProvider := providers.NewMyfxbookProvider()
		accountStats := myfxbookProvider.Get(accountConfig)

		storage := storages.NewFilesystemStorage(rcv.applicationConfig, &accountConfig)
		err := storage.Save(accountStats)
		if err != nil {
			log.Fatal("Cannot save stats", err)
			return
		}
		log.Println("Worker cycle finished for account", accountConfig.Name)

		time.Sleep(time.Duration(accountConfig.RefreshSeconds) * time.Second)
	}
}
