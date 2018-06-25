package ovsdb

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/socketplane/libovsdb"
)

// OvsDBName is the OVS database name
const OvsDBName = "Open_vSwitch"

// NuageTable represent a Nuage OVSDB table
type NuageTable struct {
	TableName string
}

// InsertRow enables insertion of a row into the Nuage OVSDB table
func (nuageTable *NuageTable) InsertRow(ovs *libovsdb.OvsdbClient, row NuageTableRow) error {

	glog.V(2).Infof("Trying to insert (%+v) into the Nuage Table (%s)", row, nuageTable.TableName)

	ovsdbRow := make(map[string]interface{})
	err := row.CreateOVSDBRow(ovsdbRow)
	if err != nil {
		glog.Errorf("Unable to create the OVSDB row %v", err)
		return err
	}

	insertOp := libovsdb.Operation{
		Op:       "insert",
		Table:    nuageTable.TableName,
		Row:      ovsdbRow,
		UUIDName: "gopher",
	}

	operations := []libovsdb.Operation{insertOp}
	reply, err := ovs.Transact(OvsDBName, operations...)

	glog.V(2).Infof("reply : (%+v) err : (%+v)", reply, err)

	if err != nil || len(reply) != 1 || reply[0].Error != "" {
		errStr := fmt.Errorf("Problem inserting row in the Nuage table row = "+
			" (%+v) ovsdbrow (%+v) err (%+v) reply (%+v)", row, ovsdbRow, err, reply)
		glog.Error(errStr)
		return (errStr)
	}

	glog.V(2).Info("Insertion into Nuage VM Table succeeded with UUID %s", reply[0].UUID)

	return nil
}

// ReadRowArgs enables a user to specific a condition and the columns of data to be read from a Nuage OVSDB table
type ReadRowArgs struct {
	Condition []string
	Columns   []string
}

// ReadRows enables reading of multiple rows from a Nuage OVSDB table.
func (nuageTable *NuageTable) ReadRows(ovs *libovsdb.OvsdbClient, readRowArgs ReadRowArgs) ([]map[string]interface{}, error) {

	condition := readRowArgs.Condition
	columns := readRowArgs.Columns

	glog.V(2).Infof("Reading rows from table %s with condition (%+v)", nuageTable.TableName, condition)

	var selectOp libovsdb.Operation

	if len(condition) == 3 {

		ovsdbCondition := libovsdb.NewCondition(condition[0], condition[1], condition[2])

		if columns == nil {
			selectOp = libovsdb.Operation{
				Op:    "select",
				Table: nuageTable.TableName,
				Where: []interface{}{ovsdbCondition},
			}
		} else {
			selectOp = libovsdb.Operation{
				Op:      "select",
				Table:   nuageTable.TableName,
				Where:   []interface{}{ovsdbCondition},
				Columns: columns,
			}
		}
	} else {
		if columns == nil {
			selectOp = libovsdb.Operation{
				Op:    "select",
				Table: nuageTable.TableName,
				Where: []interface{}{},
			}
		} else {
			selectOp = libovsdb.Operation{
				Op:      "select",
				Table:   nuageTable.TableName,
				Columns: columns,
				Where:   []interface{}{},
			}
		}
	}

	operations := []libovsdb.Operation{selectOp}
	reply, err := ovs.Transact(OvsDBName, operations...)

	glog.V(2).Infof("reply : (%+v) err : (%+v)", reply, err)

	if err != nil || len(reply) != 1 || reply[0].Error != "" {
		glog.Errorf("Problem reading row from the Nuage table %s %v %+v", nuageTable.TableName, err, reply)
		return nil, fmt.Errorf("Problem reading row from the Nuage table %s %v",
			nuageTable.TableName, err)
	}

	return reply[0].Rows, nil
}

// ReadRow enables reading of a single row from Nuage OVSDB table
func (nuageTable *NuageTable) ReadRow(ovs *libovsdb.OvsdbClient, readRowArgs ReadRowArgs) (map[string]interface{}, error) {

	condition := readRowArgs.Condition
	columns := readRowArgs.Columns

	glog.V(2).Infof("Reading row from table %s with condition (%+v)", nuageTable.TableName, condition)

	if len(condition) != 3 {
		glog.Errorf("Invalid condition %v", condition)
		return nil, fmt.Errorf("Invalid condition")
	}

	ovsdbCondition := libovsdb.NewCondition(condition[0], condition[1], condition[2])
	var selectOp libovsdb.Operation

	if columns == nil {
		selectOp = libovsdb.Operation{
			Op:    "select",
			Table: nuageTable.TableName,
			Where: []interface{}{ovsdbCondition},
		}
	} else {
		selectOp = libovsdb.Operation{
			Op:      "select",
			Table:   nuageTable.TableName,
			Where:   []interface{}{ovsdbCondition},
			Columns: columns,
		}
	}

	operations := []libovsdb.Operation{selectOp}
	reply, err := ovs.Transact(OvsDBName, operations...)

	glog.V(2).Infof("reply : (%+v) err : (%+v)", reply, err)

	if err != nil || len(reply) != 1 || reply[0].Error != "" {
		glog.Errorf("Problem reading row from the Nuage table %s %v", nuageTable.TableName, err)
		return nil, fmt.Errorf("Problem reading row from the Nuage table %s %v",
			nuageTable.TableName, err)
	}

	if len(reply[0].Rows) != 1 {
		glog.Errorf("Did not find a Nuage Table entry for table %s condition %v", nuageTable.TableName, condition)
		return nil, fmt.Errorf("Did not find a Nuage Table entry for table %s condition %v",
			nuageTable.TableName, condition)
	}

	ovsdbRow := reply[0].Rows[0]

	return ovsdbRow, nil
}

// DeleteRow is use to delete a row from the Nuage OVSDB table
func (nuageTable *NuageTable) DeleteRow(ovs *libovsdb.OvsdbClient, condition []string) error {

	glog.V(2).Infof("Delete from table %s with condition (%+v)", nuageTable.TableName, condition)

	if len(condition) != 3 {
		glog.Errorf("Invalid condition %v", condition)
		return fmt.Errorf("Invalid condition")
	}

	ovsdbCondition := libovsdb.NewCondition(condition[0], condition[1], condition[2])
	deleteOp := libovsdb.Operation{
		Op:    "delete",
		Table: nuageTable.TableName,
		Where: []interface{}{ovsdbCondition},
	}

	operations := []libovsdb.Operation{deleteOp}
	reply, err := ovs.Transact(OvsDBName, operations...)

	glog.V(2).Infof("reply : (%+v) err : (%+v)", reply, err)

	if err != nil || len(reply) != 1 || reply[0].Error != "" {
		errStr := fmt.Sprintf("Problem deleting row from the Nuage table %s (%+v) (%+v) (%+v)",
			nuageTable.TableName, ovsdbCondition, err, reply)
		glog.Errorf(errStr)
		return fmt.Errorf(errStr)
	}

	if reply[0].Count != 1 {
		glog.Errorf("Did not delete a Nuage Table entry for table %s condition %v", nuageTable.TableName, condition)
		return fmt.Errorf("Did not delete a Nuage Table entry for table %s condition %v",
			nuageTable.TableName, condition)
	}

	return nil
}

// UpdateRow updates the OVSDB table row
func (nuageTable *NuageTable) UpdateRow(ovs *libovsdb.OvsdbClient, ovsdbRow map[string]interface{}, condition []string) error {

	glog.V(2).Infof("Trying to update the row (%+v) the Nuage Table (%s)", ovsdbRow, nuageTable.TableName)

	ovsdbCondition := libovsdb.NewCondition(condition[0], condition[1], condition[2])
	updateOp := libovsdb.Operation{
		Op:    "update",
		Table: nuageTable.TableName,
		Row:   ovsdbRow,
		Where: []interface{}{ovsdbCondition},
	}

	operations := []libovsdb.Operation{updateOp}
	reply, err := ovs.Transact(OvsDBName, operations...)

	glog.V(2).Infof("reply : (%+v) err : (%+v)", reply, err)

	if err != nil || len(reply) != 1 || reply[0].Error != "" {
		glog.Errorf("Failed to update row in the Nuage table %s %v", nuageTable.TableName, err)
		return fmt.Errorf("Failed to update row in the Nuage table %s %v", nuageTable.TableName, err)
	}

	if reply[0].Count != 1 {
		glog.Errorf("Failed to update the Nuage Table entry for table %s condition %v", nuageTable.TableName, condition)
		return fmt.Errorf("Failed to update the Nuage Table entry for table %s condition %v",
			nuageTable.TableName, condition)
	}

	return nil
}
