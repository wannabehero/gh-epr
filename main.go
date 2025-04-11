package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

var emojis = []string{
	"🚀", "🤖", "🎢", "🪨", "🎨",
	"🧻", "🗞️", "📎", "🌚", "💸",
	"📝", "☠️", "🐕", "🐩", "🐃",
	"♾️", "🐜", "🦁", "🐺", "🦊",
}

func getRandomEmoji() string {
	return emojis[rand.Intn(len(emojis))]
}

func runCmd(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error running command: %v", err)
	}
	return strings.TrimSpace(string(output))
}

func defaultTitle(baseBranch string) string {
	commits := runCmd("git", "log", "--pretty=format:%s", fmt.Sprintf("origin/%s..HEAD", baseBranch))
	messages := strings.Split(commits, "\n")
	if len(messages) == 1 && messages[0] != "" {
		return messages[0]
	}

	branchName := runCmd("git", "rev-parse", "--abbrev-ref", "HEAD")
	return branchName
}

func detectBaseBranch() string {
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

func main() {
	var title string

	baseBranch := detectBaseBranch()

	if len(os.Args) >= 2 {
		title = os.Args[1]
	} else {
		title = defaultTitle(baseBranch)
	}

	fullTitle := fmt.Sprintf("%s %s", getRandomEmoji(), title)
	fmt.Println(fullTitle)
}
