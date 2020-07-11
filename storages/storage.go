package storages

import "github.com/gonamore/fxbd/account/models"

// Interface for storing account statistics
type Storage interface {
	Save(stats models.AccountStats) error
}
