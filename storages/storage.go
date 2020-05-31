package storages

import "github.com/gonamore/fxbd/models"

type Storage interface {
	Save(stats models.AccountStats) error
}
