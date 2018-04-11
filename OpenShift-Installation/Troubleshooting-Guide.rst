
=====================================================
Troubleshooting Guide for Nuage Openshift Integration
=====================================================

.. contents::
   :local:
   :depth: 3


Openshift Ansible Installation failure
======================================

If Openshift ansible installation fails, first check the nodes/inventory file as shown in the Installation guide. The below section only lists some of the mandatory parameters which we feel contribute to most of the installation issues. Read the #comments carefully.

.. Note:: This is NOT a complete nodes file. The below nodes file only lists certain mandatory parameters. Refer the Installation guide for the complete nodes file.

:: 
   
   [OSEv3:children]
   masters
   nodes
   etcd
   # Mandatory for HA/Multi-Master cluster
   lb
   [OSEv3:vars]
   # Currently Nuage Integration only works with Openshift Version 3.5. It is preferred to use 3.5.5.5
   openshift_pkg_version=-3.7.9
   # Mandatory for Nuage Integration
   openshift_use_openshift_sdn=False
   openshift_use_nuage=True
   openshift.common._use_nuage=True
   os_sdn_network_plugin_name=cni
   # https and 7443 is a must
   vsd_api_url=https://10.31.45.137:7443
   vsp_version=v5_0
   # Are the below version correct ? Are these images loaded on the nodes via docker load -i ?
   nuage_monitor_image_version=5.2.2-63
   nuage_vrs_image_version=5.2.2-70
   nuage_cni_image_version=5.2.2-63
   nuage_infra_image_version=v5.2.2
   # User must create this enterprise on VSD
   enterprise=openshift_scale
   # Below domain will be created by nuage-monitor running on the master nodes
   domain=openshift_scale
   # Is the user created, added to Administrators group and certificates generated and copied to below location ?
   vsd_user_cert_file=/opt/ose-admin.pem
   vsd_user_key_file=/opt/ose-admin-Key.pem
   vsd_user=ose-admin
   # Must for Atomic hosts
   slave_base_host_type=is_atomic
   # Must be <1500
   nuage_interface_mtu=1450
   # Mandatory to mention 9443
   nuage_openshift_monitor_rest_server_port=9443
   openshift_node_kubelet_args={'max-pods': ['110'], 'image-gc-high-threshold': ['90'], 'image-gc-low-threshold': ['80'], 'make-iptables-util-chains': ['false']}
   [masters]
   ovs-10.mvdcdev08.us.alcatel-lucent.com
   ovs-11.mvdcdev08.us.alcatel-lucent.com
   [lb]
   ovs-12.mvdcdev08.us.alcatel-lucent.com 
   [etcd]
   ovs-10.mvdcdev08.us.alcatel-lucent.com
   ovs-11.mvdcdev08.us.alcatel-lucent.com
   [nodes]
   # One of the node must be labelled as below, else docker-registry/router pod will fail to deploy
   ovs-1.mvdcdev08.us.alcatel-lucent.com openshift_schedulable=True openshift_node_labels="{'region': 'infra'}"
   ovs-2.mvdcdev08.us.alcatel-lucent.com
   # Both master nodes must be labelled as below, else nuage-monitor will fail to deploy on master nodes
   ovs-10.mvdcdev08.us.alcatel-lucent.com openshift_node_labels="{'install-monitor': 'true'}"
   ovs-11.mvdcdev08.us.alcatel-lucent.com openshift_node_labels="{'install-monitor': 'true'}"


OpenShift HA Cluster Installation failure
=========================================

If the OpenShift HA cluster installation fails with the following error, rerun the Ansible installer as done in the beginning.

   ::
   
     failed: [ovs-10.mvdcdev35.us.alcatel-lucent.com] (item=policy add-cluster-role-to-user cluster-reader system:serviceaccount:default:nuage) => {
     "changed": false, 
     "cmd": [
     "oc", 
     "adm", 
     "policy", 
     "add-cluster-role-to-user", 
     "cluster-reader", 
     "system:serviceaccount:default:nuage", 
     "--config=/tmp/openshift-ansible-gNgyWHo.kubeconfig"
     ], 
     "delta": "0:00:00.702908", 
     "end": "2016-10-23 18:01:04.446320", 
     "failed": true, 
     "failed_when_result": true, 
     "invocation": {
     "module_args":
     { "_raw_params": "oc adm policy add-cluster-role-to-user cluster-reader system:serviceaccount:default:nuage --config=/tmp/openshift-ansible-gNgyWHo.kubeconfig", "_uses_shell": false, "chdir": null, "creates": null, "executable": null, "removes": null, "warn": true    }
     , 
     "module_name": "command"
     }, 
     "item": "policy add-cluster-role-to-user cluster-reader system:serviceaccount:default:nuage", 
     "rc": 1, 
     "start": "2016-10-23 18:01:03.743412", 
     "stderr": "Error from server: Operation cannot be fulfilled on rolebinding \"cluster-readers\": the object has been modified; please apply your changes to the latest version and try again", 
     "stdout": "", 
     "stdout_lines": [], 
     "warnings": []


OpenShift HA Cluster Installation - HA proxy configuration
==========================================================

In case of OpenShift HA install, Ansible removes the nuage-monitor-server HA proxy configuration from the /etc/haproxy/haproxy.cfg file on the load balancing node. If the CNI plugin on the nodes are unable to connect to the nuage-openshift-monitor on port 9443, kindly check the below HA proxy configuration on your cluster.


1. Add the following configuration at the end of haproxy.cfg file on the lb node: 

   ::

      frontend nuage-monitor-server 
        bind *:9443 
        default_backend nuage-monitor-server 
        mode tcp 
        option tcplog 

      backend nuage-monitor-server 
        balance source 
        mode tcp 
        server master0 <master0 IP>:9443 check 
        server master1 <master1 IP>:9443 check 

2. Restart the service using command service haproxy restart and check the status service haproxy status -l.


Daemonset fails to deploy
=========================

If the nuage-vrs-ds, nuage-cni-ds, nuage-infra-ds or nuage-master-config daemonsets fail to deploy (usually with a host port error), check the Openshift Security Contstraints on the cluster. Run the below commands on the master node, delete and re-deploy the daemonsets.

   ::
   
      oc adm policy add-scc-to-user privileged system:serviceaccount:openshift-infra:daemonset-controller
      oc adm policy add-scc-to-user privileged system:serviceaccount:kube-system:daemon-set-controller


Openshift monitor fails to connect to VSD
==========================================

1. Check for the VSD certificates on the master node in /usr/share/nuage-openshift-monitor/. These certificates are copied from the location mentioned in the nodes file prior to ansible installation for 'vsd_user_cert_file' and 'vsd_user_key_file'

2. In case you missed generating the certs for VSD, re-generate the certs as mentioned in the installation guide and place the certificate and key in /usr/share/nuage-openshift-monitor/ and restart the nuage-master-config pods in the kube-system namespace.

3. Check if the user created above is added to the "Administrators" group in VSD.

4. Check if the vsdApiUrl in /usr/share/nuage-openshift-monitor/nuage-openshift-monitor.yaml is "https://<VSD-IP>:7443" and the vspVersion is correct. If not modify the daemonset file /etc/nuage-master-config-daemonset.yaml and re-deploy the /etc/nuage-master-config-daemonset.yaml daemonset.


nuage-infra-ds pods stuck in Creating/Running/Terminating state
==============================================================

If the nuage-cni-ds and nuage-vrs-ds are deleted before deleting nuage-infra-ds daemonset in the kube-system namespace, the nuage-infra-ds pods might be stuck in Creating/Running/Terminating stage when you try to delete the nuage-infra-ds post CNI or VRS delete. This happens as the infra pods are overlay pods but the user has deleted the VRS and/or the CNI network plugin. To resolve this:

1. Re-deploy the nuage-cni-ds & nuage-vrs-ds and then delete the nuage-infra-ds if it is still stuck.

   :: 
   
      [root@ovs-10 ~]# oc get all
      NAME                        READY     STATUS        RESTARTS   AGE
      nuage-infra-ds-39ldf        1/1       Terminating   0          2h
      nuage-infra-ds-n2jsb        1/1       Terminating   0          2h
      nuage-infra-ds-r516m        1/1       Terminating   0          2h
      nuage-infra-ds-rb0h7        0/1       Terminating   0          2h
      nuage-infra-ds-rv794        1/1       Terminating   0          2h
      nuage-infra-ds-tczp0        0/1       Terminating   0          2h
      nuage-master-config-chdj6   1/1       Running       0          1d
      nuage-master-config-w76wh   1/1       Running       0          1d

      [root@ovs-10 ~]# oc get ds
      NAME                  DESIRED   CURRENT   READY     NODE-SELECTOR                                                               AGE
      nuage-infra-ds        0         0         0         d68559a4-c004-11e7-b7cf-faaca60a8f00=d6855a85-c004-11e7-b7cf-faaca60a8f00   2h
      nuage-master-config   2         2         2         install-monitor=true                                                        1d
      [root@ovs-10 ~]# 

      [root@ovs-10 ~]# oc create -f /etc/nuage-node-config-daemonset.yaml 
      configmap "nuage-config" created
      daemonset "nuage-cni-ds" created
      daemonset "nuage-vrs-ds" created

      [root@ovs-10 ~]# oc get all
      NAME                        READY     STATUS    RESTARTS   AGE
      nuage-cni-ds-8j649          1/1       Running   0          1m
      nuage-cni-ds-9fbc8          1/1       Running   0          1m
      nuage-cni-ds-ds27n          1/1       Running   0          1m
      nuage-cni-ds-kv8nm          1/1       Running   0          1m
      nuage-cni-ds-s9sr8          1/1       Running   0          1m
      nuage-cni-ds-xxknc          1/1       Running   0          1m
      nuage-master-config-chdj6   1/1       Running   0          1d
      nuage-master-config-w76wh   1/1       Running   0          1d
      nuage-vrs-ds-1445b          1/1       Running   0          1m
      nuage-vrs-ds-5qqls          1/1       Running   0          1m
      nuage-vrs-ds-78hcc          1/1       Running   0          1m
      nuage-vrs-ds-92r6g          1/1       Running   0          1m
      nuage-vrs-ds-m3lqg          1/1       Running   0          1m
      nuage-vrs-ds-q2z0t          1/1       Running   0          1m


2. If the nuage-infra pods are still stuck then delete the nuage-infra-ds first and then delete the nuage-cni-ds and nuage-vrs-ds and re-deploy the nuage-cni-ds and nuage-vrs-ds first followed by nuage-infra-ds daemonsets.

   ::
   
      [root@ovs-10 ~]# oc delete -f /etc/nuage-infra-pod-config-daemonset.yaml 
      daemonset "nuage-infra-ds" deleted
      [root@ovs-10 ~]# oc delete -f /etc/nuage-node-config-daemonset.yaml 
      configmap "nuage-config" deleted
      daemonset "nuage-cni-ds" deleted
      daemonset "nuage-vrs-ds" deleted
      [root@ovs-10 ~]# oc get all
      NAME                        READY     STATUS    RESTARTS   AGE
      nuage-master-config-chdj6   1/1       Running   0          1d
      nuage-master-config-w76wh   1/1       Running   0          1d
      [root@ovs-10 ~]# 


CNI plugin pod fails to move to 'Running' state/CNI pod failure
================================================================

1. Check the Nuage CNI network plugin logs at /var/log/cni/nuage-cni.log on every node.

2. If there are any errors which point that the CNI is unable to communicate to the nuage-monitor as shown below

   :: 
   
      |2017-10-30 16:32:48.019744693 -0700 PDT|ERROR|0006|Error occured while sending POST call to Nuage K8S monitor to obtain pod metadata: Post https://ovs-12.example.com:9443/namespaces/kube-system/pods: dial tcp 10.31.45.149:9443: getsockopt: connection refused

   a. Check that the "https://" address is either the master IP/Hostname (in case of standalone master) or the Load Balancer IP/Hostname (in case of Multi-Master).
   b. Check that the port is 9443
   c. Check if the nuage-master-config pod is running on the master nodes using "oc get all -n kube-system" command
   d. Check the HA proxy/Load balancer configuration as shown in the previous section.
   e. Check that the IPTables are flushed IF using "userspace" kubeproxy.
   
2. The logs for CNI Audit daemon (takes care of clearing up any stale entries or sync up issues) can be found at /var/log/cni/nuage-daemon.log


Pods not Resolved/Pods fail to move to 'Running' state/Pods are not assigned an Overlay IP
==========================================================================================

If pods fail to move to 'Running' state or do not get an Overlay IP, check the following configurations:

1. For released prior to VSP 5.1.2, it is mandatory to enable "Underlay Support" and "Address Translation Support" for the Openshift domain on the VSD UI.

2. If the Openshift cluster is set-up with "userspace" kubeproxy, i.e if your nodes/inventory file has openshift_node_proxy_mode='userspace', after successful ansible installation, kindly flush the iptable rules (iptables -F) on all nodes to unblock traffic.

3. Check if the docker insecure-registry is setup correctly. The Openshift host preparation guide explicitly mandates to add '--insecure-registry 172.30.0.0/16' to the /etc/sysconfig/docker file. So follow the OpenShift Host Preparation Guide mentioned in the Installation section to set the insecure-registry and restart docker service on the affected nodes.

4. Check if the VRS and CNI pods are 'Running' in the kube-system namespace

   ::
   
      [root@ovs-10 ~]# oc get all -n kube-system -o wide
      NAME                        READY     STATUS              RESTARTS   AGE       IP             NODE
      nuage-cni-ds-3vmfm          1/1       Running             0          6m        10.31.45.139   ovs-2.mvdcdev08.example.com
      nuage-cni-ds-7wdkh          1/1       Running             0          6m        10.31.45.138   ovs-1.mvdcdev08.example.com
      nuage-cni-ds-bjsk0          1/1       Running             0          6m        10.31.45.140   ovs-3.mvdcdev08.example.com
      nuage-cni-ds-d18tg          1/1       Running             0          6m        10.31.45.148   ovs-11.mvdcdev08.example.com
      nuage-infra-ds-mcdzj        1/1       Running             0          6m        70.70.0.122    ovs-1.mvdcdev08.example.com
      nuage-infra-ds-qbxnj        0/1       Running             0          6m        70.70.0.105    ovs-11.mvdcdev08.example.com
      nuage-infra-ds-tww0n        1/1       Running             0          6m        70.70.0.101    ovs-3.mvdcdev08.example.com
      nuage-infra-ds-v0fvg        1/1       Running             0          6m        70.70.0.97     ovs-2.mvdcdev08.example.com
      nuage-master-config-71kcz   1/1       Running             0          6m        10.31.45.147   ovs-10.mvdcdev08.example.com
      nuage-master-config-sc09f   1/1       Running             0          6m        10.31.45.148   ovs-11.mvdcdev08.example.com
      nuage-vrs-ds-7fgcm          1/1       Running             0          6m        10.31.45.147   ovs-10.mvdcdev08.example.com
      nuage-vrs-ds-gfp0j          1/1       Running             0          6m        10.31.45.139   ovs-2.mvdcdev08.example.com
      nuage-vrs-ds-h2225          1/1       Running             0          6m        10.31.45.148   ovs-11.mvdcdev08.example.com
      nuage-vrs-ds-r8f1q          1/1       Running             0          6m        10.31.45.140   ovs-3.mvdcdev08.example.com

5. Check if all the zones/subnets are created on the VSD Architect UI for the namespaces/projects by the nuage-openshift-monitor/nuage-master-config pod. If not, then check and re-deploy the nuage-master-config daemonset.


Pods stuck in ContainerCreating Stage on the slave nodes
=========================================================

If a pod is stuck in 'ContainerCreating' stage, check the /usr/share/vsp-openshift.yaml or /var/usr/share/vsp-openshift.yaml (on Atomic hosts) for any configuration errors. It should look like below

   ::
   
      clientCert: /var/usr/share//vsp-openshift/client.crt
      # The key to the certificate in clientCert above
      clientKey: /var/usr/share//vsp-openshift/client.key
      # The certificate authority's certificate for the local kubelet.  Usually the
      # same as the CA cert used to create the client Cert/Key pair.
      CACert: /var/usr/share//vsp-openshift/ca.crt
      # Name of the enterprise in which pods will reside
      enterpriseName: openshift_scale
      # Name of the domain in which pods will reside
      domainName: openshift_scale
      # Name of the VSD user in admin group
      vsdUser: ose-admin
      # IP address and port number of master API server
      masterApiServer: https://ovs-12.mvdcdev08.us.alcatel-lucent.com:8443
      # REST server URL 
      nuageMonRestServer: https://ovs-12.mvdcdev08.us.alcatel-lucent.com:9443
      # Bridge name for the docker bridge
      dockerBridgeName: docker0
      # Certificate for connecting to the openshift monitor REST api
      nuageMonClientCert: /var/usr/share//vsp-openshift/nuageMonClient.crt
      # Key to the certificate in restClientCert
      nuageMonClientKey: /var/usr/share//vsp-openshift/nuageMonClient.key
      # CA certificate for verifying the master's rest server
      nuageMonServerCA: /var/usr/share//vsp-openshift/nuageMonCA.crt
      # Nuage vport mtu size
      interfaceMTU: 1450
      # Logging level for the plugin
      # allowed options are: "dbg", "info", "warn", "err", "emer", "off"
      logLevel: 3

If not, make necessary corrections in the /etc/nuage-node-config-daemonset.yaml and re-deploy the daemonset by using the commands "oc delete -f /etc/nuage-node-config-daemonset.yaml" and "oc create -f /etc/nuage-node-config-daemonset.yaml"


Pod to Pod/Service IP communication failure
============================================

After the Ansible Install is done, check for POSTROUTING rules in the iptables -t nat -nvL which should look like this:

   ::
   
      Chain POSTROUTING (policy ACCEPT 6 packets, 360 bytes)
      pkts bytes target     prot opt in     out     source               destination         
      0     0 MASQUERADE  all  --  *      svc-pat-tap  0.0.0.0/0            0.0.0.0/0            mark match 0x2/0x3
      0     0 MASQUERADE  all  --  *      svc-pat-tap  0.0.0.0/0            0.0.0.0/0            mark match 0x3/0x3
      113  8324 MASQUERADE  all  --  *      eth0    0.0.0.0/0            0.0.0.0/0            mark match 0x2/0x3

If you do not see the above rules, restart VRS by deleting the VRS pod on the affected node. The nuage-vrs-ds daemonset should re-deploy the VRS pod on the node. 


Unable to ping the OpenShift Master or a Public IP
=======================================================
If a deployed pod is unable to ping the OpenShift Master or a public IP like 8.8.8.8, check for the following:

1. The bridge being used by Docker: When OpenShift is installed with the default redhat/OpenShift-sdn-subnet plugin it uses the lbr0 bridge, but once the nuage-vsp-openshift plugin is put in Docker, Docker may get out of sync. Restart the Docker service. Otherwise, reboot the node.

2. Routes on the OpenShift Node: Make sure the svc-pat-tap routes and rules are added in nat table as indicated above. Otherwise, restart service openvswitch.


Docker Registry Pod stuck in Pending stage after Ansible install
=================================================================

With the Nuage Installation of OSE 3.5, it is noticed that the masters do not get listed as nodes and are by default marked as unschedulable.

   ::
      
      NAME                                    STATUS    AGE
      ovs-1.mvdcdev44.us.alcatel-lucent.com   Ready     4d
      ovs-2.mvdcdev44.us.alcatel-lucent.com   Ready     4d
      ovs-3.mvdcdev44.us.alcatel-lucent.com   Ready     4d
      ovs-4.mvdcdev44.us.alcatel-lucent.com   Ready     4d
      ovs-5.mvdcdev44.us.alcatel-lucent.com   Ready     4d
      ovs-6.mvdcdev44.us.alcatel-lucent.com   Ready     4d
      ovs-7.mvdcdev44.us.alcatel-lucent.com   Ready     4d
      ovs-8.mvdcdev44.us.alcatel-lucent.com   Ready     4d
      ovs-9.mvdcdev44.us.alcatel-lucent.com   Ready     4d

Due to this issue, the docker registry and router pod may be stuck in **Pending** stage as shown in the following example. 

   :: 
  
      [root@ovs-10 openshift-ansible]# oc get all
      NAME                 REVISION   DESIRED   CURRENT   TRIGGERED BY
      dc/docker-registry   1          0         0         config
      dc/router            2          1         0         config
      
      NAME                   DESIRED   CURRENT   READY     AGE
      rc/docker-registry-1   0         0         0         2d
      rc/router-1            0         0         0         3d
      rc/router-2            0         0         0         2d

      NAME                  CLUSTER-IP      EXTERNAL-IP   PORT(S)                   AGE
      svc/docker-registry   172.30.18.205   <none>        5000/TCP                  2d
      svc/kubernetes        172.30.0.1      <none>        443/TCP,53/UDP,53/TCP     3d
      svc/router            172.30.169.77   <none>        80/TCP,443/TCP,1936/TCP   3d

      NAME                          READY     STATUS    RESTARTS   AGE
      po/docker-registry-1-deploy   0/1       Pending   0          2d
      po/router-2-deploy            0/1       Pending   0          2d

      
      [root@ovs-10 openshift-ansible]# oc describe pod router-3-deploy
      Name:			router-3-deploy
      Namespace:		default
      Security Policy:	restricted
      Node:			/
      Labels:			openshift.io/deployer-pod-for.name=router-3
      Status:			Pending
      IP:			
      Controllers:		<none>
      Containers:
      deployment:
      Image:	openshift3/ose-deployer:v3.5.5.5
      Port:	
      Volume Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from deployer-token-vvgjl (ro)
      Environment Variables:
      KUBERNETES_MASTER:	https://ovs-10.mvdcdev44.us.alcatel-lucent.com:8443
      OPENSHIFT_MASTER:		https://ovs-10.mvdcdev44.us.alcatel-lucent.com:8443
      BEARER_TOKEN_FILE:	/var/run/secrets/kubernetes.io/serviceaccount/token
      OPENSHIFT_CA_DATA:	-----BEGIN CERTIFICATE-----

       -----END CERTIFICATE-----

      OPENSHIFT_DEPLOYMENT_NAME:	router-3
      OPENSHIFT_DEPLOYMENT_NAMESPACE:	default
      Conditions:
      Type		Status
      PodScheduled 	False 
      Volumes:
      deployer-token-vvgjl:
      Type:	Secret (a volume populated by a Secret)
      SecretName:	deployer-token-vvgjl
      QoS Class:	BestEffort
      Tolerations:	<none>
      Events:
      FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason			Message
      ---------	--------	-----	----			-------------	--------	------			-------
      52s		21s		7	{default-scheduler }			Warning		FailedScheduling	pod (router-3-deploy) failed to fit in any node
      fit failure summary on nodes : CheckServiceAffinity (9), MatchNodeSelector (9)

Delete and redeploy the registry pod. 


.. Note:: The OpenShift recommended solution is to mark the master as schedulable. But it is not listed in oc get nodes when installation is done using Nuage openshift-ansible. So it cannot be marked as schedulable. Instead, follow the work around:

1. Delete the docker-registry and router dc, rc, pod, and svc:

   ::  
   
      oc delete dc/router; oc delete svc router; oc delete pod router-1-deploy
      
      oc delete dc/docker-registry; oc delete svc docker-registry

.. Note:: The above delete commands may vary based on the pod names in your cluster.

2. Delete and re-create the service accounts and role bindings.

   ::  
   
      oc delete serviceaccount router;
      
      oadm policy remove-cluster-role-from-user cluster-reader system:serviceaccount:default:router
      
      oc delete clusterrolebinding router-router-role
      
      oc delete serviceaccount registry
      
      oc delete clusterrolebinding registry-registry-role

      oadm policy add-cluster-role-to-user \
      cluster-reader \
      system:serviceaccount:default:router

3. Redeploy the docker-registry and router pod using the following commands:

   ::  
   
      oadm router router --replicas=1  --service-account=router
      
      oadm registry --config=/etc/origin/master/admin.kubeconfig

4. Restart the Openshift master service

   :: 
   
      systemctl restart atomic-openshift-master-api atomic-openshift-master-controllers

.. Note:: Every time the registry is recreated, its service IP changes and you need to restart the OpenShift master service.


Docker Registry Pod fails to deploy due to Issues in DNS
============================================================

When the Docker registry pod fails to deploy with a "dial i/o timeout," it could be a DNS issue. You need to verify the following:

1. Make sure the OpenShift master when started is listening at 0.0.0.0:53 (in /var/log/messages). If not, there is a port conflict and needs to be resolved.

2. If the DNS is functioning fine, then check the iptables MASQUERADE rules. Default needs to be present. You can check using the command iptables -t nat -nvL. An "iptables -F" and "iptables stop/start" followed by a openvswitch restart will help.


Verify the Rules on the OpenShift Nodes after Ansible Installation
=====================================================================
After the OpenShift nodes come up after the Ansible installation, perform the following steps to verify the rules on the nodes.

1. Check the rules on the nodes on a VRS by entering the following command:

         ::
         
             ovs-appctl bridge/dump-flows alubr0 | grep table_id=4
             
     
      The ideal output displays the following rules for OpenShift to function with the Nuage vsp-openshift plugin:

         ::
         
             [root@ovs-1 ~]# ovs-appctl bridge/dump-flows alubr0 | grep table_id=4
             table_id=4, duration=9s, n_packets=7425, cookie:0x2 n_bytes=644627, priority=0,actions=resubmit(,5)
             table_id=4, duration=9s, n_packets=6, cookie:0x2 n_bytes=509, priority=32768,ip,tun_id=0,nw_src=169.254.3.3,actions=resubmit(,5)
             table_id=4, duration=9s, n_packets=0, cookie:0x2 n_bytes=0, priority=32767,nw_src=169.254.3.3,nw_proto=17,actions=move:NXM_NX_TUN_IPV4_SRC[]->NXM_OF_IP_SRC[],learn(table=4,idle_timeout=60,priority=1,eth_type=0x800,nw_proto=17,NXM_OF_IP_SRC[]=NXM_OF_IP_DST[],NXM_OF_IP_DST[]=NXM_OF_IP_SRC[],NXM_OF_UDP_SRC[]=NXM_OF_UDP_DST[],NXM_OF_UDP_DST[]=NXM_OF_UDP_SRC[],load:0xa9fe0303->NXM_OF_IP_DST[],output:NXM_OF_IN_PORT[]),resubmit(,5)
             table_id=4, duration=9s, n_packets=0, cookie:0x2 n_bytes=0, priority=32767,nw_src=169.254.3.3,nw_proto=6,actions=move:NXM_NX_TUN_IPV4_SRC[]->NXM_OF_IP_SRC[],learn(table=4,idle_timeout=180,priority=1,eth_type=0x800,nw_proto=6,NXM_OF_IP_SRC[]=NXM_OF_IP_DST[],NXM_OF_IP_DST[]=NXM_OF_IP_SRC[],NXM_OF_TCP_SRC[]=NXM_OF_TCP_DST[],NXM_OF_TCP_DST[]=NXM_OF_TCP_SRC[],load:0xa9fe0303->NXM_OF_IP_DST[],output:NXM_OF_IN_PORT[]),resubmit(,5)
             
             
      If the above rules are missing from the OVS, and the output is shown as the following display, you need to perform the workaround provided in Step 2
    
         ::
         
             [root@ovs-1 ~]# ovs-appctl bridge/dump-flows alubr0 | grep table_id=4
             table_id=4, duration=6697s, n_packets=7171, cookie:0x1 n_bytes=621164, priority=0,actions=resubmit(,5)
             
2. Perform the workaround on the primary controller for the rules to appear:

         ::
         
             *A:Dut-H# configure vswitch-controller shutdown 
             *A:Dut-H# configure vswitch-controller no shutdown

Uninstall a Nuage Openshift Cluster
===================================================

To uninstall the Nuage Openshift Cluster, delete the Nuage components manually and then uninstall Openshift using the ansible playbook.

.. Note:: It is mandatory to delete all user created zones, subnets, policies etc on the Openshift Cluster and VSD Domain manually before uninstalling the Openshift cluster. This is to ensure consistency between Nuage VSD & the Openshift cluster data.

Follow the steps below to uninstall a Nuage Openshift cluster running on RHEL worker & RHEL master nodes.

1. Delete any user created projects on the cluster and/or corresponding zones on the VSD

         ::
            
               [root@ovs-master ~]# oc delete project sales
               project "sales" deleted

2. Delete the Nuage daemonsets

         ::
         
               oc delete -f /etc/nuage-infra-pod-config-daemonset.yaml
               oc delete -f /etc/nuage-node-config-daemonset.yaml
               oc delete -f /etc/nuage-master-config-daemonset.yaml
               
3. Delete the files and directories listed below on the master nodes

         :: 
		
               /usr/share/vsp-openshift/
               /etc/default/nuage-cni.yaml 
               /etc/nuage-infra-pod-config-daemonset.yaml
               /etc/nuage-node-config-daemonset.yaml
               /etc/nuage-master-config-daemonset.yaml
               /usr/share/nuage-openshift-ca/           
               /usr/share/nuage-openshift-certificates/ 
               /usr/share/nuage-openshift-monitor/
               /opt/cni/
               /etc/cni/ 
          
4. Delete the files and directories listed below on the worker nodes

           ::
		
               /usr/share/vsp-openshift/
               /etc/default/openvswitch
               /etc/default/nuage-cni.yaml
               /opt/cni/
               /etc/cni/
               
5. Run the Openshift ansible uninstall playbook
 
           ::
            
               ansible-playbook -vvvv -i nodes playbooks/adhoc/uninstall.yml

               PLAY RECAP 
               *********************************************
               ovs-master.mvdcdev08.us.alcatel-lucent.com : ok=56   changed=18   unreachable=0    failed=0   
               ovs-worker1.mvdcdev08.us.alcatel-lucent.com : ok=33   changed=8    unreachable=0    failed=0   
               ovs-worker2.mvdcdev08.us.alcatel-lucent.com : ok=33   changed=8    unreachable=0    failed=0   
       
6. Run the Openshift ansible install playbook
      
            ::
               
               ansible-playbook -vvvv -i nodes playbooks/byo/config.yml
                
               PLAY RECAP 
               ************************************************
               localhost                  : ok=12   changed=0    unreachable=0    failed=0   
               ovs-master.mvdcdev08.us.alcatel-lucent.com : ok=577  changed=227  unreachable=0    failed=0   
               ovs-worker1.mvdcdev08.us.alcatel-lucent.com : ok=199  changed=59   unreachable=0    failed=0   
               ovs-worker2.mvdcdev08.us.alcatel-lucent.com : ok=189  changed=55   unreachable=0    failed=0   


Follow the steps below to uninstall a Nuage Openshift cluster running on Atomic worker & RHEL master nodes.

1. Delete any user created projects on the cluster and/or corresponding zones on the VSD

         ::
            
               [root@ovs-master ~]# oc delete project sales
               project "sales" deleted

2. Delete the Nuage daemonsets

         ::
         
               oc delete -f /etc/nuage-infra-pod-config-daemonset.yaml
               oc delete -f /etc/nuage-node-config-daemonset.yaml
               oc delete -f /etc/nuage-master-config-daemonset.yaml
               
3. Delete the files and directories listed below on the master nodes

         :: 
         
               /var/usr/share/vsp-openshift/
               /usr/share/vsp-openshift/
               /etc/default/nuage-cni.yaml 
               /etc/nuage-infra-pod-config-daemonset.yaml
               /etc/nuage-node-config-daemonset.yaml
               /etc/nuage-master-config-daemonset.yaml
               /usr/share/nuage-openshift-ca/           
               /usr/share/nuage-openshift-certificates/ 
               /usr/share/nuage-openshift-monitor/
               /opt/cni/
               /etc/cni/ 
          
4. Delete the files and directories listed below on the worker nodes

           ::
           
               /var/usr/share/vsp-openshift/
               /etc/default/openvswitch
               /etc/default/nuage-cni.yaml
               /opt/cni/
               /etc/cni/
               
5. Run the Openshift ansible uninstall playbook
 
           ::
            
               ansible-playbook -vvvv -i nodes playbooks/adhoc/uninstall.yml

               PLAY RECAP 
               *********************************************
               ovs-master.mvdcdev08.us.alcatel-lucent.com : ok=56   changed=18   unreachable=0    failed=0   
               ovs-worker1.mvdcdev08.us.alcatel-lucent.com : ok=33   changed=8    unreachable=0    failed=0   
               ovs-worker2.mvdcdev08.us.alcatel-lucent.com : ok=33   changed=8    unreachable=0    failed=0   
       
6. Run the Openshift ansible install playbook
      
            ::
               
               ansible-playbook -vvvv -i nodes playbooks/byo/config.yml
                
               PLAY RECAP 
               ************************************************
               localhost                  : ok=12   changed=0    unreachable=0    failed=0   
               ovs-master.mvdcdev08.us.alcatel-lucent.com : ok=577  changed=227  unreachable=0    failed=0   
               ovs-worker1.mvdcdev08.us.alcatel-lucent.com : ok=199  changed=59   unreachable=0    failed=0   
               ovs-worker2.mvdcdev08.us.alcatel-lucent.com : ok=189  changed=55   unreachable=0    failed=0   
