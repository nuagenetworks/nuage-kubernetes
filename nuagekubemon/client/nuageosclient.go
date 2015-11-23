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
	"github.com/golang/glog"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/api"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/config"
	oscache "github.com/openshift/origin/pkg/client/cache"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	kclient "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/watch"
	"net"
	"net/http"
	"strings"
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
		glog.Infof("Got an error: %s while loading the kube config", err)
	}
	// This is an internal client which is shared by most controllers, so boost default QPS
	// TODO: this should be configured by the caller, not in this method.
	kubeConfig.QPS = 100.0
	kubeConfig.Burst = 200
	kubeConfig.WrapTransport = DefaultClientTransport
	nosc.kubeConfig = kubeConfig
	kubeClient, err := kclient.New(nosc.kubeConfig)
	if err != nil {
		glog.Infof("Got an error: %s while creating the kube client", err)
	}
	nosc.kubeClient = kubeClient
}

func (nosc *NuageOsClient) RunNamespaceWatcher(nsChannel chan *api.NamespaceEvent, stop chan bool) {
	//we will use the kube client APIs than interfacing with the REST API
	nsList, err := nosc.GetNamespaces()
	if err != nil {
		glog.Infof("Got an error: %s while getting namespaces list from kube client", err)
		return
	}
	for _, ns := range *nsList {
		nsChannel <- ns
	}
	nosc.WatchNamespaces(nsChannel, stop)
}

func (nosc *NuageOsClient) RunServiceWatcher(serviceChannel chan *api.ServiceEvent, stop chan bool) {
	//we will use the kube client APIs than interfacing with the REST API
	serviceList, err := nosc.GetServices()
	if err != nil {
		glog.Infof("Got an error: %s while getting services list from kube client", err)
		return
	}
	for _, service := range *serviceList {
		serviceChannel <- service
	}
	nosc.WatchServices(serviceChannel, stop)
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
	cache.NewReflector(listWatch, &kapi.Namespace{}, nsEventQueue, 0).Run()
	for {
		eventType, obj, err := nsEventQueue.Pop()
		if err != nil {
			return err
		}
		switch eventType {
		case watch.Added:
			fallthrough
		case watch.Deleted:
			ns := obj.(*kapi.Namespace)
			receiver <- &api.NamespaceEvent{Type: api.EventType(eventType), Name: ns.ObjectMeta.Name}
		}
	}
}

func (nosc *NuageOsClient) GetServices() (*[]*api.ServiceEvent, error) {
	services, err := nosc.kubeClient.Services(kapi.NamespaceAll).List(labels.Everything())
	if err != nil {
		return nil, err
	}
	servicesList := make([]*api.ServiceEvent, 0)
	for _, service := range services.Items {
		annotations := GetNuageAnnotations(&service)
		if annotation, exists := annotations["private-service"]; !exists || strings.ToLower(annotation) == "false" {

			servicesList = append(servicesList, &api.ServiceEvent{Type: api.Added, Name: service.ObjectMeta.Name, ClusterIP: service.Spec.ClusterIP, Namespace: service.ObjectMeta.Namespace, NuageAnnotations: annotations})
		}
	}
	return &servicesList, nil
}

func (nosc *NuageOsClient) WatchServices(receiver chan *api.ServiceEvent, stop chan bool) error {
	serviceEventQueue := oscache.NewEventQueue(cache.MetaNamespaceKeyFunc)
	listWatch := &cache.ListWatch{
		ListFunc: func() (runtime.Object, error) {
			return nosc.kubeClient.Services(kapi.NamespaceAll).List(labels.Everything())
		},
		WatchFunc: func(resourceVersion string) (watch.Interface, error) {
			return nosc.kubeClient.Services(kapi.NamespaceAll).Watch(labels.Everything(), fields.Everything(), resourceVersion)
		},
	}
	cache.NewReflector(listWatch, &kapi.Service{}, serviceEventQueue, 0).Run()
	for {
		eventType, obj, err := serviceEventQueue.Pop()
		if err != nil {
			return err
		}
		switch eventType {
		case watch.Added:
			fallthrough
		case watch.Deleted:
			service := obj.(*kapi.Service)
			annotations := GetNuageAnnotations(service)
			if annotation, exists := annotations["private-service"]; !exists || strings.ToLower(annotation) == "false" {
				receiver <- &api.ServiceEvent{Type: api.EventType(eventType), Name: service.ObjectMeta.Name, ClusterIP: service.Spec.ClusterIP, Namespace: service.ObjectMeta.Namespace, NuageAnnotations: annotations}
			}
		}
	}
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

func GetNuageAnnotations(input *kapi.Service) map[string]string {
	annotations := input.Annotations
	nuageAnnotations := make(map[string]string)
	for k, v := range annotations {
		if strings.HasPrefix(k, "nuage.io") {
			tokens := strings.Split(k, "/")
			nuageAnnotations[tokens[1]] = v
		}
	}
	return nuageAnnotations
}
