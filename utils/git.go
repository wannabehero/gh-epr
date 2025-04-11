package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func runCmd(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error running command: %v", err)
	}
	return strings.TrimSpace(string(output))
}

func GetCommitsHistory(baseBranch string) []string {
	commits := runCmd("git", "log", "--pretty=format:%s", fmt.Sprintf("origin/%s..HEAD", baseBranch))
	messages := strings.Split(commits, "\n")
	return messages
}

func GetDefaultTitle(baseBranch string) string {
	messages := GetCommitsHistory(baseBranch)
	if len(messages) == 1 && messages[0] != "" {
		return messages[0]
	}

	branchName := runCmd("git", "rev-parse", "--abbrev-ref", "HEAD")
	return branchName
}

func DetectBaseBranch() string {
	branches := runCmd("git", "branch", "-r")
	if strings.Contains(branches, "origin/main") {
		return "main"
	}
	if strings.Contains(branches, "origin/master") {
		return "master"
	}

	// Default to main branch if no other branch is found
	return "main"
}
