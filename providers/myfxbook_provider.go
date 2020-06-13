package providers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/gonamore/fxbd/account/models"
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

	var currencySymbol string
	switch accountConfig.Currency {
	case "USD":
		currencySymbol = "$"
	default:
		currencySymbol = ""
	}

	//go geziyor.NewGeziyor(&geziyor.Options{
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{accountConfig.Location},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			//fmt.Println(string(r.Body))
			//defer wg.Done()

			r.HTMLDoc.Find("li").Each(func(_ int, s *goquery.Selection) {
				balance := rcv.balance(s, currencySymbol)
				if balance != nil {
					accountStats.Balance = rcv.normalizeCurrency(*balance, accountConfig.CurrencyDivider)
				}

				equity := rcv.equity(s, currencySymbol)
				if equity != nil {
					accountStats.Equity = rcv.normalizeCurrency(*equity, accountConfig.CurrencyDivider)
				}

				profit := rcv.profit(s, currencySymbol)
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

func (rcv *MyfxbookProvider) balance(s *goquery.Selection, currencySymbol string) *float64 {
	rawValue := rcv.rawValue(s, "Balance", "span.floatLeft", "span.floatNone")
	if rawValue == nil {
		log.Println("Balance not fetched")
		return nil
	}
	log.Println("Balance raw:", *rawValue)

	//trimmed := strings.TrimSpace(*rawValue)
	//log.Println("Balance trimmed:", trimmed)
	//
	//regex, _ := regexp.Compile("[^\\-\\w.]")
	//normalized := regex.ReplaceAllString(trimmed, "")
	//log.Println("Balance normalized:", normalized)
	//
	//result, err := strconv.ParseFloat(normalized, 64)
	//if err != nil {
	//	log.Println(err)
	//	return nil
	//}
	//log.Println("Balance parsed:", result)
	//return &result

	result, err := rcv.numericValue(*rawValue)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Balance parsed:", result)
	return result
}

func (rcv *MyfxbookProvider) equity(s *goquery.Selection, currencySymbol string) *float64 {
	rawValue := rcv.rawValue(s, "Equity", "span.floatLeft", "span.floatNone")
	if rawValue == nil {
		return nil
	}
	log.Println("Equity raw:", *rawValue)

	trimmed := strings.TrimSpace(*rawValue)
	log.Println("Equity trimmed:", trimmed)

	index := strings.LastIndex(trimmed, " ")
	if index == -1 {
		log.Println("Equity has not expected value:", trimmed)
		return nil
	}
	splitted := trimmed[index+1:]
	//splitted := strings.Split(trimmed, " ")
	log.Println("Equity splitted:", splitted)

	regex, _ := regexp.Compile("[^\\-\\w.]")
	normalized := regex.ReplaceAllString(splitted, "")
	log.Println("Equity normalized:", normalized)

	result, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Equity parsed:", result)
	return &result
}

func (rcv *MyfxbookProvider) profit(s *goquery.Selection, currencySymbol string) *float64 {
	rawValue := rcv.rawValue(s, "Profit", "span.floatLeft", "span.floatNone")
	if rawValue == nil {
		return nil
	}
	log.Println("Profit raw:", *rawValue)

	trimmed := strings.TrimSpace(*rawValue)
	log.Println("Profit trimmed:", trimmed)

	regex, _ := regexp.Compile("[^\\-\\w.]")
	normalized := regex.ReplaceAllString(trimmed, "")
	log.Println("Profit normalized:", normalized)

	result, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Profit parsed:", result)
	return &result
}

func (rcv *MyfxbookProvider) rawValue(s *goquery.Selection, nameMarker string, nameSelector string, valueSelector string) *string {
	name := s.Find(nameSelector).Text()
	if strings.Contains(name, nameMarker) {
		result := s.Find(valueSelector).Text()
		//log.Println(nameMarker, result)
		return &result
	}
	return nil
}

func (rcv *MyfxbookProvider) numericValue(rawValue string) (*float64, error) {
	regex, _ := regexp.Compile("[^\\-\\w.]")
	normalized := regex.ReplaceAllString(rawValue, "")
	//log.Println("Normalized for numeric:", normalized)

	result, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return nil, err
	}
	//log.Println("Parsed for numeric:", result)
	return &result, nil
}

func (rcv *MyfxbookProvider) moneyValue(s *goquery.Selection, currencySymbol string, nameSelector string, valueSelector string, nameMarker string) *float64 {
	name := s.Find(nameSelector).Text()
	if strings.Contains(name, nameMarker) {
		rawValue := s.Find(valueSelector).Text()
		log.Println(nameMarker, rawValue)

		//regex, _ := regexp.Compile("[^\\w."+currencySymbol+"]")
		regex, _ := regexp.Compile("[^\\w.]")
		value := regex.ReplaceAllString(rawValue, "")

		//index := strings.LastIndex(rawValue, "$")
		//if index != -1 {
		//	rawValue = rawValue[index+1:]
		//}

		//value := strings.TrimSpace(rawValue)
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
