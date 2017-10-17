
=======================
Network Policy API
=======================

.. contents::
   :local:
   :depth: 3
   


Network Policy API Support 
===========================

Nuage VSP Release 4.0.R8 and later adds support for the Kubernetes Network Policy API https://kubernetes.io/docs/user-guide/networkpolicies with OpenShift. Specifically the Policy API enables an OpenShift cluster administrator to specify which group of pods can communicate with each other. The grouping of pods is achieved via the use of labels while the network policy is specified in terms of a YAML file. 

The pods and network policies can be created in any order and are independent of each other. When a network policy is created, it is instantly applied to all the pods with the matching labels. Subsequently, when new pods with labels get spawned after the creation of the network policy, the policy is applied to the new pods with matching labels.

The labels attached to the pod are used to create Nuage Policy Groups into which the vports of the pods are added. While the policy rules mentioned in the policy YAML are implemented as ACLs on the VSD. 

The above concepts can be understood with the following two examples.

Example 1. Enable Intra-namespace Pod Communication
---------------------------------------------------

By default Nuage enables all the pods in a namespace to communicate with each other. However this functionality can be disabled using the VSD and adding appropriate ingress and egress ACLs. The following example assumes the blocking ACLs have been added using the VSD.

**Example json File - 1**

   ::

      hello-openshift-client.json 
      {
        "kind": "Pod",
        "apiVersion": "v1",
        "metadata": {
          "name": "hello-openshift-client",
          "creationTimestamp": null,
          "labels": {
            "name": "hello-openshift-client",
            "role": "client"
          }
        },
        "spec": {
          "containers": [
            {
              "name": "hello-openshift-client",
              "image": "openshift/hello-openshift",
              "ports": [
                {
                  "containerPort": 8080,
                  "protocol": "TCP"
                }
              ],
              "resources": {
              },
              "terminationMessagePath": "/dev/termination-log",
              "imagePullPolicy": "IfNotPresent",
              "capabilities": {},
              "securityContext": {
                "capabilities": {},
                "privileged": false
              }
            }
          ],
          "restartPolicy": "Always",
          "dnsPolicy": "ClusterFirst",
          "serviceAccount": ""
        },
        "status": {}
      }

**Example json File - 2**

   ::

      hello-openshift-server.json 
      {
        "kind": "Pod",
        "apiVersion": "v1",
        "metadata": {
          "name": "hello-openshift-server",
          "creationTimestamp": null,
          "labels": {
            "name": "hello-openshift-server",
            "role": "server"
          }
        },
        "spec": {
          "containers": [
            {
              "name": "hello-openshift-server",
              "image": "openshift/hello-openshift",
              "ports": [
                {
                  "containerPort": 8080,
                  "protocol": "TCP"
                }
              ],
              "resources": {
              },
              "terminationMessagePath": "/dev/termination-log",
              "imagePullPolicy": "IfNotPresent",
              "capabilities": {},
              "securityContext": {
                "capabilities": {},
                "privileged": false
              }
            }
          ],
          "restartPolicy": "Always",
          "dnsPolicy": "ClusterFirst",
          "serviceAccount": ""
        },
        "status": {}
      }
      oc create -f hello-openshift-client.json
      oc create -f hello-openshift-server.json

The above commands spawns sample hello-openshift client and server pods in the default namespace and if a user tries to wget from the client pod to the server pod, the wget can fail. However the communication can be enabled by specifying the following policy:
   
   ::  

      kind: NetworkPolicy
      apiVersion: extensions/v1beta1
      metadata:
        name: allow-tcp-8080
        Namespace: default
        labels:
          "nuage.io/priority": "100"
      spec:
        podSelector:
          matchLabels:
            role: server
        ingress:
          - ports:
            - protocol: TCP
              port: 8080    
            from:
              - podSelector:
                  matchLabels:
                    role: client

The network policy is created on Kubernetes using the following command:

   ::

      oc create -f network-policy.yaml

The above command creates a policy group for the client pods and another policy group for the server. Appropriate ACLs are added to enable the http communication between the two set of pods.

.. Note:: The priority for each of the network policy should be unique, as it maps to the ACL template priorities on the VSD. Also lower value of the priority implies a high priority in Nuage VSP. For more information on ACL priorities, see the "ACL Sandwich" section of the VSP User Guide.

Example 2. Enable Inter-namespace Pod Communication
----------------------------------------------------

   ::

      oc new-project server
      oc create -f hello-openshift-server.json

      oc new-project client
      oc create -f hello-openshift-client.json

      

By default Nuage blocks all of the communication between pods across namespaces. As a result a "wget" to the server pod in the server project from the client pod in the client project can fail. In order to enable the communication between the two set of pods a network policy can be defined in a YAML as shown below:

   ::

      ---
      kind: NetworkPolicy
      apiVersion: extensions/v1beta1
      metadata:
        name: allow-inter-ns-web-access
        namespace: server
        labels:
          "nuage.io/priority": "500"
      spec:
        podSelector:  
          matchLabels:  
            role: server
        ingress:
          - ports:
              - protocol: TCP
                port: 8080
            from:
              - podSelector:  
                  matchLabels:
                    role: client  

The network policy is created on OpenShift using the following command:

   ::

      oc create -f network-policy.yaml

The above command creates the policy-groups (corresponding to the pod labels) and addition of appropriate ACLs between them to facilitate the communication.

Example 3. Namespace Annotations 
--------------------------------

The Nuage OpenShift solution also supports namespace annotations, which allows an administrator to block all of the intra namespace communication by default when the isolation attribute for the namespace annotation is set to `DefaultDeny`. The namespace annotation can be enabled using the following command:

    ::

        oc annotate ns <namespace> "net.beta.kubernetes.io/network-policy={\"ingress\": {\"isolation\": \"DefaultDeny\"}}"

When the namespace is annotated, ingress and egress ACLs are added to the Nuage VSD to block all of the intra namespace communication. The namespace annotations for a namespace can be removed using the following command.

    ::

        oc annotate --overwrite ns <namespace> "net.beta.kubernetes.io/network-policy=-"

