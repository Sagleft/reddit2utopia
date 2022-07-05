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
	Utopia        utopiaConfig `json:"utopia"`
	ShowSource    bool         `json:"show_source"`
	PostsPerQuery int          `json:"posts_per_query"`

	// data from args
	FromSubreddit   string
	UtopiaChannelID string
}

type utopiaConfig struct {
	Token        string `json:"token"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	HTTPSEnabled bool   `json:"enable_https"`
}
