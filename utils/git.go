package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func runCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running command '%s %s': %w", name, strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(output)), nil
}

func findRemoteRef(baseBranch string) (string, error) {
	remotes, err := runCmd("git", "remote")
	if err != nil {
		return "", err
	}

	remotesList := strings.Split(remotes, "\n")

	remoteRefs, err := runCmd("git", "branch", "-r")
	if err != nil {
		return "", err
	}

	for _, remote := range remotesList {
		if remote == "" {
			continue
		}

		remoteRef := fmt.Sprintf("%s/%s", remote, baseBranch)

		if strings.Contains(remoteRefs, remoteRef) {
			return remoteRef, nil
		}
	}

	return "", fmt.Errorf("no remote reference found for branch %s", baseBranch)
}

func GetCommitsHistory(baseBranch string) ([]string, error) {
	remoteRef, err := findRemoteRef(baseBranch)
	if err != nil {
		return []string{}, nil
	}

	commits, err := runCmd("git", "log", "--pretty=format:%s", fmt.Sprintf("%s..HEAD", remoteRef))
	if err != nil || commits == "" {
		return []string{}, nil
	}

	return strings.Split(commits, "\n"), nil
}

func GetDefaultTitle(baseBranch string) (string, error) {
	messages, err := GetCommitsHistory(baseBranch)
	if err != nil {
		return "", err
	}

	if len(messages) == 1 && messages[0] != "" {
		return messages[0], nil
	}

	branchName, err := runCmd("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}

	return branchName, nil
}

func DetectBaseBranch() (string, error) {
	branches, err := runCmd("git", "branch", "-r")
	if err != nil {
		return "", err
	}

	baseBranchNames := []string{"main", "master", "dev", "develop", "trunk"}

	remotes, err := runCmd("git", "remote")
	if err != nil {
		return "", err
	}

	remotesList := strings.Split(remotes, "\n")

	for _, baseName := range baseBranchNames {
		for _, remote := range remotesList {
			if remote == "" {
				continue
			}

			remoteRef := fmt.Sprintf("%s/%s", remote, baseName)
			if strings.Contains(branches, remoteRef) {
				return baseName, nil
			}
		}
	}

	return "main", nil
}

func GetDiff(baseBranch string) (string, error) {
	remoteRef, err := findRemoteRef(baseBranch)
	if err != nil {
		return "", nil
	}

	diff, err := runCmd("git", "diff", remoteRef, "HEAD")
	if err != nil {
		return "", err
	}

	return diff, nil
}

func GetPRTemplate() (string, error) {
	repoRoot, err := runCmd("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}

	githubDir := filepath.Join(repoRoot, ".github")
	if _, err := os.Stat(githubDir); os.IsNotExist(err) {
		return "", nil
	}

	dirEntries, err := os.ReadDir(githubDir)
	if err != nil {
		return "", nil
	}

	for _, entry := range dirEntries {
		if !entry.IsDir() && strings.EqualFold(entry.Name(), "pull_request_template.md") {
			content, err := os.ReadFile(filepath.Join(githubDir, entry.Name()))
			if err == nil {
				return string(content), nil
			}
		}
	}

	return "", nil
}
