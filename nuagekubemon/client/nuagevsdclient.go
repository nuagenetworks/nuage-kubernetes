/*
###########################################################################
#
#   Filename:           nuagevsdclient.go
#
#   Author:             Aniket Bhat
#   Created:            July 20, 2015
#
#   Description:        NuageKubeMon Vsd Client Interface
#
###########################################################################
#
#              Copyright (c) 2015 Nuage Networks
#
###########################################################################
*/

package client

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/jmcvetta/napping"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/api"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/config"
	"io/ioutil"
	"net/http"
	"strings"
)

type NuageVsdClient struct {
	url          string
	version      string
	username     string
	password     string
	enterprise   string
	session      napping.Session
	enterpriseID string
	domainID     string
	zones        map[string]string      //project name -> zone id mapping
	subnets      map[string]*SubnetList //zone id -> list of subnets mapping
	pool         IPv4SubnetPool
	subnetSize   int //the size in bits of the subnets we allocate (i.e. size 8 produces /24 subnets).
}

type SubnetList struct {
	SubnetID string
	Subnet   *IPv4Subnet
	Next     *SubnetList
}

const clusterEnterpriseName = "Openshift-Enterprise"
const clusterDomainTemplateName = "Openshift-Domain-Template"
const clusterDomainName = "Openshift-Domain"

func NewNuageVsdClient(nkmConfig *config.NuageKubeMonConfig) *NuageVsdClient {
	nvsdc := new(NuageVsdClient)
	nvsdc.Init(nkmConfig)
	return nvsdc
}

func (nvsdc *NuageVsdClient) GetAuthorizationToken() error {
	h := nvsdc.session.Header
	h.Add("X-Nuage-Organization", nvsdc.enterprise)
	h.Add("Authorization", "XREST "+base64.StdEncoding.EncodeToString([]byte(nvsdc.username+":"+nvsdc.password)))
	var result [1]api.VsdAuthToken
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"me", nil, &result, &e)
	if err != nil {
		glog.Error("Error when requesting authorization token", err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status())
	if resp.Status() == 200 {
		h.Set("Authorization", "XREST "+base64.StdEncoding.EncodeToString([]byte(nvsdc.username+":"+result[0].APIKey)))
		return nil
	} else {
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) CreateEnterprise(enterpriseName string) (string, error) {
	payload := api.VsdEnterprise{
		Name:        enterpriseName,
		Description: "Auto-generated enterprise for Openshift Cluster",
	}
	result := make([]api.VsdEnterprise, 1)
	e := api.RESTError{}
	resp, err := nvsdc.session.Post(nvsdc.url+"enterprises", &payload, &result, &e)
	if err != nil {
		glog.Error("Error when creating enterprise", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when creating the enterprise")
	switch resp.Status() {
	case 201:
		glog.Infoln("Created the enterprise: ", result[0].ID)
		return result[0].ID, nil
	case 409:
		glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
		//Enterprise already exists, call Get to retrieve the ID
		id, err := nvsdc.GetEnterpriseID(enterpriseName)
		if err != nil {
			glog.Errorf("Error when getting enterprise ID: %s", err)
			return "", err
		} else {
			return id, nil
		}
	default:
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return "", errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) CreateAdminUser(enterpriseID, user, password string) (string, error) {
	passwd := fmt.Sprintf("%x", sha1.Sum([]byte(password)))
	payload := api.VsdUser{
		UserName:  user,
		Password:  passwd,
		FirstName: "Admin",
		LastName:  "Admin",
		Email:     "admin@localhost",
	}
	result := make([]api.VsdUser, 1)
	e := api.RESTError{}
	//Get admin ID after creating the admin user
	var adminId string
	resp, err := nvsdc.session.Post(nvsdc.url+"enterprises/"+enterpriseID+"/users", &payload, &result, &e)
	if err != nil {
		glog.Error("Error when creating admin user", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when creating the admin user")
	switch resp.Status() {
	case 201:
		glog.Infoln("Created the admin user: ", result[0].ID)
		adminId = result[0].ID
	case 409:
		//Enterprise already exists, call Get to retrieve the ID
		id, err := nvsdc.GetAdminID(enterpriseID, "admin")
		if err != nil {
			glog.Errorf("Error when getting admin user's ID: %s", err)
		} else {
			adminId = id
		}
	default:
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return "", errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
	}
	//Get admin group ID and add the admin id to the admin group
	groupId, err := nvsdc.GetAdminGroupID(enterpriseID)
	if err != nil {
		glog.Errorf("Error when getting admin group ID: %s", err)
		return "", err
	}
	groupPayload := []string{adminId}
	e = api.RESTError{}
	resp, err = nvsdc.session.Put(nvsdc.url+"groups/"+groupId+"/users", &groupPayload, nil, &e)
	if err != nil {
		glog.Error("Error when adding admin user to the admin group", err)
		return "", err
	} else {
		glog.Infoln("Got a reponse status", resp.Status(), "when adding user to the admin group")
		switch resp.Status() {
		case 204:
			glog.Infoln("Added the admin user to the admin group")
		case 409:
			glog.Infoln("Admin user already in admin group")
		default:
			glog.Errorln("Bad response status from VSD Server")
			glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
			glog.Errorf("\t Status:  %v\n", resp.Status())
			glog.Errorf("\t Message: %v\n", e.Message)
			glog.Errorf("\t Errors: %v\n", e.Message)
			return "", errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
		}
	}
	return adminId, nil
}

func (nvsdc *NuageVsdClient) GetAdminID(enterpriseID, name string) (string, error) {
	result := make([]api.VsdUser, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `userName == "`+name+`"`)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"enterprises/"+enterpriseID+"/users", nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting admin user ID %s", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting user ID")
	if resp.Status() == 200 {
		// Status code 200 is returned even if there's no results.  If
		// the filter didn't match anything (or there was nothing to
		// return), the result object will just be empty.
		if result[0].UserName == name {
			return result[0].ID, nil
		} else if result[0].UserName == "" {
			return "", errors.New("User not found")
		} else {
			return "", errors.New(fmt.Sprintf(
				"Found %q instead of %q", result[0].UserName, name))
		}
	} else {
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return "", errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) GetAdminGroupID(enterpriseID string) (string, error) {
	result := make([]api.VsdGroup, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `role == "ORGADMIN"`)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"enterprises/"+enterpriseID+"/groups", nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting admin group ID %s", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting ID of group ORGADMIN")
	if resp.Status() == 200 {
		// Status code 200 is returned even if there's no results.  If
		// the filter didn't match anything (or there was nothing to
		// return), the result object will just be empty.
		if result[0].Role == "ORGADMIN" {
			return result[0].ID, nil
		} else if result[0].ID == "" {
			return "", errors.New("Admin Group not found")
		} else {
			return "", errors.New(fmt.Sprintf(
				"Found %q instead of \"ORGADMIN\"", result[0].Role))
		}
	} else {
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return "", errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) GetEnterpriseID(name string) (string, error) {
	result := make([]api.VsdObject, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `name == "`+name+`"`)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"enterprises", nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting enterprise ID %s", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting enterprise ID")
	if resp.Status() == 200 {
		// Status code 200 is returned even if there's no results.  If
		// the filter didn't match anything (or there was nothing to
		// return), the result object will just be empty.
		if result[0].Name == name {
			return result[0].ID, nil
		} else if result[0].Name == "" {
			return "", errors.New("Enterprise not found")
		} else {
			return "", errors.New(fmt.Sprintf(
				"Found %q instead of %q", result[0].Name, name))
		}
	} else {
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return "", errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) CreateSession() {
	nvsdc.username = "csproot"
	nvsdc.password = "csproot"
	nvsdc.enterprise = "csp"
	nvsdc.session = napping.Session{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		Header: &http.Header{},
	}
	nvsdc.session.Header.Add("Content-Type", "application/json")
}

func (nvsdc *NuageVsdClient) LoginAsAdmin(user, password, enterpriseName string) error {
	nvsdc.username = user
	nvsdc.password = password
	nvsdc.enterprise = enterpriseName
	h := nvsdc.session.Header
	h.Del("X-Nuage-Organization")
	h.Del("Authorization")
	return nvsdc.GetAuthorizationToken()
}

func (nvsdc *NuageVsdClient) Init(nkmConfig *config.NuageKubeMonConfig) {
	nvsdc.version = nkmConfig.NuageVspVersion
	nvsdc.url = nkmConfig.NuageVsdApiUrl + "/nuage/api/" + nvsdc.version + "/"
	ipPool, err := IPv4SubnetFromString(nkmConfig.OsMasterConfig.NetworkConfig.ClusterCIDR)
	if err != nil {
		glog.Fatalf("Failure in init: %s\n", err)
	}
	nvsdc.subnetSize = nkmConfig.OsMasterConfig.NetworkConfig.SubnetLength
	if nvsdc.subnetSize < 0 || nvsdc.subnetSize > 32 {
		glog.Errorf("Invalid hostSubnetLength of %d.  Using default value of 8",
			nvsdc.subnetSize)
		nvsdc.subnetSize = 8
	}
	if nvsdc.subnetSize > (32 - ipPool.CIDRMask) {
		// If the size of the subnet (in bits) is larger than the total pool
		// size (in bits), we can't even allocate 1 subnet.  Default to using
		// half the remaining bits per subnet, rounded down (/24 has 8 bits
		// remaining, so use 4 bits per subnet).
		newSize := (32 - ipPool.CIDRMask) / 2
		glog.Fatalf("Cannot allocate %d bit subnets from %s.  Using %d bits per subnet.",
			nvsdc.subnetSize, ipPool.String(), newSize)
		nvsdc.subnetSize = newSize
	}
	// A null IPv4SubnetPool acts like all addresses are allocated, so we can
	// initialize it to have the available cluster address space by just
	// Free()-ing it.
	nvsdc.pool.Free(ipPool)
	nvsdc.zones = make(map[string]string)
	nvsdc.subnets = make(map[string]*SubnetList)
	nvsdc.CreateSession()
	err = nvsdc.GetAuthorizationToken()
	if err != nil {
		glog.Fatal(err)
	}
	nvsdc.enterpriseID, err = nvsdc.CreateEnterprise(clusterEnterpriseName)
	if err != nil {
		glog.Fatal(err)
	}
	_, err = nvsdc.CreateAdminUser(nvsdc.enterpriseID, "admin", "admin")
	if err != nil {
		glog.Fatal(err)
	}
	err = nvsdc.InstallLicense(nkmConfig.LicenseFile)
	if err != nil {
		glog.Fatal(err)
	}
	err = nvsdc.LoginAsAdmin("admin", "admin", clusterEnterpriseName)
	if err != nil {
		glog.Fatal(err)
	}
	domainTemplateID, err := nvsdc.CreateDomainTemplate(nvsdc.enterpriseID,
		clusterDomainTemplateName)
	if err != nil {
		glog.Fatal(err)
	}
	err = nvsdc.ApplyAclTemplates(domainTemplateID)
	if err != nil {
		glog.Fatal(err)
	}
	nvsdc.domainID, err = nvsdc.CreateDomain(nvsdc.enterpriseID,
		domainTemplateID, clusterDomainName)
	if err != nil {
		glog.Fatal(err)
	}
}

func (nvsdc *NuageVsdClient) InstallLicense(licensePath string) error {
	if licensePath == "" {
		glog.Error("No license file specified")
		//check if a license already exists.
		// if it does then its not an error
		return nvsdc.GetLicense()
	}
	//try installing the license file
	license, err := ioutil.ReadFile(licensePath)
	if err != nil {
		glog.Error("Failed to read license file", err)
		return err
	}
	licenseString := strings.TrimSpace(string(license))
	payload := api.VsdLicense{
		License: licenseString,
	}
	result := make([]api.VsdLicense, 1)
	e := api.RESTError{}
	glog.Info("Attempting to install license file", licensePath)
	resp, err := nvsdc.session.Post(nvsdc.url+"licenses", &payload, &result, &e)
	if err != nil {
		glog.Error("Error when installing license", err)
		return err
	}
	glog.Infoln("License Install: reponse status", resp.Status())
	switch resp.Status() {
	case 201:
		glog.Infoln("Installed the license: ", result[0].LicenseId)
	case 409:
		//TODO: license already exists, call Get to retrieve the ID? Do we need to delete the existing license?
		glog.Info("License already exists")
	default:
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
	}
	return nil
}

func (nvsdc *NuageVsdClient) GetLicense() error {
	result := make([]api.VsdLicense, 1)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"licenses", nil, &result, &e)
	if err != nil {
		glog.Error("Error when requesting license", err)
		return err
	}
	glog.Infoln("GetLicense() got a reponse status", resp.Status())
	if resp.Status() == 200 {
		return nil
	} else {
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) CreateDomainTemplate(enterpriseID, domainTemplateName string) (string, error) {
	result := make([]api.VsdObject, 1)
	payload := api.VsdObject{
		Name:        domainTemplateName,
		Description: "Auto-generated default domain template",
	}
	e := api.RESTError{}
	resp, err := nvsdc.session.Post(nvsdc.url+"enterprises/"+enterpriseID+"/domaintemplates", &payload, &result, &e)
	if err != nil {
		glog.Error("Error when creating domain template", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when creating domain template")
	switch resp.Status() {
	case 201:
		glog.Infoln("Created the domain: ", result[0].ID)
		return result[0].ID, nil
	case 409:
		//Enterprise already exists, call Get to retrieve the ID
		id, err := nvsdc.GetDomainTemplateID(enterpriseID, domainTemplateName)
		if err != nil {
			glog.Errorf("Error when getting domain template ID: %s", err)
			return "", err
		}
		return id, nil
	default:
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return "", errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) GetDomainTemplateID(enterpriseID, name string) (string, error) {
	result := make([]api.VsdObject, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `name == "`+name+`"`)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"enterprises/"+enterpriseID+"/domaintemplates", nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting domain template ID %s", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting domain template ID")
	if resp.Status() == 200 {
		// Status code 200 is returned even if there's no results.  If
		// the filter didn't match anything (or there was nothing to
		// return), the result object will just be empty.
		if result[0].Name == name {
			return result[0].ID, nil
		} else if result[0].Name == "" {
			return "", errors.New("Domain Template not found")
		} else {
			return "", errors.New(fmt.Sprintf(
				"Found %q instead of %q", result[0].Name, name))
		}
	} else {
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return "", errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) ApplyAclTemplates(domainTemplateID string) error {
	result := make([]api.VsdObject, 1)
	payload := api.VsdAclTemplate{
		Name:              "Auto-generated Ingress Policies",
		DefaultAllowIP:    true,
		DefaultAllowNonIP: true,
	}
	e := api.RESTError{}
	resp, err := nvsdc.session.Post(
		nvsdc.url+"domaintemplates/"+domainTemplateID+"/ingressacltemplates",
		&payload, &result, &e)
	if err != nil {
		glog.Error("Error when applying ingress acl template", err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status(),
		"when creating ingress acl template")
	switch resp.Status() {
	case 201:
		fallthrough
	case 409:
		glog.Infoln("Applied default ingress ACL")
	default:
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
	}
	// Change the name of the payload to represent that the next policy
	// template is an egress one, but otherwise reuse the existing template
	// definition
	payload.Name = "Auto-generated Egress Policies"
	resp, err = nvsdc.session.Post(
		nvsdc.url+"domaintemplates/"+domainTemplateID+"/egressacltemplates",
		&payload, &result, &e)
	if err != nil {
		glog.Error("Error when applying egress acl template", err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status(),
		"when creating egress acl template")
	switch resp.Status() {
	case 201:
		fallthrough
	case 409:
		glog.Infoln("Applied default egress ACL")
	default:
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
	}
	return nil
}

func (nvsdc *NuageVsdClient) GetZoneID(domainID, name string) (string, error) {
	result := make([]api.VsdObject, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `name == "`+name+`"`)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"domains/"+domainID+"/zones", nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting zone ID %s", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting zone ID")
	if resp.Status() == 200 {
		// Status code 200 is returned even if there's no results.  If
		// the filter didn't match anything (or there was nothing to
		// return), the result object will just be empty.
		if result[0].Name == name {
			return result[0].ID, nil
		} else if result[0].Name == "" {
			return "", errors.New("Zone not found")
		} else {
			return "", errors.New(fmt.Sprintf(
				"Found %q instead of %q", result[0].Name, name))
		}
	} else {
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return "", errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) CreateDomain(enterpriseID, domainTemplateID, name string) (string, error) {
	result := make([]api.VsdDomain, 1)
	payload := api.VsdDomain{
		Name:        name,
		Description: "Auto-generated for OpenShift containers",
		TemplateID:  domainTemplateID,
		PATEnabled:  api.PATEnabled,
	}
	e := api.RESTError{}
	resp, err := nvsdc.session.Post(nvsdc.url+"enterprises/"+enterpriseID+"/domains", &payload, &result, &e)
	if err != nil {
		glog.Error("Error when creating domain", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when creating domain")
	switch resp.Status() {
	case 201:
		glog.Infoln("Created the domain:", result[0].ID)
		return result[0].ID, nil
	case 409:
		//Domain already exists, call Get to retrieve the ID
		id, err := nvsdc.GetDomainID(enterpriseID, name)
		if err != nil {
			glog.Errorf("Error when getting domain ID: %s", err)
			return "", err
		} else {
			return id, nil
		}
	default:
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return "", errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) DeleteDomain(id string) error {
	result := make([]struct{}, 1)
	e := api.RESTError{}
	resp, err := nvsdc.session.Delete(nvsdc.url+"domains/"+id+"?responseChoice=1", &result, &e)
	if err != nil {
		glog.Errorf("Error when deleting domain with ID %s: %s", id, err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when deleting domain")
	switch resp.Status() {
	case 204:
		return nil
	default:
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) CreateZone(domainID, name string) (string, error) {
	result := make([]api.VsdObject, 1)
	payload := api.VsdObject{
		Name:        name,
		Description: "Auto-generated for OpenShift project \"" + name + "\"",
	}
	e := api.RESTError{}
	resp, err := nvsdc.session.Post(nvsdc.url+"domains/"+domainID+"/zones", &payload, &result, &e)
	if err != nil {
		glog.Error("Error when creating zone", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when creating zone")
	switch resp.Status() {
	case 201:
		glog.Infoln("Created the zone:", result[0].ID)
		return result[0].ID, nil
	case 409:
		//Zone already exists, call Get to retrieve the ID
		id, err := nvsdc.GetZoneID(domainID, name)
		if err != nil {
			glog.Errorf("Error when getting zone ID: %s", err)
			return "", err
		} else {
			return id, nil
		}
	default:
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return "", errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) DeleteZone(id string) error {
	// Delete subnets in this zone
	result := make([]struct{}, 1)
	e := api.RESTError{}
	resp, err := nvsdc.session.Delete(nvsdc.url+"zones/"+id+"?responseChoice=1", &result, &e)
	if err != nil {
		glog.Errorf("Error when deleting zone with ID %s: %s", id, err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when deleting zone")
	switch resp.Status() {
	case 204:
		return nil
	default:
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) CreateSubnet(name, zoneID string, subnet *IPv4Subnet) (string, error) {
	result := make([]api.VsdSubnet, 1)
	payload := api.VsdSubnet{
		IPType:      "IPV4",
		Address:     subnet.Address.String(),
		Netmask:     subnet.Netmask().String(),
		Description: "Auto-generated subnet",
		Name:        name,
		PATEnabled:  api.PATInherited,
	}
	e := api.RESTError{}
	resp, err := nvsdc.session.Post(nvsdc.url+"zones/"+zoneID+"/subnets", &payload, &result, &e)
	if err != nil {
		glog.Error("Error when creating subnet", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when creating subnet")
	switch resp.Status() {
	case 201:
		glog.Infoln("Created the subnet:", result[0].ID)
	case 409:
		//Subnet already exists, call Get to retrieve the ID
		if id, err := nvsdc.GetSubnetID(zoneID, subnet); err != nil {
			glog.Errorf("Error when getting subnet ID: %s", err)
			return "", err
		} else {
			return id, nil
		}
	default:
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return "", errors.New("Unexpected error code: " + string(resp.Status()))
	}
	return result[0].ID, nil
}

func (nvsdc *NuageVsdClient) DeleteSubnet(id string) error {
	result := make([]struct{}, 1)
	e := api.RESTError{}
	resp, err := nvsdc.session.Delete(nvsdc.url+"subnets/"+id+"?responseChoice=1", &result, &e)
	if err != nil {
		glog.Errorf("Error when deleting subnet with ID %s: %s", id, err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when deleting subnet")
	if resp.Status() != 204 {
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return errors.New("Unexpected error code: " + string(resp.Status()))
	}
	return nil
}

func (nvsdc *NuageVsdClient) GetSubnetID(zoneID string, subnet *IPv4Subnet) (string, error) {
	result := make([]api.VsdSubnet, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `address == "`+subnet.Address.String()+`"`)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"zones/"+zoneID+"/subnets", nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting subnet ID %s", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting subnet ID")
	if resp.Status() == 200 && result[0].Address == subnet.Address.String() {
		return result[0].ID, nil
	} else {
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return "", errors.New("Unexpected error code: " + string(resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) GetDomainID(enterpriseID, name string) (string, error) {
	result := make([]api.VsdObject, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `name == "`+name+`"`)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"enterprises/"+enterpriseID+"/domains", nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting domain ID %s", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting domain ID")
	if resp.Status() == 200 {
		// Status code 200 is returned even if there's no results.  If
		// the filter didn't match anything (or there was nothing to
		// return), the result object will just be empty.
		if result[0].Name == name {
			return result[0].ID, nil
		} else if result[0].Name == "" {
			return "", errors.New("Domain not found")
		} else {
			return "", errors.New(fmt.Sprintf(
				"Found %q instead of %q", result[0].Name, name))
		}
	} else {
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return "", errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) Run(nsChannel chan *api.NamespaceEvent, stop chan bool) {
	//we will use the kube client APIs than interfacing with the REST API
	for {
		select {
		case nsEvent := <-nsChannel:
			nvsdc.HandleNsEvent(nsEvent)
		}
	}
}

func (nvsdc *NuageVsdClient) HandleNsEvent(nsEvent *api.NamespaceEvent) error {
	glog.Infoln("Received a namespace event: Namespace: ", nsEvent.Name, nsEvent.Type)
	switch nsEvent.Type {
	case api.Added:
		if _, exists := nvsdc.zones[nsEvent.Name]; !exists {
			zoneID, err := nvsdc.CreateZone(nvsdc.domainID, nsEvent.Name)
			if err != nil {
				return err
			}
			nvsdc.zones[nsEvent.Name] = zoneID
			// subnetSize is guaranteed to be between 0 and 32 (inclusive) by
			// the Init() function defined above, so (32 - subnetSize) will
			// also produce a number between 0 and 32 (inclusive).
			subnet, err := nvsdc.pool.Alloc(32 - nvsdc.subnetSize)
			if err != nil {
				return err
			}
			if subnetID, err := nvsdc.CreateSubnet(nsEvent.Name+"-0", zoneID, subnet); err != nil {
				nvsdc.pool.Free(subnet)
				return err
			} else {
				nvsdc.subnets[zoneID] = &SubnetList{SubnetID: subnetID, Subnet: subnet, Next: nil}
			}
			return nil
		}
		id, err := nvsdc.GetZoneID(nvsdc.domainID, nsEvent.Name)
		switch {
		case id == "" && err == nil:
			err = errors.New("Invalid zone ID returned")
			fallthrough
		case err != nil:
			glog.Errorf("Invalid ID for zone %s", nsEvent.Name)
			return err
		case id != nvsdc.zones[nsEvent.Name]:
			glog.Warningf("Mismatched IDs for zone %s: local %s, configured %s", nsEvent.Name, nvsdc.zones[nsEvent.Name], id)
			nvsdc.zones[nsEvent.Name] = id
			return nil
		}
	case api.Deleted:
		if id, exists := nvsdc.zones[nsEvent.Name]; exists {
			// Delete subnets that we've created, and free them back into the pool
			if subnetsHead, exists := nvsdc.subnets[id]; exists {
				subnet := subnetsHead
				for subnet != nil {
					err := nvsdc.DeleteSubnet(subnet.SubnetID)
					if err != nil {
						glog.Warningf("Failed to delete subnet %q in zone %q",
							subnet.SubnetID, nsEvent.Name)
					}
					err = nvsdc.pool.Free(subnet.Subnet)
					if err != nil {
						glog.Warningf("Failed to free subnet %q from zone %q",
							subnet.Subnet.String(), nsEvent.Name)
					}
					subnet = subnet.Next
				}
				// Now that all subnets are deleted, remove the list associated
				// with this zone
				delete(nvsdc.subnets, id)
			}
			delete(nvsdc.zones, nsEvent.Name)
			return nvsdc.DeleteZone(id)
		}
		id, err := nvsdc.GetZoneID(nvsdc.domainID, nsEvent.Name)
		switch {
		case id == "" && err == nil:
			glog.Warningf("Got delete namespace event for non-existant zone %s", nsEvent.Name)
			return nil
		case err != nil:
			glog.Errorf("Error getting ID of zone %s", nsEvent.Name)
			return err
		case id != "":
			glog.Infof("Deleting zone %s which was not found locally", nsEvent.Name)
			return nvsdc.DeleteZone(id)
		}
	}
	return nil
}
