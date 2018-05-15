/*
###########################################################################
#
#   Filename:           api.go
#
#   Author:             Aniket Bhat
#   Created:            July 20, 2015
#
#   Description:        NuageKubeMon event API
#
###########################################################################
#
#              Copyright (c) 2015 Nuage Networks
#
###########################################################################

*/

package api

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	networkingV1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net"
)

type PgInfo struct {
	PgName   string
	PgId     string
	Selector metav1.LabelSelector
}

type EventType string

const (
	Added    EventType = "ADDED"
	Deleted  EventType = "DELETED"
	Modified EventType = "MODIFIED"
)

const (
	EncryptionEnabled  = "ENABLED"
	EncryptionDisabled = "DISABLED"
)

const (
	UnderlaySupportEnabled   = "ENABLED"
	UnderlaySupportInherited = "INHERITED"
	UnderlaySupportDisabled  = "DISABLED"
)

const (
	IngressAclTemplateName     = "Auto-generated Ingress Policies"
	EgressAclTemplateName      = "Auto-generated Egress Policies"
	ZoneAnnotationTemplateName = "Namespace Annotations"
)

type Namespace string

type NamespaceEvent struct {
	UID         string
	Type        EventType
	Name        string
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type ServiceEvent struct {
	Type        EventType
	Name        string
	ClusterIP   string
	Namespace   string
	NuageLabels map[string]string `json:"labels,omitempty"`
}

type NetworkPolicyEvent struct {
	Type      EventType
	Name      string
	Namespace string
	Policy    networkingV1.NetworkPolicySpec
	Labels    map[string]string `json:"labels,omitempty"`
}

type PodEventResp struct {
	Data  interface{}
	Error error
}

type PodEvent struct {
	Type      EventType
	Name      string
	Namespace string
	Labels    map[string]string `json:"labels,omitempty"`
	RespChan  chan *PodEventResp
}

const (
	EtcdAddSubnet        EventType = "ETCD_ADD_SUBNET"
	EtcdDelSubnet        EventType = "ETCD_DEL_SUBNET"
	EtcdIncActiveIPCount EventType = "ETCD_INC_IP_COUNT"
	EtcdDecActiveIPCount EventType = "ETCD_DEC_IP_COUNT"
	EtcdAllocSubnetCIDR  EventType = "ETCD_ALLOC_SUBNET_CIDR"
	EtcdFreeSubnetCIDR   EventType = "ETCD_FREE_SUBNET_CIDR"
	EtcdUpdateSubnetID   EventType = "ETCD_UPDATE_SUBNET_ID"
	EtcdGetSubnetID      EventType = "ETCD_GET_SUBNET_ID"
	EtcdGetSubnetInfo    EventType = "ETCD_GET_SUBNET_INFO"
	EtcdAddZone          EventType = "ETCD_ADD_ZONE"
	EtcdDeleteZone       EventType = "ETCD_DELETE_ZONE"
	EtcdUpdateZone       EventType = "ETCD_UPDATE_ZONE"
	EtcdGetZonesSubnets  EventType = "ETCD_GET_ZONES_SUBNETS"
)

type EtcdRespObject struct {
	EtcdData interface{}
	Error    error
}

type EtcdEvent struct {
	Type               EventType
	EtcdReqObject      interface{}
	EtcdRespObjectChan chan *EtcdRespObject
}

type EtcdPodMetadata struct {
	PodName       string
	NamespaceName string
}

type EtcdSubnetMetadata struct {
	Name      string
	ID        string
	Namespace string
	CIDR      string
}

type EtcdZoneMetadata struct {
	Name string
	ID   string
}
type EtcdPodSubnet struct {
	ToUse    string
	ToCreate string
}

type GetPod func(string, string) (*PodEvent, error)
type FilterPods func(*metav1.ListOptions, string) (*[]*PodEvent, error)
type FilterNamespaces func(*metav1.ListOptions) (*[]*NamespaceEvent, error)
type FilterServices func(*metav1.ListOptions) (*[]*ServiceEvent, error)
type FilterNetworkPolicies func(*metav1.ListOptions) (*[]*NetworkPolicyEvent, error)

type ClusterClientCallBacks struct {
	FilterPods       FilterPods
	FilterNamespaces FilterNamespaces
	GetPod           GetPod
}

type RESTError struct {
	Errors []struct {
		Property     string `json:"property"`
		Descriptions []struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		} `json:"descriptions"`
	} `json:"errors"`
	InternalErrorCode int `json:"internalErrorCode"`
}

//EtcdChanRequest make a request on Etcd Channel
func EtcdChanRequest(receiver chan *EtcdEvent, event EventType, params interface{}) *EtcdRespObject {
	etcdReq := &EtcdEvent{
		Type:          event,
		EtcdReqObject: params,
	}
	etcdReq.EtcdRespObjectChan = make(chan *EtcdRespObject)
	receiver <- etcdReq
	etcdResp := <-etcdReq.EtcdRespObjectChan
	return etcdResp
}

func (restErr RESTError) String() string {
	outString := fmt.Sprintf("InternalErrorCode: %d\n",
		restErr.InternalErrorCode)
	for _, vsdErr := range restErr.Errors {
		outString += "\tProperty: " + vsdErr.Property + "\n"
		for _, description := range vsdErr.Descriptions {
			outString += "\t\tTitle: " + description.Title +
				"\n\t\tDescription: " + description.Description + "\n"
		}
	}
	return outString
}

type VsdEnterprise struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	ID          string
	ExternalID  string `json:"externalID"`
}

type VsdUser struct {
	ID         string
	UserName   string `json:"userName"`
	Password   string `json:"password"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	ExternalID string `json:"externalID"`
}

type VsdGroup struct {
	ID         string
	Role       string `json:"role"`
	ExternalID string `json:"externalID"`
}

type VsdLicense struct {
	ID         string
	License    string `json:"license"`
	LicenseId  int    `json:"licenseID"`
	ExternalID string `json:"externalID"`
}

type VsdSubnet struct {
	ID              string
	IPType          string
	Name            string `json:"name"`
	Address         string `json:"address"`
	Netmask         string `json:"netmask"`
	Description     string `json:"description"`
	UnderlayEnabled string `json:"underlayEnabled,omitempty"`
	ExternalID      string `json:"externalID"`
}

// Generic VSD object. Most json objects returned by the VSD REST API will fit
// this "interface"
type VsdObject struct {
	ID          string
	Name        string `json:"name"`
	Description string `json:"description"`
	ExternalID  string `json:"externalID"`
}

type VsdDomainTemplate struct {
	ID          string
	Name        string `json:"name"`
	Description string `json:"description"`
	ExternalID  string `json:"externalID"`
	Encryption  string `json:"encryption,omitempty"`
}

type VsdDomain struct {
	ID              string
	Name            string `json:"name"`
	Description     string `json:"description"`
	TemplateID      string `json:"templateID"`
	UnderlayEnabled string `json:"underlayEnabled,omitempty"`
	Encryption      string `json:"encryption,omitempty"`
	ExternalID      string `json:"externalID"`
}

type VsdAuthToken struct {
	APIKey       string
	APIKeyExpiry int64
	ID           string
	Email        string `json:"email"`
	EnterpriseID string `json:"enterpriseID"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Role         string `json:"role"`
	UserName     string `json:"userName"`
}

type VsdAclTemplate struct {
	ID                string
	Name              string `json:"name"`
	DefaultAllowIP    bool   `json:"defaultAllowIP"`
	DefaultAllowNonIP bool   `json:"defaultAllowNonIP"`
	Active            bool   `json:"active"`
	Priority          int    `json:"priority"`
	ExternalID        string `json:"externalID"`
}

type VsdAclEntry struct {
	DSCP                string `json:"DSCP,omitempty"`
	ID                  string
	Action              string `json:"action"`
	Description         string `json:"description"`
	EntityScope         string `json:"entityScope"`
	EtherType           string `json:"etherType"`
	LocationID          string `json:"locationID"`
	LocationType        string `json:"locationType"`
	NetworkID           string `json:"networkID"`
	NetworkType         string `json:"networkType"`
	PolicyState         string `json:"policyState"`
	Priority            int    `json:"priority"`
	Protocol            string `json:"protocol"`
	Stateful            bool   `json:"stateful"`
	ExternalID          string `json:"externalID"`
	StatsLoggingEnabled bool   `json:"statsLoggingEnabled"`
}

const MAX_VSD_ACL_PRIORITY = 1000000000 //the maximum priority allowed in VSD is 1 billion.

type VsdNetworkMacro struct {
	ID         string
	Name       string `json:"name"`
	IPType     string `json:"IPType"`
	Address    string `json:"address"`
	Netmask    string `json:"netmask"`
	ExternalID string `json:"externalID"`
}

func (lhs *VsdNetworkMacro) IsEqual(rhs *VsdNetworkMacro) bool {
	if lhs.Name != rhs.Name {
		return false
	}
	if lhs.IPType != rhs.IPType {
		return false
	}
	if lhs.Address != rhs.Address {
		return false
	}
	if lhs.Netmask != rhs.Netmask {
		return false
	}
	return true
}

func (acl *VsdAclEntry) TryNextAclPriority() {
	if acl.Priority == MAX_VSD_ACL_PRIORITY {
		acl.Priority = acl.Priority - 1
	} else {
		acl.Priority = acl.Priority + 1
	}
}

func (lhs *VsdAclEntry) IsEqual(rhs *VsdAclEntry) bool {
	if lhs.DSCP != "" && lhs.DSCP != rhs.DSCP {
		glog.Info("DSCP for 2 ACLs dont match")
		glog.Infof("LHS: %s, RHS: %s", lhs.DSCP, rhs.DSCP)
		return false
	}
	if lhs.Action != "" && lhs.Action != rhs.Action {
		glog.Info("Action for 2 ACLs dont match")
		glog.Infof("LHS: %s, RHS: %s", lhs.Action, rhs.Action)
		return false
	}
	if lhs.EntityScope != "" && lhs.EntityScope != rhs.EntityScope {
		glog.Info("Entity Scope for 2 ACLs dont match")
		glog.Infof("LHS: %s, RHS: %s", lhs.EntityScope, rhs.EntityScope)
		return false
	}
	if lhs.EtherType != "" && lhs.EtherType != rhs.EtherType {
		glog.Info("Ether Type for 2 ACLs dont match")
		glog.Infof("LHS: %s, RHS: %s", lhs.EtherType, rhs.EtherType)
		return false
	}
	if lhs.LocationID != "" && lhs.LocationID != rhs.LocationID {
		glog.Info("Location ID for 2 ACLs dont match")
		glog.Infof("LHS: %s, RHS: %s", lhs.LocationID, rhs.LocationID)
		return false
	}
	if lhs.LocationType != "" && lhs.LocationType != rhs.LocationType {
		glog.Info("Location Type for 2 ACLs dont match")
		glog.Infof("LHS: %s, RHS: %s", lhs.LocationType, rhs.LocationType)
		return false
	}
	if lhs.NetworkID != "" && lhs.NetworkID != rhs.NetworkID {
		glog.Info("Network ID for 2 ACLs dont match")
		glog.Infof("LHS: %s, RHS: %s", lhs.NetworkID, rhs.NetworkID)
		return false
	}
	if lhs.NetworkType != "" && lhs.NetworkType != rhs.NetworkType {
		glog.Info("Network Type for 2 ACLs dont match")
		glog.Infof("LHS: %s, RHS: %s", lhs.NetworkType, rhs.NetworkType)
		return false
	}
	if lhs.PolicyState != "" && lhs.PolicyState != rhs.PolicyState {
		glog.Info("Policy State for 2 ACLs dont match")
		glog.Infof("LHS: %s, RHS: %s", lhs.PolicyState, rhs.PolicyState)
		return false
	}
	if lhs.Protocol != "" && lhs.Protocol != rhs.Protocol {
		glog.Info("Protocol for 2 ACLs dont match")
		glog.Infof("LHS: %s, RHS: %s", lhs.Protocol, rhs.Protocol)
		return false
	}
	return true
}

func (lhs *VsdAclEntry) BuildFilter() string {
	filter := ""
	if lhs.DSCP != "" {
		dscpClause := `DSCP == "` + lhs.DSCP + `"`
		filter = dscpClause
	}
	if lhs.Action != "" {
		actionClause := `action == "` + lhs.Action + `"`
		if filter != "" {
			filter = filter + ` and ` + actionClause
		} else {
			filter = actionClause
		}
	}
	// The entity scope is not filterable currently in VSD

	// if lhs.EntityScope != "" {
	// 	entityScopeClause := `entityScope == "` + lhs.EntityScope + `"`
	// 	if filter != "" {
	// 		filter = filter + ` and ` + entityScopeClause
	// 	} else {
	// 		filter = entityScopeClause
	// 	}
	// }
	if lhs.EtherType != "" {
		etherTypeClause := `etherType == "` + lhs.EtherType + `"`
		if filter != "" {
			filter = filter + ` and ` + etherTypeClause
		} else {
			filter = etherTypeClause
		}
	}
	if lhs.LocationID != "" {
		locationIDClause := `locationID == "` + lhs.LocationID + `"`
		if filter != "" {
			filter = filter + ` and ` + locationIDClause
		} else {
			filter = locationIDClause
		}
	}
	if lhs.LocationType != "" {
		locationTypeClause := `locationType == "` + lhs.LocationType + `"`
		if filter != "" {
			filter = filter + ` and ` + locationTypeClause
		} else {
			filter = locationTypeClause
		}
	}
	if lhs.NetworkID != "" {
		networkIDClause := `networkID == "` + lhs.NetworkID + `"`
		if filter != "" {
			filter = filter + ` and ` + networkIDClause
		} else {
			filter = networkIDClause
		}
	}
	if lhs.NetworkType != "" {
		networkTypeClause := `networkType == "` + lhs.NetworkType + `"`
		if filter != "" {
			filter = filter + ` and ` + networkTypeClause
		} else {
			filter = networkTypeClause
		}
	}
	if lhs.PolicyState != "" {
		policyStateClause := `policyState == "` + lhs.PolicyState + `"`
		if filter != "" {
			filter = filter + ` and ` + policyStateClause
		} else {
			filter = policyStateClause
		}
	}
	if lhs.Protocol != "" {
		protocolClause := `protocol == "` + lhs.Protocol + `"`
		if filter != "" {
			filter = filter + ` and ` + protocolClause
		} else {
			filter = protocolClause
		}
	}
	return filter
}

func (lhs *VsdAclEntry) String() string {
	return fmt.Sprintf(`\nACL Entry: ID: %v, Description: %v,\n`+
		`Priority: %v, Action: %v,\n`+
		`DSCP: %v, EntityScope: %v, EtherType: %v, Protocol: %v\n`+
		`LocationID: %v, LocationType: %v\n`+
		`NetworkID: %v, NetworkType: %v, PolicyState: %v, Stateful %v`,
		lhs.ID, lhs.Description, lhs.Priority, lhs.Action, lhs.DSCP,
		lhs.EntityScope, lhs.EtherType, lhs.Protocol, lhs.LocationID,
		lhs.LocationType, lhs.NetworkID, lhs.NetworkType, lhs.PolicyState,
		lhs.Stateful)
}

func (svc *ServiceEvent) String() string {
	return fmt.Sprintf(`\nService: Name: %v, Namespace %v,\n`+
		`ClusterIP: %v, Labels: %v,\n`+
		`EventType: %v`, svc.Name, svc.Namespace, svc.ClusterIP,
		svc.NuageLabels, svc.Type)
}

func (lhs *VsdNetworkMacro) String() string {
	return fmt.Sprintf(`\nNetwork Macro: Name: %v, ID: %v,\n`+
		`IPType: %v, Address: %v,\n`+
		`Netmask: %v`, lhs.Name, lhs.ID, lhs.IPType, lhs.Address, lhs.Netmask)
}

type IPv4Address [4]byte

func (addr IPv4Address) String() string {
	return fmt.Sprintf("%v.%v.%v.%v", addr[0], addr[1], addr[2], addr[3])
}

func (addr IPv4Address) ToUint() uint {
	return uint(addr[0])<<24 +
		uint(addr[1])<<16 +
		uint(addr[2])<<8 +
		uint(addr[3])
}

type IPv4Subnet struct {
	Address  IPv4Address
	CIDRMask int //e.g. 24, not 255.255.255.0
}

func IPv4SubnetFromString(input string) (*IPv4Subnet, error) {
	output := &IPv4Subnet{}
	n, err := fmt.Sscanf(input, "%d.%d.%d.%d/%d", &output.Address[0],
		&output.Address[1], &output.Address[2], &output.Address[3],
		&output.CIDRMask)
	if err != nil {
		return nil, err
	}
	if n != 5 {
		return nil, errors.New(fmt.Sprintf(
			"Invalid syntax in input string %q", input))
	}
	return output, nil
}

/* Parse address and netmask separately, and generate an IPv4Subnet
 * datastructure.
 * @Param address - IPv4 Address as a string
 * @Param netmask - Netmask in IPv4 address format (e.g. 255.255.0.0) as a
 *                  string
 */
func IPv4SubnetFromAddrNetmask(address, netmaskStr string) (*IPv4Subnet, error) {
	// Parse address
	output := &IPv4Subnet{}
	n, err := fmt.Sscanf(address, "%d.%d.%d.%d", &output.Address[0],
		&output.Address[1], &output.Address[2], &output.Address[3])
	if err != nil {
		return nil, err
	}
	if n != 4 {
		return nil, errors.New(fmt.Sprintf(
			"Invalid syntax in address string %q", address))
	}
	// Parse netmask
	var nmA, nmB, nmC, nmD byte
	n, err = fmt.Sscanf(netmaskStr, "%d.%d.%d.%d", &nmA, &nmB, &nmC, &nmD)
	if err != nil {
		return nil, err
	}
	if n != 4 {
		return nil, errors.New(fmt.Sprintf(
			"Invalid syntax in netmask string %q", netmaskStr))
	}
	var length int
	output.CIDRMask, length = net.IPv4Mask(nmA, nmB, nmC, nmD).Size()
	if length != 32 {
		return nil, errors.New(fmt.Sprintf("Failed to parse netmask %q",
			netmaskStr))
	}
	return output, nil
}

func (subnet IPv4Subnet) String() string {
	return fmt.Sprintf("%v/%v", subnet.Address, subnet.CIDRMask)
}

func (subnet IPv4Subnet) Netmask() IPv4Address {
	// returns the traditional IPv4 netmask instead of the CIDR
	// e.g. .../24 would return 255.255.255.0
	if subnet.CIDRMask >= 32 {
		return IPv4Address{255, 255, 255, 255}
	}
	fullmask := uint((1 << 32) - (1 << uint(32-subnet.CIDRMask)))
	return IPv4Address{
		byte((fullmask / uint(1<<24)) % 256),
		byte((fullmask / uint(1<<16)) % 256),
		byte((fullmask / 256) % 256),
		byte(fullmask % 256),
	}
}

func (subnet *IPv4Subnet) Split() (*IPv4Subnet, *IPv4Subnet, error) {
	if subnet.CIDRMask >= 32 {
		return nil, nil, errors.New("Cannot split /32 address space")
	}
	loSubnet, hiSubnet := &IPv4Subnet{}, &IPv4Subnet{}
	for i, mask := 0, subnet.CIDRMask; i < 4; i++ {
		switch {
		case mask >= 8:
			loSubnet.Address[i] = subnet.Address[i]
			hiSubnet.Address[i] = subnet.Address[i]
			mask -= 8
		case mask > 0:
			bitmask := byte(uint(256-(1<<uint(8-mask))) % 256)
			loSubnet.Address[i] = subnet.Address[i] & bitmask
			hiSubnet.Address[i] = subnet.Address[i] & bitmask
			mask = 0
		}
	}
	loSubnet.CIDRMask = subnet.CIDRMask + 1
	hiSubnet.CIDRMask = subnet.CIDRMask + 1
	index := subnet.CIDRMask / 8
	offset := uint(subnet.CIDRMask % 8)
	bit := byte(128) >> offset
	loSubnet.Address[index] &= ^bit
	hiSubnet.Address[index] |= bit
	return loSubnet, hiSubnet, nil
}

// Compare `a` to `b`.  If `a > b`, the result will be positive.  If `a < b`,
// the result will be negative.  If `a == b`, the result will be 0.
func (a *IPv4Subnet) Compare(b *IPv4Subnet) int {
	// For sorting purposes, a subnet with a smaller mask (larger size) will
	// always be greater than a subnet with a larger mask.
	if n := b.CIDRMask - a.CIDRMask; n != 0 {
		return n
	}
	aNetmask := a.Netmask().ToUint()
	bNetmask := b.Netmask().ToUint()
	aAddr := a.Address.ToUint()
	bAddr := b.Address.ToUint()
	// Compare only significant bits by &-ing the addresses with the netmask
	return int((aAddr & aNetmask) - (bAddr & bNetmask))
}

func (a *IPv4Subnet) Contains(b *IPv4Subnet) bool {
	if a.CIDRMask > b.CIDRMask {
		// if a is smaller than b (e.g. a is /24, but b is /16), a can't contain b
		return false
	}
	return a.Compare(&IPv4Subnet{Address: b.Address, CIDRMask: a.CIDRMask}) == 0
}

func CanMerge(a, b *IPv4Subnet) bool {
	// We can't merge the /0 address space.
	if a.CIDRMask <= 0 || b.CIDRMask <= 0 {
		return false
	}
	// An address can't be merged with itself.
	if a.Compare(b) == 0 {
		return false
	}
	// Addresses with different netmasks can't be merged.
	if a.CIDRMask != b.CIDRMask {
		return false
	}
	aCopy := &IPv4Subnet{a.Address, a.CIDRMask - 1}
	bCopy := &IPv4Subnet{b.Address, b.CIDRMask - 1}
	return aCopy.Compare(bCopy) == 0
}

func Merge(a, b *IPv4Subnet) (*IPv4Subnet, error) {
	if !CanMerge(a, b) {
		return nil, errors.New(fmt.Sprintf("Can't merge subnets %s and %s!", a, b))
	}
	newSubnet := &IPv4Subnet{a.Address, a.CIDRMask - 1}
	index := newSubnet.CIDRMask / 8
	mask := byte(uint(1<<8 - 1<<uint(8-(newSubnet.CIDRMask%8))))
	newSubnet.Address[index] &= mask
	return newSubnet, nil
}
