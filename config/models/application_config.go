package models

import "github.com/gonamore/fxbd/account/models"

// Describes all possible application configuration options
type ApplicationConfig struct {
	// Port for output statistics result
	Port int `json:"port",yaml:"port",envconfig:"PORT"`

	// Where statistics is saved
	StatsDir string `json:"stats_dir",yaml:"stats-dir",envconfig:"STATS_DIR"`

	Accounts []models.AccountConfig `json:"accounts",yaml:"accounts"`
}
