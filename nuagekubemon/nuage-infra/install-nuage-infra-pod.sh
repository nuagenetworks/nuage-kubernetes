#!/bin/sh

# This prevents Kubernetes from restarting the Nuage infra 
# pod repeatedly.
should_sleep=${SLEEP:-"true"}
echo "Spawning Nuage Infra Pod.  Sleep=$should_sleep"
while [ "$should_sleep" == "true"  ]; do
	sleep 10;
done
