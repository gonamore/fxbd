package reporters

import (
	"encoding/json"
	accmodels "github.com/gonamore/fxbd/account/models"
	cfgmodels "github.com/gonamore/fxbd/config/models"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type HtmlReporter struct {
	applicationConfig *cfgmodels.ApplicationConfig

	Reporter
}

type ReportData struct {
	Title    string
	Accounts []AccountData
}

type AccountData struct {
	Config accmodels.AccountConfig
	Stats  accmodels.AccountStats
}

func NewHtmlComposer(applicationConfig *cfgmodels.ApplicationConfig) *HtmlReporter {
	return &HtmlReporter{applicationConfig: applicationConfig}
}

func (rcv *HtmlReporter) Assemble() {
	accountData := make([]AccountData, 0)
	for _, accountConfig := range rcv.applicationConfig.Accounts {
		accountStats, err := rcv.readAccountStats(accountConfig.Name)
		if err != nil {
			log.Fatal("Cannot read account stats", err)
		}
		accountData = append(accountData, AccountData{Config: accountConfig, Stats: *accountStats})
	}

	reportData := ReportData{
		Title:    "Report stats",
		Accounts: accountData,
	}

	filepath := path.Join(rcv.applicationConfig.StatsDir, "index.html")
	myFile, err := os.Create(filepath)
	if err != nil {
		log.Println("Cannot create report file: ", err)
	}

	t, _ := template.ParseFiles("webserver/templates/index.html")
	//err = t.Execute(os.Stdout, reportData)
	err = t.Execute(myFile, reportData)
	if err != nil {
		log.Println("Cannot create report from template: ", err)
	}
}

func (rcv *HtmlReporter) readAccountStats(accountName string) (*accmodels.AccountStats, error) {
	filepath := path.Join(rcv.applicationConfig.StatsDir, accountName+".json")

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	stats := &accmodels.AccountStats{}
	err = json.Unmarshal(bytes, stats)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
