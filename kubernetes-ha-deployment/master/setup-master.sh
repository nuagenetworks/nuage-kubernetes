#!/bin/bash
# Setup Master

# Config
DOCKER_VERSION="1.12.6"
KUBE_VERSION="v1.6.4"

# Load per node config
source master.env

# Check envs
if [[ -z "$ETCD_1"  ]] || [[ -z "$ETCD_2"  ]] || [[ -z "$ETCD_3"  ]] || [[ -z "$APISERVER_1"  ]] || [[ -z "$APISERVER_2"  ]] || [[ -z "$APISERVER_3"  ]] || [[ -z "$APISERVER_LB"  ]] || [[ -z "$SERVICE_SUBNET"  ]] || [[ -z "$POD_SUBNET"  ]] || [[ -z "$CURRENT_NODE"  ]]; then
  echo "Not all variables set in master.env"
  exit 1
fi

# Check if init node or not
if [[ $CURRENT_NODE == "2" ]] || [[ $CURRENT_NODE == "3" ]]; then
  if [[ -z "$( find /etc/kubernetes -maxdepth 1 -mmin -10 2>/dev/null )" ]]; then
    echo "Files not copied from initilization master. Copy over files!"
    echo "Copy /etc/kubernetes/* to this node."
    exit 1
  elif [[ -z "$TOKEN"  ]]; then
    echo "TOKEN not set."
    exit 1
  fi
fi

# Get current hostname
APISERVER_CURRENT=APISERVER_$CURRENT_NODE
rx='([1-9]?[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])'
if [[ $APISERVER_1 =~ ^$rx\.$rx\.$rx\.$rx$ ]]; then
  echo "APISERVER needs to be a hostname"
  exit 1
fi

# k8s requirement for RHEL
if ! grep -q net.bridge.bridge-nf-call-ip /etc/sysctl.conf; then
  for x in 6tables tables; do
    echo "net.bridge.bridge-nf-call-ip$x = 1" >> /etc/sysctl.conf
  done
fi
sysctl -p

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

mkdir -p /etc/systemd/system/docker.service.d
#cat <<EOF > /etc/systemd/system/docker.service.d/http-proxy.conf
#[Service]
#Environment="HTTP_PROXY=http://proxy-dc.example.com:80"
#Environment="HTTPS_PROXY=http://proxy-dc.example.com:80"
#Environment="NO_PROXY=localhost,127.0.0.1"
#EOF

cat <<EOF > /etc/yum.repos.d/docker.repo
[dockerrepo]
name=Docker Repository
baseurl=https://yum.dockerproject.org/repo/main/centos/7/
enabled=1
gpgcheck=1
gpgkey=https://yum.dockerproject.org/gpg
EOF

# Update
yum update -y
yum install -y docker-engine-$DOCKER_VERSION yum-versionlock
yum versionlock add docker-engine

yum install -y bind-utils kubelet kubeadm kubectl kubernetes-cni 

#unset http_proxy
#unset https_proxy

# Pick DNS service IP from Service Subnet
IP="${SERVICE_SUBNET%\/*}"
DNS="${IP%\.*}.$((${IP##*\.}+10))"

sed -i -e "s/--cgroup-driver=systemd/--cgroup-driver=cgroupfs/g" /etc/systemd/system/kubelet.service.d/10-kubeadm.conf
sed -i -e "s/--cluster-dns=10.96.0.10/--cluster-dns=$DNS/g" /etc/systemd/system/kubelet.service.d/10-kubeadm.conf
sed -i -e "s/SELINUX=\(.*\)/SELINUX=disabled/g" /etc/sysconfig/selinux

# Replace current hostname
sed -i -e "s/$APISERVER_1/${!APISERVER_CURRENT}/g" /etc/kubernetes/kubelet.conf

# Copy admin config
mkdir -p ~/.kube
cp /etc/kubernetes/admin.conf ~/.kube/config
chown "$(id -nu)": ~/.kube/config

systemctl daemon-reload
systemctl enable docker && systemctl start docker
systemctl enable kubelet && systemctl start kubelet

# Disable SELinux to enable host filesystem access and networking
setenforce 0
