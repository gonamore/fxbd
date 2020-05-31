package main

import (
	"fmt"
	"github.com/gonamore/fxbd/account"
	"github.com/gonamore/fxbd/config"
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

	for _, accountConfig := range applicationConfig.Accounts {
		worker := account.NewWorker(applicationConfig)
		go worker.Start(accountConfig)
	}

	fmt.Print("Enter text: \n")
	var input string
	fmt.Scanln(&input)
	fmt.Print(input)

	println("finish")
}
