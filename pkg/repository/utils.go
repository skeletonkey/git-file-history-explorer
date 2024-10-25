package repository

import (
	"os/exec"
	"strings"

	"github.com/skeletonkey/git-file-history-explorer/pkg/report"
)

func executeCmd(cmdName string, args ...string) string {
	var out strings.Builder
	cmd := exec.Command(cmdName, args...)
	cmd.Stdout = &out
	err := cmd.Run()
	report.PanicOnError(err)
	return strings.TrimRight(out.String(), "\n")
}
