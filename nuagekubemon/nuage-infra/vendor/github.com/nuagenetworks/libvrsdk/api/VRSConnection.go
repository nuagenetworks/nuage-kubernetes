package api

import (
	"errors"
	"github.com/nuagenetworks/libvrsdk/ovsdb"
	"github.com/socketplane/libovsdb"
)

type portNameChannelMap map[string]chan *PortIPv4Info

type portNamePortInfoMap map[string]PortIPv4Info

// Registration will help to register for VRS
// port table updates
type Registration struct {
	Brport   string
	Channel  chan *PortIPv4Info
	Register bool
}

// VRSConnection represent the OVSDB connection to the VRS
type VRSConnection struct {
	ovsdbClient         *libovsdb.OvsdbClient
	vmTable             ovsdb.NuageTable
	portTable           ovsdb.NuageTable
	updatesChan         chan *libovsdb.TableUpdates
	pncTable            portNameChannelMap
	pnpTable            portNamePortInfoMap
	stopChannel         chan bool
	registrationChannel chan *Registration
}

// Disconnected will retry connecting to OVSDB
// and continue to register for OVSDB updates
func (vrsConnection VRSConnection) Disconnected(ovsClient *libovsdb.OvsdbClient) {
}

// Locked is a placeholder function for table updates
func (vrsConnection VRSConnection) Locked([]interface{}) {
}

// Stolen is a placeholder function for table updates
func (vrsConnection VRSConnection) Stolen([]interface{}) {
}

// Echo is a placeholder function for table updates
func (vrsConnection VRSConnection) Echo([]interface{}) {
}

// Update will provide updates on OVSDB table updates
func (vrsConnection VRSConnection) Update(context interface{}, tableUpdates libovsdb.TableUpdates) {
	vrsConnection.updatesChan <- &tableUpdates
}

// NewUnixSocketConnection creates a connection to the VRS Server using Unix sockets
func NewUnixSocketConnection(socketfile string) (VRSConnection, error) {
	var vrsConnection VRSConnection
	var err error

	if vrsConnection.ovsdbClient, err = libovsdb.ConnectWithUnixSocket(socketfile); err != nil {
		return vrsConnection, err
	}

	vrsConnection.vmTable.TableName = ovsdb.NuageVMTable
	vrsConnection.portTable.TableName = ovsdb.NuagePortTable
	vrsConnection.pncTable = make(portNameChannelMap)
	vrsConnection.pnpTable = make(portNamePortInfoMap)
	vrsConnection.registrationChannel = make(chan *Registration)
	vrsConnection.updatesChan = make(chan *libovsdb.TableUpdates)
	vrsConnection.stopChannel = make(chan bool)
	err = vrsConnection.monitorTable()

	return vrsConnection, err
}

func (vrsConnection *VRSConnection) monitorTable() error {
	// Setting a monitor on Nuage_Port_Table in VRS connection
	vrsConnection.ovsdbClient.Register(vrsConnection)
	tablesOfInterest := map[string]empty{"Nuage_Port_Table": {}}
	monitorRequests := make(map[string]libovsdb.MonitorRequest)
	schema, ok := vrsConnection.ovsdbClient.Schema["Open_vSwitch"]
	if !ok {
		return errors.New("Cannot read database schema")
	}

	for table, tableSchema := range schema.Tables {
		if _, interesting := tablesOfInterest[table]; interesting {
			var columns []string
			for column := range tableSchema.Columns {
				if column == "ip_addr" || column == "subnet_mask" || column == "gateway" || column == "name" || column == "mac" {
					columns = append(columns, column)
				}
			}
			monitorRequests[table] = libovsdb.MonitorRequest{
				Columns: columns,
				Select: libovsdb.MonitorSelect{
					Initial: true,
					Modify:  true,
					Delete:  true}}
		}
	}
	initialData, err := vrsConnection.ovsdbClient.Monitor("Open_vSwitch", nil, monitorRequests)
	err = vrsConnection.processUpdates(initialData)
	if err != nil {
		return errors.New("Couldn't process initial updates")
	}
	go func() {
		for {
			select {
			case registration := <-vrsConnection.registrationChannel:
				vrsConnection.handlePortRegistration(registration)
			case currentUpdate := <-vrsConnection.updatesChan:
				vrsConnection.processUpdates(currentUpdate)
			case <-vrsConnection.stopChannel:
				return
			}
		}
	}()
	return nil
}

// Disconnect closes the connection to the VRS server
func (vrsConnection VRSConnection) Disconnect() {
	vrsConnection.ovsdbClient.Disconnect()
	vrsConnection.stopChannel <- true
}
