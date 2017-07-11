# Setup etcd

 1. Copy the contents of the etcd folder to all the etcd nodes
 2. Fill in the etcd cluster specific information in the etcd.env file
 3. For the first node, generate the token manually using the link in the etcd.env file
 4. For the 2nd & 3rd node in the cluster, set the CURRENT_NODE accordingly
 5. Once etcd.env is populated correctl, run the setup-etcd.sh script as shown below.
 6. Status of the etcd service can be checked using `service etcd status` command
 
 ```
 cp etcd_template.env etcd.env
 _Fill in node/cluster specific configuration_
 vim etcd.env
 chmod +x setup-etcd.sh
 ./setup-etcd.sh
```
# Setup the Load balancer for the Master

1. Install haproxy on the node acting as the Load balancer
2. Modify the haproxy.cfg to balance the 3 masters with the config shown in masters directory
3. Restart the haproxy service for the configuration to take effect

# Setup base for the first master

1. Copy the contents from the master directory to the first master
2. Fill in the contents in the master.env file
3. The Token field will be set to blank for the first master
4. Run init-master.sh script as shown below

```
cp master_template.env master.env
_Fill in node/cluster specific configuration_
vim master.env
chmod +x setup-master.sh
_only needed on first master server and on first configuration_
chmod +x init-master.sh
./init-master.sh
```

# Configure other master servers

1. Create /etc/kubernetes directory on 2nd & 3rd master
2. Copy the contents from the first master in /etc/kubernetes directory to this master in /etc/kubernetes directory
3. Modify the master.env with the proper CURRENT_NODE value
4. Set the Token value by getting the token from the first master using `kubeadm token list` command 
5. Run the setup-master.sh script as shown below

```
_copy needed files from first master_
_copy /etc/kubernetes/*_
vim master.env
chmod +x setup-master.sh
./setup-master.sh
```

# Configure the nodes

1. Copy the contents from the node directory on all the nodes
2. Copy the token from the master.env file and populate it in the worker.env file
3. Set the correct SERVICE_SUBNET CIDR & POD_SUBNET CIDR
4. Run the setup-worker.sh script as shown below

```
_Fill in node/cluster specific configuration_
vim worker.env
chmod +x setup-worker.sh
./setup-worker.sh

```
