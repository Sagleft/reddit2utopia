package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

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
