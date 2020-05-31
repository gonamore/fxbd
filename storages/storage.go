package storages

import "github.com/gonamore/fxbd/account/models"

type Storage interface {
	Save(stats models.AccountStats) error
}
