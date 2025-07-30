package data

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/double-tilde/glo/pkg/config"
	"github.com/double-tilde/glo/pkg/model"
)

// TODO: Return errors, do not rely on other packages
// var clog *logger.Clogger

func GitInfo(dir string) ([]byte, error) {
	cmd := exec.Command(config.GitCommand[0], config.GitCommand[1:]...)
	cmd.Dir = dir

	out, err := cmd.Output()
	if err != nil {
		// clog.Warn(
		// 	"Could not find git history. Hint: Check there are commits in this repo: "+dir,
		// 	err,
		// )
		log.Fatal(err)
		return nil, err
	}

	return out, nil
}

func CollectCommits(dirs []string) []*model.GitCommit {
	var commits []*model.GitCommit

	for _, dir := range dirs {
		output, err := GitInfo(dir)
		if err != nil {
			// clog.Warn(
			// 	"Could not find git history. Hint: Check there are commits in this repo: "+dir,
			// 	err,
			// )
			log.Fatal(err)
			continue
		}

		formattedCommits := FormatCommit(output, dir)
		commits = append(commits, formattedCommits...)
	}

	return commits
}

func FormatCommit(out []byte, dirTree string) []*model.GitCommit {
	if len(out) == 0 {
		return nil
	}

	commits := []*model.GitCommit{}

	for block := range strings.SplitSeq(string(out), config.CommandSeperator+"\n") {
		lines := strings.Split(block, "\n")
		if len(lines) < config.CommandLines {
			continue
		}

		dir := strings.Split(dirTree, "/")
		directory := dir[len(dir)-1]

		date, err := time.Parse(config.TimeFormat, lines[2])
		if err != nil {
			// clog.Error("Failed to format time of commit", err)
			log.Fatal(err)
			continue
		}

		gc := model.GitCommit{
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

func WriteJSONFile(commits []*model.GitCommit, dataHome string) {
	path := filepath.Join(dataHome, "commits.json")

	data, err := json.MarshalIndent(commits, "", "  ")
	if err != nil {
		// clog.Fatal("Failed to marshal commits", err)
		log.Fatal(err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		// clog.Fatal("Failed to write to file", err)
		log.Fatal(err)
	}
}
