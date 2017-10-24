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
    docker run --rm -v `pwd`:/go -w /go golang:1.8 go build -v -o $binary_name github.com/nuagenetworks/nuage-kubernetes/nuagekubemon
    mv $binary_name $GOPATH/src/github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/
done

cd $GOPATH/src/github.com/nuagenetworks/nuage-kubernetes/nuagekubemon

sudo docker build -t nuage/master:${version} .
docker save nuage/master:${version} > nuage-master-docker-${version}.tar
