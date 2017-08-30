.. _kubernetes-8-troubleshooting:

.. include:: ../lib/doc-includes/VSDA-icons.inc



=======================
Troubleshooting
=======================

.. contents::
   :local:
   :depth: 3
   
Kubeadm join or Kubeadm init fails
======================================

If the Kubeadm join or Kubeadm init fails with the following error:

    ::
    
         kubeadm join --token 54ff8e.be449c7fded0d336 10.10.100.1:6443
         kubeadm] WARNING: kubeadm is in beta, please do not use it for production clusters.
         [preflight] Running pre-flight checks
         [preflight] Some fatal errors occurred:
 	         /proc/sys/net/bridge/bridge-nf-call-iptables contents are not set to 1
         [preflight] If you know what you are doing, you can skip pre-flight checks with `--skip-preflight-checks`

You can use the following command to set /proc/sys/net/bridge/bridge-nf-call-iptables to 1:
       
        sysctl net.bridge.bridge-nf-call-iptables=1


Kubelet service is failing
======================================

If the kubelet service is failing with the following error:

    ::
      
        [root@ovs-1 ~]# service kubelet status -l
        Redirecting to /bin/systemctl status  -l kubelet.service
        \u25cf kubelet.service - kubelet: The Kubernetes Node Agent
        Loaded: loaded (/etc/systemd/system/kubelet.service; enabled; vendor preset: disabled)
        Drop-In: /etc/systemd/system/kubelet.service.d
           \u2514\u250010-kubeadm.conf
         Active: active (running) since Fri 2017-06-23 20:25:43 UTC; 1min 5s ago
         Docs: http://kubernetes.io/docs/
         Main PID: 690 (kubelet)
         CGroup: /system.slice/kubelet.service
           \u251c\u2500690 /usr/bin/kubelet --kubeconfig=/etc/kubernetes/kubelet.conf --require-kubeconfig=true --pod-manifest-path=/etc/kubernetes/manifests --allow-privileged=true --network-plugin=cni --cni-bin-dir=/usr/bin/ --make-iptables-util-chains=false --cluster-dns=192.168.0.10 --cluster-domain=cluster.local --cgroup-driver=systemd
           \u2514\u2500708 journalctl -k -f

           Jun 23 20:26:13 ovs-1.k8s.nuagenetworks.com kubelet[690]: E0623 20:26:13.358555     690 kubelet_network.go:193] checkLimitsForResolvConf: Resolv.conf file '/etc/resolv.conf' contains search line consisting of more than 3 domains!
           Jun 23 20:26:13 ovs-1.k8s.nuagenetworks.com kubelet[690]: I0623 20:26:13.360813     690 kubelet.go:1752] skipping pod synchronization - [Failed to start ContainerManager failed to initialise top level QOS containers: failed to create top level Burstable QOS cgroup : Unit kubepods-burstable.slice already exists.]
           Jun 23 20:26:18 ovs-1.k8s.nuagenetworks.com kubelet[690]: I0623 20:26:18.361115     690 kubelet.go:1752] skipping pod synchronization - [Failed to start ContainerManager failed to initialise top level QOS containers: failed to create top level Burstable QOS cgroup : Unit kubepods-burstable.slice already exists.]
           Jun 23 20:26:23 ovs-1.k8s.nuagenetworks.com kubelet[690]: I0623 20:26:23.361450     690 kubelet.go:1752] skipping pod synchronization - [Failed to start ContainerManager failed to initialise top level QOS containers: failed to create top level Burstable QOS cgroup : Unit kubepods-burstable.slice already exists.]
           Jun 23 20:26:28 ovs-1.k8s.nuagenetworks.com kubelet[690]: I0623 20:26:28.361790     690 kubelet.go:1752] skipping pod synchronization - [Failed to start ContainerManager failed to initialise top level QOS containers: failed to create top level Burstable QOS cgroup : Unit kubepods-burstable.slice already exists.]
           Jun 23 20:26:33 ovs-1.k8s.nuagenetworks.com kubelet[690]: I0623 20:26:33.362101     690 kubelet.go:1752] skipping pod synchronization - [Failed to start ContainerManager failed to initialise top level QOS containers: failed to create top level Burstable QOS cgroup : Unit kubepods-burstable.slice already exists.]
           Jun 23 20:26:38 ovs-1.k8s.nuagenetworks.com kubelet[690]: I0623 20:26:38.362391     690 kubelet.go:1752] skipping pod synchronization - [Failed to start ContainerManager failed to initialise top level QOS containers: failed to create top level Burstable QOS cgroup : Unit kubepods-burstable.slice already exists.]
           Jun 23 20:26:43 ovs-1.k8s.nuagenetworks.com kubelet[690]: E0623 20:26:43.358816     690 kubelet_network.go:193] checkLimitsForResolvConf: Resolv.conf file '/etc/resolv.conf' contains search line consisting of more than 3 domains!
           Jun 23 20:26:43 ovs-1.k8s.nuagenetworks.com kubelet[690]: I0623 20:26:43.363039     690 kubelet.go:1752] skipping pod synchronization - [Failed to start ContainerManager failed to initialise top level QOS containers: failed to create top level Burstable QOS cgroup : Unit kubepods-burstable.slice already exists.]
           Jun 23 20:26:48 ovs-1.k8s.nuagenetworks.com kubelet[690]: I0623 20:26:48.363334     690 kubelet.go:1752] skipping pod synchronization - [Failed to start ContainerManager failed to initialise top level QOS containers: failed to create top level Burstable QOS cgroup : Unit kubepods-burstable.slice already exists.]

You can execute the following command:

            for i in $(/usr/bin/systemctl list-unit-files --no-legend --no-pager -l | grep --color=never -o .*.slice | grep kubepod);do systemctl stop $i;done

            service kubelet restart


Stuck in ContainerCreating Stage
======================================

If a pod is stuck in 'ContainerCreating' stage, check the nuage-cni.yaml for any configuration errors.

Unable to ping the Kubernetes Master or a Public IP
=======================================================

If a pod deployed is unable to ping the Kubernetes Master or a public IP like 8.8.8.8, check for the following:

* The bridge being used by Docker: When Kubernetes is installed with the default redhat/kubernetes-sdn-subnet plugin it uses the lbr0 bridge, but once the nuage-cni-k8s plugin is put in Docker, Docker may get out of sync. A restart of the Docker service should help, if not reboot the node.
* Routes on the Kubernetes Node: Make sure the svc-pat-tap routes and rules are added. If not a service openvswitch restart should help.

Verify the Rules on the Kubernetes Nodes after Ansible Installation
=====================================================================

After the Kubernetes nodes come up after the Ansible installation, perform the following steps to verify the rules on the nodes.

:Step 1: Check the rules on the nodes on a VRS by entering the following command:

         ::
         
             ovs-appctl bridge/dump-flows alubr0 | grep table_id=4
             
     
      The ideal output displays the following rules for Kubernetes to function with the nuage-cni-k8s plugin:

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

The nuagekubemon process crashes on the master node after the Ansible Installation
===================================================================================

1. Check the permissions of the Nuage specific key and certificate under /etc/kubernetes/certs. Make sure the permission on the key and certificate is set to 0666.

2. The iptable rules need to be added on the master and slave nodes to permit the access to the VXLAN port (4789) and nuagekubemon port (9443).


Verifying if Nuage CNI, nuagekubemon & VRS pods got deployed correctly
======================================================================

::
 
      		[root@ovs-10 ~]# kubectl get ds -n kube-system
               NAME                      DESIRED   CURRENT   READY     UP-TO-DATE   AVAILABLE    NODE-SELECTOR          AGE
		kube-proxy                4         4         4         4            4           <none>                 2d
		nuage-cni-ds              4         4         4         4            4           <none>                 1d
		nuage-master-monitor-ds   1         1         1         1            1           install-monitor=true   1d
		nuage-vrs-ds              4         4         4         4            4           <none>                 1d

		[root@ovs-10 ~]# kubectl get pods -n kube-system -o wide
		NAME                                                             READY     STATUS    RESTARTS   AGE       IP            NODE
		etcd-ovs-10.k8s.test.com                                         1/1       Running   0          2d        10.100.100.42   ovs-10.k8s.test.com
		kube-apiserver-ovs-10.k8s.test.com                               1/1       Running   0          2d        10.100.100.42   ovs-10.k8s.test.com
		kube-controller-manager-ovs-10.k8s.test.com                      1/1       Running   1          2d        10.100.100.42   ovs-10.k8s.test.com
		kube-dns-2425271678-hxb7k                                        3/3       Running   0          1d        70.70.10.68     ovs-1.k8s.test.com
		kube-proxy-1nlgq                                                 1/1       Running   0          2d        10.100.100.36   ovs-2.k8s.test.com
		kube-proxy-2vm4p                                                 1/1       Running   0          2d        10.100.100.39   ovs-3.k8s.test.com
		kube-proxy-2zbx0                                                 1/1       Running   0          2d        10.100.100.42   ovs-10.k8s.test.com
		kube-proxy-76g24                                                 1/1       Running   0          2d        10.100.100.34   ovs-1.k8s.test.com
		kube-scheduler-ovs-10.k8s.test.com                               1/1       Running   1          2d        10.100.100.42   ovs-10.k8s.test.com
		nuage-cni-ds-1tmcr                                               1/1       Running   0          1d        10.100.100.36   ovs-2.k8s.test.com
		nuage-cni-ds-g40gk                                               1/1       Running   0          1d        10.100.100.39   ovs-3.k8s.test.com
		nuage-cni-ds-vz8h8                                               1/1       Running   0          1d        10.100.100.42   ovs-10.k8s.test.com
		nuage-cni-ds-w2ml8                                               1/1       Running   0          1d        10.100.100.34   ovs-1.k8s.test.com
		nuage-master-monitor-ds-vvh1p                                    1/1       Running   0          1d        10.100.100.42   ovs-10.k8s.test.com
		nuage-vrs-ds-g4csq                                               1/1       Running   0          1d        10.100.100.36   ovs-2.k8s.test.com
		nuage-vrs-ds-hhdpq                                               1/1       Running   0          1d        10.100.100.34   ovs-1.k8s.test.com
		nuage-vrs-ds-l1n6r                                               1/1       Running   0          1d        10.100.100.42   ovs-10.k8s.test.com
		nuage-vrs-ds-qbl3g                                               1/1       Running   0          1d        10.100.100.39   ovs-3.k8s.test.com



Nuage CNI Logs
===============

* Detailed logs for Nuage CNI plugin can be found at /var/log/cni/nuage-cni.log

* Detailed logs for Nuage CNI audit daemon can be found at /var/log/cni/nuage-daemon.log

Nuage VRS access
=================

* With Nuage VRS running as a pod, in order to execute ovs commands inside the container, follow the steps given below:

::
    
    
      root@ovs-8 ~]# docker ps
      CONTAINER ID        IMAGE                         COMMAND                  CREATED             STATUS              PORTS               NAMES
      7f8a90a4da8f        nuage/vrs:5.1.1-3             "/usr/bin/supervisord"   5 hours ago         Up 5 hours                              k8s_install-nuage-vrs.3ad682e4_nuage-vrs-ds-hwf64_kube-system_5c2fdb7c-8da8-11e7-a9f8-faaca6105000_f5eda7f1
      9b15cbc5321f        nuage/cni:0.0.1-PR12_5ea4f6   "/install-cni.sh nuag"   5 hours ago         Up 5 hours                              k8s_install-nuage-cni.c496599_nuage-cni-ds-m2qvh_kube-system_5c2f6734-8da8-11e7-a9f8-faaca6105000_01c2a1bd
      cf7659038bbc        openshift3/ose-pod:v3.5.5.5   "/pod"                   5 hours ago         Up 5 hours                              k8s_POD.ef3fdbfd_nuage-vrs-ds-hwf64_kube-system_5c2fdb7c-8da8-11e7-a9f8-faaca6105000_2c72e720
      c6528662250a        openshift3/ose-pod:v3.5.5.5   "/pod"                   5 hours ago         Up 5 hours                              k8s_POD.ef3fdbfd_nuage-cni-ds-m2qvh_kube-system_5c2f6734-8da8-11e7-a9f8-faaca6105000_a0d9c7d6
 
      
      [root@ovs-8 ~]# docker exec -it 7f8a90a4da8f /bin/bash
      [root@ovs-8 /]# ovs-vsctl show
      338bac29-a82a-450a-b98f-c80945ef3ecd
         Bridge "alubr0"
            Controller "ctrl2"
                  target: "tcp:10.100.100.102:6633"
                  role: master
                  is_connected: true
            Controller "ctrl1"
                  target: "tcp:10.100.100.101:6633"
                  role: slave
                  is_connected: true
            Port "svc-rl-tap2"
                  Interface "svc-rl-tap2"
            Port "alubr0"
                  Interface "alubr0"
                     type: internal
            Port "svc-rl-tap1"
                  Interface "svc-rl-tap1"
            Port svc-spat-tap
                  Interface svc-spat-tap
                     type: internal
            Port svc-pat-tap
                  Interface svc-pat-tap
                     type: internal
            ovs_version: "5.1.1-3-nuage"
     
      [root@ovs-8 /]# ovs-appctl container/port-show
      Name: my-nginx-3147148373-j8mls	UUID: 5ac67fbb6c72850d236e378dd495934d13b386931f9d4e12942df653be4e5b18
	      port-UUID: e1dcd5ef-bed2-4aee-8e51-a841aa9b9d57	Name: 07db556744f4c70	MAC: 7a:66:fc:89:e4:d0
	      Bridge: alubr0	port: 23	flags: 0x0	stats-interval: 60
	      vrf_id: 2134951344	evpn_id: 1206200683	flow_flags: 0x21e64004	flood_gen_id: 0x1
	      IP: 70.70.2.145	 subnet: 255.255.252.0	 GW: 70.70.0.1
	      rate: 4294967295 kbit/s	burst:4294967295 kB	class:0	mac_count: 1
	      BUM rate: 4294967295 kbit/s	BUM peak: 4294967295 kbit/s	BUM burst: 4294967295 kB
	      FIP rate: 4294967295 kbit/s	FIP peak: 4294967295 kbit/s	FIP burst: 4294967295 kB
	      FIP Egress rate: 4294967295 kbit/s	FIP Egress peak: 4294967295 kbit/s	FIP Egress burst: 4294967295 kB
	      Trusted: false	Rewrite: false
	      RX packets:8 errors:0 dropped:4 rl_dropped:0 
	      TX packets:11 errors:0 dropped:0
	      RX bytes:648      TX bytes:882
	      anti-spoof: Enabled
	      policy group tags: 0x177877e8 0x65df994a 0x5b5ac7d0 0x4000
	      policy group domain_id: 0xcefe3
	      route_id: 0x30
	      class_id: 13(0xd)
         
*  Detailed logs for Nuage VRS can be found at /var/log/openvswitch/ovs-vswitchd.log
