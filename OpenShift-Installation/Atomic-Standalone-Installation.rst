
=======================================
OpenShift Installation on Atomic Hosts
=======================================

.. contents::
   :local:
   :depth: 3
   

Supported Platforms
====================

The VSP integration with OpenShift is supported on RHEL Atomic hosts (VERSION 7.3).

The installation procedure in this section is for VSP integration with OpenShift when the master nodes are RHEL Servers and slave nodes are RHEL Atomic Hosts.

.. Note:: For information on other supported platforms and distributions, see the *Nuage VSP Release Notes*.


Ansible Installation
==========================

You only need to install Ansible once on a machine and it can manage the master and all the other remote nodes. Ansible manages nodes using the SSH protocol, therefore it is required that SSH is setup so that the master and nodes are accessible from the host running the Ansible scripts.

.. Note:: SSH Protocol does not require a password.

OpenShift DaemonSet for Nuage Installation
===========================================

DaemonSet is used for installation of Nuage containerized services as part of the OpenShift Ansible Installation for Nuage. The nuage-openshift-monitor, CNI plugin, and VRS operate as DaemonSet pods on master and slave nodes.

.. Note:: Nuage recommends using DaemonSets for installation of Nuage services.

Pre-Installation Steps in VSD
-----------------------------
1. Login to VSD UI as csproot and create an  "openshift" Enterprise.

2. Under the "openshift" Enterprise, create a user "ose-admin" and add the user to the "Administrators" group.

   .. Note:: The steps to create the user and adding the user to a particular group can be found in the "CSP User Management" section in the "Nuage VSP User Guide."

3. Login to the VSD node using the CLI and execute the following command. The certs must be placed in the /usr/local/ or any specified directory on the host where ansible is run:

    ::

         /opt/vsd/ejbca/deploy/certMgmt.sh -a generate -u ose-admin -c ose-admin -o openshift -f pem -t client -s root@<ose-ansible-IP>:/usr/local/

         Where:
         -a <action>         Action [generate|revoke|delete|renew|bscopy]
         -u <username>       End-entity Username
         -c <commonName>     Common Name. Needs to match VSD Username
         -o <organization>   Organization. Needs to match VSD Organization
         -f <format>         Certificate Format [pem|jks|p12]
         -t <type>           Certificate Type [client|server|vsc|vrs]
         -s <scpUrl>         The remote scp url. (eg. root@myhost://home/certs/)


   .. Note:: The above command generates the client certificates for the "ose-admin" user and copies it to the /usr/local/ or any specified directory of the OSE node where Ansible is run. This certificate information is used by the nuage-openshift-monitor to securely communicate with the VSD.

Git Ansible-Copy-Install
------------------------

You need to have Git installed on your Ansible machine. Perform the following tasks:

1. Make sure https://github.com is reachable from your Ansible machine.

2. Setup SSH and access the master and the minion nodes, using the ssh command.

   .. Note:: set-up passwordless ssh between Ansible node and cluster nodes.

3. Copy the nuage-ose-atomic-install-5-1-2.tar.gz file shipped with Nuage 5.1.2 Release to a host machine where Ansible is run.

4. Unzip and Untar the above image

  ::
      
       [root@ansible-host ~]# tar -xvf nuage-ose-atomic-install.tar 
       etcd_certificates.yml
       main.yml
       nuage-master-config-daemonset.j2
       nuage-node-config-daemonset.j2
       nuage-openshift-ansible.diff
       patch-nuage-openshift-ansible.sh

       [root@ansible-host ~]# ls
       etcd_certificates.yml  nuage-master-config-daemonset.j2  nuage-openshift-ansible.diff  patch-nuage-openshift-ansible.sh
       main.yml               nuage-node-config-daemonset.j2    nuage-ose-atomic-install.tar

   
3. Run the patch-nuage-openshift-ansible.sh script to clone the ansible repo and set up Nuage changes.

   ::
   
       [root@ansible-host ~]# ./patch-nuage-openshift-ansible.sh 
       Cloning into 'openshift-ansible'...
       remote: Counting objects: 71754, done.
       remote: Compressing objects: 100% (11/11), done.
       remote: Total 71754 (delta 0), reused 6 (delta 0), pack-reused 71742
       Receiving objects: 100% (71754/71754), 18.28 MiB | 2.48 MiB/s, done.
       Resolving deltas: 100% (44453/44453), done.
       Checking connectivity... done.
       Note: checking out 'tags/openshift-ansible-3.6.128-1'.

       You are in 'detached HEAD' state. You can look around, make experimental
       changes and commit them, and you can discard any commits you make in this
       state without impacting any branches by performing another checkout.

       If you want to create a new branch to retain commits you create, you may
       do so (now or later) by using -b with the checkout command again. Example:

       git checkout -b <new-branch-name>

       HEAD is now at 2d7e10b... Automatic commit of package [openshift-ansible] release [3.6.128-1].
       Successfully patched Nuage ansible changes into openshift-ansible
       You may now use the openshift-ansible folder for your ansible installation
      

Setup
----------

1. To prepare the OpenShift cluster for installation, follow the OpenShift Host Preparation guide `here <https://docs.openshift.com/container-platform/3.6/install_config/install/host_preparation.html/>`_.

   .. Note:: Skip the yum update part in the OpenShift Host Preparation guide.

2. Load the following docker images on your master node:

   ::
   
       nuage-master-docker.tar
       nuage-cni-docker.tar
       nuage-vrs-docker.tar

3. Load the following docker images on your slave nodes:

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

  .. Note:: The following nodes file is provided as a sample. Please update the values with your actual deployment. The below nodes file deploys OpenShift version 3.6
  
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
    openshift_pkg_version=-3.6.173.0.5
    slave_base_host_type=is_atomic
    openshift_disable_check=disk_availability,memory_availability,package_version,docker_storage,docker_image_availability
    
    # If ansible_ssh_user is not root, ansible_sudo must be set to true
    #ansible_sudo=true 
    
    openshift_deployment_type=openshift-enterprise
    
    # Nuage specific parameters
    openshift_use_openshift_sdn=False
    openshift_use_nuage=True
    openshift.common._use_nuage=True
    os_sdn_network_plugin_name=cni
    vsd_api_url=https://<VSD-IP/VSD-Hostname>:7443
    vsp_version=v5_0
    nuage_monitor_image_version=v5.1.2-1
    nuage_vrs_image_version=v5.1.2-1
    nuage_cni_image_version=v5.1.2-1
    enterprise=openshift
    domain=openshift
    vsc_active_ip=10.100.100.101
    vsc_standby_ip=10.100.100.102
    uplink_interface=eth0
    nuage_openshift_monitor_log_dir=/var/log/nuage-openshift-monitor
    nuage_interface_mtu=1500
    # auto scale subnets feature
    # 0 => disabled(default)
    # 1 => enabled
    auto_scale_subnets=0
        
    # VSD user in the admin group
    vsd_user=ose-admin
    # Complete local host path to the VSD user certificate file
    vsd_user_cert_file=/usr/local/ose-admin.pem
    # Complete local host path to the VSD user key file
    vsd_user_key_file=/usr/local/ose-admin-Key.pem
   
    # Set 'make-iptables-util-chains' flag as 'false' while starting kubelet
    # NOTE: This is a mandatory parameter and Nuage Integration does not work if not set
    openshift_node_kubelet_args={'max-pods': ['110'], 'image-gc-high-threshold': ['90'], 'image-gc-low-threshold': ['80'], 'make-iptables-util-chains': ['false']}
    openshift_master_cluster_method=native
    
    # Required for Nuage Monitor REST server 
    openshift_master_cluster_hostname=master.nuageopenshift.com
    openshift_master_cluster_public_hostname=master.nuageopenshift.com
    nuage_openshift_monitor_rest_server_port=9443
    
    # host group for masters
    [masters]
    master.nuageopenshift.com
    
    # etcd 
    [etcd]
    etcd.nuageopenshift.com
    
    # host group for nodes, includes region info
    [nodes]
    node1.nuageopenshift.com openshift_schedulable=True openshift_node_labels="{'region': 'infra'}"
    node2.nuageopenshift.com
    master.nuageopenshift.com openshift_node_labels="{'install-monitor': 'true'}"


.. Note:: It is mandatory to specify the openshift_node_labels="{'install-monitor': 'true'}" parameter for the master node for Nuage OpenShift master to be deployed.

Installing the VSP Components for the Single Master
----------------------------------------------------

1. Run the following command to install the VSP components:

   ::
   
       cd openshift-ansible
       ansible-playbook -vvvv -i nodes playbooks/byo/config.yml
 
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


Deploying the Nuage DaemonSet
--------------------------------

The Ansible installer will automatically label the master nodes and deploy the nuage-master-config, nuage-cni-ds and nuage-cni-ds daemonsets. In case of any failures, use the appropriate commands to correct or verify the daemonset files and re-deploy.

The nuage-master-config-daemonset.yaml for openshift-monitor deployment and nuage-node-config-daemonset.yaml for VRS and CNI plugin deployment is copied to /etc/ directory as part of Ansible installation. 

The daemonset files are pre-populated using the values provided in the 'nodes' file during Ansible installation. You may modify the image versions or other relevant parameters in the yaml file. However, it is advised to take a back-up of the yaml files before any modification.

1. Verify the daemonset deployment.

   ::   
       
       [root@master]# oc get ds -n kube-system
        NAME                  DESIRED   CURRENT   READY     NODE-SELECTOR          AGE
        nuage-cni-ds           3        3         3         <none>                 7m
        nuage-master-config    1        1         1         install-monitor=true   7m
        nuage-vrs-ds           3        3         3         <none>                 7m
        
2. Verify that the REST server URL value is correct in the /etc/nuage-node-config-daemonset.yaml file. The 'nuageMonRestServer' should be configured with openshift_master_cluster_hostname value specified in the nodes files during Ansible installation. Modify the value and save the file if this field has incorrect values. Delete and re-deploy the node daemonset as shown in the following steps. 

   ::
   
        # REST server URL
        nuageMonRestServer: https://master.nuageopenshift.com:9443

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
        nuage-cni-ds             3        3        3        <none>                 1m
        nuage-master-config      1        1        1        install-monitor=true   1m
        nuage-vrs-ds             3        3        3        <none>                 1m
        

3. The master daemonset deploys the nuage-master-config(nuage-openshift-monitor) pod on the master node and the node daemonset deploys the CNI plugin pod and Nuage VRS pod on every slave node. Following is the output of successfully deployed master and node daemonsets.

   ::
        
        [root@master]# oc get all -n kube-system
        NAME                        READY     STATUS    RESTARTS   AGE
        nuage-cni-ds-04s43          1/1       Running   0          7m
        nuage-cni-ds-81mnp          1/1       Running   0          7m
        nuage-cni-ds-f4q2k          1/1       Running   0          7m
        nuage-master-config-0d95v   1/1       Running   0          7m
        nuage-vrs-ds-0v9sq          1/1       Running   0          7m
        nuage-vrs-ds-c0kt5          1/1       Running   0          7m
        nuage-vrs-ds-d4h7m          1/1       Running   0          7m
        
    
   
