package implementer

import (
	"fmt"
	"time"

	"github.com/nuagenetworks/go-bambou/bambou"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/policyapi/policies"
	"github.com/nuagenetworks/vspk-go/vspk"
)

func newACLTemplates(policy *policies.NuagePolicy) (*vspk.IngressACLTemplate, *vspk.EgressACLTemplate) {

	ingressACL := vspk.NewIngressACLTemplate()
	ingressACL.Name = policy.Name
	ingressACL.Description = fmt.Sprintf("Ingress ACL for %s", policy.Name)
	ingressACL.DefaultAllowIP = true
	ingressACL.DefaultAllowNonIP = true
	ingressACL.Active = true
	ingressACL.Priority = policy.Priority

	egressACL := vspk.NewEgressACLTemplate()
	egressACL.Name = policy.Name
	egressACL.Description = fmt.Sprintf("Egress ACL for %s", policy.Name)
	egressACL.DefaultAllowIP = true
	egressACL.DefaultAllowNonIP = true
	egressACL.Active = true
	egressACL.Priority = policy.Priority

	return ingressACL, egressACL
}

func (implementer *PolicyImplementer) getACLEntriesFromPolicy(policy *policies.NuagePolicy) ([]*vspk.IngressACLEntryTemplate,
	[]*vspk.EgressACLEntryTemplate, error) {

	var ingressACLEnteries []*vspk.IngressACLEntryTemplate
	var egressACLEnteries []*vspk.EgressACLEntryTemplate

	policyElements, ok := policy.PolicyElements.([]policies.DefaultPolicyElement)
	if !ok {
		return ingressACLEnteries, egressACLEnteries, fmt.Errorf("Incorrect policy type")
	}

	enterprise, err := implementer.getEnterprise(policy.Enterprise)
	if err != nil || enterprise == nil {
		return ingressACLEnteries, egressACLEnteries, fmt.Errorf("Problem fetching the enterprise")
	}

	domain, err := implementer.getDomain(enterprise, policy.Domain)
	if err != nil || domain == nil {
		return ingressACLEnteries, egressACLEnteries, fmt.Errorf("Problem fetching the domain")
	}

	for _, defaultPolicyElement := range policyElements {

		var toID, fromID string

		toType, terr := policies.ConvertPolicyEndPointStringToEndPointType(string(defaultPolicyElement.To.Type))
		if terr != nil {
			return ingressACLEnteries, egressACLEnteries, terr
		}

		fromType, ferr := policies.ConvertPolicyEndPointStringToEndPointType(string(defaultPolicyElement.From.Type))
		if ferr != nil {
			return ingressACLEnteries, egressACLEnteries, ferr
		}

		switch toType {
		case policies.EndPointZone:
		case policies.Zone:
			toZone, err := implementer.getZone(domain, defaultPolicyElement.To.Name)
			if err != nil || toZone == nil {
				return ingressACLEnteries, egressACLEnteries, fmt.Errorf("Unable to get the destination zone %+v", err)
			}
			toID = toZone.ID
		case policies.Subnet:
			toSubnet, err := implementer.getSubnet(domain, defaultPolicyElement.To.Name)
			if err != nil || toSubnet == nil {
				return ingressACLEnteries, egressACLEnteries, fmt.Errorf("Unable to get the destination subnet %+v", err)
			}
			toID = toSubnet.ID
		case policies.PolicyGroup:
			toPG, err := implementer.getPolicyGroup(domain, defaultPolicyElement.To.Name)
			if err != nil || toPG == nil {
				return ingressACLEnteries, egressACLEnteries, fmt.Errorf("Unable to get the destination PG %+v", err)
			}
			toID = toPG.ID
		case policies.NetworkMacro:
			toNWMacro, err := implementer.getEnterpriseNetworkMacro(enterprise, defaultPolicyElement.To.Name)
			if err != nil || toNWMacro == nil {
				return ingressACLEnteries, egressACLEnteries, fmt.Errorf("Unable to get the destination networkmacro %+v", err)
			}
			toID = toNWMacro.ID
		default:
			panic(fmt.Sprintf("Not implemented %+v", defaultPolicyElement.To))
		}

		switch fromType {
		case policies.EndPointZone:
		case policies.Zone:
			fromZone, err := implementer.getZone(domain, defaultPolicyElement.From.Name)
			if err != nil || fromZone == nil {
				return ingressACLEnteries, egressACLEnteries, fmt.Errorf("Unabe to get the source Zone %+v", err)
			}
			fromID = fromZone.ID
		case policies.Subnet:
			fromSubnet, err := implementer.getSubnet(domain, defaultPolicyElement.From.Name)
			if err != nil || fromSubnet == nil {
				return ingressACLEnteries, egressACLEnteries, fmt.Errorf("Unabe to get the source Subnet %+v", err)
			}
			fromID = fromSubnet.ID
		case policies.PolicyGroup:
			fromPG, err := implementer.getPolicyGroup(domain, defaultPolicyElement.From.Name)
			if err != nil || fromPG == nil {
				return ingressACLEnteries, egressACLEnteries, fmt.Errorf("Unabe to get the source PG %+v", err)
			}
			fromID = fromPG.ID
		case policies.NetworkMacro:
			fromNWMacro, err := implementer.getEnterpriseNetworkMacro(enterprise, defaultPolicyElement.From.Name)
			if err != nil || fromNWMacro == nil {
				return ingressACLEnteries, egressACLEnteries, fmt.Errorf("Unable to get the source networkmacro %+v", err)
			}
			fromID = fromNWMacro.ID
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
		if defaultPolicyElement.Action == policies.Allow {
			ingressACLEntry.Stateful = true
		}
		ingressACLEntry.SourcePort = defaultPolicyElement.NetworkParameters.SourcePortRange.String()
		ingressACLEntry.DestinationPort = defaultPolicyElement.NetworkParameters.DestinationPortRange.String()

		egressACLEntry := vspk.NewEgressACLEntryTemplate()
		egressACLEntry.Action = policies.ConvertPolicyActionToNuageAction(defaultPolicyElement.Action)
		egressACLEntry.Description = "egress rule"
		egressACLEntry.LocationType = string(toType)
		egressACLEntry.LocationID = toID
		egressACLEntry.NetworkType = string(fromType)
		egressACLEntry.NetworkID = fromID
		egressACLEntry.Protocol = defaultPolicyElement.NetworkParameters.Protocol.String()
		egressACLEntry.PolicyState = "LIVE"
		if defaultPolicyElement.Action == policies.Allow {
			egressACLEntry.Stateful = true
		}
		egressACLEntry.SourcePort = defaultPolicyElement.NetworkParameters.SourcePortRange.String()
		egressACLEntry.DestinationPort = defaultPolicyElement.NetworkParameters.DestinationPortRange.String()

		ingressACLEnteries = append(ingressACLEnteries, ingressACLEntry)
		egressACLEnteries = append(egressACLEnteries, egressACLEntry)
	}

	return ingressACLEnteries, egressACLEnteries, nil
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

// CompareIngressACLEntries returns 0 is the two ingress ACLs are the same
func (implementer *PolicyImplementer) CompareIngressACLEntries(acl1 *vspk.IngressACLEntryTemplate, acl2 *vspk.IngressACLEntryTemplate) int {
	if acl1.NetworkID == acl2.NetworkID &&
		acl1.LocationID == acl2.LocationID &&
		acl1.NetworkType == acl2.NetworkType &&
		acl1.LocationType == acl2.LocationType &&
		acl1.Protocol == acl2.Protocol &&
		acl1.SourcePort == acl2.SourcePort &&
		acl1.DestinationPort == acl2.DestinationPort {
		return 0
	}
	return -1
}

// CompareEgressACLEntries returns 0 is the two ingress ACLs are the same
func (implementer *PolicyImplementer) CompareEgressACLEntries(acl1 *vspk.EgressACLEntryTemplate, acl2 *vspk.EgressACLEntryTemplate) int {
	if acl1.NetworkID == acl2.NetworkID &&
		acl1.LocationID == acl2.LocationID &&
		acl1.NetworkType == acl2.NetworkType &&
		acl1.LocationType == acl2.LocationType &&
		acl1.Protocol == acl2.Protocol &&
		acl1.SourcePort == acl2.SourcePort &&
		acl1.DestinationPort == acl2.DestinationPort {
		return 0
	}
	return -1
}

func (implementer *PolicyImplementer) getEnterpriseNetworkMacro(enterprise *vspk.Enterprise, macroName string) (*vspk.EnterpriseNetwork, *bambou.Error) {
	var enterpriseNetwork *vspk.EnterpriseNetwork
	enterpriseNetworkFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + macroName + "\""}
	enterpriseNetworks, err := enterprise.EnterpriseNetworks(enterpriseNetworkFetchingInfo)
	if err != nil || len(enterpriseNetworks) == 0 {
		return nil, err
	}

	enterpriseNetwork = enterpriseNetworks[0]
	return enterpriseNetwork, err
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

func (implementer *PolicyImplementer) getIngressACLList(domain *vspk.Domain, policyName string) (*vspk.IngressACLTemplatesList, *bambou.Error) {
	ingressACLFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + policyName + "\""}
	ingressACLList, err := domain.IngressACLTemplates(ingressACLFetchingInfo)
	if err != nil || len(ingressACLList) == 0 {
		return nil, err
	}

	return &ingressACLList, nil
}

func (implementer *PolicyImplementer) getEgressACLList(domain *vspk.Domain, policyName string) (*vspk.EgressACLTemplatesList, *bambou.Error) {
	egressACLFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + policyName + "\""}
	egressACLList, err := domain.EgressACLTemplates(egressACLFetchingInfo)
	if err != nil || len(egressACLList) == 0 {
		return nil, err
	}

	return &egressACLList, nil
}

func (implementer *PolicyImplementer) getDraftIngressACLList(domain *vspk.Domain, policyName string) (*vspk.IngressACLTemplatesList, *bambou.Error) {
	ingressACLFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + policyName + "\" AND policyState == \"DRAFT\""}
	ingressACLList, err := domain.IngressACLTemplates(ingressACLFetchingInfo)
	if err != nil || len(ingressACLList) == 0 {
		return nil, err
	}

	return &ingressACLList, nil
}

func (implementer *PolicyImplementer) getDraftEgressACLList(domain *vspk.Domain, policyName string) (*vspk.EgressACLTemplatesList, *bambou.Error) {
	egressACLFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + policyName + "\" AND policyState == \"DRAFT\""}
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

	jerr := domain.CreateJob(job)
	if jerr != nil {
		return
	}

	jerr = waitForJob(job)
	if jerr != nil {
		return
	}
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
