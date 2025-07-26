package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
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
		fmt.Printf(
			"Error getting repo info for directory %s: %v\nHint: Check there are commits in this Repo.\n",
			dir,
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
			fmt.Printf("Error parsing date: %v\n", err)
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
		log.Fatalf("failed to marshal commits: %v", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		log.Fatalf("failed to write file %s: %v", path, err)
	}
}

func main() {
	homeDir := getHomeDir()
	dataHome := getDataHome(homeDir)

	dirs := []string{}
	dirs = findGitDirs(homeDir, dirs)

	commits := collectCommits(dirs)

	writeJSONFile(commits, dataHome)
}
