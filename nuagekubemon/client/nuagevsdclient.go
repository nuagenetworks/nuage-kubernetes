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
	vsdUrl          string
	vsdVersion      string
	vsdUsername     string
	vsdPassword     string
	vsdEnterprise   string
	vsdSession      napping.Session
	vsdEnterpriseID string
}

const clusterEnterpriseName = "Openshift-Enterprise"

func NewNuageVsdClient(nkmConfig *config.NuageKubeMonConfig) *NuageVsdClient {
	nvsdc := new(NuageVsdClient)
	nvsdc.Init(nkmConfig)
	return nvsdc
}

func (nvsdc *NuageVsdClient) GetAuthorizationToken() error {
	h := nvsdc.vsdSession.Header
	h.Add("X-Nuage-Organization", nvsdc.vsdEnterprise)
	h.Add("Authorization", "XREST "+base64.StdEncoding.EncodeToString([]byte(nvsdc.vsdUsername+":"+nvsdc.vsdPassword)))

	type ResponseElement struct {
		APIKey       string
		APIKeyExpiry int64
		ID           string
		email        string
		enterpriseID string
		firstName    string
		lastName     string
		role         string
		userName     string
	}

	var respArray [1]ResponseElement

	resp, err := nvsdc.vsdSession.Get(nvsdc.vsdUrl+"me", nil, &respArray, &api.RESTErrorResponse{})

	if err != nil {
		fmt.Println("Got an error", err)
		glog.Error("Error when requesting authorization token", err)
		return err
	} else {
		fmt.Println("Got a reponse status", resp.Status())
		if resp.Status() == 200 {
			h.Set("Authorization", "XREST "+base64.StdEncoding.EncodeToString([]byte(nvsdc.vsdUsername+":"+respArray[0].APIKey)))
		}
		return nil
	}
	return nil
}

func (nvsdc *NuageVsdClient) CreateEnterprise() error {

	payload := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{
		Name:        clusterEnterpriseName,
		Description: "Auto-generated enterprise for Openshift Cluster",
	}

	result := make([]api.VsdEnterpriseResponse, 1)
	e := api.RESTErrorResponse{}

	resp, err := nvsdc.vsdSession.Post(nvsdc.vsdUrl+"enterprises", &payload, &result, &e)
	if err != nil {
		fmt.Println("Got an error", err)
		glog.Error("Error when creating enterprise", err)

	} else {
		fmt.Println("Got a reponse status", resp.Status())
		switch resp.Status() {
		case 201:
			{
				fmt.Println("Created the enterprise: ", result[0].ID)
				nvsdc.vsdEnterpriseID = result[0].ID

			}
		case 409:
			{
				//Enterprise already exists, call Get to retrieve the ID
				id, err := nvsdc.GetEnterpriseID(clusterEnterpriseName)
				if err != nil {
					glog.Errorf("Error when getting enterprise ID: %s", err)
				} else {
					nvsdc.vsdEnterpriseID = id
				}
			}
		default:
			{
				fmt.Println("Bad response status from VSD Server")
				fmt.Printf("\t Status:  %v\n", resp.Status())
				fmt.Printf("\t Message: %v\n", e.Message)
				fmt.Printf("\t Errors: %v\n", e.Message)
				return errors.New("Unexpected error code: " + string(resp.Status()))
			}
		}
	}
	return nil
}

func (nvsdc *NuageVsdClient) CreateAdminUser() error {
	passwd := fmt.Sprintf("%x", sha1.Sum([]byte("admin")))
	payload := struct {
		UserName  string `json:"userName"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
	}{
		UserName:  "admin",
		Password:  passwd,
		FirstName: "Admin",
		LastName:  "Admin",
		Email:     "admin@localhost",
	}
	result := make([]api.VsdAdminCreateResponse, 1)
	e := api.RESTErrorResponse{}
	//Get admin ID after creating the admin user
	adminId := ""
	resp, err := nvsdc.vsdSession.Post(nvsdc.vsdUrl+"enterprises/"+nvsdc.vsdEnterpriseID+"/users", &payload, &result, &e)
	if err != nil {
		fmt.Println("Got an error", err)
		glog.Error("Error when creating admin user", err)

	} else {
		fmt.Println("Got a reponse status", resp.Status())
		switch resp.Status() {
		case 201:
			{
				fmt.Println("Created the admin user: ", result[0].ID)
				adminId = result[0].ID
			}
		case 409:
			{
				//Enterprise already exists, call Get to retrieve the ID
				id, err := nvsdc.GetAdminID("admin")
				if err != nil {
					glog.Errorf("Error when getting admin users' ID: %s", err)
				} else {
					adminId = id
				}
			}
		default:
			{
				fmt.Println("Bad response status from VSD Server")
				fmt.Printf("\t Status:  %v\n", resp.Status())
				fmt.Printf("\t Message: %v\n", e.Message)
				fmt.Printf("\t Errors: %v\n", e.Message)
				return errors.New("Unexpected error code: " + string(resp.Status()))
			}
		}
	}
	//Get admin group ID and add the admin id to the admin group
	groupId, err := nvsdc.GetAdminGroupID()
	if err != nil {
		glog.Errorf("Error when getting admin group ID: %s", err)
	} else {
		groupPayload := []string{adminId}
		e := api.RESTErrorResponse{}
		groupResp, groupErr := nvsdc.vsdSession.Put(nvsdc.vsdUrl+"groups/"+groupId+"/users", groupPayload, nil, &e)
		if groupErr != nil {
			fmt.Println("Got an error", groupErr)
			glog.Error("Error when adding admin user to the admin group", groupErr)

		} else {
			fmt.Println("Got a reponse status", groupResp.Status())
			switch groupResp.Status() {
			case 204:
				{
					fmt.Println("Added the admin user to the admin group")
				}
			case 409:
				{
					fmt.Println("Admin user already in admin group")
				}
			default:
				{
					fmt.Println("Bad response status from VSD Server")
					fmt.Printf("\t Status:  %v\n", groupResp.Status())
					fmt.Printf("\t Message: %v\n", e.Message)
					fmt.Printf("\t Errors: %v\n", e.Message)
					return errors.New("Unexpected error code: " + string(groupResp.Status()))
				}
			}
		}

	}
	return nil
}

func (nvsdc *NuageVsdClient) GetAdminID(name string) (string, error) {
	respArray := make([]api.VsdAdminCreateResponse, 1)
	h := nvsdc.vsdSession.Header

	h.Add("X-Nuage-Filter", `userName == "`+name+`"`)
	resp, err := nvsdc.vsdSession.Get(nvsdc.vsdUrl+"enterprises/"+nvsdc.vsdEnterpriseID+"/users", nil, &respArray, &api.RESTErrorResponse{})
	h.Del("X-Nuage-Filter")
	if err != nil {
		fmt.Println("Got an error", err)
		glog.Errorf("Error when getting enterprise ID %s", err)
		return "", err

	} else {
		fmt.Println("Got a reponse status", resp.Status())
		if resp.Status() == 200 && respArray[0].UserName == name {
			return respArray[0].ID, nil
		}
	}
	return "", errors.New("Unexpected error occured: " + string(resp.Status()))
}

func (nvsdc *NuageVsdClient) GetAdminGroupID() (string, error) {
	respArray := make([]api.VsdAdminGroupResponse, 1)
	h := nvsdc.vsdSession.Header

	h.Add("X-Nuage-Filter", `role == "ORGADMIN"`)
	resp, err := nvsdc.vsdSession.Get(nvsdc.vsdUrl+"enterprises/"+nvsdc.vsdEnterpriseID+"/groups", nil, &respArray, &api.RESTErrorResponse{})
	h.Del("X-Nuage-Filter")
	if err != nil {
		fmt.Println("Got an error", err)
		glog.Errorf("Error when getting enterprise ID %s", err)
		return "", err

	} else {
		fmt.Println("Got a reponse status", resp.Status())
		if resp.Status() == 200 && respArray[0].Role == "ORGADMIN" {
			return respArray[0].ID, nil
		}
	}
	return "", errors.New("Unexpected error occured: " + string(resp.Status()))
}

func (nvsdc *NuageVsdClient) GetEnterpriseID(name string) (string, error) {
	respArray := make([]api.VsdEnterpriseResponse, 1)
	h := nvsdc.vsdSession.Header

	h.Add("X-Nuage-Filter", `name == "`+name+`"`)
	resp, err := nvsdc.vsdSession.Get(nvsdc.vsdUrl+"enterprises", nil, &respArray, &api.RESTErrorResponse{})
	h.Del("X-Nuage-Filter")
	if err != nil {
		fmt.Println("Got an error", err)
		glog.Errorf("Error when getting enterprise ID %s", err)
		return "", err

	} else {
		fmt.Println("Got a reponse status", resp.Status())
		if resp.Status() == 200 && respArray[0].Name == name {
			return respArray[0].ID, nil
		}
	}
	return "", errors.New("Unexpected error occured: " + string(resp.Status()))
}

func (nvsdc *NuageVsdClient) CreateSession() {
	nvsdc.vsdUsername = "csproot"
	nvsdc.vsdPassword = "csproot"
	nvsdc.vsdEnterprise = "csp"

	nvsdc.vsdSession = napping.Session{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		Header: &http.Header{},
	}

	h := nvsdc.vsdSession.Header
	h.Add("Content-Type", "application/json")
}

func (nvsdc *NuageVsdClient) LoginAsAdmin() {
	nvsdc.vsdUsername = "admin"
	nvsdc.vsdPassword = "admin"
	nvsdc.vsdEnterprise = clusterEnterpriseName
	h := nvsdc.vsdSession.Header
	h.Del("X-Nuage-Organization")
	h.Del("Authorization")
	nvsdc.GetAuthorizationToken()
}

func (nvsdc *NuageVsdClient) Init(nkmConfig *config.NuageKubeMonConfig) {
	nvsdc.vsdVersion = nkmConfig.NuageVspVersion
	nvsdc.vsdUrl = nkmConfig.NuageVsdApiUrl + "/nuage/api/" + nvsdc.vsdVersion + "/"
	nvsdc.CreateSession()
	nvsdc.GetAuthorizationToken()
	nvsdc.CreateEnterprise()
	nvsdc.CreateAdminUser()
	nvsdc.InstallLicense(nkmConfig.LicenseFile)
	nvsdc.LoginAsAdmin()
}

func (nvsdc *NuageVsdClient) InstallLicense(licensePath string) error {
	if licensePath != "" {
		//try installing the license file
		license, err := ioutil.ReadFile(licensePath)
		if err != nil {
			glog.Error("Failed to read license file", err)
			return err
		} else {
			licenseString := strings.TrimSpace(string(license))
			payload := struct {
				License string `json:"license"`
			}{
				License: licenseString,
			}
			result := make([]api.VsdCreateLicenseResponse, 1)
			e := api.RESTErrorResponse{}
			glog.Info("Attempting to install license file", licensePath)
			resp, err := nvsdc.vsdSession.Post(nvsdc.vsdUrl+"licenses", &payload, &result, &e)
			if err != nil {
				fmt.Println("Got an error", err)
				glog.Error("Error when installing license", err)

			} else {
				fmt.Println("License Install: reponse status", resp.Status())
				switch resp.Status() {
				case 201:
					{
						fmt.Println("Installed the license: ", result[0].LicenseId)
					}
				case 409:
					{
						//TODO: license already exists, call Get to retrieve the ID? Do we need to delete the existing license?
						glog.Info("License already exists")
					}
				default:
					{
						fmt.Println("Bad response status from VSD Server")
						fmt.Printf("\t Status:  %v\n", resp.Status())
						fmt.Printf("\t Message: %v\n", e.Message)
						fmt.Printf("\t Errors: %v\n", e.Message)
						return errors.New("Unexpected error code: " + string(resp.Status()))
					}
				}
			}
		}
	} else {
		glog.Error("No license file specified")
		//check if a license already exists.
		// if it does then its not an error
		return nvsdc.GetLicense()
	}
	return nil
}

func (nvsdc *NuageVsdClient) GetLicense() error {

	result := make([]api.VsdCreateLicenseResponse, 1)
	resp, err := nvsdc.vsdSession.Get(nvsdc.vsdUrl+"licenses", nil, &result, &api.RESTErrorResponse{})

	if err != nil {
		glog.Error("Error when requesting license", err)
		return err
	} else {
		fmt.Println("GetLicense() got a reponse status", resp.Status())
		if resp.Status() == 200 {
			return nil
		} else {
			return errors.New("GetLicense() failed: " + string(resp.Status()))
		}

	}
	return nil
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
	fmt.Println("Received a namespace event: Namespace: ", nsEvent.Name, nsEvent.Type)
	return nil
}
