/*
###########################################################################
#
#   Filename:           nuageClusterclient.go
#
#   Author:             Aniket Bhat
#   Created:            July 20, 2015
#
#   Description:        Nuage VSP Cluster Client Interface
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
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/config"
	oscache "github.com/openshift/origin/pkg/client/cache"
	kapi "k8s.io/kubernetes/pkg/api"
	kextensions "k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/cache"
	krestclient "k8s.io/kubernetes/pkg/client/restclient"
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

type NuageClusterClient struct {
	kubeConfig *krestclient.Config
	kubeClient *kclient.Client
}

func NewNuageOsClient(nkmConfig *config.NuageKubeMonConfig) *NuageClusterClient {
	nosc := new(NuageClusterClient)
	nosc.Init(nkmConfig)
	return nosc
}

func (nosc *NuageClusterClient) GetClusterClientCallBacks() *api.ClusterClientCallBacks {
	return &api.ClusterClientCallBacks{
		FilterPods: nosc.GetPods,
		GetPod:     nosc.GetPod,
	}
}

func (nosc *NuageClusterClient) Init(nkmConfig *config.NuageKubeMonConfig) {
	loadingRules := &clientcmd.ClientConfigLoadingRules{}
	loadingRules.ExplicitPath = nkmConfig.KubeConfigFile
	loader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	kubeConfig, err := loader.ClientConfig()
	if err != nil {
		glog.Errorf("Got an error: %s while loading the kube config", err)
	}
	// This is an internal client which is shared by most controllers, so boost default QPS
	// TODO: this should be configured by the caller, not in this method.
	if kubeConfig != nil {
		kubeConfig.QPS = 100.0
		kubeConfig.Burst = 200
		kubeConfig.WrapTransport = DefaultClientTransport
	}
	nosc.kubeConfig = kubeConfig
	kubeClient, err := kclient.New(nosc.kubeConfig)
	if err != nil {
		glog.Errorf("Got an error: %s while creating the kube client", err)
	}
	nosc.kubeClient = kubeClient
}

func (nosc *NuageClusterClient) GetExistingEvents(nsChannel chan *api.NamespaceEvent, serviceChannel chan *api.ServiceEvent, policyEventChannel chan *api.NetworkPolicyEvent) {
	//we will use the kube client APIs than interfacing with the REST API
	listOpts := kapi.ListOptions{LabelSelector: labels.Everything(), FieldSelector: fields.Everything()}
	nsList, err := nosc.GetNamespaces(&listOpts)
	if err != nil {
		glog.Infof("Got an error: %s while getting namespaces list from kube client", err)
		return
	}
	for _, ns := range *nsList {
		nsChannel <- ns
	}
	//we will use the kube client APIs than interfacing with the REST API
	serviceList, err := nosc.GetServices(&listOpts)
	if err != nil {
		glog.Infof("Got an error: %s while getting services list from kube client", err)
		return
	}
	for _, service := range *serviceList {
		serviceChannel <- service
	}
	//get policies
	policiesList, err := nosc.GetNetworkPolicies(&listOpts)
	if err != nil {
		glog.Infof("Got an error: %s while getting network policies list from kube client", err)
	}
	for _, policy := range *policiesList {
		policyEventChannel <- policy
	}
}

func (nosc *NuageClusterClient) RunNetworkPolicyWatcher(policyChannel chan *api.NetworkPolicyEvent, stop chan bool) {
	nosc.WatchNetworkPolicies(policyChannel, stop)
}

func (nosc *NuageClusterClient) RunNamespaceWatcher(nsChannel chan *api.NamespaceEvent, stop chan bool) {
	nosc.WatchNamespaces(nsChannel, stop)
}

func (nosc *NuageClusterClient) RunServiceWatcher(serviceChannel chan *api.ServiceEvent, stop chan bool) {
	nosc.WatchServices(serviceChannel, stop)
}

func (nosc *NuageClusterClient) GetNamespaces(listOpts *kapi.ListOptions) (*[]*api.NamespaceEvent, error) {
	namespaces, err := nosc.kubeClient.Namespaces().List(*listOpts)
	if err != nil {
		return nil, err
	}
	namespaceList := make([]*api.NamespaceEvent, 0)
	for _, obj := range namespaces.Items {
		namespaceList = append(namespaceList, &api.NamespaceEvent{Type: api.Added, Name: obj.ObjectMeta.Name, Labels: obj.ObjectMeta.Labels, Annotations: obj.GetAnnotations()})
	}
	return &namespaceList, nil
}

func (nosc *NuageClusterClient) WatchNamespaces(receiver chan *api.NamespaceEvent, stop chan bool) error {
	nsEventQueue := oscache.NewEventQueue(cache.MetaNamespaceKeyFunc)
	listWatch := &cache.ListWatch{
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return nosc.kubeClient.Namespaces().List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return nosc.kubeClient.Namespaces().Watch(options)
		},
	}
	cache.NewReflector(listWatch, &kapi.Namespace{}, nsEventQueue, 0).Run()
	for {
		evt, obj, err := nsEventQueue.Pop()
		if err != nil {
			return err
		}

		eventType := watch.EventType(evt)
		ns := obj.(*kapi.Namespace)

		switch eventType {
		case watch.Added:
			receiver <- &api.NamespaceEvent{Type: api.Added, Name: ns.ObjectMeta.Name, Labels: ns.ObjectMeta.Labels, Annotations: ns.GetAnnotations()}
		case watch.Modified:
			receiver <- &api.NamespaceEvent{Type: api.Modified, Name: ns.ObjectMeta.Name, Labels: ns.ObjectMeta.Labels, Annotations: ns.GetAnnotations()}
		case watch.Deleted:
			receiver <- &api.NamespaceEvent{Type: api.Deleted, Name: ns.ObjectMeta.Name, Labels: ns.ObjectMeta.Labels, Annotations: ns.GetAnnotations()}
		}
	}
}

func (nosc *NuageClusterClient) GetServices(listOpts *kapi.ListOptions) (*[]*api.ServiceEvent, error) {
	services, err := nosc.kubeClient.Services(kapi.NamespaceAll).List(*listOpts)
	if err != nil {
		return nil, err
	}
	servicesList := make([]*api.ServiceEvent, 0)
	for _, service := range services.Items {
		labels := GetNuageLabels(&service)
		if label, exists := labels["private-service"]; !exists || strings.ToLower(label) == "false" {
			servicesList = append(servicesList, &api.ServiceEvent{Type: api.Added, Name: service.ObjectMeta.Name, ClusterIP: service.Spec.ClusterIP, Namespace: service.ObjectMeta.Namespace, NuageLabels: labels})
		}
	}
	return &servicesList, nil
}

func (nosc *NuageClusterClient) WatchServices(receiver chan *api.ServiceEvent, stop chan bool) error {
	serviceEventQueue := oscache.NewEventQueue(cache.MetaNamespaceKeyFunc)
	listWatch := &cache.ListWatch{
		ListFunc: func(rv kapi.ListOptions) (runtime.Object, error) {
			return nosc.kubeClient.Services(kapi.NamespaceAll).List(kapi.ListOptions{LabelSelector: labels.Everything(), FieldSelector: fields.Everything()})
		},
		WatchFunc: func(rv kapi.ListOptions) (watch.Interface, error) {
			return nosc.kubeClient.Services(kapi.NamespaceAll).Watch(kapi.ListOptions{LabelSelector: labels.Everything(), FieldSelector: fields.Everything()})
		},
	}
	cache.NewReflector(listWatch, &kapi.Service{}, serviceEventQueue, 0).Run()
	for {
		evt, obj, err := serviceEventQueue.Pop()
		if err != nil {
			return err
		}
		eventType := watch.EventType(evt)
		switch eventType {
		case watch.Added:
			fallthrough
		case watch.Deleted:
			service := obj.(*kapi.Service)
			labels := GetNuageLabels(service)
			if label, exists := labels["private-service"]; !exists || strings.ToLower(label) == "false" {
				receiver <- &api.ServiceEvent{Type: api.EventType(eventType), Name: service.ObjectMeta.Name, ClusterIP: service.Spec.ClusterIP, Namespace: service.ObjectMeta.Namespace, NuageLabels: labels}
			}
		}
	}
}

func (nosc *NuageClusterClient) GetPod(name string, ns string) (*api.PodEvent, error) {
	if ns == "" {
		ns = kapi.NamespaceAll
	}
	pod, err := nosc.kubeClient.Pods(ns).Get(name)
	if err != nil {
		return nil, err
	}
	return &api.PodEvent{Type: api.Added, Name: pod.Name, Namespace: pod.Namespace, Labels: pod.Labels}, nil
}

func (nosc *NuageClusterClient) GetPods(listOpts *kapi.ListOptions, ns string) (*[]*api.PodEvent, error) {
	if ns == "" {
		ns = kapi.NamespaceAll
	}
	pods, err := nosc.kubeClient.Pods(ns).List(*listOpts)
	if err != nil {
		return nil, err
	}
	podsList := make([]*api.PodEvent, 0)
	for _, pod := range pods.Items {
		podsList = append(podsList, &api.PodEvent{Type: api.Added, Name: pod.Name, Namespace: pod.Namespace, Labels: pod.Labels})

	}
	return &podsList, nil
}

func (nosc *NuageClusterClient) WatchPods(receiver chan *api.PodEvent, stop chan bool) error {
	podEventQueue := oscache.NewEventQueue(cache.MetaNamespaceKeyFunc)
	listWatch := &cache.ListWatch{
		ListFunc: func(rv kapi.ListOptions) (runtime.Object, error) {
			return nosc.kubeClient.Pods(kapi.NamespaceAll).List(kapi.ListOptions{LabelSelector: labels.Everything(), FieldSelector: fields.Everything()})
		},
		WatchFunc: func(rv kapi.ListOptions) (watch.Interface, error) {
			return nosc.kubeClient.Pods(kapi.NamespaceAll).Watch(kapi.ListOptions{LabelSelector: labels.Everything(), FieldSelector: fields.Everything()})
		},
	}
	cache.NewReflector(listWatch, &kapi.Pod{}, podEventQueue, 0).Run()
	for {
		evt, obj, err := podEventQueue.Pop()
		if err != nil {
			return err
		}
		eventType := watch.EventType(evt)
		switch eventType {
		case watch.Added:
			fallthrough
		case watch.Deleted:
			pod := obj.(*kapi.Pod)
			receiver <- &api.PodEvent{Type: api.EventType(eventType), Name: pod.Name, Namespace: pod.Namespace, Labels: pod.Labels}
		}
	}
}

func (nosc *NuageClusterClient) GetNetworkPolicies(listOpts *kapi.ListOptions) (*[]*api.NetworkPolicyEvent, error) {
	policies, err := nosc.kubeClient.NetworkPolicies(kapi.NamespaceAll).List(*listOpts)
	if err != nil {
		return nil, err
	}
	policiesList := make([]*api.NetworkPolicyEvent, 0)
	for _, policy := range policies.Items {
		policiesList = append(policiesList, &api.NetworkPolicyEvent{Type: api.Added, Name: policy.Name, Namespace: policy.Namespace, Policy: policy.Spec, Labels: policy.Labels})

	}
	return &policiesList, nil
}

func (nosc *NuageClusterClient) WatchNetworkPolicies(receiver chan *api.NetworkPolicyEvent, stop chan bool) error {
	policyEventQueue := oscache.NewEventQueue(cache.MetaNamespaceKeyFunc)
	listWatch := &cache.ListWatch{
		ListFunc: func(rv kapi.ListOptions) (runtime.Object, error) {
			return nosc.kubeClient.NetworkPolicies(kapi.NamespaceAll).List(kapi.ListOptions{LabelSelector: labels.Everything(), FieldSelector: fields.Everything()})
		},
		WatchFunc: func(rv kapi.ListOptions) (watch.Interface, error) {
			return nosc.kubeClient.NetworkPolicies(kapi.NamespaceAll).Watch(kapi.ListOptions{LabelSelector: labels.Everything(), FieldSelector: fields.Everything()})
		},
	}
	cache.NewReflector(listWatch, &kextensions.NetworkPolicy{}, policyEventQueue, 0).Run()
	for {
		evt, obj, err := policyEventQueue.Pop()
		if err != nil {
			return err
		}
		eventType := watch.EventType(evt)
		switch eventType {
		case watch.Added:
			fallthrough
		case watch.Deleted:
			policy := obj.(*kextensions.NetworkPolicy)
			receiver <- &api.NetworkPolicyEvent{Type: api.EventType(eventType), Name: policy.Name, Namespace: policy.Namespace, Policy: policy.Spec, Labels: policy.Labels}
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

func GetNuageLabels(input *kapi.Service) map[string]string {
	labels := input.Labels
	nuageLabels := make(map[string]string)
	for k, v := range labels {
		if strings.HasPrefix(k, "nuage.io") {
			tokens := strings.Split(k, "/")
			nuageLabels[tokens[1]] = v
		}
	}
	return nuageLabels
}
