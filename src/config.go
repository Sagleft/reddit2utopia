package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type solutionConfig struct {
	Utopia utopiaConfig `json:"utopia"` // data from config file

	FromSubreddit   string // data from args
	UtopiaChannelID string // data from args
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
