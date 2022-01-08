package redwall

import "path/filepath"

const (
	databaseFilename string = "redwall.db"
)

type Config struct {
	Reddit  redditInfo
	Redwall redwallInfo
}

func (c *Config) DatabasePath() string {
	return filepath.Join(c.Redwall.DataDir, databaseFilename)
}

type redditInfo struct {
	UserAgent string `toml:"user_agent"`
}

type redwallInfo struct {
	DataDir         string   `toml:"data_dir"`
	SubmissionLimit int      `toml:"submission_limit"`
	TimeFilter      string   `toml:"time_filter"`
	Subreddits      []string `toml:"subreddits"`
}
