.. _Kubernetes-4-installation:

.. include:: ../lib/doc-includes/VSDA-icons.inc

===========================
Kubernetes Installation
===========================

.. contents::
   :local:
   :depth: 3

Kubernetes Ansible Installation
====================================

The following sections provide instructions for installing Kubernetes using Kubeadm followed by insertion of Nuage VSP components.

.. Note:: If you have already setup the Kubernetes cluster using Kubeadm, you can skip the Pre-Installation steps and go to the `Install Git Repository`_ section.

You only need to install Ansible once on a machine and it can manage the master and all the other remote nodes. Ansible manages nodes using the SSH protocol, therefore it is required that SSH is setup so that the master and nodes are accessible from the host running the Ansible scripts.

Kubernetes DaemonSet for Nuage installation
===========================================

As a part of the Kubernetes Ansible Installation for Nuage, Kubernetes DaemonSet will be used for installation of Nuage components. DaemonSet will be responsible for installation & maintenance of containerized monitor (nuagekubemon) and containerized CNI plugin with containerized Nuage VRS on master and slave nodes respectively.

.. Note:: All Nuage services like nuagekubemon, CNI plugin and VRS will be operating as DaemonSet pods on master and slave nodes. This is the recommended method of installing Nuage components with Kubernetes.

Pre-Installation Steps in VSD
-----------------------------
1. Login to VSD UI as csproot and create an Enterprise "kubernetes".

2. Under the "kubernetes" Enterprise, create a user "k8s-admin" and add the user to the "Administrators" group.

   .. Note:: The steps to create the user and adding the user to a particular group can be found in the "CSP User Management" section in the "Nuage VSP User Guide."

3. Login to the VSD node using the CLI and execute the following command:

    ::

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

4. To complete the steps provided in the Kubeadm installer guide, go `here <https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/>`_. 

  .. Note:: Kubectl needs the kube config to be copied to a specific location after the master is initialized using kubeadm init. Ansible scripts used to install Nuage components also rely on kubectl being available to the ansible user. In order to achieve that, execute the following commands on the master node after kubeadm init:
    ::
          mkdir -p ~/.kube
          cp /etc/kubernetes/admin.conf ~/.kube/config
          chown "$(id -nu)": ~/.kube/config
  
5. Follow the steps 1, 2 & 4 of the document provided in the above link. For the pod network, install Nuage using the Ansible installer mentioned below. 

  .. Note:: By default, Kubeadm uses 10.96.0.0/12 as the service CIDR. Make sure this service CIDR does not overlap with your existing underlay network CIDR. If it does then, please run step 2 from the above guide as follows so as to change the service CIDR:
         kubeadm init --service-cidr=192.168.0.0/16 --kubernetes-version <k8s-version>

  .. Note:: With new version of Kubernetes, pod network is required by default and nodes won't move to ready state until pod network is installed. So, after kubeadm join is done on all nodes, your kubectl get nodes can be as follows:

  ::

          kubectl get nodes
          NAME                            STATUS     AGE       VERSION
          ovs-1.test.nuagenetworks.com    NotReady   3h        v1.7.0
          ovs-10.test.nuagenetworks.com   NotReady   3h        v1.7.0
          ovs-2.test.nuagenetworks.com    NotReady   3h        v1.7.0
          ovs-3.test.nuagenetworks.com    NotReady   3h        v1.7.0
          ovs-4.test.nuagenetworks.com    NotReady   3h        v1.7.0
          ovs-5.test.nuagenetworks.com    NotReady   3h        v1.7.0
 
6. Update the cluster-dns in the 10-kubeadm.conf file on all nodes and master as follows:
  
    On the Node:
    
    ::
        
    	cat /etc/systemd/system/kubelet.service.d/10-kubeadm.conf 
    	[Service]
        Environment="KUBELET_KUBECONFIG_ARGS=--kubeconfig=/etc/kubernetes/kubelet.conf --require-kubeconfig=true"
        Environment="KUBELET_SYSTEM_PODS_ARGS=--pod-manifest-path=/etc/kubernetes/manifests --allow-privileged=true"
        Environment="KUBELET_DNS_ARGS=--cluster-dns=192.168.0.10 --cluster-domain=cluster.local"
        Environment="KUBELET_AUTHZ_ARGS=--authorization-mode=Webhook --client-ca-file=/etc/kubernetes/pki/ca.crt"
        Environment="KUBELET_CGROUP_ARGS=--cgroup-driver=systemd" 
        ExecStart=
        Environment="KUBELET_NETWORK_ARGS=--network-plugin=cni --cni-bin-dir=/usr/bin/ --make-iptables-util-chains=false"
        ExecStart=/usr/bin/kubelet $KUBELET_KUBECONFIG_ARGS $KUBELET_SYSTEM_PODS_ARGS $KUBELET_NETWORK_ARGS $KUBELET_DNS_ARGS $KUBELET_EXTRA_ARGS $KUBELET_CGROUP_ARGS

    
    On the Master:
    
    ::
    
    	cat /etc/systemd/system/kubelet.service.d/10-kubeadm.conf 
        [Service]
        Environment="KUBELET_KUBECONFIG_ARGS=--kubeconfig=/etc/kubernetes/kubelet.conf --require-kubeconfig=true"
        Environment="KUBELET_SYSTEM_PODS_ARGS=--pod-manifest-path=/etc/kubernetes/manifests --allow-privileged=true"
        Environment="KUBELET_DNS_ARGS=--cluster-dns=192.168.0.10 --cluster-domain=cluster.local"
        Environment="KUBELET_AUTHZ_ARGS=--authorization-mode=Webhook --client-ca-file=/etc/kubernetes/pki/ca.crt"
        Environment="KUBELET_CGROUP_ARGS=--cgroup-driver=systemd" 
        ExecStart=
        Environment="KUBELET_NETWORK_ARGS=--network-plugin=cni --cni-bin-dir=/usr/bin/ --make-iptables-util-chains=false"
        ExecStart=/usr/bin/kubelet $KUBELET_KUBECONFIG_ARGS $KUBELET_SYSTEM_PODS_ARGS $KUBELET_NETWORK_ARGS $KUBELET_DNS_ARGS $KUBELET_EXTRA_ARGS $KUBELET_CGROUP_ARGS

  .. Note:: With new version of Kubernetes, KUBELET_CGROUP_ARGS are required to be added to ExecStart or else it will cause the kubelet to fail. The cgroup-driver can be set to cgroupfs or systemd depending on the driver used for docker installation. By default, its set to systemd
    
   This service CIDR also gets updated in another file as explained in the "Installation for a Single Master" section.

Install Git Repository
-----------------------

You need to have Git installed on your Ansible host machine. Perform the following tasks:

1. Access Git 
2. Setup SSH and access the master and the minion nodes, using the **ssh** command.
3. Clone the Ansible git repository, by entering the **git clone** command as shown in the example below and checkout the tag **nuage-kubernetes-<version>** corresponding to the VSP version. 

.. Note:: For the required versions, see the `Requirements <kubernetes-1-overview.html#requirements>`_ section in the "Overview" chapter of this guide.

   ::
   
        git clone https://github.com/nuagenetworks/nuage-kubernetes.git
        git checkout tags/nuage-kubernetes-5.1.1-1
        cd nuage-kubernetes/ansible

4. Load the following docker images on your master node:

::

    docker load -i nuage-master-docker.tar
    docker load -i nuage-cni-docker.tar
    docker load -i nuage-vrs-docker.tar

5. Load the following docker images on your slave nodes:

::

    docker load -i nuage-cni-docker.tar
    docker load -i nuage-vrs-docker.tar

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

Make sure the **image** parameter is correctly set to the Nuagekubemon docker images version pre-loaded on master nodes:

::

      containers:
        # This container configures Nuage Master node
        - name: install-nuage-master-config
          image: nuage/master:5.1.1

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
            nuageMonRestServer: https://<Master-Node-IP or hostname>:9443
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
          image: nuage/vrs:5.1.1

      containers:
        # This container installs Nuage CNI binaries
        # and CNI network config file on each node.
        - name: install-nuage-cni
          image: nuage/cni:5.1.1


Installation for a Single Master 
-----------------------------------

Create a inventory file for Ansible configuration for a single master in the nuage-kubernetes/ansible/inventory directory with the contents shown below.

::

    # Create an k8s group that contains the masters and nodes groups
    [k8s:children]
    masters
    nodes
    
    # Set variables common for all k8s hosts
    [k8s:vars]
    # SSH user, this user should allow ssh based auth without requiring a password
    ansible_ssh_user=root

    nuage_cluster_network_CIDR=80.80.0.0/16

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
    master.nuageKubernetes.com
    
    # host group for nodes, includes region info
    [nodes]
    node1.nuageKubernetes.com
    node2.nuageKubernetes.com
    master.nuageKubernetes.com
     

Installing the VSP Components for the Single Master
----------------------------------------------------

1. Make sure you are in the nuage-kubernetes/ansible directory. 
2. Run the following command to install the VSP components:

   ::
   
      cd nuage-kubernetes/ansible/scripts
      ./deploy-cluster.sh --tags=nuage
 
  A successful installation displays the following output:
   ::
   
       
       2016-02-11 22:01:49,891 p=16545 u=root |  PLAY RECAP *********************************************************************
       2016-02-11 22:01:49,892 p=16545 u=root |  localhost                : ok=20   changed=0   unreachable=0  failed=0
       2016-02-11 22:01:49,892 p=16545 u=root |  master.nuageKubernetes.com: ok=247  changed=22  unreachable=0  failed=0
       2016-02-11 22:01:49,893 p=16545 u=root |  node1.nuageKubernetes.com : ok=111  changed=21  unreachable=0  failed=0
       2016-02-11 22:01:49,894 p=16545 u=root |  node2.nuageKubernetes.com : ok=111  changed=21  unreachable=0  failed=0
       
3. Verify that the Master-Node connectivity is up and all nodes are running:

   ::
   
      kubectl get nodes       
      

.. Note:: To troubleshoot any installation issues, see the `Troubleshooting <kubernetes-4-troubleshooting.html>`_ section of this document.


