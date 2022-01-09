package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	databaseFilename string = "redwall.db"
)

type Config struct {
	Reddit  redditInfo
	Redwall redwallInfo
}

func (c *Config) DataDir() string {
	return c.Redwall.DataDir
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

func LoadTOML(configPath string) (*Config, error) {
	if configPath == "" {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return &Config{}, err
		}

		configPath = filepath.Join(userHome, ".config", "redwall.toml")
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return &Config{}, fmt.Errorf("config: failed to read file %q: %w", configPath, err)
	}

	config := &Config{}
	_, err = toml.Decode(string(configBytes), config)
	if err != nil {
		return &Config{}, fmt.Errorf("config: failed to decode TOML from %q: %w", configPath, err)
	}

	return config, nil
}
