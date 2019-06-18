package policy

import (
	"github.com/golang/glog"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//NamespaceLabelsMap label to list of namespaces mapping
type NamespaceLabelsMap map[string][]string

func (rm *ResourceManager) findNamespacesWithLabel(selectorLabel *metav1.LabelSelector) error {
	var err error
	var namespaces *[]*api.NamespaceEvent
	namespaceList := []string{}

	if selectorLabel == nil {
		return nil
	}

	nsSelectorLabel, err := metav1.LabelSelectorAsSelector(selectorLabel)
	if err != nil {
		glog.Errorf("Extracting namespace label failed %v", err)
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

	rm.nsLabelsMap[nsSelectorStr] = namespaceList
	return nil
}
