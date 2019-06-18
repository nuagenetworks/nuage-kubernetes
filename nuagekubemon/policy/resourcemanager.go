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

package policy

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/policyapi/implementer"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/policyapi/policies"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/policy/translator"
	networkingV1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

//ResourceManager for network policy constructs
type ResourceManager struct {
	pgMap                  PgMap
	networkMacroMap        NWMacroMap
	nsLabelsMap            NamespaceLabelsMap
	callBacks              *CallBacks
	clusterClientCallBacks *api.ClusterClientCallBacks
	vsdMeta                VSDMetaData
	lock                   sync.Mutex
	implementer            implementer.PolicyImplementer
	externalID             string
}

//NewResourceManager creates a new resource manager
func NewResourceManager(callBacks *CallBacks, clusterCbs *api.ClusterClientCallBacks, vsdMeta *VSDMetaData) (*ResourceManager, error) {
	rm := &ResourceManager{}
	rm.Init(callBacks, clusterCbs, vsdMeta)
	return rm, nil
}

//Init initializes the resource manager
func (rm *ResourceManager) Init(callBacks *CallBacks, clusterCbs *api.ClusterClientCallBacks, vsdMeta *VSDMetaData) {
	rm.initPGInfoMap()
	rm.callBacks = callBacks
	rm.clusterClientCallBacks = clusterCbs
	rm.vsdMeta = *vsdMeta
	rm.externalID, _ = rm.vsdMeta["externalID"]
}

//InitPolicyImplementer initializes the policy implementer
func (rm *ResourceManager) InitPolicyImplementer() error {
	var url string
	var usercert string
	var userkey string
	var ok bool

	url, ok = rm.vsdMeta["vsdUrl"]
	if !ok {
		glog.Error("Couldn't initialize a implementer for vspk policies: vsdURL absent")
		return fmt.Errorf("VSD URL absent")
	}

	if rm.vsdMeta["username"] == "" || rm.vsdMeta["password"] == "" {
		usercert, ok = rm.vsdMeta["usercertfile"]
		if !ok {
			glog.Error("Couldn't initialize a implementer for vspk policies: user certificate file absent")
			return fmt.Errorf("VSD User certificate file absent")
		}

		userkey, ok = rm.vsdMeta["userkeyfile"]
		if !ok {
			glog.Error("Couldn't initialize a implementer for vspk policies: user key file absent")
			return fmt.Errorf("VSD User key file absent")
		}
	}
	vsdCredentials := implementer.VSDCredentials{
		URL:          url,
		UserCertFile: usercert,
		UserKeyFile:  userkey,
		Username:     rm.vsdMeta["username"],
		Password:     rm.vsdMeta["password"],
		Organization: rm.vsdMeta["organization"],
	}

	return rm.implementer.Init(&vsdCredentials)
}

//GetPolicyGroupsForPod fetches the policy groups that can be applied to a pod
func (rm *ResourceManager) GetPolicyGroupsForPod(podName string, podNs string) (*[]string, error) {
	var pgList []string
	if pod, err := rm.clusterClientCallBacks.GetPod(podName, podNs); err == nil {
		rm.lock.Lock()
		defer rm.lock.Unlock()
		for _, pgInfo := range rm.pgMap[podNs] {
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

	enterprise, ok := rm.vsdMeta["enterpriseName"]
	if !ok {
		glog.Error("Failed to get enterprise for namespace annotation operation")
		return errors.New("Failed to get enterprise for namespace annotation operation")
	}

	domain, ok := rm.vsdMeta["domainName"]
	if !ok {
		glog.Error("Failed to get domain for namespace annotation operation")
		return errors.New("Failed to get domain for namespace annotation operation")
	}

	nuagePolicy := policies.NuagePolicy{
		Version:    policies.V1Alpha,
		Type:       policies.Default,
		Enterprise: enterprise,
		Domain:     domain,
		Name:       api.ZoneAnnotationTemplateName,
	}

	defaultPolicyElementTCP := policies.DefaultPolicyElement{
		Name:   fmt.Sprintf("Namespace annotation for %s - TCP", namespace),
		From:   policies.EndPoint{Name: namespace, Type: policies.Zone},
		To:     policies.EndPoint{Name: namespace, Type: policies.Zone},
		Action: policies.Deny,
		NetworkParameters: policies.NetworkParameters{
			Protocol:             policies.TCP,
			DestinationPortRange: policies.PortRange{StartPort: 1, EndPort: 65535},
		},
	}

	defaultPolicyElementUDP := policies.DefaultPolicyElement{
		Name:   fmt.Sprintf("Namespace annotation for %s - UDP", namespace),
		From:   policies.EndPoint{Name: namespace, Type: policies.Zone},
		To:     policies.EndPoint{Name: namespace, Type: policies.Zone},
		Action: policies.Deny,
		NetworkParameters: policies.NetworkParameters{
			Protocol:             policies.UDP,
			DestinationPortRange: policies.PortRange{StartPort: 1, EndPort: 65535},
		},
	}

	nuagePolicy.PolicyElements = []policies.DefaultPolicyElement{defaultPolicyElementTCP, defaultPolicyElementUDP}

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
	if err := rm.createPgAddVports(&pe.Policy.PodSelector, pe); err != nil {
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

	if nuagePolicy, err := translator.CreateNuagePGPolicy(pe, rm.pgMap[pe.Namespace], rm.vsdMeta, rm.nsLabelsMap); err == nil {
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
	if _, ok := rm.pgMap[pe.Namespace]; !ok {
		glog.Info("No policy group map entry found for this policy")
		return errors.New("No policy group map entry found")
	}

	if err := rm.destroyPgRemoveVports(&pe.Policy.PodSelector, pe); err != nil {
		glog.Errorf("removing vports and deleting pg failed: %v", err)
		return err
	}

	for _, ingressRule := range pe.Policy.Ingress {
		for _, from := range ingressRule.From {
			if err := rm.destroyPgRemoveVports(from.PodSelector, pe); err != nil {
				glog.Errorf("removing vports and deleting pg failed: %v", err)
			}
		}
	}

	for _, egressRule := range pe.Policy.Egress {
		for _, to := range egressRule.To {
			if err := rm.destroyPgRemoveVports(to.PodSelector, pe); err != nil {
				glog.Errorf("removing vports and deleting pg failed: %v", err)
			}
		}
	}

	rm.deletePGInfo(pe.Namespace, pe.Name)
	enterprise, ok := rm.vsdMeta["enterpriseName"]

	if !ok {
		glog.Error("Failed to get enterprise when deleting policy")
		return errors.New("Failed to get enterprise when deleting policy")
	}

	domain, ok := rm.vsdMeta["domainName"]

	if !ok {
		glog.Error("Failed to get domain when deleting policy")
		return errors.New("Failed to get domain when deleting policy")
	}

	glog.Infof("Trying to delete policy %+v", pe.Name)

	if err := rm.implementer.DeletePolicy(pe.Name, enterprise, domain); err != nil {
		return errors.New("Got an error when deleting nuage policy")
	}
	return nil
}

func (rm *ResourceManager) translatePeerPolicy(peer networkingV1.NetworkPolicyPeer, pe *api.NetworkPolicyEvent) error {

	if peer.IPBlock != nil {
		if (peer.PodSelector != nil) || (peer.NamespaceSelector != nil) {
			return fmt.Errorf("unsupported network policy. ip block and pod/namespace selector cannot be specified together")
		}
	}

	if peer.PodSelector != nil && peer.NamespaceSelector != nil {
		if err := rm.createPgAddVports(peer.PodSelector, pe); err != nil {
			glog.Errorf("converting pod label to vports and adding them to pg failed: %v", err)
			return err
		}
	} else if peer.NamespaceSelector != nil {
		if err := rm.findNamespacesWithLabel(peer.NamespaceSelector); err != nil {
			glog.Errorf("finding namespaces from selector label %s failed: %v", peer.NamespaceSelector.String(), err)
			return err
		}
	} else if peer.PodSelector != nil {
		if err := rm.createPgAddVports(peer.PodSelector, pe); err != nil {
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
