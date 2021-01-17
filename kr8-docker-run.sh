#! /bin/sh

set -ex

docker container run --name go-gokit-gorilla-restsvc -p 8080:8080 go-gokit-gorilla-restsvc:1.0