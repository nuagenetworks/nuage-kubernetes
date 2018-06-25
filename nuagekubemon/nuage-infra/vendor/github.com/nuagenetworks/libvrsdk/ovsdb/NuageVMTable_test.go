package ovsdb

import (
	"fmt"
	"github.com/docker/distribution/uuid"
	"github.com/nuagenetworks/libvrsdk/api/entity"
	"github.com/socketplane/libovsdb"
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"testing"
)

const OvsdbServerIP = "127.0.0.1"
const OvsdbServerPort = 6640

func createTestVMTableRow() NuageVMTableRow {

	metaData := make(map[string]string)
	metaData["m1"] = "v1"
	metaData["m2"] = "v2"

	port1Uuid := uuid.Generate().String()
	port2Uuid := uuid.Generate().String()

	ports := []string{port1Uuid, port2Uuid}

	nuageVMTableRow := NuageVMTableRow{
		Type:            11,
		State:           12,
		Reason:          0,
		Event:           0,
		EventType:       1,
		VMName:          fmt.Sprintf("vm-%d", rand.Int()),
		VMUuid:          uuid.Generate().String(),
		Domain:          entity.LXC,
		NuageUser:       "NuageUser",
		NuageEnterprise: "NuageEnterprise",
		Metadata:        metaData,
		Ports:           ports,
		Dirty:           0,
	}

	return nuageVMTableRow
}

func TestOVSConnection(t *testing.T) {
	ovs, err := libovsdb.Connect(OvsdbServerIP, OvsdbServerPort)
	if err != nil {
		t.Fatal("Failed to Connect. error:", err)
		panic(err)
	}
	defer ovs.Disconnect()

	dbs, ok := ovs.ListDbs()

	if ok != nil {
		t.Fatal("Unable to get the OVS database")
	}

	t.Logf("(%+v)", dbs)
}

func TestNuageVMTableInsert(t *testing.T) {

	ovs, err := libovsdb.Connect(OvsdbServerIP, OvsdbServerPort)
	if err != nil {
		t.Fatal("Failed to Connect. error:", err)
		panic(err)
	}
	defer ovs.Disconnect()

	nuageVMTable := new(NuageTable)
	nuageVMTable.TableName = NuageVMTable
	nuageVMTableRow := createTestVMTableRow()

	err = nuageVMTable.InsertRow(ovs, &nuageVMTableRow)
	if err != nil {
		t.Fatalf("Unable to insert row into nuage vm table %v", err)
	}

	var readRowArgs ReadRowArgs
	readRowArgs.Condition = []string{"vm_uuid", "==", nuageVMTableRow.VMUuid}
	ovsdbRow, err := nuageVMTable.ReadRow(ovs, readRowArgs)
	if err != nil {
		t.Fatalf("Unable to read row from the nuage vm table %v", err)
	}

	// Verify the vm name for the inserted row
	if vmName, ok := ovsdbRow["vm_name"].(string); ok {
		if strings.Compare(vmName, nuageVMTableRow.VMName) != 0 {
			t.Fatalf("vm name mismatch after insertion %s %s (%+v)", vmName,
				nuageVMTableRow.VMName, ovsdbRow)
		}
	} else {
		t.Fatalf("Problem getting vm name")
	}

	condition := []string{"vm_uuid", "==", nuageVMTableRow.VMUuid}
	err = nuageVMTable.DeleteRow(ovs, condition)
	if err != nil {
		t.Fatalf("Problem deleting the row")
	}
}

func TestNuageVMTableUpdate(t *testing.T) {

	ovs, err := libovsdb.Connect(OvsdbServerIP, OvsdbServerPort)
	if err != nil {
		t.Fatal("Failed to Connect. error:", err)
		panic(err)
	}
	defer ovs.Disconnect()

	nuageVMTable := new(NuageTable)
	nuageVMTable.TableName = NuageVMTable
	nuageVMTableRow := createTestVMTableRow()
	vmUUID := nuageVMTableRow.VMUuid

	err = nuageVMTable.InsertRow(ovs, &nuageVMTableRow)
	if err != nil {
		t.Fatalf("Unable to insert row into nuage vm table %v", err)
	}

	nuageVMTableChangedColumns := make(map[string]interface{})
	nuageVMTableChangedColumns["vm_name"] = "ChangedVMName"
	condition := []string{"vm_uuid", "==", vmUUID}
	err = nuageVMTable.UpdateRow(ovs, nuageVMTableChangedColumns, condition)
	if err != nil {
		t.Fatalf("Unable to insert row into nuage vm table %v", err)
	}

	var readRowArgs ReadRowArgs
	readRowArgs.Condition = []string{"vm_uuid", "==", nuageVMTableRow.VMUuid}
	ovsdbRow, err := nuageVMTable.ReadRow(ovs, readRowArgs)
	if err != nil {
		t.Fatalf("Unable to read row from the nuage vm table %v", err)
	}

	// Verify the updated vm name
	if vmName, ok := ovsdbRow["vm_name"].(string); ok {
		if strings.Compare(vmName, "ChangedVMName") != 0 {
			t.Fatalf("vm name mismatch after update %s %s (%+v)", vmName, "ChangedVMName", ovsdbRow)
		}
	} else {
		t.Fatalf("Problem getting vm name")
	}

	for i := 0; i < 2; i++ {
		// Change the port names
		nuageVMTableChangedColumns = make(map[string]interface{})
		ports := []string{}
		port1 := fmt.Sprintf("veth1.%04d", rand.Int())
		ports = append(ports, port1)
		if i == 1 {
			port2 := fmt.Sprintf("veth1.%04d", rand.Int())
			ports = append(ports, port2)

			port3 := fmt.Sprintf("veth1.%04d", rand.Int())
			ports = append(ports, port3)
		}

		// Update the new ports
		nuageVMTableChangedColumns["ports"], err = libovsdb.NewOvsSet(ports)
		if err != nil {
			t.Fatalf("Unable to create a port set")
		}

		err = nuageVMTable.UpdateRow(ovs, nuageVMTableChangedColumns, condition)
		if err != nil {
			t.Fatalf("Unable to update row into nuage vm table %v", err)
		}

		// Read and verify the ports
		readRowArgs.Condition = []string{"vm_uuid", "==", vmUUID}
		ovsdbRow, err = nuageVMTable.ReadRow(ovs, readRowArgs)
		if err != nil {
			t.Fatalf("Unable to read row from the nuage vm table %v", err)
		}

		if readPorts, ok := ovsdbRow["ports"]; ok {
			finalPorts, err := UnMarshallOVSStringSet(readPorts)
			if err != nil {
				t.Fatalf("Problem unmarshalling port data %+v", ovsdbRow["ports"])
			}

			sort.Strings(finalPorts)
			sort.Strings(ports)

			if !reflect.DeepEqual(finalPorts, ports) {
				t.Fatalf("Did not find expected ports %+v after an update %+v", finalPorts, ports)
			}
		} else {
			t.Fatalf("Unable to read ports ")
		}
	}
}

func TestNuageVMTableDelete(t *testing.T) {

	ovs, err := libovsdb.Connect(OvsdbServerIP, OvsdbServerPort)
	if err != nil {
		t.Fatal("Failed to Connect. error:", err)
		panic(err)
	}
	defer ovs.Disconnect()

	nuageVMTable := new(NuageTable)
	nuageVMTable.TableName = NuageVMTable
	nuageVMTableRow := createTestVMTableRow()

	err = nuageVMTable.InsertRow(ovs, &nuageVMTableRow)
	if err != nil {
		t.Fatalf("Unable to insert row into nuage vm table %v", err)
	}

	condition := []string{"vm_uuid", "==", nuageVMTableRow.VMUuid}
	err = nuageVMTable.DeleteRow(ovs, condition)
	if err != nil {
		t.Fatalf("Unable to delete row from the Nuage VM Table (%+v)", err)
	}
}
