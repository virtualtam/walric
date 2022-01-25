package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	databaseFilename string = "walric.db"
)

type Config struct {
	Reddit redditInfo
	Walric walricInfo
}

func (c *Config) DataDir() string {
	return c.Walric.DataDir
}

func (c *Config) DatabasePath() string {
	return filepath.Join(c.Walric.DataDir, databaseFilename)
}

type redditInfo struct {
	UserAgent string `toml:"user_agent"`
}

type walricInfo struct {
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

		configPath = filepath.Join(userHome, ".config", "walric.toml")
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
