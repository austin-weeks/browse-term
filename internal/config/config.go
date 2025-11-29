// Package config provides configuration functionality for browse-term.
package config

import (
	"errors"
	"os"
	"path"

	"github.com/austin-weeks/browse-term/internal/themes"
	"github.com/charmbracelet/log"
	"github.com/goccy/go-yaml"
)

var configPath = path.Join(".config", "browse-term", "config.yaml")

type Config struct {
	Theme themes.Theme `yaml:"theme"`
}

func LoadConfig() Config {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Debug("Home directory not set", "error", err)
		return defaultConfig()
	}
	path := path.Join(home, configPath)
	b, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Debug("Failed to open existing config file", "error", err, "path", path)
		}
		return defaultConfig()
	}

	var userConfig Config
	err = yaml.Unmarshal(b, &userConfig)
	if err != nil {
		log.Error("Failed to parse config file, using default configuration", "error", err, "file", path)
		return defaultConfig()
	}

	return populateMissingFields(userConfig)
}

func defaultConfig() Config {
	return Config{
		Theme: themes.TokyoNight,
	}
}

// Apply default values for unset config fields
func populateMissingFields(c Config) Config {
	d := defaultConfig()
	if c.Theme == "" {
		c.Theme = d.Theme
	}
	if err := themes.ValidateTheme(c.Theme); err != nil {
		log.Error("Invalid theme, using default configuration", "error", err, "config", c)
		return d
	}
	return c
}
