package implementer

import (
	"fmt"
	"github.com/nuagenetworks/go-bambou/bambou"
	"github.com/nuagenetworks/nuagepolicyapi/policies"
	"github.com/nuagenetworks/vspk-go/vspk"
)

func (implementer *PolicyImplementer) updateDefaultPolicy(policy *policies.NuagePolicy, op policies.PolicyUpdateOperation) error {

	enterprise, err := implementer.getEnterprise(policy.Enterprise)
	if err != nil || enterprise == nil {
		return fmt.Errorf("Problem fetching the enterprise")
	}

	domain, err := implementer.getDomain(enterprise, policy.Domain)
	if err != nil || domain == nil {
		return fmt.Errorf("Problem fetching the domain")
	}

	var policyTransac policyTransaction
	policyTransac.err = startpolicyTransaction(domain)

	if policyTransac.err != nil {
		return fmt.Errorf("Unable to start policy transaction %+v", policyTransac.err)
	}

	defer endpolicyTransaction(domain, &policyTransac)

	ingressACLList, err := implementer.getDraftIngressACLList(domain, policy.Name)
	if err != nil || ingressACLList == nil {
		policyTransac.err = bambou.NewError(500, "Error getting ingress ACL enteries in draft mode")
		return fmt.Errorf("Problem fetching existing ingress ACLs")
	}

	if len(*ingressACLList) != 1 {
		policyTransac.err = bambou.NewError(500, "Unexpected ingress ACL template list size")
		return fmt.Errorf("Unexpected ingress ACL template list size : %d", len(*ingressACLList))
	}

	ingressACL := (*ingressACLList)[0]
	var ingressACLEntryList []*vspk.IngressACLEntryTemplate
	if ingressACLEntryList, err = ingressACL.IngressACLEntryTemplates(&bambou.FetchingInfo{}); err != nil {
		policyTransac.err = bambou.NewError(500, "Unable to fetch existing ingress ACL enteries")
		return fmt.Errorf("Unable to fetch existing ingress ACL enteries")
	}

	egressACLList, err := implementer.getDraftEgressACLList(domain, policy.Name)
	if err != nil || egressACLList == nil {
		policyTransac.err = bambou.NewError(500, "Error getting egress ACL enteries in draft mode")
		return fmt.Errorf("Problem fetching existing egress ACLs")
	}

	if len(*egressACLList) != 1 {
		policyTransac.err = bambou.NewError(500, "Unexpected egress ACL template list size")
		return fmt.Errorf("Unexpected egress ACL template list size : %d", len(*ingressACLList))
	}

	egressACL := (*egressACLList)[0]
	var egressACLEntryList []*vspk.EgressACLEntryTemplate
	if egressACLEntryList, err = egressACL.EgressACLEntryTemplates(&bambou.FetchingInfo{}); err != nil {
		return fmt.Errorf("Unable to fetch existing egress ACL enteries")
	}

	if op == policies.UpdateAdd {
		var newIngressACLEnteries []*vspk.IngressACLEntryTemplate
		var newEgressACLEnteries []*vspk.EgressACLEntryTemplate
		var aclerr error

		if newIngressACLEnteries, newEgressACLEnteries, aclerr = implementer.getACLEntriesFromPolicy(policy); aclerr != nil {
			policyTransac.err = bambou.NewError(500, "Error parsing the policy ACL enteries")
			return fmt.Errorf("Unable to decode the ACL enteries %+v", aclerr)
		}

	IngressACLLoop:
		for _, ingressACLEntry := range newIngressACLEnteries {
			// Skip if an entry already exists
			for _, existingACLEntry := range ingressACLEntryList {
				if implementer.CompareIngressACLEntries(ingressACLEntry, existingACLEntry) == 0 {
					continue IngressACLLoop
				}
			}

			err := ingressACL.CreateIngressACLEntryTemplate(ingressACLEntry)
			if err != nil {
				policyTransac.err = bambou.NewError(500, "Error creating ingress ACL enteries")
				return fmt.Errorf("Unable to create ingress ACL entry %+v\n %+v\n %+v\n", err, ingressACL, ingressACLEntry)
			}
		}

	EgressACLLoop:
		for _, egressACLEntry := range newEgressACLEnteries {
			// Skip if an entry already exists
			for _, existingACLEntry := range egressACLEntryList {
				if implementer.CompareEgressACLEntries(egressACLEntry, existingACLEntry) == 0 {
					continue EgressACLLoop
				}
			}

			err := egressACL.CreateEgressACLEntryTemplate(egressACLEntry)
			if err != nil {
				policyTransac.err = bambou.NewError(500, "Error creating egress ACL enteries")
				return fmt.Errorf("Unable to create egress ACL entry %+v\n %+v\n %+v\n", err, egressACL, egressACLEntry)
			}
		}
	}

	if op == policies.UpdateRemove {
		var staleIngressACLEnteries []*vspk.IngressACLEntryTemplate
		var staleEgressACLEnteries []*vspk.EgressACLEntryTemplate
		var aclerr error

		if staleIngressACLEnteries, staleEgressACLEnteries, aclerr = implementer.getACLEntriesFromPolicy(policy); aclerr != nil {
			policyTransac.err = bambou.NewError(500, "Error parsing the policy ACL enteries")
			return fmt.Errorf("Unable to decode the ACL enteries")
		}

		for _, ingressACLEntry := range ingressACLEntryList {
			for _, staleACLEntry := range staleIngressACLEnteries {
				if implementer.CompareIngressACLEntries(ingressACLEntry, staleACLEntry) == 0 {
					if policyTransac.err = ingressACLEntry.Delete(); policyTransac.err != nil {
						return fmt.Errorf("Unable to purge ingress ACL entry")
					}
				}
			}
		}

		for _, egressACLEntry := range egressACLEntryList {
			for _, staleACLEntry := range staleEgressACLEnteries {
				if implementer.CompareEgressACLEntries(egressACLEntry, staleACLEntry) == 0 {
					if policyTransac.err = egressACLEntry.Delete(); policyTransac.err != nil {
						return fmt.Errorf("Unable to purge egress ACL entry")
					}
				}
			}
		}
	}

	return nil
}

// UpdatePolicy updates the default policy
func (implementer *PolicyImplementer) UpdatePolicy(policy *policies.NuagePolicy, op policies.PolicyUpdateOperation) error {
	if policy.Type != policies.Default {
		return fmt.Errorf("Invalid policy")
	}

	switch policy.PolicyElements.(type) {
	case []policies.DefaultPolicyElement:
		return implementer.updateDefaultPolicy(policy, op)
	default:
		return fmt.Errorf("Invalid policy elements")
	}

}
