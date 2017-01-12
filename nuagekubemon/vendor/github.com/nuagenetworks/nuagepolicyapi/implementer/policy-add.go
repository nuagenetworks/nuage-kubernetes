package implementer

import (
	"fmt"
	"github.com/nuagenetworks/nuagepolicyapi/policies"
	"github.com/nuagenetworks/vspk-go/vspk"
)

func (implementer *PolicyImplementer) processDefaultPolicy(policy *policies.NuagePolicy) error {

	enterprise, err := implementer.getEnterprise(policy.Enterprise)
	if err != nil || enterprise == nil {
		return fmt.Errorf("Problem fetching the enterprise")
	}

	domain, err := implementer.getDomain(enterprise, policy.Domain)
	if err != nil || domain == nil {
		return fmt.Errorf("Problem fetching the domain")
	}

	ingressACL, egressACL := createACLTemplates(policy)

	var ingressACLEnteries []*vspk.IngressACLEntryTemplate
	var egressACLEnteries []*vspk.EgressACLEntryTemplate
	var aclerr error
	if ingressACLEnteries, egressACLEnteries, aclerr = implementer.getACLEntriesFromPolicy(policy); aclerr != nil {
		return err
	}

	var policyTransac policyTransaction
	policyTransac.err = startpolicyTransaction(domain)

	if policyTransac.err != nil {
		return fmt.Errorf("Unable to start policy transaction %+v", policyTransac.err)
	}

	defer endpolicyTransaction(domain, &policyTransac)

	policyTransac.err = domain.CreateIngressACLTemplate(ingressACL)
	if policyTransac.err != nil {
		return fmt.Errorf("Unable to create Ingress ACL template %+v", policyTransac.err)
	}

	policyTransac.err = domain.CreateEgressACLTemplate(egressACL)
	if policyTransac.err != nil {
		return fmt.Errorf("Unable to create Egress ACL template %+v", policyTransac.err)
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
			return fmt.Errorf("Unable to create egress ACL entry %+v %+v", policyTransac.err, egressACLEntry)
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
