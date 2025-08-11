package main

import (
	"log/slog"
	"os"

	"github.com/double-tilde/glo/pkg/config"
	"github.com/double-tilde/glo/pkg/data"
	"github.com/double-tilde/glo/pkg/fs"
	"github.com/double-tilde/glo/pkg/logger"
	"github.com/double-tilde/glo/pkg/ui"
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

	commits, errs := data.CollectCommits(dirs)
	if errs != nil && cfg.LogMessages {
		for _, err := range errs {
			slog.Warn("does this repository have any commits?", "error", err)
		}
	}

	err = data.WriteJSONFile(commits, dataHome)
	if err != nil {
		slog.Warn("warn", "error", err)
	}

	commits, err = data.ReadJSONFile(dataHome)
	if err != nil {
		slog.Warn("warn", "error", err)
	}

	sortedCommits := data.GetYearOfCommits(commits)

	collectedDates, monthLabels, err := ui.CollectDates(sortedCommits)
	if err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}

	err = ui.Display(collectedDates, monthLabels)
	if err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}
