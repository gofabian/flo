# `flo` = Let Concourse <span style="color:green">**fl**</span>y with Dr<span style="color:green">**o**</span>ne

`flo` is an adapter that runs [Drone] pipelines on a [Concourse] CI server.

> Did you ever shout because the syntax of Concourse pipelines is too cumbersome?

> Did you ever wonder why Concourse does not support multibranch pipelines like [Jenkins] or Drone CI?

Then give `flo` a try and run Drone pipelines in your Concourse environment.

[Drone]: https://www.drone.io/
[Concourse]: https://concourse-ci.org/
[Jenkins]: https://www.jenkins.io/

## Usage

At first sign in to your Concourse team:

    $ fly -t mytarget login \
        -c https://ci-server \
        -n myteam

We highly recommend that you use a separate Concourse team for each multibranch pipeline (= for each Git repository).

The easiest way to run the `flo` binary is to use the [`gofabian/flo` Docker image from Dockerhub][gofabian/flo], e.g.

[gofabian/flo]: https://hub.docker.com/repository/docker/gofabian/flo


    $ docker run -it gofabian/flo:0 flo --help

Setup a multibranch pipeline...

    $ docker run -it gofabian/flo:0 \
        flo setup-pipeline \
            --style multibranch \
            --target mytarget \
            --git-url https://github.com/org/repo.git

![A multibranch pipeline in Concourse](https://github.com/gofabian/flo/raw/main/doc/multibranch.png "A multibranch pipeline in Concourse")

... or setup a pipeline for a single branch:

    $ docker run -it gofabian/flo:0 \
        flo setup-pipeline \
            --style branch \
            --target mytarget \
            --git-url https://github.com/org/repo.git \
            --branch develop

![A branch pipeline in Concourse](https://github.com/gofabian/flo/raw/main/doc/branch.png "A branch pipeline in Concourse")

## Goals

- Use the simple Drone syntax to specify build pipelines.
- Run Drone pipelines on Concourse workers.
- Pipelines are based on a `.drone.yml` file in your Git repository.
- Pipelines are generated from your Git branches.
- Pipelines update themselves from `.drone.yml`.
- No maintenance after initial setup.
- Shared workspace.

## No-Goals

- Support all aspects of Drone.

## Compatibility

Status | Drone feature
---|---
:heavy_check_mark: | Public Git repositories
:x: | Custom Git credentials/config
:heavy_check_mark: | [Simple Drone steps](https://docs.drone.io/pipeline/docker/syntax/steps/)
:x: | [Drone plugins](https://docs.drone.io/pipeline/docker/syntax/plugins/)
:x: | [Conditional steps](https://docs.drone.io/pipeline/docker/syntax/conditions/)
:x: | [Drone services](https://docs.drone.io/pipeline/docker/syntax/services/)
:heavy_check_mark: | [Docker pipelines](https://docs.drone.io/pipeline/docker/overview/)
:x: | [other pipeline types](https://docs.drone.io/pipeline/overview/)
:heavy_check_mark: | [Linux amd64 platform](https://docs.drone.io/pipeline/docker/syntax/platform/)
:x: | [Windows or MacOS platform](https://docs.drone.io/pipeline/docker/syntax/platform/)

Work in progress: https://github.com/gofabian/flo/issues


## Build binary

    $ go build flo.go
    $ ./flo -h

## Build Docker image

    $ docker build -t gofabian/flo:0 .
    $ docker run -it gofabian/flo:0 flo -h
