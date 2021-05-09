package main

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

func isPhotoInURL(url string) bool {
	if strings.Contains(url, ".png") {
		return true
	}
	if strings.Contains(url, ".jpg") {
		return true
	}
	return false
}

func isRemoteFileExists(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		return false
	}
	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func getRemoteFileBytes(url string) ([]byte, error) {
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New("Received non 200 response code")
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
