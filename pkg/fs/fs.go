package fs

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/double-tilde/glo/pkg/config"
)

var dirs []string

func GetHomeDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("no home directory found")
	}

	return homeDir, nil
}

func GetDataHomeDir(homeDir string) (string, error) {
	var dataHomeDir string

	if dataHomeDir = os.Getenv("XDG_DATA_HOME"); dataHomeDir != "" {
		dataHomeDir = filepath.Join(dataHomeDir, config.GloDirectory)
	} else if dataHomeDir = os.Getenv("LOCALAPPDATA"); dataHomeDir != "" {
		dataHomeDir = filepath.Join(dataHomeDir, config.GloDirectory)
	} else if homeDir != "" {
		dataHomeDir = filepath.Join(homeDir, ".local", "share", config.GloDirectory)
	} else {
		dataHomeDir = filepath.Join(homeDir, config.GloHomeDirectory)
	}

	err := os.MkdirAll(dataHomeDir, 0755)
	if err != nil {
		return "", errors.New("failed to create data directory")
	}

	return dataHomeDir, nil
}

func GetUserConfigHomeDir(homeDir string) (string, error) {
	var configHomeDir string

	if configHomeDir = os.Getenv("XDG_CONFIG_HOME"); configHomeDir != "" {
		configHomeDir = filepath.Join(configHomeDir, config.GloDirectory)
	} else if homeDir != "" {
		configHomeDir = filepath.Join(homeDir, ".config", config.GloDirectory)
	} else {
		configHomeDir = filepath.Join(homeDir, config.GloHomeDirectory)
	}

	err := os.MkdirAll(configHomeDir, 0755)
	if err != nil {
		return "", errors.New("failed to create config directory")
	}

	return configHomeDir, nil
}

func FindGitDirs(cfg *config.Config, startingDir string) ([]string, error) {
	contents, err := os.ReadDir(startingDir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

outside:
	for _, content := range contents {
		for _, ignoredDir := range cfg.IgnoredDirs {
			if ignoredDir == content.Name() {
				if cfg.LogMessages {
					slog.Info("ignoring directory " + ignoredDir)
				}
				continue outside
			}
		}

		if content.IsDir() {
			path := filepath.Join(startingDir, content.Name())

			gitPath := filepath.Join(path, config.GitDirectory)
			if _, err := os.Stat(gitPath); err == nil {
				dirs = append(dirs, path)
				continue
			}

			subDirs, err := FindGitDirs(cfg, path)
			if err != nil {
				return nil, err
			}

			dirs = subDirs
		}
	}

	return dirs, nil
}
