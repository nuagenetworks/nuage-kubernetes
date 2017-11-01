
# Multi Master Kubernetes Setup

.. Note:: The steps documented to install the Multi Master Kubernetes Setup work with Kubernetes v1.6.6 or below

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


         /opt/vsd/ejbca/deploy/certMgmt.sh -a generate -u k8s-admin -c k8s-admin -o kubernetes -f pem -t client -s root@<ansible-node-IP>:/usr/local/
         
         Where:
         -a <action>         Action [generate|revoke|delete|renew|bscopy]
         -u <username>       End-entity Username
         -c <commonName>     Common Name. Needs to match VSD Username
         -o <organization>   Organization. Needs to match VSD Organization
         -f <format>         Certificate Format [pem|jks|p12]
         -t <type>           Certificate Type [client|server|vsc|vrs]
         -s <scpUrl>         The remote scp url. (eg. root@myhost://home/certs/)


	.. Note:: The above command generates the client certificates for the "k8s-admin" user and copies it to the /usr/local/ or any specified directory of the k8s node where Ansible is run and this file path is also specified in the nodes file. This certificate information is used by the nuagekubemon (nuage k8S monitor) to securely communicate with the VSD.

4. To complete the steps provided in the Kubeadm installer guide, go `here <https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/>`_. 

  .. Note:: Kubectl needs the kube config to be copied to a specific location after the master is initialized using kubeadm init. Ansible scripts used to install Nuage components also rely on kubectl being available to the ansible user. In order to achieve that, execute the following commands on the master node after kubeadm init:
    ::
          sudo mkdir -p ~/.kube
          sudo cp /etc/kubernetes/admin.conf ~/.kube/config
          sudo chown $(id -u):$(id -g) $HOME/.kube/config


As a part of the Kubernetes Ansible Installation for Nuage, Kubernetes DaemonSet will be used for installation of Nuage components. DaemonSet will be responsible for installation & maintenance of containerized monitor (nuagekubemon) and containerized CNI plugin with containerized Nuage VRS on master and slave nodes respectively.

.. Note:: All Nuage services like nuagekubemon, CNI plugin and VRS will be operating as DaemonSet pods on master and slave nodes. This is the recommended method of installing Nuage components with Kubernetes.

# Installing Nuage Kubernetes components using Daemonsets

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
   
.. Note::  Post Nuage release 5.1.1, rpm based install is not supported using ansible. Daemonsets is the recommended mode of installing Nuage components

4. Load the following docker images on your master node:

::

    docker load -i nuage-master-docker.tar
    docker load -i nuage-cni-docker.tar
    docker load -i nuage-vrs-docker.tar
    docker load -i nuage-infra-docker.tar

5. Load the following docker images on your slave nodes:

::

    docker load -i nuage-cni-docker.tar
    docker load -i nuage-vrs-docker.tar
    docker load -i nuage-infra-docker.tar

6. Update the following parameters in ConfigMap section of **nuage-kubernetes/ansible/roles/nuage-daemonset/files/nuage-master-config-daemonset.yaml** file as per your environment configuration:

::

      # This will generate the required Nuage monitor configuration
      # on master nodes
      monitor_yaml_config: |
          kubeConfig: /usr/share/nuagekubemon/nuage.kubeconfig

          masterConfig: /usr/share/nuagekubemon/net-config.yaml
          # URL of the VSD Architect
          vsdApiUrl: https://xmpp.example.com:7443
          # API version to query against
          vspVersion: v5_0
          # Name of the enterprise in which pods will reside
          enterpriseName: kubernetes
          # Name of the domain in which pods will reside
          domainName: kubernetes
          # VSD generated user certificate file location on master node
          userCertificateFile: /usr/share/nuagekubemon/k8s-admin.pem
          # VSD generated user key file location on master node
          userKeyFile: /usr/share/nuagekubemon/k8s-admin-Key.pem
          # Location where logs should be saved
          log_dir: /var/log/nuagekubemon
          # Monitor rest server paramters
          # Logging level for the nuage monitor
          # allowed options are: 0 => INFO, 1 => WARNING, 2 => ERROR, 3 => FATAL
          logLevel: 0
          # Parameters related to the nuage monitor REST server
          nuageMonServer:
              URL: 0.0.0.0:9443
              certificateDirectory: /usr/share/nuagekubemon

      # This will generate the required Nuage network configuration
      # on master nodes
      net_yaml_config: |
          networkConfig:
            clusterNetworkCIDR: 70.70.0.0/16
            serviceNetworkCIDR: 192.168.0.0/16
            hostSubnetLength: 8

Make sure to set the etcd config correctly if there is an external etcd cluster. If the etcd cluster is not using TLS certificates, do not set the ca, certFile & keyFile parameters. Also, if etcd is running locally on the master, use the localhost IP as shown below. If the etcd cluster is setup using FQDN, set the URL to the FQDN hostname. Also, make sure to check the protocol for the etcd cluster and set http or https accordingly.

# etcd config required for HA etcdClientConfig:

::

        ca: "" 
        certFile: "" 
        keyFile: "" 
        urls:
            - http://127.0.0.1:2379

Set the parameter to 1 in order to allow nuagekubemon to automagically create a new subnet when the existing subnet gets depleted. Threshold for new subnet creation is set to 70% utilization of namespace/zone subnet pool. It will also delete additional subnets if the namespace subnet pool utilization falls below 25%

::


         #auto scale subnets feature 
	 # 0 => disabled(default) 
	 # 1 => enabled 
	 autoScaleSubnets: 1
         #This will generate the required Nuage network configuration 
	 # on master nodes 
	 net_yaml_config: |
	     networkConfig:
	       clusterNetworkCIDR: 70.70.0.0/16 
	       serviceNetworkCIDR: 192.168.0.0/16 
	       hostSubnetLength: 8


Make sure the **image** parameter is correctly set to the Nuagekubemon docker images version pre-loaded on master nodes:

::

      containers:
        # This container configures Nuage Master node
        - name: install-nuage-master-config
          image: nuage/master:<nuage-version>

7. Update the following parameters in **nuage-kubernetes/ansible/roles/nuage-daemonset/files/nuage-node-config-daemonset.yaml** file as per your environment configuration:

::

        # This will generate the required Nuage vsp-k8s.yaml
        # config on each slave node
        plugin_yaml_config: |
            # Path to Nuage kubeconfig
            kubeConfig: /usr/share/vsp-k8s/nuage.kubeconfig
            # Name of the enterprise in which pods will reside
            enterpriseName: kubernetes
            # Name of the domain in which pods will reside
            domainName: kubernetes
            # Name of the VSD user in admin group
            vsdUser: k8s-admin
            # REST server URL
            nuageMonRestServer: https://<Load-balancer IP>:9443
            # Bridge name for the docker bridge
            dockerBridgeName: docker0
            # Certificate for connecting to the kubemon REST API
            nuageMonClientCert: /usr/share/vsp-k8s/nuageMonClient.crt
            # Key to the certificate in restClientCert
            nuageMonClientKey: /usr/share/vsp-k8s/nuageMonClient.key
            # CA certificate for verifying the master's rest server
            nuageMonServerCA: /usr/share/vsp-k8s/nuageMonCA.crt
            # Nuage vport mtu size
            interfaceMTU: 1460
            # Service CIDR
            serviceCIDR: 192.168.0.0/16
            # Logging level for the plugin
            # allowed options are: "dbg", "info", "warn", "err", "emer", "off"
            logLevel: dbg
	    
	# This will generate the required Nuage CNI yaml configuration
        cni_yaml_config: |
            vrsendpoint: "/var/run/openvswitch/db.sock"
            vrsbridge: "alubr0"
            monitorinterval: 60
            cniversion: 0.2.0
            loglevel: "debug"
            portresolvetimer: 60
            logfilesize: 1
            vrsconnectionchecktimer: 180
            mtu: 1460
            staleentrytimeout: 600

Update the following environment variables in the DaemonSet section for **nuage-cni-ds** with the value set in clusterNetworkCIDR in the nuage-master-config-daemonset.yaml above    

         # Nuage cluster network CIDR for iptables configuration
            - name: NUAGE_CLUSTER_NW_CIDR
              value: "70.70.0.0/16"


Update the following environment variables in DaemonSet section for **nuage-vrs-ds** with Active and Standby Nuage VSC IP addresses for containerized Nuage VRS and NUAGE_K8S_SERVICE_IPV4_SUBNET with the value for serviceNetworkCIDR set in nuage-master-config-daemonset.yaml above

::

      env:
        # Configure parameters for VRS openvswitch file
        - name: NUAGE_ACTIVE_CONTROLLER
          value: "10.100.100.100"
        - name: NUAGE_STANDBY_CONTROLLER
          value: "10.100.100.101"
        - name: NUAGE_PLATFORM
          value: '"kvm, k8s"'
        - name: NUAGE_K8S_SERVICE_IPV4_SUBNET
          value: '192.168.0.0\/16'


Make sure the **image** parameter is correctly set to the Nuage VRS and CNI docker images version pre-loaded on slave nodes:

::

      containers:
        # This container installs Nuage VRS running as a
        # container on each worker node
        - name: install-nuage-vrs
          image: nuage/vrs:<nuage-version>

      containers:
        # This container installs Nuage CNI binaries
        # and CNI network config file on each node.
        - name: install-nuage-cni
          image: nuage/cni:<nuage-version>

Update the **image** parameter in **nuage-kubernetes/ansible/roles/nuage-daemonset/files/nuage-infra-pod-config-daemonset.yaml** file and make sure that it is correctly set to the Nuage infra pod image version pre-loaded on the slave nodes:

::

      containers:
        # This container spawns a Nuage Infra pod
        # on each worker node
        - name: install-nuage-infra
          image: nuage/infra:<nuage-version>


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

    nuage_cluster_network_CIDR=70.70.0.0/16
       
    # Complete local host path to the VSD user certificate file
    vsd_user_cert_file=/usr/local/k8s-admin.pem
    # Complete local host path to the VSD user key file
    vsd_user_key_file=/usr/local/k8s-admin-Key.pem

    # Required for Nuage Monitor REST server. In case of HA, this should be set to the hostname of the LB node
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
      



