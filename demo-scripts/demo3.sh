nclude the magic
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
    echo -e "{CYAN}All pods are still not in running state. Waiting..."
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
    echo -e "${CYAN} All pods are still not in running state. Waiting..."
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

echo -e "${CYAN} Welcome again to the Nuage Kubernetes Integration demo"
echo -e "${CYAN} Let's start with a same cluster for Operations workflow"
echo -e "${CYAN} Creating Zone, Subnet and corresponding ACLs on the Nuage VSD for Operational workflow"
p  " python nuage-vsdk.py"
echo -e "${CYAN} Creating pod which has nuage labels as shown"
pe " cat kubernetes/docs/user-guide/new-nginx-deployment.yaml"

pe " kubectl create -f kubernetes/docs/user-guide/new-nginx-deployment.yaml"

get_pods_2

pe " kubectl get deployments"
pe " kubectl create -f kubernetes/docs/user-guide/new-nginx-service.yaml"


echo -e "${CYAN} Creating another nginx pod in a different Nuage Zone"

pe " cat kubernetes/docs/user-guide/pod.yaml "

kubectl create namespace demo

pe " kubectl create -f kubernetes/docs/user-guide/pod.yaml --namespace=demo"
pe " kubectl create -f kubernetes/docs/user-guide/pod-no-pg.yaml --namespace=demo"

get_pods_4

pe " kubectl get services  "

echo -e "${CYAN} Let's see if the pod placed in Policy group demo can reach the service IP in Zone tfd"

pe " kubectl exec --namespace=demo $(kubectl get pod --namespace=demo $(kubectl get pods --namespace=demo | grep -v NAME | awk 'NR==1{ print $1}') -o 'jsonpath={.metadata.name}') -- curl $(kubectl get services | grep -i $(kubectl get deployments | grep -v NAME | awk 'NR==1{ print $1}') | awk 'NR==1{print $2}') "

p " ##################################################"
p " Pod 2 (not in Policy Group) --> Service IP"
p " ##################################################"

pe " kubectl exec --namespace=demo $(kubectl get pod --namespace=demo $(kubectl get pods --namespace=demo | grep -v NAME | awk 'NR==2{ print $1}') -o 'jsonpath={.metadata.name}') -- curl --connect-timeout 5 $(kubectl get services | grep -i $(kubectl get deployments | grep -v NAME | awk 'NR==1{ print $1}') | awk 'NR==1{print $2}') "


echo -e "${CYAN} Cleaning up pods and deployments"


pe " kubectl delete service  $(kubectl get services | grep -v NAME | awk 'NR==2{ print $1}') "
pe " kubectl delete deployments $(kubectl get deployments | grep -v NAME | awk 'NR==1{ print $1}') "
pe " kubectl delete pods nginx nginx-2 --namespace=demo"
pe " kubectl delete namespace demo "

