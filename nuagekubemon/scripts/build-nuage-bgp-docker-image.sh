#!/bin/bash

set -e

if [ -z ${version} ]; then
    echo "\"version\" environmental variable is not set";
    exit 1
fi

sudo docker build -t nuage/bgp:${version} -f nuage-bgp/Dockerfile nuage-bgp/
docker save nuage/bgp:${version} > nuage-bgp-docker-${version}.tar
docker rmi nuage/bgp:${version}
