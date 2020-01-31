/*
###########################################################################
#
#   Filename:           resourcemanager.go
#
#   Author:             Aniket Bhat
#   Created:            October 27, 2016
#
#   Description:        Resource manager for policy objects specific to
#
###########################################################################
#
#              Copyright (c) 2015 Nuage Networks
#
###########################################################################
*/

package translator

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	xlateApi "github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/apis/translate"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/policyapi/implementer"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/policyapi/policies"
	kapi "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
)

//ResourceManager for network policy constructs
type ResourceManager struct {
	vsdObjsMap             *xlateApi.VSDObjsMap
	callBacks              *xlateApi.CallBacks
	clusterClientCallBacks *api.ClusterClientCallBacks
	vsdMetaData            xlateApi.VSDMetaData
	lock                   sync.Mutex
	implementer            implementer.PolicyImplementer
	externalID             string
}

//NewResourceManager creates a new resource manager
func NewResourceManager(callBacks *xlateApi.CallBacks,
	clusterCbs *api.ClusterClientCallBacks,
	vsdMeta *xlateApi.VSDMetaData) (*ResourceManager, error) {
	rm := &ResourceManager{}
	rm.Init(callBacks, clusterCbs, vsdMeta)
	return rm, nil
}

//Init initializes the resource manager
func (rm *ResourceManager) Init(callBacks *xlateApi.CallBacks,
	clusterCbs *api.ClusterClientCallBacks,
	vsdMeta *xlateApi.VSDMetaData) {
	rm.vsdObjsMap = xlateApi.InitVSDObjsMap()
	rm.callBacks = callBacks
	rm.clusterClientCallBacks = clusterCbs
	rm.vsdMetaData = *vsdMeta
	rm.externalID, _ = rm.vsdMetaData["externalID"]
}

//InitPolicyImplementer initializes the policy implementer
func (rm *ResourceManager) InitPolicyImplementer() error {
	var url string
	var usercert string
	var userkey string
	var ok bool

	url, ok = rm.vsdMetaData["vsdUrl"]
	if !ok {
		glog.Error("Couldn't initialize a implementer for vspk policies: vsdURL absent")
		return fmt.Errorf("VSD URL absent")
	}

	if rm.vsdMetaData["username"] == "" || rm.vsdMetaData["password"] == "" {
		usercert, ok = rm.vsdMetaData["usercertfile"]
		if !ok {
			glog.Error("Couldn't initialize a implementer for vspk policies: user certificate file absent")
			return fmt.Errorf("VSD User certificate file absent")
		}

		userkey, ok = rm.vsdMetaData["userkeyfile"]
		if !ok {
			glog.Error("Couldn't initialize a implementer for vspk policies: user key file absent")
			return fmt.Errorf("VSD User key file absent")
		}
	}
	vsdCredentials := implementer.VSDCredentials{
		URL:          url,
		UserCertFile: usercert,
		UserKeyFile:  userkey,
		Username:     rm.vsdMetaData["username"],
		Password:     rm.vsdMetaData["password"],
		Organization: rm.vsdMetaData["organization"],
	}

	return rm.implementer.Init(&vsdCredentials)
}

//GetPolicyGroupsForPod fetches the policy groups that can be applied to a pod
func (rm *ResourceManager) GetPolicyGroupsForPod(podName string, podNs string) (*[]string, error) {
	var pgList []string
	if pod, err := rm.clusterClientCallBacks.GetPod(podName, podNs); err == nil {
		rm.lock.Lock()
		defer rm.lock.Unlock()
		for _, pgInfo := range rm.vsdObjsMap.PGMap[podNs] {
			if selector, err := metav1.LabelSelectorAsSelector(&pgInfo.Selector); err == nil {
				if selector.Matches(labels.Set(pod.Labels)) {
					pgList = append(pgList, pgInfo.PgName)
				}
			}
		}
	}
	return &pgList, nil
}

func (rm *ResourceManager) updateZoneAnnotationTemplate(namespace string,
	updateOp policies.PolicyUpdateOperation) error {

	p := &xlateApi.PolicyData{
		Name:       fmt.Sprintf("Namespace annotation for %s - TCP", namespace),
		SourceName: namespace, SourceType: policies.Zone,
		TargetName: namespace, TargetType: policies.Zone,
		Action: policies.Deny,
	}

	nuagePolicy := policies.NuagePolicy{
		Version:        policies.V1Alpha,
		Type:           policies.Default,
		Enterprise:     rm.vsdMetaData["enterprise"],
		Domain:         rm.vsdMetaData["domain"],
		Name:           api.ZoneAnnotationTemplateName,
		PolicyElements: createNuagePolicyElements([]networkingV1.NetworkPolicyPort{}, p),
	}

	err := rm.InitPolicyImplementer()
	if err != nil {
		return fmt.Errorf("Unable to initialize policy implementer %+v", err)
	}

	if notImplemented := rm.implementer.UpdatePolicy(&nuagePolicy, updateOp); notImplemented != nil {
		glog.Errorf("Got a %s error when implementing policy", notImplemented)
		return notImplemented
	}

	return nil
}

//HandleNsEvent handles default annotations on a namespace
func (rm *ResourceManager) HandleNsEvent(nsEvent *api.NamespaceEvent) error {
	glog.Infof("Received namespace event for policy parsing %+v", nsEvent)

	switch nsEvent.Type {
	case api.Added:
		fallthrough
	case api.Modified:
		if nsEvent.Annotations != nil {
			if annotation, ok := nsEvent.Annotations["net.beta.kubernetes.io/network-policy"]; ok {
				if strings.Compare(annotation, "{\"ingress\": {\"isolation\": \"DefaultDeny\"}}") == 0 {
					glog.Info("Implementing Network policy to block intra zone communication")
					err := rm.updateZoneAnnotationTemplate(nsEvent.Name, policies.UpdateAdd)
					if err != nil {
						glog.Errorf("Unable to add annotations for namespace %s", nsEvent.Name)
						return err
					}
					glog.Infof("Successfully implemented namespace annotations for %s", nsEvent.Name)
				} else {
					err := rm.updateZoneAnnotationTemplate(nsEvent.Name, policies.UpdateRemove)
					if err != nil {
						glog.Warningf("Unable to remove annotations from namespace %s", nsEvent.Name)
						return err
					}
				}
			} else {
				err := rm.updateZoneAnnotationTemplate(nsEvent.Name, policies.UpdateRemove)
				if err != nil {
					glog.Warningf("Unable to remove annotations from namespace %s", nsEvent.Name)
					return err
				}
			}
		}

	case api.Deleted:
		err := rm.updateZoneAnnotationTemplate(nsEvent.Name, policies.UpdateRemove)
		if err != nil {
			glog.Errorf("Unable to remove annotations from namespace %s", nsEvent.Name)
			return err
		}
	}

	//TODO: Handle the case where ns labels are updated after policy is created

	return nil
}

//HandlePolicyEvent handles policy add/del event
func (rm *ResourceManager) HandlePolicyEvent(pe *api.NetworkPolicyEvent) error {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	//TODO: Verify the spec
	//TODO: Fill the defaults

	// Since vspk session times out after X amount of time, re-init the policy
	// implementer each time a new policy needs to be implemented
	err := rm.InitPolicyImplementer()
	if err != nil {
		return fmt.Errorf("Unable to initialize policy implementer %+v", err)
	}

	switch pe.Type {
	case api.Added:
		rm.policyAddEvent(pe)
	case api.Deleted:
		rm.policyDelEvent(pe)
	}
	return nil
}

func (rm *ResourceManager) policyAddEvent(pe *api.NetworkPolicyEvent) error {
	if err := rm.createPgAddVports(&pe.Policy.PodSelector, pe.Namespace, pe.Name); err != nil {
		glog.Errorf("converting pod label to vports and adding them to pg failed: %v", err)
		return err
	}

	for _, ingressRule := range pe.Policy.Ingress {
		for _, from := range ingressRule.From {
			rm.translatePeerPolicy(from, pe)
		}
	}

	for _, egressRule := range pe.Policy.Egress {
		for _, to := range egressRule.To {
			rm.translatePeerPolicy(to, pe)
		}
	}

	if nuagePolicy, err := rm.CreateNuagePGPolicy(pe); err == nil {
		if notImplemented := rm.implementer.ImplementPolicy(nuagePolicy); notImplemented != nil {
			glog.Errorf("Got a %s error when implementing policy", notImplemented)
		}
	} else {
		glog.Errorf("Got an error %s when creating nuage policy", err)
		return errors.New("Got an error when creating nuage policy")
	}

	return nil
}

func (rm *ResourceManager) policyDelEvent(pe *api.NetworkPolicyEvent) error {
	glog.Infof("Starting deletion of policy %+v", pe.Name)
	if _, ok := rm.vsdObjsMap.PGMap[pe.Namespace]; !ok {
		glog.Info("No policy group map entry found for this policy")
		return errors.New("No policy group map entry found")
	}

	if err := rm.destroyPgRemoveVports(&pe.Policy.PodSelector, pe.Namespace); err != nil {
		glog.Errorf("removing vports and deleting pg failed: %v", err)
		return err
	}

	for _, ingressRule := range pe.Policy.Ingress {
		for _, from := range ingressRule.From {
			if err := rm.deletePeerPolicy(from, pe); err != nil {
				glog.Errorf("removing vports and deleting pg failed: %v", err)
			}
		}
	}

	for _, egressRule := range pe.Policy.Egress {
		for _, to := range egressRule.To {
			if err := rm.deletePeerPolicy(to, pe); err != nil {
				glog.Errorf("removing vports and deleting pg failed: %v", err)
			}
		}
	}

	rm.deletePGInfo(pe.Namespace, pe.Name)

	if err := rm.implementer.DeletePolicy(pe.Name,
		rm.vsdMetaData["enterprise"],
		rm.vsdMetaData["domain"]); err != nil {
		return errors.New("Got an error when deleting nuage policy")
	}
	glog.Infof("policy %+v deletion completed", pe.Name)
	return nil
}

func (rm *ResourceManager) translatePeerPolicy(peer networkingV1.NetworkPolicyPeer, pe *api.NetworkPolicyEvent) error {

	if peer.IPBlock != nil {
		if (peer.PodSelector != nil) || (peer.NamespaceSelector != nil) {
			return fmt.Errorf("unsupported network policy. ip block and pod/namespace selector cannot be specified together")
		}
	}

	if peer.PodSelector != nil && peer.NamespaceSelector != nil {
		if err := rm.populateNamespacesWithLabel(peer.NamespaceSelector); err != nil {
			glog.Errorf("finding namespaces from selector label %s failed: %v", peer.NamespaceSelector.String(), err)
			return err
		}
		nsList, err := rm.getNamespacesWithLabel(peer.NamespaceSelector)
		if err != nil {
			glog.Errorf("cannot find namespaces with label %v", peer.NamespaceSelector)
			return err
		}
		for _, ns := range nsList {
			if err := rm.createPgAddVports(peer.PodSelector, ns, pe.Name); err != nil {
				glog.Errorf("converting pod label to vports and adding them to pg failed: %v", err)
				return err
			}
		}
	} else if peer.NamespaceSelector != nil {
		if err := rm.populateNamespacesWithLabel(peer.NamespaceSelector); err != nil {
			glog.Errorf("finding namespaces from selector label %s failed: %v", peer.NamespaceSelector.String(), err)
			return err
		}
	} else if peer.PodSelector != nil {
		if err := rm.createPgAddVports(peer.PodSelector, pe.Namespace, pe.Name); err != nil {
			glog.Errorf("converting pod label to vports and adding them to pg failed: %v", err)
			return err
		}
	}

	if err := rm.createNetworkMacros(peer.IPBlock, pe); err != nil {
		glog.Errorf("creating network macros failed: %v", err)
		return err
	}

	return nil
}

func (rm *ResourceManager) deletePeerPolicy(peer networkingV1.NetworkPolicyPeer, pe *api.NetworkPolicyEvent) error {
	if peer.PodSelector != nil && peer.NamespaceSelector != nil {
		nsList, err := rm.getNamespacesWithLabel(peer.NamespaceSelector)
		if err != nil {
			glog.Errorf("cannot find namespaces with label %v", peer.NamespaceSelector)
			return err
		}
		for _, ns := range nsList {
			if err := rm.destroyPgRemoveVports(peer.PodSelector, ns); err != nil {
				glog.Errorf("converting pod label to vports and adding them to pg failed: %v", err)
				return err
			}
		}
	} else if peer.NamespaceSelector != nil {
		//Do nothing
	} else if peer.PodSelector != nil {
		if err := rm.destroyPgRemoveVports(peer.PodSelector, pe.Namespace); err != nil {
			glog.Errorf("converting pod label to vports and adding them to pg failed: %v", err)
			return err
		}
	} else if peer.IPBlock != nil {
		if err := rm.deleteNetworkMacros(peer.IPBlock, pe); err != nil {
			glog.Errorf("deleting network macros failed for ip block %v with error %v", peer.IPBlock, err)
			return err
		}
	}

	return nil
}

func (rm *ResourceManager) fillDefaultPorts(ports []networkingV1.NetworkPolicyPort) {
	if len(ports) == 0 {
		protocol := kapi.ProtocolTCP
		port := intstr.FromInt(0)
		//if nothing is specified, this is the default for policy
		ports = append(ports, networkingV1.NetworkPolicyPort{
			Protocol: &protocol,
			Port:     &port,
		})
	}
}

func (rm *ResourceManager) validateNetworkSpecPorts(ports []networkingV1.NetworkPolicyPort) error {
	for _ = range ports {
		return fmt.Errorf("invalid port information")
	}
	return nil
}
