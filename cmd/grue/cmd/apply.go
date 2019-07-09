package cmd

import (
	"os/exec"

	"github.com/algolia/grue/pkg/schema"
	"github.com/algolia/grue/pkg/util/utilcmd"
	"github.com/spf13/cobra"
)

var cluster string

func init() {
	applyCmd.PersistentFlags().StringVar(&cluster, "cluster", "", "Specify which cluster to target. If none, apply manifests for all clusters.")
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply manifests using kubectl",
	RunE:  apply,
}

func apply(cmd *cobra.Command, args []string) error {
	config, err := schema.New("images.yaml")
	if err != nil {
		return err
	}

	for _, c := range config.Deploy.Clusters {
		if cluster != "" && cluster != c.Name {
			continue
		}
		err := applyCluster(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func applyCluster(c schema.Cluster) error {
	err := auth(c)
	if err != nil {
		return err
	}

	for _, manifest := range c.Manifests {
		err := applyManifest(manifest)
		if err != nil {
			return err
		}
	}
	return nil
}

func applyManifest(file string) error {
	cmd := exec.Command("kubectl", "apply", "-f", file)
	return utilcmd.Run(cmd)
}

func auth(c schema.Cluster) error {
	cmd := exec.Command("gcloud", "container", "clusters", "get-credentials", c.Name, "--region", c.Region, "--project", c.Project)
	return utilcmd.Run(cmd)
}
