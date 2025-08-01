package main

import (
	"errors"
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

	dataHome, err := fs.GetDataHome(homeDir)
	if err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}

	if err := logger.Init(homeDir, dataHome, config.LogFileName); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}

	dirs, err := fs.FindGitDirs(homeDir)
	if err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}

	// TODO: remove
	slog.Info(dataHome)

	// TODO: remove
	for _, dir := range dirs {
		slog.Info(dir)
	}

	err = errors.New("the bad thing")
	slog.Error("something went wrong", "error", err)
	os.Exit(1)

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
