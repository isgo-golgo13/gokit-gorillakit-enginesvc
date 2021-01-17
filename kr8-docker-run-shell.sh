#! /bin/sh

set -ex

# The following docker container sub-command ONLY works if container is running

docker container exec -it go-gokit-gorilla-restsvc /bin/sh  