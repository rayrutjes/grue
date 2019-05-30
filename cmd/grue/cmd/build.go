package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/algolia/grue/pkg/schema"
	"github.com/algolia/grue/pkg/util"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var publish bool
var image string

func init() {
	buildCmd.PersistentFlags().StringVar(&image, "image", "", "Specificy which image to build. If none, builds all configured images.")
	buildCmd.PersistentFlags().BoolVar(&publish, "publish", false, "If true, also publishes the image(s).")
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build images",
	Long:  "Build images using Docker CLI.",
	RunE:  build,
}

func build(cmd *cobra.Command, args []string) error {
	f, err := ioutil.ReadFile("images.yaml")
	if err != nil {
		return err
	}
	var config schema.Config
	err = yaml.Unmarshal(f, &config)
	if err != nil {
		return err
	}

	for _, a := range config.Build.Artifacts {
		if image != "" && a.Image != image {
			continue
		}
		err = buildArtifact(a, config.Build.Script, publish)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildArtifact(a schema.Artifact, script string, publish bool) error {
	tag, err := gitTag()
	if err != nil {
		return err
	}

	err = buildImage(a, script, tag)
	if err != nil {
		return err
	}
	if !publish {
		return nil
	}
	return publishImage(a, tag)
}

func buildImage(a schema.Artifact, script string, tag string) error {
	if script != "" {
		err := runBuildScript(a, script)
		if err != nil {
			return err
		}
	}
	return runDockerBuild(a, tag)
}

func publishImage(a schema.Artifact, tag string) error {
	image := fmt.Sprintf("%s:%s", a.Image, tag)
	cmd := exec.Command("docker", "push", image)
	cmd.Dir = a.Context
	return util.RunCmdOut(cmd)
}

func runBuildScript(a schema.Artifact, script string) error {
	env := append(os.Environ(),
		fmt.Sprintf("%s=%s", "CONTEXT", a.Context),
		fmt.Sprintf("%s=%s", "IMAGE_BASE_NAME", filepath.Base(a.Image)),
	)
	cmd := exec.Command("/bin/sh", script)
	cmd.Env = env
	return util.RunCmdOut(cmd)
}

func runDockerBuild(a schema.Artifact, tag string) error {
	image := fmt.Sprintf("%s:%s", a.Image, tag)
	path, err := filepath.Abs(a.Context)
	if err != nil {
		return err
	}
	cmd := exec.Command("docker", "build", "-t", image, path)
	cmd.Dir = a.Context
	return util.RunCmdOut(cmd)
}

func gitTag() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(out)), nil
}
