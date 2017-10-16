.. _openshift-4-installation:

.. include:: ../lib/doc-includes/VSDA-icons.inc

===========================
OpenShift Installation
===========================

.. contents::
   :local:
   :depth: 3
   

Supported Platforms
====================

The VSP integration with OpenShift works for both virtual machines (VMs) and bare metal installations of OpenShift.

.. Note:: For information on other supported platforms and distributions, see the *Nuage VSP Release Notes*.


Ansible Installation
==========================

You only need to install Ansible once on a machine and it can manage the master and all the other remote nodes. Ansible manages nodes using the SSH protocol, therefore it is required that SSH is setup so that the master and nodes are accessible from the host running the Ansible scripts.

.. Note:: SSH Protocol does not require a password.

OpenShift DaemonSet for Nuage installation
===========================================

DaemonSet is used for installation of Nuage containerized services as part of the OpenShift Ansible Installation for Nuage. The nuage-openshift-monitor, CNI plugin, and VRS operate as DaemonSet pods on master and slave nodes.

.. Note:: Nuage recommends using Daemonsets for installation of Nuage services.

Pre-Installation Steps in VSD
-----------------------------
1. Login to VSD UI as csproot and create an  "openshift" Enterprise.

2. Under the "openshift" Enterprise, create a user "ose-admin" and add the user to the "Administrators" group.

   .. Note:: The steps to create the user and adding the user to a particular group can be found in the "CSP User Management" section in the "Nuage VSP User Guide."

3. Login to the VSD node using the CLI and execute the following command:

    ::

         /opt/vsd/ejbca/deploy/certMgmt.sh -a generate -u ose-admin -c ose-admin -o openshift -f pem -t client -s root@<ose-master-IP>:/usr/local/

         Where:
         -a <action>         Action [generate|revoke|delete|renew|bscopy]
         -u <username>       End-entity Username
         -c <commonName>     Common Name. Needs to match VSD Username
         -o <organization>   Organization. Needs to match VSD Organization
         -f <format>         Certificate Format [pem|jks|p12]
         -t <type>           Certificate Type [client|server|vsc|vrs]
         -s <scpUrl>         The remote scp url. (eg. root@myhost://home/certs/)


    .. Note:: The above command generates the client certificates for the "ose-admin" user and copies it to the /usr/local/ or any specified directory of the OSE node where Ansible is run. This certificate information is used by the nuage-openshift-monitor to securely communicate with the VSD.

Using Git
-------------

You need to have Git installed on your Ansible machine. Perform the following tasks:

1. Make sure https://github.com is reachable from your Ansible machine.

2. Setup SSH and access the master and the minion nodes, using the **ssh** command.

.. Note:: set-up passwordless **ssh** between Ansible node and cluster nodes.
   
3. Clone the Ansible git repository, by entering the **git clone** command as shown in the following example:

   ::
   
       [root@ansible-mc ~]# git clone https://github.com/openshift/openshift-ansible.git
       Cloning into 'openshift-ansible'...
       remote: Counting objects: 65216, done.
       remote: Total 65216 (delta 0), reused 0 (delta 0), pack-reused 65215
       Receiving objects: 100% (65216/65216), 16.61 MiB | 1.77 MiB/s, done.
       Resolving deltas: 100% (40178/40178), done.
      

      
4. Checkout tag openshift-ansible-3.7.0-0.116.0, by entering the **git checkout** command as shown in the following example:

   ::

      [root@ansible-mc ~]# cd openshift-ansible/
      [root@ansible-mc openshift-ansible]# git checkout tags/openshift-ansible-3.7.0-0.116.0
      Note: checking out 'tags/openshift-ansible-3.7.0-0.116.0'.

      You are in 'detached HEAD' state. You can look around, make experimental
      changes and commit them, and you can discard any commits you make in this
      state without impacting any branches by performing another checkout.

      If you want to create a new branch to retain commits you create, you may
      do so (now or later) by using -b with the checkout command again. Example:

      git checkout -b new_branch_name

      HEAD is now at cc47755... Automatic commit of package [openshift-ansible] release [3.7.0-0.116.0].
      [root@ansible-mc openshift-ansible]#

Setup
----------

1. To prepare the OpenShift cluster for installation, follow the OpenShift Host Preparation guide `here <https://docs.openshift.com/container-platform/3.5/install_config/install/host_preparation.html/>`_.

2. Remove the `Defaults requiretty` from /etc/sudoers on all of the nodes (including master). This is required for accelerated installation mode.

3. Load the following docker images on your master node:

   ::
   
       nuage-master-docker.tar
       nuage-cni-docker.tar
       nuage-vrs-docker.tar

4. Load the following docker images on your slave nodes:

   ::
   
       nuage-cni-docker.tar
       nuage-vrs-docker.tar
       
   
Including the ansible.cfg File
--------------------------------

1. Add a file ansible.cfg in openshift-ansible directory with the following contents:

   ::
   
       [defaults]
       # Add the roles directory to the roles path
       roles_path = roles/
       
       # Set the log_path
       log_path = ~/ansible_logs/ansible.log
       
       [ssh_connection]
       pipelining = True
       

2. Make sure the directory specified for the log_path exists.


Installation for a Single Master
-----------------------------------

1. Create a nodes file for Ansible configuration for a single master in the openshift-ansible directory with the contents shown below.

2. Verify that the image versions are accurate by checking the TAG displayed by 'docker images' output for successful deployment of Nuage daemonsets: 

  .. Note:: The following nodes file is provided as a sample. Please update the values with your actual deployment.
::

    # Create an OSEv3 group that contains the masters and nodes groups
    [OSEv3:children]
    masters
    nodes
    etcd 
    
    # Set variables common for all OSEv3 hosts
    [OSEv3:vars]
    # SSH user, this user should allow ssh based auth without requiring a password
    ansible_ssh_user=root
    openshift_master_portal_net=172.30.0.0/16
    osm_cluster_network_cidr=70.70.0.0/16
    deployment_type=openshift-enterprise
    osm_host_subnet_length=10
    openshift_pkg_version=-3.5.5.5

    # If ansible_ssh_user is not root, ansible_sudo must be set to true
    #ansible_sudo=true 
    
    deployment_type=openshift-enterprise
    
    # Nuage specific parameters
    openshift_use_openshift_sdn=False
    openshift_use_nuage=True
    os_sdn_network_plugin_name=cni
    vsd_api_url=https://<VSD-IP/VSD-Hostname>:7443
    vsp_version=v5_0
    nuage_monitor_image_version=v5.1.1-1
    nuage_vrs_image_version=v5.1.1-1
    nuage_cni_image_version=v5.1.1-1
    enterprise=openshift
    domain=openshift
    vsc_active_ip=10.100.100.101
    vsc_standby_ip=10.100.100.102
    uplink_interface=eth0
    nuage_openshift_monitor_log_dir=/var/log/nuage-openshift-monitor
    nuage_interface_mtu=1500
    
    
    # VSD user in the admin group
    vsd_user=ose-admin
    # Complete local host path to the VSD user certificate file
    vsd_user_cert_file=/usr/local/ose-admin.pem
    # Complete local host path to the VSD user key file
    vsd_user_key_file=/usr/local/ose-admin-Key.pem
   
    
    # Set 'make-iptables-util-chains' flag as 'false' while starting kubelet
    # NOTE: This is a mandatory parameter and Nuage Integration does not work if not set
    openshift_node_kubelet_args={'max-pods': ['110'], 'image-gc-high-threshold': ['90'], 'image-gc-low-threshold': ['80'], 'make-iptables-util-chains': ['false']}
    
    # Required for Nuage Monitor REST server 
    openshift_master_cluster_hostname=master.nuageopenshift.com
    nuage_openshift_monitor_rest_server_port=9443
    
    # host group for masters
    [masters]
    master.nuageopenshift.com
    
    # etcd 
    [etcd]
    etcd.nuageopenshift.com
    
    # host group for nodes, includes region info
    [nodes]
    node1.nuageopenshift.com
    node2.nuageopenshift.com
    master.nuageopenshift.com openshift_node_labels="{'install-monitor': 'true'}"


.. Note:: It is mandatory to specify the openshift_node_labels="{'install-monitor': 'true'}" parameter for the master node for Nuage OpenShift master to be deployed.

Installing the VSP Components for the Single Master
----------------------------------------------------

1. Run the following command to install the VSP components:

   ::
   
       cd openshift-ansible
       ansible-playbook -vvvv -e openshift_disable_check=disk_availability,docker_storage,package_version,memory_availability,package_availability -i nodes playbooks/byo/config.yml
 
  A successful installation displays the following output:
   ::
   
       
       2017-08-11 22:01:49,891 p=16545 u=root |  PLAY RECAP *********************************************************************
       2017-08-11 22:01:49,892 p=16545 u=root |  localhost                : ok=20   changed=0   unreachable=0  failed=0
       2017-08-11 22:01:49,893 p=16545 u=root |  master.nuageopenshift.com: ok=247  changed=22  unreachable=0  failed=0
       2017-08-11 22:01:49,894 p=16545 u=root |  etcd.nuageopenshift.com: ok=247  changed=22  unreachable=0  failed=0
       2017-08-11 22:01:49,895 p=16545 u=root |  node1.nuageopenshift.com : ok=111  changed=21  unreachable=0  failed=0
       2017-08-11 22:01:49,896 p=16545 u=root |  node2.nuageopenshift.com : ok=111  changed=21  unreachable=0  failed=0
       
2. Verify that the Master-Node connectivity is up and all nodes are running:

   ::
   
       oc login -u system:admin
       oc get nodes


Installation for Multiple Masters
----------------------------------

A High Availability (HA) environment can be configured with multiple masters and multiple nodes.

Nuage OpenShift only supports HA configuration method described in this section. This can be combined with any load balancing solution, the default being HAProxy. In the inventory file, there are two master hosts, the nodes, an etcd server and a host that functions as the HAProxy to balance the master API on all master hosts. The HAProxy host is defined in the [lb] section of the inventory file enabling Ansible to automatically install and configure HAProxy as the load balancing solution.

1. Create the nodes file for Ansible configuration for multiple masters in the openshift-ansible directory with the contents shown below.

2. Verify that the image versions are accurate by checking the TAG displayed by 'docker images' output for successful deployment of Nuage daemonsets.

   .. Note:: The following nodes file is provided as a sample. Please update the values with your actual deployment.

    ::
    
        # Create an OSEv3 group that contains the masters and nodes groups
        [OSEv3.1:children]
        masters
        nodes
        etcd
        lb
        
        # Set variables common for all OSEv3 hosts
        [OSEv3:vars]
        # SSH user, this user should allow ssh based auth without requiring a password
        ansible_ssh_user=root
        openshift_master_portal_net=172.30.0.0/16
        osm_cluster_network_cidr=70.70.0.0/16
        deployment_type=openshift-enterprise
        osm_host_subnet_length=10
        openshift_pkg_version=-3.5.5.5
    
        # If ansible_ssh_user is not root, ansible_sudo must be set to true
        #ansible_sudo=true 
        
        deployment_type=openshift-enterprise
        
        # Nuage specific parameters
        openshift_use_openshift_sdn=False
        openshift_use_nuage=True
        os_sdn_network_plugin_name=cni
        vsd_api_url=https://<VSD-IP/VSD-Hostname>:7443
        vsp_version=v5_0
        nuage_monitor_image_version=v5.1.1-1
        nuage_vrs_image_version=v5.1.1-1
        nuage_cni_image_version=v5.1.1-1
        enterprise=openshift
        domain=openshift
        vsc_active_ip=10.100.100.101
        vsc_standby_ip=10.100.100.102
        uplink_interface=eth0
        nuage_openshift_monitor_log_dir=/var/log/nuage-openshift-monitor
        nuage_interface_mtu=1500
        nuage_openshift_monitor_rest_server_port=9443
        
        # VSD user in the admin group
        vsd_user=ose-admin
        # Complete local host path to the VSD user certificate file
        vsd_user_cert_file=/usr/local/ose-admin.pem
        # Complete local host path to the VSD user key file
        vsd_user_key_file=/usr/local/ose-admin-Key.pem
    
        # Set 'make-iptables-util-chains' flag as 'false' while starting kubelet
        # NOTE: This is a mandatory parameter and Nuage Integration does not work if not set
        openshift_node_kubelet_args={'max-pods': ['110'], 'image-gc-high-threshold': ['90'], 'image-gc-low-threshold': ['80'], 'make-iptables-util-chains': ['false']}
    
        # Required for Nuage Monitor REST server and HA
        openshift_master_cluster_method=native
        openshift_master_cluster_hostname=lb.nuageopenshift.com
        openshift_master_cluster_public_hostname=lb.nuageopenshift.com
        
        # host group for masters
        [masters]
        master1.nuageopenshift.com
        master2.nuageopenshift.com
        
        # Specify load balancer host
        [lb]
        lb.nuageopenshift.com
        
        [etcd]
        etcd.nuageopenshift.com
        
        # host group for nodes
        [nodes]
        node1.nuageopenshift.com
        node2.nuageopenshift.com
        master1.nuageopenshift.com openshift_node_labels="{'install-monitor': 'true'}"
        master2.nuageopenshift.com openshift_node_labels="{'install-monitor': 'true'}"
        

.. Note:: It is mandatory to specify the openshift_node_labels="{'install-monitor': 'true'}" parameter for every master node for Nuage OpenShift master to be deployed.


Installing the VSP Components for Multiple Masters
---------------------------------------------------

1. Run the following command to install the VSP components:

   ::
   
       cd openshift-ansible
       ansible-playbook -vvvv -e openshift_disable_check=disk_availability,docker_storage,package_version,memory_availability,package_availability -i nodes playbooks/byo/config.yml

  A successful installation displays the following output:

   ::
   
       
       2017-08-11 22:01:49,891 p=16545 u=root | PLAY RECAP *********************************************************************
       2017-08-11 22:01:49,892 p=16545 u=root | localhost             : ok=20  changed=0  unreachable=0 failed=0
       2017-08-11 22:01:49,892 p=16545 u=root | master1.nuageopenshift.com : ok=247 changed=22 unreachable=0  failed=0
       2017-08-11 22:01:49,893 p=16545 u=root | master2.nuageopenshift.com : ok=248 changed=22 unreachable=0  failed=0
       2017-08-11 22:01:49,894 p=16545 u=root | node1.nuageopenshift.com : ok=111 changed=21 unreachable=0  failed=0
       2017-08-11 22:01:49,895 p=16545 u=root | node2.nuageopenshift.com : ok=111 changed=21 unreachable=0  failed=0 
       
2. Verify that the Master-Node connectivity is up and all nodes are running:

   ::
   
       oc login -u system:admin
       oc get nodes
   
   .. Note:: Both the masters should display all nodes as connected.

3. Ansible configures the loadbalancer to balance the Openshift Master's 9443 port. 

Deploying the Nuage DaemonSet
--------------------------------

The Ansible installer with automatically label the master nodes and deploy the nuage-master-config, nuage-cni-ds and nuage-cni-ds daemonsets. In case of any failures, use the appropriate commands to correct or verify the daemonset files and re-deploy.

The nuage-master-config-daemonset.yaml for openshift-monitor deployment and nuage-node-config-daemonset.yaml for VRS and CNI plugin deployment is copied to /etc/ directory as part of Ansible installation. 

The daemonset files are pre-populated using the values provided in the 'nodes' file during Ansible installation. You may modify the image versions or other relevant parameters in the yaml file. However, it is advised to take a back-up of the yaml files before any modification.

1. Verify the daemonset deployment.

   ::   
       
       [root@master]# oc get ds -n kube-system
        NAME                  DESIRED   CURRENT   READY     NODE-SELECTOR          AGE
        nuage-cni-ds          10        10        10        <none>                 7m
        nuage-master-config   1         1         1         install-monitor=true   1h
        nuage-vrs-ds          10        10        10        <none>                 7m
        
2. Verify that the REST server URL value is correct in the /etc/nuage-node-config-daemonset.yaml file. The 'nuageMonRestServer' should be configured with openshift_master_cluster_hostname value specified in the nodes files during Ansible installation. Modify the value and save the file if this field has incorrect values. Delete and re-deploy the node daemonset as shown in the following steps. 

   ::
   
        # REST server URL
        nuageMonRestServer: https://master.nuageopenshift.com:9443

   .. Note:: If 'nuageMonRestServer' has the value 0.0.0.0:9443, it is incorrect. Please change the value and re-deploy.

2. If you modify the daemonset files, delete and re-deploy the master or node daemonsets respectively using the following commands.

   ::
    
        [root@master]# oc delete -f /etc/nuage-master-config-daemonset.yaml
        configmap "nuage-master-config" deleted
        daemonset "nuage-master-config" deleted
        
        [root@master]# oc delete -f /etc/nuage-node-config-daemonset.yaml 
        configmap "nuage-config" deleted
        daemonset "nuage-cni-ds" deleted
        daemonset "nuage-vrs-ds" deleted
   
        [root@master]# oc create -f /etc/nuage-master-config-daemonset.yaml 
        configmap "nuage-master-config" created
        daemonset "nuage-master-config" created
   
        [root@master]# oc get ds -n kube-system
        NAME                  DESIRED   CURRENT   READY     NODE-SELECTOR          AGE
        nuage-master-config   1         1         1         install-monitor=true   1m
        
        [root@master]# oc create -f /etc/nuage-node-config-daemonset.yaml 
        configmap "nuage-config" created
        daemonset "nuage-cni-ds" created
        daemonset "nuage-vrs-ds" created
        
        [root@master]# oc get ds -n kube-system
        NAME                  DESIRED   CURRENT   READY     NODE-SELECTOR          AGE
        nuage-cni-ds          10        10        10        <none>                 7m
        nuage-master-config   1         1         1         install-monitor=true   1h
        nuage-vrs-ds          10        10        10        <none>                 7m
        

3. The master daemonset deploys the nuage-master-config(nuage-openshift-monitor) pod on the master node and the node daemonset deploys the CNI plugin pod and Nuage VRS pod on every slave node. Following is the output of successfully deployed master and node daemonsets.

   ::
        
        [root@master]# oc get all -n kube-system
        NAME                        READY     STATUS    RESTARTS   AGE
        nuage-cni-ds-04s43          1/1       Running   0          7m
        nuage-cni-ds-81mnp          1/1       Running   0          7m
        nuage-cni-ds-f4q2k          1/1       Running   0          7m
        nuage-cni-ds-hgrjt          1/1       Running   0          7m
        nuage-cni-ds-j0g2k          1/1       Running   0          7m
        nuage-cni-ds-k6df0          1/1       Running   0          7m
        nuage-cni-ds-kclh5          1/1       Running   0          7m
        nuage-cni-ds-l2ftp          1/1       Running   0          7m
        nuage-cni-ds-q68s3          1/1       Running   0          7m
        nuage-cni-ds-zkdv4          1/1       Running   0          7m
        nuage-master-config-0d95v   1/1       Running   0          1h
        nuage-vrs-ds-0v9sq          1/1       Running   0          7m
        nuage-vrs-ds-c0kt5          1/1       Running   0          7m
        nuage-vrs-ds-d4h7m          1/1       Running   0          7m
        nuage-vrs-ds-kqmhf          1/1       Running   0          7m
        nuage-vrs-ds-qcq65          1/1       Running   0          7m
        nuage-vrs-ds-qkxv3          1/1       Running   0          7m
        nuage-vrs-ds-qpp21          1/1       Running   0          7m
        nuage-vrs-ds-rg9w1          1/1       Running   0          7m
        nuage-vrs-ds-w05bw          1/1       Running   0          7m
        nuage-vrs-ds-w5v9r          1/1       Running   0          7m
   
    

