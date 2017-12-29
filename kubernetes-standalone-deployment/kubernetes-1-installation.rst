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

With Nuage release 5.2.1 and above, installation of Nuage components will be done using daemonset files only and ansible playbooks are no longer required to be run by the user.

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

         /opt/vsd/ejbca/deploy/certMgmt.sh -a generate -u k8s-admin -c k8s-admin -o kubernetes -f pem -t client -s root@<k8s-master-IP>:/etc/kubernetes/pki/

         
         Where:
         -a <action>         Action [generate|revoke|delete|renew|bscopy]
         -u <username>       End-entity Username
         -c <commonName>     Common Name. Needs to match VSD Username
         -o <organization>   Organization. Needs to match VSD Organization
         -f <format>         Certificate Format [pem|jks|p12]
         -t <type>           Certificate Type [client|server|vsc|vrs]
         -s <scpUrl>         The remote scp url. (eg. root@myhost://home/certs/)


	.. Note:: The above command generates the client certificates for the "k8s-admin" user and copies it to the /etc/kubernetes/pki or any specified directory of the k8s master node where daemonsets are going to be created. This certificate information is used by the nuagekubemon (nuage k8S monitor) to securely communicate with the VSD.

4. To complete the steps provided in the Kubeadm installer guide, go `here <https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/>`_. 

  .. Note:: The default directory for kubernetes certs is in /etc/kubernetes/pki on the master node. If you specify any other directory while doing kubeadm init then that directory needs to be specified above in the VSD cert generation step

  .. Note:: Kubectl needs the kube config to be copied to a specific location after the master is initialized using kubeadm init. Ansible scripts used to install Nuage components also rely on kubectl being available to the ansible user. In order to achieve that, execute the following commands on the master node after kubeadm init:
    ::
          sudo mkdir -p ~/.kube
          sudo cp /etc/kubernetes/admin.conf ~/.kube/config
          sudo chown $(id -u):$(id -g) $HOME/.kube/config

5. Follow the steps 1, 2 & 4 of the document provided in the above link. For the pod network, install Nuage using the Ansible installer mentioned below. 

  .. Note:: By default, Kubeadm uses 10.96.0.0/12 as the service CIDR. Make sure this service CIDR does not overlap with your existing underlay network CIDR. If it does then, please run step 2 from the above guide as follows so as to change the service CIDR:
         kubeadm init --service-cidr=192.168.0.0/16 --kubernetes-version <k8s-version>

  .. Note:: With new version of Kubernetes, pod network is required by default and nodes won't move to ready state until pod network is installed. So, after kubeadm join is done on all nodes, your kubectl get nodes can be as follows:

  ::

          kubectl get nodes
          NAME                            STATUS     AGE       VERSION

          ovs-1.test.nuagenetworks.com    NotReady   3h        v1.9.0
          ovs-10.test.nuagenetworks.com   NotReady   3h        v1.9.0
          ovs-2.test.nuagenetworks.com    NotReady   3h        v1.9.0
          ovs-3.test.nuagenetworks.com    NotReady   3h        v1.9.0
          ovs-4.test.nuagenetworks.com    NotReady   3h        v1.9.0
          ovs-5.test.nuagenetworks.com    NotReady   3h        v1.9.0

6. Update the cluster-dns in the 10-kubeadm.conf file on all nodes and master as follows:
  
    On the Node:
    
    ::
        
    	cat /etc/systemd/system/kubelet.service.d/10-kubeadm.conf 
    	[Service]
        Environment="KUBELET_KUBECONFIG_ARGS=--bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf -- kubeconfig=/etc/kubernetes/kubelet.conf"
        Environment="KUBELET_SYSTEM_PODS_ARGS=--pod-manifest-path=/etc/kubernetes/manifests --allow-privileged=true"
        Environment="KUBELET_NETWORK_ARGS=--network-plugin=cni --cni-conf-dir=/etc/cni/net.d --cni-bin-dir=/opt/cni/bin"
        Environment="KUBELET_DNS_ARGS=--cluster-dns=192.168.0.10 --cluster-domain=cluster.local"
        Environment="KUBELET_AUTHZ_ARGS=--authorization-mode=Webhook --client-ca-file=/etc/kubernetes/pki/ca.crt"
        Environment="KUBELET_CADVISOR_ARGS=--cadvisor-port=0"
        Environment="KUBELET_CGROUP_ARGS=--cgroup-driver=cgroupfs"
        Environment="KUBELET_CERTIFICATE_ARGS=--rotate-certificates=true --cert-dir=/var/lib/kubelet/pki"
        Environment="KUBELET_EXTRA_ARGS=--fail-swap-on=false"
        ExecStart=

        ExecStart=/usr/bin/kubelet $KUBELET_KUBECONFIG_ARGS $KUBELET_SYSTEM_PODS_ARGS $KUBELET_NETWORK_ARGS $KUBELET_DNS_ARGS $KUBELET_AUTHZ_ARGS $KUBELET_CADVISOR_ARGS $KUBELET_CGROUP_ARGS $KUBELET_CERTIFICATE_ARGS $KUBELET_EXTRA_ARGS

    
    On the Master:
    
    ::
    
    	cat /etc/systemd/system/kubelet.service.d/10-kubeadm.conf
	[Service]
        Environment="KUBELET_KUBECONFIG_ARGS=--bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf"
        Environment="KUBELET_SYSTEM_PODS_ARGS=--pod-manifest-path=/etc/kubernetes/manifests --allow-privileged=true"
        Environment="KUBELET_NETWORK_ARGS=--network-plugin=cni --cni-conf-dir=/etc/cni/net.d --cni-bin-dir=/opt/cni/bin"
        Environment="KUBELET_DNS_ARGS=--cluster-dns=192.168.0.10 --cluster-domain=cluster.local"
        Environment="KUBELET_AUTHZ_ARGS=--authorization-mode=Webhook --client-ca-file=/etc/kubernetes/pki/ca.crt"
        Environment="KUBELET_CADVISOR_ARGS=--cadvisor-port=0"
        Environment="KUBELET_CGROUP_ARGS=--cgroup-driver=cgroupfs"
        Environment="KUBELET_CERTIFICATE_ARGS=--rotate-certificates=true --cert-dir=/var/lib/kubelet/pki"
        Environment="KUBELET_EXTRA_ARGS=--fail-swap-on=false"
        ExecStart=

        ExecStart=/usr/bin/kubelet $KUBELET_KUBECONFIG_ARGS $KUBELET_SYSTEM_PODS_ARGS $KUBELET_NETWORK_ARGS $KUBELET_DNS_ARGS $KUBELET_AUTHZ_ARGS $KUBELET_CADVISOR_ARGS $KUBELET_CGROUP_ARGS $KUBELET_CERTIFICATE_ARGS $KUBELET_EXTRA_ARGS

  .. Note:: With new version of Kubernetes, KUBELET_CGROUP_ARGS are required to be added to ExecStart or else it will cause the kubelet to fail. The cgroup-driver can be set to cgroupfs or systemd depending on the driver used for docker installation. By default, its set to systemd. Also, if swap is enabled on the node/master then set --fail-swap-on=false in KUBELET_EXTRA_ARGS and if swap is not required then turn it off using command `swapoff -a` or else it will result in kubelet failing on the node/master.
    
   This service CIDR also gets updated in another file as explained in the "Installation for a Single Master" section.

Install Git Repository
-----------------------

You need to have Git installed on your Ansible host machine. Perform the following tasks:

1. Access Git 
2. Setup SSH and access the master and the minion nodes, using the **ssh** command.
3. Clone the Ansible git repository, by entering the **git clone** command as shown in the example below and checkout the tag **nuage-kubernetes-<version>** corresponding to the VSP version. 

   ::
   
        git clone https://github.com/nuagenetworks/nuage-kubernetes.git
        git checkout tags/<Nuage-release>
        cd nuage-kubernetes/daemonset-templates
	
  .. Note:: Post Nuage 5.1.1, rpm based install is not supported using ansible. Daemonsets is the recommended mode of installing Nuage components 


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

6. Update the following parameters in ConfigMap section of **nuage-kubernetes/daemonset-templates/nuage-master-config-daemonset.yaml** file as per your environment configuration:

::

      # This will generate the required Nuage monitor configuration
      # on master nodes
      monitor_yaml_config: |
          kubeConfig: /usr/share/vsp-k8s/nuage.kubeconfig

          masterConfig: /usr/share/nuagekubemon/net-config.yaml
          # URL of the VSD Architect
          vsdApiUrl: https://xmpp.example.com:7443                    <--- Set this to the VSD IP or hostname for the cluster
          # API version to query against
          vspVersion: v5_0
          # Name of the enterprise in which pods will reside
          enterpriseName: kubernetes
          # Name of the domain in which pods will reside
          domainName: kubernetes
          # VSD generated user certificate file location on master node
          userCertificateFile: /etc/kubernetes/pki/k8s-admin.pem      <--- Set this to the cert-dir which holds the client certificates
          # VSD generated user key file location on master node            generated above in the pre installation steps
          userKeyFile: /etc/kubernetes/pki/k8s-admin-Key.pem
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
              clientCA: ""
              serverCertificate: ""
              serverKey: ""

Make sure to set the etcd config correctly if there is an external etcd cluster. If the etcd cluster is not using TLS certificates, do not set the ca, certFile & keyFile parameters. Also, if etcd is running locally on the master, use the localhost IP as shown below. If the etcd cluster is setup using FQDN, set the URL to the FQDN hostname. Also, make sure to check the protocol for the etcd cluster and set http or https for the URL accordingly.

::

          # etcd config required for HA
          etcdClientConfig:
              ca: ""
              certFile: ""
              keyFile: ""
              urls:
                 - http://127.0.0.1:2379

Set the parameter to 1 in order to allow nuagekubemon to automagically create a new subnet when the existing subnet gets depleted. Threshold for new subnet creation is set to 70% namespace/zone allocation. It will also delete additional subnets if the namespace usage falls below 25%

::

          # auto scale subnets feature
          # 0 => disabled(default)
          # 1 => enabled
          autoScaleSubnets: 1

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
          image: nuage/master:<nuage-release>


The kubernetes-cert-dir mount path needs to be set if the user has specified the --cert-dir option during kubeadm init phase. The default cert-dir is set to /etc/kubernetes/pki

::

	  volumeMounts:
            - mountPath: /var/log
              name: cni-log-dir
            - mountPath: /usr/share
              name: usr-share-dir
            - mountPath: /etc/kubernetes/pki/
              name: kubernetes-cert-dir
      volumes:
        - name: cni-log-dir
          hostPath:
            path: /var/log
        - name: usr-share-dir
          hostPath:
            path: /usr/share
        - name: kubernetes-cert-dir
          hostPath:
            path: /etc/kubernetes/pki/

7. Update the following parameters in **nuage-kubernetes/daemonset-templates/nuage-node-config-daemonset.yaml** file as per your environment configuration:

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
            nuageMonClientCert: /var/lib/kubelet/pki/kubelet-client.crt  <--- Nuage will be re-using kubernetes generated certs for
            # Key to the certificate in restClientCert                        communication between the CNI plugin and the monitor
            nuageMonClientKey: /var/lib/kubelet/pki/kubelet-client.key        daemon running the on the masters. Set this parameter to
	    # CA certificate for verifying the master's rest server           the location where the kubelet-client.crt and key are
            nuageMonServerCA: /etc/kubernetes/pki/ca.crt                      stored on the nodes. 
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


Make sure the **image** parameter is correctly set to the Nuage CNI docker image version pre-loaded on the slave nodes:

::
      containers:
        # This container installs Nuage CNI binaries
        # and CNI network config file on each node.
        - name: install-nuage-cni
          image: nuage/cni:<nuage-version>

Update the following environment variables in the DaemonSet section for **nuage-cni-ds** with the value set in clusterNetworkCIDR in the nuage-master-config-daemonset.yaml above    

::

            # Nuage cluster network CIDR for iptables configuration
              - name: NUAGE_CLUSTER_NW_CIDR
                value: "70.70.0.0/16"


Steps to be followed to generate the Nuage Token:

i. On the master node, create a nuage-serviceaccount.yaml file as shown below:

::

     cat nuage-serviceaccount.yaml 
     
     apiVersion: v1
     kind: ServiceAccount
     metadata:
         name: nuage

ii. Create the nuage-serviceaccount on the master using the command `kubectl create -f nuage-serviceaccount.yaml`

iii. Execute the following commands on the master:

::

     kubectl create clusterrolebinding add-on-cluster-admin --clusterrole=cluster-admin --serviceaccount=default:nuage --namespace=default

     kubectl describe secret $(kubectl get serviceaccounts/nuage -o yaml | grep -o "nuage-token.*") | grep "token:" | awk '{print $2}'

iv. The output of the above command is the Nuage Token secret which needs to be added below

Update the master IP and Nuage Token in the daemonset file and also the kubernetes-cert-dir if it is set to a specific value other than /etc/kubernetes/pki which is used by default

::

            # Kubernetes Master api-server URL
            - name: MASTER_API_SERVER_URL
              value: "https://<master-ip>:6443"  <-- Set the master-ip to the Kubernetes Master IP or hostname
            # nuage user service account token string
            - name: NUAGE_TOKEN
              value: "Add Kubernetes generated nuage service account token here" <--- The Nuage Token secret needs to be generated as 
          volumeMounts:                                                               per the steps mentioned above and then the token
            - mountPath: /host/opt                                                    secret needs to be added here
              name: cni-bin-dir
            - mountPath: /host/etc
              name: cni-yaml-dir
            - mountPath: /var/run
              name: var-run-dir
            - mountPath: /var/log
              name: cni-log-dir
            - mountPath: /usr/share
              name: usr-share-dir
            - mountPath: /etc/kubernetes/pki/
              name: kubernetes-ca-dir
            - mountPath: /var/lib/kubelet/pki/
              name: kubernetes-cert-dir
      volumes:
        - name: cni-bin-dir
          hostPath:
            path: /opt
        - name: cni-yaml-dir
          hostPath:
            path: /etc
        - name: var-run-dir
          hostPath:
            path: /var/run
        - name: cni-log-dir
          hostPath:
            path: /var/log
        - name: usr-share-dir
          hostPath:
            path: /usr/share
        - name: kubernetes-ca-dir
          hostPath:
            path: /etc/kubernetes/pki/
        - name: kubernetes-cert-dir
          hostPath:
            path: /var/lib/kubelet/pki/

Make sure the **image** parameter is correctly set to the Nuage VRS docker image version pre-loaded on the slave nodes:

::

      containers:
        # This container installs Nuage VRS running as a
        # container on each worker node
        - name: install-nuage-vrs
          image: nuage/vrs:<nuage-version>

Update the following environment variables in DaemonSet section for **nuage-vrs-ds** with Active and Standby Nuage VSC IP addresses for containerized Nuage VRS and NUAGE_K8S_SERVICE_IPV4_SUBNET with the value for serviceNetworkCIDR set in nuage-master-config-daemonset.yaml above

::

      env:
        # Configure parameters for VRS openvswitch file
        - name: NUAGE_ACTIVE_CONTROLLER
          value: "10.10.10.10"
        - name: NUAGE_STANDBY_CONTROLLER
          value: "20.20.20.20"
        - name: NUAGE_PLATFORM
          value: '"kvm, k8s"'
        - name: NUAGE_K8S_SERVICE_IPV4_SUBNET
          value: '192.168.0.0\/16'

Update the **image** parameter in **nuage-kubernetes/daemonset-templates/nuage-infra-pod-config-daemonset.yaml** file and make sure that it is correctly set to the Nuage infra pod image version pre-loaded on the slave nodes:

::

      containers:
        # This container spawns a Nuage Infra pod
        # on each worker node
        - name: install-nuage-infra
          image: nuage/infra:<nuage-version>

       

Installing the VSP components for a Single Master 
--------------------------------------------------

1. After updating all the 3 daemonset files with the correct values, the daemonsets can be created using the following commands:

::

    kubectl apply -f nuage-kubernetes/daemonset-templates/nuage-node-config-daemonset.yaml
    kubectl apply -f nuage-kubernetes/daemonset-templates/nuage-master-config-daemonset.yaml
    kubectl apply -f nuage-kubernetes/daemonset-templates/nuage-infra-pod-config-daemonset.yaml
    
2. Verify that all the Nuage monitor, CNI and VRS pods are up and running:

   ::
   
      `kubectl get pods -n kube-system`      
     

