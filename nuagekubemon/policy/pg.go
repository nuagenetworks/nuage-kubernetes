package policy

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
)

//SelectorMap map of selectors to policy groups
type SelectorMap map[string]*api.PgInfo

//PgMap map of namespace to policies in a namespace
type PgMap map[string]SelectorMap

func (rm *ResourceManager) initPGInfoMap() {
	rm.pgMap = make(PgMap)
}

func (rm *ResourceManager) deletePGInfo(ns, selector string) {
	pgInfo, _ := rm.pgMap[ns][selector]
	pgInfo.RefCount = pgInfo.RefCount - 1
	if pgInfo.RefCount == 0 {
		delete(rm.pgMap[ns], selector)
	}
}

func (rm *ResourceManager) findPGInfo(ns, selector string) (*api.PgInfo, bool) {
	pgInfo, ok := rm.pgMap[ns][selector]
	return pgInfo, ok
}

func (rm *ResourceManager) appendPGInfo(ns, selector string, pgInfo *api.PgInfo) {
	if _, ok := rm.pgMap[ns]; !ok {
		rm.pgMap[ns] = make(map[string]*api.PgInfo)
	}
	rm.pgMap[ns][selector] = pgInfo
}

func (rm *ResourceManager) incPGInfoRefCount(ns, selector string) {
	pgInfo, _ := rm.pgMap[ns][selector]
	pgInfo.RefCount = pgInfo.RefCount + 1
}

func (rm *ResourceManager) createPgAddVports(selectorLabel *metav1.LabelSelector, pe *api.NetworkPolicyEvent) error {
	var pgID string
	var err error

	if pgID, err = rm.createPG(selectorLabel, pe); err != nil {
		glog.Infof("creating pg failed %v", err)
		return err
	} else if pgID == "" {
		glog.Infof("pg already exists")
		return nil
	}

	if err = rm.addVPorts(pgID, selectorLabel, pe); err != nil {
		glog.Errorf("add vport to pg %s failed: %v", pgID, err)
		return err
	}

	return nil
}

func (rm *ResourceManager) createPG(selectorLabel *metav1.LabelSelector, pe *api.NetworkPolicyEvent) (string, error) {

	var err error
	var pgID string
	var podSelectorLabel labels.Selector

	if selectorLabel == nil {
		return "", errors.New("selector label is passed as nil")
	}

	podSelectorLabel, err = metav1.LabelSelectorAsSelector(selectorLabel)
	if err != nil {
		glog.Errorf("error extracting label %v", err)
		return "", err
	}

	podSelectorStr := podSelectorLabel.String()
	pgName := pe.Namespace + " " + pe.Name + " " + podSelectorStr

	if pgInfo, found := rm.findPGInfo(pe.Namespace, podSelectorStr); found {
		rm.incPGInfoRefCount(pe.Namespace, podSelectorStr)
		glog.Infof("Policy group for selector %s exists already. nothing to do", podSelectorStr)
		return "", nil
	}

	_, pgID, err = rm.callBacks.AddPg(pgName, podSelectorStr)
	if pgID == "" && err != nil {
		glog.Errorf("creating policy group for %s failed %v", podSelectorStr, err)
		return "", err
	}

	rm.appendPGInfo(pe.Namespace, podSelectorStr, &api.PgInfo{
		PgName:     pgName,
		PgId:       pgID,
		Selector:   *selectorLabel,
		PolicyName: pe.Name,
		RefCount:   1,
	})

	return pgID, nil
}

func (rm *ResourceManager) addVPorts(id string, selectorLabel *metav1.LabelSelector, pe *api.NetworkPolicyEvent) error {
	var err error
	var podList []string
	var pods *[]*api.PodEvent
	var podSelectorLabel labels.Selector

	podSelectorLabel, err = metav1.LabelSelectorAsSelector(selectorLabel)
	if err != nil {
		glog.Errorf("error extracting label %v", err)
		return err
	}

	podSelectorStr := podSelectorLabel.String()
	//get pods for this selector and add them to pg.
	if pods, err = rm.clusterClientCallBacks.FilterPods(&metav1.ListOptions{LabelSelector: podSelectorStr,
		FieldSelector: fields.Everything().String()}, pe.Namespace); err != nil {
		glog.Errorf("retrieving pods from the cluster client failed: %v", err)
		return err
	}

	for _, pod := range *pods {
		podList = append(podList, pod.Name)
	}

	if err = rm.callBacks.AddPortsToPg(id, podList); err != nil {
		glog.Errorf("adding ports %s to policy group %s failed: %v", podList, id, err)
		return err
	}

	return nil
}

func (rm *ResourceManager) destroyPgRemoveVports(selectorLabel *metav1.LabelSelector, pe *api.NetworkPolicyEvent) error {
	var err error
	var found bool
	var pgInfo *api.PgInfo
	var podSelectorLabel labels.Selector

	if selectorLabel == nil {
		return fmt.Errorf("empty label found when cleaning up pg and vports")
	}

	podSelectorLabel, err = metav1.LabelSelectorAsSelector(selectorLabel)
	if err != nil {
		glog.Errorf("error extracting label %v", err)
		return err
	}
	podSelectorStr := podSelectorLabel.String()

	if pgInfo, found = rm.findPGInfo(pe.Namespace, podSelectorStr); !found {
		glog.Errorf("Policy group for podSelectorStr %s is not found", podSelectorStr)
		return err
	}

	if err = rm.callBacks.DeletePortsFromPg(pgInfo.PgId); err != nil {
		glog.Errorf("Removing ports from policy group %s failed: %v", pgInfo.PgId, err)
		return err
	}

	if err = rm.callBacks.DeletePg(pgInfo.PgId); err != nil {
		glog.Errorf("deleting policy group %s failed %v", pgInfo.PgId, err)
		return err
	}

	rm.deletePGInfo(pe.Namespace, podSelectorStr)
	return nil
}
