package translator

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	xlateApi "github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/apis/translate"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNamespacesWithLabel(t *testing.T) {
	g := NewGomegaWithT(t)

	callBacks := &xlateApi.CallBacks{}
	vsdMetaData := make(xlateApi.VSDMetaData)
	clusterCallBacks := &api.ClusterClientCallBacks{}

	selectorLabel := &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}}

	empty := []string{}
	expected := []string{"a", "b", "c"}
	nsEvents := []*api.NamespaceEvent{
		&api.NamespaceEvent{
			Name: "a",
		},
		&api.NamespaceEvent{
			Name: "b",
		},
		&api.NamespaceEvent{
			Name: "c",
		},
	}

	filterNamespaces := func(lo *metav1.ListOptions) (*[]*api.NamespaceEvent, error) {
		selector, _ := metav1.ParseToLabelSelector(lo.LabelSelector)
		if selector.String() == selectorLabel.String() {
			return &nsEvents, nil
		}
		return nil, fmt.Errorf("could not match label")
	}
	clusterCallBacks.FilterNamespaces = filterNamespaces

	rm, err := NewResourceManager(callBacks, clusterCallBacks, &vsdMetaData)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(rm).NotTo(BeNil())

	err = rm.populateNamespacesWithLabel(nil)
	g.Expect(err).To(HaveOccurred())

	err = rm.populateNamespacesWithLabel(selectorLabel)
	g.Expect(err).NotTo(HaveOccurred())

	found := []string{}
	found, err = rm.getNamespacesWithLabel(selectorLabel)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(expected).To(Equal(found))

	_, err = rm.getNamespacesWithLabel(nil)
	g.Expect(err).To(HaveOccurred())

	label2 := &metav1.LabelSelector{MatchLabels: map[string]string{"foo": "bar1"}}
	found2, err2 := rm.getNamespacesWithLabel(label2)
	g.Expect(err2).NotTo(HaveOccurred())
	g.Expect(len(empty)).To(Equal(len(found2)))
}
