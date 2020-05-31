package models

import "github.com/gonamore/fxbd/models"

// Describes all possible application configuration options
type ApplicationConfig struct {
	// Where statistics is saved
	StatsDir string `json:"stats_dir",yaml:"stats-dir",envconfig:"STATS_DIR"`

	Accounts []models.AccountConfig `json:"accounts",yaml:"accounts"`
}
