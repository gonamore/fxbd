package providers

import "github.com/gonamore/fxbd/account/models"

// Interface for obtaining account statistics
type Provider interface {
	Get(url string) models.AccountStats
}
