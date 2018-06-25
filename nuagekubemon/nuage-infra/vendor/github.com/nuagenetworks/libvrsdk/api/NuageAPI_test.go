package api

import (
	"bytes"
	"fmt"
	"github.com/docker/distribution/uuid"
	"github.com/nuagenetworks/go-bambou/bambou"
	"github.com/nuagenetworks/libvrsdk/api/entity"
	"github.com/nuagenetworks/libvrsdk/api/port"
	"github.com/nuagenetworks/libvrsdk/test/util"
	"github.com/nuagenetworks/vspk-go/vspk"
	"golang.org/x/crypto/ssh"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	VRS1            = "127.0.0.1"
	VRS2            = "127.0.0.1"
	VRSPort         = 6640
	User            = "sdkuser"
	Enterprise      = "vrsdk"
	Domain          = "vrsdk-domain"
	Zone            = "vrsdk-zone"
	Network1        = "vrsdk-subnet-1"
	Network2        = "vrsdk-subnet-2"
	ReconfNWPrefix  = "192.168.178"
	Bridge          = "alubr0"
	VSDIP           = "10.10.10.10"
	VSDPort         = "8443"
	VSDURL          = "https://" + VSDIP + ":" + VSDPort
	VSDUsername     = "*******"
	VSDPassword     = "*******"
	VSDOrganization = "***"
	VRSUser         = "root"
	VRSPassword     = "******"
	UnixSocketFile  = "/var/run/openvswitch/db.sock"
)

var VSDConnection *bambou.Session
var Root *vspk.Me

func init() {
	VSDConnection, Root = vspk.NewSession(VSDUsername, VSDPassword, VSDOrganization, VSDURL)
	if err1, err2 := VSDConnection.SetInsecureSkipVerify(true), VSDConnection.Start(); err1 != nil || err2 != nil {
		panic("Unable to connect to the VSD")
	}
}

func execCMDOnRemoteHost(cmd string, host string) error {

	config := &ssh.ClientConfig{
		User: VRSUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(VRSPassword),
		},
	}

	conn, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return fmt.Errorf("Error while establishing connection to remote VRS %v", err)
	}

	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("Error while establishing remote session to VRS %v", err)
	}

	defer func() {
		err := session.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("Error while executing config commands on remote VRS %v", err)
	}

	return nil
}

// getPortInfo will register for updates from VRS for entity
// port information
func getPortInfo(vrsConnection VRSConnection, port string) (*PortIPv4Info, error) {

	portInfoUpdateChan := make(chan *PortIPv4Info)
	vrsConnection.RegisterForPortUpdates(port, portInfoUpdateChan)
	ticker := time.NewTicker(time.Duration(10) * time.Second)
	portInfo := &PortIPv4Info{}
	select {
	case portInfo = <-portInfoUpdateChan:
		return portInfo, nil
	case <-ticker.C:
		return portInfo, fmt.Errorf("Unable to obtain port information from VRS")
	}
}

// cleanup will be a template to perform post
// API test execution entity port cleanup on VRS
func cleanup(vrsConnection VRSConnection, vmInfo map[string]string) error {

	var err error
	err = vrsConnection.DeregisterForPortUpdates(vmInfo["entityport"])
	if err != nil {
		return fmt.Errorf("Unable to deregister for port updates")
	}

	err = vrsConnection.DestroyEntity(vmInfo["vmuuid"])
	if err != nil {
		return fmt.Errorf("Unable to remove the entity from OVSDB table %v", err)
	}

	// Performing cleanup of port/entity on VRS
	err = vrsConnection.DestroyPort(vmInfo["entityport"])
	if err != nil {
		return fmt.Errorf("Unable to delete port from OVSDB table %v", err)
	}

	err = vrsConnection.RemovePortFromAlubr0(vmInfo["entityport"])
	if err != nil {
		return fmt.Errorf("Unable to delete veth port as part of cleanup from alubr0 %v", err)
	}

	// Cleaning up veth paired ports from VRS
	err = util.DeleteVETHPair(vmInfo["entityport"], vmInfo["brport"])
	if err != nil {
		return fmt.Errorf("Unable to delete veth pairs as a part of cleanup on VRS %v", err)
	}

	return nil
}

// createVM will be a template to create dummy VM entries per
// API test execution
func createVM(vrsConnection VRSConnection, vmInfo map[string]string, domain entity.Domain,
	eventCategory entity.EventCategory, eventType entity.Event) error {

	var err error
	// Create Port Attributes
	portAttributes := port.Attributes{
		Platform: domain,
		MAC:      vmInfo["mac"],
		Bridge:   Bridge,
	}

	// Create Port Metadata
	portMetadata := make(map[port.MetadataKey]string)
	portMetadata[port.MetadataKeyDomain] = Domain
	portMetadata[port.MetadataKeyNetwork] = Network1
	portMetadata[port.MetadataKeyZone] = Zone
	portMetadata[port.MetadataKeyNetworkType] = "ipv4"

	// Associate one veth port to entity
	err = vrsConnection.CreatePort(vmInfo["entityport"], portAttributes, portMetadata)
	if err != nil {
		return fmt.Errorf("Unable to create entity port %v", err)
	}

	// Create VM metadata
	vmMetadata := make(map[entity.MetadataKey]string)
	vmMetadata[entity.MetadataKeyUser] = User
	vmMetadata[entity.MetadataKeyEnterprise] = Enterprise

	// Define ports associated with the VM
	ports := []string{vmInfo["entityport"]}
	entityInfo := EntityInfo{}
	entityInfo.UUID = vmInfo["vmuuid"]
	entityInfo.Name = vmInfo["name"]
	entityInfo.Domain = domain
	entityInfo.Ports = ports
	entityInfo.Metadata = vmMetadata
	if domain != entity.Docker {
		entityInfo.Type = entity.VM
	} else {
		entityInfo.Type = entity.Container
		events := &entity.EntityEvents{}
		events.EntityEventCategory = entity.EventCategoryStarted
		events.EntityEventType = entity.EventStartedBooted
		events.EntityState = entity.Running
		events.EntityReason = entity.RunningBooted
		entityInfo.Events = events
	}
	err = vrsConnection.CreateEntity(entityInfo)
	if err != nil {
		return fmt.Errorf("Unable to add entity to VRS %v", err)
	}

	if domain != entity.Docker {
		// Notify VRS that VM has completed booted
		err = vrsConnection.PostEntityEvent(vmInfo["vmuuid"], eventCategory, eventType)
		if err != nil {
			return fmt.Errorf("Problem sending VRS notification %v", err)
		}
	}
	return nil
}

// splitActivationCreateContainer will be a template to create dummy VM entries per
// API test execution
func splitActivationCreateContainer(vrsConnection VRSConnection, vmInfo map[string]string, domain entity.Domain,
	eventCategory entity.EventCategory, eventType entity.Event) error {

	var err error
	// Create Port Attributes
	portAttributes := port.Attributes{
		Platform: domain,
		MAC:      vmInfo["mac"],
		Bridge:   Bridge,
	}

	// Create Port Metadata
	portMetadata := make(map[port.MetadataKey]string)
	portMetadata[port.MetadataKeyDomain] = ""
	portMetadata[port.MetadataKeyNetwork] = ""
	portMetadata[port.MetadataKeyZone] = ""
	portMetadata[port.MetadataKeyNetworkType] = ""

	// Associate one veth port to entity
	err = vrsConnection.CreatePort(vmInfo["entityport"], portAttributes, portMetadata)
	if err != nil {
		return fmt.Errorf("Unable to create entity port %v", err)
	}

	// Create VM metadata
	vmMetadata := make(map[entity.MetadataKey]string)
	vmMetadata[entity.MetadataKeyUser] = ""
	vmMetadata[entity.MetadataKeyEnterprise] = ""
	vmMetadata[entity.MetadataKeyExtension] = "true"

	// Define ports associated with the VM
	ports := []string{vmInfo["entityport"]}

	events := &entity.EntityEvents{}
	events.EntityEventCategory = entity.EventCategoryStarted
	events.EntityEventType = entity.EventStartedBooted
	events.EntityState = entity.Running
	events.EntityReason = entity.RunningBooted
	// Add entity to the VRS
	entityInfo := EntityInfo{
		UUID:     vmInfo["vmuuid"],
		Name:     vmInfo["name"],
		Type:     entity.Container,
		Domain:   domain,
		Ports:    ports,
		Metadata: vmMetadata,
		Events:   events,
	}
	err = vrsConnection.CreateEntity(entityInfo)
	if err != nil {
		return fmt.Errorf("Unable to add entity to VRS %v", err)
	}

	return nil
}

// TestGetAllVMsVports queries and gets all existing VMs as well as vports
func TestGetAllVMsVports(t *testing.T) {

	var vrsConnection VRSConnection
	var err error

	err = util.EnableOVSDBRPCSocket(VRSPort)
	if err != nil {
		t.Fatal("Unable to add an interface to the ovsdb-server on VRS to make it accept RPCs via TCP socket")
	}

	if vrsConnection, err = NewUnixSocketConnection(UnixSocketFile); err != nil {
		t.Fatal("Unable to connect to the VRS")
	}

	vmInfo := make(map[string]string)
	vmInfo["name"] = fmt.Sprintf("Test-VM-%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	vmInfo["mac"] = util.GenerateMAC()
	vmInfo["vmuuid"] = uuid.Generate().String()
	vmInfo["entityport"] = fmt.Sprintf("veth.%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	vmInfo["brport"] = fmt.Sprintf("vethbr.%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	portList := []string{vmInfo["entityport"], vmInfo["brport"]}
	err = util.CreateVETHPair(portList)
	if err != nil {
		t.Fatal("Unable to create veth pairs on VRS")
	}

	// Add the paired veth port to alubr0 on VRS
	var entityInfo EntityInfo
	entityInfo.Name = vmInfo["name"]
	entityInfo.UUID = vmInfo["vmuuid"]
	err = vrsConnection.AddPortToAlubr0(vmInfo["entityport"], entityInfo)
	if err != nil {
		t.Errorf("Error inserting row in alubr0: %v", err)
		t.Fatal("Unable to add veth port to alubr0")
	}

	err = createVM(vrsConnection, vmInfo, entity.KVM, entity.EventCategoryStarted, entity.EventStartedBooted)
	if err != nil {
		t.Fatal("Unable to create a test VM")
	}

	t.Logf("Waiting for 15 seconds before querying for all existing VMs")
	time.Sleep(time.Duration(15) * time.Second)

	vms, err := vrsConnection.GetAllEntities()
	if err != nil {
		t.Fatal("Unable to get existing VMs")
	}

	t.Logf("Successfully obtained all existing VMs")
	fmt.Println(vms)

	vports, err := vrsConnection.GetAllPorts()
	if err != nil {
		t.Fatal("Unable to get existing vports")
	}

	t.Logf("Successfully obtained all existing vports")
	fmt.Println(vports)

	// Performing cleanup of port/entity on VRS
	err = cleanup(vrsConnection, vmInfo)
	if err != nil {
		t.Fatal("Unable to delete port from OVSDB table")
	}
}

// TestContainerCreateDelete tests that a VM and an associated port gets resolved
// in VRS-VM as well as on the VSD and gets removed from VRS and VSD when deleted
func TestContainerCreateDelete(t *testing.T) {

	var vrsConnection VRSConnection
	var err error

	err = util.EnableOVSDBRPCSocket(VRSPort)
	if err != nil {
		t.Fatal("Unable to add an interface to the ovsdb-server on VRS to make it accept RPCs via TCP socket")
	}

	if vrsConnection, err = NewUnixSocketConnection(UnixSocketFile); err != nil {
		t.Fatal("Unable to connect to the VRS")
	}

	vmInfo := make(map[string]string)
	vmInfo["name"] = fmt.Sprintf("Test-Container-%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	vmInfo["mac"] = util.GenerateMAC()
	containerUUID := strings.Replace(uuid.Generate().String(), "-", "", -1)
	containerUUID = containerUUID + strings.Replace(uuid.Generate().String(), "-", "", -1)
	vmInfo["vmuuid"] = containerUUID
	vmInfo["entityport"] = fmt.Sprintf("veth.%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	vmInfo["brport"] = fmt.Sprintf("vethbr.%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	portList := []string{vmInfo["entityport"], vmInfo["brport"]}
	err = util.CreateVETHPair(portList)
	if err != nil {
		t.Fatal("Unable to create veth pairs on VRS")
	}

	var entityInfo EntityInfo
	entityInfo.Name = vmInfo["name"]
	entityInfo.UUID = vmInfo["vmuuid"]
	// Add the paired veth port to alubr0 on VRS
	err = vrsConnection.AddPortToAlubr0(vmInfo["entityport"], entityInfo)
	if err != nil {
		t.Errorf("Error inserting row in alubr0: %v", err)
		t.Fatal("Unable to add veth port to alubr0")
	}

	err = createVM(vrsConnection, vmInfo, entity.Docker, entity.EventCategoryStarted, entity.EventStartedBooted)
	if err != nil {
		t.Fatal("Unable to create a test VM")
	}
	// Obtaining entity port information from VRS using update mechanism
	portInfo, err := getPortInfo(vrsConnection, vmInfo["entityport"])
	if err != nil {
		t.Fatal("Unable to obtain an IP address for the VM from VRS")
	}

	if portInfo.IPAddr == "0.0.0.0" || portInfo.IPAddr == "" {
		t.Fatalf("Unable to resolve VM %s ", vmInfo["name"])
	}

	t.Logf("VM %s got resolved with an IP address %s on VRS", vmInfo["name"], portInfo.IPAddr)

	// Verifying if entity exists in OVSDB
	entityExists, err := vrsConnection.CheckEntityExists(vmInfo["vmuuid"])
	if err != nil {
		t.Fatal("Unable to verify if entity exists in OVSDB")
	}

	if entityExists {
		t.Logf("VM %s with UUID %s exists in OVSDB", vmInfo["name"], vmInfo["vmuuid"])
	} else {
		t.Fatalf("VM %s with UUID %s does not exist in OVSDB", vmInfo["name"], vmInfo["vmuuid"])
	}

	// Performing cleanup of port/entity on VRS
	err = cleanup(vrsConnection, vmInfo)
	if err != nil {
		t.Fatal("Unable to delete port from OVSDB table")
	}

	t.Logf("Waiting for 300 seconds before verifying port gets removed from VRS")

	portState, _ := vrsConnection.GetPortState(vmInfo["entityport"])
	if _, ok := portState[port.StateKeyIPAddress]; ok {
		t.Fatal("Entry for deleted VM Port still present in OVSDB table")
	}

	t.Logf("VM %s got removed from VRS successfully", vmInfo["name"])

	vrsConnection.Disconnect()
}

// TestVMCreateDelete tests that a VM and an associated port gets resolved
// in VRS-VM as well as on the VSD and gets removed from VRS and VSD when deleted
func TestVMCreateDelete(t *testing.T) {

	var vrsConnection VRSConnection
	var err error

	err = util.EnableOVSDBRPCSocket(VRSPort)
	if err != nil {
		t.Fatal("Unable to add an interface to the ovsdb-server on VRS to make it accept RPCs via TCP socket")
	}

	if vrsConnection, err = NewUnixSocketConnection(UnixSocketFile); err != nil {
		t.Fatal("Unable to connect to the VRS")
	}

	vmInfo := make(map[string]string)
	vmInfo["name"] = fmt.Sprintf("Test-VM-%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	vmInfo["mac"] = util.GenerateMAC()
	vmInfo["vmuuid"] = uuid.Generate().String()
	vmInfo["entityport"] = fmt.Sprintf("veth.%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	vmInfo["brport"] = fmt.Sprintf("vethbr.%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	portList := []string{vmInfo["entityport"], vmInfo["brport"]}
	err = util.CreateVETHPair(portList)
	if err != nil {
		t.Fatal("Unable to create veth pairs on VRS")
	}

	var entityInfo EntityInfo
	entityInfo.Name = vmInfo["name"]
	entityInfo.UUID = vmInfo["vmuuid"]
	// Add the paired veth port to alubr0 on VRS
	err = vrsConnection.AddPortToAlubr0(vmInfo["entityport"], entityInfo)
	if err != nil {
		t.Errorf("Error inserting row in alubr0: %v", err)
		t.Fatal("Unable to add veth port to alubr0")
	}

	err = createVM(vrsConnection, vmInfo, entity.KVM, entity.EventCategoryStarted, entity.EventStartedBooted)
	if err != nil {
		t.Fatal("Unable to create a test VM")
	}
	// Obtaining entity port information from VRS using update mechanism
	portInfo, err := getPortInfo(vrsConnection, vmInfo["entityport"])
	if err != nil {
		t.Fatal("Unable to obtain an IP address for the VM from VRS")
	}

	// Verifying port got an IP on VSD
	portIPOnVSD, vsdError := util.VerifyVSDPortResolution(Root, Enterprise, Domain, Zone, vmInfo["entityport"])
	if vsdError != nil || portIPOnVSD == "" || portIPOnVSD == "0.0.0.0" {
		t.Fatal("IP resolution for port " + vmInfo["entityport"] + " failed on VSD.")
	} else {
		t.Logf("VM %s port %s got resolved with an IP address %s on VSD", vmInfo["name"], vmInfo["entityport"], portIPOnVSD)
	}

	if portInfo.IPAddr == "0.0.0.0" || portInfo.IPAddr == "" {
		t.Fatalf("Unable to resolve VM %s ", vmInfo["name"])
	}

	portIPOnVRS := portInfo.IPAddr
	t.Logf("VM %s got resolved with an IP address %s on VRS", vmInfo["name"], portIPOnVRS)

	// Comparing port's IP address on VRS and VSD
	if portIPOnVSD != portIPOnVRS {
		t.Fatal("Port IPs on VSD and VRS do not match.")
	} else {
		t.Logf("Port IPs on VSD and VRS match.")
	}

	// Verifying if entity exists in OVSDB
	entityExists, err := vrsConnection.CheckEntityExists(vmInfo["vmuuid"])
	if err != nil {
		t.Fatal("Unable to verify if entity exists in OVSDB")
	}

	if entityExists {
		t.Logf("VM %s with UUID %s exists in OVSDB", vmInfo["name"], vmInfo["vmuuid"])
	} else {
		t.Fatalf("VM %s with UUID %s does not exist in OVSDB", vmInfo["name"], vmInfo["vmuuid"])
	}

	// Performing cleanup of port/entity on VRS
	err = cleanup(vrsConnection, vmInfo)
	if err != nil {
		t.Fatal("Unable to delete port from OVSDB table")
	}

	t.Logf("Waiting for 300 seconds before verifying port gets removed from VRS")
	time.Sleep(time.Duration(300) * time.Second)

	// Verifying port deletion on VSD
	if portDeletionFailure, vsdErr := util.VerifyVSDPortDeletion(Root, Enterprise, Domain, Zone, vmInfo["entityport"]); vsdErr != nil || portDeletionFailure {
		t.Fatal("Deleted VM port still present on VSD")
	}

	portState, _ := vrsConnection.GetPortState(vmInfo["entityport"])

	if _, ok := portState[port.StateKeyIPAddress]; ok {
		t.Fatal("Entry for deleted VM Port still present in OVSDB table")
	}

	t.Logf("VM %s got removed from VRS successfully", vmInfo["name"])

	vrsConnection.Disconnect()
}

// TestVMMigrate tests that a VM and its ports get resolved on VRS as well as on VSD during a migration
func TestVMMigrate(t *testing.T) {

	if os.Getenv("RUN_MIGRATION_TEST") == "" {
		t.Skip("Skipping execution of migration test; $RUN_MIGRATION_TEST not set")
	}

	var sourceVrsConnection VRSConnection
	var destinationVrsConnection VRSConnection
	var err error

	// Building VM Info
	vmInfo := make(map[string]string)
	vmInfo["name"] = fmt.Sprintf("Test-VM-%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	vmInfo["mac"] = util.GenerateMAC()
	vmInfo["vmuuid"] = uuid.Generate().String()
	vmInfo["entityport"] = fmt.Sprintf("veth.%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	vmInfo["brport"] = fmt.Sprintf("vethbr.%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	portList := []string{vmInfo["entityport"], vmInfo["brport"]}

	err = util.EnableOVSDBRPCSocket(VRSPort)
	if err != nil {
		t.Fatal("Unable to add an interface to the ovsdb-server on VRS to make it accept RPCs via TCP socket")
	}

	if sourceVrsConnection, err = NewUnixSocketConnection(UnixSocketFile); err != nil {
		t.Fatalf("Unable to connect to the VRS : %s", VRS1)
	}

	err = util.CreateVETHPair(portList)
	if err != nil {
		t.Fatal("Unable to create veth pairs on VRS")
	}

	// Creating a test VM on source VRS
	err = createVM(sourceVrsConnection, vmInfo, entity.ESXI, entity.EventCategoryStarted, entity.EventStartedBooted)
	if err != nil {
		t.Fatal("Unable to create a test VM")
	}

	// Waiting for 15 seconds before verifying port got resolved with an IP address
	time.Sleep(time.Duration(15) * time.Second)

	// Verifying the VM gets resolved with an IP address on VRS-VM
	portState := make(map[port.StateKey]interface{})
	portState, err = sourceVrsConnection.GetPortState(vmInfo["entityport"])

	if err != nil {
		t.Fatal("Unable to query port state on VRS")
	}

	if portState[port.StateKeyIPAddress] == "0.0.0.0" || portState[port.StateKeyIPAddress] == "" {
		t.Fatalf("Unable to resolve VM %s port %s", vmInfo["name"], vmInfo["entityport"])
	}

	portIP := portState[port.StateKeyIPAddress]
	fmt.Printf("VM %s port %s got resolved with IP address %s\n", vmInfo["name"], vmInfo["entityport"], portIP)

	err = util.EnableOVSDBRPCSocket(VRSPort)
	if err != nil {
		t.Fatal("Unable to add an interface to the ovsdb-server on destination VRS to make it accept RPCs via TCP socket")
	}

	if destinationVrsConnection, err = NewUnixSocketConnection(UnixSocketFile); err != nil {
		t.Fatalf("Unable to connect to the VRS : %s", VRS2)
	}

	// Create veth ports on destination VRS-VM
	cmdstr := fmt.Sprintf("ip link %s %s type veth peer name %s", "add", portList[0], portList[1])
	err = execCMDOnRemoteHost(cmdstr, VRS2)
	if err != nil {
		t.Fatal("Error while creating veth ports on destination VRS-VM")
	}

	for index := range portList {
		cmdstr = fmt.Sprintf("ip link set dev %s up", portList[index])
		err = execCMDOnRemoteHost(cmdstr, VRS2)

		if err != nil {
			t.Fatal("Error while bringing up veth ports on destination VRS-VM")
		}
	}

	cmdstr = fmt.Sprintf("/usr/bin/ovs-vsctl --no-wait --if-exists del-port alubr0 %s -- %s-port alubr0 %s -- set interface %s 'external-id={vm-uuid=%s,vm-name=%s}'", vmInfo["entityport"], "add", vmInfo["entityport"], vmInfo["entityport"], vmInfo["vmuuid"], vmInfo["name"])
	err = execCMDOnRemoteHost(cmdstr, VRS2)
	if err != nil {
		t.Fatal("Error while adding veth ports to alubr0 on destination VRS-VM")
	}

	// Creating the same test VM on destination VRS
	err = createVM(destinationVrsConnection, vmInfo, entity.ESXI, entity.EventCategoryStarted, entity.EventStartedBooted)
	if err != nil {
		t.Fatal("Unable to boot a migrating VM")
	}

	// Remove the VM from source VRS
	err = sourceVrsConnection.PostEntityEvent(vmInfo["vmuuid"], entity.EventCategoryStopped,
		entity.EventStoppedMigrated)
	if err != nil {
		t.Fatal("Error sending event notification to VRS")
	}

	// Removing VM port from alubr0
	err = cleanup(sourceVrsConnection, vmInfo)
	if err != nil {
		t.Fatal("Unable to remove port from OVSDB table")
	}

	// Waiting for 15 seconds before verifying migrated VM port got removed from VRS
	time.Sleep(time.Duration(15) * time.Second)

	portState, _ = sourceVrsConnection.GetPortState(vmInfo["entityport"])

	if _, ok := portState[port.StateKeyIPAddress]; ok {
		t.Fatal("Entry for migrated VM Port still present in OVSDB table")
	}

	// Verifying if the VM got resolved on destination VRS after migration
	// Waiting for 15 seconds before verifying port got resolved with an IP address
	time.Sleep(time.Duration(15) * time.Second)

	portState, err = destinationVrsConnection.GetPortState(vmInfo["entityport"])
	if err != nil {
		t.Fatal("Unable to query port state on VRS")
	}
	if portState[port.StateKeyIPAddress] == "0.0.0.0" || portState[port.StateKeyIPAddress] == "" {
		t.Fatalf("Unable to resolve migrating VM %s ", vmInfo["name"])
	}

	portIPPostVMMigration := portState[port.StateKeyIPAddress]

	if portIP != portIPPostVMMigration {
		t.Fatal("Migrated VM booted with a different IP")
	}

	// Performing cleanup of port/entity on VRS
	err = cleanup(destinationVrsConnection, vmInfo)
	if err != nil {
		t.Fatal("Unable to delete port from OVSDB table for migrated VM")
	}

	// Cleaning up veth ports created on destination VRS-VM
	cmdstr = fmt.Sprintf("/usr/bin/ovs-vsctl --no-wait %s-port alubr0 %s", "delete", vmInfo["entityport"])
	err = execCMDOnRemoteHost(cmdstr, VRS2)
	if err != nil {
		t.Fatal("Error while removing veth port from alubr0 on destination VRS-VM")
	}

	cmdstr = fmt.Sprintf("ip link %s %s type veth peer name %s", "delete", vmInfo["entityport"], vmInfo["brport"])
	err = execCMDOnRemoteHost(cmdstr, VRS2)
	if err != nil {
		t.Fatal("Error while cleaning up veth ports on destination VRS-VM")
	}
}

// TestVMHotNICAdd tests hot NIC addition on a VM
func TestVMHotNICAdd(t *testing.T) {

	var vrsConnection VRSConnection
	var err error

	err = util.EnableOVSDBRPCSocket(VRSPort)
	if err != nil {
		t.Fatal("Unable to add an interface to the ovsdb-server on VRS to make it accept RPCs via TCP socket")
	}

	if vrsConnection, err = NewUnixSocketConnection(UnixSocketFile); err != nil {
		t.Fatal("Unable to connect to the VRS")
	}

	vmInfo := make(map[string]string)
	vmInfo["name"] = fmt.Sprintf("Test-VM-%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	vmInfo["mac"] = util.GenerateMAC()
	vmInfo["vmuuid"] = uuid.Generate().String()
	vmInfo["entityport"] = fmt.Sprintf("veth.%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	vmInfo["brport"] = fmt.Sprintf("vethbr.%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	portList := []string{vmInfo["entityport"], vmInfo["brport"]}
	err = util.CreateVETHPair(portList)
	if err != nil {
		t.Fatal("Unable to create veth pairs on VRS")
	}

	var entityInfo EntityInfo
	entityInfo.Name = vmInfo["name"]
	entityInfo.UUID = vmInfo["vmuuid"]
	// Add the paired veth port to alubr0 on VRS
	err = vrsConnection.AddPortToAlubr0(vmInfo["entityport"], entityInfo)
	if err != nil {
		t.Errorf("Error inserting row in alubr0: %v", err)
		t.Fatal("Unable to add veth port to alubr0")
	}

	err = createVM(vrsConnection, vmInfo, entity.KVM, entity.EventCategoryStarted, entity.EventStartedBooted)
	if err != nil {
		t.Fatal("Unable to create a test VM")
	}

	// Obtaining entity port information from VRS using update mechanism
	portInfo, err := getPortInfo(vrsConnection, vmInfo["entityport"])
	if err != nil {
		t.Fatal("Unable to obtain an IP address for the VM from VRS")
	}

	portIPOnVSD, vsdError := util.VerifyVSDPortResolution(Root, Enterprise, Domain, Zone, vmInfo["entityport"])
	if vsdError != nil || portIPOnVSD == "" || portIPOnVSD == "0.0.0.0" {
		t.Fatal("IP resolution for port " + vmInfo["entityport"] + " failed on VSD.")
	} else {
		t.Logf("VM %s port %s got resolved with an IP address %s on VSD", vmInfo["name"], vmInfo["entityport"], portIPOnVSD)
	}

	if portInfo.IPAddr == "0.0.0.0" || portInfo.IPAddr == "" {
		t.Fatalf("Unable to resolve VM %s ", vmInfo["name"])
	}

	portIPOnVRS := portInfo.IPAddr
	t.Logf("VM %s got resolved with an IP address %s successfully", vmInfo["name"], portIPOnVRS)

	// Comparing port's IP address on VRS and VSD
	if portIPOnVSD != portIPOnVRS {
		t.Fatal("Port IPs on VSD and VRS do not match.")
	}

	// Adding a NIC to an existing entity to test HOT NIC addition
	hotNICEntityPort := fmt.Sprintf("veth.%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	hotNICBridgePort := fmt.Sprintf("vethbr.%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	portListHotNIC := []string{hotNICEntityPort, hotNICBridgePort}
	err = util.CreateVETHPair(portListHotNIC)
	if err != nil {
		t.Fatal("Unable to create veth pairs on VRS")
	}

	portAttributes := port.Attributes{
		Platform: entity.KVM,
		MAC:      util.GenerateMAC(),
		Bridge:   Bridge,
	}

	// Create HOT NIC Port Metadata
	portMetadata := make(map[port.MetadataKey]string)
	portMetadata[port.MetadataKeyDomain] = Domain
	portMetadata[port.MetadataKeyNetwork] = Network1
	portMetadata[port.MetadataKeyZone] = Zone
	portMetadata[port.MetadataKeyNetworkType] = "ipv4"

	// Associate one veth port to entity
	err = vrsConnection.CreatePort(hotNICEntityPort, portAttributes, portMetadata)
	if err != nil {
		t.Fatal("Unable to create entity new NIC port")
	}

	// Add the paired veth port to alubr0 on VRS
	err = vrsConnection.AddPortToAlubr0(hotNICEntityPort, entityInfo)
	if err != nil {
		t.Errorf("Error inserting row in alubr0: %v", err)
		t.Fatal("Unable to add veth port to alubr0")
	}

	err = vrsConnection.AddEntityPort(vmInfo["vmuuid"], hotNICEntityPort)
	if err != nil {
		t.Fatal("Unable to add Port to entity")
	}

	// Obtaining entity port information from VRS using update mechanism
	portInfo, err = getPortInfo(vrsConnection, hotNICEntityPort)
	if err != nil {
		t.Fatal("Unable to obtain an IP address for the VM from VRS")
	}

	// Verifying port got an IP on VSD
	hotNICIPOnVSD, vsdError := util.VerifyVSDPortResolution(Root, Enterprise, Domain, Zone, hotNICEntityPort)
	if vsdError != nil || hotNICIPOnVSD == "" || hotNICIPOnVSD == "0.0.0.0" {
		t.Fatal("IP resolution for port " + hotNICEntityPort + " failed on VSD.")
	} else {
		t.Logf("VM %s port %s got resolved with an IP address %s on VSD", vmInfo["name"], hotNICEntityPort, portIPOnVSD)
	}

	if portInfo.IPAddr == "0.0.0.0" || portInfo.IPAddr == "" {
		t.Fatalf("Unable to resolve VM %s with new port", vmInfo["name"])
	}

	hotNICIPOnVRS := portInfo.IPAddr
	t.Logf("VM %s got resolved with an IP address %s successfully", vmInfo["name"], hotNICIPOnVRS)

	// Comparing port's IP address on VRS and VSD
	if hotNICIPOnVSD != hotNICIPOnVRS {
		t.Fatal("Port IPs on VSD and VRS do not match.")
	}

	err = cleanup(vrsConnection, vmInfo)
	if err != nil {
		t.Fatal("Unable to delete port from OVSDB table")
	}

	err = vrsConnection.DeregisterForPortUpdates(hotNICEntityPort)
	if err != nil {
		t.Fatal("Unable to deregister for hot nic port updates")
	}

	// Performing cleanup of newly added HOT NIC port on VRS
	err = vrsConnection.DestroyPort(hotNICEntityPort)
	if err != nil {
		t.Fatal("Unable to delete newly added port from OVSDB table")
	}

	// Purging out the newly added HOT NIC veth port from VRS alubr0
	err = vrsConnection.RemovePortFromAlubr0(hotNICEntityPort)
	if err != nil {
		t.Fatalf("Unable to delete newly added veth port as part of cleanup from alubr0 %v", err)
	}

	// Cleaning up veth paired ports created for HOT NIC addition from VRS
	err = util.DeleteVETHPair(hotNICEntityPort, hotNICBridgePort)
	if err != nil {
		t.Fatal("Unable to delete veth pairs as a part of cleanup on VRS")
	}

	t.Logf("Waiting for 300 seconds before verifying port gets removed from VSD")
	time.Sleep(time.Duration(300) * time.Second)

	// Verifying port deletion on VSD
	if portDeletionFailure, vsdErr := util.VerifyVSDPortDeletion(Root, Enterprise, Domain, Zone, hotNICEntityPort); vsdErr != nil || portDeletionFailure {
		t.Fatal("Port did not get removed on VSD")
	}

	vrsConnection.Disconnect()
}

// TestVMReconfigure tests that a VM and an associated port gets resolved
// in VRS-VM as well as on the VSD on VM reconfigure event
func TestVMReconfigure(t *testing.T) {

	var vrsConnection VRSConnection
	var err error

	err = util.EnableOVSDBRPCSocket(VRSPort)
	if err != nil {
		t.Fatal("Unable to add an interface to the ovsdb-server on VRS to make it accept RPCs via TCP socket")
	}

	if vrsConnection, err = NewUnixSocketConnection(UnixSocketFile); err != nil {
		t.Fatal("Unable to connect to the VRS")
	}

	vmInfo := make(map[string]string)
	vmInfo["name"] = fmt.Sprintf("Test-VM-%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	vmInfo["mac"] = util.GenerateMAC()
	vmInfo["vmuuid"] = uuid.Generate().String()
	vmInfo["entityport"] = fmt.Sprintf("veth.%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	vmInfo["brport"] = fmt.Sprintf("vethbr.%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	portList := []string{vmInfo["entityport"], vmInfo["brport"]}
	err = util.CreateVETHPair(portList)
	if err != nil {
		t.Fatal("Unable to create veth pairs on VRS")
	}

	var entityInfo EntityInfo
	entityInfo.Name = vmInfo["name"]
	entityInfo.UUID = vmInfo["vmuuid"]
	// Add the paired veth port to alubr0 on VRS
	err = vrsConnection.AddPortToAlubr0(vmInfo["entityport"], entityInfo)
	if err != nil {
		t.Errorf("Error inserting row in alubr0: %v", err)
		t.Fatal("Unable to add veth port to alubr0")
	}

	err = createVM(vrsConnection, vmInfo, entity.ESXI, entity.EventCategoryStarted, entity.EventStartedBooted)
	if err != nil {
		t.Fatal("Unable to create a test VM")
	}

	// Obtaining entity port information from VRS using update mechanism
	portInfo, err := getPortInfo(vrsConnection, vmInfo["entityport"])
	if err != nil {
		t.Fatal("Unable to obtain an IP address for the VM from VRS")
	}

	// Verifying port got an IP on VSD
	portIPOnVSD, vsdError := util.VerifyVSDPortResolution(Root, Enterprise, Domain, Zone, vmInfo["entityport"])
	if vsdError != nil || portIPOnVSD == "" || portIPOnVSD == "0.0.0.0" {
		t.Fatal("IP resolution for port " + vmInfo["entityport"] + " failed on VSD.")
	} else {
		t.Logf("VM %s port %s got resolved with an IP address %s on VSD", vmInfo["name"], vmInfo["entityport"], portIPOnVSD)
	}

	if portInfo.IPAddr == "0.0.0.0" || portInfo.IPAddr == "" {
		t.Fatalf("Unable to resolve VM %s ", vmInfo["name"])
	}

	portIPOnVRS := portInfo.IPAddr
	t.Logf("VM %s got resolved with an IP address %s successfully", vmInfo["name"], portIPOnVRS)

	// Comparing port's IP address on VRS and VSD
	if portIPOnVSD != portIPOnVRS {
		t.Fatal("Port IPs on VSD and VRS do not match.")
	}

	// Notify VRS that VM has been re-configured
	err = vrsConnection.PostEntityEvent(vmInfo["vmuuid"], entity.EventCategoryStopped, entity.EventStoppedShutdown)
	if err != nil {
		t.Fatal("Unable to notify VRS regarding VM shutdown event")
	}

	// Reconfigure VM NIC to be in another subnet
	portReconfigure := make(map[string]string)
	portReconfigure[string(port.MetadataKeyNetwork)] = Network2

	err = vrsConnection.UpdatePortMetadata(vmInfo["entityport"], portReconfigure)
	if err != nil {
		t.Fatal("Unable to re-configure VM port")
	}

	// Notify VRS that VM has been powered on
	err = vrsConnection.PostEntityEvent(vmInfo["vmuuid"], entity.EventCategoryStarted, entity.EventStartedBooted)
	if err != nil {
		t.Fatal("Unable to notify VRS regarding VM powered ON event")
	}

	// Notify VRS that VM has been re-configured
	err = vrsConnection.PostEntityEvent(vmInfo["vmuuid"], entity.EventCategoryStarted, entity.EventDefinedUpdated)
	if err != nil {
		t.Fatal("Unable to notify VRS regarding VM re-configure event")
	}

	t.Logf("Waiting for 5 seconds before verifying re-configured VM port got resolved")
	time.Sleep(time.Duration(5) * time.Second)

	// Verifying port got an IP on VSD
	reconfiguredPortIPOnVSD, vsdError := util.VerifyVSDPortResolution(Root, Enterprise, Domain, Zone, vmInfo["entityport"])
	if vsdError != nil || reconfiguredPortIPOnVSD == "" || reconfiguredPortIPOnVSD == "0.0.0.0" {
		t.Fatal("IP resolution for port " + vmInfo["entityport"] + " failed on VSD.")
	} else {
		t.Logf("VM %s port %s got resolved with an IP address %s on VSD", vmInfo["name"], vmInfo["entityport"], reconfiguredPortIPOnVSD)
	}

	portState, err := vrsConnection.GetPortState(vmInfo["entityport"])

	if err != nil {
		t.Fatal("Unable to query re-configured VM port state on VRS")
	}

	if portState[port.StateKeyIPAddress] == "0.0.0.0" || portState[port.StateKeyIPAddress] == "" {
		t.Fatalf("Unable to resolve re-configured VM %s ", vmInfo["name"])
	}

	reconfiguredPortIPOnVRS := portState[port.StateKeyIPAddress]

	if reconfiguredPortIPOnVRS == portIPOnVRS {
		t.Fatal("VM port failed to re-configure and resolve with an IP in the new VSD network")
	}

	// Comparing port's IP address on VRS and VSD
	if reconfiguredPortIPOnVSD != reconfiguredPortIPOnVRS {
		t.Fatal("Port IPs on VSD and VRS do not match.")
	}

	if strings.Contains((reconfiguredPortIPOnVRS.(string)), ReconfNWPrefix) {
		t.Logf("Re-configured VM %s got resolved with an IP address %s successfully", vmInfo["name"], reconfiguredPortIPOnVRS)
	} else {
		t.Fatal("VM port failed to re-configure and resolve with an IP in the new VSD network")
	}

	// Performing cleanup of port/entity on VRS
	err = cleanup(vrsConnection, vmInfo)
	if err != nil {
		t.Fatal("Unable to delete port from OVSDB table")
	}

	t.Logf("Waiting for 300 seconds before verifying deleted port entry gets removed")
	time.Sleep(time.Duration(300) * time.Second)

	portState, _ = vrsConnection.GetPortState(vmInfo["entityport"])

	if _, ok := portState[port.StateKeyIPAddress]; ok {
		t.Fatal("Entry for deleted VM Port still present in OVSDB table")
	}

	// Verifying port deletion on VSD
	if portDeletionFailure, vsdErr := util.VerifyVSDPortDeletion(Root, Enterprise, Domain, Zone, vmInfo["entityport"]); vsdErr != nil || portDeletionFailure {
		t.Fatal("Port did not get removed on VSD")
	}

	t.Logf("VM %s got removed from VRS successfully", vmInfo["name"])
}

// TestVMPowerOff tests that a VM and an associated port gets resolved
// in VRS-VM as well as on the VSD and gets removed from VRS and VSD when deleted
func TestVMPowerOff(t *testing.T) {

	var vrsConnection VRSConnection
	var err error

	err = util.EnableOVSDBRPCSocket(VRSPort)
	if err != nil {
		t.Fatal("Unable to add an interface to the ovsdb-server on VRS to make it accept RPCs via TCP socket")
	}

	if vrsConnection, err = NewUnixSocketConnection(UnixSocketFile); err != nil {
		t.Fatal("Unable to connect to the VRS")
	}

	vmInfo := make(map[string]string)
	vmInfo["name"] = fmt.Sprintf("Test-VM-%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	vmInfo["mac"] = util.GenerateMAC()
	vmInfo["vmuuid"] = uuid.Generate().String()
	vmInfo["entityport"] = fmt.Sprintf("veth.%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	vmInfo["brport"] = fmt.Sprintf("vethbr.%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	portList := []string{vmInfo["entityport"], vmInfo["brport"]}
	err = util.CreateVETHPair(portList)
	if err != nil {
		t.Fatal("Unable to create veth pairs on VRS")
	}

	var entityInfo EntityInfo
	entityInfo.Name = vmInfo["name"]
	entityInfo.UUID = vmInfo["vmuuid"]
	// Add the paired veth port to alubr0 on VRS
	err = vrsConnection.AddPortToAlubr0(vmInfo["entityport"], entityInfo)
	if err != nil {
		t.Errorf("Error inserting row in alubr0: %v", err)
		t.Fatal("Unable to add veth port to alubr0")
	}

	err = createVM(vrsConnection, vmInfo, entity.ESXI, entity.EventCategoryStarted, entity.EventStartedBooted)
	if err != nil {
		t.Fatal("Unable to create a test VM")
	}

	// Obtaining entity port information from VRS using update mechanism
	portInfo, err := getPortInfo(vrsConnection, vmInfo["entityport"])
	if err != nil {
		t.Fatal("Unable to obtain an IP address for the VM from VRS")
	}

	// Verifying port got an IP on VSD
	portIPOnVSD, vsdError := util.VerifyVSDPortResolution(Root, Enterprise, Domain, Zone, vmInfo["entityport"])
	if vsdError != nil || portIPOnVSD == "" || portIPOnVSD == "0.0.0.0" {
		t.Fatal("IP resolution for port " + vmInfo["entityport"] + " failed on VSD.")
	} else {
		t.Logf("VM %s port %s got resolved with an IP address %s on VSD", vmInfo["name"], vmInfo["entityport"], portIPOnVSD)
	}

	if portInfo.IPAddr == "0.0.0.0" || portInfo.IPAddr == "" {
		t.Fatalf("Unable to resolve VM %s ", vmInfo["name"])
	}

	portIPOnVRS := portInfo.IPAddr
	t.Logf("VM %s got resolved with an IP address %s successfully", vmInfo["name"], portIPOnVRS)

	// Comparing port's IP address on VRS and VSD
	if portIPOnVSD != portIPOnVRS {
		t.Fatal("Port IPs on VSD and VRS do not match.")
	} else {
		t.Logf("Port IPs on VSD and VRS match.")
	}

	// Notify VRS that VM has been powered off
	err = vrsConnection.SetEntityState(vmInfo["vmuuid"], entity.Shutoff, entity.ShutoffShutdown)
	if err != nil {
		t.Fatal("Unable to shut down entity")
	}

	t.Logf("Waiting for 30 seconds before sending powered off event for shutdown VM")
	time.Sleep(time.Duration(30) * time.Second)

	// Notify VRS that VM has been powered off
	err = vrsConnection.PostEntityEvent(vmInfo["vmuuid"], entity.EventCategoryShutdown, entity.EventStoppedShutdown)
	if err != nil {
		t.Fatal("Unable to notify VRS regarding VM shutdown event")
	}
	err = vrsConnection.DeregisterForPortUpdates(vmInfo["entityport"])
	if err != nil {
		t.Fatal("Unable to deregister for port updates")
	}
	// Removing powered off VM port
	err = vrsConnection.DestroyPort(vmInfo["entityport"])
	if err != nil {
		t.Fatal("Unable to delete port from OVSDB table")
	}

	// Purging out the veth port from VRS alubr0
	err = vrsConnection.RemovePortFromAlubr0(vmInfo["entityport"])
	if err != nil {
		t.Fatalf("Unable to delete veth port for powered off VM %v", err)
	}

	err = vrsConnection.DestroyEntity(vmInfo["vmuuid"])
	if err != nil {
		t.Fatal("Unable to remove the entity from OVSDB table")
	}

	// t.Logf("Waiting for 60 seconds before verifying port gets removed from VRS")
	// time.Sleep(time.Duration(60) * time.Second)

	portState := make(map[port.StateKey]interface{})
	portState, _ = vrsConnection.GetPortState(vmInfo["entityport"])

	if _, ok := portState[port.StateKeyIPAddress]; ok {
		t.Fatal("Entry for deleted VM Port still present in OVSDB table")
	}

	t.Logf("VM %s got removed from VRS successfully", vmInfo["name"])

	// Cleaning up veth paired ports from VRS
	err = util.DeleteVETHPair(vmInfo["entityport"], vmInfo["brport"])
	if err != nil {
		t.Fatal("Unable to delete veth pairs as a part of cleanup on VRS")
	}

	vrsConnection.Disconnect()
}

//TestSplitActivation tests the split activation mode using SDK
func TestSplitActivation(t *testing.T) {

	vrsConnection, err1 := NewUnixSocketConnection(UnixSocketFile)
	if err1 != nil {
		t.Fatal("Unable to connect to the VRS")
	}

	enterprise, err := util.FetchEnterprise(Root, Enterprise)
	if err != nil {
		t.Fatalf("%v", err)
	}

	domain, err := util.FetchDomain(enterprise, Domain)
	if err != nil {
		t.Fatalf("%v", err)
	}

	subnet, err := util.FetchSubnet(domain, Network1)
	if err != nil {
		t.Fatalf("%v", err)
	}

	vportName := fmt.Sprintf("vethentity-%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	containerName := fmt.Sprintf("test_container_%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	containerUUID := strings.Replace(uuid.Generate().String(), "-", "", -1)
	containerUUID = containerUUID + strings.Replace(uuid.Generate().String(), "-", "", -1)
	macAddress := util.GenerateMAC()
	t.Logf("vport = %s, container name = %s, containeruuid = %s, macAddress = %s", vportName, containerName, containerUUID, macAddress)
	// new container interface that matches to vport
	containerInterface := vspk.NewContainerInterface()
	containerInterface.Name = vportName
	containerInterface.MAC = macAddress
	containerInterface.AttachedNetworkID = subnet.ID
	interfaceList := make([]interface{}, 1)
	interfaceList[0] = containerInterface
	// create a new container under a user
	container := vspk.NewContainer()
	container.UUID = containerUUID
	container.Name = containerName
	container.Interfaces = interfaceList
	err = Root.CreateContainer(container)
	if err != nil {
		t.Fatalf("Creating container failed with error: %v", err)
	}

	containerInfo := make(map[string]string)
	containerInfo["name"] = containerName
	containerInfo["mac"] = macAddress
	containerInfo["vmuuid"] = containerUUID
	containerInfo["entityport"] = vportName
	containerInfo["brport"] = fmt.Sprintf("vethbr-%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100))
	portList := []string{containerInfo["brport"], containerInfo["entityport"]}
	if err := util.CreateVETHPair(portList); err != nil {
		t.Fatal("Unable to create veth pairs on VRS")
	}

	var entityInfo EntityInfo
	entityInfo.Name = containerInfo["name"]
	entityInfo.UUID = containerInfo["vmuuid"]
	var err2 error
	// Add the paired veth port to alubr0 on VRS
	err2 = vrsConnection.AddPortToAlubr0(containerInfo["entityport"], entityInfo)
	if err2 != nil {
		t.Errorf("Error inserting row in alubr0: %v", err)
		t.Fatal("Unable to add veth port to alubr0")
	}

	if err := splitActivationCreateContainer(vrsConnection, containerInfo, entity.Container, entity.EventCategoryStarted, entity.EventStartedBooted); err != nil {
		t.Fatal("Unable to create a test VM")
	}

	portInfo, err1 := getPortInfo(vrsConnection, containerInfo["entityport"])
	if err1 != nil {
		t.Errorf("Getting port info failed with error : %v", err1)
	}

	if portInfo.IPAddr != containerInterface.IPAddress {
		t.Errorf("Container IP address does not match in port table")
	}

	err = container.Delete()
	if err != nil {
		t.Fatalf("Container deletion failed with error: %v", err)
	}
	if err := cleanup(vrsConnection, containerInfo); err != nil {
		t.Fatalf("cleanup failed with error: %v", err)
	}
	vrsConnection.Disconnect()
}
