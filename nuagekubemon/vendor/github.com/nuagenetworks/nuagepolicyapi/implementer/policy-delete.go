package implementer

import (
	"bytes"
	"fmt"
)

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
