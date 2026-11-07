package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/unix1/blug/internal/site"
	"go.yaml.in/yaml/v4"
)

type Config struct {
	SiteHeader      string `yaml:"site_header"`
	SiteTitle       string `yaml:"site_title"`
	SiteSubtitle    string `yaml:"site_subtitle"`
	SiteDescription string `yaml:"site_description"`
	SiteFooter      string `yaml:"site_footer"`
	ListingHeading  string `yaml:"listing_heading"`
	ListingIntro    string `yaml:"listing_intro"`
}

func Load(siteRoot string) (Config, error) {
	path := filepath.Join(siteRoot, site.ConfigFile)
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, fmt.Errorf("%s not found in %s (run blug init)", site.ConfigFile, siteRoot)
		}
		return Config{}, err
	}
	var cfg Config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return Config{}, fmt.Errorf("%s: %w", site.ConfigFile, err)
	}
	return cfg, nil
}
