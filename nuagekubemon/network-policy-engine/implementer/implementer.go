package implementer

import (
	"bytes"
	"fmt"
	"github.com/golang/glog"
	"github.com/nuagenetworks/go-bambou/bambou"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/network-policy-engine/policies"
	"github.com/nuagenetworks/vspk-go/vspk"
	"gopkg.in/yaml.v2"
	"time"
)

type VSDCredentials struct {
	Username     string
	Password     string
	Organization string
	URL          string
}

type PolicyImplementer struct {
	vsdSession *bambou.Session
	vsdRoot    *vspk.Me
}

type PolicyTransaction struct {
	err *bambou.Error
}

func (implementer *PolicyImplementer) Init(vsdCredentials *VSDCredentials) error {

	glog.Infof("VSD credentials %+v", vsdCredentials)

	if vsdCredentials == nil || vsdCredentials.Username == "" ||
		vsdCredentials.Password == "" || vsdCredentials.Organization == "" ||
		vsdCredentials.URL == "" {
		return fmt.Errorf("Invalid VSD credentials %+v", vsdCredentials)
	}

	implementer.vsdSession, implementer.vsdRoot =
		vspk.NewSession(vsdCredentials.Username,
			vsdCredentials.Password,
			vsdCredentials.Organization,
			vsdCredentials.URL)

	if implementer.vsdRoot == nil || implementer.vsdSession == nil {
		return fmt.Errorf("Unable to establish session to VSD")
	}

	implementer.vsdSession.SetInsecureSkipVerify(true)
	implementer.vsdSession.Start()
	glog.Infof("Connected to vsd with credentials: Username: %s, Password: %s, Organization: %s, URL: %s", vsdCredentials.Username, vsdCredentials.Password, vsdCredentials.Organization, vsdCredentials.URL)
	return nil
}

func (implementer *PolicyImplementer) processDefaultPolicy(policy *policies.NuagePolicy) error {

	glog.Infof("Implementing default policy %+v", policy)

	d, merr := yaml.Marshal(policy)
	if merr != nil {
		glog.Errorf("Error marshaling the policy %+v", policy)
	}

	glog.Infof("Implementing default policy %s", string(d))

	policyElements, ok := policy.PolicyElements.([]policies.DefaultPolicyElement)
	if !ok {
		return fmt.Errorf("Incorrect policy type")
	}

	enterprise, err := implementer.getEnterprise(policy.Enterprise)
	if err != nil || enterprise == nil {
		return fmt.Errorf("Problem fetching the enterprise")
	}

	domain, err := implementer.getDomain(enterprise, policy.Domain)
	if err != nil || domain == nil {
		return fmt.Errorf("Problem fetching the domain")
	}

	ingressACL := vspk.NewIngressACLTemplate()
	ingressACL.Name = fmt.Sprintf("Policy-%s-ingress", policy.Name)
	ingressACL.Description = fmt.Sprintf("Ingress ACL for %s", policy.Name)
	ingressACL.DefaultAllowIP = true
	ingressACL.DefaultAllowNonIP = true
	ingressACL.Active = true
	ingressACL.ExternalID = policy.ID
	ingressACL.Priority = policy.Priority
	ingressACL.ExternalID = policy.ID

	egressACL := vspk.NewEgressACLTemplate()
	egressACL.Name = fmt.Sprintf("Policy-%s-egress", policy.Name)
	egressACL.Description = fmt.Sprintf("Egress ACL for %s", policy.Name)
	egressACL.DefaultAllowIP = true
	egressACL.DefaultAllowNonIP = true
	egressACL.Active = true
	egressACL.ExternalID = policy.ID
	egressACL.Priority = policy.Priority
	egressACL.ExternalID = policy.ID

	var ingressACLEnteries []*vspk.IngressACLEntryTemplate
	var egressACLEnteries []*vspk.EgressACLEntryTemplate

	for _, defaultPolicyElement := range policyElements {

		var toID, fromID string

		switch defaultPolicyElement.To.Type {
		case policies.POLICY_GROUP:
			toPG, err := implementer.getPolicyGroup(domain, defaultPolicyElement.To.Name)
			if err != nil || toPG == nil {
				return fmt.Errorf("Unable to get the destination PG %+v", err)
			}
			toID = toPG.ID
		default:
			panic("Not implemented")
		}

		switch defaultPolicyElement.From.Type {
		case policies.POLICY_GROUP:
			fromPG, err := implementer.getPolicyGroup(domain, defaultPolicyElement.From.Name)
			if err != nil || fromPG == nil {
				return fmt.Errorf("Unabe to get the source PG")
			}
			fromID = fromPG.ID
		default:
			panic("Not implemented")
		}

		ingressACLEntry := vspk.NewIngressACLEntryTemplate()

		ingressACLEntry.Action = string(defaultPolicyElement.Action)
		ingressACLEntry.Description = "ingress rule"

		ingressACLEntry.LocationType = string(defaultPolicyElement.From.Type)
		ingressACLEntry.LocationID = fromID

		ingressACLEntry.NetworkType = string(defaultPolicyElement.To.Type)
		ingressACLEntry.NetworkID = toID

		ingressACLEntry.Protocol = defaultPolicyElement.NetworkParameters.Protocol.String()
		ingressACLEntry.PolicyState = "LIVE"
		ingressACLEntry.Stateful = true
		ingressACLEntry.Reflexive = true
		if ingressACLEntry.Protocol != policies.ANY.String() {
			ingressACLEntry.SourcePort = defaultPolicyElement.NetworkParameters.SourcePortRange.String()
			ingressACLEntry.DestinationPort = defaultPolicyElement.NetworkParameters.DestinationPortRange.String()
		}

		egressACLEntry := vspk.NewEgressACLEntryTemplate()

		egressACLEntry.Action = string(defaultPolicyElement.Action)
		egressACLEntry.Description = "egress rule"
		egressACLEntry.LocationType = string(defaultPolicyElement.To.Type)
		egressACLEntry.LocationID = toID

		egressACLEntry.NetworkType = string(defaultPolicyElement.From.Type)
		egressACLEntry.NetworkID = fromID

		egressACLEntry.Protocol = defaultPolicyElement.NetworkParameters.Protocol.String()
		egressACLEntry.PolicyState = "LIVE"
		egressACLEntry.Stateful = true
		egressACLEntry.Reflexive = true

		if egressACLEntry.Protocol != policies.ANY.String() {

			egressACLEntry.SourcePort = defaultPolicyElement.NetworkParameters.DestinationPortRange.String()
			egressACLEntry.DestinationPort = defaultPolicyElement.NetworkParameters.SourcePortRange.String()
		}
		ingressACLEnteries = append(ingressACLEnteries, ingressACLEntry)
		egressACLEnteries = append(egressACLEnteries, egressACLEntry)

	}

	var policyTransaction PolicyTransaction
	policyTransaction.err = startPolicyTransaction(domain)

	if policyTransaction.err != nil {
		return fmt.Errorf("Unable to start policy transaction %+v", policyTransaction.err)
	}

	defer endPolicyTransaction(domain, &policyTransaction)

	policyTransaction.err = domain.CreateIngressACLTemplate(ingressACL)
	if policyTransaction.err != nil {
		return fmt.Errorf("Unable to create Ingress ACL %+v", policyTransaction.err)
	}

	policyTransaction.err = domain.CreateEgressACLTemplate(egressACL)
	if policyTransaction.err != nil {
		return fmt.Errorf("Unable to create Egress ACL %+v", policyTransaction.err)
	}

	for _, ingressACLEntry := range ingressACLEnteries {
		policyTransaction.err = ingressACL.CreateIngressACLEntryTemplate(ingressACLEntry)
		if policyTransaction.err != nil {
			glog.Errorf("Unable to create ingress ACL entry %+v %+v", policyTransaction.err, ingressACLEntry)
			return fmt.Errorf("Unable to create ingress ACL entry %+v %+v", policyTransaction.err, ingressACLEntry)
		}
	}

	for _, egressACLEntry := range egressACLEnteries {
		policyTransaction.err = egressACL.CreateEgressACLEntryTemplate(egressACLEntry)
		if policyTransaction.err != nil {
			glog.Errorf("Unable to create egress ACL entry %+v %+v", policyTransaction.err, egressACLEntry)
			return fmt.Errorf("Unable to create egress ACL entry %+v", policyTransaction.err)
		}
	}

	policyTransaction.err = nil
	return nil
}

func (implementer *PolicyImplementer) ImplementPolicy(policy *policies.NuagePolicy) error {

	if policy.Type != policies.DEFAULT {
		return fmt.Errorf("Invalid policy")
	}

	switch policy.PolicyElements.(type) {
	case []policies.DefaultPolicyElement:
		return implementer.processDefaultPolicy(policy)
	default:
		return fmt.Errorf("Invalid policy elements")
	}
	return nil
}

func (implementer *PolicyImplementer) DeletePolicy(policyID string, enterpriseName string,
	domainName string) error {

	var deleteErr bool
	var deleteErrBuf bytes.Buffer

	enterprise, err := implementer.getEnterprise(enterpriseName)
	if err != nil || enterprise == nil {
		return fmt.Errorf("Problem fetching the enterprise")
	}

	domain, err := implementer.getDomain(enterprise, domainName)
	if err != nil || domain == nil {
		return fmt.Errorf("Problem fetching the domain")
	}

	ingressACLList, err := implementer.getIngressACLList(domain, policyID)
	if err == nil && ingressACLList != nil {
		for _, ingressACLTemplate := range *ingressACLList {
			if err := ingressACLTemplate.Delete(); err != nil {
				deleteErr = true
				deleteErrBuf.WriteString(fmt.Sprintf("Error deleting ingress template %s\n",
					ingressACLTemplate.ID))
			}
		}
	}

	egressACLList, err := implementer.getEgressACLList(domain, policyID)
	if err == nil && egressACLList != nil {
		for _, egressACLTemplate := range *egressACLList {
			if err := egressACLTemplate.Delete(); err != nil {
				deleteErr = true
				deleteErrBuf.WriteString(fmt.Sprintf("Error deleting egress template %s\n",
					egressACLTemplate.ID))
			}
		}
	}

	var retErr error
	retErr = nil
	if deleteErr {
		retErr = fmt.Errorf(deleteErrBuf.String())
	}

	return retErr
}

func (implementer *PolicyImplementer) getEnterprise(enterpriseName string) (*vspk.Enterprise, *bambou.Error) {
	var enterprise *vspk.Enterprise
	enterpriseFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + enterpriseName + "\""}
	enterprises, err := implementer.vsdRoot.Enterprises(enterpriseFetchingInfo)
	if err != nil || len(enterprises) == 0 {
		return nil, err
	}

	enterprise = enterprises[0]
	return enterprise, err
}

func (implementer *PolicyImplementer) getDomain(enterprise *vspk.Enterprise, domainName string) (*vspk.Domain, *bambou.Error) {
	var domain *vspk.Domain
	domainFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + domainName + "\""}
	domains, err := enterprise.Domains(domainFetchingInfo)
	if err != nil || len(domains) == 0 {
		return nil, err
	}

	domain = domains[0]
	return domain, err
}

func (implementer *PolicyImplementer) getPolicyGroup(domain *vspk.Domain, pgName string) (*vspk.PolicyGroup, *bambou.Error) {
	var pg *vspk.PolicyGroup
	pgFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + pgName + "\""}
	pgs, err := domain.PolicyGroups(pgFetchingInfo)
	if err != nil || len(pgs) == 0 {
		return nil, err
	}
	pg = pgs[0]
	return pg, err
}

func (implementer *PolicyImplementer) getIngressACLList(domain *vspk.Domain, policyID string) (*vspk.IngressACLTemplatesList, *bambou.Error) {
	ingressACLFetchingInfo := &bambou.FetchingInfo{Filter: "externalID == \"" + policyID + "\""}
	ingressACLList, err := domain.IngressACLTemplates(ingressACLFetchingInfo)
	if err != nil || len(ingressACLList) == 0 {
		return nil, err
	}

	return &ingressACLList, nil
}

func (implementer *PolicyImplementer) getEgressACLList(domain *vspk.Domain, policyID string) (*vspk.EgressACLTemplatesList, *bambou.Error) {
	egressACLFetchingInfo := &bambou.FetchingInfo{Filter: "externalID == \"" + policyID + "\""}
	egressACLList, err := domain.EgressACLTemplates(egressACLFetchingInfo)
	if err != nil || len(egressACLList) == 0 {
		return nil, err
	}

	return &egressACLList, nil
}

func startPolicyTransaction(domain *vspk.Domain) *bambou.Error {
	job := vspk.NewJob()
	job.Command = "BEGIN_POLICY_CHANGES"
	jerr := domain.CreateJob(job)

	if jerr != nil {
		return jerr
	}

	jerr = waitForJob(job)
	if jerr != nil {
		return jerr
	}

	return nil
}

func endPolicyTransaction(domain *vspk.Domain, transaction *PolicyTransaction) *bambou.Error {
	var command string
	if transaction.err != nil {
		command = "DISCARD_POLICY_CHANGES"
	} else {
		command = "APPLY_POLICY_CHANGES"
	}

	fmt.Printf("Applying transaction %s\n", command)
	job := vspk.NewJob()
	job.Command = command
	jerr := domain.CreateJob(job)

	if jerr != nil {
		return jerr
	}

	jerr = waitForJob(job)
	if jerr != nil {
		return jerr
	}

	return nil
}

func waitForJob(job *vspk.Job) *bambou.Error {
	for {
		err := job.Fetch()
		if err != nil {
			return err
		}

		switch job.Status {
		case "SUCCESS":
			return nil
		case "FAILED":
			return bambou.NewError(500, "Error waiting for job")
		}

		time.Sleep(time.Second)
	}
}
