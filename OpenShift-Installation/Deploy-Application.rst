
===============================
Deploying Application Pods
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

The nuage-openshift-monitor creates a zone and subnet for each project in the OpenShift cluster. Whenever a pod is deployed in a project, it gets an IP from the corresponding subnet and a virtual port is created on the VSD.

1. Deploy the sample hello-openshift pod. The openshift-origin repository can be cloned to get the hello-pod.json.

   ::

        oc create -f /dir/origin/examples/hello-openshift/hello-pod.json
        
2. Deploy the docker registry service in the default project.

   ::

        oadm registry --config=/etc/origin/master/admin.kubeconfig
        
3. Deploy the `OpenShift 3 Application Lifecycle Sample <https://github.com/openshift/origin/tree/master/examples/sample-app>`_. 

   ::

        oc new-app -f origin/examples/sample-app/application-template-stibuild.json
        


   This creates a service, build config, replication controller, deployment config and the number of pods mentioned inthe json for the frontend and backend applications respectively as can be seen in the following output.

    .. Note:: The following output is too wide and may not display correctly in the PDF format.
    

    
   ::

        [root@master-ovs1 nuage-openshift-monitor]# oc get all
       NAME                        TYPE                                              FROM      LATEST
       ruby-sample-build           Source                                            Git       1
       NAME                        TYPE                                              FROM      STATUS     STARTED            DURATION
       ruby-sample-build-1         Source                                            Git       Complete   45 minutes ago   53s
       NAME                        DOCKER REPO                                       TAGS      UPDATED
       origin-ruby-sample          172.30.24.50:5000/sample-app/origin-ruby-sample   latest    44 minutes ago
       ruby-22-centos7             centos/ruby-22-centos7                            latest    45 minutes ago
       NAME                        TRIGGERS                                          LATEST
       database                    ConfigChange                                      1
       frontend                    ConfigChange, ImageChange                         1
       CONTROLLER                  CONTAINER(S)                                      IMAGE(S)                                                                                                                 SELECTOR                                                        REPLICAS                                                AGE
       database-1                  ruby-helloworld-database                          openshift/mysql-55-centos7:latest                                                                                           deployment=database-1,deploymentconfig=database,name=database   1                                                       45m
         frontend-1                  ruby-helloworld                                   172.30.24.50:5000/sample-app/origin-ruby-sample@sha256:2524e5f38d6a50e38ebb6ce0e5595669a3bb3d57c1b6c0f229b04cc581267ab3   deployment=frontend-1,deploymentconfig=frontend,name=frontend   2                                                       44m
       NAME                        HOST/PORT                                         PATH                                                                                                                      SERVICE                                                         LABELS                                                  INSECURE POLICY   TLS TERMINATION
       route-edge                  www.demo-test.com                                                                                                                                                             frontend                                                        handle=testing,template=application-template-stibuild                     edge
       NAME                        CLUSTER_IP                                        EXTERNAL_IP                                                                                                               PORT(S)                                                         SELECTOR                                                AGE
       database                    172.30.63.77                                      <none>                                                                                                                    5434/TCP                                                        name=database                                           45m
       frontend                    172.30.190.80                                     <none>                                                                                                                    5432/TCP                                                        name=frontend                                           45m
       NAME                        READY                                             STATUS                                                                                                                    RESTARTS                                                        AGE
       database-1-32g23            1/1                                               Running                                                                                                                   0                                                               44m
       frontend-1-s6isp            1/1                                               Running                                                                                                                   0                                                               44m
       frontend-1-uk171            1/1                                               Running                                                                                                                   0                                                               44m
       ruby-sample-build-1-build   0/1                                               Completed          
       

   It also creates a route to the frontend service. Each frontend and backend pod goes through the various lifecycles (prehook, posthook).

Operations Workflow
====================

1. Create a zone and subnet manually on the VSD for pods deployed in the Operations workflow. Deploy the sample hello-openshift pod by adding the Nuage specific labels in the "metadata" section of the json as shown below.

   ::

        {
          "kind": "Pod",
          "apiVersion": "v1",
          "metadata": {
            "name": "pod-wf2",
            "creationTimestamp": null,
            "labels": {
              "name": "pod-wf2",
              _**"nuage.io/subnet": "wf2_subnet",**_
              _**"nuage.io/zone": "wf2_zone",**_
              _**"nuage.io/user": "admin",**_
              _**"nuage.io/policy-group": "wf2_policy-group"**_
            }
          },
          "spec": {
            "containers": [
              {
                "name": "pod-wf2",
                "image": "openshift/hello-openshift",
                "ports": [
                  {
                    "containerPort": 8080,
                    "protocol": "TCP"
                  }
                ],
                "resources": {},
                "volumeMounts": [
                  {
                    "name":"tmp",
                    "mountPath":"/tmp"
                  }
                ],
                "terminationMessagePath": "/dev/termination-log",
                "imagePullPolicy": "IfNotPresent",
                "capabilities": {},
                "securityContext": {
                  "capabilities": {},
                  "privileged": false
                }
              }
            ],
            "volumes": [
              {
                "name":"tmp",
                "emptyDir": {}
              }
            ],
            "restartPolicy": "Always",
            "dnsPolicy": "ClusterFirst",
            "serviceAccount": ""
          },
          "status": {}
        }
        


   Once the pod is deployed, a vPort can be seen resolved in the manually created zone (wf2_zone) and subnet (wf2_subnet). With the policy-group label, the vPort also gets added to the policy group mentioned provided it is already created on the VSD before deployment.

2. A Network Macro Group can be manually created on the VSD for services deployed in the Operations workflow. Deploy the sample hello-openshift service by adding the Nuage specific labels in the "metadata" section of the json as shown below. The original json can be found in the openshift-origin repo.

   ::

        {
          "kind": "List",
          "apiVersion": "v1",
          "metadata": {
            "name": "hello-service-complete-example"
         },
          "items": [
            {
              "kind": "Service",
              "apiVersion": "v1",
              "metadata": {
                "name": "hello-openshift",
                "labels": {
                  "name": "hello-openshift",
                  _**"nuage.io/private-service": "false",**_
                  _**"nuage.io/network-macro-group.name": "workflow-2-manual"**_
                }
              },
              "spec": {
                "selector": {
                  "name": "hello-openshift"
                },
                "ports": [
                  {
                    "protocol": "TCP",
                    "port": 27017,
                    "targetPort": 8080
                  }
                ],
                "portalIP": "",
                "type": "ClusterIP",
                "sessionAffinity": "None"
              }
            },
            {
              "kind": "Route",
              "apiVersion": "v1",
              "metadata": {
                "name": "hello-openshift-route",
                "labels": {
                  "name": "hello-openshift"
                }
              },
              "spec": {
                "host": "hello-openshift.example.com",
                "to": {
                  "name": "hello-openshift-service"
                },
                "tls": {
                  "termination": "edge"
                }
              }
            },
            {
              "kind": "DeploymentConfig",
              "apiVersion": "v1",
              "metadata": {
                "name": "hello-openshift",
                "labels": {
                  "name": "hello-openshift"
                }
              },
              "spec": {
                "strategy": {
                  "type": "Recreate",
                  "resources": {}
                },
                "triggers": [
                  {
                    "type": "ConfigChange"
                  }
                ],
                "replicas": 3,
                "selector": {
                  "name": "hello-openshift"
                },
                "template": {
                  "metadata": {
                    "creationTimestamp": null,
                    "labels": {
                      "name": "hello-openshift"
                    }
                  },
                  "spec": {
                    "containers": [
                      {
                        "name": "hello-openshift",
                        "image": "openshift/hello-openshift:v1.0.6",
                        "ports": [
                          {
                            "name": "http",
                            "containerPort": 8080,
                            "protocol": "TCP"
                          }
                        ],
                        "resources": {
                           "limits": {
                            "cpu": "10m",
                            "memory": "16Mi"
                          }
                        },
                        "terminationMessagePath": "/dev/termination-log",
                        "imagePullPolicy": "IfNotPresent",
                        "capabilities": {},
                        "securityContext": {
                          "capabilities": {},
                          "privileged": false
                        },
                        "livenessProbe": {
                          "tcpSocket": {
                            "port": 8080
                          },
                          "timeoutSeconds": 1,
                          "initialDelaySeconds": 10
                        }
                      }
                    ],
                    "restartPolicy": "Always",
                    "dnsPolicy": "ClusterFirst",
                    "serviceAccount": ""
                  }
                }
              }
            }
          ]
        }


   Once the service is created, the network macro that is automatically created by nuage-openshift-monitor gets added to the network macro group (workflow-2-manual) mentioned in the json. The network macro group ID can also be used instead of the name in the json.



