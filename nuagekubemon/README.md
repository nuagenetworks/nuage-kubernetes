NuageKubeMon
=====================

NuageKubeMon is a daemon that runs on the master node of a Openshift
installation and monitors the Openshift ReST API for creation of Openshift
entities that are interesting from network orchestration perspective. It is
also responsible for communicating with the VSD and providing VSD state
information to all the nodes of the Openshift cluster.


Installation Instructions
--------------------------

From source:
------------

1. First git clone the repository
2. Make sure GOPATH is set appropriately.
3. Install godep using `go get github.com\tools\godep`
4. Go to openshift-integration/nuagekubemon directory and run `godep restore`
5. Finally do `godep go build` and `godep go install`
6. Run nuagekubemon on the server passing the path to the kubeconfig file
   used for connecting to the kubernetes client (admin.kubeconfig) and the
   master config.yaml file. There are other configurable parameters associated
   with nuagekubemon that can be seen using --help option.

From binary:
------------

1. TBD after we can do go get.
