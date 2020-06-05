package workers

import (
	cfgmodels "github.com/gonamore/fxbd/config/models"
	"github.com/gonamore/fxbd/reporters"
	"log"
	"time"
)

// Periodically regenerates reports
type ReporterWorker struct {
	applicationConfig *cfgmodels.ApplicationConfig
}

func NewReporterWorker(applicationConfig *cfgmodels.ApplicationConfig) *ReporterWorker {
	return &ReporterWorker{applicationConfig: applicationConfig}
}

func (rcv *ReporterWorker) Start() {
	for true {
		log.Println("ReporterWorker cycle started")

		reporter := reporters.NewHtmlComposer(rcv.applicationConfig)
		reporter.Assemble()

		log.Println("ReporterWorker cycle finished")

		time.Sleep(time.Duration(60) * time.Second)
	}
}
