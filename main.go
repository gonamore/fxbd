package main

import (
	"github.com/gonamore/fxbd/account"
	"github.com/gonamore/fxbd/config"
	"github.com/gonamore/fxbd/webserver"
	"github.com/gonamore/fxbd/workers"
	"log"
)

func main() {
	// read config
	applicationConfigResolver := config.NewResolver()
	applicationConfig, err := applicationConfigResolver.GetConfig()
	if err != nil {
		log.Fatal("Please specify the application config file")
		return
	}
	log.Println(applicationConfig)

	// start workers
	for _, accountConfig := range applicationConfig.Accounts {
		worker := account.NewWorker(applicationConfig)
		go worker.Start(accountConfig)
	}

	// generate results
	//reporter := reporters.NewHtmlComposer(applicationConfig)
	//go reporter.Assemble()
	reporterWorker := workers.NewReporterWorker(applicationConfig)
	go reporterWorker.Start()

	// start web server
	////http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
	////	http.ServeFile(res, req, path.Join(applicationConfig.StatsDir, "my-account.json"))
	////})
	////log.Fatal(http.ListenAndServe(":8080", nil))
	//
	//server := &http.Server{Addr: ":8080", Handler: http.HandlerFunc(
	//	func(res http.ResponseWriter, req *http.Request) {
	//		http.ServeFile(res, req, path.Join(applicationConfig.StatsDir, "my-account.json"))
	//	},
	//)}
	//
	//go func() {
	//	if err := server.ListenAndServe(); err != nil {
	//		// handle err
	//	}
	//}()
	//
	//// Setting up signal capturing
	//stop := make(chan os.Signal, 1)
	//signal.Notify(stop, os.Interrupt)
	//
	//// Waiting for SIGINT (pkill -2)
	//<-stop
	//
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//if err := server.Shutdown(ctx); err != nil {
	//	// handle err
	//}
	//
	//// Wait for ListenAndServe goroutine to close.

	server := webserver.NewWebServer(applicationConfig)
	server.Start()
}
