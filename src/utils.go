package main

import (
	"errors"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
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
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, errors.New("received non 200 response code")
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetRandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func GetRandomArrString(arr []string) string {
	return arr[GetRandomInt(0, len(arr)-1)]
}
