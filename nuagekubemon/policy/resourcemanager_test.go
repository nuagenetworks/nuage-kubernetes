package policy

import (
	"bufio"
	log "github.com/Sirupsen/logrus"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	"github.com/nuagenetworks/vspk-go/vspk"
	networkingV1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strings"
	"testing"
)

/*
* Unit test file for policy creation and deletion
*
 */

const (
	ENTERPRISE = "test-enterprise"
	DOMAIN     = "test-domain"
	URL        = "https://127.0.0.1:8443"
	USERNAME   = "*****"
	PASSWORD   = "*****"
	ORG        = "*****"
)

type objIds struct {
	enterpriseID     string
	domainTemplateID string
	domainID         string
	zoneIDs          []string
}

type testingT struct {
	t  *testing.T
	rm *ResourceManager
}

var ids objIds
var ZONES = []string{"zone1", "zone2", "zone3"}

func (tt *testingT) init() {

	vsdCallBacks := &CallBacks{
		AddPg:             tt.createPolicyGroup,
		DeletePg:          tt.deletePolicyGroup,
		AddPortsToPg:      tt.addPortsToPg,
		DeletePortsFromPg: tt.deletePortsFromPg,
	}

	k8sCallBacks := &api.ClusterClientCallBacks{
		FilterPods:       tt.getPods,
		FilterNamespaces: tt.getNamespaces,
		GetPod:           tt.getPod,
	}

	vsdMeta := &VsdMetaData{
		"enterpriseName": ENTERPRISE,
		"domainName":     DOMAIN,
		"vsdUrl":         URL,
		"username":       USERNAME,
		"password":       PASSWORD,
		"organization":   ORG,
	}

	rm, err := NewResourceManager(vsdCallBacks, k8sCallBacks, vsdMeta)
	if err != nil {
		tt.t.Fatalf("creating policy resouce manager failed %v", err)
		return
	}

	tt.rm = rm
	tt.rm.InitPolicyImplementer()
	tt.createEnterprise()
	tt.createDomainTemplate()
	tt.createDomain()
	tt.createZones()

}

func (tt *testingT) deinit() {
	tt.rm.InitPolicyImplementer()
	tt.deleteZones()
	tt.deleteDomain()
	tt.deleteDomainTemplate()
	tt.deleteEnterprise()
}

func (tt *testingT) checkIfPolicyCreated(p *api.NetworkPolicyEvent) {
	reader := bufio.NewReader(os.Stdin)
	log.Infof("check on vsd if policy is created. After that press any key and hit enter")
	reader.ReadString('\n')
	log.Infof("continuing...")
}

func (tt *testingT) checkIfPolicyRemoved(p *api.NetworkPolicyEvent) {
	reader := bufio.NewReader(os.Stdin)
	log.Infof("check on vsd if policy is deleted. After that press any key and hit enter")
	reader.ReadString('\n')
	log.Infof("continuing...")
}

func TestPolicyFramework(t *testing.T) {
	policyName := "test-np"
	policyNamespace := ZONES[0]
	policyLabels := map[string]string{
		"nuage.io/priority": "500",
	}

	allPolicies := []*networkingV1.NetworkPolicy{
		&networkingV1.NetworkPolicy{
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
							},
							{
								NamespaceSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{"zone2": "zone2"},
								},
							},
						},
					},
				},
				Egress: []networkingV1.NetworkPolicyEgressRule{
					{
						To: []networkingV1.NetworkPolicyPeer{
							{
								PodSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{"c": "d"},
								},
							},
							{
								NamespaceSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{"zone3": "zone3"},
								},
							},
						},
					},
				},
				PolicyTypes: []networkingV1.PolicyType{networkingV1.PolicyTypeIngress, networkingV1.PolicyTypeEgress},
			},
		},
	}
	/*allPolicies := []*networkingV1.NetworkPolicy{
		&networkingV1.NetworkPolicy{
			Spec: networkingV1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: map[string]string{"a": "b"},
				},
				PolicyTypes: []networkingV1.PolicyType{networkingV1.PolicyTypeIngress},
			},
		},
		&networkingV1.NetworkPolicy{
			Spec: networkingV1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: map[string]string{"a": "b"},
				},
				Ingress:     []networkingV1.NetworkPolicyIngressRule{},
				PolicyTypes: []networkingV1.PolicyType{networkingV1.PolicyTypeIngress},
			},
		},
		&networkingV1.NetworkPolicy{
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
	}*/

	tt := &testingT{
		t: t,
	}

	tt.init()

	for _, policy := range allPolicies {
		nuagePolicyEvent := &api.NetworkPolicyEvent{
			Type:      api.Added,
			Name:      policyName,
			Namespace: policyNamespace,
			Labels:    policyLabels,
			Policy:    policy.Spec,
		}

		tt.rm.HandlePolicyEvent(nuagePolicyEvent)
		tt.checkIfPolicyCreated(nuagePolicyEvent)

		nuagePolicyEvent.Type = api.Deleted
		tt.rm.HandlePolicyEvent(nuagePolicyEvent)
		tt.checkIfPolicyRemoved(nuagePolicyEvent)
	}

	tt.deinit()
}

func (tt *testingT) createEnterprise() {
	enterprise := vspk.NewEnterprise()
	enterprise.Name = ENTERPRISE
	if err := enterprise.Save(); err != nil {
		tt.t.Fatalf("creating enterprise failed with error %v", err)
	}
	ids.enterpriseID = enterprise.ID
}

func (tt *testingT) deleteEnterprise() {
	enterprise := vspk.NewEnterprise()
	enterprise.ID = ids.enterpriseID
	if err := enterprise.Delete(); err != nil {
		tt.t.Fatalf("deleting enterprise failed with error %v", err)
	}
}

func (tt *testingT) createDomainTemplate() {
	domainTemplate := vspk.NewDomainTemplate()
	domainTemplate.Name = DOMAIN
	domainTemplate.ParentID = ids.enterpriseID
	if err := domainTemplate.Save(); err != nil {
		tt.t.Fatalf("creating domain template failed with error %v", err)
	}
	ids.domainTemplateID = domainTemplate.ID
}

func (tt *testingT) deleteDomainTemplate() {
	domainTemplate := vspk.NewDomainTemplate()
	domainTemplate.ID = ids.domainTemplateID
	if err := domainTemplate.Delete(); err != nil {
		tt.t.Fatalf("deleting domain template failed with error %v", err)
	}
}

func (tt *testingT) createDomain() {
	domain := vspk.NewDomain()
	domain.Name = DOMAIN
	domain.ParentID = ids.enterpriseID
	domain.TemplateID = ids.domainTemplateID
	if err := domain.Save(); err != nil {
		tt.t.Fatalf("creating domain failed with error %v", err)
	}
	ids.domainID = domain.ID
}

func (tt *testingT) deleteDomain() {
	domain := vspk.NewDomain()
	domain.ID = ids.domainID
	if err := domain.Delete(); err != nil {
		tt.t.Fatalf("deleting domain failed with error %v", err)
	}
}

func (tt *testingT) createZones() {
	for idx, zoneName := range ZONES {
		zone := vspk.NewZone()
		zone.Name = zoneName
		zone.ParentID = ids.domainID
		if err := zone.Save(); err != nil {
			tt.t.Fatalf("creating zone(%s) failed with error %v", zoneName, err)
		}
		ids.zoneIDs[idx] = zone.ID
	}
}

func (tt *testingT) deleteZones() {
	for idx, zoneID := range ids.zoneIDs {
		zone := vspk.NewZone()
		zone.ID = zoneID
		if err := zone.Delete(); err != nil {
			tt.t.Fatalf("deleting zone(%s) failed with error %v", ZONES[idx], err)
		}
	}
}

func (tt *testingT) createPolicyGroup(name string, desc string) (string, string, error) {
	pg := vspk.NewPolicyGroup()
	pg.Name = name
	pg.Description = desc
	pg.ParentID = ids.domainID

	if err := pg.Save(); err != nil {
		tt.t.Fatalf("saving pg(%s) failed with error %v", name, err)
		return "", "", err
	}

	return name, pg.ID, nil
}

func (tt *testingT) deletePolicyGroup(id string) error {
	pg := vspk.NewPolicyGroup()
	pg.ID = id

	if err := pg.Delete(); err != nil {
		tt.t.Fatalf("deleting pg failed with error %v", err)
		return err
	}
	return nil
}

func (tt *testingT) addPortsToPg(pgId string, podsList []string) error {
	// Just a stub. Will not do anything here
	return nil
}

func (tt *testingT) deletePortsFromPg(pgId string) error {
	// Just a stub. Will not do anything here
	return nil
}

func (tt *testingT) getPods(listOpts *metav1.ListOptions, ns string) (*[]*api.PodEvent, error) {
	// Just a stub. Will not do anything here
	pods := &[]*api.PodEvent{}
	return pods, nil
}

func (tt *testingT) getPod(name string, ns string) (*api.PodEvent, error) {
	// Just a stub. Will not do anything here
	return &api.PodEvent{}, nil
}

func (tt *testingT) getNamespaces(listOpts *metav1.ListOptions) (*[]*api.NamespaceEvent, error) {
	// For testing we pass zone name within the label
	// so we return the zone name that matches with label
	namespaces := []*api.NamespaceEvent{}
	for _, zoneName := range ZONES {
		if strings.Contains(listOpts.LabelSelector, zoneName) {
			ns := &api.NamespaceEvent{
				Name: zoneName,
			}
			namespaces = append(namespaces, ns)
			break
		}
	}
	return &namespaces, nil
}
