package models

type AccountConfig struct {
	Name            string `json:"name,omitempty"`
	Location        string `json:"location,omitempty"`
	RefreshSeconds  int64  `json:"refresh_seconds,omitempty"`
	CurrencyDivider int64  `json:"currency_divider,omitempty"`
}
