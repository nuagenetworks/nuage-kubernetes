yum update kubelet kubernetes-cni docker-engine
run kubeadm init --config on clean machine
copy generated manifest files to old master nodes
systemctl restart kubelet (should pickup new manifest files)
