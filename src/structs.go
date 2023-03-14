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

	FromSubreddits   []string
	IsProcessingPost bool
}

type solutionConfig struct {
	// data from config file
	Utopia utopiaConfig `ignored:"true"`
	Reddit redditConfig `ignored:"true"`
	Main   mainConfig   `ignored:"true"`
}

type mainConfig struct {
	BotNickname           string `envconfig:"BOT_NICKNAME" default:"UnboundMedia"`
	ShowSource            bool   `envconfig:"SHOW_SOURCE" default:"false"`
	MaxPostsPerQuery      int    `envconfig:"MAX_POSTS_PER_QUERY" default:"1"`
	UsePostsPerQuery      int    `envconfig:"POSTS_PER_QUERY" default:"5"`
	Cron                  string `envconfig:"CRON_SPEC" default:"@every 1h"`
	UtopiaChannelID       string `envconfig:"UTOPIA_CHANNEL_ID" required:"true"`
	UtopiaChannelPassword string `envconfig:"UTOPIA_CHANNEL_PASS"`
	FromSubredditsRaw     string `envconfig:"FROM_SUBREDDITS" required:"true"`
}

type utopiaConfig struct {
	Token        string `envconfig:"UTOPIA_TOKEN" required:"true"`
	Host         string `envconfig:"UTOPIA_HOST" default:"127.0.0.1"`
	Port         int    `envconfig:"UTOPIA_PORT" default:"20000"`
	HTTPSEnabled bool   `envconfig:"UTOPIA_USE_HTTPS" default:"false"`
}

type redditConfig struct {
	APIKeyID  string `envconfig:"REDDIT_KEY_ID" required:"true"`
	APISecret string `envconfig:"REDDIT_SECRET" required:"true"`
	User      string `envconfig:"REDDIT_USER" required:"true"`
	Password  string `envconfig:"REDDIT_PASS" required:"true"`
}
