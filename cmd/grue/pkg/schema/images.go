package schema

// Config defines the structure of the images.yaml config file.
type Config struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	// Build defines how images are built.
	Build Build `yaml:"build"`
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
	Image string `yaml:"image"`
	// Context is the directory containing the image's Dockerfile.
	Context string `yaml:"context"`
}
