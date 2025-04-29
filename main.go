package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/wannabehero/gh-aipr/config"
	"github.com/wannabehero/gh-aipr/git"
	"github.com/wannabehero/gh-aipr/llm"
	"github.com/wannabehero/gh-aipr/utils"
)

func getTitleAndBody(commits []string, diff string, template string, ctx context.Context) (*string, *string) {
	title, body := llm.GenerateTitleAndBody(commits, diff, template, ctx)

	if title == nil {
		if defaultTitle := git.GenerateTitle(commits); defaultTitle != nil {
			value := fmt.Sprintf("%s %s", utils.GetRandomEmoji(), *defaultTitle)
			title = &value
		}
	}

	return title, body
}

func main() {
	config.LoadConfig()

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

	ctx := context.Background()
	title, body := getTitleAndBody(commits, diff, template, ctx)

	if title != nil {
		args = append(args, "--title", *title)
	}

	if body != nil {
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
