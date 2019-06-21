package translator

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	xlateApi "github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/apis/translate"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
)

func (rm *ResourceManager) deletePGInfo(ns, selector string) {
	pgInfo, _ := rm.vsdObjsMap.PGMap[ns][selector]
	pgInfo.RefCount = pgInfo.RefCount - 1
	if pgInfo.RefCount == 0 {
		delete(rm.vsdObjsMap.PGMap[ns], selector)
	}
}

func (rm *ResourceManager) findPGInfo(ns, selector string) (*xlateApi.PgInfo, bool) {
	pgInfo, ok := rm.vsdObjsMap.PGMap[ns][selector]
	return pgInfo, ok
}

func (rm *ResourceManager) appendPGInfo(ns, selector string, pgInfo *xlateApi.PgInfo) {
	if _, ok := rm.vsdObjsMap.PGMap[ns]; !ok {
		rm.vsdObjsMap.PGMap[ns] = make(map[string]*xlateApi.PgInfo)
	}
	rm.vsdObjsMap.PGMap[ns][selector] = pgInfo
}

func (rm *ResourceManager) incPGInfoRefCount(ns, selector string) {
	pgInfo, _ := rm.vsdObjsMap.PGMap[ns][selector]
	pgInfo.RefCount = pgInfo.RefCount + 1
}

func (rm *ResourceManager) createPgAddVports(selectorLabel *metav1.LabelSelector, ns, name string) error {
	var pgID string
	var err error

	if pgID, err = rm.createPG(selectorLabel, ns, name); err != nil {
		glog.Infof("creating pg failed %v", err)
		return err
	} else if pgID == "" {
		glog.Infof("pg already exists")
		return nil
	}

	if err = rm.addVPorts(pgID, selectorLabel, ns); err != nil {
		glog.Errorf("add vport to pg %s failed: %v", pgID, err)
		return err
	}

	return nil
}

func (rm *ResourceManager) createPG(selectorLabel *metav1.LabelSelector, ns, name string) (string, error) {

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
	pgName := ns + " " + name + " " + podSelectorStr

	if _, found := rm.findPGInfo(ns, podSelectorStr); found {
		rm.incPGInfoRefCount(ns, podSelectorStr)
		glog.Infof("Policy group for selector %s exists already. nothing to do", podSelectorStr)
		return "", nil
	}

	_, pgID, err = rm.callBacks.AddPg(pgName, podSelectorStr)
	if pgID == "" && err != nil {
		glog.Errorf("creating policy group for %s failed %v", podSelectorStr, err)
		return "", err
	}

	rm.appendPGInfo(ns, podSelectorStr, &xlateApi.PgInfo{
		PgName:     pgName,
		PgID:       pgID,
		Selector:   *selectorLabel,
		PolicyName: name,
		RefCount:   1,
	})

	return pgID, nil
}

func (rm *ResourceManager) addVPorts(id string, selectorLabel *metav1.LabelSelector, ns string) error {
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
		FieldSelector: fields.Everything().String()}, ns); err != nil {
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
	var pgInfo *xlateApi.PgInfo
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

	if err = rm.callBacks.DeletePortsFromPg(pgInfo.PgID); err != nil {
		glog.Errorf("Removing ports from policy group %s failed: %v", pgInfo.PgID, err)
		return err
	}

	if err = rm.callBacks.DeletePg(pgInfo.PgID); err != nil {
		glog.Errorf("deleting policy group %s failed %v", pgInfo.PgID, err)
		return err
	}

	rm.deletePGInfo(pe.Namespace, podSelectorStr)
	return nil
}
