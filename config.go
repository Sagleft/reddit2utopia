package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type solutionConfig struct {
	// data from config file
	Utopia     utopiaConfig `json:"utopia"`
	ShowSource bool         `json:"show_source"`

	// data from args
	FromSubreddit   string
	UtopiaChannelID string
}

type utopiaConfig struct {
	Token        string `json:"token"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	HTTPSEnabled bool   `json:"enable_https"`
}

func parseConfig(sol *solution) error {
	if _, err := os.Stat(configJSONPath); os.IsNotExist(err) {
		return errors.New("failed to find config json")
	}
	jsonBytes, err := ioutil.ReadFile(configJSONPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, &sol.Config)
}
