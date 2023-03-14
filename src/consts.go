package main

import "time"

const (
	redditHost            = "https://www.reddit.com"
	defaultConfigJSONPath = "config.json"
	cacheFolderPath       = "cache"
	getSubredditPostsBy   = "day"
	subredditPostsLimit   = 24
	cacheFilename         = "cache.json"
	defaultAccountName    = "account.db"
	botNickname           = "UnboundMedia"
	botLogName            = "R2U bot"
	donateAddress         = "F50AF5410B1F3F4297043F0E046F205BCBAA76BEC70E936EB0F3AB94BF316804"
	coinTag               = "CRP"
	welcomeMessage        = "Hi. I'm just a bot that works with content"
	errConnectionMessage  = "failed to connect to Utopia Client"
)

const (
	getPostsTimeout = time.Second * 5
)

var cachedElementValue = struct{}{}
