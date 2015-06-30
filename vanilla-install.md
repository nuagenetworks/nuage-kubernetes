Openshift Install Guide
=========================

On all nodes:
-------------
```
export ORIGIN_MASTER1=openshift-origin-master
export ORIGIN_MINION1=openshift-origin-minion1
export ORIGIN_MINION2=openshift-origin-minion2
export MASTER1_IP=192.168.33.4
export MINION1_IP=192.168.33.5
export MINION2_IP=192.168.33.6
```

Install go from the tar-ball on golang.org/dl/

Setup for master
-----------------

1. Stop iptables

  ```
  systemctl stop iptables.service
  ```

2. Set CERT_DIR and SERVER_CONFIG_DIR appropriately and then create master certificates and node configs. 
  ```
  export CERT_DIR=$HOME/os-certs/master
  export SERVER_CONFIG_DIR=$HOME/os-certs
  
  openshift admin create-master-certs --overwrite=false --cert-dir=${CERT_DIR} --master=https://${MASTER1_IP}:8443 --hostnames=${MASTER1_IP},${ORIGIN_MASTER1}
  
  openshift admin create-node-config --node-dir="${SERVER_CONFIG_DIR}/node-${ORIGIN_MINION1}" --node=${ORIGIN_MINION1} --hostnames="${ORIGIN_MINION1},${MINION1_IP}" --master="https://${MASTER1_IP}:8443" --node-client-certificate-authority="${CERT_DIR}/ca.crt" --signer-cert="${CERT_DIR}/ca.crt" --signer-key="${CERT_DIR}/ca.key" --signer-serial="${CERT_DIR}/ca.serial.txt" --volume-dir="/root/openshift.local.volumes_${ORIGIN_MINION1}" --certificate-authority="${CERT_DIR}/ca.crt"
  
  openshift admin create-node-config --node-dir="${SERVER_CONFIG_DIR}/node-${ORIGIN_MINION2}" --node=${ORIGIN_MINION2} --hostnames="${ORIGIN_MINION2},${MINION2_IP}" --master="https://${MASTER1_IP}:8443" --node-client-certificate-authority="${CERT_DIR}/ca.crt" --signer-cert="${CERT_DIR}/ca.crt" --signer-key="${CERT_DIR}/ca.key" --signer-serial="${CERT_DIR}/ca.serial.txt" --volume-dir="/root/openshift.local.volumes_${ORIGIN_MINION2}" --certificate-authority="${CERT_DIR}/ca.crt"
  ```

3. Tar and scp all the os-certs files to the nodes. 
  ```
  tar cvzf os-certs.tar.gz os-certs
  ```
4. Install openshift-sdn (for using overlay networks for pod inter-networking)
  ```
  git clone https://github.com/openshift/openshift-sdn
  cd openshift-sdn
  make clean;
  make;
  make install;
  ```

5. Create a template configuration file. Note that the network-plugin is specified to be openshift-sdn. This will ensure that we will use openshift-sdn as the networking technology for inter-pod networking. The parameters of interest are clusterCIDR and networkPluginName in the config.yaml.

  ```
  openshift start master --master=https://${MASTER1_IP}:8443 --nodes=${ORIGIN_MINION1},${ORIGIN_MINION2} --create-certs=false --network-plugin=redhat/openshift-ovs-subnet --write-config='openshift.local.config'
  ```

6. Copy the master-config.yaml and policy.json files to the os-certs/master directory. This is a bit of a hack. You have to start the master first and let it create the policy.json and then kill it immediately.
    ```
    openshift start master --master=https://${MASTER1_IP}:8443 --nodes=${ORIGIN_MINION1},${ORIGIN_MINION2} --create-certs=false --network-plugin=redhat/openshift-ovs-subnet
    ```
  
    Then copy the policy.json and master-config.yaml:
    
    ```
    cp openshift.local.config/master-config.yaml /root/os-certs/master/.
    cp ~/openshift.local.config/policy.json /root/os-certs/master/.
    ```

7. Finally start the master using the config file in the new location. If the start complains on any files missing. Just copy them from openshift.local.config directory to the master directory.
  ```
  openshift start master --config=/root/os-certs/master/master-config.yaml &
  ```

8. Set KUBECONFIG to the appropriate admin.kubeconfig file in /root/os-certs/master/admin.kubeconfig. Alternatively also set it in bash_profile of the current user. The KUBECONFIG env variable is only used by osc CLI.

Setup on node
--------------

1. Stop iptables
  ```
  systemctl stop iptables.service
  ```

2. Untar the os-certs.tar.gz to the os-certs directory and set CERT_DIR appropriately.
  ```
  export CERT_DIR=$HOME/os-certs/
  ```

3. Make sure the requirements below are met

  #### Requirements:
  
  Need docker 1.6.2-8 on the minions/nodes for getting around the issue of not being able to create the /name/secrets correctly.
  Can be installed using the cbs centos yum repo:
  
  ###### Add the following to a new /etc/yum.repos.d/cbs.repo file
      ```
      [virt7-testing]
      name=virt7-testing
      baseurl=http://cbs.centos.org/repos/virt7-testing/x86_64/os/
      enabled=1
      gpgcheck=0
      ```
  
  Then install docker using yum
      ```
      yum install docker
      ```
  
  Check the version using docker version or rpm -q docker.

4. Make sure the following insecure-registry config option is turned ON in /etc/sysconfig/docker.
  ```
  INSECURE_REGISTRY='--insecure-registry 172.30.0.0/16'
  ```

5. Start docker service: 

  ```
  systemctl restart docker.service
  ```

6. Install openshift-sdn (for using overlay networks for pod inter-networking)

  ```
  git clone https://github.com/openshift/openshift-sdn
  cd openshift-sdn
  make clean;
  make;
  make install;
  ```

7. Start node
  
  ```
  openshift start node --config=/root/os-certs/node-openshift-origin-minion1/node-config.yaml
  ```

8. Set KUBECONFIG to the appropriate admin.kubeconfig file in /root/os-certs/master/admin.kubeconfig. Alternatively also set it in bash_profile of the current user. The KUBECONFIG env variable is only used by osc CLI.


Back on master
---------------

9. For launching applications from source, Openshift builds images and puts them in a repository that is created inside a private docker-registry, create this registry first. Registry should be created using the """default project"""

  ```
  chmod +r /root/os-certs/master/openshift-registry.kubeconfig
  oadm registry --create --credentials=/root/os-certs/master/openshift-registry.kubeconfig --config=/root/os-certs/master/admin.kubeconfig
  ```

10. Make sure the registry is up and ready

  ```
  oc describe service docker-registry --config=/root/os-certs/master/admin.kubeconfig
  ```

If "Endpoints" is listed as `<none>`, your registry hasn't started yet. You can run 
  ```
  oc get pods 
  ``` 

to see the registry pod and if there are any issues. Once the pod has started, the IP of the pod will be added to the docker-registry service list so that it's reachable from other places.

Follow the sample-app guide for deployment of your first app

https://github.com/openshift/origin/tree/master/examples/sample-app

##### Some tips

###### To clean up exited containers.
```
for i in `docker ps -a | cut -f1 -d ' '`; do docker rm -f $i; done
```

###### To clean up openshift entities
  ```
  for i in `osc get pods | cut -f1 -d ' '`; do osc delete pod $i; done
  for i in `osc get deploymentConfigs | cut -f1 -d ' '`; do osc delete deploymentConfigs $i; done
  for i in `osc get service | cut -f1 -d ' '`; do osc delete service $i; done
  for i in `osc get builds | cut -f1 -d ' '`; do osc delete builds $i; done
  for i in `osc get imagestreams | cut -f1 -d ' '`; do osc delete imagestreams $i; done
  for i in `osc get buildConfigs | cut -f1 -d ' '`; do osc delete buildConfigs $i; done
  for i in `osc get routes | cut -f1 -d ' '`; do osc delete routes $i; done
  ```

###### To check if the service forwarding is correct or not:
  ```
  sudo iptables -L -n -v -t nat
  oc get endpoints -l name='app-name-label' -o json
  oc get se -l name='app-name-label' -o json
  ```

###### General training material on openshift is available here:

https://github.com/openshift/training/blob/master/beta-4-setup.md

###### SSH Port Forwarding
To set up port forwarding once the app is up and running for the service IP (this is the magic that makes service_ip reachable from any outside host (strictly speaking we are cheating))
  ```
  ssh -L :8000:service_ip:service_port -N -f root@192.168.33.5 
  ```

Note that without a router for the service, request needs to be forwarded to one of the nodes. In this case we just chose 192.168.33.5. If you route it to any other remote host (even the master for that matter) which is not a node, since kube_proxy is not running there, it will not respond. In case you mess up and you need to re-configure the port forwarding, you will have to first delete the existing port forwarding process. In order to delete an existing ssh port forwarding, so that new one can be set, look for the process PID that is running the ssh tunnel
  ```
  fuser 8000/tcp
  ```

This should return a PID. Kill the process using kill -9 pid. And then re-enable port forwarding to a valid remote host.

###### Router creation and setup.

https://github.com/openshift/origin/blob/master/docs/routing.md

To verify if the route is attached to the router and the route resolves the service alias:
  
  ```
  curl -k --resolve www.aniketsrubysample.com:443:192.168.33.5 https://www.aniketsrubysample.com
  ```






