package git

import (
	"fmt"
	"strings"
)

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

func GenerateTitle(commits []string) *string {
	if len(commits) == 1 && commits[0] != "" {
		return &commits[0]
	}

	branchName, err := runCmd("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return nil
	}

	return &branchName
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
