package policy

import (
	"bufio"
	log "github.com/Sirupsen/logrus"
	"github.com/nuagenetworks/go-bambou/bambou"
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
 */

const (
	ENTERPRISE = "test-enterprise"
	DOMAIN     = "test-domain"
	URL        = "https://127.0.0.1:8443"
	USERNAME   = "******"
	PASSWORD   = "******"
	ORG        = "******"
)

type objIds struct {
	enterpriseID     string
	domainTemplateID string
	domainID         string
	zoneIDs          []string
}

type testingT struct {
	t          *testing.T
	rm         *ResourceManager
	vsdSession *bambou.Session
	vsdRoot    *vspk.Me
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

	vsdSession, vsdRoot := vspk.NewSession(USERNAME, PASSWORD, ORG, URL)
	if err := vsdSession.Start(); err != nil {
		tt.t.Fatalf("Unable to connect to Nuage VSD: " + err.Description)
	}

	tt.rm = rm
	tt.vsdSession = vsdSession
	tt.vsdRoot = vsdRoot
	tt.createEnterprise()
	tt.createDomainTemplate()
	tt.createDomain()
	tt.createZones()

}

func (tt *testingT) deinit() {
	if err := tt.rm.InitPolicyImplementer(); err != nil {
		tt.t.Fatalf("initializing policy implementer failed %v", err)
		return
	}
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
		// ingress and egress acl templates each should have two acl entries under them
		// two with policy groups. "a=b" can receive traffic from "c=d" and can send
		// traffic to "e=f"
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
						},
					},
				},
				Egress: []networkingV1.NetworkPolicyEgressRule{
					{
						To: []networkingV1.NetworkPolicyPeer{
							{
								PodSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{"e": "f"},
								},
							},
						},
					},
				},
				PolicyTypes: []networkingV1.PolicyType{networkingV1.PolicyTypeIngress, networkingV1.PolicyTypeEgress},
			},
		},
		// ingress and egress acl templates each should have two acl entries under them
		// two with zones. "a=b" can receive traffic from "zone2" and send traffic to "zone3"
		&networkingV1.NetworkPolicy{
			Spec: networkingV1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: map[string]string{"a": "b"},
				},
				Ingress: []networkingV1.NetworkPolicyIngressRule{
					{
						From: []networkingV1.NetworkPolicyPeer{
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
		// ingress and egress acl templates each should have two acl entries under them
		// one with policy group and one with zone. "a=b" can receive traffic from "c=d"
		// and can send traffic to "zone3"
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
						},
					},
				},
				Egress: []networkingV1.NetworkPolicyEgressRule{
					{
						To: []networkingV1.NetworkPolicyPeer{
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
		// ingress and egress acl templates each should have two acl entries under them
		// one with policy group and one with zone. "a=b" can receive traffic from "c=d"
		// and "zone2"
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
			},
		},
		// ingress and egress acl templates each should have two acl entries under them
		// one with policy group and one with zone. "a=b" can send traffic to "e=f" and
		// zone3
		&networkingV1.NetworkPolicy{
			Spec: networkingV1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: map[string]string{"a": "b"},
				},
				Egress: []networkingV1.NetworkPolicyEgressRule{
					{
						To: []networkingV1.NetworkPolicyPeer{
							{
								PodSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{"e": "f"},
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
		// ingress and egress acl templates each should have four acl entries under them
		// two with policy groups and two with zones
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
									MatchLabels: map[string]string{"e": "f"},
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
	if err := tt.vsdRoot.CreateEnterprise(enterprise); err != nil {
		tt.t.Fatalf("creating enterprise failed with error %v", err)
	}
	ids.enterpriseID = enterprise.ID
}

func (tt *testingT) fetchEnterprise() (*vspk.Enterprise, error) {
	enterprise := vspk.NewEnterprise()
	enterprise.ID = ids.enterpriseID
	if err := enterprise.Fetch(); err != nil {
		return nil, err
	}
	return enterprise, nil
}

func (tt *testingT) deleteEnterprise() {
	enterprise := vspk.NewEnterprise()
	enterprise.ID = ids.enterpriseID
	if err := enterprise.Delete(); err != nil {
		tt.t.Fatalf("deleting enterprise failed with error %v", err)
	}
}

func (tt *testingT) createDomainTemplate() {
	enterprise, err := tt.fetchEnterprise()
	if err != nil {
		tt.t.Fatalf("fetching enterprise failed %v", err)
	}
	domainTemplate := vspk.NewDomainTemplate()
	domainTemplate.Name = DOMAIN
	domainTemplate.ParentID = ids.enterpriseID
	if err := enterprise.CreateDomainTemplate(domainTemplate); err != nil {
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
	enterprise, err := tt.fetchEnterprise()
	if err != nil {
		tt.t.Fatalf("fetching enterprise failed %v", err)
	}
	domain := vspk.NewDomain()
	domain.Name = DOMAIN
	domain.ParentID = ids.enterpriseID
	domain.TemplateID = ids.domainTemplateID
	if err := enterprise.CreateDomain(domain); err != nil {
		tt.t.Fatalf("creating domain failed with error %v", err)
	}
	ids.domainID = domain.ID
}

func (tt *testingT) fetchDomain() (*vspk.Domain, error) {
	domain := vspk.NewDomain()
	domain.ID = ids.domainID
	if err := domain.Fetch(); err != nil {
		return nil, err
	}
	return domain, nil
}

func (tt *testingT) deleteDomain() {
	domain := vspk.NewDomain()
	domain.ID = ids.domainID
	if err := domain.Delete(); err != nil {
		tt.t.Fatalf("deleting domain failed with error %v", err)
	}
}

func (tt *testingT) createZones() {
	domain, err := tt.fetchDomain()
	if err != nil {
		tt.t.Fatalf("fetching domain failed %v", err)
	}

	for _, zoneName := range ZONES {
		zone := vspk.NewZone()
		zone.Name = zoneName
		zone.ParentID = ids.domainID
		if err := domain.CreateZone(zone); err != nil {
			tt.t.Fatalf("creating zone(%s) failed with error %v", zoneName, err)
		}
		ids.zoneIDs = append(ids.zoneIDs, zone.ID)
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
	domain, err := tt.fetchDomain()
	if err != nil {
		tt.t.Fatalf("fetching domain failed %v", err)
	}

	pg := vspk.NewPolicyGroup()
	pg.Name = name
	pg.Description = desc
	pg.ParentID = ids.domainID

	if err := domain.CreatePolicyGroup(pg); err != nil {
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
