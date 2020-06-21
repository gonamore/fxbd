package models

type SymbolStats struct {
	Name          string   `json:"name,omitempty"`
	Profit        *float64 `json:"profit,omitempty"`
	ProfitPercent *float64 `json:"profit_percent,omitempty"`
}
