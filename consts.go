package main

const (
	redditHost          = "https://www.reddit.com"
	configJSONPath      = "config.json"
	cacheFolderPath     = "cache"
	getSubredditPostsBy = "day"
	subredditPostsLimit = 24
	cacheFilename       = "cache.json"
)

var cachedElementValue = struct{}{}
