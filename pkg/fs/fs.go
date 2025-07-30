package fs

import (
	"log"
	"os"
	"path/filepath"

	"github.com/double-tilde/glo/pkg/config"
)

// TODO: Return errors, do not rely on other packages

func GetHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("No home directory found", err)
	}

	return homeDir
}

func GetDataHome(homeDir string) string {
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
		log.Fatal("Failed to create directory: ", err)
	}

	return dataHome
}

func FindGitDirs(startingDir string, dirs []string) []string {
	contents, err := os.ReadDir(startingDir)
	if err != nil {
		log.Fatal(err)
	}

outside:
	for _, content := range contents {
		for _, ignoreDir := range config.IgnoreDirs {
			if ignoreDir == content.Name() {
				// clog.Info("Ignoring " + ignoreDir)
				continue outside
			}
		}

		if content.IsDir() {
			path := filepath.Join(startingDir, content.Name())

			dir, err := os.ReadDir(path)
			if err != nil {
				log.Fatal(err)
			}

			for _, entry := range dir {
				if entry.Name() == config.GitDirectory {
					dirs = append(dirs, path)
					break
				}
			}

			dirs = FindGitDirs(path, dirs)
		}
	}

	return dirs
}
