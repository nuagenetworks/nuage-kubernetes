package implementer

import (
	"bytes"
	"fmt"
	"github.com/nuagenetworks/go-bambou/bambou"
	"github.com/nuagenetworks/nuagepolicyapi/policies"
	"github.com/nuagenetworks/vspk-go/vspk"
	"time"
)

// VSDCredentials stores information required to establish session with VSD
type VSDCredentials struct {
	Username     string
	Password     string
	Organization string
	URL          string
}

// PolicyImplementer is used to implement a Nuage Policy
type PolicyImplementer struct {
	vsdSession *bambou.Session
	vsdRoot    *vspk.Me
}

type policyTransaction struct {
	err *bambou.Error
}

// Init initializes the policy implementer and establishes session with VSD
func (implementer *PolicyImplementer) Init(vsdCredentials *VSDCredentials) error {

	if vsdCredentials == nil || vsdCredentials.Username == "" ||
		vsdCredentials.Password == "" || vsdCredentials.Organization == "" ||
		vsdCredentials.URL == "" {
		return fmt.Errorf("Invalid VSD credentials %+v", vsdCredentials)
	}

	if implementer.vsdSession != nil {
		implementer.vsdSession.Reset()
	}

	implementer.vsdSession, implementer.vsdRoot =
		vspk.NewSession(vsdCredentials.Username,
			vsdCredentials.Password,
			vsdCredentials.Organization,
			vsdCredentials.URL)

	if implementer.vsdRoot == nil || implementer.vsdSession == nil {
		return fmt.Errorf("Unable to establish session to VSD")
	}

	if err := implementer.vsdSession.SetInsecureSkipVerify(true); err != nil {
		return fmt.Errorf("Error establishing connection to vsd %+v", err)
	}

	if err := implementer.vsdSession.Start(); err != nil {
		return fmt.Errorf("Error starting vsd session %+v", err)
	}

	return nil
}

func (implementer *PolicyImplementer) processDefaultPolicy(policy *policies.NuagePolicy) error {

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

		toType, terr := policies.ConvertPolicyEndPointStringToEndPointType(string(defaultPolicyElement.To.Type))
		if terr != nil {
			return terr
		}

		fromType, ferr := policies.ConvertPolicyEndPointStringToEndPointType(string(defaultPolicyElement.From.Type))
		if ferr != nil {
			return ferr
		}

		switch toType {
		case policies.Zone:
			toZone, err := implementer.getZone(domain, defaultPolicyElement.To.Name)
			if err != nil || toZone == nil {
				return fmt.Errorf("Unable to get the destination zone %+v", err)
			}
			toID = toZone.ID
		case policies.Subnet:
			toSubnet, err := implementer.getSubnet(domain, defaultPolicyElement.To.Name)
			if err != nil || toSubnet == nil {
				return fmt.Errorf("Unable to get the destination subnet %+v", err)
			}
			toID = toSubnet.ID
		case policies.PolicyGroup:
			toPG, err := implementer.getPolicyGroup(domain, defaultPolicyElement.To.Name)
			if err != nil || toPG == nil {
				return fmt.Errorf("Unable to get the destination PG %+v", err)
			}
			toID = toPG.ID
		default:
			panic(fmt.Sprintf("Not implemented %+v", defaultPolicyElement.To))
		}

		switch fromType {
		case policies.Zone:
			fromZone, err := implementer.getZone(domain, defaultPolicyElement.From.Name)
			if err != nil || fromZone == nil {
				return fmt.Errorf("Unabe to get the source Zone %+v", err)
			}
			fromID = fromZone.ID
		case policies.Subnet:
			fromSubnet, err := implementer.getSubnet(domain, defaultPolicyElement.From.Name)
			if err != nil || fromSubnet == nil {
				return fmt.Errorf("Unabe to get the source Subnet %+v", err)
			}
			fromID = fromSubnet.ID
		case policies.PolicyGroup:
			fromPG, err := implementer.getPolicyGroup(domain, defaultPolicyElement.From.Name)
			if err != nil || fromPG == nil {
				return fmt.Errorf("Unabe to get the source PG %+v", err)
			}
			fromID = fromPG.ID

		default:
			panic(fmt.Sprintf("Not implemented %+v", defaultPolicyElement.From))
		}

		ingressACLEntry := vspk.NewIngressACLEntryTemplate()

		ingressACLEntry.Action = policies.ConvertPolicyActionToNuageAction(defaultPolicyElement.Action)
		ingressACLEntry.Description = "ingress rule"
		ingressACLEntry.LocationType = string(fromType)
		ingressACLEntry.LocationID = fromID
		ingressACLEntry.NetworkType = string(toType)
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
		egressACLEntry.Action = policies.ConvertPolicyActionToNuageAction(defaultPolicyElement.Action)
		egressACLEntry.Description = "egress rule"
		egressACLEntry.LocationType = string(toType)
		egressACLEntry.LocationID = toID
		egressACLEntry.NetworkType = string(fromType)
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

	var policyTransac policyTransaction
	policyTransac.err = startpolicyTransaction(domain)

	if policyTransac.err != nil {
		return fmt.Errorf("Unable to start policy transaction %+v", policyTransac.err)
	}

	defer endpolicyTransaction(domain, &policyTransac)

	policyTransac.err = domain.CreateIngressACLTemplate(ingressACL)
	if policyTransac.err != nil {
		return fmt.Errorf("Unable to create Ingress ACL %+v", policyTransac.err)
	}

	policyTransac.err = domain.CreateEgressACLTemplate(egressACL)
	if policyTransac.err != nil {
		return fmt.Errorf("Unable to create Egress ACL %+v", policyTransac.err)
	}

	for _, ingressACLEntry := range ingressACLEnteries {
		policyTransac.err = ingressACL.CreateIngressACLEntryTemplate(ingressACLEntry)
		if policyTransac.err != nil {
			return fmt.Errorf("Unable to create ingress ACL entry %+v %+v", policyTransac.err, ingressACLEntry)
		}
	}

	for _, egressACLEntry := range egressACLEnteries {
		policyTransac.err = egressACL.CreateEgressACLEntryTemplate(egressACLEntry)
		if policyTransac.err != nil {
			return fmt.Errorf("Unable to create egress ACL entry %+v", policyTransac.err)
		}
	}

	policyTransac.err = nil
	return nil
}

// ImplementPolicy implements a Nuage policy
func (implementer *PolicyImplementer) ImplementPolicy(policy *policies.NuagePolicy) error {

	if policy.Type != policies.Default {
		return fmt.Errorf("Invalid policy")
	}

	switch policy.PolicyElements.(type) {
	case []policies.DefaultPolicyElement:
		return implementer.processDefaultPolicy(policy)
	default:
		return fmt.Errorf("Invalid policy elements")
	}
}

// DeletePolicy deletes a Nuage policy using the policy ID
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
			if err = ingressACLTemplate.Delete(); err != nil {
				deleteErr = true
				if _, berr := deleteErrBuf.WriteString(fmt.Sprintf("Error deleting ingress template %s\n",
					ingressACLTemplate.ID)); berr != nil {
				}
			}
		}
	}

	egressACLList, err := implementer.getEgressACLList(domain, policyID)
	if err == nil && egressACLList != nil {
		for _, egressACLTemplate := range *egressACLList {
			if err = egressACLTemplate.Delete(); err != nil {
				deleteErr = true
				if _, berr := deleteErrBuf.WriteString(fmt.Sprintf("Error deleting egress template %s\n",
					egressACLTemplate.ID)); berr != nil {
				}
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

func (implementer *PolicyImplementer) getZone(domain *vspk.Domain, zoneName string) (*vspk.Zone, *bambou.Error) {
	var zone *vspk.Zone
	zoneFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + zoneName + "\""}
	zones, err := domain.Zones(zoneFetchingInfo)
	if err != nil || len(zones) == 0 {
		return nil, err
	}
	zone = zones[0]
	return zone, err
}

func (implementer *PolicyImplementer) getSubnet(domain *vspk.Domain, subnetName string) (*vspk.Subnet, *bambou.Error) {
	var subnet *vspk.Subnet
	subnetFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + subnetName + "\""}
	subnets, err := domain.Subnets(subnetFetchingInfo)
	if err != nil || len(subnets) == 0 {
		return nil, err
	}
	subnet = subnets[0]
	return subnet, err
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

func startpolicyTransaction(domain *vspk.Domain) *bambou.Error {
	job := vspk.NewJob()
	job.Command = "BEGIN_POLICY_CHANGES"
	jerr := domain.CreateJob(job)

	if jerr != nil {
		return jerr
	}

	return waitForJob(job)
}

func endpolicyTransaction(domain *vspk.Domain, transaction *policyTransaction) {
	var command string
	if transaction.err != nil {
		command = "DISCARD_POLICY_CHANGES"
	} else {
		command = "APPLY_POLICY_CHANGES"
	}

	job := vspk.NewJob()
	job.Command = command
	domain.CreateJob(job)
	waitForJob(job)
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
