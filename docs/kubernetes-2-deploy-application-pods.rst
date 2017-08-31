.. _kubernetes-6-deploy-application-pods:

.. include:: ../lib/doc-includes/VSDA-icons.inc

===============================
Deployment of Application Pods
===============================

.. contents::
   :local:
   :depth: 3
   

Overview
========
There are two workflows for deployment of application pods:

1. Developer workflow
2. Operations workflow 


Developer Workflow
===================

The nuagekubemon creates a zone and subnet for each namespace in the Kubernetes cluster. Whenever a pod is deployed in a namespace, it gets an IP from the corresponding subnet and a virtual port is created on the VSD.

1. Deploy the sample nginx pod using the command as shown below. 
   ::

        kubectl run my-nginx --image=nginx --replicas=2 --port=80
        

2. Deploy the `kubernetes Guestbook app as described in <http://thockin.github.io/kubernetes/v1.0/examples/guestbook-go/README.html>`_. 

    The following example shows how to build a simple multi-tier web application using Kubernetes and Docker. The application consists of a web front-end, Redis master for storage, and replicated set of Redis slaves, all for which you need to create Kubernetes replication controllers, pods, and services.

    
   ::


        [root@ovs-10 ansible]# kubectl get all --namespace=demo
        NAME                            CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
        frontend                        10.254.236.50   <none>        80/TCP     3d
        redis-master                    10.254.70.91    <none>        6379/TCP   3d
        redis-slave                     10.254.76.102   <none>        6379/TCP   3d
        NAME                            READY           STATUS        RESTARTS   AGE
        frontend-440143558-03kmv        1/1             Running       0          3d
        frontend-440143558-7ammk        1/1             Running       0          3d
        frontend-440143558-jbtye        1/1             Running       0          3d
        redis-master-2353460263-nugq1   1/1             Running       0          3d
        redis-slave-1691881626-cgde7    1/1             Running       0          3d
        redis-slave-1691881626-opcf8    1/1             Running       0          3d

Operations Workflow
====================

1. Create a zone and subnet manually on the VSD for pods deployed in the Operations workflow. Deploy the sample  nginx pod by adding the Nuage specific labels in the "metadata" section of the yaml as shown below.

   ::

      kind: ReplicationController
      metadata:
         name: nginx
      spec:
         replicas: 2
         selector:
            app: nginx
         template:
            metadata:
               name: nginx
               labels:
                  app: nginx
                  nuage.io/subnet: operations-wf
                  nuage.io/zone: operations-wf
                  nuage.io/user: admin
                  nuage.io/policy-group: test

            spec:
               containers:
               - name: nginx
                 image: nginx
                 ports:
                 - containerPort: 80


   Once the pod is deployed, a vPort can be seen resolved in the manually created zone (operations-wf) and subnet (operations-wf). With the policy-group label, the vPort also gets added to the policy group mentioned provided it is already created on the VSD before deployment.
   
2. A Network Macro Group can be manually created on the VSD for services deployed in the Operations workflow. Deploy the nginx service by adding the Nuage specific labels in the "metadata" section of the yaml as shown below. 

   ::

      apiVersion: v1
      kind: Service
      metadata:
         name: nginx-service
         labels:
            nuage.io/network-macro-group.name: Operations_WF_Macro_Group
            "nuage.io/private-service": "false"
      spec:
         ports:
         - port: 8000 # the port that this service should serve on
           # the container on each pod to connect to, can be a name
           # (e.g. 'www') or a number (e.g. 80)
           targetPort: 80
           protocol: TCP
         # just like the selector in the replication controller,
         # but this time it identifies the set of pods to load balance
         # traffic to.
         selector:
            app: nginx

Once the service is created, the network macro is automatically created by nuage-openshift-monitor and gets added to the network macro group (Operations_WF_Macro_Group) mentioned in the yaml. The network macro group ID can also be used instead of the name in the yaml.


