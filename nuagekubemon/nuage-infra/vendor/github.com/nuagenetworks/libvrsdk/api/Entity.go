package api

import (
	"fmt"
	"strings"

	"github.com/nuagenetworks/libvrsdk/api/entity"
	"github.com/nuagenetworks/libvrsdk/ovsdb"
	"github.com/socketplane/libovsdb"
)

// EntityInfo represents the information about an entity that needs to provided by the user to VRS
type EntityInfo struct {
	UUID     string
	Name     string
	Type     entity.Type
	Domain   entity.Domain
	Ports    []string
	Metadata map[entity.MetadataKey]string
	Events   *entity.EntityEvents
}

// CreateEntity adds an entity to the Nuage VRS
func (vrsConnection *VRSConnection) CreateEntity(info EntityInfo) error {

	if len(info.UUID) == 0 {
		return fmt.Errorf("Uuid absent")
	}

	if len(info.Name) == 0 {
		return fmt.Errorf("Name absent")
	}

	// The Nuage_VM_Table has separate columns for enterprise and user.
	// Hence make a copy of the metadata and delete these keys.
	var metadata map[string]string
	if info.Metadata != nil {
		metadata = make(map[string]string)
		for k, v := range info.Metadata {
			metadata[string(k)] = v
		}
	}
	//delete(metadata, string(entity.MetadataKeyEnterprise))
	delete(metadata, string(entity.MetadataKeyUser))

	nuageVMTableRow := ovsdb.NuageVMTableRow{
		Type:            int(info.Type),
		VMName:          info.Name,
		VMUuid:          info.UUID,
		Domain:          info.Domain,
		NuageUser:       info.Metadata[entity.MetadataKeyUser],
		NuageEnterprise: info.Metadata[entity.MetadataKeyEnterprise],
		Metadata:        metadata,
		Ports:           info.Ports,
		Event:           int(entity.EventCategoryDefined),
		EventType:       int(entity.EventDefinedAdded),
		State:           int(entity.Running),
		Reason:          int(entity.RunningUnknown),
	}

	if info.Events != nil {
		nuageVMTableRow.Event = int(info.Events.EntityEventCategory)
		nuageVMTableRow.EventType = int(info.Events.EntityEventType)
		nuageVMTableRow.State = int(info.Events.EntityState)
		nuageVMTableRow.Reason = int(info.Events.EntityReason)
	}

	if err := vrsConnection.vmTable.InsertRow(vrsConnection.ovsdbClient, &nuageVMTableRow); err != nil {
		return fmt.Errorf("Problem adding entity info to VRS %v", err)
	}

	return nil
}

// DestroyEntity removes an entity from the Nuage VRS
func (vrsConnection *VRSConnection) DestroyEntityByVMName(vm_name string) error {

	condition := []string{ovsdb.NuageVMTableColumnVMName, "==", vm_name}
	if err := vrsConnection.vmTable.DeleteRow(vrsConnection.ovsdbClient, condition); err != nil {
		return fmt.Errorf("Unable to delete the entity from VRS %v", err)
	}

	return nil
}

// DestroyEntity removes an entity from the Nuage VRS
func (vrsConnection *VRSConnection) DestroyEntity(uuid string) error {

	condition := []string{ovsdb.NuageVMTableColumnVMUUID, "==", uuid}
	if err := vrsConnection.vmTable.DeleteRow(vrsConnection.ovsdbClient, condition); err != nil {
		return fmt.Errorf("Unable to delete the entity from VRS %v", err)
	}

	return nil
}

// AddEntityPort adds a port to the Entity
func (vrsConnection *VRSConnection) AddEntityPort(uuid string, portName string) error {
	var ports []string
	var err error
	if ports, err = vrsConnection.GetEntityPorts(uuid); err != nil {
		return fmt.Errorf("Unable to get existing ports %s %s", uuid, err)
	}

	ports = append(ports, portName)
	row := make(map[string]interface{})
	row[ovsdb.NuageVMTableColumnPorts], err = libovsdb.NewOvsSet(ports)
	if err != nil {
		return err
	}

	condition := []string{ovsdb.NuageVMTableColumnVMUUID, "==", uuid}

	if err = vrsConnection.vmTable.UpdateRow(vrsConnection.ovsdbClient, row, condition); err != nil {
		return fmt.Errorf("Unable to add port %s %s %s", uuid, portName, err)
	}

	return nil
}

// RemoveEntityPort removes port from the Entity
func (vrsConnection *VRSConnection) RemoveEntityPort(uuid string, portName string) error {
	var ports []string
	var err error
	if ports, err = vrsConnection.GetEntityPorts(uuid); err != nil {
		return fmt.Errorf("Unable to get existing ports %s %s", uuid, err)
	}

	portIndex := -1
	for i, port := range ports {
		if strings.Compare(port, portName) == 0 {
			portIndex = i
			break
		}
	}

	if portIndex == -1 {
		return fmt.Errorf("%s port %s not found", uuid, portName)
	}

	ports = append(ports[:portIndex], ports[(portIndex+1):]...)

	row := make(map[string]interface{})
	row[ovsdb.NuageVMTableColumnPorts], err = libovsdb.NewOvsSet(ports)
	if err != nil {
		return err
	}

	condition := []string{ovsdb.NuageVMTableColumnVMUUID, "==", uuid}

	if err = vrsConnection.vmTable.UpdateRow(vrsConnection.ovsdbClient, row, condition); err != nil {
		return fmt.Errorf("Unable to remove port %s %s %s", uuid, portName, err)
	}

	return nil
}

// GetEntityPorts retrives the list of all of the attached ports
func (vrsConnection *VRSConnection) GetEntityPorts(uuid string) ([]string, error) {

	readRowArgs := ovsdb.ReadRowArgs{
		Columns:   []string{ovsdb.NuageVMTableColumnPorts},
		Condition: []string{ovsdb.NuageVMTableColumnVMUUID, "==", uuid},
	}

	row, err := vrsConnection.vmTable.ReadRow(vrsConnection.ovsdbClient, readRowArgs)
	if err != nil {
		return []string{}, fmt.Errorf("Unable to get port information for the VM")
	}

	return ovsdb.UnMarshallOVSStringSet(row[ovsdb.NuageVMTableColumnPorts])
}

// SetEntityState sets the entity state
func (vrsConnection *VRSConnection) SetEntityState(uuid string, state entity.State, subState entity.SubState) error {

	row := make(map[string]interface{})
	row[ovsdb.NuageVMTableColumnState] = int(state)
	row[ovsdb.NuageVMTableColumnReason] = int(subState)

	condition := []string{ovsdb.NuageVMTableColumnVMUUID, "==", uuid}

	if err := vrsConnection.vmTable.UpdateRow(vrsConnection.ovsdbClient, row, condition); err != nil {
		return fmt.Errorf("Unable to update the state %s %v %v %v", uuid, state, subState, err)
	}

	return nil
}

// PostEntityEvent posts a new event to entity
func (vrsConnection *VRSConnection) PostEntityEvent(uuid string, evtCategory entity.EventCategory, evt entity.Event) error {

	if !entity.ValidateEvent(evtCategory, evt) {
		return fmt.Errorf("Invalid event %v for event category %v", evt, evtCategory)
	}

	row := make(map[string]interface{})
	row[ovsdb.NuageVMTableColumnEventCategory] = int(evtCategory)
	row[ovsdb.NuageVMTableColumnEventType] = int(evt)

	condition := []string{ovsdb.NuageVMTableColumnVMUUID, "==", uuid}

	if err := vrsConnection.vmTable.UpdateRow(vrsConnection.ovsdbClient, row, condition); err != nil {
		return fmt.Errorf("Unable to send the state %s %v %v %v", uuid, evtCategory, evt, err)
	}

	return nil
}

// SetEntityMetadata applies Nuage specific metadata to the Entity
func (vrsConnection *VRSConnection) SetEntityMetadata(uuid string, metadata map[entity.MetadataKey]string) error {
	row := make(map[string]interface{})
	row[ovsdb.NuageVMTableColumnMetadata] = metadata

	condition := []string{ovsdb.NuageVMTableColumnVMUUID, "==", uuid}

	if err := vrsConnection.vmTable.UpdateRow(vrsConnection.ovsdbClient, row, condition); err != nil {
		return fmt.Errorf("Unable to update the metadata %s %v %v", uuid, metadata, err)
	}

	return nil
}

// GetAllEntities retrives a slice of all the UUIDs of the entities associated with the VRS
func (vrsConnection *VRSConnection) GetAllEntities() ([]string, error) {
	readRowArgs := ovsdb.ReadRowArgs{
		Condition: []string{ovsdb.NuageVMTableColumnVMUUID, "!=", "xxxx"},
		Columns:   []string{ovsdb.NuageVMTableColumnVMUUID},
	}

	var uuidRows []map[string]interface{}
	var err error
	if uuidRows, err = vrsConnection.vmTable.ReadRows(vrsConnection.ovsdbClient, readRowArgs); err != nil {
		return []string{}, fmt.Errorf("Unable to obtain the entity uuids %v", err)
	}

	var uuids []string
	for _, uuid := range uuidRows {
		uuids = append(uuids, uuid[ovsdb.NuageVMTableColumnVMUUID].(string))
	}

	return uuids, nil
}

// CheckEntityExists verifies if a specified entity exists in VRS
func (vrsConnection *VRSConnection) CheckEntityExists(id string) (bool, error) {
	readRowArgs := ovsdb.ReadRowArgs{
		Condition: []string{ovsdb.NuageVMTableColumnVMUUID, "==", id},
		Columns:   []string{ovsdb.NuageVMTableColumnVMUUID},
	}

	var idRows []map[string]interface{}
	var err error
	if idRows, err = vrsConnection.vmTable.ReadRows(vrsConnection.ovsdbClient, readRowArgs); err != nil {
		return false, fmt.Errorf("OVSDB read error %v", err)
	}

	var ids []string
	for _, row := range idRows {
		ids = append(ids, row[ovsdb.NuageVMTableColumnVMUUID].(string))
	}

	if len(ids) == 1 && id == ids[0] {
		return true, err
	}

	return false, err
}
