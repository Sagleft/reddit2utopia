package main

import (
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
