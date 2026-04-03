package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

const defaultSnippetDirRelative = "~/.local/share/snx/snippets"

// Config holds the SNX configuration.
type Config struct {
	SnippetDir string `toml:"snippet_dir"`
}

// Load reads the SNX config from ~/.config/snx/config.toml.
// If the file does not exist, a default config is returned with no error.
func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not determine home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".config", "snx", "config.toml")

	cfg := &Config{
		SnippetDir: defaultSnippetDirRelative,
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// No config file — use defaults.
			cfg.SnippetDir = expandTilde(defaultSnippetDirRelative, homeDir)
			return cfg, nil
		}
		return nil, fmt.Errorf("could not read config file %s: %w", configPath, err)
	}

	if _, err := toml.Decode(string(data), cfg); err != nil {
		return nil, fmt.Errorf("could not parse config file %s: %w", configPath, err)
	}

	cfg.SnippetDir = expandTilde(cfg.SnippetDir, homeDir)
	return cfg, nil
}

// expandTilde replaces a leading ~ with the user's home directory.
func expandTilde(path, homeDir string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}
	if path == "~" {
		return homeDir
	}
	return path
}
