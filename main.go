package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/wannabehero/gh-epr/utils"
)

func main() {
	baseBranch, err := utils.DetectBaseBranch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting base branch: %v\n", err)
		os.Exit(1)
	}
	
	commits, err := utils.GetCommitsHistory(baseBranch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting commit history: %v\n", err)
		os.Exit(1)
	}

	var fullTitle string

	if generatedTitle := utils.GenerateTitle(commits); generatedTitle != nil {
		fullTitle = *generatedTitle
	} else {
		defaultTitle, err := utils.GetDefaultTitle(baseBranch)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting default title: %v\n", err)
			os.Exit(1)
		}

		fullTitle = fmt.Sprintf("%s %s", utils.GetRandomEmoji(), defaultTitle)
	}

	extraArgs := os.Args[1:]
	args := []string{"pr", "create", "--title", fullTitle}

	if generatedBody := utils.GenerateBody(commits); generatedBody != nil {
		args = append(args, "--body", *generatedBody)
	}

	args = append(args, extraArgs...)

	cmd := exec.Command("gh", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating pull request: %v\n", err)
		os.Exit(1)
	}
}
