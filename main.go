package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/wannabehero/gh-epr/git"
	"github.com/wannabehero/gh-epr/llm"
	"github.com/wannabehero/gh-epr/utils"
)

func getTitle(commits []string) *string {
	if generatedTitle := llm.GenerateTitle(commits); generatedTitle != nil {
		return generatedTitle
	} else if defaultTitle := git.GenerateTitle(commits); defaultTitle != nil {
		title := fmt.Sprintf("%s %s", utils.GetRandomEmoji(), *defaultTitle)
		return &title
	}
	return nil
}

func getBody(commits []string, diff string, template string) *string {
	return llm.GenerateBody(commits, diff, template)
}

func main() {
	baseBranch, err := git.DetectBaseBranch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting base branch: %v\n", err)
		os.Exit(1)
	}

	commits, err := git.GetCommitsHistory(baseBranch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting commit history: %v\n", err)
		os.Exit(1)
	}

	diff, err := git.GetDiff(baseBranch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting diff: %v\n", err)
		os.Exit(1)
	}

	template, err := git.GetPRTemplate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting PR template: %v\n", err)
		os.Exit(1)
	}

	args := []string{"pr", "create"}

	if title := getTitle(commits); title != nil {
		args = append(args, "--title", *title)
	}

	if body := getBody(commits, diff, template); body != nil {
		args = append(args, "--body", *body)
	}

	extraArgs := os.Args[1:]
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
