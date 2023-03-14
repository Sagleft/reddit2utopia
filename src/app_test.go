package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	cfgData, err := parseConfig()
	require.NoError(t, err)

	assert.NotEmpty(t, cfgData.Main.Cron)
	assert.NotEqual(t, 0, cfgData.Main.UsePostsPerQuery)
}
