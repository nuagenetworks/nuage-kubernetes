
# Multi Master Kubernetes Setup

Setup etcd
-----------

 1. Copy the contents of the etcd folder from this repository to all the etcd nodes in the cluster ( Recommended etcd cluster of 3 nodes)
 2. Fill in the etcd cluster specific information in the etcd.env file
 3. For the first node, generate the token manually using the link in the etcd.env file
 4. For the 2nd & 3rd node in the cluster, set the CURRENT_NODE accordingly
 5. Once etcd.env is populated correctl, run the setup-etcd.sh script as shown below.
 6. Status of the etcd service can be checked using `service etcd status` command
 
 ```
 cp etcd_template.env etcd.env
 _Fill in node/cluster specific configuration_
 vim etcd.env
 chmod +x setup-etcd.sh
 ./setup-etcd.sh
```
Setup the Load balancer for the Kubernetes Masters
--------------------------------------------------

1. Install haproxy on the node acting as the Load balancer
2. Modify the haproxy.cfg to balance the 3 masters with the config shown in masters directory
3. Restart the haproxy service for the configuration to take effect

Setup base for the first Kubernetes Master
--------------------------------------------

1. Copy the contents from the master directory in this repository to the first master
2. Fill in the contents in the master.env file
3. The Token field will be set to blank for the first master
4. Run init-master.sh script as shown below

```
cp master_template.env master.env
_Fill in node/cluster specific configuration_
vim master.env
chmod +x setup-master.sh
_only needed on first master server and on first configuration_
chmod +x init-master.sh
./init-master.sh
```

Configure other Kubernetes Master servers
------------------------------------------

1. Create /etc/kubernetes directory on the 2nd & 3rd Masters
2. Copy the contents from the first Master in /etc/kubernetes directory to this Master in /etc/kubernetes directory
3. Modify the master.env with the proper CURRENT_NODE value
4. Set the Token value by getting the token from the first master using `kubeadm token list` command 
5. Run the setup-master.sh script as shown below

```
_copy needed files from first master_
_copy /etc/kubernetes/*_
vim master.env
chmod +x setup-master.sh
./setup-master.sh
```

Configure the Kubernetes Nodes
------------------------------

1. Copy the contents from the node directory in this repository on all the nodes
2. Copy the token from the master.env file and populate it in the worker.env file
3. Set the correct SERVICE_SUBNET CIDR & POD_SUBNET CIDR
4. Run the setup-worker.sh script as shown below

```
_Fill in node/cluster specific configuration_
vim worker.env
chmod +x setup-worker.sh
./setup-worker.sh

```

# Pre-Installation Steps on VSD

1. Login to VSD UI as csproot and create an Enterprise "kubernetes".

2. Under the "kubernetes" Enterprise, create a user "k8s-admin" and add the user to the "Administrators" group.

   .. Note:: The steps to create the user and adding the user to a particular group can be found in the "CSP User Management" section in the "Nuage VSP User Guide".
   
3. Login to the VSD node using the CLI and execute the following command:


         /opt/vsd/ejbca/deploy/certMgmt.sh -a generate -u k8s-admin -c k8s-admin -o kubernetes -f pem -t client -s root@<k8s-master-IP>:/usr/local/
         
         Where:
         -a <action>         Action [generate|revoke|delete|renew|bscopy]
         -u <username>       End-entity Username
         -c <commonName>     Common Name. Needs to match VSD Username
         -o <organization>   Organization. Needs to match VSD Organization
         -f <format>         Certificate Format [pem|jks|p12]
         -t <type>           Certificate Type [client|server|vsc|vrs]
         -s <scpUrl>         The remote scp url. (eg. root@myhost://home/certs/)


	.. Note:: The above command generates the client certificates for the "k8s-admin" user and copies it to the /usr/local/ or any specified directory of the k8s node where Ansible is run. This certificate information is used by the nuagekubemon (nuage k8S monitor) to securely communicate with the VSD.


# Installing Nuage Kubernetes components

Clone the Nuage Ansible Git Repository
---------------------------------------

You need to have Git installed on your Ansible host machine. Perform the following tasks:

1. Access Git 
2. Setup SSH and access the master and the minion nodes, using the **ssh** command.

   .. Note:: You do not need a password to use **ssh**.

3. Clone the Ansible git repository, by entering the **git clone** command as shown in the example below and checkout the branch corresponding to the VSP version. 

.. Note:: kubernetes HA install is supported in VSP version 5.0 & above  
   
        git clone https://github.com/nuagenetworks/nuage-kubernetes.git
        git checkout origin/<vsp-version> -b <vsp-version>
        cd nuage-kubernetes/ansible

Create the configuration for Ansible
-------------------------------------

Create a inventory file for Ansible configuration in the nuage-kubernetes/ansible/inventory directory with the contents shown below.

    # Create an k8s group that contains the masters and nodes groups
    [k8s:children]
    masters
    nodes
    
    # Set variables common for all k8s hosts
    [k8s:vars]
    # SSH user, this user should allow ssh based auth without requiring a password
    ansible_ssh_user=root

    vsd_api_url=https://192.168.103.200:7443
    vsp_version=v5_0
    enterprise=kubernetes
    domain=Kubernetes
    
    vsc_active_ip=10.168.103.201
    vsc_standby_ip=10.168.103.202
    uplink_interface=eth0
    nuage_host_subnet_length=10
    nuage_cluster_network_CIDR=70.70.0.0/16

    nuage_monitor_rpm=http://172.22.61.12/Kubernetes/RPMS/x86_64/nuagekubemon-5.0.x.el7.centos.x86_64.rpm
    vrs_rpm=http://172.22.61.12/Kubernetes/RPMS/x86_64/nuage-openvswitch-5.0.x.x86_64.rpm
    plugin_rpm=http://172.22.61.12/Kubernetes/RPMS/x86_64/nuage-cni-k8s-5.0.x.el7.centos.x86_64.rpm
    
    # Complete local host path to the k8S loopback CNI plugin
    k8s_cni_loopback_plugin=/tmp/loopback
    
    # VSD user in the admin group
    vsd_user=k8s-admin
    # Complete local host path to the VSD user certificate file
    vsd_user_cert_file=/usr/local/k8s-admin.pem
    # Complete local host path to the VSD user key file
    vsd_user_key_file=/usr/local/k8s-admin-Key.pem

    # Required for Nuage Monitor REST server 
    Kubernetes_master_cluster_hostname=master.nuageKubernetes.com
    nuagekubemon_rest_server_port=9443
    
    # Optional
    nuage_interface_mtu=1460
    nuagekubemon_log_dir=/var/log/nuagekubemon
    
    # host group for masters
    [masters]
    master1.k8s.test.com
    master2.k8s.test.com
    master3.k8s.test.com
    
    # host group for nodes, includes region info
    [nodes]
    node1.k8s.test.com
    node2.k8s.test.com
    node3.k8s.test.com
    master1.k8s.test.com
    master2.k8s.test.com
    master3.k8s.test.com
      
    # host group for etcd cluster
    [etcd]
    etcd1.k8s.test.com
    etcd2.k8s.test.com
    etcd3.k8s.test.com
        
    # host group for LB
    [lb]
    lb.k8s.test.com

Modify the kube_service_addresses in the  nuage-kubernetes/ansible/inventory/group_vars/all.yml file to the service CIDR used to initialize the cluster.If any service CIDR is not specified during install, then kube_service_addresses should be updated to 10.96.0.0/12 which is the default service CIDR used by kubeadm. Also, configure the LB node as decribed in the section above


    # Kubernetes internal network for services.
    # Kubernetes services will get fake IP addresses from this range.
    # This range must not conflict with anything in your infrastructure. These
    # addresses do not need to be routable and must just be an unused block of space.
    # kube_service_addresses: 192.168.0.0/16


Installing the VSP Components 
------------------------------

This will install the following Nuage components:

 - Nuage Kubernetes Monitor (nuagekubemon) on the Kubernetes Masters
 - VRS on the Kubernetes Nodes
 - Nuage CNI plugin on the Kubernetes Nodes

This will also generate the certificates required for communication between the CNI plugin and nuagekubemon.

1. Make sure you are in the nuage-kubernetes/ansible directory. 
2. Run the following command to install the VSP components:

   
   ```
      cd nuage-kubernetes/ansible/scripts
      ./deploy-cluster.sh --tags=nuage
   ```
 
  A successful installation displays the following output:
   
       2017-07-11 22:01:49,891 p=16545 u=root |  PLAY RECAP *********************************************************************
       2017-07-11 22:01:49,892 p=16545 u=root |  localhost : ok=20   changed=0   unreachable=0  failed=0
       2017-07-11 22:01:49,892 p=16545 u=root |  master1.k8s.test.com : ok=247  changed=22  unreachable=0  failed=0
       2017-07-11 22:01:49,893 p=16545 u=root |  master2.k8s.test.com : ok=247  changed=22  unreachable=0  failed=0
       2017-07-11 22:01:49,894 p=16545 u=root |  master3.k8s.test.com : ok=247  changed=22  unreachable=0  failed=0
       2017-07-11 22:01:49,895 p=16545 u=root |  node1.k8s.test.com : ok=111  changed=21  unreachable=0  failed=0
       2017-07-11 22:01:49,896 p=16545 u=root |  node2.k8s.test.com : ok=111  changed=21  unreachable=0  failed=0
       2017-07-11 22:01:49,897 p=16545 u=root |  node3.k8s.test.com : ok=111  changed=21  unreachable=0  failed=0
       2017-07-11 22:01:49,895 p=16545 u=root |  etcd1.k8s.test.com : ok=111  changed=21  unreachable=0  failed=0
       2017-07-11 22:01:49,896 p=16545 u=root |  etcd2.k8s.test.com : ok=111  changed=21  unreachable=0  failed=0
       2017-07-11 22:01:49,897 p=16545 u=root |  etcd3.k8s.test.com : ok=111  changed=21  unreachable=0  failed=0
       
3. Verify that the Master-Node connectivity is up and all nodes are running using the following command on the master:

      `kubectl get nodes`       
      



