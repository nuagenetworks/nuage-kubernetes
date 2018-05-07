#!/bin/bash

/usr/bin/nuage-bgp-startup.py

# This prevents Kubernetes/OpenShift from restarting the Nuage BGP 
# pod repeatedly.
should_sleep=${SLEEP:-"true"}
echo "Spawning Nuage BGP Pod.  Sleep=$should_sleep"
while [ "$should_sleep" == "true"  ]; do
        sleep 10;
done
