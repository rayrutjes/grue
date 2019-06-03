# Grue

Grue is a command line tool that facilitates building and publishing docker images, as well as deploying manifests to GKE clusters.

## Why Grue

Grue is developed to accomodate the following workflow:
- you have a monorepo that contains code for several docker images and their respective kubernetes manifests.
- you want to build and publish all images after each commits: see `grue build`.
- you deploy images through `kubectl apply` (you first update manifests with the new image name) and want a way to apply everything at once, possibly across clusters: see `grue apply`.

Grue is similar to [skaffold](https://skaffold.dev/): config is through a top level yaml, and `grue build`, like `skaffold build`, will look for all images in the repo to build them. The similarities stop there though: `grue` is not meant to auto deploy images. In our workflow, manifests are first manually updated, and then applied. Therefore, `grue apply` is only meant to be a helper around `kubectl apply`.

## Limitations

Grue is very young and should be considered experimental at this point.

Grue calls `git`, `docker`, `gcloud`, and `kubectl`. Therefore today Grue is only useful to deploy docker images to GKE clusters from a git repository.

## Usage

### Build

`grue build` build docker images found in the repository.

- Images are automatically tagged with the current git commit hash.
- Start by adding a `images.yaml` file at the root of the directory to list the images you want to build, e.g:

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

Examples:
- `grue build`: builds all the images referenced under `build.artifacts`, executing `build.script` (if any) before each.
- `grue build --image=gcr.io/my-project/my-service`: builds the specified image. The image must be referenced in the `images.yaml`.
- `grue build --publish`: builds all the images, and publishes them using `docker push`.

### Apply

`grue apply` apply manifests using kubectl.

- Add a `deploy` section in the `images.yaml` to list clusters you want to deploy to, and where to find their respective deployment manifests.
- `grue apply` runs `gloud` to authenticate against each clusters then looks for all its manifests (`.yaml` or `.yml`) before running `kubectl apply` on them.


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


Examples:
- `grue apply`: apply manifests for all clusters found in `images.yaml`.
- `grue apply --cluster=my-project`: apply all manifests for the specified cluster. The cluster must exist in `images.yaml`.
