#!/bin/sh
set -e

cd $(dirname $0)

docker build -t gofabian/flo:0 .
docker push gofabian/flo:0
