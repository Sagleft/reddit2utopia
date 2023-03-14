package main

const (
	redditHost            = "https://www.reddit.com"
	defaultConfigJSONPath = "config.json"
	cacheFolderPath       = "cache"
	getSubredditPostsBy   = "day"
	subredditPostsLimit   = 24
	cacheFilename         = "cache.json"
	defaultAccountName    = "account.db"
	botNickname           = "UnboundMedia"
)

var cachedElementValue = struct{}{}
