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
	url              string
	version          string
	username         string
	password         string
	enterprise       string
	session          napping.Session
	enterpriseID     string
	domainTemplateID string
	zoneTemplateID   string
	domains          map[string]string
}

const clusterEnterpriseName = "Openshift-Enterprise"
const clusterDomainTemplateName = "Openshift-Domain-Template"
const clusterZoneTemplateName = "default"

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
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return errors.New("Unexpected error code: " + string(resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) CreateEnterprise() error {
	payload := api.VsdEnterprise{
		Name:        clusterEnterpriseName,
		Description: "Auto-generated enterprise for Openshift Cluster",
	}
	result := make([]api.VsdEnterprise, 1)
	e := api.RESTError{}
	resp, err := nvsdc.session.Post(nvsdc.url+"enterprises", &payload, &result, &e)
	if err != nil {
		glog.Error("Error when creating enterprise", err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when creating the enterprise")
	switch resp.Status() {
	case 201:
		glog.Infoln("Created the enterprise: ", result[0].ID)
		nvsdc.enterpriseID = result[0].ID
	case 409:
		//Enterprise already exists, call Get to retrieve the ID
		id, err := nvsdc.GetEnterpriseID(clusterEnterpriseName)
		if err != nil {
			glog.Errorf("Error when getting enterprise ID: %s", err)
			return err
		} else {
			nvsdc.enterpriseID = id
		}
	default:
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return errors.New("Unexpected error code: " + string(resp.Status()))
	}
	return nil
}

func (nvsdc *NuageVsdClient) CreateAdminUser() error {
	passwd := fmt.Sprintf("%x", sha1.Sum([]byte("admin")))
	payload := api.VsdUser{
		UserName:  "admin",
		Password:  passwd,
		FirstName: "Admin",
		LastName:  "Admin",
		Email:     "admin@localhost",
	}
	result := make([]api.VsdUser, 1)
	e := api.RESTError{}
	//Get admin ID after creating the admin user
	var adminId string
	resp, err := nvsdc.session.Post(nvsdc.url+"enterprises/"+nvsdc.enterpriseID+"/users", &payload, &result, &e)
	if err != nil {
		glog.Error("Error when creating admin user", err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when creating the admin user")
	switch resp.Status() {
	case 201:
		glog.Infoln("Created the admin user: ", result[0].ID)
		adminId = result[0].ID
	case 409:
		//Enterprise already exists, call Get to retrieve the ID
		id, err := nvsdc.GetAdminID("admin")
		if err != nil {
			glog.Errorf("Error when getting admin user's ID: %s", err)
		} else {
			adminId = id
		}
	default:
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return errors.New("Unexpected error code: " + string(resp.Status()))
	}
	//Get admin group ID and add the admin id to the admin group
	groupId, err := nvsdc.GetAdminGroupID()
	if err != nil {
		glog.Errorf("Error when getting admin group ID: %s", err)
		return err
	}
	groupPayload := []string{adminId}
	e = api.RESTError{}
	resp, err = nvsdc.session.Put(nvsdc.url+"groups/"+groupId+"/users", &groupPayload, nil, &e)
	if err != nil {
		glog.Error("Error when adding admin user to the admin group", err)
		return err
	} else {
		glog.Infoln("Got a reponse status", resp.Status(), "when adding user to the admin group")
		switch resp.Status() {
		case 204:
			glog.Infoln("Added the admin user to the admin group")
		case 409:
			glog.Infoln("Admin user already in admin group")
		default:
			glog.Errorln("Bad response status from VSD Server")
			glog.Errorf("\t Status:  %v\n", resp.Status())
			glog.Errorf("\t Message: %v\n", e.Message)
			glog.Errorf("\t Errors: %v\n", e.Message)
			return errors.New("Unexpected error code: " + string(resp.Status()))
		}
	}
	return nil
}

func (nvsdc *NuageVsdClient) GetAdminID(name string) (string, error) {
	result := make([]api.VsdUser, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `userName == "`+name+`"`)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"enterprises/"+nvsdc.enterpriseID+"/users", nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting admin user ID %s", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting user ID")
	if resp.Status() == 200 && result[0].UserName == name {
		return result[0].ID, nil
	} else {
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return "", errors.New("Unexpected error code: " + string(resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) GetAdminGroupID() (string, error) {
	result := make([]api.VsdGroup, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `role == "ORGADMIN"`)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"enterprises/"+nvsdc.enterpriseID+"/groups", nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting admin group ID %s", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting ID of group ORGADMIN")
	if resp.Status() == 200 && result[0].Role == "ORGADMIN" {
		return result[0].ID, nil
	} else {
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return "", errors.New("Unexpected error code: " + string(resp.Status()))
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
	if resp.Status() == 200 && result[0].Name == name {
		return result[0].ID, nil
	} else {
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return "", errors.New("Unexpected error code: " + string(resp.Status()))
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

func (nvsdc *NuageVsdClient) LoginAsAdmin() {
	nvsdc.username = "admin"
	nvsdc.password = "admin"
	nvsdc.enterprise = clusterEnterpriseName
	h := nvsdc.session.Header
	h.Del("X-Nuage-Organization")
	h.Del("Authorization")
	nvsdc.GetAuthorizationToken()
}

func (nvsdc *NuageVsdClient) Init(nkmConfig *config.NuageKubeMonConfig) {
	nvsdc.version = nkmConfig.NuageVspVersion
	nvsdc.url = nkmConfig.NuageVsdApiUrl + "/nuage/api/" + nvsdc.version + "/"
	nvsdc.domains = make(map[string]string)
	nvsdc.CreateSession()
	nvsdc.GetAuthorizationToken()
	nvsdc.CreateEnterprise()
	nvsdc.CreateAdminUser()
	nvsdc.InstallLicense(nkmConfig.LicenseFile)
	nvsdc.LoginAsAdmin()
	nvsdc.CreateTemplates()
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
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return errors.New("Unexpected error code: " + string(resp.Status()))
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
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return errors.New("Unexpected error code: " + string(resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) CreateTemplates() error {
	err := nvsdc.CreateDomainTemplate()
	if err != nil {
		return err
	}
	err = nvsdc.CreateZoneTemplate()
	if err != nil {
		return err
	}
	return nil
}

func (nvsdc *NuageVsdClient) CreateDomainTemplate() error {
	result := make([]api.VsdObject, 1)
	payload := api.VsdObject{
		Name:        clusterDomainTemplateName,
		Description: "Auto-generated default domain template",
	}
	e := api.RESTError{}
	resp, err := nvsdc.session.Post(nvsdc.url+"enterprises/"+nvsdc.enterpriseID+"/domaintemplates", &payload, &result, &e)
	if err != nil {
		glog.Error("Error when creating domain template", err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when creating domain template")
	switch resp.Status() {
	case 201:
		glog.Infoln("Created the domain: ", result[0].ID)
		nvsdc.domainTemplateID = result[0].ID
	case 409:
		//Enterprise already exists, call Get to retrieve the ID
		id, err := nvsdc.GetDomainTemplateID(clusterDomainTemplateName)
		if err != nil {
			glog.Errorf("Error when getting domain template ID: %s", err)
			return err
		}
		nvsdc.domainTemplateID = id
	default:
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return errors.New("Unexpected error code: " + string(resp.Status()))
	}
	return nil
}

func (nvsdc *NuageVsdClient) GetDomainTemplateID(name string) (string, error) {
	result := make([]api.VsdObject, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `name == "`+name+`"`)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"enterprises/"+nvsdc.enterpriseID+"/domaintemplates", nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting domain template ID %s", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting domain template ID")
	if resp.Status() == 200 && result[0].Name == name {
		return result[0].ID, nil
	} else {
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return "", errors.New("Unexpected error code: " + string(resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) CreateZoneTemplate() error {
	result := make([]api.VsdObject, 1)
	payload := api.VsdObject{
		Name:        clusterZoneTemplateName,
		Description: "Auto-generated default zone template",
	}
	e := api.RESTError{}
	resp, err := nvsdc.session.Post(nvsdc.url+"domaintemplates/"+nvsdc.domainTemplateID+"/zonetemplates", &payload, &result, &e)
	if err != nil {
		glog.Error("Error when creating zone template", err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when creating zone template")
	switch resp.Status() {
	case 201:
		glog.Infoln("Created the zone: ", result[0].ID)
		nvsdc.zoneTemplateID = result[0].ID
	case 409:
		//Enterprise already exists, call Get to retrieve the ID
		id, err := nvsdc.GetZoneTemplateID(clusterZoneTemplateName)
		if err != nil {
			glog.Errorf("Error when getting zone template ID: %s", err)
			return err
		}
		nvsdc.zoneTemplateID = id
	default:
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return errors.New("Unexpected error code: " + string(resp.Status()))
	}
	return nil
}

func (nvsdc *NuageVsdClient) GetZoneTemplateID(name string) (string, error) {
	result := make([]api.VsdObject, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `name == "`+name+`"`)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"domaintemplates/"+nvsdc.domainTemplateID+"/zonetemplates", nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting zone template ID %s", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting zone template ID")
	if resp.Status() == 200 && result[0].Name == name {
		return result[0].ID, nil
	} else {
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return "", errors.New("Unexpected error code: " + string(resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) CreateDomain(name string) error {
	result := make([]api.VsdObjectInstance, 1)
	payload := api.VsdObjectInstance{
		Name:        name,
		Description: "Auto-generated domain for " + name,
		TemplateID:  nvsdc.domainTemplateID,
	}
	e := api.RESTError{}
	resp, err := nvsdc.session.Post(nvsdc.url+"enterprises/"+nvsdc.enterpriseID+"/domains", &payload, &result, &e)
	if err != nil {
		glog.Error("Error when creating domain", err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when creating domain")
	switch resp.Status() {
	case 201:
		glog.Infoln("Created the domain:", result[0].ID)
		nvsdc.domains[name] = result[0].ID
	case 409:
		//Domain already exists, call Get to retrieve the ID
		id, err := nvsdc.GetDomainID(name)
		if err != nil {
			glog.Errorf("Error when getting domain ID: %s", err)
			return err
		} else {
			nvsdc.domains[name] = id
		}
	default:
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return errors.New("Unexpected error code: " + string(resp.Status()))
	}
	return nil
}

func (nvsdc *NuageVsdClient) DeleteDomain(name, id string) error {
	result := make([]struct{}, 1)
	e := api.RESTError{}
	resp, err := nvsdc.session.Delete(nvsdc.url+"domains/"+id+"?responseChoice=1", &result, &e)
	if err != nil {
		glog.Errorf("Error when deleting domain with ID %s: %s", id, err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when delete domain")
	switch resp.Status() {
	case 204:
		delete(nvsdc.domains, name)
		return nil
	default:
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return errors.New("Unexpected error code: " + string(resp.Status()))
	}
}

func (nvsdc *NuageVsdClient) GetDomainID(name string) (string, error) {
	result := make([]api.VsdObject, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `name == "`+name+`"`)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"enterprises/"+nvsdc.enterpriseID+"/domains", nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting domain ID %s", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting domain ID")
	if resp.Status() == 200 && result[0].Name == name {
		return result[0].ID, nil
	} else {
		glog.Errorln("Bad response status from VSD Server")
		glog.Errorf("\t Status:  %v\n", resp.Status())
		glog.Errorf("\t Message: %v\n", e.Message)
		glog.Errorf("\t Errors: %v\n", e.Message)
		return "", errors.New("Unexpected error code: " + string(resp.Status()))
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
		if _, exists := nvsdc.domains[nsEvent.Name]; !exists {
			return nvsdc.CreateDomain(nsEvent.Name)
		}
		id, err := nvsdc.GetDomainID(nsEvent.Name)
		switch {
		case id == "" && err == nil:
			err = errors.New("Invalid domain ID returned")
			fallthrough
		case err != nil:
			glog.Errorf("Invalid ID for domain %s", nsEvent.Name)
			return err
		case id != nvsdc.domains[nsEvent.Name]:
			glog.Warningf("Mismatched IDs for domain %s: local %s, configured %s", nsEvent.Name, nvsdc.domains[nsEvent.Name], id)
			nvsdc.domains[nsEvent.Name] = id
			return nil
		}
	case api.Deleted:
		if id, exists := nvsdc.domains[nsEvent.Name]; exists {
			return nvsdc.DeleteDomain(nsEvent.Name, id)
		}
		id, err := nvsdc.GetDomainID(nsEvent.Name)
		switch {
		case id == "" && err == nil:
			glog.Warningf("Got delete namespace event for non-existant domain %s", nsEvent.Name)
			return nil
		case err != nil:
			glog.Errorf("Error getting ID of domain %s", nsEvent.Name)
			return err
		case id != "":
			glog.Infof("Deleting domain %s which was not found locally", nsEvent.Name)
			return nvsdc.DeleteDomain(nsEvent.Name, id)
		}
	}
	return nil
}
