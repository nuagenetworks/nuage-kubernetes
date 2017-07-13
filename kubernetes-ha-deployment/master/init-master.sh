#!/bin/bash
# Init Master

source setup-master.sh

# Check for already initialize master
if [ -f /etc/kubernetes/admin.conf ]; then
  echo "Rerunning init master on top of an initialized master. Exiting!"
  exit 1
fi

# Generate IPs
APISERVER_IP_1=$(nslookup $APISERVER_1 | awk '/^Address: / { print $2 ; exit }')
APISERVER_IP_2=$(nslookup $APISERVER_2 | awk '/^Address: / { print $2 ; exit }')
APISERVER_IP_3=$(nslookup $APISERVER_3 | awk '/^Address: / { print $2 ; exit }')
APISERVER_IP_LB=$(nslookup $APISERVER_LB | awk '/^Address: / { print $2 ; exit }')

if [[ -z "$TOKEN" ]]; then
  TOKEN=$(kubeadm token generate)
fi

# Get base config
cp kubeadm-master_template.conf kubeadm-master.conf

# Generate specific config from template
sed -i -e "s/<etcd-IP-1>/$ETCD_1/g" kubeadm-master.conf
sed -i -e "s/<etcd-IP-2>/$ETCD_2/g" kubeadm-master.conf
sed -i -e "s/<etcd-IP-3>/$ETCD_3/g" kubeadm-master.conf
sed -i -e "s/<apiserver-HN-1>/$APISERVER_1/g" kubeadm-master.conf
sed -i -e "s/<apiserver-HN-2>/$APISERVER_2/g" kubeadm-master.conf
sed -i -e "s/<apiserver-HN-3>/$APISERVER_3/g" kubeadm-master.conf
sed -i -e "s/<apiserver-IP-1>/$APISERVER_IP_1/g" kubeadm-master.conf
sed -i -e "s/<apiserver-IP-2>/$APISERVER_IP_2/g" kubeadm-master.conf
sed -i -e "s/<apiserver-IP-3>/$APISERVER_IP_3/g" kubeadm-master.conf
sed -i -e "s/<apiserver-LB>/$APISERVER_LB/g" kubeadm-master.conf
sed -i -e "s/<apiserver-IP-LB>/$APISERVER_IP_LB/g" kubeadm-master.conf
sed -i -e "s/<token>/$TOKEN/g" kubeadm-master.conf
sed -i -e "s/<kube-version>/$KUBE_VERSION/g" kubeadm-master.conf

# Custom subnets
sed -i -e "s~<pod-subnet>~$POD_SUBNET~g" kubeadm-master.conf
sed -i -e "s~<service-subnet>~$SERVICE_SUBNET~g" kubeadm-master.conf

# Use custom subnet in flannel
sed -i -e "s~10.244.0.0/16~$POD_SUBNET~g" kube-flannel.yml

# Master initialization (Only on first setup! On first server!)
kubeadm init --config kubeadm-master.conf

# Replace LB IP with DNS
sed -i -e "/    server: /c\    server: https://$APISERVER_LB:6443/" /etc/kubernetes/admin.conf
sed -i -e "/    server: /c\    server: https://$APISERVER_LB:6443/" /etc/kubernetes/kubelet.conf
sed -i -e "/    server: /c\    server: https://$APISERVER_LB:6443/" /etc/kubernetes/scheduler.conf
sed -i -e "/    server: /c\    server: https://$APISERVER_LB:6443/" /etc/kubernetes/controller-manager.conf

sed -i -e "s/$APISERVER_1/${!APISERVER_CURRENT}/" /etc/kubernetes/kubelet.conf

# Copy admin config
mkdir -p ~/.kube
cp /etc/kubernetes/admin.conf ~/.kube/config
chown "$(id -nu)": ~/.kube/config

sleep 20

# Setup cluster networking
#kubectl create -f kube-flannel.yml
#kubectl create -f kube-flannel-rbac.yml
