# nuage-kubernetes

Kubernetes Integration Related Artifacts
 - Nuage ansible repository
 - Kubernetes HA deployment guide
 - Nuage Install Documentation 
 - Manual installation instructions
 - Demo scripts
 - nuagekubemon monitor source
 - nuage-k8s-plugin source
 - Sample jsons


Getting the right Nuage release

Follow this release pattern and you can't go wrong:

 |   Nuage-VSP-version    |    Kubernetes-version	     |    Nuage-ansible-branch    |
 | -----------------------|----------------------------|----------------------------|
 |       4.0	             |         1.5.4	             |         4.0                |
 |       5.0              |         1.6.4              |         5.0                |
 |      5.1.1             |  HA 1.6.6/Standalone 1.7.4 |        tags/v5.1.1-1       |
 |      5.1.2             |  HA 1.6.6/Standalone 1.7.4 |        tags/v5.1.2-1       |
 
 Kubernetes-version implies version of kubeadm, kubelet and kubectl componenets. In case, you want to install a particular kubernetes release, you can do
 `yum install -y kubeadm-1.6.6 kubelet-1.6.6 kubectl-1.6.6 docker kubernetes-cni`
 
 This should install the required version of kubernetes components along with the required version of docker.
 
 ..Note: Ansible version supported is 2.3.0.0 for Nuage ansible playbook with Kubernetes
