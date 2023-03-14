package main

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/kelseyhightower/envconfig"
)

func parseCronSpec(spec string) string {
	runes := []rune(spec)
	if !unicode.IsDigit(runes[0]) {
		spec = "@" + spec
	}
	return spec
}

func getRedditURL(url string) string {
	if strings.Contains(url, "http://") || strings.Contains(url, "https://") {
		return url
	}

	return redditHost + url
}

func parseConfig() (solutionConfig, error) {
	var cfg solutionConfig

	if err := envconfig.Process("", &cfg.Main); err != nil {
		return cfg, fmt.Errorf("parse main envs: %w", err)
	}

	if err := envconfig.Process("", &cfg.Reddit); err != nil {
		return cfg, fmt.Errorf("parse reddit envs: %w", err)
	}

	if err := envconfig.Process("", &cfg.Utopia); err != nil {
		return cfg, fmt.Errorf("parse utopia envs: %w", err)
	}

	return cfg, nil
}
