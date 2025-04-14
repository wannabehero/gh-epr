package git

import (
	"fmt"
	"os/exec"
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
