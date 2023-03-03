package main

type mediaPost struct {
	Text         string
	ImageURL     string
	IsLocalImage bool
}

type solution struct {
	Config solutionConfig
	Utopia *utopiaService
	Cache  *CacheHandler
}

type solutionConfig struct {
	// data from config file
	Utopia           utopiaConfig `json:"utopia"`
	Reddit           redditConfig `json:"reddit"`
	ShowSource       bool         `json:"show_source"`
	MaxPostsPerQuery int          `json:"posts_per_query"`
	UsePostsPerQuery int          `json:"use_posts_per_query"`

	// data from args
	FromSubreddits  []string
	UtopiaChannelID string
}

type utopiaConfig struct {
	Token        string `json:"token"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	HTTPSEnabled bool   `json:"enable_https"`
}

type redditConfig struct {
	APIKeyID  string `json:"keyID"`
	APISecret string `json:"keySecret"`
	User      string `json:"user"`
	Password  string `json:"pass"`
}
