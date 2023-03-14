package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	// given
	os.Setenv("UTOPIA_CHANNEL_ID", "test")
	os.Setenv("FROM_SUBREDDITS", "test,test")
	os.Setenv("REDDIT_KEY_ID", "test")
	os.Setenv("REDDIT_SECRET", "test")
	os.Setenv("REDDIT_USER", "test")
	os.Setenv("REDDIT_PASS", "test")
	os.Setenv("UTOPIA_TOKEN", "test")

	// when
	cfgData, err := parseConfig()

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, cfgData.Main.Cron)
	assert.NotEqual(t, 0, cfgData.Main.UsePostsPerQuery)
}
