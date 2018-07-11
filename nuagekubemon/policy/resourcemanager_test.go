package policy

import (
	"fmt"
	networkingV1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

/*
* Unit test file for policy creation and deletion
 */

func init() {
	createEnterprise()
	createDomainTemplate()
	createDomain()
	createZones()
	createPolicyGroups()
	NewResourceManager()
}

func deinit() {
	deletePolicyGroups()
	deleteZones()
	deleteDomain()
	deleteDomainTemplate()
	deleteEnterprise()
}

func createPolicy(p *networkingV1.NetworkPolicy) {

}

func removePolicy(p *networkingV1.NetworkPolicy) {

}

func checkIfPolicyCreated(p *networkingV1.NetworkPolicy) {

}

func checkIfPolicyRemoved(p *networkingV1.NetworkPolicy) {

}

func TestPolicyFramework(t *testing.T) {
	init()

	tests := []*networkingV1.NetworkPolicy{
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

	for _, policy := range allPolicies {
		createPolicy(policy)
		checkIfPolicyCreated(policy)
		removePolicy(policy)
		checkIfPolicyRemoved(policy)
	}

	deinit()
}
