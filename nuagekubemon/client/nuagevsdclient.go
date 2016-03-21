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
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/jmcvetta/napping"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/api"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/config"
	"github.com/rfredette/sleepy"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type NuageVsdClient struct {
	url                   string
	version               string
	username              string
	password              string
	enterprise            string
	session               napping.Session
	enterpriseID          string
	domainID              string
	namespaces            map[string]NamespaceData //namespace name -> namespace data
	pods                  *PodList                 //<namespace>/<pod-name> -> subnet
	pool                  IPv4SubnetPool
	clusterNetwork        *IPv4Subnet //clusterNetworkCIDR used to generate pool
	serviceNetwork        *IPv4Subnet
	ingressAclTemplateID  string
	egressAclTemplateID   string
	nextAvailablePriority int
	subnetSize            int         //the size in bits of the subnets we allocate (i.e. size 8 produces /24 subnets).
	restAPI               *sleepy.API //TODO: split the rest server into its own package
	restServer            *http.Server
	newSubnetQueue        chan config.NamespaceUpdateRequest //list of namespaces that need new subnets
}

type NamespaceData struct {
	ZoneID              string
	Name                string
	NetworkMacroGroupID string
	NetworkMacros       map[string]string //service name (qualified with the namespace) -> network macro id
	Subnets             *SubnetNode
	NeedsNewSubnet      bool
	numSubnets          int //used for naming new subnets (nsname-0, nsname-1, etc.)
}

type SubnetNode struct {
	SubnetID   string
	Subnet     *IPv4Subnet
	SubnetName string
	ActiveIPs  int //Number of IP addresses that are accounted for in this subnet.
	Next       *SubnetNode
}

func NewNuageVsdClient(nkmConfig *config.NuageKubeMonConfig) *NuageVsdClient {
	nvsdc := new(NuageVsdClient)
	nvsdc.Init(nkmConfig)
	return nvsdc
}

func (nvsdc *NuageVsdClient) GetAuthorizationToken() error {
	h := nvsdc.session.Header
	// Add the organization header if it's not present
	if h.Get("X-Nuage-Organization") == "" {
		h.Add("X-Nuage-Organization", nvsdc.enterprise)
	}
	if h.Get("Authorization") == "" {
		h.Add("Authorization", "XREST "+base64.StdEncoding.EncodeToString([]byte(nvsdc.username+":"+nvsdc.password)))
	} else {
		h.Set("Authorization", "XREST "+base64.StdEncoding.EncodeToString([]byte(nvsdc.username+":"+nvsdc.password)))
	}
	var result [1]api.VsdAuthToken
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"me", nil, &result, &e)
	if err != nil {
		glog.Error("Error when requesting authorization token", err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status())
	if resp.Status() == http.StatusOK {
		// Launch a separate go routine to get the new token 3 minutes before
		// this one expires
		// APIKeyExpiry is in milliseconds, while Unix time is in seconds.
		delay := time.Duration(result[0].APIKeyExpiry-(time.Now().Unix()*1000))*time.Millisecond - 3*time.Minute
		// If there's less than 3 minutes until expiration, get the new key when
		// the old one expires.
		if delay < 0 {
			delay = time.Duration(result[0].APIKeyExpiry-(time.Now().Unix()*1000)) * time.Millisecond
		}
		time.AfterFunc(delay, func() { nvsdc.GetAuthorizationToken() })
		h.Set("Authorization", "XREST "+base64.StdEncoding.EncodeToString([]byte(nvsdc.username+":"+result[0].APIKey)))
		return nil
	} else {
		return VsdErrorResponse(resp, &e)
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
	case http.StatusCreated:
		glog.Infoln("Created the enterprise: ", result[0].ID)
		return result[0].ID, nil
	case http.StatusConflict:
		glog.Infoln("Error from VSD:\n", e)
		//Enterprise already exists, call Get to retrieve the ID
		id, err := nvsdc.GetEnterpriseID(enterpriseName)
		if err != nil {
			glog.Errorf("Error when getting enterprise ID: %s", err)
			return "", err
		} else {
			return id, nil
		}
	default:
		return "", VsdErrorResponse(resp, &e)
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
	case http.StatusCreated:
		glog.Infoln("Created the admin user: ", result[0].ID)
		adminId = result[0].ID
	case http.StatusConflict:
		glog.Infoln("Error from VSD:\n", e)
		//Enterprise already exists, call Get to retrieve the ID
		id, erradminID := nvsdc.GetAdminID(enterpriseID, "admin")
		if erradminID != nil {
			glog.Errorf("Error when getting admin user's ID: %s", erradminID)
		} else {
			adminId = id
		}
	default:
		return "", VsdErrorResponse(resp, &e)
	}
	//Get admin group ID and add the admin id to the admin group
	groupId, err := nvsdc.GetAdminGroupID(enterpriseID)
	if err != nil {
		glog.Errorf("Error when getting admin group ID: %s", err)
		return "", err
	}
	err = nvsdc.AddUserToGroup(adminId, groupId)
	if err != nil {
		glog.Errorf("Failed to add %s to admin group", user)
		return "", err
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
	if resp.Status() == http.StatusOK {
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
		return "", VsdErrorResponse(resp, &e)
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
	if resp.Status() == http.StatusOK {
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
		return "", VsdErrorResponse(resp, &e)
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
	if resp.Status() == http.StatusOK {
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
		return "", VsdErrorResponse(resp, &e)
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
	var err error
	nvsdc.clusterNetwork, err = IPv4SubnetFromString(nkmConfig.MasterConfig.NetworkConfig.ClusterCIDR)
	if err != nil {
		glog.Fatalf("Failure in getting cluster CIDR: %s\n", err)
	}
	nvsdc.serviceNetwork, err = IPv4SubnetFromString(nkmConfig.MasterConfig.NetworkConfig.ServiceCIDR)
	if err != nil {
		glog.Fatalf("Failure in getting service CIDR: %s\n", err)
	}
	nvsdc.subnetSize = nkmConfig.MasterConfig.NetworkConfig.SubnetLength
	if nvsdc.subnetSize < 0 || nvsdc.subnetSize > 32 {
		glog.Errorf("Invalid hostSubnetLength of %d.  Using default value of 8",
			nvsdc.subnetSize)
		nvsdc.subnetSize = 8
	}
	if nvsdc.subnetSize > (32 - nvsdc.clusterNetwork.CIDRMask) {
		// If the size of the subnet (in bits) is larger than the total pool
		// size (in bits), we can't even allocate 1 subnet.  Default to using
		// half the remaining bits per subnet, rounded down (/24 has 8 bits
		// remaining, so use 4 bits per subnet).
		newSize := (32 - nvsdc.clusterNetwork.CIDRMask) / 2
		glog.Fatalf("Cannot allocate %d bit subnets from %s.  Using %d bits per subnet.",
			nvsdc.subnetSize, nvsdc.clusterNetwork.String(), newSize)
		nvsdc.subnetSize = newSize
	}
	// A null IPv4SubnetPool acts like all addresses are allocated, so we can
	// initialize it to have the available cluster address space by just
	// Free()-ing it.
	nvsdc.pool.Free(nvsdc.clusterNetwork)
	nvsdc.namespaces = make(map[string]NamespaceData)
	nvsdc.newSubnetQueue = make(chan config.NamespaceUpdateRequest)
	nvsdc.pods = NewPodList(nvsdc.namespaces, nvsdc.newSubnetQueue)
	nvsdc.CreateSession()
	nvsdc.nextAvailablePriority = 0

	err = nvsdc.GetAuthorizationToken()
	if err != nil {
		glog.Fatal(err)
	}
	nvsdc.enterpriseID, err = nvsdc.CreateEnterprise(nkmConfig.EnterpriseName)
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
	err = nvsdc.LoginAsAdmin("admin", "admin", nkmConfig.EnterpriseName)
	if err != nil {
		glog.Fatal(err)
	}
	domainTemplateID, err := nvsdc.CreateDomainTemplate(nvsdc.enterpriseID,
		nkmConfig.DomainName+"-Template")
	if err != nil {
		glog.Fatal(err)
	}
	nvsdc.domainID, err = nvsdc.CreateDomain(nvsdc.enterpriseID,
		domainTemplateID, nkmConfig.DomainName)
	if err != nil {
		glog.Fatal(err)
	}
	_, err = nvsdc.CreateIngressAclTemplate(nvsdc.domainID)
	if err != nil {
		glog.Fatal(err)
	}
	_, err = nvsdc.CreateEgressAclTemplate(nvsdc.domainID)
	if err != nil {
		glog.Fatal(err)
	}
	err = nvsdc.StartRestServer(nkmConfig.RestServer)
	if err != nil {
		glog.Fatal(err)
	}
}

func (nvsdc *NuageVsdClient) StartRestServer(restServerCfg config.RestServerConfig) error {
	// Process config options
	url := restServerCfg.Url
	if url == "" {
		url = "0.0.0.0:9443"
	}
	certDir := restServerCfg.CertificateDirectory
	if certDir == "" {
		certDir = "/usr/share/" + os.Args[0]
	}
	clientCA := restServerCfg.ClientCA
	if clientCA == "" {
		clientCA = certDir + "/nuageMonCA.crt"
	}
	glog.Infof("Using %s as rest server CA", clientCA)
	serverCert := restServerCfg.ServerCertificate
	if serverCert == "" {
		serverCert = certDir + "/nuageMonServer.crt"
	}
	glog.Infof("Using %s as rest server cert", serverCert)
	serverKey := restServerCfg.ServerKey
	if serverKey == "" {
		serverKey = certDir + "/nuageMonServer.key"
	}
	glog.Infof("Using %s as rest server key", serverKey)
	CAPool := x509.NewCertPool()
	// Read in the CA certificate, and add it to the pool of valid CAs
	clientCAData, err := ioutil.ReadFile(clientCA)
	if err != nil {
		return err
	}
	/*clientCACert, err := x509.ParseCertificate(clientCAData)
	if err != nil {
		return err
	}*/
	CAPool.AppendCertsFromPEM(clientCAData)
	// Create the rest API router, and add endpoints
	nvsdc.restAPI = sleepy.NewAPI()
	nvsdc.restAPI.AddResource(nvsdc.pods, "/namespaces/{namespace}/pods",
		"/namespaces/{namespace}/pods/{podName}")
	// Create the server config
	nvsdc.restServer = &http.Server{
		Addr:           url,
		Handler:        nvsdc.restAPI.Mux(),
		MaxHeaderBytes: 1 << 20, // not sure if this is necessary
		TLSConfig: &tls.Config{
			Certificates: make([]tls.Certificate, 1),
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    CAPool,
			RootCAs:      CAPool,
			MinVersion:   tls.VersionTLS10,
		},
	}
	// Add the server certificate to the certificate chain
	nvsdc.restServer.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(serverCert, serverKey)
	if err != nil {
		return err
	}
	// TODO: if TLS setup is unsucessful, serve over http instead
	go nvsdc.restServer.ListenAndServeTLS(serverCert, serverKey)
	return nil
}

func (nvsdc *NuageVsdClient) InstallLicense(licensePath string) error {
	if licensePath == "" {
		glog.Info("No license file specified")
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
	case http.StatusCreated:
		glog.Infoln("Installed the license: ", result[0].LicenseId)
	case http.StatusConflict:
		//TODO: license already exists, call Get to retrieve the ID? Do we need to delete the existing license?
		glog.Info("License already exists")
	default:
		return VsdErrorResponse(resp, &e)
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
	if resp.Status() == http.StatusOK {
		return nil
	} else {
		return VsdErrorResponse(resp, &e)
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
	case http.StatusCreated:
		glog.Infoln("Created the domain: ", result[0].ID)
		return result[0].ID, nil
	case http.StatusConflict:
		//Enterprise already exists, call Get to retrieve the ID
		id, err := nvsdc.GetDomainTemplateID(enterpriseID, domainTemplateName)
		if err != nil {
			glog.Errorf("Error when getting domain template ID: %s", err)
			return "", err
		}
		return id, nil
	default:
		return "", VsdErrorResponse(resp, &e)
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
	if resp.Status() == http.StatusOK {
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
		return "", VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) GetIngressAclTemplate(domainID, name string) (*api.VsdAclTemplate, error) {
	result := make([]api.VsdAclTemplate, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `name == "`+name+`"`)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"domains/"+domainID+"/ingressacltemplates", nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting ingress ACL template ID %s", err)
		return nil, err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting ingress ACL template ID")
	if resp.Status() == http.StatusOK {
		// Status code 200 is returned even if there's no results.  If
		// the filter didn't match anything (or there was nothing to
		// return), the result object will just be empty.
		if result[0].Name == name {
			return &result[0], nil
		} else if result[0].Name == "" {
			return nil, errors.New("Ingress ACL Template not found")
		} else {
			return nil, errors.New(fmt.Sprintf(
				"Found %q instead of %q", result[0].Name, name))
		}
	} else {
		return nil, VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) GetAclTemplateByID(templateID string, ingress bool) (*api.VsdAclTemplate, error) {
	result := make([]api.VsdAclTemplate, 1)
	e := api.RESTError{}
	url := nvsdc.url + "egressacltemplates/" + templateID
	if ingress {
		url = nvsdc.url + "ingressacltemplates/" + templateID
	}
	resp, err := nvsdc.session.Get(url, nil, &result, &e)
	if err != nil {
		glog.Errorf("Error when getting ACL template with ID %s: %s", templateID, err)
		return nil, err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting ACL template")
	if resp.Status() == http.StatusOK {
		// Status code 200 is returned even if there's no results.  If
		// the filter didn't match anything (or there was nothing to
		// return), the result object will just be empty.
		return &result[0], nil
	} else {
		return nil, VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) GetEgressAclTemplate(domainID, name string) (*api.VsdAclTemplate, error) {
	result := make([]api.VsdAclTemplate, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `name == "`+name+`"`)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"domains/"+domainID+"/egressacltemplates", nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting egress ACL template ID %s", err)
		return nil, err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting egress ACL template ID")
	if resp.Status() == http.StatusOK {
		// Status code 200 is returned even if there's no results.  If
		// the filter didn't match anything (or there was nothing to
		// return), the result object will just be empty.
		if result[0].Name == name {
			return &result[0], nil
		} else if result[0].Name == "" {
			return nil, errors.New("Egress ACL Template not found")
		} else {
			return nil, errors.New(fmt.Sprintf(
				"Found %q instead of %q", result[0].Name, name))
		}
	} else {
		return nil, VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) CreateIngressAclEntries() error {
	aclEntry := api.VsdAclEntry{
		Action:       "FORWARD",
		DSCP:         "*",
		Description:  "Allow Intra-Zone Traffic",
		EntityScope:  "ENTERPRISE",
		EtherType:    "0x0800",
		LocationType: "ANY",
		NetworkType:  "ENDPOINT_ZONE",
		PolicyState:  "LIVE",
		Priority:     0,
		Protocol:     "ANY",
		Reflexive:    false,
	}
	_, err := nvsdc.CreateAclEntry(true, &aclEntry)
	if err != nil {
		glog.Error("Error when creating ingress acl entry", err)
		return err
	}
	aclEntry.Action = "DROP"
	aclEntry.Description = "Drop intra-domain traffic"
	aclEntry.NetworkType = "ENDPOINT_DOMAIN"
	aclEntry.Priority = api.MAX_VSD_ACL_PRIORITY
	_, err = nvsdc.CreateAclEntry(true, &aclEntry)
	if err != nil {
		glog.Error("Error when creating ingress acl entry", err)
	}
	networkMacro := &api.VsdNetworkMacro{
		Name:    `NetworkMacro for Service CIDR`,
		IPType:  "IPV4",
		Address: nvsdc.serviceNetwork.Address.String(),
		Netmask: nvsdc.serviceNetwork.Netmask().String(),
	}
	networkMacroID, err := nvsdc.CreateNetworkMacro(nvsdc.enterpriseID, networkMacro)
	if err != nil {
		glog.Error("Error when creating the network macro for service CIDR")
	} else {
		//
		aclEntry.Priority = aclEntry.Priority - 1
		aclEntry.NetworkType = "ENTERPRISE_NETWORK"
		aclEntry.NetworkID = networkMacroID
		aclEntry.Description = "Drop traffic from domain to the service CIDR"
		_, err = nvsdc.CreateAclEntry(true, &aclEntry)
		if err != nil {
			glog.Error("Error when creating ingress acl entry", err)
		}
	}
	return nil
}

func (nvsdc *NuageVsdClient) CreateEgressAclEntries() error {
	aclEntry := api.VsdAclEntry{
		Action:       "FORWARD",
		DSCP:         "*",
		Description:  "Allow Intra-Zone Traffic",
		EntityScope:  "ENTERPRISE",
		EtherType:    "0x0800",
		LocationType: "ANY",
		NetworkType:  "ENDPOINT_ZONE",
		PolicyState:  "LIVE",
		Priority:     0,
		Protocol:     "ANY",
		Reflexive:    false,
	}
	_, err := nvsdc.CreateAclEntry(false, &aclEntry)
	if err != nil {
		glog.Error("Error when creating egress acl entry", err)
		return err
	}
	aclEntry.Action = "DROP"
	aclEntry.Description = "Drop intra-domain traffic"
	aclEntry.NetworkType = "ENDPOINT_DOMAIN"
	aclEntry.Priority = api.MAX_VSD_ACL_PRIORITY
	_, err = nvsdc.CreateAclEntry(false, &aclEntry)
	if err != nil {
		glog.Error("Error when creating egress acl entry", err)
	}
	networkMacro := &api.VsdNetworkMacro{
		Name:    `NetworkMacro for Service CIDR`,
		IPType:  "IPV4",
		Address: nvsdc.serviceNetwork.Address.String(),
		Netmask: nvsdc.serviceNetwork.Netmask().String(),
	}
	networkMacroID, err := nvsdc.CreateNetworkMacro(nvsdc.enterpriseID, networkMacro)
	if err != nil {
		glog.Error("Error when creating the network macro for service CIDR")
	} else {
		//
		aclEntry.Priority = aclEntry.Priority - 1
		aclEntry.NetworkType = "ENTERPRISE_NETWORK"
		aclEntry.NetworkID = networkMacroID
		aclEntry.Description = "Drop traffic from domain to the service CIDR"
		_, err = nvsdc.CreateAclEntry(false, &aclEntry)
		if err != nil {
			glog.Error("Error when creating ingress acl entry", err)
		}
	}
	return nil
}

func (nvsdc *NuageVsdClient) CreateIngressAclTemplate(domainID string) (string, error) {
	result := make([]api.VsdAclTemplate, 1)
	payload := api.VsdAclTemplate{
		Name:              api.IngressAclTemplateName,
		DefaultAllowIP:    true,
		DefaultAllowNonIP: true,
		Active:            true,
		Priority:          0,
	}
	e := api.RESTError{}
	for {
		resp, err := nvsdc.session.Post(
			nvsdc.url+"domains/"+domainID+"/ingressacltemplates",
			&payload, &result, &e)
		if err != nil {
			glog.Error("Error when applying ingress acl template", err)
			return "", err
		}
		glog.Infoln("Got a reponse status", resp.Status(),
			"when creating ingress acl template")
		switch resp.Status() {
		case http.StatusCreated:
			nvsdc.ingressAclTemplateID = result[0].ID
			glog.Infoln("Applied default ingress ACL")
			err := nvsdc.CreateIngressAclEntries()
			if err != nil {
				return "", err
			}
			return nvsdc.ingressAclTemplateID, nil
		case http.StatusConflict:
			if e.InternalErrorCode == 2533 {
				ingressAclTemplate, err := nvsdc.GetIngressAclTemplate(domainID, payload.Name)
				if err != nil {
					return "", err
				}
				nvsdc.ingressAclTemplateID = ingressAclTemplate.ID
				glog.Infoln("Applied default ingress ACL")
				err = nvsdc.CreateIngressAclEntries()
				if err != nil {
					return "", err
				}
				return nvsdc.ingressAclTemplateID, nil
			} else {
				// Increment priority, and retry
				payload.Priority++
			}
		default:
			return "", VsdErrorResponse(resp, &e)
		}
	}
}

func (nvsdc *NuageVsdClient) CreateEgressAclTemplate(domainID string) (string, error) {
	result := make([]api.VsdAclTemplate, 1)
	payload := api.VsdAclTemplate{
		Name:              api.EgressAclTemplateName,
		DefaultAllowIP:    true,
		DefaultAllowNonIP: true,
		Active:            true,
		Priority:          0,
	}
	e := api.RESTError{}
	for {
		resp, err := nvsdc.session.Post(
			nvsdc.url+"domains/"+domainID+"/egressacltemplates",
			&payload, &result, &e)
		if err != nil {
			glog.Error("Error when applying egress acl template", err)
			return "", err
		}
		glog.Infoln("Got a reponse status", resp.Status(),
			"when creating egress acl template")
		switch resp.Status() {
		case http.StatusCreated:
			nvsdc.egressAclTemplateID = result[0].ID
			glog.Infoln("Applied default egress ACL")
			err := nvsdc.CreateEgressAclEntries()
			if err != nil {
				return "", err
			}
			return nvsdc.egressAclTemplateID, nil
		case http.StatusConflict:
			if e.InternalErrorCode == 2533 {
				egressAclTemplate, err := nvsdc.GetEgressAclTemplate(domainID, payload.Name)
				if err != nil {
					return "", err
				}
				nvsdc.egressAclTemplateID = egressAclTemplate.ID
				glog.Infoln("Applied default egress ACL")
				err = nvsdc.CreateEgressAclEntries()
				if err != nil {
					return "", err
				}
				return nvsdc.egressAclTemplateID, nil
			} else {
				// Increment priority, and retry
				payload.Priority++
			}
		default:
			return "", VsdErrorResponse(resp, &e)
		}
	}
}

func (nvsdc *NuageVsdClient) UpdateAclTemplate(aclTemplate *api.VsdAclTemplate, ingress bool) error {
	url := nvsdc.url + "egressacltemplates/" + aclTemplate.ID
	if ingress {
		url = nvsdc.url + "ingressacltemplates/" + aclTemplate.ID
	}
	e := api.RESTError{}
	resp, err := nvsdc.session.Put(
		url, aclTemplate, nil, &e)
	if err != nil || resp.Status() != http.StatusNoContent {
		VsdErrorResponse(resp, &e)
		return err
	}
	return nil
}

func (nvsdc *NuageVsdClient) GetAclEntryByPriority(ingress bool, aclEntryPriority int) (*api.VsdAclEntry, error) {
	result := make([]api.VsdAclEntry, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `priority == `+fmt.Sprintf("%v", aclEntryPriority))
	e := api.RESTError{}
	url := nvsdc.url + "egressacltemplates/" + nvsdc.egressAclTemplateID + "/egressaclentrytemplates"
	if ingress {
		url = nvsdc.url + "ingressacltemplates/" + nvsdc.ingressAclTemplateID + "/ingressaclentrytemplates"
	}
	resp, err := nvsdc.session.Get(url, nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting ACL entry with Priority %s: %d", err, aclEntryPriority)
		return nil, err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting ACL entry with priority", aclEntryPriority)
	if resp.Status() == http.StatusOK {
		// Status code 200 is returned even if there's no results.  If
		// the filter didn't match anything (or there was nothing to
		// return), the result object will just be empty.
		if result[0].Priority == aclEntryPriority {
			return &result[0], nil
		} else if result[0].Priority == 0 && result[0].ID == "" && result[0].Description == "" {
			return nil, errors.New("ACL entry not found")
		} else {
			return nil, errors.New(fmt.Sprintf(
				"Found %q instead of %q", result[0].Priority, aclEntryPriority))
		}
	} else {
		return nil, VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) GetAclEntry(ingress bool, aclEntry *api.VsdAclEntry) (*api.VsdAclEntry, error) {
	result := make([]api.VsdAclEntry, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", aclEntry.BuildFilter())
	glog.Infoln("Build filter is set to", aclEntry.BuildFilter())
	e := api.RESTError{}
	url := nvsdc.url + "egressacltemplates/" + nvsdc.egressAclTemplateID + "/egressaclentrytemplates"
	if ingress {
		url = nvsdc.url + "ingressacltemplates/" + nvsdc.ingressAclTemplateID + "/ingressaclentrytemplates"
	}
	resp, err := nvsdc.session.Get(url, nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting ACL entry %v: %s", aclEntry, err)
		return nil, err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting ACL entry: ", aclEntry)
	if resp.Status() == http.StatusOK {
		// Status code 200 is returned even if there's no results.  If
		// the filter didn't match anything (or there was nothing to
		// return), the result object will just be empty.
		if aclEntry.IsEqual(&result[0]) {
			return &result[0], nil
		} else if result[0].ID == "" {
			glog.Error("Acl Entry not found")
			return nil, errors.New("ACL entry not found")
		} else {
			glog.Error("Found an ACL entry that doesn't match the requested one")
			return nil, errors.New(fmt.Sprintf("Found ACL entry %v instead of %v", &result[0], aclEntry))
		}
	} else if resp.Status() == http.StatusNotFound {
		VsdErrorResponse(resp, &e)
		if ingress {
			aclTemplate, err := nvsdc.GetIngressAclTemplate(nvsdc.domainID, api.IngressAclTemplateName)
			if err != nil {
				glog.Error("Failed to fetch the ingress acl template ID from VSD")
				return nil, err
			}
			nvsdc.ingressAclTemplateID = aclTemplate.ID
			glog.Infoln("Refreshed ingress ACL template")
		} else {
			aclTemplate, err := nvsdc.GetEgressAclTemplate(nvsdc.domainID, api.EgressAclTemplateName)
			if err != nil {
				glog.Error("Failed to fetch the egress acl template ID from VSD")
				return nil, err
			}
			nvsdc.egressAclTemplateID = aclTemplate.ID
			glog.Infoln("Refreshed egress ACL template")
		}
		return nvsdc.GetAclEntry(ingress, aclEntry)
	} else {
		return nil, VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) CreateAclEntry(ingress bool, aclEntry *api.VsdAclEntry) (string, error) {
	//check if any entry matches the desired semantics with a different priority
	if acl, err := nvsdc.GetAclEntry(ingress, aclEntry); err == nil && acl != nil {
		return acl.ID, nil
	} else {
		result := make([]api.VsdObject, 1)
		e := api.RESTError{}
		url := nvsdc.url + "egressacltemplates/" + nvsdc.egressAclTemplateID + "/egressaclentrytemplates"
		if ingress {
			url = nvsdc.url + "ingressacltemplates/" + nvsdc.ingressAclTemplateID + "/ingressaclentrytemplates"
		}
		resp, err := nvsdc.session.Post(url+"?responseChoice=1", &aclEntry, &result, &e)
		if err != nil {
			glog.Error("Error when adding acl template entry", err)
			return "", err
		}
		glog.Infoln("Got a reponse status", resp.Status(),
			"when creating acl template entry")
		switch resp.Status() {
		case http.StatusCreated:
			glog.Infoln("Created ACL entry with priority: ", aclEntry.Priority)
			return result[0].ID, nil
		case http.StatusConflict:
			VsdErrorResponse(resp, &e)
			acl, err := nvsdc.GetAclEntryByPriority(ingress, aclEntry.Priority)
			if err != nil {
				return "", err
			}
			glog.Infoln("Applied ACL entry with priority: ", aclEntry.Priority)
			if aclEntry.IsEqual(acl) {
				return acl.ID, nil
			} else {
				aclEntry.TryNextAclPriority()
				return nvsdc.CreateAclEntry(ingress, aclEntry)
			}
		case http.StatusNotFound:
			VsdErrorResponse(resp, &e)
			if ingress {
				aclTemplate, err := nvsdc.GetIngressAclTemplate(nvsdc.domainID, api.IngressAclTemplateName)
				if err != nil {
					glog.Error("Failed to fetch the ingress acl template ID from VSD")
					return "", err
				}
				nvsdc.ingressAclTemplateID = aclTemplate.ID
				glog.Infoln("Refreshed ingress ACL template")
			} else {
				aclTemplate, err := nvsdc.GetEgressAclTemplate(nvsdc.domainID, api.EgressAclTemplateName)
				if err != nil {
					glog.Error("Failed to fetch the egress acl template ID from VSD")
					return "", err
				}
				nvsdc.egressAclTemplateID = aclTemplate.ID
				glog.Infoln("Refreshed egress ACL template")
			}
			return nvsdc.CreateAclEntry(ingress, aclEntry)
		default:
			return "", VsdErrorResponse(resp, &e)
		}
	}
}

func (nvsdc *NuageVsdClient) DeleteAclEntry(ingress bool, aclID string) error {
	// Delete subnets in this zone
	result := make([]struct{}, 1)
	e := api.RESTError{}
	url := nvsdc.url + "egressaclentrytemplates/" + aclID + "?responseChoice=1"
	if ingress {
		url = nvsdc.url + "ingressaclentrytemplates/" + aclID + "?responseChoice=1"
	}
	resp, err := nvsdc.session.Delete(url, nil, &result, &e)
	if err != nil {
		glog.Errorf("Error when deleting acl with ID %s: %s", aclID, err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when deleting acl")
	switch resp.Status() {
	case http.StatusNoContent:
		return nil
	default:
		return VsdErrorResponse(resp, &e)
	}
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
	if resp.Status() == http.StatusOK {
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
		return "", VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) CreateDomain(enterpriseID, domainTemplateID, name string) (string, error) {
	result := make([]api.VsdDomain, 1)
	payload := api.VsdDomain{
		Name:            name,
		Description:     "Auto-generated for OpenShift containers",
		TemplateID:      domainTemplateID,
		PATEnabled:      api.PATEnabled,
		UnderlayEnabled: api.PATEnabled,
	}
	e := api.RESTError{}
	resp, err := nvsdc.session.Post(nvsdc.url+"enterprises/"+enterpriseID+"/domains", &payload, &result, &e)
	if err != nil {
		glog.Error("Error when creating domain", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when creating domain")
	switch resp.Status() {
	case http.StatusCreated:
		glog.Infoln("Created the domain:", result[0].ID)
		return result[0].ID, nil
	case http.StatusConflict:
		//Domain already exists, call Get to retrieve the ID
		id, err := nvsdc.GetDomainID(enterpriseID, name)
		if err != nil {
			glog.Errorf("Error when getting domain ID: %s", err)
			return "", err
		} else {
			return id, nil
		}
	default:
		return "", VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) DeleteDomain(id string) error {
	result := make([]struct{}, 1)
	e := api.RESTError{}
	resp, err := nvsdc.session.Delete(nvsdc.url+"domains/"+id+"?responseChoice=1", nil, &result, &e)
	if err != nil {
		glog.Errorf("Error when deleting domain with ID %s: %s", id, err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when deleting domain")
	switch resp.Status() {
	case http.StatusNoContent:
		return nil
	default:
		return VsdErrorResponse(resp, &e)
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
	case http.StatusCreated:
		glog.Infoln("Created the zone:", result[0].ID)
		return result[0].ID, nil
	case http.StatusConflict:
		//Zone already exists, call Get to retrieve the ID
		id, err := nvsdc.GetZoneID(domainID, name)
		if err != nil {
			glog.Errorf("Error when getting zone ID: %s", err)
			return "", err
		} else {
			return id, nil
		}
	default:
		return "", VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) DeleteZone(id string) error {
	// Delete subnets in this zone
	result := make([]struct{}, 1)
	e := api.RESTError{}
	resp, err := nvsdc.session.Delete(nvsdc.url+"zones/"+id+"?responseChoice=1", nil, &result, &e)
	if err != nil {
		glog.Errorf("Error when deleting zone with ID %s: %s", id, err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when deleting zone")
	switch resp.Status() {
	case http.StatusNoContent:
		return nil
	default:
		return VsdErrorResponse(resp, &e)
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
	case http.StatusCreated:
		glog.Infoln("Created the subnet:", result[0].ID)
	case http.StatusConflict:
		glog.Infoln("Error from VSD:\n", e)
		// Subnet already exists, call Get to retrieve the ID
		if id, err := nvsdc.GetSubnetID(zoneID, name); err != nil {
			if e.InternalErrorCode == 2504 {
				// The network is overlapping with an existing one
				return "", errors.New("Overlapping Subnet")
			} else {
				glog.Errorf("Error when getting subnet ID: %s", err)
				return "", err
			}
		} else {
			return id, nil
		}
	default:
		return "", VsdErrorResponse(resp, &e)
	}
	return result[0].ID, nil
}

func (nvsdc *NuageVsdClient) DeleteSubnet(id string) error {
	result := make([]struct{}, 1)
	e := api.RESTError{}
	resp, err := nvsdc.session.Delete(nvsdc.url+"subnets/"+id+"?responseChoice=1", nil, &result, &e)
	if err != nil {
		glog.Errorf("Error when deleting subnet with ID %s: %s", id, err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when deleting subnet")
	if resp.Status() != http.StatusNoContent {
		return VsdErrorResponse(resp, &e)
	}
	return nil
}

func (nvsdc *NuageVsdClient) GetSubnet(zoneID, subnetName string) (*api.VsdSubnet, error) {
	result := make([]api.VsdSubnet, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `name == "`+subnetName+`"`)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"zones/"+zoneID+"/subnets", nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting subnet ID %s", err)
		return nil, err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting subnet ID")
	if resp.Status() == http.StatusOK {
		if result[0].Name == subnetName {
			return &result[0], nil
		} else {
			return nil, errors.New("Subnet not found")
		}
	} else {
		return nil, VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) GetSubnetID(zoneID, subnetName string) (string, error) {
	if vsdSubnet, err := nvsdc.GetSubnet(zoneID, subnetName); err == nil {
		return vsdSubnet.ID, nil
	} else {
		return "", err
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
	if resp.Status() == http.StatusOK {
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
		return "", VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) Run(nsChannel chan *api.NamespaceEvent, serviceChannel chan *api.ServiceEvent, stop chan bool) {
	//we will use the kube client APIs than interfacing with the REST API
	for {
		select {
		case nsEvent := <-nsChannel:
			nvsdc.HandleNsEvent(nsEvent)
		case serviceEvent := <-serviceChannel:
			nvsdc.HandleServiceEvent(serviceEvent)
		case nsUpdateReq := <-nvsdc.newSubnetQueue:
			switch nsUpdateReq.Event {
			case config.AddSubnet:
				nvsdc.CreateAdditionalSubnet(nsUpdateReq.NamespaceID)
			}
		}
	}
}

func (nvsdc *NuageVsdClient) CreateAdditionalSubnet(namespaceID string) {
	nvsdc.pods.editLock.Lock()
	defer nvsdc.pods.editLock.Unlock()
	namespace, exists := nvsdc.namespaces[namespaceID]
	if !exists {
		glog.Warningf("Got event to update namespace %q, but no namespace "+
			"found with that name", namespaceID)
		return
	}
	subnetName := namespace.Name + "-" + strconv.Itoa(namespace.numSubnets)
	glog.Infof("Creating subnet %q", subnetName)
	subnet, err := nvsdc.pool.Alloc(32 - nvsdc.subnetSize)
	if err != nil {
		glog.Warningf(
			"Error allocating new subnet for namespace %s: %s",
			namespace.Name, err.Error())
		return
	}
	for {
		subnetID, err := nvsdc.CreateSubnet(subnetName, namespace.ZoneID, subnet)
		if err == nil {
			// Step through the existing list to get to the tail
			subnetNode := namespace.Subnets
			for subnetNode.Next != nil {
				subnetNode = subnetNode.Next
			}
			// Then add the new subnet at the end of the list
			subnetNode.Next = &SubnetNode{
				SubnetID:   subnetID,
				Subnet:     subnet,
				SubnetName: subnetName,
				ActiveIPs:  0,
				Next:       nil,
			}
			namespace.numSubnets++
			namespace.NeedsNewSubnet = false
			nvsdc.namespaces[namespaceID] = namespace
			glog.Infof("Successfully created subnet %q", subnetName)
			return
		} else if err.Error() == "Overlapping Subnet" {
			// If this subnet overlaps with one in the VSD, don't return
			// it to the pool (since it won't work anyway), but get a
			// new subnet and retry.
			subnet, err = nvsdc.pool.Alloc(32 - nvsdc.subnetSize)
			if err != nil {
				glog.Warningf(
					"Error allocating new subnet for namespace %s: %s",
					namespace.Name, err.Error())
				return
			}
		} else {
			glog.Warningf(
				"Error allocating new subnet for namespace %s: %s",
				namespace.Name, err.Error())
			return
		}
	}
}

func (nvsdc *NuageVsdClient) HandleServiceEvent(serviceEvent *api.ServiceEvent) error {
	glog.Infoln("Received a service event: Service: ", serviceEvent)
	switch serviceEvent.Type {
	case api.Added:
		zone := serviceEvent.Namespace
		nmgID := ""
		err := errors.New("")
		exists := false
		if nmgID, exists = serviceEvent.NuageLabels[`network-macro-group.id`]; !exists {
			if nmgName, exists := serviceEvent.NuageLabels[`network-macro-group.name`]; exists {
				//use the label provided name to get network macro group ID and use that to create the network macro association
				nmgID, err = nvsdc.GetNetworkMacroGroupID(nvsdc.enterpriseID, nmgName)
				if err != nil {
					glog.Error("Label provided for network macro group name, but no network macro group identified", serviceEvent)
					return errors.New("Incorrect label information for creating service network macro")
				}
			}
		}
		if v, exists := serviceEvent.NuageLabels[`zone`]; exists {
			if _, exists = nvsdc.namespaces[v]; exists {
				if v != serviceEvent.Namespace {
					//label specified for a zone that is managed by nuagekubemon but for a different namespace
					glog.Errorf("Not authorized to create a service with zone label %v, in namespace %v", v, serviceEvent.Namespace)
					return errors.New("Incorrect label information for creating service network macro")
				}
			} else if nmgID == "" {
				// zone label is specified, but nuagekubemon doesn't manage this zone; and network macro group ID or Name are missing
				glog.Error("Label provided for a zone, but no network macro group identified", serviceEvent)
				return errors.New("Insufficient label information for creating service network macro")
			}
		}
		//default to using the validated zone's network macro group; if no specific labels are present.
		if nmgID == "" {
			nmgID = nvsdc.namespaces[zone].NetworkMacroGroupID
			//if we don't have a cached version, get the ID from the VSD
			if nmgID == "" {
				nmgID, err = nvsdc.GetNetworkMacroGroupID(nvsdc.enterpriseID, "Service Group For Zone - "+zone)
				if err != nil {
					glog.Error("Failed to get Network Macro Group ID: ", err)
				}
			}
		}
		networkMacro := &api.VsdNetworkMacro{
			Name:    `NetworkMacro for service ` + serviceEvent.Namespace + "--" + serviceEvent.Name,
			IPType:  "IPV4",
			Address: serviceEvent.ClusterIP,
			Netmask: "255.255.255.255",
		}
		networkMacroID, err := nvsdc.CreateNetworkMacro(nvsdc.enterpriseID, networkMacro)
		if err != nil {
			glog.Error("Error when creating the network macro for service", serviceEvent)
		} else {
			//add the network macro to the cached datastructure and also to the network macro group obtained via labels/default group
			nvsdc.namespaces[serviceEvent.Namespace].NetworkMacros[serviceEvent.Name] = networkMacroID
			err = nvsdc.AddNetworkMacroToNMG(networkMacroID, nmgID)
			if err != nil {
				glog.Error("Error when adding network macro to network macro group:", err)
			}
		}
	case api.Deleted:
		zone := serviceEvent.Namespace
		if _, exists := nvsdc.namespaces[zone]; exists {
			if nmID, exists := nvsdc.namespaces[zone].NetworkMacros[serviceEvent.Name]; exists {
				err := nvsdc.DeleteNetworkMacro(nmID)
				if err != nil {
					glog.Error("Error when deleting network macro with ID: ", nmID)
					return err
				} else {
					delete(nvsdc.namespaces[zone].NetworkMacros, nmID)
				}
			} else {
				glog.Warning("Could not retrieve network macro ID for the service that is being deleted", serviceEvent)
			}
		} else {
			glog.Warning("Could not retrieve namespace for the service that is being deleted", serviceEvent)
		}
	}
	return nil
}

func (nvsdc *NuageVsdClient) HandleNsEvent(nsEvent *api.NamespaceEvent) error {
	glog.Infoln("Received a namespace event: Namespace: ", nsEvent.Name, nsEvent.Type)
	switch nsEvent.Type {
	case api.Added:
		namespace, exists := nvsdc.namespaces[nsEvent.Name]
		if !exists {
			zoneID, err := nvsdc.CreateZone(nvsdc.domainID, nsEvent.Name)
			if err != nil {
				return err
			}
			namespace := NamespaceData{
				ZoneID:        zoneID,
				Name:          nsEvent.Name,
				NetworkMacros: make(map[string]string),
			}
			nvsdc.namespaces[nsEvent.Name] = namespace
			var subnetID string
			var subnet *IPv4Subnet
			// Check if subnet already exists
			subnetName := nsEvent.Name + "-0"
			vsdSubnet, err := nvsdc.GetSubnet(zoneID, subnetName)
			if err != nil {
				// If it didn't exist, allocate an appropriate subnet, then
				// create it
				if err.Error() == "Subnet not found" {
					// subnetSize is guaranteed to be between 0 and 32
					// (inclusive) by the Init() function defined above, so
					// (32 - subnetSize) will also produce a number between 0
					// and 32 (inclusive)
					subnet, err = nvsdc.pool.Alloc(32 - nvsdc.subnetSize)
					if err != nil {
						return err
					}
					for {
						subnetID, err = nvsdc.CreateSubnet(subnetName, zoneID,
							subnet)
						if err == nil {
							break
						} else if err.Error() == "Overlapping Subnet" {
							// If the error was "Overlapping Subnet," get a new
							// subnet and retry
							subnet, err = nvsdc.pool.Alloc(32 - nvsdc.subnetSize)
							if err != nil {
								return errors.New(
									"Error getting subnet for allocation: " +
										err.Error())
							}
						} else {
							// Only free the subnet if it wasn't overlapping so
							// that overlapping subnets are not retried
							nvsdc.pool.Free(subnet)
							return err
						}
					}
				} else {
					// There was an error, but it wasn't "Subnet not found".
					// Abort with that error.
					return err
				}
			} else {
				// If it exists, remove the subnet allocated to it from the pool
				subnet, err = IPv4SubnetFromAddrNetmask(vsdSubnet.Address,
					vsdSubnet.Netmask)
				if err != nil {
					return err
				}
				err = nvsdc.pool.AllocSpecific(subnet)
				if err != nil {
					if !nvsdc.clusterNetwork.Contains(subnet) {
						glog.Error(subnet.String(),
							" is not a member of clusterNetworkCIDR ",
							nvsdc.clusterNetwork.String())
					}
				}
				subnetID = vsdSubnet.ID
			}
			namespace.Subnets = &SubnetNode{
				SubnetID:   subnetID,
				SubnetName: subnetName,
				Subnet:     subnet,
				Next:       nil,
			}
			namespace.numSubnets++
			nvsdc.namespaces[nsEvent.Name] = namespace
			if nsEvent.Name == "default" {
				err = nvsdc.CreateDefaultZoneAcls(zoneID)
				if err != nil {
					glog.Error("Got an error when creating default zone's ACL entries")
					return err
				}
			} else {
				err = nvsdc.CreateSpecificZoneAcls(nsEvent.Name, zoneID)
				if err != nil {
					glog.Error("Got an error when creating zone specific ACLs: ", nsEvent.Name)
					return err
				}
			}
			return nil
		}
		// else (nvsdc.namespaces[nsEvent.Name] exists)
		id, err := nvsdc.GetZoneID(nvsdc.domainID, nsEvent.Name)
		switch {
		case id == "" && err == nil:
			err = errors.New("Invalid zone ID returned")
			fallthrough
		case err != nil:
			glog.Errorf("Invalid ID for zone %s", nsEvent.Name)
			return err
		case id != "" && err == nil:
			if nsEvent.Name == "default" {
				err = nvsdc.CreateDefaultZoneAcls(id)
				if err != nil {
					glog.Error("Got an error when creating default zone's ACL entries")
					return err
				}
			} else {
				err = nvsdc.CreateSpecificZoneAcls(nsEvent.Name, id)
				if err != nil {
					glog.Error("Got an error when creating zone specific ACLs: ", nsEvent.Name)
					return err
				}
			}
			namespace.ZoneID = id
			if namespace.NetworkMacros == nil {
				namespace.NetworkMacros = make(map[string]string)
			}
			return nil
		}
	case api.Deleted:
		if zone, exists := nvsdc.namespaces[nsEvent.Name]; exists {
			// Delete subnets that we've created, and free them back into the pool
			if nsEvent.Name == "default" {
				err := nvsdc.DeleteDefaultZoneAcls(zone.ZoneID)
				if err != nil {
					// Log the error, but continue to delete subnets/zone
					glog.Error("Got an error when deleting default zone's ACL entries")
				}
			} else {
				err := nvsdc.DeleteSpecificZoneAcls(nsEvent.Name)
				if err != nil {
					// Log the error, but continue to delete subnets/zone
					glog.Error("Got an error when deleting network macro group for zone: ", nsEvent.Name)
				}
			}
			if subnetsHead := zone.Subnets; subnetsHead != nil {
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
				// with this zone (remove the reference to the list so it can be
				// garbage collected)
				zone.Subnets = nil
				zone.numSubnets = 0
			}
			delete(nvsdc.namespaces, nsEvent.Name)
			return nvsdc.DeleteZone(zone.ZoneID)
		}
		id, err := nvsdc.GetZoneID(nvsdc.domainID, nsEvent.Name)
		switch {
		case id == "" && err == nil:
			glog.Warningf("Got delete namespace event for non-existant zone %s", nsEvent.Name)
			return nil
		case err != nil:
			glog.Errorf("Error getting ID of zone %s", nsEvent.Name)
			return err
		case id != "" && err == nil:
			glog.Infof("Deleting zone %s which was not found locally", nsEvent.Name)
			if nsEvent.Name == "default" {
				err = nvsdc.DeleteDefaultZoneAcls(id)
				if err != nil {
					// Log the error, but continue to delete subnets/zone
					glog.Error("Got an error when deleting default zone's ACL entries")
				}
			} else {
				err = nvsdc.DeleteSpecificZoneAcls(nsEvent.Name)
				if err != nil {
					// Log the error, but continue to delete subnets/zone
					glog.Error("Got an error when deleting network macro group for zone", nsEvent.Name)
				}
			}
			return nvsdc.DeleteZone(id)
		}
	}
	return nil
}

func (nvsdc *NuageVsdClient) CreateDefaultZoneAcls(zoneID string) error {
	nmgid, err := nvsdc.CreateNetworkMacroGroup(nvsdc.enterpriseID, "default")
	if err != nil {
		glog.Error("Error when creating the network macro group for zone", "default")
		return err
	} else {
		if nsd, exists := nvsdc.namespaces["default"]; exists {
			nsd.NetworkMacroGroupID = nmgid
			nvsdc.namespaces["default"] = nsd
		} else {
			nvsdc.namespaces["default"] = NamespaceData{
				ZoneID:              zoneID,
				Name:                "default",
				NetworkMacroGroupID: nmgid,
				NetworkMacros:       make(map[string]string),
			}
		}
	}
	//add ingress and egress ACL entries for allowing zone to default zone communication
	aclEntry := api.VsdAclEntry{
		Action:       "FORWARD",
		DSCP:         "*",
		Description:  "Allow Traffic Between All Zones and Default Zone",
		EntityScope:  "ENTERPRISE",
		EtherType:    "0x0800",
		LocationID:   "",
		LocationType: "ANY",
		NetworkType:  "NETWORK_MACRO_GROUP",
		NetworkID:    nmgid,
		PolicyState:  "LIVE",
		Priority:     1,
		Protocol:     "ANY",
		Reflexive:    false,
	}
	_, err = nvsdc.CreateAclEntry(true, &aclEntry)
	if err != nil {
		glog.Error("Error when creating the ACL rules for the default zone")
		return err
	}
	_, err = nvsdc.CreateAclEntry(false, &aclEntry)
	if err != nil {
		glog.Error("Error when creating the ACL rules for the default zone")
		return err
	}
	//default to any ACL rule
	aclEntry.LocationID = zoneID
	aclEntry.LocationType = "ZONE"
	aclEntry.NetworkType = "ANY"
	aclEntry.NetworkID = ""
	aclEntry.Priority = 2
	_, err = nvsdc.CreateAclEntry(true, &aclEntry)
	if err != nil {
		glog.Error("Error when creating the ACL rules for the default zone")
		return err
	}
	_, err = nvsdc.CreateAclEntry(false, &aclEntry)
	if err != nil {
		glog.Error("Error when creating the ACL rules for the default zone")
		return err
	}
	return nil
}

func (nvsdc *NuageVsdClient) CreateSpecificZoneAcls(zoneName string, zoneID string) error {
	//first create the network macro group for the zone.
	nmgid, err := nvsdc.CreateNetworkMacroGroup(nvsdc.enterpriseID, zoneName)
	if err != nil {
		glog.Error("Error when creating the network macro group for zone", zoneName)
		return err
	} else {
		if nsd, exists := nvsdc.namespaces[zoneName]; exists {
			nsd.NetworkMacroGroupID = nmgid
			nvsdc.namespaces[zoneName] = nsd
		} else {
			nvsdc.namespaces[zoneName] = NamespaceData{
				ZoneID:              zoneID,
				Name:                zoneName,
				NetworkMacroGroupID: nmgid,
				NetworkMacros:       make(map[string]string),
			}
		}
	}
	//add ingress and egress ACL entries for allowing zone to default zone communication
	aclEntry := api.VsdAclEntry{
		Action:       "FORWARD",
		DSCP:         "*",
		Description:  "Allow Traffic Between Zone - " + zoneName + " And Its Services",
		EntityScope:  "ENTERPRISE",
		EtherType:    "0x0800",
		LocationID:   nvsdc.namespaces[zoneName].ZoneID,
		LocationType: "ZONE",
		NetworkID:    nmgid,
		NetworkType:  "NETWORK_MACRO_GROUP",
		PolicyState:  "LIVE",
		Priority:     300 + nvsdc.NextAvailablePriority(),
		Protocol:     "ANY",
		Reflexive:    false,
	}
	_, err = nvsdc.CreateAclEntry(true, &aclEntry)
	if err != nil {
		glog.Error("Error when creating the ACL rules for the zone: ", zoneName)
		return err
	} else {
		nvsdc.SetNextAvailablePriority(aclEntry.Priority + 1 - 300)
	}
	_, err = nvsdc.CreateAclEntry(false, &aclEntry)
	if err != nil {
		glog.Error("Error when creating the ACL rules for the zone: ", zoneName)
		return err
	} else {
		nvsdc.SetNextAvailablePriority(aclEntry.Priority + 1 - 300)
	}
	return nil
}

func (nvsdc *NuageVsdClient) NextAvailablePriority() int {
	defer nvsdc.IncrementNextAvailablePriority()
	return nvsdc.nextAvailablePriority
}

func (nvsdc *NuageVsdClient) IncrementNextAvailablePriority() {
	nvsdc.nextAvailablePriority++
}

func (nvsdc *NuageVsdClient) SetNextAvailablePriority(val int) {
	nvsdc.nextAvailablePriority = val
}

func (nvsdc *NuageVsdClient) CreateNetworkMacroGroup(enterpriseID string, zoneName string) (string, error) {
	result := make([]api.VsdObject, 1)
	payload := api.VsdObject{
		Name:        "Service Group For Zone - " + zoneName,
		Description: "Auto-generated network macro group for zone - " + zoneName,
	}
	e := api.RESTError{}
	resp, err := nvsdc.session.Post(nvsdc.url+"enterprises/"+enterpriseID+"/networkmacrogroups", &payload, &result, &e)
	if err != nil {
		glog.Error("Error when creating network macro group for zone: ", zoneName, err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when creating network macro group")
	switch resp.Status() {
	case http.StatusCreated:
		return result[0].ID, nil
	case http.StatusConflict:
		//Network Macro Group already exists, call Get to retrieve the ID
		id, err := nvsdc.GetNetworkMacroGroupID(enterpriseID, payload.Name)
		if err != nil {
			glog.Errorf("Error when getting network macro group ID for zone: %s - %s", zoneName, err)
			return "", err
		}
		return id, nil
	default:
		return "", VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) GetNetworkMacroGroupID(enterpriseID, nmgName string) (string, error) {
	result := make([]api.VsdObject, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `name == "`+nmgName+`"`)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"enterprises/"+enterpriseID+"/networkmacrogroups", nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting network macro group ID with name: %s - %s", nmgName, err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting network macro group ID")
	if resp.Status() == http.StatusOK {
		// Status code 200 is returned even if there's no results.  If
		// the filter didn't match anything (or there was nothing to
		// return), the result object will just be empty.
		if result[0].Name == nmgName {
			return result[0].ID, nil
		} else if result[0].Name == "" {
			return "", errors.New("Network Macro Group not found")
		} else {
			return "", errors.New(fmt.Sprintf(
				"Found %q instead of %q", result[0].Name, nmgName))
		}
	} else {
		return "", VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) DeleteNetworkMacroGroup(networkMacroGroupID string) error {
	// Delete network macro group
	result := make([]struct{}, 1)
	e := api.RESTError{}
	url := nvsdc.url + "networkmacrogroups/" + networkMacroGroupID + "?responseChoice=1"
	resp, err := nvsdc.session.Delete(url, nil, &result, &e)
	if err != nil {
		glog.Errorf("Error when deleting network macro group with ID %s: %s", networkMacroGroupID, err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when deleting network macro group")
	switch resp.Status() {
	case http.StatusNoContent:
		return nil
	default:
		return VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) DeleteSpecificZoneAcls(zoneName string) error {
	//add ingress and egress ACL entries for allowing zone to default zone communication
	// aclEntry := api.VsdAclEntry{
	// 	Action:       "FORWARD",
	// 	DSCP:         "*",
	// 	Description:  "Allow Traffic Between Zone - " + zoneName + " And Its Services",
	// 	EntityScope:  "ENTERPRISE",
	// 	EtherType:    "0x0800",
	// 	LocationID:   nvsdc.namespaces[zoneName].ZoneID,
	// 	LocationType: "ZONE",
	// 	NetworkID:    nvsdc.namespaces[zoneName].NetworkMacroGroupID,
	// 	NetworkType:  "NETWORK_MACRO_GROUP",
	// 	PolicyState:  "LIVE",
	// 	Protocol:     "ANY",
	// 	Reflexive:    false,
	// }
	// if acl, err := nvsdc.GetAclEntry(true, &aclEntry); err == nil && acl != nil {
	// 	err = nvsdc.DeleteAclEntry(true, acl.ID)
	// 	if err != nil {
	// 		glog.Error("Error when deleting the ingress ACL rules for the zone: ", zoneName, aclEntry)
	// 		return err
	// 	}
	// } else {
	// 	glog.Error("Failed to get ingress acl entry to delete", aclEntry)
	// 	return err
	// }
	// if acl, err := nvsdc.GetAclEntry(false, &aclEntry); err == nil && acl != nil {
	// 	err = nvsdc.DeleteAclEntry(false, acl.ID)
	// 	if err != nil {
	// 		glog.Error("Error when deleting the egress ACL rules for the zone: ", zoneName, aclEntry)
	// 		return err
	// 	}
	// } else {
	// 	glog.Error("Failed to get egress acl entry to delete", aclEntry)
	// 	return err
	// }
	glog.Info("Looking up zone specific network macro group")
	if nvsdc.namespaces[zoneName].NetworkMacroGroupID != "" {
		glog.Infof("Found zone specific network macro group with ID: %s for zone name: %s", nvsdc.namespaces[zoneName].NetworkMacroGroupID, zoneName)
		err := nvsdc.DeleteNetworkMacroGroup(nvsdc.namespaces[zoneName].NetworkMacroGroupID)
		if err != nil {
			glog.Error("Failed to delete network macro group for zone: ", zoneName)
			return err
		} else {
			glog.Infof("Deleted network macro group with ID: %s for zone name: %s", nvsdc.namespaces[zoneName].NetworkMacroGroupID, zoneName)
			if nsd, exists := nvsdc.namespaces[zoneName]; exists {
				nsd.NetworkMacroGroupID = ""
				nvsdc.namespaces[zoneName] = nsd
			}
		}
	}
	glog.Info("Succeeded in deleting the network macro group")
	return nil
}

func (nvsdc *NuageVsdClient) DeleteDefaultZoneAcls(zoneID string) error {
	// aclEntry := api.VsdAclEntry{
	// 	Action:       "FORWARD",
	// 	DSCP:         "*",
	// 	Description:  "Allow Traffic Between All Zones and Default Zone",
	// 	EntityScope:  "ENTERPRISE",
	// 	EtherType:    "0x0800",
	// 	LocationID:   "",
	// 	LocationType: "ANY",
	// 	NetworkID:    nvsdc.namespaces["default"].NetworkMacroGroupID,
	// 	NetworkType:  "NETWORK_MACRO_GROUP",
	// 	PolicyState:  "LIVE",
	// 	Protocol:     "ANY",
	// 	Reflexive:    false,
	// }
	// if acl, err := nvsdc.GetAclEntry(true, &aclEntry); err == nil && acl != nil {
	// 	err = nvsdc.DeleteAclEntry(true, acl.ID)
	// 	if err != nil {
	// 		glog.Error("Error when deleting the ingress ACL rules for the default zone", aclEntry)
	// 		return err
	// 	}
	// } else {
	// 	glog.Error("Failed to get ingress acl entry to delete", aclEntry)
	// 	return err
	// }
	// if acl, err := nvsdc.GetAclEntry(false, &aclEntry); err == nil && acl != nil {
	// 	err = nvsdc.DeleteAclEntry(false, acl.ID)
	// 	if err != nil {
	// 		glog.Error("Error when deleting the egress ACL rules for the default zone", aclEntry)
	// 		return err
	// 	}
	// } else {
	// 	glog.Error("Failed to get egress acl entry to delete", aclEntry)
	// 	return err
	// }
	if nvsdc.namespaces["default"].NetworkMacroGroupID != "" {
		err := nvsdc.DeleteNetworkMacroGroup(nvsdc.namespaces["default"].NetworkMacroGroupID)
		if err != nil {
			glog.Error("Failed to delete network macro group for default zone")
			return err
		} else {
			if nsd, exists := nvsdc.namespaces["default"]; exists {
				nsd.NetworkMacroGroupID = ""
				nvsdc.namespaces["default"] = nsd
			}
		}
	}
	return nil
}

func (nvsdc *NuageVsdClient) CreateNetworkMacro(enterpriseID string, networkMacro *api.VsdNetworkMacro) (string, error) {
	result := make([]api.VsdNetworkMacro, 1)
	e := api.RESTError{}
	resp, err := nvsdc.session.Post(nvsdc.url+"enterprises/"+enterpriseID+"/enterprisenetworks", networkMacro, &result, &e)
	if err != nil {
		glog.Error("Error when creating network macro", networkMacro, err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when creating network macro")
	switch resp.Status() {
	case http.StatusCreated:
		return result[0].ID, nil
	case http.StatusConflict:
		//Network Macro already exists, call Get to retrieve the ID
		id, err := nvsdc.GetNetworkMacroID(enterpriseID, networkMacro)
		if err != nil {
			glog.Errorf("Error when getting network macro ID: %v - %v", networkMacro, err)
			return "", err
		}
		return id, nil
	default:
		return "", VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) GetNetworkMacroID(enterpriseID string, networkMacro *api.VsdNetworkMacro) (string, error) {
	result := make([]api.VsdNetworkMacro, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `name == "`+networkMacro.Name+`" and IPType =="`+networkMacro.IPType+`" and address == "`+networkMacro.Address+
		`" and netmask == "`+networkMacro.Netmask+`"`)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"enterprises/"+enterpriseID+"/enterprisenetworks", nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting network macro ID for network macro: %v - %v", networkMacro, err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting network macro ID")
	if resp.Status() == http.StatusOK {
		// Status code 200 is returned even if there's no results.  If
		// the filter didn't match anything (or there was nothing to
		// return), the result object will just be empty.
		if result[0].Name == networkMacro.Name {
			return result[0].ID, nil
		} else if result[0].Name == "" {
			return "", errors.New("Network Macro not found")
		} else {
			return "", errors.New(fmt.Sprintf(
				"Found %q instead of %q", result[0].Name, networkMacro.Name))
		}
	} else {
		return "", VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) DeleteNetworkMacro(networkMacroID string) error {
	// Delete network macro
	result := make([]struct{}, 1)
	e := api.RESTError{}
	url := nvsdc.url + "enterprisenetworks/" + networkMacroID + "?responseChoice=1"
	resp, err := nvsdc.session.Delete(url, nil, &result, &e)
	if err != nil {
		glog.Errorf("Error when deleting network macro with ID %s: %s", networkMacroID, err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when deleting network macro")
	switch resp.Status() {
	case http.StatusNoContent:
		return nil
	default:
		return VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) AddUserToGroup(userID, groupID string) error {
	result := make([]api.VsdUser, 0, 100)
	e := api.RESTError{}
	nvsdc.session.Header.Add("X-Nuage-PageSize", "100")
	page := 0
	nvsdc.session.Header.Add("X-Nuage-Page", strconv.Itoa(page))
	// guarantee that the headers are cleared so that we don't change the
	// behavior of other functions
	defer nvsdc.session.Header.Del("X-Nuage-PageSize")
	defer nvsdc.session.Header.Del("X-Nuage-Page")
	userIDList := []string{userID}
	for {
		resp, err := nvsdc.session.Get(nvsdc.url+"groups/"+groupID+"/users",
			nil, &result, &e)
		if err != nil {
			glog.Errorf("Error when adding user %s to group %s: %s", userID,
				groupID, err)
			return err
		}
		if resp.Status() == http.StatusNoContent || resp.HttpResponse().Header.Get("x-nuage-count") == "0" {
			break
		} else if resp.Status() == http.StatusOK {
			// Add all the items on this page to the list
			for _, user := range result {
				if user.ID == userID {
					// The user is already a part of the intended group, so no
					// extra work is necessary.
					return nil
				}
				userIDList = append(userIDList, user.ID)
			}
			// If there's less than 100 items in the page, we must've reached
			// the last page.  Break here instead of getting the next
			// (guaranteed empty) page.
			if count, err := strconv.Atoi(resp.HttpResponse().Header.Get("x-nuage-count")); err == nil {
				if count < 100 {
					break
				}
			} else {
				// Something went wrong with parsing the x-nuage-count header
				return errors.New("Invalid X-Nuage-Count: " + err.Error())
			}
			// Update headers to get the next page
			page++
			nvsdc.session.Header.Set("X-Nuage-Page", strconv.Itoa(page))
		} else {
			// Something went wrong
			return VsdErrorResponse(resp, &e)
		}
	}
	// Delete headers.  Calling Header.Del(...) on a non-existent header is a
	// no-op, so the `defer ...Header.Del(...)` calls above are still valid.
	nvsdc.session.Header.Del("X-Nuage-PageSize")
	nvsdc.session.Header.Del("X-Nuage-Page")
	resp, err := nvsdc.session.Put(nvsdc.url+"groups/"+
		groupID+"/users", &userIDList, nil, &e)
	if err != nil {
		glog.Errorf("Error when adding user %s to group %s: %s", userID,
			groupID, err)
		return err
	} else {
		glog.Infoln("Got a reponse status", resp.Status(),
			"when adding user to group")
		switch resp.Status() {
		case http.StatusNoContent:
			glog.Infof("Added user %s to group %s", userID, groupID)
		default:
			return VsdErrorResponse(resp, &e)
		}
	}
	return nil
}

func (nvsdc *NuageVsdClient) AddNetworkMacroToNMG(networkMacroID, networkMacroGroupID string) error {
	result := make([]api.VsdObject, 0, 100)
	e := api.RESTError{}
	nvsdc.session.Header.Add("X-Nuage-PageSize", "100")
	page := 0
	nvsdc.session.Header.Add("X-Nuage-Page", strconv.Itoa(page))
	// guarantee that the headers are cleared so that we don't change the
	// behavior of other functions
	defer nvsdc.session.Header.Del("X-Nuage-PageSize")
	defer nvsdc.session.Header.Del("X-Nuage-Page")
	networkMacroIDList := []string{networkMacroID}
	for {
		resp, err := nvsdc.session.Get(nvsdc.url+"networkmacrogroups/"+
			networkMacroGroupID+"/enterprisenetworks", nil, &result, &e)
		if err != nil {
			glog.Errorf("Error when adding network macro with ID %s: %s", networkMacroID, err)
			return err
		}
		// Using if...else here instead of switch because you can't use 'break'
		// inside the switch to break from the infinite for-loop
		if resp.Status() == http.StatusNoContent || resp.HttpResponse().Header.Get("x-nuage-count") == "0" {
			break
		} else if resp.Status() == http.StatusOK {
			// The response contains a list of network macros.  Add them to the
			// list
			for _, networkMacro := range result {
				if networkMacro.ID == networkMacroID {
					// The network macro we're trying to add already exists.  No
					// REST call is necessary.
					return nil
				}
				networkMacroIDList = append(networkMacroIDList, networkMacro.ID)
			}
			// Increment the page number for the next call
			page++
			nvsdc.session.Header.Set("X-Nuage-Page", strconv.Itoa(page))
		} else {
			// Something went wrong
			return VsdErrorResponse(resp, &e)
		}
	}
	nvsdc.session.Header.Del("X-Nuage-PageSize")
	nvsdc.session.Header.Del("X-Nuage-Page")
	resp, err := nvsdc.session.Put(nvsdc.url+"networkmacrogroups/"+
		networkMacroGroupID+"/enterprisenetworks", &networkMacroIDList, nil, &e)
	if err != nil {
		glog.Error("Error when adding network macro to the network macro group", err)
		return err
	} else {
		glog.Infoln("Got a reponse status", resp.Status(),
			"when adding network macro to the network macro group")
		switch resp.Status() {
		case http.StatusNoContent:
			glog.Infoln("Added the network macro to the network macro group")
		default:
			return VsdErrorResponse(resp, &e)
		}
	}
	return nil
}

func VsdErrorResponse(resp *napping.Response, e *api.RESTError) error {
	glog.Errorln("Bad response status from VSD Server")
	glog.Errorf("\t Raw Text:\n%v\n", resp.RawText())
	glog.Errorf("\t Status:  %v\n", resp.Status())
	glog.Errorf("\t VSD Error: %v\n", e)
	return errors.New("Unexpected error code: " + fmt.Sprintf("%v", resp.Status()))
}
