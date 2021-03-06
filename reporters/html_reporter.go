package reporters

import (
	"encoding/json"
	"fmt"
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
			log.Println("Cannot read account stats: ", err)
		} else {
			accountData = append(accountData, AccountData{Config: accountConfig, Stats: *accountStats})
		}
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

	t := template.Must(template.New("index.html").Funcs(template.FuncMap{
		"DerefFloat": func(number *float64) float64 {
			if number != nil {
				return *number
			} else {
				return 0
			}
		},
		"ColorOf": func(number *float64) string {
			if number == nil || *number == 0.0 {
				return ""
			} else if *number > 0 {
				return "green"
			} else {
				return "red"
			}
		},
		"ValueOf": func(number *float64) string {
			if *number == 0.0 {
				return "0.00"
			} else if *number > 0 {
				return "+" + fmt.Sprintf("%.2f", *number)
			} else {
				return fmt.Sprintf("%.2f", *number)
			}
		},
		"NoSignValueOf": func(number *float64) string {
			return fmt.Sprintf("%.2f", *number)
		},
	}).ParseFiles("webserver/templates/index.html"))

	err = t.Execute(myFile, reportData)
	if err != nil {
		log.Println("Cannot create report from template: ", err)
	}
	//err = t.Execute(os.Stdout, reportData)
	//if err != nil {
	//	log.Println("Cannot console report from template: ", err)
	//}
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
