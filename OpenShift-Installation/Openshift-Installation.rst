
=============================================
OpenShift Installation - Standalone and HA
=============================================

.. contents::
   :local:
   :depth: 3
   

Supported Platforms
====================

The VSP integration with OpenShift is supported on RHEL Server hosts (Version 7.3).

.. Note:: For information on other supported platforms and distributions, see the *Nuage VSP Release Notes*.


Ansible Installation
==========================

You only need to install Ansible once on a machine and it can manage the master and all the other remote nodes. Ansible manages nodes using the SSH protocol, therefore it is required that SSH is setup so that the master and nodes are accessible from the host running the Ansible scripts.

.. Note:: SSH Protocol does not require a password.

The Ansible Version to be used is 2.3. To check or install ansible version 2.3, follow the below steps:

    ::

         [root@ansible-host ~]# ansible --version
         ansible 2.4.2.0
         
         [root@ansible-host ~]# sudo pip install ansible==2.3
         
         
OpenShift DaemonSet for Nuage installation
===========================================

DaemonSet is used for installation of Nuage containerized services as part of the OpenShift Ansible Installation for Nuage. The nuage-openshift-monitor, CNI plugin, nuage-infra pod and VRS operate as DaemonSet pods on master and slave nodes.

.. Note:: Nuage recommends using Daemonsets for installation of Nuage services.

Pre-Installation Steps in VSD
-----------------------------
1. Login to VSD UI as csproot and create an  "openshift" Enterprise.

2. Under the "openshift" Enterprise, create a user "ose-admin" and add the user to the "Administrators" group. It is mandatory to add the user to the "Administrators" group on the VSD.

   .. Note:: The steps to create the user and adding the user to a particular group can be found in the "CSP User Management" section in the "Nuage VSP User Guide."

3. Login to the VSD node using the CLI and execute the following command. The certificates must be placed in the /usr/local or any specified directory on the host where ansible is run.

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

Using Git
-------------

You need to have Git installed on your Ansible machine. Perform the following tasks:

1. Make sure https://github.com is reachable from your Ansible machine.

2. Setup SSH and access the master and the minion nodes, using the **ssh** command.

.. Note:: set-up passwordless **ssh** between Ansible node and cluster nodes.
   
3. Copy the nuage-ose-rhel-install-<version>.tar file shipped with the specific Nuage Release to a host machine where Ansible is run.

4. Untar the above image

   ::
   
       [root@ovs-10 ~]# tar -xvf nuage-ose-rhel-install.tar 
       etcd_certificates.yml
       nuage-infra-pod-config-daemonset.j2
       nuage-openshift-ansible.diff
       patch-nuage-openshift-ansible.sh

      
5. Run the patch-nuage-openshift-ansible.sh script to clone the ansible repo and set up Nuage changes.

   .. Note:: The tag version may change based on the Nuage release you use.
   
   ::

      [root@ansible-mc ~]# ./patch-nuage-openshift-ansible.sh 
      Cloning into 'openshift-ansible'...
      remote: Counting objects: 74439, done.
      remote: Compressing objects: 100% (7/7), done.
      remote: Total 74439 (delta 4), reused 2 (delta 2), pack-reused 74430
      Receiving objects: 100% (74439/74439), 18.85 MiB | 3.60 MiB/s, done.
      Resolving deltas: 100% (46284/46284), done.
      Checking connectivity... done.
      Note: checking out 'tags/openshift-ansible-3.7.0-0.116.0'.

      You are in 'detached HEAD' state. You can look around, make experimental
      changes and commit them, and you can discard any commits you make in this
      state without impacting any branches by performing another checkout.

      If you want to create a new branch to retain commits you create, you may
      do so (now or later) by using -b with the checkout command again. Example:

      git checkout -b <new-branch-name>

      HEAD is now at cc47755... Automatic commit of package [openshift-ansible] release [3.7.0-0.116.0].
      Successfully patched Nuage ansible changes into openshift-ansible
      You may now use the openshift-ansible folder for your ansible installation


Setup
----------

1. To prepare the OpenShift cluster for installation, follow the OpenShift Host Preparation guide.
   
   For Nuage releases prior to 5.2.1, go `here <https://docs.openshift.com/container-platform/3.5/install_config/install/host_preparation.html>`_.
   
   For Nuage releases 5.2.1 and later, go `here <https://docs.openshift.com/container-platform/3.6/install_config/install/host_preparation.html>`_. 

   .. Note:: Skip the yum update part in the OpenShift Host Preparation guide.


1. Load the following docker images on your master node:

   ::
   
       nuage-master-docker-<version>.tar
       nuage-cni-docker-<version>.tar
       nuage-vrs-docker-<version>.tar
       nuage-infra-docker-<version>.tar

3. Load the following docker images on your slave nodes:

   ::
   
       nuage-cni-docker-<version>.tar
       nuage-vrs-docker-<version>.tar
       nuage-infra-docker-<version>.tar

4. By loading the images, we mean loading the images to docker using 'docker load -i' command. Example shown below

   ::
   
      [root@node-1 .ssh]# docker load -i nuage-vrs-docker-<version>.tar 
      b431d6b0d399: Loading layer [==================================================>] 7.591 MB/7.591 MB
      3936811d0a81: Loading layer [==================================================>]   173 MB/173 MB
      Loaded image: nuage/vrs:<version>

      [root@node-1 ~]# docker load -i nuage-infra-docker-<version>.tar 
      6a749002dd6a: Loading layer [==================================================>] 1.338 MB/1.338 MB
      6b59b94504a9: Loading layer [==================================================>] 2.048 kB/2.048 kB
      Loaded image: nuage/infra:<version>

      [root@node-1 ~]# docker load -i nuage-cni-docker-<version>.tar
      99b28d9413e4: Loading layer [==================================================>] 200.2 MB/200.2 MB
      1541333c4fbd: Loading layer [==================================================>]  63.9 MB/63.9 MB
      523358a7deb2: Loading layer [==================================================>]  63.9 MB/63.9 MB
      62e0df2908be: Loading layer [==================================================>] 3.174 MB/3.174 MB
      a658b822d29a: Loading layer [==================================================>] 5.632 kB/5.632 kB
      b2914c7a133a: Loading layer [==================================================>] 2.048 kB/2.048 kB
      bb72aaeb25b7: Loading layer [==================================================>] 2.048 kB/2.048 kB
      4defe2b005cb: Loading layer [==================================================>] 75.86 MB/75.86 MB
      Loaded image: nuage/cni:<version>

      [root@ovs-1 ~]# docker images
      REPOSITORY                                                 TAG                 IMAGE ID            CREATED             SIZE
      nuage/vrs                                                  <version>           0f83ba129dc2        14 hours ago        505.8 MB
      nuage/infra                                                <version>           53580dde0343        13 days ago         1.13 MB
      nuage/cni                                                  <version>           01be44d6d037        5 weeks ago         399.1 MB
      registry.access.redhat.com/openshift3/ose-haproxy-router   v3.5.5.5            0da1e83a3dc6        6 months ago        775.2 MB
      registry.access.redhat.com/openshift3/node                 v3.5.5.5            2a27de91c442        6 months ago        981.4 MB
      registry.access.redhat.com/openshift3/ose-deployer         v3.5.5.5            fe58ee8f8ced        6 months ago        755.7 MB
      registry.access.redhat.com/openshift3/ose-pod              v3.5.5.5            b717234be287        6 months ago        205.6 MB
   
   
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

  .. Note:: The following nodes file is provided as a sample. Please update the values with your actual deployment. The below nodes file deploys OpenShift version 3.5. To deploy OpenShift version 3.6, use 'openshift_pkg_version=-3.6.173.0.5'
  
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
    openshift_disable_check=disk_availability,memory_availability,docker_storage,docker_image_availability,package_version,package_availability
    
    # Nuage specific parameters
    openshift_use_openshift_sdn=False
    openshift_use_nuage=True
    openshift.common._use_nuage=True
    os_sdn_network_plugin_name=cni
    vsd_api_url=https://<VSD-IP/VSD-Hostname>:7443
    vsp_version=v5_0
    
    # The below versions should match the TAG version in the output of 'docker images' on the nodes. See point 2 above
    # Example: nuage_monitor_image_version=5.1.2-70
    nuage_monitor_image_version=<version>
    nuage_vrs_image_version=<version>
    nuage_cni_image_version=<version>
    nuage_infra_container_image_version=<version>
    
    enterprise=openshift
    domain=openshift
    vsc_active_ip=10.100.100.101
    vsc_standby_ip=10.100.100.102
    uplink_interface=eth0
    nuage_openshift_monitor_log_dir=/var/log/nuage-openshift-monitor
    nuage_interface_mtu=1450
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
   
.. Note:: Make sure that all the images are loaded on the nodes & masters using 'docker load -i <docker-image.tar>' command as shown in the Setup section above. If the images are not loaded, the deployment of daemonsets will fail.

2. Verify that the Master-Node connectivity is up and all nodes are running:

   ::
   
       oc login -u system:admin
       oc get nodes


Installation for Multiple Masters
----------------------------------

A High Availability (HA) environment can be configured with multiple masters and multiple nodes.

Nuage OpenShift only supports HA configuration method described in this section. This can be combined with any load balancing solution, the default being HAProxy. In the inventory file, there are two master hosts, the nodes, an etcd server and a host that functions as the HAProxy to balance the master API on all master hosts. The HAProxy host is defined in the [lb] section of the inventory file enabling Ansible to automatically install and configure HAProxy as the load balancing solution.

1. Create the nodes file for Ansible configuration for multiple masters in the openshift-ansible directory with the content shown below.

2. Verify that the image versions are accurate by checking the TAG displayed by 'docker images' output for successful deployment of Nuage daemonsets.

   .. Note:: The following nodes file is provided as a sample. Please update the values with your actual deployment. The below nodes file deploys OpenShift version 3.5. To deploy OpenShift version 3.6, use 'openshift_pkg_version=-3.6.173.0.5'
  
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
        openshift_disable_check=disk_availability,memory_availability,package_version,docker_storage,docker_image_availability
        
        # Nuage specific parameters
        openshift_use_openshift_sdn=False
        openshift_use_nuage=True
        openshift.common._use_nuage=True
        os_sdn_network_plugin_name=cni
        vsd_api_url=https://<VSD-IP/VSD-Hostname>:7443
        vsp_version=v5_0
        
        # The below versions should match the TAG version in the output of 'docker images' on the nodes. See point 2 above
        # Example: nuage_monitor_image_version=5.1.2-70
        nuage_monitor_image_version=<version>
        nuage_vrs_image_version=<version>
        nuage_cni_image_version=<version>
        nuage_infra_image_version=<version>
        
        enterprise=openshift
        domain=openshift
        vsc_active_ip=10.100.100.101
        vsc_standby_ip=10.100.100.102
        uplink_interface=eth0
        nuage_openshift_monitor_log_dir=/var/log/nuage-openshift-monitor
        nuage_interface_mtu=1450
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
    
        # Required for Nuage Monitor REST server and HA
        openshift_master_cluster_method=native
        nuage_openshift_monitor_rest_server_port=9443
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
        node1.nuageopenshift.com openshift_schedulable=True openshift_node_labels="{'region': 'infra'}"
        node2.nuageopenshift.com
        master1.nuageopenshift.com openshift_node_labels="{'install-monitor': 'true'}"
        master2.nuageopenshift.com openshift_node_labels="{'install-monitor': 'true'}"
        

.. Note:: It is mandatory to specify the openshift_node_labels="{'install-monitor': 'true'}" parameter for every master node for Nuage OpenShift master to be deployed.


Installing the VSP Components for Multiple Masters
---------------------------------------------------

1. Run the following command to install the VSP components:

   ::
   
       cd openshift-ansible
       ansible-playbook -vvvv -i nodes playbooks/byo/config.yml

  A successful installation displays the following output:

   ::
   
       2017-08-11 22:01:49,891 p=16545 u=root | PLAY RECAP *********************************************************************
       2017-08-11 22:01:49,892 p=16545 u=root | localhost             : ok=20  changed=0  unreachable=0 failed=0
       2017-08-11 22:01:49,892 p=16545 u=root | master1.nuageopenshift.com : ok=247 changed=22 unreachable=0  failed=0
       2017-08-11 22:01:49,893 p=16545 u=root | master2.nuageopenshift.com : ok=248 changed=22 unreachable=0  failed=0
       2017-08-11 22:01:49,894 p=16545 u=root | node1.nuageopenshift.com : ok=111 changed=21 unreachable=0  failed=0
       2017-08-11 22:01:49,895 p=16545 u=root | node2.nuageopenshift.com : ok=111 changed=21 unreachable=0  failed=0 

.. Note:: Make sure that all the images are loaded on the nodes & masters using 'docker load -i <docker-image.tar>' command as shown in the Setup section above. If the images are not loaded, the deployment of daemonsets will fail.

2. Verify that the Master-Node connectivity is up and all nodes are running:

   ::
   
       oc login -u system:admin
       oc get nodes
   
   .. Note:: Both the masters should display all nodes as connected.

3. Ansible configures the loadbalancer to balance the Openshift Master's 9443 port. 

Deploying the Nuage DaemonSet
--------------------------------

The Ansible installer with automatically label the master nodes and deploy the nuage-master-config, nuage-vrs-ds, nuage-infra-ds and nuage-cni-ds daemonsets. In case of any failures, use the appropriate commands to correct or verify the daemonset files and re-deploy.

The nuage-master-config-daemonset.yaml for openshift-monitor deployment and nuage-node-config-daemonset.yaml for VRS and CNI plugin deployment and nuage-infra-pod-config-daemonset.yaml for nuage-infra pod is copied to /etc/ directory as part of Ansible installation. 
The Nuage infra pod now runs on all nodes to enable access to the service IP from underlay nodes.

The daemonset files are pre-populated using the values provided in the 'nodes' file during Ansible installation. You may modify the image versions or other relevant parameters in the yaml file. However, it is advised to take a back-up of the yaml files before any modification.

1. Verify the daemonset deployment.

   ::   
       
       [root@master]# oc get ds -n kube-system
        NAME                  DESIRED   CURRENT   READY     NODE-SELECTOR          AGE
        nuage-cni-ds             3        3        3        <none>                 7m
        nuage-infra-ds           4        4        2        <none>                 7m
        nuage-master-config      1        1        1        install-monitor=true   7m
        nuage-vrs-ds             3        3        3        <none>                 7m
        
2. Verify that the REST server URL value is correct in the /etc/nuage-node-config-daemonset.yaml file. The 'nuageMonRestServer' should be configured with openshift_master_cluster_hostname value specified in the nodes files during Ansible installation. Modify the value and save the file if this field has incorrect values. Delete and re-deploy the node daemonset as shown in the following steps. 

   ::
   
        # REST server URL
        nuageMonRestServer: https://master.nuageopenshift.com:9443

   .. Note:: If 'nuageMonRestServer' has the value 0.0.0.0:9443, it is incorrect. Please change the value and re-deploy.

3. If you modify the daemonset files, delete and re-deploy the master or node daemonsets respectively using the following commands.

.. Note:: It is mandatory to delete the nuage-infra-ds using the command 'oc delete -f /etc/nuage-infra-pod-config-daemonset.yaml' before deleting nuage-cni-ds or nuage-vrs-ds i.e before you do 'oc delete -f /etc/nuage-node-config-daemonset.yaml'. In case you skipped doing this and there are stale nuage-infra pods in kube-system namespace, refer to the troubleshooting guide.

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
        
         [root@master]# oc create -f /etc/nuage-infra-pod-config-daemonset.yaml 
        daemonset "nuage-infra-ds" created

        [root@master]# oc delete -f /etc/nuage-infra-pod-config-daemonset.yaml 
        daemonset "nuage-infra-ds" deleted
        
        [root@master]# oc get ds -n kube-system
        NAME                  DESIRED   CURRENT   READY     NODE-SELECTOR          AGE
        nuage-cni-ds            3        3        3        <none>                 7m
        nuage-master-config     1        1        1        install-monitor=true   7m
        nuage-vrs-ds            3        3        3        <none>                 7m
        nuage-infra-ds         3        3         3         <none>                 7m
         
4. The master daemonset deploys the nuage-master-config(nuage-openshift-monitor) pod on the master node and the node daemonset deploys the CNI plugin pod and Nuage VRS pod on every slave node. Following is the output of successfully deployed master and node daemonsets.
The Nuage infra pod now runs on all nodes to enable access to the service IP from underlay nodes. 

   ::
        
        [root@master]# oc get all -n kube-system
        NAME                        READY     STATUS    RESTARTS   AGE
        nuage-cni-ds-04s43          1/1       Running   0          7m
        nuage-cni-ds-81mnp          1/1       Running   0          7m
        nuage-cni-ds-f4q2k          1/1       Running   0          7m
        nuage-master-config-0d95v   1/1       Running   0          7m
        nuage-infra-ds-sftn2        1/1       Running   0          7m
        nuage-infra-ds-x6fmr        1/1       Running   0          7m
        nuage-vrs-ds-0v9sq          1/1       Running   0          7m
        nuage-vrs-ds-c0kt5          1/1       Running   0          7m
        nuage-vrs-ds-d4h7m          1/1       Running   0          7m
   
5. If the nuage-infra daemonset is stuck in 'ContainerCreating' stage on the master nodes, you can ignore as the pods are unable to get an overlay IP as the master nodes are probably not being used to actively schedule pods or services. The infra pods are not restricted from running on the masters due a fact that some customers might be interested in using the master nodes to schedule pods or services.    

Post Installation
-----------------------

1. Check the docker-registry and router pods in the default namespace. If they have failed to deploy, delete and re-deploy the docker-registry and router pods. Check the troubleshooting guide for more information.


