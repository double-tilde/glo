package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/double-tilde/glo/clogger"
)

func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("No home directory found", err)
	}

	return homeDir
}

func getDataHome(homeDir string) string {
	var dataHome string

	if dataHome = os.Getenv("XDG_DATA_HOME"); dataHome != "" {
		dataHome = filepath.Join(dataHome, gloDirectory)
	} else if dataHome = os.Getenv("LOCALAPPDATA"); dataHome != "" {
		dataHome = filepath.Join(dataHome, gloDirectory)
	} else if homeDir != "" {
		dataHome = filepath.Join(homeDir, ".local", "share", gloDirectory)
	} else {
		dataHome = filepath.Join(homeDir, "."+gloDirectory)
	}

	err := os.MkdirAll(dataHome, 0755)
	if err != nil {
		log.Fatal("Failed to create directory: ", err)
	}

	return dataHome
}

func findGitDirs(startingDir string, dirs []string) []string {
	contents, err := os.ReadDir(startingDir)
	if err != nil {
		log.Fatal(err)
	}

outside:
	for _, content := range contents {
		for _, ignoreDir := range ignoreDirs {
			if ignoreDir == content.Name() {
				clog.Info("Ignoring " + ignoreDir)
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
				if entry.Name() == gitDirectory {
					dirs = append(dirs, path)
					break
				}
			}

			dirs = findGitDirs(path, dirs)
		}
	}

	return dirs
}

func gitInfo(dir string) ([]byte, error) {
	cmd := exec.Command(gitCommand[0], gitCommand[1:]...)
	cmd.Dir = dir

	out, err := cmd.Output()
	if err != nil {
		clog.Warn(
			"Could not find git history. Hint: Check there are commits in this repo: "+dir,
			err,
		)
		return nil, err
	}

	return out, nil
}

func collectCommits(dirs []string) []*GitCommit {
	var commits []*GitCommit

	for _, dir := range dirs {
		output, err := gitInfo(dir)
		if err != nil {
			clog.Warn(
				"Could not find git history. Hint: Check there are commits in this repo: "+dir,
				err,
			)
			continue
		}

		formattedCommits := formatCommit(output, dir)
		commits = append(commits, formattedCommits...)
	}

	return commits
}

func formatCommit(out []byte, dirTree string) []*GitCommit {
	if len(out) == 0 {
		return nil
	}

	commits := []*GitCommit{}

	for block := range strings.SplitSeq(string(out), commandSeperator+"\n") {
		lines := strings.Split(block, "\n")
		if len(lines) < commandLines {
			continue
		}

		dir := strings.Split(dirTree, "/")
		directory := dir[len(dir)-1]

		date, err := time.Parse(timeFormat, lines[2])
		if err != nil {
			clog.Error("Failed to format time of commit", err)
			continue
		}

		gc := GitCommit{
			Hash:      lines[0],
			Author:    lines[1],
			Directory: directory,
			Date:      date,
			Message:   lines[3],
		}

		commits = append(commits, &gc)
	}

	return commits
}

func writeJSONFile(commits []*GitCommit, dataHome string) {
	path := filepath.Join(dataHome, "commits.json")

	data, err := json.MarshalIndent(commits, "", "  ")
	if err != nil {
		clog.Fatal("Failed to marshal commits", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		clog.Fatal("Failed to write to file", err)
	}
}

func main() {
	homeDir := getHomeDir()
	dataHome := getDataHome(homeDir)

	gloLogFile := filepath.Join(dataHome, "glo.log")
	clog = clogger.New(gloLogFile)
	defer func() {
		if err := clog.Close(); err != nil {
			log.Fatal("Error closing log file")
		}
	}()

	dirs := []string{}
	dirs = findGitDirs(homeDir, dirs)

	commits := collectCommits(dirs)

	writeJSONFile(commits, dataHome)
}
