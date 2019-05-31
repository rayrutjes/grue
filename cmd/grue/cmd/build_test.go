package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/algolia/grue/pkg/schema"
	"github.com/algolia/grue/pkg/util/utilcmd"

	"github.com/stretchr/testify/require"
)

func TestBuildArtifact(t *testing.T) {
	r := &utilcmd.MockRunner{}
	r.On("RunOut", exec.Command("git", "rev-parse", "--short", "HEAD")).Return([]byte("123abc"), nil)

	buildCmd := exec.Command("/bin/sh", "scripts/build.sh")
	buildCmd.Env = append(os.Environ(), "CONTEXT=cmd/foo/bar", "IMAGE_BASE_NAME=bar")
	r.On("Run", buildCmd).Return(nil)

	p, err := filepath.Abs("cmd/foo/bar")
	require.NoError(t, err)
	dockerCmd := exec.Command("docker", "build", "-t", "gcr.io/foo/bar:123abc", p)
	dockerCmd.Dir = "cmd/foo/bar"
	r.On("Run", dockerCmd).Return(nil)

	publishCmd := exec.Command("docker", "push", "gcr.io/foo/bar:123abc")
	publishCmd.Dir = "cmd/foo/bar"
	r.On("Run", publishCmd).Return(nil)

	defer func(r utilcmd.Runner) { utilcmd.DefaultRunner = r }(utilcmd.DefaultRunner)
	utilcmd.DefaultRunner = r

	a := schema.Artifact{
		Image:   "gcr.io/foo/bar",
		Context: "cmd/foo/bar",
	}
	err = buildArtifact(a, "scripts/build.sh", true)
	require.NoError(t, err)
	require.True(t, r.AssertExpectations(t))
}
