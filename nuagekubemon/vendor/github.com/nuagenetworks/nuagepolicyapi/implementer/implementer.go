package implementer

import (
	"fmt"
	"github.com/nuagenetworks/go-bambou/bambou"
	"github.com/nuagenetworks/vspk-go/vspk"
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
