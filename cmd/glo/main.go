package main

import (
	"fmt"
	"log"

	"github.com/double-tilde/glo/pkg/fs"
)

// var clog *logger.Clogger

// TODO: replace logs with clogs?

func main() {
	homeDir, err := fs.GetHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dataHome, err := fs.GetDataHome(homeDir)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: remove
	fmt.Println(dataHome)

	// gloLogFile := filepath.Join(dataHome, "glo.log")
	// clog = logger.New(gloLogFile)
	// defer func() {
	// 	if err := clog.Close(); err != nil {
	// 		log.Fatal("Error closing log file")
	// 	}
	// }()

	dirs, err := fs.FindGitDirs(homeDir)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: remove
	for _, dir := range dirs {
		fmt.Println(dir)
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
