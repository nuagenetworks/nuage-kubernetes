package implementer

import (
	"crypto/tls"
	"fmt"
	"github.com/nuagenetworks/go-bambou/bambou"
	"github.com/nuagenetworks/vspk-go/vspk"
)

// VSDCredentials stores information required to establish session with VSD
type VSDCredentials struct {
	URL          string
	UserCertFile string
	UserKeyFile  string
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

	if vsdCredentials == nil || vsdCredentials.URL == "" {
		return fmt.Errorf("Invalid VSD credentials %+v", vsdCredentials)
	}

	if implementer.vsdSession != nil {
		implementer.vsdSession.Reset()
	}

	cert, err := tls.LoadX509KeyPair(vsdCredentials.UserCertFile, vsdCredentials.UserKeyFile)
	if err != nil {
		return fmt.Errorf("Loading TLS certificate and private key failed")
	}
	implementer.vsdSession, implementer.vsdRoot = vspk.NewX509Session(&cert, vsdCredentials.URL)

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
