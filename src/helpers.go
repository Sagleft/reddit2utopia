package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/Sagleft/uchatbot-engine"
	"github.com/kelseyhightower/envconfig"
	"github.com/sagleft/go-reddit/v2/reddit"
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

func redditConnect(cfg redditConfig) (*reddit.Client, error) {
	client, err := reddit.NewClient(reddit.Credentials{
		ID:       cfg.APIKeyID,
		Secret:   cfg.APISecret,
		Username: cfg.User,
		Password: cfg.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("create reddit client: %w", err)
	}
	return client, nil
}

func parseContentRoutes(routes string) contentRoutes {
	r := make(contentRoutes)

	channels := strings.Split(routes, ";")

	for i := 0; i < len(channels); i++ {
		if channels[i] == "" {
			continue
		}

		parts := strings.Split(channels[i], ":")
		subreddits := strings.Split(parts[1], ",")
		channelParts := strings.Split(parts[0], ",")
		channelID := channelParts[0]
		channelPassword := ""

		if len(channelParts) > 1 {
			channelPassword = channelParts[1]
		}

		r[channelID] = contentRoute{
			Password:   channelPassword,
			Subreddits: subreddits,
		}
	}

	return r
}

func getChats(r contentRoutes) []uchatbot.Chat {
	c := []uchatbot.Chat{}
	for channelID, data := range r {
		c = append(c, uchatbot.Chat{
			ID:       channelID,
			Password: data.Password,
		})
	}
	return c
}

func arrayShuffle(a []string) []string {
	rand.Seed(time.Now().UnixNano())

	for i := len(a) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
	return a
}
