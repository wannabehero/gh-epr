package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/wannabehero/gh-epr/utils"
)

func main() {
	baseBranch := utils.DetectBaseBranch()
	commits := utils.GetCommitsHistory(baseBranch)

	var fullTitle string

	if generatedTitle := utils.GenerateTitle(commits); generatedTitle != nil {
		fullTitle = *generatedTitle
	} else {
		defaulTitle := utils.GetDefaultTitle(baseBranch)

		fullTitle = fmt.Sprintf("%s %s", utils.GetRandomEmoji(), defaulTitle)
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

	cmd.Run()
}
