package providers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/gonamore/fxbd/models"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type MyfxbookProvider struct {
	Provider
}

func NewMyfxbookProvider() *MyfxbookProvider {
	return &MyfxbookProvider{}
}

func (rcv *MyfxbookProvider) Get(accountConfig models.AccountConfig) models.AccountStats {
	//var wg sync.WaitGroup
	//wg.Add(1)
	accountStats := models.AccountStats{}

	//go geziyor.NewGeziyor(&geziyor.Options{
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{accountConfig.Location},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			//fmt.Println(string(r.Body))
			//defer wg.Done()

			r.HTMLDoc.Find("li").Each(func(_ int, s *goquery.Selection) {
				balance := rcv.balance(s)
				if balance != nil {
					accountStats.Balance = rcv.normalizeCurrency(*balance, accountConfig.CurrencyDivider)
				}

				equity := rcv.equity(s)
				if equity != nil {
					accountStats.Equity = rcv.normalizeCurrency(*equity, accountConfig.CurrencyDivider)
				}

				profit := rcv.profit(s)
				if profit != nil {
					accountStats.Profit = rcv.normalizeCurrency(*profit, accountConfig.CurrencyDivider)
				}
			})

			r.HTMLDoc.Find("tr").Each(func(_ int, s *goquery.Selection) {
				dayProfitMoney, dayProfitPercent := rcv.profitPeriod(s, "td", "This Week")
				if dayProfitMoney != nil && dayProfitPercent != nil {
					//log.Println(*dayProfitMoney)
					//log.Println(*dayProfitPercent)
					accountStats.DayProfitMoney = rcv.normalizeCurrency(*dayProfitMoney, accountConfig.CurrencyDivider)
					accountStats.DayProfitPercent = dayProfitPercent
				}

				weekProfitMoney, weekProfitPercent := rcv.profitPeriod(s, "td", "This Week")
				if weekProfitMoney != nil && weekProfitPercent != nil {
					//log.Println(*weekProfitMoney)
					//log.Println(*weekProfitPercent)
					accountStats.WeekProfitMoney = rcv.normalizeCurrency(*weekProfitMoney, accountConfig.CurrencyDivider)
					accountStats.WeekProfitPercent = weekProfitPercent
				}

				monthProfitMoney, monthProfitPercent := rcv.profitPeriod(s, "td", "This Month")
				if monthProfitMoney != nil && monthProfitPercent != nil {
					//log.Println(*monthProfitMoney)
					//log.Println(*monthProfitPercent)
					accountStats.MonthProfitMoney = rcv.normalizeCurrency(*monthProfitMoney, accountConfig.CurrencyDivider)
					accountStats.MonthProfitPercent = monthProfitPercent
				}

				yearProfitMoney, yearProfitPercent := rcv.profitPeriod(s, "td", "This Year")
				if yearProfitMoney != nil && yearProfitPercent != nil {
					//log.Println(*yearProfitMoney)
					//log.Println(*yearProfitPercent)
					accountStats.YearProfitMoney = rcv.normalizeCurrency(*yearProfitMoney, accountConfig.CurrencyDivider)
					accountStats.YearProfitPercent = yearProfitPercent
				}
			})
		},
	}).Start()

	//wg.Wait()
	//log.Println(accountStats)
	return accountStats
}

//func (rcv *MyfxbookProvider) strToFloat(name string, value string, nameMarker string) *float64 {
//	regex, err := regexp.Compile("[^0-9.]+")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if strings.Contains(name, nameMarker) {
//		result, err := strconv.ParseFloat(regex.ReplaceAllString(value, ""), 64)
//		if err != nil {
//			log.Fatal(err)
//		} else {
//			return &result
//		}
//	}
//	return nil
//}

func (rcv *MyfxbookProvider) balance(s *goquery.Selection) *float64 {
	//name := s.Find("span.floatLeft").Text()
	//marker := "Balance"
	//if strings.Contains(name, marker) {
	//	raw := s.Find("span.floatNone").Text()
	//	log.Println(marker, raw)
	//	value := strings.TrimLeft(strings.TrimSpace(raw), "$")
	//	result, err := strconv.ParseFloat(value, 64)
	//	if err != nil {
	//		log.Fatal(err)
	//	} else {
	//		return &result
	//	}
	//}
	return rcv.moneyValue(s, "span.floatLeft", "span.floatNone", "Balance")
}

func (rcv *MyfxbookProvider) equity(s *goquery.Selection) *float64 {
	//name := s.Find("span.floatLeft").Text()
	//marker := "Equity"
	//if strings.Contains(name, marker) {
	//	raw := s.Find("span.floatNone").Text()
	//	log.Println(marker, raw)
	//	value := strings.TrimLeft(strings.TrimSpace(raw), "$")
	//	result, err := strconv.ParseFloat(value, 64)
	//	if err != nil {
	//		log.Fatal(err)
	//	} else {
	//		return &result
	//	}
	//}
	//return nil
	return rcv.moneyValue(s, "span.floatLeft", "span.floatNone", "Equity")
}

func (rcv *MyfxbookProvider) profit(s *goquery.Selection) *float64 {
	return rcv.moneyValue(s, "span.floatLeft", "span.floatNone", "Profit")
}

func (rcv *MyfxbookProvider) moneyValue(s *goquery.Selection, nameSelector string, valueSelector string, nameMarker string) *float64 {
	name := s.Find(nameSelector).Text()
	if strings.Contains(name, nameMarker) {
		rawValue := s.Find(valueSelector).Text()
		log.Println(nameMarker, rawValue)

		index := strings.LastIndex(rawValue, "$")
		if index != -1 {
			rawValue = rawValue[index+1:]
		}

		value := strings.TrimSpace(rawValue)
		result, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Println(err)
		} else {
			return &result
		}
	}
	return nil
}

func (rcv *MyfxbookProvider) normalizeCurrency(value float64, divider int64) *float64 {
	if divider == 0 || divider == 1 {
		return &value
	}
	divided := value / float64(divider)
	rounded := math.Round(divided*100) / 100
	return &rounded
}

func (rcv *MyfxbookProvider) profitPeriod(s *goquery.Selection, nameSelector string, nameMarker string) (*float64, *float64) {
	selection := s.Find(nameSelector)
	name := selection.Text()
	if strings.Contains(name, nameMarker) {
		//log.Println(name)

		profitPercentAsString := selection.Next().Find("span").First().Text()
		profitMoneyAsString := selection.Next().Next().Find("span").First().Text()
		//log.Println(profitPercentAsString)
		//log.Println(profitMoneyAsString)

		return rcv.strToFloat(profitMoneyAsString), rcv.strToFloat(profitPercentAsString)
	}
	return nil, nil
}

func (rcv *MyfxbookProvider) strToFloat(value string) *float64 {
	regex, err := regexp.Compile("[^+\\-0-9.]+")
	if err != nil {
		log.Fatal(err)
	}
	result, err := strconv.ParseFloat(regex.ReplaceAllString(value, ""), 64)
	if err != nil {
		log.Println(err)
	} else {
		return &result
	}
	return nil
}
