package policy

import (
	"fmt"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	"github.com/nuagenetworks/vspk-go/vspk"
	networkingV1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

/*
* Unit test file for policy creation and deletion
*
 */

const (
	ENTERPRISE = "test-enterprise"
	DOMAIN     = "test-domain"
	ZONE1      = "zone1"
	ZONE2      = "zone2"
)

type objIds struct {
	enterpriseID     string
	domainTemplateID string
	domainID         string
	zone1ID          string
	zone2ID          string
}

var ids objIds

func (t *testing.T) init() {
	t.createEnterprise()
	t.createDomainTemplate()
	t.createDomain()
	t.createZones()
	NewResourceManager()
}

func (t *testing.T) deinit() {
	t.deleteZones()
	t.deleteDomain()
	t.deleteDomainTemplate()
	t.deleteEnterprise()
}

func createPolicy(p *api.NetworkPolicyEvent) {

}

func removePolicy(p *api.NetworkPolicyEvent) {

}

func checkIfPolicyCreated(p *api.NetworkPolicyEvent) {

}

func checkIfPolicyRemoved(p *api.NetworkPolicyEvent) {

}

func TestPolicyFramework(t *testing.T) {

	allPolicies := []*networkingV1.NetworkPolicy{
		&networkingV1.NetworkPolicy{
			Name:      "test-np",
			Namespace: "test-ns",
			Spec: networkingV1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: map[string]string{"a": "b"},
				},
				PolicyTypes: []networkingV1.PolicyType{networkingV1.PolicyTypeIngress},
			},
		},
		&networkingV1.NetworkPolicy{
			Name:      "test-np",
			Namespace: "test-ns",
			Spec: networkingV1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: map[string]string{"a": "b"},
				},
				Ingress:     []networkingV1.NetworkPolicyIngressRule{},
				PolicyTypes: []networkingV1.PolicyType{networkingV1.PolicyTypeIngress},
			},
		},
		&networkingV1.NetworkPolicy{
			Name:      "test-np",
			Namespace: "test-ns",
			Spec: networkingV1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: map[string]string{"a": "b"},
				},
				Ingress: []networkingV1.NetworkPolicyIngressRule{
					{
						From: []networkingV1.NetworkPolicyPeer{
							{
								PodSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{"c": "d"},
								},
								NamespaceSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{"c": "d"},
								},
							},
						},
					},
				},
				Egress: []networkingV1.NetworkPolicyEgressRule{
					{
						To: []networkingV1.NetworkPolicyPeer{
							{
								NamespaceSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{"c": "d"},
								},
							},
						},
					},
				},
				PolicyTypes: []networkingV1.PolicyType{networkingV1.PolicyTypeIngress, networkingV1.PolicyTypeEgress},
			},
		},
		&networkingV1.NetworkPolicy{
			Name:      "test-np",
			Namespace: "test-ns",
			Spec: networkingV1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: map[string]string{"a": "b"},
				},
				Egress: []networkingV1.NetworkPolicyEgressRule{
					{
						To: []networkingV1.NetworkPolicyPeer{
							{
								NamespaceSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{"c": "d"},
								},
							},
						},
					},
				},
				PolicyTypes: []networkingV1.PolicyType{networkingV1.PolicyTypeIngress, networkingV1.PolicyTypeEgress},
			},
		},
		&networkingV1.NetworkPolicy{
			Name:      "test-np",
			Namespace: "test-ns",
			Spec: networkingV1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: map[string]string{"a": "b"},
				},
				Egress: []networkingV1.NetworkPolicyEgressRule{
					{
						To: []networkingV1.NetworkPolicyPeer{
							{
								NamespaceSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{"Egress": "only"},
								},
							},
						},
					},
				},
				PolicyTypes: []networkingV1.PolicyType{networkingV1.PolicyTypeEgress},
			},
		},
	}

	init()

	for _, policy := range allPolicies {
		nuagePolicyEvent := &api.NetworkPolicyEvent{
			Type:      api.Added,
			Name:      policy.Name,
			Namespace: policy.Namespace,
			Policy:    policy.Spec,
			Labels:    policy.Labels,
		}

		createPolicy(nuagePolicyEvent)
		checkIfPolicyCreated(nuagePolicyEvent)

		nuagePolicyEvent.Type = api.Deleted
		removePolicy(nuagePolicyEvent)
		checkIfPolicyRemoved(nuagePolicyEvent)
	}

	deinit()
}

func (t *testing.T) createEnterprise() {
	enterprise := vspk.NewEnterprise()
	enterprise.Name = ENTERPRISE
	if err := enterprise.Save(); err != nil {
		t.Fatalf("creating enterprise failed with error %v", err)
	}
	ids.enterpriseID = enterprise.ID
}

func (t *testing.T) deleteEnterprise() {
	enterprise := vspk.NewEnterprise()
	enterprise.ID = ids.enterpriseID
	if err := enterprise.Delete(); err != nil {
		t.Fatalf("deleting enterprise failed with error %v", err)
	}
}

func (t *testing.T) createDomainTemplate() {
	domainTemplate := vspk.NewDomainTemplate()
	domainTemplate.Name = DOMAIN
	domainTemplate.ParentID = ids.enterpriseID
	if err := domainTemplate.Save(); err != nil {
		t.Fatalf("creating domain template failed with error %v", err)
	}
	ids.domainTemplateID = domainTemplate.ID
}

func (t *testing.T) deleteDomainTemplate() {
	domainTemplate := vspk.NewDomainTemplate()
	domainTemplate.ID = ids.domainTemplateID
	if err := domainTemplate.Delete(); err != nil {
		t.Fatalf("deleting domain template failed with error %v", err)
	}
}

func (t *testing.T) createDomain() {
	domain := vspk.NewDomain()
	domain.Name = DOMAIN
	domain.ParentID = ids.enterpriseID
	domain.TemplateID = ids.domainTemplateID
	if err := domain.Save(); err != nil {
		t.Fatalf("creating domain failed with error %v", err)
	}
	ids.domainID = domain.ID
}

func (t *testing.T) deleteDomain() {
	domain := vspk.NewDomain()
	domain.ID = ids.domainID
	if err := domain.Delete(); err != nil {
		t.Fatalf("deleting domain failed with error %v", err)
	}
}

func (t *testing.T) createZones() {
	zone1 := vspk.NewZone()
	zone1.Name = ZONE1
	zone1.ParentID = ids.domainID
	if err := zone1.Save(); err != nil {
		t.Fatalf("creating zone(%s) failed with error %v", ZONE1, err)
	}
	ids.zone1ID = zone1.ID

	zone2 := vspk.NewZone()
	zone2.Name = ZONE2
	zone2.ParentID = ids.domainID
	if err := zone2.Save(); err != nil {
		t.Fatalf("creating zone(%s) failed with error %v", ZONE2, err)
	}
	ids.zone2ID = zone2.ID
}

func (t *testing.T) deleteZones() {
	zone1 := vspk.NewZone()
	zone1.ID = ids.zone1ID
	if err := zone1.Delete(); err != nil {
		t.Fatalf("deleting zone(%s) failed with error %v", ZONE1, err)
	}

	zone2 := vspk.NewZone()
	zone2.ID = ids.zone2ID
	if err := zone2.Delete(); err != nil {
		t.Fatalf("deleting zone(%s) failed with error %v", ZONE2, err)
	}
}
