#!/bin/bash

set -e

if [ -z ${GOPATH} ]; then
    echo "\"GOPATH\" environmental variable is not set";
    exit 1
fi

if [ -z ${version} ]; then
    echo "\"version\" environmental variable is not set";
    exit 1
fi

cd $GOPATH/src/github.com/nuagenetworks/nuage-kubernetes/nuagekubemon
go build -o nuagekubemon
go build -o nuage-openshift-monitor
sudo docker build -t nuage/master:${version} .
docker save nuage/master:${version} > nuage-master-docker.tar
