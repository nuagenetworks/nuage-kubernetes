#!/bin/bash

########################
# include the magic
########################
. demo-magic.sh

# hide the evidence
clear

get_pods()
{
  pe " kubectl get pods"
}


get_pods_2()
{
  y=`kubectl get pods $(kubectl get pods | grep -v NAME | awk '{ print $1}') -o 'jsonpath={.items[*].status.phase}'`
   while [ "$y" != "Running Running Running" ]
   do
    sleep 5
    echo -e "${CYAN}All pods are still not in running state. Waiting..."
    y=`kubectl get pods $(kubectl get pods | grep -v NAME | awk '{ print $1}') -o 'jsonpath={.items[*].status.phase}'`
   done
    rm -rf pods
    x=`kubectl get pods | grep -v NAME | awk '{ print $1}'`
    for i in $x
    do
      echo "pod $i info" >> pods
      echo "status = $(kubectl get pod $i -o 'jsonpath={.status.phase}')" >> pods
      echo "name = $(kubectl get pod $i -o 'jsonpath={.metadata.name}')" >> pods
      echo "hostIP = $(kubectl get pod $i -o 'jsonpath={.status.hostIP}')" >> pods
      echo "podIP = $(kubectl get pod $i -o 'jsonpath={.status.podIP}')" >> pods

    done
    sleep 1
    cat pods
}

get_pods_3()
{
 y=`kubectl get pods $(kubectl get pods --namespace=guestbook | grep -v NAME | awk '{ print $1}') --namespace=guestbook -o 'jsonpath={.items[*].status.phase}'`
   while [ "$y" != "Running Running Running Running Running Running" ]
   do
    sleep 5
    echo -e "${CYAN}All pods are still not in running state. Waiting..."
    y=`kubectl get pods $(kubectl get pods --namespace=guestbook | grep -v NAME | awk '{ print $1}') --namespace=guestbook -o 'jsonpath={.items[*].status.phase}'`
   done
    rm -rf pods
    x=`kubectl get pods --namespace=guestbook | grep -v NAME | awk '{ print $1}'`
    for i in $x
    do
      echo "pod $i info" >> pods
      echo "status = $(kubectl get pod $i --namespace=guestbook -o 'jsonpath={.status.phase}')" >> pods
      echo "name = $(kubectl get pod $i --namespace=guestbook -o 'jsonpath={.metadata.name}')" >> pods
      echo "hostIP = $(kubectl get pod $i --namespace=guestbook -o 'jsonpath={.status.hostIP}')" >> pods
      echo "podIP = $(kubectl get pod $i --namespace=guestbook -o 'jsonpath={.status.podIP}')" >> pods

    done
    sleep 1
    cat pods
}

get_pods_4()
{
 y=`kubectl get pods $(kubectl get pods --namespace=demo | grep -v NAME | awk '{ print $1}') --namespace=demo -o 'jsonpath={.items[*].status.phase}'`
   while [ "$y" != "Running Running" ]
   do
    sleep 5
    echo -e "${CYAN}All pods are still not in running state. Waiting..."
    y=`kubectl get pods $(kubectl get pods --namespace=demo | grep -v NAME | awk '{ print $1}') --namespace=demo -o 'jsonpath={.items[*].status.phase}'`
   done
    rm -rf pods
    x=`kubectl get pods --namespace=demo | grep -v NAME | awk '{ print $1}'`
    for i in $x
    do
      echo "pod $i info" >> pods
      echo "status = $(kubectl get pod $i --namespace=demo -o 'jsonpath={.status.phase}')" >> pods
      echo "name = $(kubectl get pod $i --namespace=demo -o 'jsonpath={.metadata.name}')" >> pods
      echo "hostIP = $(kubectl get pod $i --namespace=demo -o 'jsonpath={.status.hostIP}')" >> pods
      echo "podIP = $(kubectl get pod $i --namespace=demo -o 'jsonpath={.status.podIP}')" >> pods

    done
    sleep 1
    cat pods
}

echo -e "${CYAN} Welcome to the Nuage Kubernetes Integration demo. This should be fun!"
echo -e "${CYAN} Let's start with a small cluster"
pe " kubectl get nodes "
echo -e "${CYAN} Let's see if we have any pods here"
get_pods
pe " kubectl run my-nginx --image=rstarmer/nginx-curl --replicas=3 --port=80"

get_pods_2

pe " kubectl get deployments"
pe " kubectl expose deployment $(kubectl get deployments | grep -v NAME | awk 'NR==1{ print $1}') --port=80 --type=NodePort"


#Add a get service
pe " kubectl get services"

echo -e "${CYAN} Let's see if the pods can reach the service IP"

pe " kubectl exec $(kubectl get pod $(kubectl get pods | grep -v NAME | awk 'NR==1{ print $1}') -o 'jsonpath={.metadata.name}') -- curl $(kubectl get services | grep -i $(kubectl get deployments | grep -v NAME | awk 'NR==1{ print $1}') | awk 'NR==1{print $2}') "

echo -e "${CYAN} Let's make sure pods can reach each other directly"
p " ################################"
p " Pod 1 --> Pod 2"
p " ################################"

pe " kubectl exec $(kubectl get pod $(kubectl get pods | grep -v NAME | awk 'NR==1{ print $1}') -o 'jsonpath={.metadata.name}') -- curl $(kubectl get pod $(kubectl get pods | grep -v NAME | awk 'NR==2{ print $1}') -o 'jsonpath={.status.podIP}')"

echo -e "${CYAN} Let's also deploy a guestbook app in another namespace"
pe " kubectl create namespace guestbook"
pe " kubectl create -f kubernetes/examples/guestbook/all-in-one/guestbook-all-in-one.yaml --namespace=guestbook"

get_pods_3
pe " kubectl get all --namespace=guestbook"


echo -e "${CYAN} Let's see how by default inter-namespace service and pod communication is disabled"

p " ###########################################"
p " nginx Pod 1 --> Guestbook Frontend Service"
p " ###########################################"

pe " kubectl exec $(kubectl get pod $(kubectl get pods | grep -v NAME | awk 'NR==1{ print $1}') -o 'jsonpath={.metadata.name}') -- curl --connect-timeout 5 $(kubectl --namespace=guestbook get services $(kubectl --namespace=guestbook get deployments | grep -v NAME | awk 'NR==1{ print $1}') -o jsonpath={.spec.clusterIP})"

p " ###########################################"
p " nginx Pod 1 --> One of the frontend pods"
p " ###########################################"

pe " kubectl exec $(kubectl get pod $(kubectl get pods | grep -v NAME | awk 'NR==1{ print $1}') -o 'jsonpath={.metadata.name}') -- curl --connect-timeout 5 $(kubectl get pods --namespace=guestbook $(kubectl get pods --namespace=guestbook -L frontend | grep -v NAME | awk 'NR==1{ print $1}') -o 'jsonpath={.status.podIP}')"


echo -e "${CYAN} Cleaning up pods and deployments"

pe " kubectl delete deployments $(kubectl get deployments | grep -v NAME | awk 'NR==1{ print $1}') "
#pe " kubectl delete services $(kubectl get deployments | grep -v NAME | awk 'NR==1{ print $1}') "
pe " kubectl delete service $(kubectl get services | grep -v NAME | awk 'NR==2{print $1}') "
pe " kubectl delete deployments frontend redis-master redis-slave --namespace=guestbook"
pe " kubectl delete services frontend redis-master redis-slave --namespace=guestbook"
#pe " kubectl delete namespace demo"
#Use case 3
: <<'END'

echo -e "${CYAN} Creating Zone, Subnet and corresponding ACLs on the Nuage VSD for Operational workflow"
p  " python nuage-vsdk.py"
echo -e "${CYAN} Creating pod which has nuage labels as shown"
pe " cat kubernetes/docs/user-guide/new-nginx-deployment.yaml"

pe " kubectl create -f kubernetes/docs/user-guide/new-nginx-deployment.yaml"

get_pods_2

pe " kubectl get deployments"
pe " kubectl create -f kubernetes/docs/user-guide/new-nginx-service.yaml"

#pe " kubectl expose deployment $(kubectl get deployments | grep -v NAME | awk 'NR==1{ print $1}') --port=80 --type=NodePort"

echo -e "${CYAN} Creating another nginx pod in a different Nuage Zone"

pe " cat kubernetes/docs/user-guide/pod.yaml "

#kubectl delete namespace demo
kubectl create namespace demo

pe " kubectl create -f kubernetes/docs/user-guide/pod.yaml --namespace=demo"
pe " kubectl create -f kubernetes/docs/user-guide/pod-no-pg.yaml --namespace=demo"

get_pods_4


echo -e "${CYAN} Let's see if the pods in Zone demo can reach the service IP in Zone tfd"

pe " kubectl exec --namespace=demo $(kubectl get pod --namespace=demo $(kubectl get pods --namespace=demo | grep -v NAME | awk 'NR==1{ print $1}') -o 'jsonpath={.metadata.name}') -- curl $(kubectl get services | grep -i $(kubectl get deployments | grep -v NAME | awk 'NR==1{ print $1}') | awk 'NR==1{print $2}') "


echo -e "${CYAN} Let's see if the pods in Zone demo can reach the service IP in Zone tfd"


echo -e "${CYAN} Cleaning up pods, services, namespaces and deployments"

pe " kubectl delete service  $(kubectl get services | grep -v NAME | awk 'NR==2{ print $1}') "
pe " kubectl delete deployments $(kubectl get deployments | grep -v NAME | awk 'NR==1{ print $1}') "
pe " kubectl delete pods nginx nginx-2 --namespace=demo"

END

pe " kubectl delete namespace guestbook "
