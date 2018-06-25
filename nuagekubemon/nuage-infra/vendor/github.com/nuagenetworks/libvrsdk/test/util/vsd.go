package util

import (
	"github.com/nuagenetworks/go-bambou/bambou"
	"github.com/nuagenetworks/vspk-go/vspk"
)

// VerifyVSDPortResolution will verify if the given port is present on VSD. If yes, returns the IP
func VerifyVSDPortResolution(root *vspk.Me, vsdEnterprise string, vsdDomain string, vsdZone string, vsdPort string) (string, *bambou.Error) {

	var enterprise *vspk.Enterprise
	var domain *vspk.Domain
	var zone *vspk.Zone
	var ipAddress string
	var err *bambou.Error

	//Fetching enterprise object
	enterprise, err = FetchEnterprise(root, vsdEnterprise)
	if err != nil {
		return ipAddress, err
	}

	//Fetching domain object
	domain, err = FetchDomain(enterprise, vsdDomain)
	if err != nil {
		return ipAddress, err
	}

	//Fetching zone object
	zone, err = FetchZone(domain, vsdZone)
	if err != nil {
		return ipAddress, err
	}

	//Fetching all port objects
	var ports vspk.VPortsList
	ports, err = FetchAllPorts(zone)
	if err != nil {
		return ipAddress, err
	}

	for i := 0; i < len(ports); i++ {

		vmInterfacesFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + vsdPort + "\""}
		var vmInterfaces vspk.VMInterfacesList
		vmInterfaces, err = ports[i].VMInterfaces(vmInterfacesFetchingInfo)

		if err == nil && len(vmInterfaces) > 0 {
			vmInterface := vmInterfaces[0]
			ipAddress = vmInterface.IPAddress
		}
	}
	return ipAddress, err
}

// VerifyVSDPortDeletion will verify if the given port is removed from VSD or not
func VerifyVSDPortDeletion(root *vspk.Me, vsdEnterprise string, vsdDomain string, vsdZone string, vsdPort string) (bool, *bambou.Error) {

	var enterprise *vspk.Enterprise
	var domain *vspk.Domain
	var zone *vspk.Zone
	var portDeletionFailure bool
	var err *bambou.Error

	// Fetching enterprise object
	enterprise, err = FetchEnterprise(root, vsdEnterprise)
	if err != nil {
		return portDeletionFailure, err
	}

	// Fetching domain object
	domain, err = FetchDomain(enterprise, vsdDomain)
	if err != nil {
		return portDeletionFailure, err
	}

	// Fetching zone object
	zone, err = FetchZone(domain, vsdZone)
	if err != nil {
		return portDeletionFailure, err
	}

	// Fetching all port objects
	var ports vspk.VPortsList
	ports, err = FetchAllPorts(zone)
	if err != nil {
		return portDeletionFailure, err
	}

	for i := 0; i < len(ports); i++ {

		vmInterfacesFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + vsdPort + "\""}
		var vmInterfaces vspk.VMInterfacesList
		vmInterfaces, err = ports[i].VMInterfaces(vmInterfacesFetchingInfo)

		if err != nil || len(vmInterfaces) > 0 {
			portDeletionFailure = true
		}
	}
	return portDeletionFailure, err
}

// FetchEnterprise fetches enterprise object
func FetchEnterprise(root *vspk.Me, vsdEnterprise string) (*vspk.Enterprise, *bambou.Error) {

	enterpriseFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + vsdEnterprise + "\""}
	enterprises, enterpriseErr := root.Enterprises(enterpriseFetchingInfo)
	if enterpriseErr != nil {
		return nil, enterpriseErr
	}
	if len(enterprises) == 0 {
		return nil, bambou.NewError(101, "no matching enterprise found on VSD")
	}
	return enterprises[0], nil
}

// FetchDomain fetches domain object
func FetchDomain(enterprise *vspk.Enterprise, vsdDomain string) (*vspk.Domain, *bambou.Error) {

	domainFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + vsdDomain + "\""}
	domains, domainErr := enterprise.Domains(domainFetchingInfo)
	if domainErr != nil {
		return nil, domainErr
	}
	if len(domains) == 0 {
		return nil, bambou.NewError(102, "no matching domain found on VSD")
	}
	return domains[0], nil
}

// FetchZone fetches zone object
func FetchZone(domain *vspk.Domain, vsdZone string) (*vspk.Zone, *bambou.Error) {

	zoneFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + vsdZone + "\""}
	zones, zonesErr := domain.Zones(zoneFetchingInfo)
	if zonesErr != nil {
		return nil, zonesErr
	}
	if len(zones) == 0 {
		return nil, bambou.NewError(103, "no matching zone found on VSD")
	}
	return zones[0], nil
}

// FetchSubnete fetches subnet object
func FetchSubnet(domain *vspk.Domain, vsdSubnet string) (*vspk.Subnet, *bambou.Error) {

	subnetFetchingInfo := &bambou.FetchingInfo{Filter: "name == \"" + vsdSubnet + "\""}
	subnets, subnetsErr := domain.Subnets(subnetFetchingInfo)
	if subnetsErr != nil {
		return nil, subnetsErr
	}
	if len(subnets) == 0 {
		return nil, bambou.NewError(103, "no matching subnets found on VSD")
	}
	return subnets[0], nil
}

// FetchAllPorts fetches all port objects
func FetchAllPorts(zone *vspk.Zone) (vspk.VPortsList, *bambou.Error) {

	portsFetchingInfo := &bambou.FetchingInfo{Filter: "type == \"VM\""}
	return zone.VPorts(portsFetchingInfo)
}
