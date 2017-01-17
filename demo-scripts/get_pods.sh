#!/bin/sh

get_pods() {

 y=`kubectl get pods $(kubectl get pods | grep -v NAME | awk '{ print $1}') -o 'jsonpath={.items[*].status.phase}'`
   while [ "$y" != "Running Running Running" ]
   do
    sleep 5
    y=`kubectl get pods $(kubectl get pods | grep -v NAME | awk '{ print $1}') -o 'jsonpath={.items[*].status.phase}'`
   done
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

 get_pods

