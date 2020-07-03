package models

type AccountStats struct {
	// general
	Balance         *float64 `json:"balance,omitempty"`
	Equity          *float64 `json:"equity,omitempty"`
	Profit          *float64 `json:"profit,omitempty"`
	Deposits        *float64 `json:"deposits,omitempty"`
	Withdrawals     *float64 `json:"withdrawals,omitempty"`
	Drawdown        *float64 `json:"drawdown,omitempty"`
	OverallDrawdown *float64 `json:"overall_drawdown,omitempty"`

	// profit by period
	DayProfitMoney           *float64 `json:"day_profit_money,omitempty"`
	DayProfitMoneyPrevious   *float64 `json:"day_profit_money_previous,omitempty"`
	DayProfitPercent         *float64 `json:"day_profit_percent,omitempty"`
	DayProfitPercentPrevious *float64 `json:"day_profit_percent_previous,omitempty"`

	WeekProfitMoney           *float64 `json:"week_profit_money,omitempty"`
	WeekProfitMoneyPrevious   *float64 `json:"week_profit_money_previous,omitempty"`
	WeekProfitPercent         *float64 `json:"week_profit_percent,omitempty"`
	WeekProfitPercentPrevious *float64 `json:"week_profit_percent_previous,omitempty"`

	MonthProfitMoney           *float64 `json:"month_profit_money,omitempty"`
	MonthProfitMoneyPrevious   *float64 `json:"month_profit_money_previous,omitempty"`
	MonthProfitPercent         *float64 `json:"month_profit_percent,omitempty"`
	MonthProfitPercentPrevious *float64 `json:"month_profit_percent_previous,omitempty"`

	YearProfitMoney           *float64 `json:"year_profit_money,omitempty"`
	YearProfitMoneyPrevious   *float64 `json:"year_profit_money_previous,omitempty"`
	YearProfitPercent         *float64 `json:"year_profit_percent,omitempty"`
	YearProfitPercentPrevious *float64 `json:"year_profit_percent_previous,omitempty"`

	// other
	SymbolStats []SymbolStats `json:"symbols,omitempty"`

	UpdateTime *string `json:"update_time,omitempty"`
}
