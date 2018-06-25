package ovsdb

import (
	"github.com/golang/glog"
	"github.com/nuagenetworks/libvrsdk/api/entity"
	"github.com/socketplane/libovsdb"
	"reflect"
	"strings"
)

// Constants needed to describe Nuage_Port_Table
const (
	NuagePortTable           = "Nuage_Port_Table"
	NuagePortTableColumnName = "name"

	// Attributes
	NuagePortTableColumnMAC      = "mac"
	NuagePortTableColumnBridge   = "bridge"
	NuagePortTableColumnVMDomain = "vm_domain"

	// State
	NuagePortTableColumnNuageDomain  = "nuage_domain"
	NuagePortTableColumnNuageZone    = "nuage_zone"
	NuagePortTableColumnNuageNetwork = "nuage_network"
	NuagePortTableColumnIPAddress    = "ip_addr"
	NuagePortTableColumnSubnetMask   = "subnet_mask"
	NuagePortTableColumnGateway      = "gateway"
	NuagePortTableColumnVRFId        = "vrf_id"
	NuagePortTableColumnEVPNID       = "evpn_id"

	NuagePortTableColumnMetadata = "metadata"
)

// NuagePortTableRow represents a row in Nuage_Port_Table
type NuagePortTableRow struct {
	Name             string
	Mac              string
	IPAddr           string
	SubnetMask       string
	Gateway          string
	Bridge           string
	Alias            string
	NuageDomain      string
	NuageNetwork     string
	NuageZone        string
	NuageNetworkType string
	EVPNId           int
	VRFId            int
	VMDomain         entity.Domain
	Metadata         map[string]string
	Dirty            int
}

// Equals checks for equality of two Nuage_Port_Table rows
func (row *NuagePortTableRow) Equals(otherRow interface{}) bool {

	nuagePortTableRow, ok := otherRow.(NuagePortTableRow)

	if !ok {
		return false
	}

	if strings.Compare(row.Name, nuagePortTableRow.Name) != 0 {
		return false
	}

	if strings.Compare(row.Mac, nuagePortTableRow.Mac) != 0 {
		return false
	}

	if strings.Compare(row.IPAddr, nuagePortTableRow.IPAddr) != 0 {
		return false
	}

	if strings.Compare(row.SubnetMask, nuagePortTableRow.SubnetMask) != 0 {
		return false
	}

	if strings.Compare(row.Gateway, nuagePortTableRow.Gateway) != 0 {
		return false
	}

	if strings.Compare(row.Bridge, nuagePortTableRow.Bridge) != 0 {
		return false
	}

	if strings.Compare(row.Alias, nuagePortTableRow.Alias) != 0 {
		return false
	}

	if strings.Compare(row.NuageDomain, nuagePortTableRow.NuageDomain) != 0 {
		return false
	}

	if strings.Compare(row.NuageNetwork, nuagePortTableRow.NuageNetwork) != 0 {
		return false
	}

	if strings.Compare(row.NuageNetworkType, nuagePortTableRow.NuageNetworkType) != 0 {
		return false
	}

	if strings.Compare(row.NuageZone, nuagePortTableRow.NuageZone) != 0 {
		return false
	}

	if row.EVPNId != nuagePortTableRow.EVPNId {
		return false
	}

	if row.VRFId != nuagePortTableRow.VRFId {
		return false
	}

	if row.VMDomain != nuagePortTableRow.VMDomain {
		return false
	}

	if !reflect.DeepEqual(row.Metadata, nuagePortTableRow.Metadata) {
		return false
	}

	if row.Dirty != nuagePortTableRow.Dirty {
		return false
	}

	return true
}

// CreateOVSDBRow creates a new row in Nuage_Port_Table using the data provided by the user
func (row *NuagePortTableRow) CreateOVSDBRow(ovsdbRow map[string]interface{}) error {

	metadataMap, err := libovsdb.NewOvsMap(row.Metadata)
	if err != nil {
		glog.Errorf("error while creating metadataMap (%+v)", ovsdbRow)
		return err
	}

	ovsdbRow["name"] = row.Name
	ovsdbRow["mac"] = row.Mac
	ovsdbRow["ip_addr"] = row.IPAddr
	ovsdbRow["subnet_mask"] = row.SubnetMask
	ovsdbRow["gateway"] = row.Gateway
	ovsdbRow["bridge"] = row.Bridge
	ovsdbRow["alias"] = row.Alias
	ovsdbRow["nuage_domain"] = row.NuageDomain
	ovsdbRow["nuage_network"] = row.NuageNetwork
	ovsdbRow["nuage_zone"] = row.NuageZone
	ovsdbRow["nuage_network_type"] = row.NuageNetworkType
	ovsdbRow["evpn_id"] = row.EVPNId
	ovsdbRow["vrf_id"] = row.VRFId
	ovsdbRow["vm_domain"] = row.VMDomain
	ovsdbRow["metadata"] = metadataMap
	ovsdbRow["dirty"] = row.Dirty
	return nil
}
