package translator

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (rm *ResourceManager) populateNamespacesWithLabel(selectorLabel *metav1.LabelSelector) error {
	var err error
	var namespaces *[]*api.NamespaceEvent
	namespaceList := []string{}

	if selectorLabel == nil {
		return nil
	}

	nsSelectorLabel, err := metav1.LabelSelectorAsSelector(selectorLabel)
	if err != nil {
		glog.Errorf("extracting namespace label failed %v", err)
		return err
	}
	nsSelectorStr := nsSelectorLabel.String()

	if namespaces, err = rm.clusterClientCallBacks.FilterNamespaces(&metav1.ListOptions{LabelSelector: nsSelectorStr}); err != nil {
		glog.Errorf("call to cluster client to filter namespaces %s failed: %v", nsSelectorStr, err)
		return err
	}

	for _, namespace := range *namespaces {
		namespaceList = append(namespaceList, namespace.Name)
	}

	rm.vsdObjsMap.NSLabelsMap[nsSelectorStr] = namespaceList
	return nil
}

func (rm *ResourceManager) getNamespacesWithLabel(selectorLabel *metav1.LabelSelector) ([]string, error) {

	if selectorLabel == nil {
		return nil, fmt.Errorf("selector label is nil")
	}

	nsSelectorLabel, err := metav1.LabelSelectorAsSelector(selectorLabel)
	if err != nil {
		glog.Errorf("extracting namespace label failed %v", err)
		return nil, err
	}
	nsSelectorStr := nsSelectorLabel.String()

	nsList, _ := rm.vsdObjsMap.NSLabelsMap[nsSelectorStr]
	return nsList, nil
}
