/*
###########################################################################
#
#   Filename:           nuageosclient.go
#
#   Author:             Aniket Bhat
#   Created:            July 20, 2015
#
#   Description:        NuageKubeMon Openshift Client Interface
#
###########################################################################
#
#              Copyright (c) 2015 Nuage Networks
#
###########################################################################

*/

package client

import (
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	kclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client/cache"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client/clientcmd"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
	"github.com/golang/glog"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/api"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/config"
	oscache "github.com/openshift/origin/pkg/client/cache"
	"net"
	"net/http"
	"time"
)

type NuageOsClient struct {
	kubeConfig *kclient.Config
	kubeClient *kclient.Client
}

func NewNuageOsClient(nkmConfig *config.NuageKubeMonConfig) *NuageOsClient {
	nosc := new(NuageOsClient)
	nosc.Init(nkmConfig)
	return nosc
}

func (nosc *NuageOsClient) Init(nkmConfig *config.NuageKubeMonConfig) {
	loadingRules := &clientcmd.ClientConfigLoadingRules{}
	loadingRules.ExplicitPath = nkmConfig.KubeConfigFile
	loader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	kubeConfig, err := loader.ClientConfig()
	if err != nil {
		glog.Info("Got an error %s while loading the kube config", err)
	}
	// This is an internal client which is shared by most controllers, so boost default QPS
	// TODO: this should be configured by the caller, not in this method.
	kubeConfig.QPS = 100.0
	kubeConfig.Burst = 200
	kubeConfig.WrapTransport = DefaultClientTransport
	nosc.kubeConfig = kubeConfig
	kubeClient, err := kclient.New(nosc.kubeConfig)
	if err != nil {
		glog.Info("Got an error %s while creating the kube client", err)
	}
	nosc.kubeClient = kubeClient
}

func (nosc *NuageOsClient) Run(nsChannel chan *api.NamespaceEvent, stop chan bool) {
	//we will use the kube client APIs than interfacing with the REST API
	nsList, err := nosc.GetNamespaces()
	if err != nil {
		return
	}
	for _, ns := range *nsList {
		nsChannel <- ns
	}
	nosc.WatchNamespaces(nsChannel, stop)
}

func (nosc *NuageOsClient) GetNamespaces() (*[]*api.NamespaceEvent, error) {
	namespaces, err := nosc.kubeClient.Namespaces().List(labels.Everything(), fields.Everything())
	if err != nil {
		return nil, err
	}
	namespaceList := make([]*api.NamespaceEvent, 0)
	for _, obj := range namespaces.Items {
		namespaceList = append(namespaceList, &api.NamespaceEvent{Type: api.Added, Name: obj.ObjectMeta.Name})
	}
	return &namespaceList, nil
}

func (nosc *NuageOsClient) WatchNamespaces(receiver chan *api.NamespaceEvent, stop chan bool) error {
	nsEventQueue := oscache.NewEventQueue(cache.MetaNamespaceKeyFunc)
	listWatch := &cache.ListWatch{
		ListFunc: func() (runtime.Object, error) {
			return nosc.kubeClient.Namespaces().List(labels.Everything(), fields.Everything())
		},
		WatchFunc: func(resourceVersion string) (watch.Interface, error) {
			return nosc.kubeClient.Namespaces().Watch(labels.Everything(), fields.Everything(), resourceVersion)
		},
	}
	cache.NewReflector(listWatch, &kapi.Namespace{}, nsEventQueue, 4*time.Second).Run()
	for {
		eventType, obj, err := nsEventQueue.Pop()
		if err != nil {
			return err
		}
		switch eventType {
		case watch.Added:
			// we should ignore the modified event because status updates cause unnecessary noise
			// the only time we would care about modified would be if the minion changes its IP address
			// and hence all nodes need to update their vtep entries for the respective subnet
			// create minionEvent
			ns := obj.(*kapi.Namespace)
			receiver <- &api.NamespaceEvent{Type: api.Added, Name: ns.ObjectMeta.Name}
		case watch.Deleted:
			// TODO: There is a chance that a Delete event will not get triggered.
			// Need to use a periodic sync loop that lists and compares.
			ns := obj.(*kapi.Namespace)
			receiver <- &api.NamespaceEvent{Type: api.Deleted, Name: ns.ObjectMeta.Name}
		}
	}
	return nil
}

// DefaultClientTransport sets defaults for a client Transport that are suitable
// for use by infrastructure components.
func DefaultClientTransport(rt http.RoundTripper) http.RoundTripper {
	transport := rt.(*http.Transport)
	// TODO: this should be configured by the caller, not in this method.
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	transport.Dial = dialer.Dial
	// Hold open more internal idle connections
	// TODO: this should be configured by the caller, not in this method.
	transport.MaxIdleConnsPerHost = 100
	return transport
}
