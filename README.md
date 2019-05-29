# Grue

Grue is a command line tool that facilitates building and publishing docker images.

## Getting Started

Add a `images.yaml` in the root directory:

```
apiVersion: grue/v1
kind: Config
build:
  script: scripts/build.sh
  artifacts:
  - image: gcr.io/my-project/my-service
    context: cmd/foo/my-service
```

- `script`: points to a script called before building images. optional.
- `artifacts`: lists all the images that are buildable.
- `image`: the name of the image to build.
- `context`: the directory containing the image's Dockerfile.

## Usage

### Build

- Grue can only build using `docker` CLI. `docker` must already be installed.
- Images are automatically tagged with the current git commit. `git` must already be installed.

```
$ grue build --help
Build images using Docker CLI.

Usage:
  grue build [flags]

Flags:
  -h, --help           help for build
      --image string   Specificy which image to build. If none, builds all configured images.
      --publish        If true, also publishes the image(s).
```

Examples:
- `grue build`: builds all the images referenced under `build.artifacts`, executing `build.script` before each.
- `grue build --image=gcr.io/my-project/my-service`: builds the specified image. The image must be referenced in the `images.yaml`.
- `grue build --publish`: builds all the images, and publishes them using `docker push`
