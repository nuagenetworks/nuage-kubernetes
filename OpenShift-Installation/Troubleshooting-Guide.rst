
=======================
Troubleshooting
=======================

.. contents::
   :local:
   :depth: 3

OpenShift iptables rules
==========================

If the Openshift Installation is with Nuage release 5.0.2 or prior, default proxy mode is set to userspace but Openshift install itself adds iptables rules on the master and the nodes. Since, Nuage is managing the datapath forwarding, it is necessary to flush the iptables rules on all the nodes and the master using the command ``iptables -F`` 

.. Note:: This applies only to Nuage Installation with userspace kube-proxy. With Nuage release 5.1.1 and beyond, iptables based kube-proxy is used by default so iptables should not be flushed.

Insecure registry config in /etc/sysconfig/docker
==================================================

One of the most important steps in host preparation for Openshift is adding the insecure-registry option in the /etc/sysconfig/docker on the masters & nodes being used for Openshift. If the docker-registry if failing to deploy, missing insecure-registry option might be the root cause. Please follow the steps for host preparation mentioned in the Openshift Installation guide on the Openshift website

Enabling Address translation support & Underlay support on the VSD
===================================================================

If the Openshift Installation is with Nuage release 5.1.1 or prior, Address translation support & Underlay support needs to be enabled on the Openshift domain in the VSD Architect manually. This is required to allow the pods to reach the internet via the underlay in order to pull and push images from the docker registry


OpenShift HA Cluster Installation Fails
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

OpenShift HA Cluster Installation HA proxy config
=================================================

In case of OpenShift HA install, Ansible removes the nuage-monitor-server HA proxy configuration from the /etc/haproxy/haproxy.cfg file on the lb node. 

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

   
POSTROUTING Rules for OSE
==============================

After the Ansible Install is done, check for POSTROUTING rules in the iptables -t nat -nvL which should look like this:

   ::
   
      Chain POSTROUTING (policy ACCEPT 6 packets, 360 bytes)
      pkts bytes target     prot opt in     out     source               destination         
      0     0 MASQUERADE  all  --  *      svc-pat-tap  0.0.0.0/0            0.0.0.0/0            mark match 0x2/0x3
      0     0 MASQUERADE  all  --  *      svc-pat-tap  0.0.0.0/0            0.0.0.0/0            mark match 0x3/0x3
      113  8324 MASQUERADE  all  --  *      eth0    0.0.0.0/0            0.0.0.0/0            mark match 0x2/0x3

If you do not see the above rules, restart openvswitch. 

Pod Stuck in ContainerCreating Stage
======================================
If a pod is stuck in 'ContainerCreating' stage, check the vsp-openshift.yaml for any configuration errors.

Unable to ping the OpenShift Master or a Public IP
=======================================================
If a deployed pod is unable to ping the OpenShift Master or a public IP like 8.8.8.8, check for the following:

* The bridge being used by Docker: When OpenShift is installed with the default redhat/OpenShift-sdn-subnet plugin it uses the lbr0 bridge, but once the nuage-vsp-openshift plugin is put in Docker, Docker may get out of sync. Restart the Docker service. Otherwise, reboot the node.
* Routes on the OpenShift Node: Make sure the svc-pat-tap routes and rules are added in nat table as indicated above. Otherwise, restart service openvswitch.

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

1. Delete the router dc, rc, pod, and svc and re-deploy using the following commands:

      oc delete dc/router; oc delete svc router; oc delete pod router-1-deploy
      
      oc delete dc/docker-registry; oc delete svc docker-registry

2. Delete and re-create the service accounts and role bindings.

      oc delete serviceaccount router;
      
      oadm policy remove-cluster-role-from-user cluster-reader system:serviceaccount:default:router
      
      oc delete clusterrolebinding router-router-role
      
      oc delete serviceaccount registry
      
      oc delete clusterrolebinding registry-registry-role

      oadm policy add-cluster-role-to-user \
      cluster-reader \
      system:serviceaccount:default:router

3. Redeploy the docker-registry and router pod using the following commands:

      oadm router router --replicas=1  --service-account=router
      
      oadm registry --config=/etc/origin/master/admin.kubeconfig

.. Note:: Every time the registry is recreated, its service IP changes and you need to restart the OpenShift master service.


Docker Registry Pod fails to deploy due to Issues in DNS
============================================================

When the Docker registry pod fails to deploy with a "dial i/o timeout," it could be a DNS issue. You need to verify the following:

* Make sure the OpenShift master when started is listening at 0.0.0.0:53 (in /var/log/messages). If not, there is a port conflict and needs to be resolved.
* If the DNS is functioning fine, then check the iptables MASQUERADE rules. Default needs to be present. You can check using the command iptables -t nat -nvL. An "iptables -F" and "iptables stop/start" followed by a openvswitch restart will help.


Verify the Rules on the OpenShift Nodes after Ansible Installation
=====================================================================
After the OpenShift nodes come up after the Ansible installation, perform the following steps to verify the rules on the nodes.

:Step 1: Check the rules on the nodes on a VRS by entering the following command:

         ::
         
             ovs-appctl bridge/dump-flows alubr0 | grep table_id=4
             
     
      The ideal output displays the following rules for OpenShift to function with the Nuage vsp-openshift plugin:

         ::
         
             [root@ovs-1 ~]# ovs-appctl bridge/dump-flows alubr0 | grep table_id=4
             table_id=4, duration=9s, n_packets=7425, cookie:0x2 n_bytes=644627, priority=0,actions=resubmit(,5)
             table_id=4, duration=9s, n_packets=6, cookie:0x2 n_bytes=509, priority=32768,ip,tun_id=0,nw_src=169.254.3.3,actions=resubmit(,5)
             table_id=4, duration=9s, n_packets=0, cookie:0x2 n_bytes=0, priority=32767,nw_src=169.254.3.3,nw_proto=17,actions=move:NXM_NX_TUN_IPV4_SRC[]->NXM_OF_IP_SRC[],learn(table=4,idle_timeout=60,priority=1,eth_type=0x800,nw_proto=17,NXM_OF_IP_SRC[]=NXM_OF_IP_DST[],NXM_OF_IP_DST[]=NXM_OF_IP_SRC[],NXM_OF_UDP_SRC[]=NXM_OF_UDP_DST[],NXM_OF_UDP_DST[]=NXM_OF_UDP_SRC[],load:0xa9fe0303->NXM_OF_IP_DST[],output:NXM_OF_IN_PORT[]),resubmit(,5)
             table_id=4, duration=9s, n_packets=0, cookie:0x2 n_bytes=0, priority=32767,nw_src=169.254.3.3,nw_proto=6,actions=move:NXM_NX_TUN_IPV4_SRC[]->NXM_OF_IP_SRC[],learn(table=4,idle_timeout=180,priority=1,eth_type=0x800,nw_proto=6,NXM_OF_IP_SRC[]=NXM_OF_IP_DST[],NXM_OF_IP_DST[]=NXM_OF_IP_SRC[],NXM_OF_TCP_SRC[]=NXM_OF_TCP_DST[],NXM_OF_TCP_DST[]=NXM_OF_TCP_SRC[],load:0xa9fe0303->NXM_OF_IP_DST[],output:NXM_OF_IN_PORT[]),resubmit(,5)
             
             
      If the above rules are missing from the OVS, and the output is shown as the following display, you need to perform the workaround provided in Step 2:
    
         ::
         
             [root@ovs-1 ~]# ovs-appctl bridge/dump-flows alubr0 | grep table_id=4
             table_id=4, duration=6697s, n_packets=7171, cookie:0x1 n_bytes=621164, priority=0,actions=resubmit(,5)
             
:Step 2: Perform the workaround on the primary controller for the rules to appear:

         ::
         
             *A:Dut-H# configure vswitch-controller shutdown 
             *A:Dut-H# configure vswitch-controller no shutdown
 
Nuage CNI Logs
===============

* Detailed logs for Nuage CNI plugin can be found at /var/log/cni/nuage-cni.log

* Detailed logs for Nuage CNI audit daemon can be found at /var/log/cni/nuage-daemon.log

