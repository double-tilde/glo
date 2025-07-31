package fs

import (
	"errors"
	"fmt"
	"log"
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

func GetDataHome(homeDir string) (string, error) {
	var dataHome string

	if dataHome = os.Getenv("XDG_DATA_HOME"); dataHome != "" {
		dataHome = filepath.Join(dataHome, config.GloDirectory)
	} else if dataHome = os.Getenv("LOCALAPPDATA"); dataHome != "" {
		dataHome = filepath.Join(dataHome, config.GloDirectory)
	} else if homeDir != "" {
		dataHome = filepath.Join(homeDir, ".local", "share", config.GloDirectory)
	} else {
		dataHome = filepath.Join(homeDir, "."+config.GloDirectory)
	}

	err := os.MkdirAll(dataHome, 0755)
	if err != nil {
		return "", errors.New("failed to create data directory")
	}

	return dataHome, nil
}

func FindGitDirs(startingDir string) ([]string, error) {
	contents, err := os.ReadDir(startingDir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

outside:
	for _, content := range contents {
		for _, ignoreDir := range config.IgnoreDirs {
			if ignoreDir == content.Name() {
				if config.LogIgnoreDirs {
					log.Println("ignoring directory", ignoreDir)
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

			subDirs, err := FindGitDirs(path)
			if err != nil {
				return nil, err
			}

			dirs = subDirs
		}
	}

	return dirs, nil
}
