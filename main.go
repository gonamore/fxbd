package main

import (
	"github.com/gonamore/fxbd/config"
	"github.com/gonamore/fxbd/providers"
	"github.com/gonamore/fxbd/storages"
	"log"
)

func main() {
	println("start")

	applicationConfigResolver := config.NewResolver()
	applicationConfig, err := applicationConfigResolver.GetConfig()
	if err != nil {
		log.Fatal("Please specify the application config file")
		return
	}
	log.Println(applicationConfig)

	myfxbookProvider := providers.NewMyfxbookProvider()
	accountStats := myfxbookProvider.Get(applicationConfig.Accounts[0])

	storage := storages.NewFilesystemStorage(applicationConfig, &applicationConfig.Accounts[0])
	err = storage.Save(accountStats)
	if err != nil {
		log.Fatal("Cannot save stats", err)
		return
	}

	println("finish")
}
