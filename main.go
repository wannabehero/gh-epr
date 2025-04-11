package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gh epr <pr title>")
		return
	}

	title := os.Args[1]

	fullTitle := fmt.Sprintf("%s %s", getRandomEmoji(), title)

	extraArgs := os.Args[2:]

	args := append([]string{"pr", "create", "--title", fullTitle}, extraArgs...)

	cmd := exec.Command("gh", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()
}
