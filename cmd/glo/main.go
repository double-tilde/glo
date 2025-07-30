package main

import (
	"github.com/double-tilde/glo/pkg/data"
	"github.com/double-tilde/glo/pkg/fs"
)

// TODO: main calls will go here

// var clog *logger.Clogger

func main() {
	homeDir := fs.GetHomeDir()
	dataHome := fs.GetDataHome(homeDir)

	// gloLogFile := filepath.Join(dataHome, "glo.log")
	// clog = logger.New(gloLogFile)
	// defer func() {
	// 	if err := clog.Close(); err != nil {
	// 		log.Fatal("Error closing log file")
	// 	}
	// }()

	dirs := []string{}
	dirs = fs.FindGitDirs(homeDir, dirs)

	commits := data.CollectCommits(dirs)

	data.WriteJSONFile(commits, dataHome)
}
