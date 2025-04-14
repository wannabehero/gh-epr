package git

import (
	"os"
	"path/filepath"
	"strings"
)

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
