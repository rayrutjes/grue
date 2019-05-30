# Grue

Grue is a command line tool that facilitates building and publishing docker images, as well as deploying manifests to GKE clusters.

## Workflow

Grue is developed to accomodate the following workflow:
- you have a monorepo that holds code for several docker images and their respective k8s manifests, you deploy images to multiple clusters on GKE.
- you want to build and publish all images after each commits (see `grue build`).
- deployment is achieved through `kubectl apply`: manifests are meant to be updated first with the new image name and then applied (see `grue apply`).

## Usage

### Build

- Grue can only build using `docker` CLI. `docker` must already be installed.
- Images are automatically tagged with the current git commit. `git` must already be installed.
- Grue requires a `images.yaml` file at the root of the directory:


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

### Deploy

- Grue can only deploy using `gcloud` and `kubectl`. They both must already be installed.
- Deploy by using the `grue apply` command, which will look for manifests and apply them using `kubectl apply`.
- Grue is meant to work with multiple clusters, and will authenticate using `gcloud`.
- Grue requires a `images.yaml` file at the root of the directory:


```
apiVersion: grue/v1
kind: Config
build:
  script: scripts/build.sh
  artifacts:
  - image: gcr.io/my-project/my-service
    context: cmd/foo/my-service
deploy:
  clusters:
  - name: my-cluster
    project: my-gcp-project
    region: us-east1
    manifests: k8s/
```

- `clusters`: lists clusters to deploy to.
- `name`: name of the cluster. will be used for authentication using gcloud.
- `project`: the gcp project of the cluster. will be used for authentication using gcloud.
- `region`: region of the cluster. will be used for authentication using gcloud.
- `manifests`: directory where the manifests are located. Will be scanned recursively for .yaml or .yml files.

```
$ grue apply --help
Apply manifests using kubectl

Usage:
  grue apply [flags]

Flags:
      --cluster string   Specify which cluster to target. If none, apply manifests for all clusters.
  -h, --help             help for apply
```

Examples:
- `grue apply`: apply all manifests for all clusters found in `images.yaml`.
- `grue apply --cluster=my-project`: apply all manifests for the specified cluster. The cluster must exist in `images.yaml`.
