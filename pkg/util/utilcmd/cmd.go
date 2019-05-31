package utilcmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/stretchr/testify/mock"
)

// DefaultRunner defines the default implementation to execute a Cmd.
var DefaultRunner Runner = &CmdRunner{}

// Run runs the given Cmd.
func Run(cmd *exec.Cmd) error {
	return DefaultRunner.Run(cmd)
}

// RunOut runs the given Cmd and returns its output.
func RunOut(cmd *exec.Cmd) ([]byte, error) {
	return DefaultRunner.RunOut(cmd)
}

// Runner is an interface for running exec.Cmd.
type Runner interface {
	RunOut(cmd *exec.Cmd) ([]byte, error)
	Run(cmd *exec.Cmd) error
}

// CmdRunner is a Runner that executes exec.Cmd.
type CmdRunner struct{}

// Run executes the  given command, piping outputs to stdout and stderr.
func (*CmdRunner) Run(cmd *exec.Cmd) error {
	fmt.Println(strings.Join(cmd.Args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}

// Run executes the given command and returns its output.
func (*CmdRunner) RunOut(cmd *exec.Cmd) ([]byte, error) {
	fmt.Println(strings.Join(cmd.Args, " "))
	return cmd.CombinedOutput()
}

// MockRunner is a Runner that fakes the execution of exec.Cmd.
type MockRunner struct {
	mock.Mock
}

// Run fakes the execution of the given command.
func (r *MockRunner) Run(cmd *exec.Cmd) error {
	args := r.Called(cmd)
	return args.Error(0)
}

// Run fakes the execution of the given command.
func (r *MockRunner) RunOut(cmd *exec.Cmd) ([]byte, error) {
	args := r.Called(cmd)
	return args.Get(0).([]byte), args.Error(1)
}
