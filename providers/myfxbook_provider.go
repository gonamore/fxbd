package providers

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/gonamore/fxbd/account/models"
	"io/ioutil"
	"log"
	"math"
	"net/http"
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
	accountStats := models.AccountStats{}

	//var currencySymbol string
	//switch accountConfig.Currency {
	//case "USD":
	//	currencySymbol = "$"
	//default:
	//	currencySymbol = ""
	//}

	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{accountConfig.Location},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			//fmt.Println(string(r.Body))

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

				deposits := rcv.deposits(s)
				if deposits != nil {
					accountStats.Deposits = rcv.normalizeCurrency(*deposits, accountConfig.CurrencyDivider)
				}

				withdrawals := rcv.withdrawals(s)
				if withdrawals != nil {
					accountStats.Withdrawals = rcv.normalizeCurrency(*withdrawals, accountConfig.CurrencyDivider)
				}
			})

			drawdown := rcv.drawdown(accountStats.Balance, accountStats.Equity)
			if drawdown != nil {
				accountStats.Drawdown = drawdown
			}

			overallDrawdown := rcv.overallDrawdown(accountStats.Deposits, accountStats.Withdrawals, accountStats.Equity)
			if overallDrawdown != nil {
				accountStats.OverallDrawdown = overallDrawdown
			}

			r.HTMLDoc.Find("tr").Each(func(_ int, s *goquery.Selection) {
				dayProfitMoney, dayProfitPercent, err := rcv.profitPeriod(s, "Today", "td")
				if err != nil {
					log.Println("Cannot fetch day profit", err)
				} else if dayProfitMoney != nil && dayProfitPercent != nil {
					//log.Println(*dayProfitMoney)
					//log.Println(*dayProfitPercent)
					accountStats.DayProfitMoney = rcv.normalizeCurrency(*dayProfitMoney, accountConfig.CurrencyDivider)
					accountStats.DayProfitPercent = dayProfitPercent
				}

				weekProfitMoney, weekProfitPercent, err := rcv.profitPeriod(s, "This Week", "td")
				if err != nil {
					log.Println("Cannot fetch week profit", err)
				} else if weekProfitMoney != nil && weekProfitPercent != nil {
					//log.Println(*weekProfitMoney)
					//log.Println(*weekProfitPercent)
					accountStats.WeekProfitMoney = rcv.normalizeCurrency(*weekProfitMoney, accountConfig.CurrencyDivider)
					accountStats.WeekProfitPercent = weekProfitPercent
				}

				monthProfitMoney, monthProfitPercent, err := rcv.profitPeriod(s, "This Month", "td")
				if err != nil {
					log.Println("Cannot fetch month profit", err)
				} else if monthProfitMoney != nil && monthProfitPercent != nil {
					//log.Println(*monthProfitMoney)
					//log.Println(*monthProfitPercent)
					accountStats.MonthProfitMoney = rcv.normalizeCurrency(*monthProfitMoney, accountConfig.CurrencyDivider)
					accountStats.MonthProfitPercent = monthProfitPercent
				}

				yearProfitMoney, yearProfitPercent, err := rcv.profitPeriod(s, "This Year", "td")
				if err != nil {
					log.Println("Cannot fetch year profit", err)
				} else if yearProfitMoney != nil && yearProfitPercent != nil {
					//log.Println(*yearProfitMoney)
					//log.Println(*yearProfitPercent)
					accountStats.YearProfitMoney = rcv.normalizeCurrency(*yearProfitMoney, accountConfig.CurrencyDivider)
					accountStats.YearProfitPercent = yearProfitPercent
				}
			})

			symbolStats := rcv.symbolStats(accountConfig, r.HTMLDoc)
			if symbolStats != nil {
				accountStats.SymbolStats = symbolStats
			}
		},
	}).Start()

	return accountStats
}

func (rcv *MyfxbookProvider) balance(s *goquery.Selection) *float64 {
	rawValue := rcv.rawValue(s, "Balance", "span.floatLeft", "span.floatNone")
	if rawValue == nil {
		return nil
	}
	//log.Println("Balance raw:", *rawValue)

	result, err := rcv.numericValue(*rawValue)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Balance:", *result)
	return result
}

func (rcv *MyfxbookProvider) equity(s *goquery.Selection) *float64 {
	rawValue := rcv.rawValue(s, "Equity", "span.floatLeft", "span.floatNone")
	if rawValue == nil {
		return nil
	}
	//log.Println("Equity raw:", *rawValue)

	trimmed := strings.TrimSpace(*rawValue)
	//log.Println("Equity trimmed:", trimmed)

	index := strings.LastIndex(trimmed, " ")
	if index == -1 {
		log.Println("Equity has not expected value:", trimmed)
		return nil
	}
	splitted := trimmed[index+1:]
	//log.Println("Equity splitted:", splitted)

	result, err := rcv.numericValue(splitted)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Equity:", *result)
	return result
}

func (rcv *MyfxbookProvider) profit(s *goquery.Selection) *float64 {
	rawValue := rcv.rawValue(s, "Profit", "span.floatLeft", "span.floatNone")
	if rawValue == nil {
		return nil
	}
	//log.Println("Profit raw:", *rawValue)

	result, err := rcv.numericValue(*rawValue)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Profit:", *result)
	return result
}

func (rcv *MyfxbookProvider) deposits(s *goquery.Selection) *float64 {
	rawValue := rcv.rawValue(s, "Deposits", "span.floatLeft", "span.floatNone")
	if rawValue == nil {
		return nil
	}
	//log.Println("Profit raw:", *rawValue)

	result, err := rcv.numericValue(*rawValue)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Deposits:", *result)
	return result
}

func (rcv *MyfxbookProvider) withdrawals(s *goquery.Selection) *float64 {
	rawValue := rcv.rawValue(s, "Withdrawals", "span.floatLeft", "span.floatNone")
	if rawValue == nil {
		return nil
	}
	//log.Println("Profit raw:", *rawValue)

	result, err := rcv.numericValue(*rawValue)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Withdrawals:", *result)
	return result
}

func (rcv *MyfxbookProvider) drawdown(balance *float64, equity *float64) *float64 {
	if balance != nil && equity != nil {
		// balance -> 100%
		// equity -> x
		drawdown := (100 - (100 * *equity / *balance)) * -1.0

		rounded := math.Round(drawdown*100) / 100
		return &rounded
	}
	return nil
}

func (rcv *MyfxbookProvider) overallDrawdown(deposits *float64, withdrawals *float64, equity *float64) *float64 {
	if deposits != nil && withdrawals != nil && equity != nil {
		adjustedDeposit := *deposits - *withdrawals
		if adjustedDeposit == 0 {
			return nil
		}

		// deposit -> 100%
		// equity -> x
		drawdown := (100 - (100 * *equity / adjustedDeposit)) * -1.0

		rounded := math.Round(drawdown*100) / 100
		return &rounded
	}
	return nil
}

func (rcv *MyfxbookProvider) drawdown2(s *goquery.Selection, equity *float64) *float64 {
	depositsRawValue := rcv.rawValue(s, "Deposits", "span.floatLeft", "span.floatNone")
	if depositsRawValue == nil {
		return nil
	}
	//log.Println("Profit raw:", *rawValue)

	deposits, err := rcv.numericValue(*depositsRawValue)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Deposits:", *deposits)

	withdrawalsRawValue := rcv.rawValue(s, "Withdrawals", "span.floatLeft", "span.floatNone")
	if withdrawalsRawValue == nil {
		return nil
	}
	//log.Println("Profit raw:", *rawValue)

	withdrawals, err := rcv.numericValue(*withdrawalsRawValue)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Withdrawals:", *withdrawals)

	//equity := rcv.equity(s)
	if equity == nil {
		return nil
	}

	adjustedDeposit := *deposits - *withdrawals
	drawdown := 100 * *equity / adjustedDeposit

	rounded := math.Round(drawdown*100) / 100
	return &rounded
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

	if normalized == "" {
		return nil, nil
	}

	result, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return nil, err
	}
	//log.Println("Parsed for numeric:", result)
	return &result, nil
}

func (rcv *MyfxbookProvider) normalizeCurrency(value float64, divider int64) *float64 {
	if divider == 0 || divider == 1 {
		return &value
	}
	divided := value / float64(divider)
	rounded := math.Round(divided*100) / 100
	return &rounded
}

func (rcv *MyfxbookProvider) profitPeriod(s *goquery.Selection, nameMarker string, nameSelector string) (*float64, *float64, error) {
	selection := s.Find(nameSelector)
	name := selection.Text()
	if strings.Contains(name, nameMarker) {
		//log.Println(name)

		profitMoneyAsString := selection.Next().Next().Find("span").First().Text()
		//log.Println(profitMoneyAsString)
		profitMoney, err := rcv.numericValue(profitMoneyAsString)
		if err != nil {
			return nil, nil, err
		}

		profitPercentAsString := selection.Next().Find("span").First().Text()
		//log.Println(profitPercentAsString)
		profitPercent, err := rcv.numericValue(profitPercentAsString)
		if err != nil {
			return nil, nil, err
		}

		return profitMoney, profitPercent, nil
	}
	return nil, nil, nil
}

func (rcv *MyfxbookProvider) symbolStats(accountConfig models.AccountConfig, doc *goquery.Document) []models.SymbolStats {
	pageCount := 1
	doc.Find("#openTrades .paging").Each(func(_ int, s *goquery.Selection) {
		page := s.Text()
		//log.Println(page)

		if pageAsInt, err := strconv.Atoi(page); err == nil {
			pageCount = pageAsInt
		}
	})
	log.Println("page count:", pageCount)

	split := strings.Split(strings.TrimRight(accountConfig.Location, "/"), "/")
	accountId := split[len(split)-1]

	for pageIndex := 1; pageIndex <= pageCount; pageIndex++ {
		body, err := rcv.fetchOpenTradesForPage(pageIndex, accountId)
		if err != nil {
			log.Println("err:", err)
		}
		log.Println("body:", *body)

	}

	return nil
}

//curl 'https://www.myfxbook.com/paging.html?pt=15&p=1&ts=20000&l=x&id=5923181'    \
//	-H 'authority: www.myfxbook.com'    \
//	-H 'accept: */*'    \
//	-H 'user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36 Edg/83.0.478.50'    \
//	-H 'x-requested-with: XMLHttpRequest'    \
//	-H 'sec-fetch-site: same-origin'    \
//	-H 'sec-fetch-mode: cors'    \
//	-H 'sec-fetch-dest: empty'    \
//	-H 'referer: https://www.myfxbook.com/'    \
//	-H 'accept-language: uk,en;q=0.9,en-GB;q=0.8,en-US;q=0.7'    \
//	--compressed
func (rcv *MyfxbookProvider) fetchOpenTradesForPage(page int, accountId string) (*string, error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.myfxbook.com/paging.html?pt=15&p="+string(page)+"&ts=20000&l=x&id="+accountId, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("authority", "www.myfxbook.com")
	req.Header.Add("accept", "*/*")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36 Edg/83.0.478.50")
	req.Header.Add("x-requested-with", "XMLHttpRequest")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", "https://www.myfxbook.com/")
	req.Header.Add("accept-language", "uk,en;q=0.9,en-GB;q=0.8,en-US;q=0.7")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("not 200")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := string(body)
	return &result, nil
}

func (rcv *MyfxbookProvider) fetchSymbolStatsForPage() ([]models.SymbolStats, error) {

	return nil, nil
}
