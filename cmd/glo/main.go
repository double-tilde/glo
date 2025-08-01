package main

import (
	"log/slog"
	"os"

	"github.com/double-tilde/glo/pkg/config"
	"github.com/double-tilde/glo/pkg/fs"
	"github.com/double-tilde/glo/pkg/logger"
)

func main() {
	homeDir, err := fs.GetHomeDir()
	if err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}

	dataHome, err := fs.GetDataHomeDir(homeDir)
	if err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}

	configHome, err := fs.GetUserConfigHomeDir(homeDir)
	if err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}

	if err := logger.Setup(homeDir, dataHome, config.LogFileName); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}

	if err := config.Setup(configHome); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}

	cfg := config.New()

	dirs, err := fs.FindGitDirs(homeDir, cfg)
	if err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}

	// TODO: remove
	for _, dir := range dirs {
		slog.Info(dir)
	}

	// commits, err := data.CollectCommits(dirs)
	// if err != nil {
	// 	log.Println(err)
	// }
	//
	// err = data.WriteJSONFile(commits, dataHome)
	// if err != nil {
	// 	log.Println(err)
	// }
}
