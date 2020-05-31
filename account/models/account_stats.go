package models

type AccountStats struct {
	Balance            *float64 `json:"balance,omitempty"`
	Equity             *float64 `json:"equity,omitempty"`
	Profit             *float64 `json:"profit,omitempty"`
	DayProfitMoney     *float64 `json:"day_profit_money,omitempty"`
	DayProfitPercent   *float64 `json:"day_profit_percent,omitempty"`
	WeekProfitMoney    *float64 `json:"week_profit_money,omitempty"`
	WeekProfitPercent  *float64 `json:"week_profit_percent,omitempty"`
	MonthProfitMoney   *float64 `json:"month_profit_money,omitempty"`
	MonthProfitPercent *float64 `json:"month_profit_percent,omitempty"`
	YearProfitMoney    *float64 `json:"year_profit_money,omitempty"`
	YearProfitPercent  *float64 `json:"year_profit_percent,omitempty"`
}
