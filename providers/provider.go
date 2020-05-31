package providers

import "github.com/gonamore/fxbd/models"

type Provider interface {
	Get(url string) models.AccountStats
}
