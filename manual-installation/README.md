# Installation On Kubernetes

## Master installation
### Install `nuagekubemon` on the Master
- Put `nuagekubemon` at `/usr/bin/nuagekubemon`
- Create the `/var/log/nuagekubemon` directory
- Create the `/usr/share/nuagekubemon` directory

### Install the Configuration Files
- Get the `clusterNetworkCIDR` from the `kube/default/kube-apiserver`'s
  `$KUBE_APISERVER_OPTS` `--service-cluster-ip-range` option
- Copy the included `net-config.yaml` and `nuagekubemon.yaml` files into
  `/usr/share/nuagekubemon/`
  - Change the value of `serviceNetworkCIDR` in `net-config.yaml` to whatever
    was in `--service-cluser-ip-range`
  - Change the value of 'vsdApiUrl' to point to your VSD
 
### Generate Necessary Certificates
`gencerts.sh` is the REST server certificates part of the serviceaccount.sh
script.  In order to generate the required certificates for communication
between the plugin(s) and the master's REST API server, run `gencerts.sh` as
follows:
```Shell
$ gencerts.sh --output-cert-dir=/usr/share/nuagekubemon/
```
 
### Disable Flannel and Start `nuagekubemon`
The default command structure is:
```Shell
$ systemctl disable flanneld
$ systemctl stop flanneld
$ nuagekubemon --config=/usr/share/nuagekubemon/nuagekubemon.yaml
```

## Node Installation
### Install Nuage k8s CNI plugin
- Put `nuage-cni-k8s` at
  `/usr/bin`
- Create the `/usr/share/vsp-k8s` directory
 
### Configure vsp-k8s on each node as follows:
- Copy the included `vsp-k8s.yaml` into `/usr/share/vsp-k8s/`
  - Change the values of `masterApiServer` and `nuageMonRestServer` to point to
    your k8s master
- Copy `nuageMonCA.crt`, `nuageMonClient.crt`, and `nuageMonClient.key` from the
  master so that the plugin can communicate with the REST server provided by
  nuagekubemon.
- On this setup, the api server is configured to do http only, but we still need
  to supply certs since we currently assume https.  Copy `/srv/kubernets/ca.crt`
  from the master to `/usr/share/vsp-k8s/ca.crt`, and copy `nuageMonClient.crt`
  and `nuageMonClient.key` to `client.crt` and `client.key` (respectively) to
  satisfy the plugin

### Configure nuage-openvswitch
Configure `/etc/default/openvswitch` to use PAT and set `PLATFORM="k8s"`
   
### Configure k8s to use the vsp-k8s plugin
Add `--cni-bin-dir=/usr/bin/` and
`--network-plugin=cni` to the end of `$KUBELET_OPTS` in
`/etc/default/kubelet`

### Set kube-proxy to userspace mode
Add `--proxy-mode=userspace` to `$KUBE_PROXY_OPTS` in `/etc/default/kube-proxy`

### Restart services
```Shell
$ systemctl disable flanneld ufw # or whatever firewall you actually have (firewalld or iptables for example)
$ systemctl stop flanneld
$ systemctl restart nuage-openvswitch kubelet kube-proxy
```

## Questions

1) Are we going to change the plugin rpmspec file to fix the location where the plugin gets installed ?

Yes, that is the plan.

2) In nuagekubemon.yaml what should be the value of kubeConfig ? How is the kubeConfig generated ?

A kubeconfig can be created using the `kubectl config` subcommand.  Details on how to use `kubectl config` can be found [here](http://kubernetes.io/docs/user-guide/kubeconfig-file/)
