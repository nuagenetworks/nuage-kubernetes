/*
Package libvrsdk enables interaction of userspace applications with Nuage VRS.

Concepts

The package introduces the concepts of virtual entity (which represents a VM or a container) and a port
which represent the network endpoint for the virtual entity. A userspace application first establishes a connection
to the VRS and gets a handle for it. Using the connection handle the app can call methods to notify
the VRS about the virtual entities and their associated network ports that need to be managed by Nuage VSP.

The Nuage VRS SDK is event driven. The application needs to first create ports, followed by creation of the entity
object. The application then needs to manage the lifecycle of the virtual entity and it's associated ports by
sending VRS different events. The following unit test examples show how to handle the life cycle of VM ports using the
SDK.

Examples
        const UnixSocketFile = "/var/run/openvswitch/db.sock"  
 
	// TestAddition tests that a VM and an associated port is added to VRS successfully
	func TestAddition(t *testing.T) {

		var vrsConnection VRSConnection
		var err error
                
                
		if vrsConnection, err = NewUnixSocketConnection(UnixSocketFile); err != nil {
			t.Fatal("Unable to connect to the VRS")
		}

		// Create Port Attributes
		portAttributes := port.Attributes{
			Platform: port.PlatformDocker,
			MAC:      MAC,
			Bridge:   Bridge,
		}

		// Create Port Metadata
		portMetadata := make(map[port.MetadataKey]string)
		portMetadata[port.MetadataKeyDomain] = Domain
		portMetadata[port.MetadataKeyNetwork] = Network
		portMetadata[port.MetadataKeyZone] = Zone
		portMetadata[port.MetadataKeyPg] = "pg"
		portMetadata[port.MetadataKeyStaticIP] = "192.168.100.101"

		// Add port to VRS
		err = vrsConnection.CreatePort(PortName, portAttributes, portMetadata)
		if err != nil {
			t.Fatal("Unable to add port the VRS")
		}

		// Create VM metadata
		vmMetadata := make(map[entity.MetadataKey]string)
		vmMetadata[entity.MetadataKeyUser] = User
		vmMetadata[entity.MetadataKeyUser] = Enterprise

		// Define ports associated with the VM
		ports := []string{PortName}

		// Add entity to the VRS
		entityInfo := EntityInfo{
			UUID:     VMUUID,
			Name:     "TestVM",
			Type:     entity.NuageEntityTypeLXC,
			Ports:    ports,
			Metadata: vmMetadata,
		}
		err = vrsConnection.AddEntity(entityInfo)
		if err != nil {
			t.Fatal("Unable to add entity to VRS")
		}

		// Notify VRS that VM has completed booted
		err = vrsConnection.PostEntityEvent(VMUUID, entity.EventCategoryStarted,
			entity.EventStartedBooted)

		if err != nil {
			t.Fatal("Problem sending VRS notification")
		}

		err = vrsConnection.Disconnect()
		if err != nil {
			t.Fatal("Problem disconnecting from VRS")
		}
	}

	// TestStopped tests that a VM is stopped for migration
	func TestStopped(t *testing.T) {

		var vrsConnection VRSConnection
		var err error

		if vrsConnection, err = NewConnection(VrsHost, VrsPort); err != nil {
			t.Fatal("Unable to connect to the VRS")
		}

		// Notify VRS about an existing entity being stopped for migration
		err = vrsConnection.PostEntityEvent(VMUUID, entity.EventCategoryStopped,
			entity.EventStoppedMigrated)
		if err != nil {
			t.Fatal("VRS notification failed")
		}
	}

	// TestStartedForMigration tests VM migration on the destination
	func TestStartedForMigration(t *testing.T) {
		var vrsConnection VRSConnection
		var err error

		if vrsConnection, err = NewConnection(VrsHost, VrsPort); err != nil {
			t.Fatal("Unable to connect to the VRS")
		}

		// Notify VRS about a migrated VM. Note the entity needs to be added before sending the started event
		err = vrsConnection.PostEntityEvent(VMUUID, entity.EventCategoryStarted,
			entity.EventStartedMigrated)

		if err != nil {
			t.Fatal("VRS Notification failed")
		}
	}

	// TestRemoval Tests successful removal of an entity from VRS
	func TestRemoval(t *testing.T) {
		var vrsConnection VRSConnection
		var err error

		if vrsConnection, err = NewConnection(VrsHost, VrsPort); err != nil {
			t.Fatal("Unable to connect to the VRS")
		}

		// Remove the entity
		err = vrsConnection.RemoveEntity(VMUUID)
		if err != nil {
			t.Fatal("Unable to remove the entity")
		}
	}

*/
package libvrsdk
