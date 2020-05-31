package providers

import "github.com/gonamore/fxbd/account/models"

type Provider interface {
	Get(url string) models.AccountStats
}
