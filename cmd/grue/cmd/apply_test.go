package cmd

import (
	"os/exec"
	"testing"

	"github.com/algolia/grue/pkg/schema"
	"github.com/algolia/grue/pkg/util/utilcmd"

	"github.com/stretchr/testify/require"
)

func TestApplyCluster(t *testing.T) {
	r := &utilcmd.MockRunner{}
	r.On("Run", exec.Command("gcloud", "container", "clusters", "get-credentials", "my-cluster", "--region", "us-east1", "--project", "my-project")).Return(nil)
	r.On("Run", exec.Command("kubectl", "apply", "-f", "testdata/applytest/manifest.yml")).Return(nil)
	r.On("Run", exec.Command("kubectl", "apply", "-f", "testdata/applytest/manifest.yaml")).Return(nil)

	defer func(r utilcmd.Runner) { utilcmd.DefaultRunner = r }(utilcmd.DefaultRunner)
	utilcmd.DefaultRunner = r

	c := schema.Cluster{
		Name:      "my-cluster",
		Region:    "us-east1",
		Project:   "my-project",
		Manifests: "testdata/applytest",
	}
	err := applyCluster(c)
	require.NoError(t, err)
	require.True(t, r.AssertExpectations(t))
}
