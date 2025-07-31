package data

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/double-tilde/glo/pkg/config"
	"github.com/double-tilde/glo/pkg/model"
)

// TODO: Return errors, do not rely on other packages

// gitInfo returns the output of the git command.
// Parameters:
//
//	dir: the directory to run the command in.
//
// Returns:
//
//	[]byte: the standard output of the command.
//	error: an error if the git command cannot return an output.
func gitInfo(dir string) ([]byte, error) {
	cmd := exec.Command(config.GitCommand[0], config.GitCommand[1:]...)
	cmd.Dir = dir

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}

// TODO: Does the formatCommit function need refactoring?

// formatCommit returns the commit in a standard structure so it is easier to work with.
// Parameters:
//
//	dirTree: the full path of the directory the commits came from.
//	out: the unformatted commits for the directory.
//
// Returns:
//
//	[]*model.GitCommit: the formmatted commit.
//	error: an error if there is no output, or if the output cannot be formatted.
func formatCommit(dirTree string, out []byte) ([]*model.GitCommit, error) {
	if len(out) == 0 {
		return nil, errors.New("no output to commit")
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
			return nil, errors.New("failed to format time of commit")
			// continue
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

	return commits, nil
}

// CollectCommits returns each commit formatted and stored in a slice.
// Parameters:
//
//	dirs: the directories to collect the git commits from.
//
// Returns:
//
//	[]*model.GitCommit: the formmated commits stored in a slice.
//	error: an error if the git command cannot return an output.
func CollectCommits(dirs []string) ([]*model.GitCommit, error) {
	var commits []*model.GitCommit

	for _, dir := range dirs {
		output, err := gitInfo(dir)
		if err != nil {
			return nil, errors.New("could not find git history in directory" + dir)
		}

		// TODO: What about this error?
		formattedCommits, _ := formatCommit(dir, output)
		commits = append(commits, formattedCommits...)
	}

	return commits, nil
}

// WriteJsonFile returns each commit formatted and stored in a slice.
// Parameters:
//
//	commits: the collection of commits to write to the json file.
//	dataHome: the place to write and store the json file.
//
// Returns:
//
//	error: an error if the json cannot be mashaled or if the file cannot be written to.
func WriteJSONFile(commits []*model.GitCommit, dataHome string) error {
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
