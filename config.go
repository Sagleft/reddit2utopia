package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

func parseConfig(sol *solution) error {
	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		configPath = defaultConfigJSONPath
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return errors.New("failed to find config json")
	}
	jsonBytes, err := ioutil.ReadFile(defaultConfigJSONPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, &sol.Config)
}
