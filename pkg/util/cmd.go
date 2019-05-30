package util

import (
	"os"
	"os/exec"
)

// RunCmdOut executes the  given command, piping outputs to stdout and stderr.
func RunCmdOut(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}
