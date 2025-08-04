package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/double-tilde/glo/pkg/config"
)

type GitCommit struct {
	Hash      string    `json:"hash"`
	Author    string    `json:"author"`
	Directory string    `json:"directory"`
	Date      time.Time `json:"date"`
	Message   string    `json:"message"`
}

func GitInfo(dir string) ([]byte, error) {
	cmd := exec.Command(config.GitCommand[0], config.GitCommand[1:]...)
	cmd.Dir = dir

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}

func FormatCommit(dirTree string, out []byte) ([]*GitCommit, error) {
	if len(out) == 0 {
		return nil, errors.New("no output to commit")
	}

	commits := []*GitCommit{}

	for block := range strings.SplitSeq(string(out), config.CommandSeperator+"\n") {
		lines := strings.Split(block, "\n")
		if len(lines) < config.CommandLines {
			continue
		}

		dir := strings.Split(dirTree, "/")
		directory := dir[len(dir)-1]

		date, err := time.Parse(config.TimeFormat, lines[2])
		if err != nil {
			return nil, errors.New("failed to format time of commit")
			// continue
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

	return commits, nil
}

func CollectCommits(dirs []string) ([]*GitCommit, []error) {
	var commits []*GitCommit
	var errs []error

	for _, dir := range dirs {
		output, err := GitInfo(dir)
		if err != nil {
			errs = append(errs, fmt.Errorf("could not find git history in %s %v", dir, err))
			continue
		}

		formattedCommits, err := FormatCommit(dir, output)
		if err != nil {
			errs = append(errs, fmt.Errorf("could not format commit infomation in %s %v", dir, err))
			continue
		}

		commits = append(commits, formattedCommits...)
	}

	if len(errs) > 0 {
		return commits, errs
	}

	return commits, nil
}

func WriteJSONFile(commits []*GitCommit, dataHome string) error {
	path := filepath.Join(dataHome, config.GloCommitsFile)

	data, err := json.MarshalIndent(commits, "", "  ")
	if err != nil {
		return errors.New("failed to marshal commits")
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return errors.New("failed to write to file")
	}

	return nil
}

func ReadJSONFile(dataHome string) ([]*GitCommit, error) {
	var commits []*GitCommit
	path := filepath.Join(dataHome, config.GloCommitsFile)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("failed to read file")
	}

	if err := json.Unmarshal(data, &commits); err != nil {
		return nil, errors.New("failed to unmarshal commits")
	}

	return commits, nil
}

func GetYearOfCommits(commits []*GitCommit) []time.Time {
	var recentCommits []time.Time
	oneYearAgo := time.Now().AddDate(-1, 0, 0)

	for _, commit := range commits {
		if !commit.Date.Before(oneYearAgo) {
			recentCommits = append(recentCommits, commit.Date)
		}
	}

	sort.Slice(recentCommits, func(i, j int) bool {
		return recentCommits[i].Before(recentCommits[j])
	})

	return recentCommits
}
