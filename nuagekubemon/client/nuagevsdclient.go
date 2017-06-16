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
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/jmcvetta/napping"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/config"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/policy"
	"github.com/nuagenetworks/vspk-go/vspk"
	"github.com/rfredette/sleepy"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type NuageVsdClient struct {
	url                                string
	version                            string
	session                            napping.Session
	enterpriseID                       string
	domainID                           string
	namespaces                         map[string]NamespaceData //namespace name -> namespace data
	pods                               *PodList                 //<namespace>/<pod-name> -> subnet
	pool                               IPv4SubnetPool
	clusterNetwork                     *IPv4Subnet //clusterNetworkCIDR used to generate pool
	serviceNetwork                     *IPv4Subnet
	ingressAclTemplateID               string
	egressAclTemplateID                string
	ingressAclTemplateZoneAnnotationID string
	egressAclTemplateZoneAnnotationID  string
	nextAvailablePriority              int
	subnetSize                         int         //the size in bits of the subnets we allocate (i.e. size 8 produces /24 subnets).
	restAPI                            *sleepy.API //TODO: split the rest server into its own package
	restServer                         *http.Server
	newSubnetQueue                     chan config.NamespaceUpdateRequest //list of namespaces that need new subnets
	privilegedProjectName              string
	resourceManager                    *policy.ResourceManager
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

func NewNuageVsdClient(nkmConfig *config.NuageKubeMonConfig, clusterCallBacks *api.ClusterClientCallBacks) *NuageVsdClient {
	nvsdc := new(NuageVsdClient)
	nvsdc.Init(nkmConfig, clusterCallBacks)
	return nvsdc
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

func (nvsdc *NuageVsdClient) CreateSession(userCertFile string, userKeyFile string) {

	cert, err := tls.LoadX509KeyPair(userCertFile, userKeyFile)
	if err != nil {
		glog.Errorf("Error loading VSD generated certificates to authenticate with VSD %s", err)
	}

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()

	nvsdc.session = napping.Session{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		},
		Header: &http.Header{},
	}

	nvsdc.session.Header.Add("Content-Type", "application/json")
	// Request that the TCP connection is closed when the transaction is
	// complete
	nvsdc.session.Header.Add("Connection", "close")
}

func (nvsdc *NuageVsdClient) Init(nkmConfig *config.NuageKubeMonConfig, clusterCallBacks *api.ClusterClientCallBacks) {
	cb := &policy.CallBacks{
		AddPg:             nvsdc.CreatePolicyGroup,
		DeletePg:          nvsdc.DeletePolicyGroup,
		AddPortsToPg:      nvsdc.AddPodsToPolicyGroup,
		DeletePortsFromPg: nvsdc.RemovePortsFromPolicyGroup,
	}
	var err error
	nvsdc.version = nkmConfig.NuageVspVersion
	nvsdc.url = nkmConfig.NuageVsdApiUrl + "/nuage/api/" + nvsdc.version + "/"
	nvsdc.privilegedProjectName = nkmConfig.PrivilegedProject
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

	//initialize the resource manager
	vsdMeta := make(policy.VsdMetaData)
	vsdMeta["domainName"] = nkmConfig.DomainName
	vsdMeta["enterpriseName"] = nkmConfig.EnterpriseName
	vsdMeta["vsdUrl"] = nkmConfig.NuageVsdApiUrl
	vsdMeta["usercertfile"] = nkmConfig.UserCertificateFile
	vsdMeta["userkeyfile"] = nkmConfig.UserKeyFile
	rm, err := policy.NewResourceManager(cb, clusterCallBacks, &vsdMeta)
	if err != nil {
		glog.Error("Failed to initialize the resource manager properly")
	} else {
		nvsdc.resourceManager = rm
	}

	nvsdc.pods = NewPodList(nvsdc.namespaces, nvsdc.newSubnetQueue, nvsdc.resourceManager.GetPolicyGroupsForPod)

	nvsdc.CreateSession(nkmConfig.UserCertificateFile, nkmConfig.UserKeyFile)
	nvsdc.nextAvailablePriority = 0

	for {
		nvsdc.enterpriseID, err = nvsdc.GetEnterpriseID(nkmConfig.EnterpriseName)
		if err != nil {
			glog.Error("Failed to fetch enterprise ID from VSD. Will retry in 10 seconds")
		} else {
			break
		}
		time.Sleep(time.Duration(10) * time.Second)
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

	err = nvsdc.CreateIngressAclEntries()
	if err != nil {
		glog.Fatal(err)
	}

	_, err = nvsdc.CreateEgressAclTemplate(nvsdc.domainID)
	if err != nil {
		glog.Fatal(err)
	}

	err = nvsdc.CreateEgressAclEntries()
	if err != nil {
		glog.Fatal(err)
	}

	_, err = nvsdc.CreateIngressAclTemplateForNamespaceAnnotations(nvsdc.domainID)
	if err != nil {
		glog.Fatal(err)
	}

	_, err = nvsdc.CreateEgressAclTemplateForNamespaceAnnotations(nvsdc.domainID)
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

	glog.Infof("Getting ACL template by ID %s using URL: %s", templateID, url)

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

func (nvsdc *NuageVsdClient) GetAclTemplateID(domainID, name string, ingress bool, priority int) (string, error) {
	result := make([]api.VsdAclTemplate, 1)
	e := api.RESTError{}

	restpath := "/ingressacltemplates"
	if !ingress {
		restpath = "/egressacltemplates"
	}

	resp, err := nvsdc.session.Get(nvsdc.url+"domains/"+domainID+restpath, nil, &result, &e)

	if err != nil {
		glog.Errorf("Error when getting ACL template ID %s", err)
		return "", err
	}

	glog.Infoln("Result for ACL template obtained from VSD when getting ACL template ID: ", result)

	glog.Infoln("Got a reponse status", resp.Status(), "when getting ACL template ID")
	if resp.Status() == http.StatusOK {
		// Status code 200 is returned even if there's no results.  If
		// the filter didn't match anything (or there was nothing to
		// return), the result object will just be empty.
		for index := range result {
			if result[index].Name == name && result[index].Priority == priority && result[index].Active == true {
				glog.Infof("Found active VSD ACL template with priority %d and name %s", priority, name)
				glog.Infof("Active VSD ACL template ID is %s", result[index].ID)
				return result[index].ID, nil
			}
		}
		glog.Errorf("Could not find already existing VSD ACL active template with priority %d and name %s", priority, name)
		return "", errors.New("Active ACL template not found")
	} else {
		return "", VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) CreateAclTemplate(domainID string, name string, priority int, ingress bool) (string, error) {
	result := make([]api.VsdAclTemplate, 1)
	payload := api.VsdAclTemplate{
		Name:              name,
		DefaultAllowIP:    true,
		DefaultAllowNonIP: true,
		Active:            true,
		Priority:          priority,
	}
	e := api.RESTError{}

	restpath := "/ingressacltemplates"
	if !ingress {
		restpath = "/egressacltemplates"
	}

	for {
		id, err := nvsdc.GetAclTemplateID(domainID, name, ingress, priority)
		if err != nil {
			glog.Errorf("Error when ACL template domain ID: %s", err)
		} else {
			return id, nil
		}

		resp, err := nvsdc.session.Post(
			nvsdc.url+"domains/"+domainID+restpath,
			&payload, &result, &e)
		if err != nil {
			glog.Errorf("Error %s when creating ACL template %s", err, name)
			return "", err
		}
		glog.Infoln("Got a reponse status", resp.Status(),
			"when creating acl template")
		switch resp.Status() {
		case http.StatusCreated:
			glog.Infof("Created ACL template %s with priority %d", name, priority)
			return result[0].ID, nil
		case http.StatusConflict:
			if e.InternalErrorCode == 2533 {
				var aclTemplate *api.VsdAclTemplate
				var err error
				if ingress {
					aclTemplate, err = nvsdc.GetIngressAclTemplate(domainID, payload.Name)
				} else {
					aclTemplate, err = nvsdc.GetEgressAclTemplate(domainID, payload.Name)
				}

				if err != nil {
					return "", err
				}
				glog.Infof("Created ACL template %s with priority %d", name, priority)
				return aclTemplate.ID, nil
			} else {
				// Increment priority, and retry
				payload.Priority--
			}
		default:
			return "", VsdErrorResponse(resp, &e)
		}
	}
}
func (nvsdc *NuageVsdClient) CreateIngressAclTemplate(domainID string) (string, error) {
	id, err := nvsdc.CreateAclTemplate(domainID, api.IngressAclTemplateName, api.MAX_VSD_ACL_PRIORITY, true)
	if err != nil {
		return "", err
	}
	nvsdc.ingressAclTemplateID = id
	return id, nil
}

func (nvsdc *NuageVsdClient) CreateEgressAclTemplate(domainID string) (string, error) {
	id, err := nvsdc.CreateAclTemplate(domainID, api.EgressAclTemplateName, api.MAX_VSD_ACL_PRIORITY, false)
	if err != nil {
		return "", err
	}
	nvsdc.egressAclTemplateID = id
	return id, nil
}

func (nvsdc *NuageVsdClient) CreateIngressAclTemplateForNamespaceAnnotations(domainID string) (string, error) {
	id, err := nvsdc.CreateAclTemplate(domainID, api.ZoneAnnotationTemplateName, api.MAX_VSD_ACL_PRIORITY-1, true)
	if err != nil {
		return "", err
	}
	nvsdc.ingressAclTemplateZoneAnnotationID = id
	return id, nil
}

func (nvsdc *NuageVsdClient) CreateEgressAclTemplateForNamespaceAnnotations(domainID string) (string, error) {
	id, err := nvsdc.CreateAclTemplate(domainID, api.ZoneAnnotationTemplateName, api.MAX_VSD_ACL_PRIORITY-1, false)
	if err != nil {
		return "", err
	}
	nvsdc.egressAclTemplateZoneAnnotationID = id
	return id, nil
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

	glog.Infof("Getting ACL entry by priority %d using URL: %s ", aclEntryPriority, url)

	resp, err := nvsdc.session.Get(url, nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting ACL entry with Priority %s: %d", err, aclEntryPriority)
		return nil, err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting ACL entry with priority", aclEntryPriority)
	if resp.Status() == http.StatusOK {
		glog.Infoln("Result for ACL entry obtained from VSD for priority ACL: ", result)
		glog.Infoln("Result first element for ACL entry obtained from VSD for priority ACL: ", result[0])
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

	glog.Infof("Getting ACL entry using URL: %s ", url)

	resp, err := nvsdc.session.Get(url, nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting ACL entry %v: %s", aclEntry, err)
		return nil, err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting ACL entry: ", aclEntry)
	if resp.Status() == http.StatusOK {
		glog.Infoln("Result for ACL entry obtained from VSD: ", result)
		glog.Infoln("Result first element for ACL entry obtained from VSD: ", result[0])
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
		Name:        name,
		Description: "Auto-generated domain",
		TemplateID:  domainTemplateID,
		PATEnabled:  api.PATDisabled,
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
		Description: "Auto-generated zone for project \"" + name + "\"",
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
	if vsdSubnet, err := nvsdc.GetSubnet(zoneID, subnetName); vsdSubnet != nil {
		return vsdSubnet.ID, err
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

//get interface list for a container.
func (nvsdc *NuageVsdClient) GetPodInterfaces(podName string) (*[]vspk.ContainerInterface, error) {
	//iterates over a list of containers with name matching the podName and then gets its interface elements.
	result := make([]vspk.Container, 0, 100)
	var interfaces []vspk.ContainerInterface
	e := api.RESTError{}
	nvsdc.session.Header.Add("X-Nuage-PageSize", "100")

	page := 0
	nvsdc.session.Header.Add("X-Nuage-Page", strconv.Itoa(page))
	// guarantee that the headers are cleared so that we don't change the
	// behavior of other functions
	defer nvsdc.session.Header.Del("X-Nuage-PageSize")
	defer nvsdc.session.Header.Del("X-Nuage-Page")
	for {
		nvsdc.session.Header.Add("X-Nuage-Filter", `name == "`+podName+`"`)
		resp, err := nvsdc.session.Get(nvsdc.url+"domains/"+nvsdc.domainID+"/containers",
			nil, &result, &e)
		nvsdc.session.Header.Del("X-Nuage-Filter")
		if err != nil {
			glog.Errorf("Error when getting containers matching %s: %s", podName, err)
			return nil, err
		}
		if resp.Status() == http.StatusNoContent || resp.HttpResponse().Header.Get("x-nuage-count") == "0" {
			if page == 0 {
				break
			} else {
				return &interfaces, nil
			}
		} else if resp.Status() == http.StatusOK {
			// Add all the items on this page to the list
			for _, container := range result {
				if interfaceList, err := nvsdc.GetInterfaces(container.ID); err != nil {
					glog.Errorf("Unable to get container interfaces for container %s", container.ID)
					continue
				} else {
					for _, intf := range *interfaceList {
						interfaces = append(interfaces, intf)
					}
				}
			}
			// If there's less than 100 items in the page, we must've reached
			// the last page.  Break here instead of getting the next
			// (guaranteed empty) page.
			if count, err := strconv.Atoi(resp.HttpResponse().Header.Get("x-nuage-count")); err == nil {
				if count < 100 {
					return &interfaces, nil
				}
			} else {
				// Something went wrong with parsing the x-nuage-count header
				return nil, errors.New("Invalid X-Nuage-Count: " + err.Error())
			}
			// Update headers to get the next page
			page++
			nvsdc.session.Header.Set("X-Nuage-Page", strconv.Itoa(page))
		} else {
			// Something went wrong
			return nil, VsdErrorResponse(resp, &e)
		}
	}
	return nil, errors.New("Unable to fetch pods in the domain and their interfaces")
}

func (nvsdc *NuageVsdClient) GetInterfaces(containerId string) (*[]vspk.ContainerInterface, error) {
	var interfaces []vspk.ContainerInterface
	result := make([]vspk.ContainerInterface, 0, 100)
	e := api.RESTError{}
	nvsdc.session.Header.Add("X-Nuage-PageSize", "100")
	page := 0
	nvsdc.session.Header.Add("X-Nuage-Page", strconv.Itoa(page))
	// guarantee that the headers are cleared so that we don't change the
	// behavior of other functions
	defer nvsdc.session.Header.Del("X-Nuage-PageSize")
	defer nvsdc.session.Header.Del("X-Nuage-Page")
	for {
		resp, err := nvsdc.session.Get(nvsdc.url+"containers/"+containerId+"/containerinterfaces",
			nil, &result, &e)
		if err != nil {
			glog.Errorf("Error when getting container interfaces matching %s: %s", containerId, err)
			return nil, err
		}
		if resp.Status() == http.StatusNoContent || resp.HttpResponse().Header.Get("x-nuage-count") == "0" {
			if page == 0 {
				glog.Errorf("Got an error when getting interfaces matching %s: %s", containerId, err)
				return nil, VsdErrorResponse(resp, &e)
			} else {
				return &interfaces, nil
			}
		} else if resp.Status() == http.StatusOK {
			// Add all the items on this page to the list
			for _, intf := range result {
				interfaces = append(interfaces, intf)
			}
			// If there's less than 100 items in the page, we must've reached
			// the last page.  Break here instead of getting the next
			// (guaranteed empty) page.
			if count, err := strconv.Atoi(resp.HttpResponse().Header.Get("x-nuage-count")); err == nil {
				if count < 100 {
					return &interfaces, nil
				}
			} else {
				// Something went wrong with parsing the x-nuage-count header
				return nil, errors.New("Invalid X-Nuage-Count: " + err.Error())
			}
			// Update headers to get the next page
			page++
			nvsdc.session.Header.Set("X-Nuage-Page", strconv.Itoa(page))
		} else {
			// Something went wrong
			return nil, VsdErrorResponse(resp, &e)
		}
	}
	return nil, errors.New("Unknown error when trying to fetch container interfaces")

}

//podsList is a list of pod names that need to be added to policy group with Id pgId
func (nvsdc *NuageVsdClient) AddPodsToPolicyGroup(pgId string, podsList []string) error {
	//call GetPodInterfaces() and iterate over them to get vports and add them to policy group for each pod.
	var vportsList []string
	for _, pod := range podsList {
		if interfaceList, err := nvsdc.GetPodInterfaces(pod); err == nil {
			for _, intf := range *interfaceList {
				vportsList = append(vportsList, intf.VPortID)
			}
		} else {
			glog.Errorf("Cannot get interfaces for pod %s", pod)
			continue
		}
	}
	result := make([]vspk.VPort, 0, 100)
	e := api.RESTError{}
	nvsdc.session.Header.Add("X-Nuage-PageSize", "100")
	page := 0
	nvsdc.session.Header.Add("X-Nuage-Page", strconv.Itoa(page))
	// guarantee that the headers are cleared so that we don't change the
	// behavior of other functions
	defer nvsdc.session.Header.Del("X-Nuage-PageSize")
	defer nvsdc.session.Header.Del("X-Nuage-Page")
	glog.Infof("Got the following vports %s to add to the policy group", vportsList)
	for {
		resp, err := nvsdc.session.Get(nvsdc.url+"policygroups/"+pgId+"/vports",
			nil, &result, &e)
		if err != nil {
			glog.Errorf("Error when fetching vports for pg id %s : %s", pgId, err)
			return err
		}
		if resp.Status() == http.StatusNoContent || resp.HttpResponse().Header.Get("x-nuage-count") == "0" {
			glog.Infof("No existing vports found under policy group with id: %s", pgId)
			break
		} else if resp.Status() == http.StatusOK {
			// Add all the items on this page to the list
			for _, vport := range result {
				vportsList = append(vportsList, vport.ID)
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
	if len(vportsList) != 0 {
		glog.Infof("Adding the following %d vports %s to the policygroup with id: %s", len(vportsList), vportsList, pgId)
		resp, err := nvsdc.session.Put(nvsdc.url+"policygroups/"+
			pgId+"/vports", &vportsList, nil, &e)
		if err != nil {
			glog.Errorf("Error when adding vports to policy group %s: %s", pgId, err)
			return err
		} else {
			glog.Infoln("Got a reponse status", resp.Status(),
				"when adding vports to policy group")
			switch resp.Status() {
			case http.StatusNoContent:
				glog.Infof("Added vports %s to policy group %s", vportsList, pgId)
			default:
				return VsdErrorResponse(resp, &e)
			}
		}
	}
	return nil
}

func (nvsdc *NuageVsdClient) RemovePortsFromPolicyGroup(pgId string) error {
	vportsList := make([]string, 0)
	e := api.RESTError{}
	resp, err := nvsdc.session.Put(nvsdc.url+"policygroups/"+
		pgId+"/vports", &vportsList, nil, &e)
	if err != nil {
		glog.Errorf("Error when deleting vports from policy group %s: %s", pgId, err)
		return err
	} else {
		glog.Infoln("Got a reponse status", resp.Status(),
			"when deleting vports from policy group")
		switch resp.Status() {
		case http.StatusNoContent:
			glog.Infof("Deleted vports from policy group %s", pgId)
		default:
			return VsdErrorResponse(resp, &e)
		}
	}
	return nil
}

func (nvsdc *NuageVsdClient) GetPolicyGroup(name string) (string, error) {
	result := make([]vspk.PolicyGroup, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `name == "`+name+`"`)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"domains/"+nvsdc.domainID+"/policygroups", nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting policy group ID %s", err)
		return "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting policy group ID")
	if resp.Status() == http.StatusOK {
		if result[0].Name == name {
			return result[0].ID, nil
		} else {
			return "", errors.New("Policy group not found")
		}
	} else {
		return "", VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) CreatePolicyGroup(name string, description string) (string, string, error) {
	result := make([]vspk.PolicyGroup, 1)
	payload := vspk.PolicyGroup{
		Name:        name,
		Description: description,
		Type:        "SOFTWARE",
	}
	e := api.RESTError{}
	resp, err := nvsdc.session.Post(nvsdc.url+"domains/"+nvsdc.domainID+"/policygroups", &payload, &result, &e)
	if err != nil {
		glog.Error("Error when creating policy group", err)
		return "", "", err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when creating policy group")
	switch resp.Status() {
	case http.StatusCreated:
		glog.Infoln("Created the policy group:", result[0].ID)
	case http.StatusConflict:
		glog.Infoln("Error from VSD:\n", e)
		// Subnet already exists, call Get to retrieve the ID
		if id, err := nvsdc.GetPolicyGroup(name); err != nil {
			glog.Errorf("Error when getting policy group ID: %s", err)
			return "", "", err
		} else {
			return name, id, nil
		}
	default:
		return "", "", VsdErrorResponse(resp, &e)
	}
	return result[0].Name, result[0].ID, nil
}

func (nvsdc *NuageVsdClient) DeletePolicyGroup(id string) error {
	result := make([]struct{}, 1)
	e := api.RESTError{}
	resp, err := nvsdc.session.Delete(nvsdc.url+"policygroups/"+id+"?responseChoice=1", nil, &result, &e)
	if err != nil {
		glog.Errorf("Error when deleting policy group with ID %s: %s", id, err)
		return err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when deleting policy group")
	if resp.Status() != http.StatusNoContent {
		return VsdErrorResponse(resp, &e)
	}
	return nil
}

func (nvsdc *NuageVsdClient) Run(nsChannel chan *api.NamespaceEvent, serviceChannel chan *api.ServiceEvent, podChannel chan *api.PodEvent, policyChannel chan *api.NetworkPolicyEvent, stop chan bool) {
	//we will use the kube client APIs than interfacing with the REST API
	for {
		select {
		case nsEvent := <-nsChannel:
			nvsdc.HandleNsEvent(nsEvent)
		case serviceEvent := <-serviceChannel:
			nvsdc.HandleServiceEvent(serviceEvent)
		case podEvent := <-podChannel:
			nvsdc.HandlePodEvent(podEvent)
		case policyEvent := <-policyChannel:
			nvsdc.HandleNetworkPolicyEvent(policyEvent)
		case nsUpdateReq := <-nvsdc.newSubnetQueue:
			switch nsUpdateReq.Event {
			case config.AddSubnet:
				nvsdc.CreateAdditionalSubnet(nsUpdateReq.NamespaceID)
			}
		}
	}
}

func (nvsdc *NuageVsdClient) CreateAdditionalSubnet(namespaceID string) {
	glog.Infof("Got request to create subnet in namespace %q", namespaceID)
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
		// Even though we failed to allocate a new subnet, remove the
		// NeedsNewSubnet flag so that a new subnet can be requested again.
		namespace.NeedsNewSubnet = false
		nvsdc.namespaces[namespaceID] = namespace
		return
	}
	for {
		// Release the lock while we wait for the VSD response
		nvsdc.pods.editLock.Unlock()
		subnetID, err := nvsdc.CreateSubnet(subnetName, namespace.ZoneID, subnet)
		nvsdc.pods.editLock.Lock()
		// Reacquire it before we edit the namespace again
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

func (nvsdc *NuageVsdClient) HandlePodEvent(podEvent *api.PodEvent) error {
	glog.Infoln("Received a pod event: Pod: ", podEvent)
	switch podEvent.Type {
	case api.Added:
		glog.Infof("Pod: %s was added", podEvent.Name)
	case api.Deleted:
		glog.Infof("Pod: %s was deleted", podEvent.Name)

	}
	return nil
}

func (nvsdc *NuageVsdClient) HandleNetworkPolicyEvent(policyEvent *api.NetworkPolicyEvent) error {
	glog.Infoln("Received a policy event: Policy: ", policyEvent)
	switch policyEvent.Type {
	case api.Added:
		fallthrough
	case api.Deleted:
		err := nvsdc.resourceManager.HandlePolicyEvent(policyEvent)
		glog.Infof("Policy: %s was %s %+v", policyEvent.Name, policyEvent.Type, err)
		return err
	}
	return nil
}

func (nvsdc *NuageVsdClient) HandleServiceEvent(serviceEvent *api.ServiceEvent) error {
	glog.Infoln("Received a service event: Service: ", serviceEvent)
	switch serviceEvent.Type {
	case api.Added:
		zone := serviceEvent.Namespace
		nmgID := ""
		err := errors.New("")
		exists := false
		userSpecifiedZone := false
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
				glog.Infoln("Label provided for a zone, but no network macro group identified", serviceEvent)
				userSpecifiedZone = true
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
			if !userSpecifiedZone {
				err = nvsdc.AddNetworkMacroToNMG(networkMacroID, nmgID)
				if err != nil {
					glog.Error("Error when adding network macro to network macro group:", err)
				}
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
	//handle annotations
	nvsdc.resourceManager.HandleNsEvent(nsEvent)
	//handle regular processing
	switch nsEvent.Type {
	case api.Added:
		fallthrough
	case api.Modified:
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
			if nsEvent.Name == nvsdc.privilegedProjectName {
				err = nvsdc.CreatePrivilegedZoneAcls(zoneID)
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
			if nsEvent.Name == nvsdc.privilegedProjectName {
				err = nvsdc.CreatePrivilegedZoneAcls(id)
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
			if nsEvent.Name == nvsdc.privilegedProjectName {
				err := nvsdc.DeletePrivilegedZoneAcls(zone.ZoneID)
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
			if nsEvent.Name == nvsdc.privilegedProjectName {
				err = nvsdc.DeletePrivilegedZoneAcls(id)
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

func (nvsdc *NuageVsdClient) CreatePrivilegedZoneAcls(zoneID string) error {
	nmgid, err := nvsdc.CreateNetworkMacroGroup(nvsdc.enterpriseID, nvsdc.privilegedProjectName)
	if err != nil {
		glog.Error("Error when creating the network macro group for zone", nvsdc.privilegedProjectName)
		return err
	} else {
		if nsd, exists := nvsdc.namespaces[nvsdc.privilegedProjectName]; exists {
			nsd.NetworkMacroGroupID = nmgid
			nvsdc.namespaces[nvsdc.privilegedProjectName] = nsd
		} else {
			nvsdc.namespaces[nvsdc.privilegedProjectName] = NamespaceData{
				ZoneID:              zoneID,
				Name:                nvsdc.privilegedProjectName,
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

func (nvsdc *NuageVsdClient) DeletePrivilegedZoneAcls(zoneID string) error {
	// aclEntry := api.VsdAclEntry{
	// 	Action:       "FORWARD",
	// 	DSCP:         "*",
	// 	Description:  "Allow Traffic Between All Zones and Default Zone",
	// 	EntityScope:  "ENTERPRISE",
	// 	EtherType:    "0x0800",
	// 	LocationID:   "",
	// 	LocationType: "ANY",
	// 	NetworkID:    nvsdc.namespaces[nvsdc.privilegedProjectName].NetworkMacroGroupID,
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
	if nvsdc.namespaces[nvsdc.privilegedProjectName].NetworkMacroGroupID != "" {
		err := nvsdc.DeleteNetworkMacroGroup(nvsdc.namespaces[nvsdc.privilegedProjectName].NetworkMacroGroupID)
		if err != nil {
			glog.Error("Failed to delete network macro group for default zone")
			return err
		} else {
			if nsd, exists := nvsdc.namespaces[nvsdc.privilegedProjectName]; exists {
				nsd.NetworkMacroGroupID = ""
				nvsdc.namespaces[nvsdc.privilegedProjectName] = nsd
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
		fetchedNetworkMacro, err := nvsdc.GetNetworkMacro(enterpriseID, networkMacro.Name)
		if err != nil {
			glog.Errorf("Error when getting network macro ID: %v - %v", networkMacro, err)
			return "", err
		}
		// If we got back a network macro with the same name but different info,
		// inherit the ID of the existing macro, but overwrite the contents.
		if !networkMacro.IsEqual(fetchedNetworkMacro) {
			networkMacro.ID = fetchedNetworkMacro.ID
			err := nvsdc.UpdateNetworkMacro(networkMacro)
			if err != nil {
				glog.Error("Error when updating existing network macro: ", err)
				return "", err
			}
			return networkMacro.ID, err
		}
		return fetchedNetworkMacro.ID, nil
	default:
		return "", VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) GetNetworkMacro(enterpriseID string, networkMacroName string) (*api.VsdNetworkMacro, error) {
	result := make([]api.VsdNetworkMacro, 1)
	h := nvsdc.session.Header
	h.Add("X-Nuage-Filter", `name == "`+networkMacroName+`"`)
	e := api.RESTError{}
	resp, err := nvsdc.session.Get(nvsdc.url+"enterprises/"+enterpriseID+"/enterprisenetworks", nil, &result, &e)
	h.Del("X-Nuage-Filter")
	if err != nil {
		glog.Errorf("Error when getting network macro ID for network macro: %v - %v", networkMacroName, err)
		return nil, err
	}
	glog.Infoln("Got a reponse status", resp.Status(), "when getting network macro ID")
	if resp.Status() == http.StatusOK {
		// Status code 200 is returned even if there's no results.  If
		// the filter didn't match anything (or there was nothing to
		// return), the result object will just be empty.
		if result[0].Name == networkMacroName {
			glog.Infof("Found network macro %s when looking for %q",
				result[0].ID, networkMacroName)
			return &result[0], nil
		} else if result[0].Name == "" {
			return nil, errors.New("Network Macro not found")
		} else {
			return nil, errors.New(fmt.Sprintf(
				"Found %q instead of %q", result[0].Name, networkMacroName))
		}
	} else {
		return nil, VsdErrorResponse(resp, &e)
	}
}

func (nvsdc *NuageVsdClient) GetNetworkMacroID(enterpriseID string, networkMacroName string) (string, error) {
	if networkMacro, err := nvsdc.GetNetworkMacro(enterpriseID, networkMacroName); networkMacro != nil {
		return networkMacro.ID, err
	} else {
		return "", err
	}
}

func (nvsdc *NuageVsdClient) UpdateNetworkMacro(networkMacro *api.VsdNetworkMacro) error {
	if networkMacro == nil {
		return errors.New("No network macro specified")
	}
	url := nvsdc.url + "enterprisenetworks/" + networkMacro.ID
	e := api.RESTError{}
	resp, err := nvsdc.session.Put(url, networkMacro, nil, &e)
	if err != nil || resp.Status() != http.StatusNoContent {
		VsdErrorResponse(resp, &e)
		return err
	}
	return nil
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
