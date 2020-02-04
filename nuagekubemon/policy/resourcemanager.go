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
	"strconv"
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/policyapi/implementer"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/policyapi/policies"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/policy/translator"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
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
		podTargetSelector, err := metav1.LabelSelectorAsSelector(&pe.Policy.PodSelector)
		if err == nil {
			targetSelectorStr := podTargetSelector.String()
			if _, found := rm.policyPgMap[pe.Name][targetSelectorStr]; !found {
				name := "PG Target For " + pe.Name
				pgName, pgId, err := rm.callBacks.AddPg(name, targetSelectorStr)
				if err == nil {
					rm.policyPgMap[pe.Name][targetSelectorStr] = api.PgInfo{PgName: pgName, PgId: pgId, Selector: pe.Policy.PodSelector}
					//get pods for this selector and add them to pg.
					var podList []string
					if pods, err := rm.clusterClientCallBacks.FilterPods(&metav1.ListOptions{LabelSelector: podTargetSelector.String(), FieldSelector: fields.Everything().String()}, ""); err == nil {
						for _, pod := range *pods {
							podList = append(podList, pod.Name)
						}
						if err = rm.callBacks.AddPortsToPg(pgId, podList); err != nil {
							glog.Errorf("Couldn't add ports %s to policy group %s", podList, pgId)
						}
					} else {
						glog.Error("Couldn't retrieve pods from the cluster client")
					}

				} else {
					glog.Errorf("Couldn't create policy group for %s", targetSelectorStr)
				}
			} else {
				glog.Infof("Policy group for targetSelectorStr %s exists already", targetSelectorStr)
			}
		}
		for i, ingressRule := range pe.Policy.Ingress {
			for f, from := range ingressRule.From {
				if from.PodSelector != nil {
					sourceSelector, err := metav1.LabelSelectorAsSelector(from.PodSelector)
					if err == nil {
						sourceSelectorStr := sourceSelector.String()
						if _, found := rm.policyPgMap[pe.Name][sourceSelectorStr]; !found {
							name := "PG Source " + strconv.Itoa(i) + "-" + strconv.Itoa(f) + " For " + pe.Name
							pgName, pgId, err := rm.callBacks.AddPg(name, sourceSelectorStr)
							if err == nil {
								rm.policyPgMap[pe.Name][sourceSelectorStr] = api.PgInfo{PgName: pgName, PgId: pgId, Selector: *from.PodSelector}
								//get pods for this selector and add them to pg.
								var podList []string
								if pods, err := rm.clusterClientCallBacks.FilterPods(&metav1.ListOptions{LabelSelector: sourceSelector.String(), FieldSelector: fields.Everything().String()}, ""); err == nil {
									for _, pod := range *pods {
										podList = append(podList, pod.Name)
									}
									if err = rm.callBacks.AddPortsToPg(pgId, podList); err != nil {
										glog.Errorf("Couldn't add ports %s to policy group %s", podList, pgId)
									}
								} else {
									glog.Error("Couldn't retrieve pods from the cluster client")
								}

							} else {
								glog.Errorf("Couldn't create policy group for %s", sourceSelectorStr)

							}
						} else {
							glog.Infof("Policy group for sourceSelectorStr %s exists already", sourceSelectorStr)
						}
					}
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
			podTargetSelector, err := metav1.LabelSelectorAsSelector(&pe.Policy.PodSelector)
			if err == nil {
				targetSelectorStr := podTargetSelector.String()
				if pgInfo, found := rm.policyPgMap[pe.Name][targetSelectorStr]; !found {
					glog.Errorf("Policy group for targetSelectorStr %s is not found", targetSelectorStr)
				} else {
					//first unassign pods from pg
					if err = rm.callBacks.DeletePortsFromPg(pgInfo.PgId); err != nil {
						glog.Errorf("Couldn't remove ports from policy group %s", pgInfo.PgId)
					} else {
						if err := rm.callBacks.DeletePg(pgInfo.PgId); err != nil {
							glog.Errorf("Failed to delete policy group %s", pgInfo.PgId)
						} else {
							delete(rm.policyPgMap[pe.Name], targetSelectorStr)
						}
					}
				}
			}
			for _, ingressRule := range pe.Policy.Ingress {
				for _, from := range ingressRule.From {
					if from.PodSelector != nil {
						sourceSelector, err := metav1.LabelSelectorAsSelector(from.PodSelector)
						if err == nil {
							sourceSelectorStr := sourceSelector.String()
							if pgInfo, found := rm.policyPgMap[pe.Name][sourceSelectorStr]; !found {
								glog.Errorf("Policy group for sourceSelectorStr %s is not found", sourceSelectorStr)
							} else {
								//first unassign pods from pg
								if err = rm.callBacks.DeletePortsFromPg(pgInfo.PgId); err != nil {
									glog.Errorf("Couldn't remove ports from policy group %s", pgInfo.PgId)
								} else {
									if err := rm.callBacks.DeletePg(pgInfo.PgId); err != nil {
										glog.Errorf("Failed to delete policy group %s", pgInfo.PgId)
									} else {
										delete(rm.policyPgMap[pe.Name], sourceSelectorStr)
									}
								}
							}
						}
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
