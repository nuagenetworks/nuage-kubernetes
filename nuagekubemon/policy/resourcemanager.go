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
	"github.com/golang/glog"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/api"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/network-policy-engine/implementer"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/network-policy-engine/translator"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"strconv"
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

	user, ok := rm.vsdMeta["username"]
	if !ok {
		glog.Error("Couldn't initialize a implementer for vspk policies")
	}
	password, ok := rm.vsdMeta["password"]
	if !ok {
		glog.Error("Couldn't initialize a implementer for vspk policies")
	}
	org, ok := rm.vsdMeta["organization"]
	if !ok {
		glog.Error("Cannot get the master organization for vspk policies")
	}
	url, ok := rm.vsdMeta["vsdUrl"]
	if !ok {
		glog.Error("Couldn't initialize a implementer for vspk policies")
	}
	vsdCredentials := implementer.VSDCredentials{
		Username:     user,
		Password:     password,
		Organization: org,
		URL:          url,
	}
	rm.implementer.Init(&vsdCredentials)
	return nil
}

func (rm *ResourceManager) GetPolicyGroupsForPod(podName string, podNs string) (*[]string, error) {
	var pgList []string
	if pod, err := rm.clusterClientCallBacks.GetPod(podName, podNs); err == nil {
		rm.lock.Lock()
		defer rm.lock.Unlock()
		for _, pgMap := range rm.policyPgMap {
			for _, pgInfo := range pgMap {
				if selector, err := unversioned.LabelSelectorAsSelector(&pgInfo.Selector); err == nil {
					if selector.Matches(labels.Set(pod.Labels)) {
						pgList = append(pgList, pgInfo.PgName)
					}
				}
			}
		}
	}
	return &pgList, nil
}

func (rm *ResourceManager) HandleNsEvent(nsEvent *api.NamespaceEvent) error {
	// switch nsEvent.Type {
	// case api.Added:
	// case api.Modified:
	// 	//needs to handle a case where annotation for isolation exists/doesn't exist. Or previously existed and now doesn't exist.
	// 	if nuagePolicy, err := translator.CreateNuageNSPolicy(nsEvent, rm.vsdMeta); err == nil {
	// 		if notImplemented := rm.implementer.ImplementPolicy(nuagePolicy); notImplemented != nil {
	// 			glog.Errorf("Got a %s error when implementing namespace policy", notImplemented)
	// 		}
	// 	} else {
	// 		glog.Errorf("Got an error %s when creating nuage policy", err)
	// 		return errors.New("Got an error when creating nuage policy")
	// 	}
	// case api.Deleted:
	// 	//needs to handle a case where annotation for isolation exists/doesn't exist. Or previously existed and now doesn't exist.
	// 	if notImplemented := rm.implementer.DeletePolicy(nsEvent.Name); notImplemented != nil {
	// 		glog.Errorf("Got a %s error when deleting namespace policy", notImplemented)
	// 	} else {
	// 		glog.Errorf("Got an error %s when deleting nuage ns policy", err)
	// 		return errors.New("Got an error when deleting nuage ns policy")
	// 	}

	// }
	return nil
}

func (rm *ResourceManager) HandlePolicyEvent(pe *api.NetworkPolicyEvent) error {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	switch pe.Type {
	case api.Added:
		if _, ok := rm.policyPgMap[pe.Name]; !ok {
			rm.policyPgMap[pe.Name] = make(PgMap)
		}
		podTargetSelector, err := unversioned.LabelSelectorAsSelector(&pe.Policy.PodSelector)
		if err == nil {
			targetSelectorStr := podTargetSelector.String()
			if _, found := rm.policyPgMap[pe.Name][targetSelectorStr]; !found {
				name := "PG Target For " + pe.Name
				pgName, pgId, err := rm.callBacks.AddPg(name, targetSelectorStr)
				if err == nil {
					rm.policyPgMap[pe.Name][targetSelectorStr] = api.PgInfo{PgName: pgName, PgId: pgId, Selector: pe.Policy.PodSelector}
					//get pods for this selector and add them to pg.
					var podList []string
					if pods, err := rm.clusterClientCallBacks.FilterPods(&kapi.ListOptions{LabelSelector: podTargetSelector, FieldSelector: fields.Everything()}, ""); err == nil {
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
					sourceSelector, err := unversioned.LabelSelectorAsSelector(from.PodSelector)
					if err == nil {
						sourceSelectorStr := sourceSelector.String()
						if _, found := rm.policyPgMap[pe.Name][sourceSelectorStr]; !found {
							name := "PG Source " + strconv.Itoa(i) + "-" + strconv.Itoa(f) + " For " + pe.Name
							pgName, pgId, err := rm.callBacks.AddPg(name, sourceSelectorStr)
							if err == nil {
								rm.policyPgMap[pe.Name][sourceSelectorStr] = api.PgInfo{PgName: pgName, PgId: pgId, Selector: *from.PodSelector}
								//get pods for this selector and add them to pg.
								var podList []string
								if pods, err := rm.clusterClientCallBacks.FilterPods(&kapi.ListOptions{LabelSelector: sourceSelector, FieldSelector: fields.Everything()}, ""); err == nil {
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
		if nuagePolicy, err := translator.CreateNuagePGPolicy(&pe.Policy, pe.Name, rm.policyPgMap[pe.Name], rm.vsdMeta); err == nil {
			if notImplemented := rm.implementer.ImplementPolicy(nuagePolicy); notImplemented != nil {
				glog.Errorf("Got a %s error when implementing policy", notImplemented)
			}
		} else {
			glog.Errorf("Got an error %s when creating nuage policy", err)
			return errors.New("Got an error when creating nuage policy")
		}
	case api.Deleted:
		if _, ok := rm.policyPgMap[pe.Name]; !ok {
			glog.Info("No policy group map entry found for this policy")
			return errors.New("No policy group map entry found")
		} else {
			podTargetSelector, err := unversioned.LabelSelectorAsSelector(&pe.Policy.PodSelector)
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
						sourceSelector, err := unversioned.LabelSelectorAsSelector(from.PodSelector)
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
			if err := rm.implementer.DeletePolicy(pe.Name, enterprise, domain); err != nil {
				return errors.New("Got an error when deleting nuage policy")
			}
		}
	}
	return nil
}
