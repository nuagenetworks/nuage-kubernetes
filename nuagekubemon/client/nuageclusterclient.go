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
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/config"
	v1 "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	krestclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type NuageClusterClient struct {
	kubeConfig *krestclient.Config
	clientset  *kubernetes.Clientset
}

func NewNuageOsClient(nkmConfig *config.NuageKubeMonConfig) *NuageClusterClient {
	nosc := new(NuageClusterClient)
	nosc.Init(nkmConfig)
	return nosc
}

func (nosc *NuageClusterClient) GetClusterClientCallBacks() *api.ClusterClientCallBacks {
	return &api.ClusterClientCallBacks{
		FilterPods:       nosc.GetPods,
		FilterNamespaces: nosc.GetNamespaces,
		GetPod:           nosc.GetPod,
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
	//contain clients to various api groups including rest client
	clientset, err := kubernetes.NewForConfig(nosc.kubeConfig)
	if err != nil {
		glog.Errorf("Creating new clientset from kubeconfig failed with error: %v", err)
		return
	}

	nosc.clientset = clientset
}

func (nosc *NuageClusterClient) GetExistingEvents(nsChannel chan *api.NamespaceEvent, serviceChannel chan *api.ServiceEvent, policyEventChannel chan *api.NetworkPolicyEvent) {
	//we will use the kube client APIs than interfacing with the REST API
	listOpts := metav1.ListOptions{LabelSelector: labels.Everything().String(), FieldSelector: fields.Everything().String()}
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

func (nosc *NuageClusterClient) RunPodWatcher(podChannel chan *api.PodEvent, stop chan bool) {
	nosc.WatchPods(podChannel, stop)
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

func (nosc *NuageClusterClient) GetNamespaces(listOpts *metav1.ListOptions) (*[]*api.NamespaceEvent, error) {
	namespaces, err := nosc.clientset.CoreV1().Namespaces().List(*listOpts)
	if err != nil {
		return nil, err
	}
	namespaceList := make([]*api.NamespaceEvent, 0)
	for _, obj := range namespaces.Items {
		namespaceList = append(namespaceList, &api.NamespaceEvent{Type: api.Added, Name: obj.ObjectMeta.Name, UID: string(obj.ObjectMeta.UID), Annotations: obj.GetAnnotations()})
	}
	return &namespaceList, nil
}

func (nosc *NuageClusterClient) WatchNamespaces(receiver chan *api.NamespaceEvent, stop chan bool) error {
	source := cache.NewListWatchFromClient(
		nosc.clientset.CoreV1().RESTClient(),
		"namespaces",
		metav1.NamespaceAll,
		fields.Everything())

	_, controller := cache.NewInformer(
		source,
		&v1.Namespace{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				ns := obj.(*v1.Namespace)
				receiver <- &api.NamespaceEvent{Type: api.Added, Name: ns.ObjectMeta.Name, Annotations: ns.GetAnnotations()}
			},
			UpdateFunc: func(oldobj, newobj interface{}) {
				ns := newobj.(*v1.Namespace)
				receiver <- &api.NamespaceEvent{Type: api.Modified, Name: ns.ObjectMeta.Name, Annotations: ns.GetAnnotations()}
			},
			DeleteFunc: func(obj interface{}) {
				ns := obj.(*v1.Namespace)
				receiver <- &api.NamespaceEvent{Type: api.Deleted, Name: ns.ObjectMeta.Name, Annotations: ns.GetAnnotations()}
			},
		})
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	controller.Run(ctx.Done())
	return nil
}

func (nosc *NuageClusterClient) GetServices(listOpts *metav1.ListOptions) (*[]*api.ServiceEvent, error) {
	services, err := nosc.clientset.CoreV1().Services(metav1.NamespaceAll).List(*listOpts)
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
	source := cache.NewListWatchFromClient(
		nosc.clientset.CoreV1().RESTClient(),
		"services",
		metav1.NamespaceAll,
		fields.Everything())

	_, controller := cache.NewInformer(
		source,
		&v1.Service{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				service := obj.(*v1.Service)
				labels := GetNuageLabels(service)
				if label, exists := labels["private-service"]; !exists || strings.ToLower(label) == "false" {
					receiver <- &api.ServiceEvent{Type: api.Added, Name: service.ObjectMeta.Name, ClusterIP: service.Spec.ClusterIP, Namespace: service.ObjectMeta.Namespace, NuageLabels: labels}
				}
			},
			DeleteFunc: func(obj interface{}) {
				service := obj.(*v1.Service)
				labels := GetNuageLabels(service)
				if label, exists := labels["private-service"]; !exists || strings.ToLower(label) == "false" {
					receiver <- &api.ServiceEvent{Type: api.Deleted, Name: service.ObjectMeta.Name, ClusterIP: service.Spec.ClusterIP, Namespace: service.ObjectMeta.Namespace, NuageLabels: labels}
				}
			},
		})
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	controller.Run(ctx.Done())
	return nil
}

func (nosc *NuageClusterClient) GetPod(name string, ns string) (*api.PodEvent, error) {
	if ns == "" {
		ns = metav1.NamespaceAll
	}
	pod, err := nosc.clientset.CoreV1().Pods(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &api.PodEvent{Type: api.Added, Name: pod.Name, Namespace: pod.Namespace, Labels: pod.Labels}, nil
}

func (nosc *NuageClusterClient) GetPods(listOpts *metav1.ListOptions, ns string) (*[]*api.PodEvent, error) {
	if ns == "" {
		ns = metav1.NamespaceAll
	}
	pods, err := nosc.clientset.CoreV1().Pods(ns).List(*listOpts)
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
	source := cache.NewListWatchFromClient(
		nosc.clientset.CoreV1().RESTClient(),
		"pods",
		metav1.NamespaceAll,
		fields.Everything())

	_, controller := cache.NewInformer(
		source,
		&v1.Pod{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				receiver <- &api.PodEvent{Type: api.Added, Name: pod.Name, Namespace: pod.Namespace, Labels: pod.Labels}
			},
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				receiver <- &api.PodEvent{Type: api.Deleted, Name: pod.Name, Namespace: pod.Namespace, Labels: pod.Labels}
			},
		})
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	controller.Run(ctx.Done())
	return nil
}

func (nosc *NuageClusterClient) GetNetworkPolicies(listOpts *metav1.ListOptions) (*[]*api.NetworkPolicyEvent, error) {
	policies, err := nosc.clientset.NetworkingV1().NetworkPolicies(metav1.NamespaceAll).List(*listOpts)
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
	source := cache.NewListWatchFromClient(
		nosc.clientset.NetworkingV1().RESTClient(),
		"networkpolicies",
		metav1.NamespaceAll,
		fields.Everything())

	_, controller := cache.NewInformer(
		source,
		&networkingV1.NetworkPolicy{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				policy := obj.(*networkingV1.NetworkPolicy)
				receiver <- &api.NetworkPolicyEvent{Type: api.Added, Name: policy.Name, Namespace: policy.Namespace, Policy: policy.Spec, Labels: policy.Labels}
			},
			DeleteFunc: func(obj interface{}) {
				policy := obj.(*networkingV1.NetworkPolicy)
				receiver <- &api.NetworkPolicyEvent{Type: api.Deleted, Name: policy.Name, Namespace: policy.Namespace, Policy: policy.Spec, Labels: policy.Labels}
			},
		})
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	controller.Run(ctx.Done())
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

func GetNuageLabels(input *v1.Service) map[string]string {
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
