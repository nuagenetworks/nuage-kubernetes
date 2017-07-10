#!/bin/bash
# Setup Worker

# k8s requirement for RHEL
if ! grep -q net.bridge.bridge-nf-call-ip /etc/sysctl.conf; then
  for x in 6tables tables; do
    echo "net.bridge.bridge-nf-call-ip$x = 1" >> /etc/sysctl.conf
  done
fi
sysctl -p

# Load per node config
source worker.env

# Check envs
if [[ -z "$TOKEN"  ]] || [[ -z "$APISERVER_LB"  ]] || [[ -z "$SERVICE_SUBNET"  ]] || [[ -z "$POD_SUBNET"  ]]; then
  echo "Not all variables set in worker.env"
  exit 1
fi

# Config
DOCKER_VERSION="1.12.6"
KUBE_VERSION="v1.6.6"

# Use proxy to acces internets
#export http_proxy="http://proxy-dc.example.com:80"
#export https_proxy="http://proxy-dc.example.com:80"

# Install
#TODO: use package mirror and switch to internal package repository
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg
        https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
EOF

cat <<EOF > /etc/yum.repos.d/docker.repo
[dockerrepo]
name=Docker Repository
baseurl=https://yum.dockerproject.org/repo/main/centos/7/
enabled=1
gpgcheck=1
gpgkey=https://yum.dockerproject.org/gpg
EOF

yum update -y
yum install -y docker-engine-$DOCKER_VERSION yum-versionlock
yum versionlock add docker-engine

yum install -y kubelet-1.6.6 kubeadm-1.6.6 kubectl-1.6.6 kubernetes-cni

#unset http_proxy
#unset https_proxy

IP="${SERVICE_SUBNET%\/*}"
DNS="${IP%\.*}.$((${IP##*\.}+10))"
sed -i -e "s/--cgroup-driver=systemd/--cgroup-driver=cgroupfs/g" /etc/systemd/system/kubelet.service.d/10-kubeadm.conf
sed -i -e "s/--cluster-dns=10.96.0.10/--cluster-dns=$DNS/g" /etc/systemd/system/kubelet.service.d/10-kubeadm.conf

mkdir -p /etc/systemd/system/docker.service.d
#cat <<EOF > /etc/systemd/system/docker.service.d/http-proxy.conf
#[Service]
#Environment="HTTP_PROXY=http://proxy-dc.example.com:80"
#Environment="HTTPS_PROXY=http://proxy-dc.example.com:80"
#Environment="NO_PROXY=localhost,127.0.0.1"
#EOF

systemctl enable docker && systemctl start docker
systemctl enable kubelet && systemctl start kubelet

# Disable SELinux to enable host filesystem access and networking
setenforce 0

# Initialize
kubeadm join --token $TOKEN "$APISERVER_LB:6443"

sed -i -e "/    server: /c\    server: https://$APISERVER_LB:6443/" /etc/kubernetes/kubelet.conf
