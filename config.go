package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

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
