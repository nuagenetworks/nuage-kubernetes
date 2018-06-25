package ovsdb

import (
	"reflect"
	"strings"

	"github.com/nuagenetworks/libvrsdk/api/entity"
	"github.com/socketplane/libovsdb"
)

// These constants describe the Nuage_VM_Table
const (
	NuageVMTable                    = "Nuage_VM_Table"
	NuageVMTableColumnVMUUID        = "vm_uuid"
	NuageVMTableColumnVMName        = "vm_name"
	NuageVMTableColumnPorts         = "ports"
	NuageVMTableColumnState         = "state"
	NuageVMTableColumnReason        = "reason"
	NuageVMTableColumnEventCategory = "event"
	NuageVMTableColumnEventType     = "event_type"
	NuageVMTableColumnMetadata      = "metadata"
)

// NuageVMTableRow represents a row in the Nuage_VM_Table
type NuageVMTableRow struct {
	Type            int
	Event           int
	EventType       int
	State           int
	Reason          int
	VMUuid          string
	VMName          string
	Domain          entity.Domain
	NuageUser       string
	NuageEnterprise string
	Metadata        map[string]string
	Ports           []string
	Dirty           int
}

// Equals checks for equality of two rows in the Nuage_VM_Table
func (row *NuageVMTableRow) Equals(otherRow interface{}) bool {

	nuageVMTableRow, ok := otherRow.(NuageVMTableRow)

	if !ok {
		return false
	}

	if row.Type != nuageVMTableRow.Type {
		return false
	}

	if row.Event != nuageVMTableRow.Event {
		return false
	}

	if row.EventType != nuageVMTableRow.EventType {
		return false
	}

	if row.State != nuageVMTableRow.State {
		return false
	}

	if row.Reason != nuageVMTableRow.Reason {
		return false
	}

	if strings.Compare(row.VMUuid, nuageVMTableRow.VMUuid) != 0 {
		return false
	}

	if row.Domain != nuageVMTableRow.Domain {
		return false
	}

	if strings.Compare(row.VMName, nuageVMTableRow.VMName) != 0 {
		return false
	}

	if strings.Compare(row.NuageUser, nuageVMTableRow.NuageUser) != 0 {
		return false
	}

	if strings.Compare(row.NuageEnterprise, nuageVMTableRow.NuageEnterprise) != 0 {
		return false
	}

	if !reflect.DeepEqual(row.Metadata, nuageVMTableRow.Metadata) {
		return false
	}

	if !reflect.DeepEqual(row.Ports, nuageVMTableRow.Ports) {
		return false
	}

	if row.Dirty != nuageVMTableRow.Dirty {
		return false
	}

	return true
}

// CreateOVSDBRow creates a OVSDB row for Nuage_VM_Table
func (row *NuageVMTableRow) CreateOVSDBRow(ovsdbRow map[string]interface{}) error {

	metadataMap, err := libovsdb.NewOvsMap(row.Metadata)
	if err != nil {
		return err
	}

	portSet, err := libovsdb.NewOvsSet(row.Ports)
	if err != nil {
		return err
	}

	ovsdbRow["type"] = row.Type
	ovsdbRow["event"] = row.Event
	ovsdbRow["event_type"] = row.EventType
	ovsdbRow["state"] = row.State
	ovsdbRow["reason"] = row.Reason
	ovsdbRow["vm_uuid"] = row.VMUuid
	ovsdbRow["domain"] = row.Domain
	ovsdbRow["vm_name"] = row.VMName
	ovsdbRow["nuage_user"] = row.NuageUser
	ovsdbRow["nuage_enterprise"] = row.NuageEnterprise
	ovsdbRow["metadata"] = metadataMap
	ovsdbRow["ports"] = portSet
	ovsdbRow["dirty"] = row.Dirty

	return nil
}
