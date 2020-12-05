# Flo

## Restrictions

Some of these restrictions may be removed in future:

- only Docker pipelines (`type: docker`)
- Linux amd64
- No Drone Plugins (https://docs.drone.io/pipeline/docker/syntax/plugins/)
- No Drone Services (https://docs.drone.io/pipeline/docker/syntax/services/)
- `sh` required in Docker image for executing multiple commands in single step


## Build local binary

Build:

    $ go build flo.go

Run:

    $ ./flo -h

## Build docker image

Build:

    $ docker build -t gofabian/flo:0 .

Run:

    $ docker run -it gofabian/flo:0

... or ...

    $ docker run -it gofabian/flo:0 flo -h
