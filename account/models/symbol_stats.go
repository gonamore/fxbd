package models

type SymbolStats struct {
	Name            string  `json:"name,omitempty"`
	Profit          float64 `json:"profit,omitempty"`
	ProfitPercent   float64 `json:"profit_percent,omitempty"`
	BuyOrdersCount  int     `json:"buy_orders_count,omitempty"`
	BuyOrdersLot    float64 `json:"buy_orders_lot,omitempty"`
	SellOrdersCount int     `json:"sell_orders_count,omitempty"`
	SellOrdersLot   float64 `json:"sell_orders_lot,omitempty"`
}
