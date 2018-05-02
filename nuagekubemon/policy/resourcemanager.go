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
	"github.com/golang/glog"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/policyapi/implementer"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/policyapi/policies"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/policy/translator"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"strconv"
	"strings"
	"sync"
)

type CreatePgFunc func(string, string) (string, string, error)
type DeletePgFunc func(string) error
type AddPortsToPgFunc func(string, []string) error
type DeletePortsFromPgFunc func(string) error

//map of LabelSelector as selector as string => pg corresponding to it.

type CallBacks struct {
	AddPg             CreatePgFunc
	DeletePg          DeletePgFunc
	AddPortsToPg      AddPortsToPgFunc
	DeletePortsFromPg DeletePortsFromPgFunc
}

//map of label selector for a policy group to policy group info
type PgMap map[string]api.PgInfo

//map of policy name to policy group map
type PolicyPgMap map[string]PgMap

type VsdMetaData map[string]string

type ResourceManager struct {
	policyPgMap            PolicyPgMap
	callBacks              CallBacks
	clusterClientCallBacks api.ClusterClientCallBacks
	vsdMeta                VsdMetaData
	lock                   sync.Mutex
	implementer            implementer.PolicyImplementer
}

func NewResourceManager(callBacks *CallBacks, clusterCbs *api.ClusterClientCallBacks, vsdMeta *VsdMetaData) (*ResourceManager, error) {
	rm := new(ResourceManager)
	if err := rm.Init(callBacks, clusterCbs, vsdMeta); err != nil {
		glog.Error("Cannot instantiate a new resource manager")
		return rm, err
	} else {
		return rm, nil
	}
}

func (rm *ResourceManager) Init(callBacks *CallBacks, clusterCbs *api.ClusterClientCallBacks, vsdMeta *VsdMetaData) error {
	rm.policyPgMap = make(PolicyPgMap)
	rm.callBacks = *callBacks
	rm.clusterClientCallBacks = *clusterCbs
	rm.vsdMeta = *vsdMeta

	return nil
}

func (rm *ResourceManager) InitPolicyImplementer() error {
	url, ok := rm.vsdMeta["vsdUrl"]
	if !ok {
		glog.Error("Couldn't initialize a implementer for vspk policies: vsdURL absent")
		return fmt.Errorf("VSD URL absent")
	}

	usercert, ok := rm.vsdMeta["usercertfile"]
	if !ok {
		glog.Error("Couldn't initialize a implementer for vspk policies: user certificate file absent")
		return fmt.Errorf("VSD User certificate file absent")
	}

	userkey, ok := rm.vsdMeta["userkeyfile"]
	if !ok {
		glog.Error("Couldn't initialize a implementer for vspk policies: user key file absent")
		return fmt.Errorf("VSD User key file absent")
	}

	vsdCredentials := implementer.VSDCredentials{
		URL:          url,
		UserCertFile: usercert,
		UserKeyFile:  userkey,
	}

	return rm.implementer.Init(&vsdCredentials)
}

func (rm *ResourceManager) GetPolicyGroupsForPod(podName string, podNs string) (*[]string, error) {
	var pgList []string
	if pod, err := rm.clusterClientCallBacks.GetPod(podName, podNs); err == nil {
		rm.lock.Lock()
		defer rm.lock.Unlock()
		for _, pgMap := range rm.policyPgMap {
			for _, pgInfo := range pgMap {
				if selector, err := metav1.LabelSelectorAsSelector(&pgInfo.Selector); err == nil {
					if selector.Matches(labels.Set(pod.Labels)) {
						pgList = append(pgList, pgInfo.PgName)
					}
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
		Name:   fmt.Sprint("Namespace annotation for %s - TCP", namespace),
		From:   policies.EndPoint{Name: namespace, Type: policies.Zone},
		To:     policies.EndPoint{Name: namespace, Type: policies.Zone},
		Action: policies.Deny,
		NetworkParameters: policies.NetworkParameters{
			Protocol:             policies.TCP,
			DestinationPortRange: policies.PortRange{StartPort: 1, EndPort: 65535},
		},
	}

	defaultPolicyElementUDP := policies.DefaultPolicyElement{
		Name:   fmt.Sprint("Namespace annotation for %s - UDP", namespace),
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
	return nil
}

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
		if _, ok := rm.policyPgMap[pe.Name]; !ok {
			rm.policyPgMap[pe.Name] = make(PgMap)
		}
		pgName := "PG Target For " + pe.Name
		if err := rm.createPgAddVports(&pe.Policy.PodSelector, pe, pgName); err != nil {
			glog.Errorf("converting pod label to vports and adding them to pg failed: %v", err)
			return err
		}
		for i, ingressRule := range pe.Policy.Ingress {
			for f, from := range ingressRule.From {
				if from.PodSelector == nil {
					continue
				}
				pgName := "PG Source " + strconv.Itoa(i) + "-" + strconv.Itoa(f) + " For " + pe.Name
				if err := rm.createPgAddVports(from.PodSelector, pe, pgName); err != nil {
					glog.Errorf("converting pod label to vports and adding them to pg failed: %v", err)
					return err
				}
			}
		}
		if nuagePolicy, err := translator.CreateNuagePGPolicy(pe, rm.policyPgMap[pe.Name], rm.vsdMeta); err == nil {
			if notImplemented := rm.implementer.ImplementPolicy(nuagePolicy); notImplemented != nil {
				glog.Errorf("Got a %s error when implementing policy", notImplemented)
			}
		} else {
			glog.Errorf("Got an error %s when creating nuage policy", err)
			return errors.New("Got an error when creating nuage policy")
		}
	case api.Deleted:
		glog.Infof("Starting deletion of policy %+v", pe.Name)
		if _, ok := rm.policyPgMap[pe.Name]; !ok {
			glog.Info("No policy group map entry found for this policy")
			return errors.New("No policy group map entry found")
		} else {
			if err := rm.destroyPgRemoveVports(&pe.Policy.PodSelector, pe); err != nil {
				glog.Errorf("removing vports and deleting pg failed: %v, err")
				return err
			}
			for _, ingressRule := range pe.Policy.Ingress {
				for _, from := range ingressRule.From {
					if from.PodSelector == nil {
						continue
					}
					if err := rm.destroyPgRemoveVports(from.PodSelector, pe); err != nil {
						glog.Errorf("removing vports and deleting pg failed: %v, err")
						return err
					}
				}
			}
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
		}
	}
	return nil
}

func (rm *ResourceManager) createPgAddVports(selectorLabel *metav1.LabelSelector, pe *api.NetworkPolicyEvent, pgName string) error {
	var err error
	var pgId string
	var podList []string
	var pods *[]*api.PodEvent
	var podSelectorLabel labels.Selector

	podSelectorLabel, err = metav1.LabelSelectorAsSelector(selectorLabel)
	if err != nil {
		glog.Errorf("error extracting label failed %v", err)
		return err
	}
	podSelectorStr := podSelectorLabel.String()

	if _, found := rm.policyPgMap[pe.Name][podSelectorStr]; found {
		glog.Infof("Policy group for podSelectorStr %s exists already", podSelectorStr)
		return nil
	}

	_, pgId, err = rm.callBacks.AddPg(pgName, podSelectorStr)
	if err != nil {
		glog.Errorf("creating policy group for %s failed %v", podSelectorStr, err)
		return err
	}
	rm.policyPgMap[pe.Name][podSelectorStr] = api.PgInfo{PgName: pgName, PgId: pgId, Selector: pe.Policy.PodSelector}

	//get pods for this selector and add them to pg.
	if pods, err = rm.clusterClientCallBacks.FilterPods(&metav1.ListOptions{LabelSelector: podSelectorLabel.String(),
		FieldSelector: fields.Everything().String()}, pe.Namespace); err != nil {
		glog.Error("retrieving pods from the cluster client failed: %v", err)
		return err
	}
	for _, pod := range *pods {
		podList = append(podList, pod.Name)
	}

	if err = rm.callBacks.AddPortsToPg(pgId, podList); err != nil {
		glog.Errorf("adding ports %s to policy group %s failed: %v", podList, pgId, err)
		return err
	}
	return nil
}

func (rm *ResourceManager) destroyPgRemoveVports(selectorLabel *metav1.LabelSelector, pe *api.NetworkPolicyEvent) error {
	var err error
	var found bool
	var pgInfo api.PgInfo
	var podSelectorLabel labels.Selector

	podSelectorLabel, err = metav1.LabelSelectorAsSelector(selectorLabel)
	if err != nil {
		glog.Errorf("error extracting label %v", err)
		return err
	}
	podSelectorStr := podSelectorLabel.String()

	if pgInfo, found = rm.policyPgMap[pe.Name][podSelectorStr]; !found {
		glog.Errorf("Policy group for podSelectorStr %s is not found", podSelectorStr)
		return err
	}
	//first unassign pods from pg
	if err = rm.callBacks.DeletePortsFromPg(pgInfo.PgId); err != nil {
		glog.Errorf("Removing ports from policy group %s failed: %v", pgInfo.PgId, err)
		return err
	}
	if err = rm.callBacks.DeletePg(pgInfo.PgId); err != nil {
		glog.Errorf("deleting policy group %s failed %v", pgInfo.PgId, err)
		return err
	}
	delete(rm.policyPgMap[pe.Name], podSelectorStr)
	return nil
}
