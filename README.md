# Flo

## Story

Did you ever ...

## Goals

## No-Goals

## Restrictions

Some of these restrictions may be removed in future:

- Linux amd64
- only Docker pipelines (`type: docker`)
- No Drone Plugins (https://docs.drone.io/pipeline/docker/syntax/plugins/)
- No Drone Services (https://docs.drone.io/pipeline/docker/syntax/services/)


## Build local binary

    $ go build flo.go
    $ ./flo -h

## Build docker image

    $ docker build -t gofabian/flo:0 .
    $ docker run -it gofabian/flo:0 flo -h
