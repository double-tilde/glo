package main

import (
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

	// Linux
	if dataHome = os.Getenv("XDG_DATA_HOME"); dataHome != "" {
		return dataHome
	}

	// Windows
	if dataHome = os.Getenv("LOCALAPPDATA"); dataHome != "" {
		return dataHome
	}

	if homeDir != "" {
		return filepath.Join(homeDir, ".local", "share")
	}

	// Last resort: store in home directory
	return homeDir
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
			hash:      lines[0],
			author:    lines[1],
			directory: directory,
			date:      date,
			message:   lines[3],
		}

		commits = append(commits, &gc)
	}

	return commits
}

func main() {
	homeDir := getHomeDir()

	dataHome := getDataHome(homeDir)
	println(dataHome)

	dirs := []string{}
	dirs = findGitDirs(homeDir, dirs)

	commits := []*GitCommit{}

	for _, dir := range dirs {
		output, err := gitInfo(dir)
		if err != nil {
			break
		}

		formattedCommits := formatCommit(output, dir)
		commits = append(commits, formattedCommits...)
	}

	for _, commit := range commits {
		fmt.Println(commit.hash)
		fmt.Println(commit.author)
		fmt.Println(commit.directory)
		fmt.Println(commit.date)
		fmt.Println(commit.message)
	}
}
