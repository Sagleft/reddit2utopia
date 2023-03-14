package main

import (
	"strings"
	"unicode"
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
