package ovsdb

import (
	"fmt"
	"github.com/nuagenetworks/libvrsdk/api/entity"
	"github.com/socketplane/libovsdb"
	"math/rand"
	"strings"
	"testing"
	"time"
)

const mac1 = "76:22:F6:70:4E:47"
const Domain = "TestDomain"
const NetworkType = "vxlan"
const PortName = "port1234"

const Bridge = "alubr0"
const ChangedBridge = "alubr0_changed"

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func createTestPortTableRow() NuagePortTableRow {
	metaData := make(map[string]string)
	metaData["m1"] = "v1"
	metaData["m2"] = "v2"

	nuagePortTableRow := NuagePortTableRow{
		Name:             fmt.Sprintf("vport.%d", rand.Int()),
		Mac:              mac1,
		IPAddr:           "192.168.100.100",
		SubnetMask:       "255.255.255.0",
		Gateway:          "192.168.100.1",
		Bridge:           Bridge,
		Alias:            "portAlias1234",
		NuageDomain:      Domain,
		NuageNetwork:     "TestNetwork",
		NuageZone:        "TestZone",
		NuageNetworkType: NetworkType,
		EVPNId:           60001,
		VRFId:            40001,
		VMDomain:         entity.LXC,
		Metadata:         metaData,
		Dirty:            0,
	}

	return nuagePortTableRow
}

// TestNuagePortTableInsert tests insertion of a row in the Nuage_Port_Table
func TestNuagePortTableInsert(t *testing.T) {

	ovs, err := libovsdb.Connect(OvsdbServerIP, OvsdbServerPort)
	if err != nil {
		t.Fatal("Failed to Connect. error:", err)
		panic(err)
	}
	defer ovs.Disconnect()

	nuagePortTable := new(NuageTable)
	nuagePortTable.TableName = NuagePortTable
	nuagePortTableRow := createTestPortTableRow()

	err = nuagePortTable.InsertRow(ovs, &nuagePortTableRow)
	if err != nil {
		t.Fatalf("Unable to insert row into nuage port table %v", err)
	}

	var readRowArgs ReadRowArgs
	readRowArgs.Condition = []string{"name", "==", nuagePortTableRow.Name}
	ovsdbRow, err := nuagePortTable.ReadRow(ovs, readRowArgs)
	if err != nil {
		t.Fatalf("Unable to read row from the nuage port table %v", err)
	}

	// Verify the port name for the inserted row
	if portName, ok := ovsdbRow["name"].(string); ok {
		if strings.Compare(portName, nuagePortTableRow.Name) != 0 {
			t.Fatalf("port name mismatch after insertion %s %s (%+v)", portName, PortName, ovsdbRow)
		}
	} else {
		t.Fatalf("Problem getting port name")
	}

	condition := []string{"name", "==", nuagePortTableRow.Name}
	err = nuagePortTable.DeleteRow(ovs, condition)
	if err != nil {
		t.Fatalf("Problem deleting the row")
	}
}

// TestNuagePortTableUpdate tests an update to Nuage_Port_Table row
func TestNuagePortTableUpdate(t *testing.T) {

	ovs, err := libovsdb.Connect(OvsdbServerIP, OvsdbServerPort)
	if err != nil {
		t.Fatal("Failed to Connect. error:", err)
		panic(err)
	}
	defer ovs.Disconnect()

	nuagePortTable := new(NuageTable)
	nuagePortTable.TableName = NuagePortTable
	nuagePortTableRow := createTestPortTableRow()

	err = nuagePortTable.InsertRow(ovs, &nuagePortTableRow)
	if err != nil {
		t.Fatalf("Unable to insert row into nuage port table %v", err)
	}

	nuagePortTableChangedColumns := make(map[string]interface{})
	nuagePortTableChangedColumns["bridge"] = ChangedBridge
	condition := []string{"name", "==", nuagePortTableRow.Name}
	err = nuagePortTable.UpdateRow(ovs, nuagePortTableChangedColumns, condition)
	if err != nil {
		t.Fatalf("Unable to insert row into nuage port table %v", err)
	}

	var readRowArgs ReadRowArgs
	readRowArgs.Condition = []string{"name", "==", nuagePortTableRow.Name}
	ovsdbRow, err := nuagePortTable.ReadRow(ovs, readRowArgs)
	if err != nil {
		t.Fatalf("Unable to read row from the nuage port table %v", err)
	}

	// Verify the updated bridge name
	if bridgeName, ok := ovsdbRow["bridge"].(string); ok {
		if strings.Compare(bridgeName, ChangedBridge) != 0 {
			t.Fatalf("bridge name mismatch after update %s %s (%+v)", bridgeName, PortName, ovsdbRow)
		}
	} else {
		t.Fatalf("Problem getting bridge name")
	}

	condition = []string{"name", "==", nuagePortTableRow.Name}
	err = nuagePortTable.DeleteRow(ovs, condition)
	if err != nil {
		t.Fatalf("Problem deleting the row")
	}
}

// TestNuagePortTableDelete tests deletion of row from the Nuage_Port_Table
func TestNuagePortTableDelete(t *testing.T) {

	ovs, err := libovsdb.Connect(OvsdbServerIP, OvsdbServerPort)
	if err != nil {
		t.Fatal("Failed to Connect. error:", err)
		panic(err)
	}
	defer ovs.Disconnect()

	nuagePortTable := new(NuageTable)
	nuagePortTable.TableName = NuagePortTable
	nuagePortTableRow := createTestPortTableRow()

	err = nuagePortTable.InsertRow(ovs, &nuagePortTableRow)
	if err != nil {
		t.Fatalf("Unable to insert row into nuage port table %v", err)
	}

	condition := []string{"name", "==", nuagePortTableRow.Name}
	err = nuagePortTable.DeleteRow(ovs, condition)
	if err != nil {
		t.Fatalf("Unable to delete row from the Nuage Port Table (%+v)", err)
	}
}

// TestNuagePortTableMultipleReads tests reading of multiple rows from the Nuage_Port_Table
func TestNuagePortTableMultipleReads(t *testing.T) {
	ovs, err := libovsdb.Connect(OvsdbServerIP, OvsdbServerPort)
	if err != nil {
		t.Fatal("Failed to Connect. error:", err)
		panic(err)
	}
	defer ovs.Disconnect()

	nuagePortTable := new(NuageTable)
	nuagePortTable.TableName = NuagePortTable

	portStatus := make(map[string]bool)

	for i := 0; i < 10; i++ {
		nuagePortTableRow := createTestPortTableRow()
		portStatus[nuagePortTableRow.Name] = false
		err = nuagePortTable.InsertRow(ovs, &nuagePortTableRow)
		if err != nil {
			t.Fatalf("Unable to insert row into nuage port table %v", err)
		}
	}

	readRowArgs := ReadRowArgs{
		Condition: []string{NuagePortTableColumnName, "!=", "xxxx"},
		Columns:   []string{NuagePortTableColumnName},
	}

	ovsDbPortRows, err := nuagePortTable.ReadRows(ovs, readRowArgs)
	if err != nil {
		t.Fatalf("Unable to read multiple rows from the Nuage Port Table %v", err)
	}

	for _, row := range ovsDbPortRows {
		portRow := row[NuagePortTableColumnName]
		portName, err := portRow.(string)
		if err {
			portStatus[portName] = true
		}
	}

	for port, status := range portStatus {
		if !status {
			t.Fatalf("Unable to find a expected port %s in the Nuage Port Table", port)
		}

		condition := []string{NuagePortTableColumnName, "==", port}
		err := nuagePortTable.DeleteRow(ovs, condition)
		if err != nil {
			t.Fatal("Unable to cleanup port rows")
		}
	}
}
