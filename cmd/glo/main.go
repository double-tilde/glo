package main

import (
	"log/slog"
	"os"

	"github.com/double-tilde/glo/pkg/config"
	"github.com/double-tilde/glo/pkg/data"
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

	dirs, err := fs.FindGitDirs(cfg, homeDir)
	if err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}

	commits, logs := data.CollectCommits(dirs)
	if len(logs) > 0 && cfg.LogMessages {
		for _, log := range logs {
			slog.Warn("does this repository have any commits?", "error", log)
		}
	}

	err = data.WriteJSONFile(commits, dataHome)
	if err != nil {
		slog.Warn("warn", "error", err)
	}
}
