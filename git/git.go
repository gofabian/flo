package git

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Repository struct {
	URL    string
	Branch string
}

func GetRepository() (*Repository, error) {
	url, err := readRemoteUrl(".git/config")
	if err != nil {
		return nil, err
	}

	branch, err := readBranch(".git/HEAD")
	if err != nil {
		return nil, err
	}

	return &Repository{URL: url, Branch: branch}, nil
}

func readRemoteUrl(pathToGitConfig string) (string, error) {
	file, err := os.Open(pathToGitConfig)
	if err != nil {
		return "", fmt.Errorf("cannot open file %s, %w", pathToGitConfig, err)
	}
	defer file.Close()

	var remoteName string
	var firstRemoteURL string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(content, "[remote") {
			// found "remote" section, e.g. [remote "origin"]
			elements := strings.SplitN(content, " ", 2)
			if len(elements) != 2 {
				return "", fmt.Errorf("Invalid remote subsection in %s: %s", pathToGitConfig, content)
			}
			remoteName = strings.Trim(elements[1], "\"]")
		} else if strings.HasPrefix(content, "[") {
			// found non-"remote" section, e.g. [core]
			remoteName = ""
		} else if remoteName != "" {
			// found setting within "remote" section, e.g. url = git@github.com:org/repo.git
			elements := strings.SplitN(content, " = ", 2)
			if len(elements) != 2 {
				return "", fmt.Errorf("Invalid key/value format in %s: %s", pathToGitConfig, content)
			}
			if elements[0] == "url" {
				if remoteName == "origin" {
					return elements[1], nil
				}
				if firstRemoteURL == "" {
					firstRemoteURL = elements[1]
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("cannot read file %s, %w", pathToGitConfig, err)
	}

	if firstRemoteURL == "" {
		return "", fmt.Errorf("Could not find URL of git repo in %s", pathToGitConfig)
	}
	return firstRemoteURL, nil
}

func readBranch(pathToHeadFile string) (string, error) {
	// e.g. ref: refs/heads/develop
	file, err := os.Open(pathToHeadFile)
	if err != nil {
		return "", fmt.Errorf("cannot open file %s, %w", pathToHeadFile, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	bytes, _, err := reader.ReadLine()
	if err != nil {
		return "", fmt.Errorf("cannot read file %s, %w", pathToHeadFile, err)
	}
	line := string(bytes)

	if !strings.HasPrefix(line, "ref: refs/heads/") {
		return "", fmt.Errorf("cannot read branch from %s", pathToHeadFile)
	}
	branch := line[len("ref: refs/heads/"):]
	return branch, nil
}
