package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

func parseConfig() (solutionConfig, error) {
	var cfg solutionConfig
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, fmt.Errorf("parse envs: %w", err)
	}
	return cfg, nil
}
