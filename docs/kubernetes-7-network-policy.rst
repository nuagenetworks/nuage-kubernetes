.. _openshift-7-network-policy:

.. include:: ../lib/doc-includes/VSDA-icons.inc


   
.. ========= ======= =======  ==========
.. Date      Version Author   Reviewers
.. ========= ======= =======  ==========

.. 03/07/17  4.0.R8  Madhu    Vishal, Aniket, Harmeet [Adding namespace]


=======================
Network Policy API
=======================

.. contents::
   :local:
   :depth: 3
   


Network Policy API Support 
===========================

Nuage VSP supports the Kubernetes Network Policy API https://kubernetes.io/docs/user-guide/networkpolicies. Specifically the Policy API enables a k8s cluster administrator to specify which group of pods can communicate with each other. The grouping of pods is achieved via the use of labels while the network policy is specified in terms of a YAML file. 

The pods and network policies can be created in any order and are independent of each other. When a network policy is created, it is instantly applied to all the pods with the matching labels. Subsequently, when new pods with labels get spawned after the creation of the network policy, the policy is applied to the new pods with matching labels.

The labels attached to the pod are used to create Nuage Policy Groups into which the vports of the pods are added. While the policy rules mentioned in the policy YAML are implemented as ACLs on the VSD. 

The above concepts can be understood with the following two examples.

Example 1. Enable Intra-namespace Pod Communication
---------------------------------------------------

By default Nuage enables all the pods in a namespace to communicate with each other. However this functionality can be disabled using the VSD and adding appropriate ingress and egress ACLs. The following example assumes the blocking ACLs have been added using the VSD.

   ::

      kubectl run my-nginx --image=nginx --replicas=2 --port=80 --labels=role=web-service
      kubectl run busybox --image=busybox --replicas=1 --labels=role=client
      
The above commands spawns ngnix and busybox pods in the default namespace and if a user tries to wget from the busybox pod to the nginx pod, the wget can fail. However the communication can be enabled by specifying the following policy:
   
   ::  

      kind: NetworkPolicy
      apiVersion: extensions/v1beta1
      metadata:
        name: allow-tcp-80
        namespace: default
        labels:
          "nuage.io/priority": "500"
      spec:
        podSelector:
          matchLabels:
            role: web-service
        ingress:
          - ports:
            - protocol: TCP
              port: 80    
            from:
              - podSelector:
                  matchLabels:
                    role: client

The network policy is created on Kubernetes using the following command:

   ::

      kubectl create -f network-policy.yaml

The above command creates a policy group for the ngnix pods and another policy group for the busybox. Appropriate ACLs are added to enable the http communication between the two set of pods. 

.. Note:: The priority for each of the network policy should be unique, as it maps to the ACL template priorities on the VSD. Also lower value of the priority implies a high priority in Nuage VSP. For more information on ACL priorities, see the "ACL Sandwich" section of the VSP User Guide.

Example 2. Enable Inter-namespace Pod Communication
----------------------------------------------------

   ::

      kubectl create namespace web-ns
      kubectl run my-nginx --image=nginx --replicas=2 --port=80 --namespace=web-ns --labels=role=web-service
      
      kubectl create namespace client-ns 
      kubectl label namespace client-ns project=web-client
      kubectl run busybox --image=busybox --replicas=1 --namespace=client-ns --labels=role=web-client
      

By default Nuage blocks all of the communication between pods across namespaces. As a result a "wget" to the ngnix server in the web-ns from the busybox client in the client-ns can fail. In order to enable the communication between the two set of pods a network policy can be defined in a YAML as shown below:

   ::

      ---
      kind: NetworkPolicy
      apiVersion: extensions/v1beta1
      metadata:
        name: allow-inter-ns-web-access
        namespace: web-ns
        labels:
          "nuage.io/priority": "500"
      spec:
        podSelector:  
          matchLabels:  
            service: web
        ingress:
          - ports:
              - protocol: TCP
                port: 80
            from:
              - podSelector:  
                  matchLabels:
                    role: web-client  

The network policy is created on Kubernetes using the following command:

   ::

      kubectl create -f network-policy.yaml

The above command creates the policy-groups (corresponding to the pod labels) and addition of appropriate ACLs between them to facilitate the communication.

Example 3. Namespace Annotations 
--------------------------------

The Nuage k8s solution also supports namespace annotations, which allows an administrator to block all of the intra namespace communication by default when the isolation attribute for the namespace annotation is set to `DefaultDeny`. The namespace annotation can be enabled using the following command:

    ::

        kubectl annotate ns <namespace> "net.beta.kubernetes.io/network-policy={\"ingress\": {\"isolation\": \"DefaultDeny\"}}"

When the namespace is annotated, ingress and egress ACLs are added to the Nuage VSD to block all of the intra namespace communication. The namespace annotations for a namespace can be removed using the following command.

    ::

        kubectl annotate --overwrite ns <namespace> "net.beta.kubernetes.io/network-policy=-"
