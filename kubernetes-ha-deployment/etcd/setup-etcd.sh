#!/bin/bash
# Setup etcd

# Config
ETCD_VERSION="v3.0.17"
DOCKER_VERSION="1.12.6"

# Load per node config
source etcd.env

# Check envs
if [[ -z "$ETCD_TOKEN"  ]] || [[ -z "$ETCD_1"  ]] || [[ -z "$ETCD_2"  ]] || [[ -z "$ETCD_3"  ]] || [[ -z "$CURRENT_NODE"  ]]
then
  echo "Not all variables set in etcd.env"
  exit 1
fi

declare ETCD_NAME_1="etcd-node-1"
declare ETCD_NAME_2="etcd-node-2"
declare ETCD_NAME_3="etcd-node-3"

ETCD_THIS_NAME="etcd-node-$CURRENT_NODE"
ETCD_THIS_IP=ETCD_$CURRENT_NODE

# Use proxy to acces internets
#export http_proxy="http://proxy-dc.example.com:80"
#export https_proxy="http://proxy-dc.example.com:80"

mkdir -p /etc/systemd/system/docker.service.d
#cat <<EOF > /etc/systemd/system/docker.service.d/http-proxy.conf
#[Service]
#Environment="HTTP_PROXY=http://proxy-dc.example.com:80"
#Environment="HTTPS_PROXY=http://proxy-dc.example.com:80"
#Environment="NO_PROXY=localhost,127.0.0.1"
#EOF

# Install
cat <<EOF > /etc/yum.repos.d/docker.repo
[dockerrepo]
name=Docker Repository
baseurl=https://yum.dockerproject.org/repo/main/centos/7/
enabled=1
gpgcheck=1
gpgkey=https://yum.dockerproject.org/gpg
EOF

#cp CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo
yum makecache
yum install -y docker-engine-$DOCKER_VERSION yum-versionlock
yum versionlock add docker-engine

#unset http_proxy
#unset https_proxy

# Disable SELinux to enable host filesystem access and networking
setenforce 0

cp etcd_template.service etcd.service
sed -i -e "s/<etcd-token>/$ETCD_TOKEN/g" etcd.service
sed -i -e "s/<etcd-IP-1>/$ETCD_1/g" etcd.service
sed -i -e "s/<etcd-IP-2>/$ETCD_2/g" etcd.service
sed -i -e "s/<etcd-IP-3>/$ETCD_3/g" etcd.service
sed -i -e "s/<etcd-name-1>/$ETCD_NAME_1/g" etcd.service
sed -i -e "s/<etcd-name-2>/$ETCD_NAME_2/g" etcd.service
sed -i -e "s/<etcd-name-3>/$ETCD_NAME_3/g" etcd.service
sed -i -e "s/<etcd-this-name>/$ETCD_THIS_NAME/g" etcd.service
sed -i -e "s/<etcd-this-IP>/${!ETCD_THIS_IP}/g" etcd.service
sed -i -e "s/<etcd-version>/$ETCD_VERSION/g" etcd.service

mv etcd.service /etc/systemd/system/

systemctl daemon-reload
systemctl enable etcd && systemctl restart etcd

#TODO: secure down etcd via encryted connection and authentication for clients (apiserver)
