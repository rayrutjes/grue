package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/algolia/grue/pkg/schema"
	"github.com/algolia/grue/pkg/util"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
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
	f, err := ioutil.ReadFile("images.yaml")
	if err != nil {
		return err
	}
	var config schema.Config
	err = yaml.Unmarshal(f, &config)
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

	return filepath.Walk(c.Manifests, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// Consider only yaml or yml files.
		if !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml") {
			return nil
		}
		return applyManifest(path)
	})
}

func applyManifest(file string) error {
	fmt.Printf("kubectl apply -f %s\n", file)
	cmd := exec.Command("kubectl", "apply", "-f", file)
	return util.RunCmdOut(cmd)
}

func auth(c schema.Cluster) error {
	fmt.Printf("gcloud beta container clusters get-credentials %s --region %s --project %s\n", c.Name, c.Region, c.Project)
	cmd := exec.Command("gcloud", "beta", "container", "clusters", "get-credentials", c.Name, "--region", c.Region, "--project", c.Project)
	return util.RunCmdOut(cmd)
}
