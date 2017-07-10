#!/bin/bash

systemctl stop docker
\rm -r /var/lib/etcd/*
kubeadm reset
yum remove -y kubernetes-cni kubeadm kubelet kubectl docker-engine
ip link set docker0 down
ip link set cni0 down
ip link set flannel.1 down
ip link delete flannel.1
ip link delete cni0
ip link delete docker0
