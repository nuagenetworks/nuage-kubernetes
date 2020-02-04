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

cd $GOPATH

for binary_name in nuagekubemon nuage-openshift-monitor
do
    docker run --rm -v `pwd`:/go -w /go golang:1.13 go build -v -o $binary_name github.com/nuagenetworks/nuage-kubernetes/nuagekubemon
    mv $binary_name $GOPATH/src/github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/
done

cd $GOPATH/src/github.com/nuagenetworks/nuage-kubernetes/nuagekubemon

sudo docker build -t nuage/monitor:${version} .
docker save nuage/monitor:${version} > nuage-monitor-docker-${version}.tar
docker rmi nuage/monitor:${version}
