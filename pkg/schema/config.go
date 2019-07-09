package schema

import (
	"io/ioutil"

	"github.com/asaskevich/govalidator"
	yaml "gopkg.in/yaml.v2"
)

// Config defines the structure of the images.yaml config file.
type Config struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	// Build defines how images are built.
	Build Build `yaml:"build"`
	// Deploy defines how images can be deployed.
	Deploy Deploy `yaml:"deploy"`
}

// Build defines how images are built.
type Build struct {
	// Script points to a script (relative) called before building images.
	Script string `yaml:"script"`
	// Artifacts lists all the images that are buildable.
	Artifacts []Artifact `yaml:"artifacts"`
}

// Artifact holds information related to building a single artifact.
type Artifact struct {
	// Image defines the name of the image to build. e.g: `gcr.io/my-project/foo`.
	Image string `yaml:"image" valid:"required"`
	// Context is the directory containing the image's Dockerfile.
	Context string `yaml:"context" valid:"required"`
}

// Deploy holds information related to deploying manifests using `kubectl apply`.
type Deploy struct {
	// Clusters lists GKE clusters where manifests can be applied.
	Clusters []Cluster `yaml:"clusters"`
}

// Cluster defines a specific GKE cluster with the manifests that can be applied
// in it.
type Cluster struct {
	// Name define the cluster name
	Name string `yaml:"name" valid:"required"`
	// Project defines the GCP project
	Project string `yaml:"project" valid:"required"`
	// Region defines the cluster location
	Region string `yaml:"region" valid:"required"`
	// Manifests points to folders where manifests are located, e.g: `k8s/`
	Manifests []string `yaml:"manifests" valid:"required"`
}

// New parses a yaml file into a schema.Config.
func New(configPath string) (*Config, error) {
	f, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(f, &config)
	if err != nil {
		return nil, err
	}

	_, err = govalidator.ValidateStruct(config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
