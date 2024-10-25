package repository

import (
	"os/exec"
	"strings"
)

func executeCmd(cmdName string, args ...string) (string, error) {
	var out strings.Builder
	cmd := exec.Command(cmdName, args...)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(out.String(), "\n"), nil
}
