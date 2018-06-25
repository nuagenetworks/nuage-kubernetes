package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/nuagenetworks/libvrsdk/api"
	"github.com/nuagenetworks/libvrsdk/api/entity"
	"github.com/nuagenetworks/libvrsdk/api/port"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	yaml "gopkg.in/yaml.v2"
)

// Constants used by Util
const (
	Bridge        = "alubr0"
	vrsSocketFile = "/var/run/openvswitch/db.sock"
	LOGFILE       = "/var/log/cni/nuage-infra-pod.log"
)

type logTextFormatter log.TextFormatter

// VRSBConfig stores the config information for vrsb
type VRSBConfig struct {
	Name     string
	UUID     string
	Metadata struct {
		Username    string
		Enterprise  string
		Domain      string
		Zone        string
		Subnet      string
		NetworkType string
	}
	Interface struct {
		Veth1 string
		Veth2 string
	}
}

func main() {
	var cleanup bool
	var configYAML string

	flag.StringVar(&configYAML, "config", "", "Config file for VRSB")
	flag.BoolVar(&cleanup, "cleanup", false, "Cleanup mode")
	flag.Parse()

	setupLogging()

	yamlContent, err := ioutil.ReadFile(configYAML)
	if err != nil {
		log.Errorf("Unable to read the config YAML")
		return
	}

	log.Infof("Read yamlContent %s", string(yamlContent))

	config := VRSBConfig{}
	err = yaml.Unmarshal(yamlContent, &config)
	if err != nil {
		log.Errorf("error: %v", err)
		return
	}

	log.Infof("Configuring a container vport with %+v", config)

	vrsConnection, err := api.NewUnixSocketConnection(vrsSocketFile)
	defer vrsConnection.Disconnect()
	if err != nil {
		log.Errorf("error: %v, unable to connect to the VRS OVSDB server", err)
	}

	cleanUpStaleEntries(&config, vrsConnection)

	if cleanup {
		log.Infof("Invoked in cleanup mode. Done cleaning up.. Exiting..")
		return
	}

	localVethPair := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{Name: config.Interface.Veth1},
		PeerName:  config.Interface.Veth2,
	}

	if err := netlink.LinkAdd(localVethPair); err != nil {
		log.Errorf("Failed to create veth paired port for container %s with error %v", config.Name, err)
		return
	}

	if err := setupUpLink(config.Interface.Veth1); err != nil {
		log.Errorf("%v", err)
		return
	}
	if err := setupUpLink(config.Interface.Veth2); err != nil {
		log.Errorf("%v", err)
		return
	}
	mac, err := findMAC(config.Interface.Veth2)
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	var entityInfo api.EntityInfo
	entityInfo.Name = config.Name
	entityInfo.UUID = config.UUID
	err = vrsConnection.AddPortToAlubr0(config.Interface.Veth1, entityInfo)
	if err != nil {
		log.Fatal("Error inserting row in alubr0: %v", err)
		return
	}

	log.Infof("Added port to %+v", config)

	// Create Port Attributes
	portAttributes := port.Attributes{
		Platform: entity.Docker,
		MAC:      mac,
		Bridge:   Bridge,
	}

	// Create Port Metadata
	portMetadata := make(map[port.MetadataKey]string)
	portMetadata[port.MetadataKeyDomain] = config.Metadata.Domain
	portMetadata[port.MetadataKeyNetwork] = config.Metadata.Subnet
	portMetadata[port.MetadataKeyZone] = config.Metadata.Zone
	portMetadata[port.MetadataKeyNetworkType] = config.Metadata.NetworkType

	// Associate one veth port to entity
	err = vrsConnection.CreatePort(config.Interface.Veth1, portAttributes, portMetadata)
	if err != nil {
		log.Errorf("Unable to create entity port %v", err)
		return
	}

	// Create VM metadata
	vmMetadata := make(map[entity.MetadataKey]string)
	vmMetadata[entity.MetadataKeyUser] = config.Metadata.Username
	vmMetadata[entity.MetadataKeyEnterprise] = config.Metadata.Enterprise

	// Define ports associated with the VM
	ports := []string{config.Interface.Veth1}

	// Add entity to the VRS
	entityInfo = api.EntityInfo{
		UUID:     config.UUID,
		Name:     config.Name,
		Type:     entity.Container,
		Domain:   entity.Docker,
		Ports:    ports,
		Metadata: vmMetadata,
	}

	// Post a event while bringing up the container
	events := &entity.EntityEvents{}
	events.EntityEventCategory = entity.EventCategoryStarted
	events.EntityEventType = entity.EventStartedBooted
	events.EntityState = entity.Running
	events.EntityReason = entity.RunningBooted
	entityInfo.Events = events

	err = vrsConnection.CreateEntity(entityInfo)
	if err != nil {
		log.Errorf("Unable to add entity to VRS %v", err)
		return
	}

	// Registering for VRS port updates
	portInfoUpdateChan := make(chan *api.PortIPv4Info)
	err = vrsConnection.RegisterForPortUpdates(config.Interface.Veth1, portInfoUpdateChan)
	if err != nil {
		log.Errorf("Failed to register for updates from VRS for entity port %s", config.Interface.Veth1)
		return
	}
	ticker := time.NewTicker(10 * time.Second)
	portInfo := &api.PortIPv4Info{}
	select {
	case portInfo = <-portInfoUpdateChan:
	case <-ticker.C:
		log.Errorf("Failed to receive an update from VRS for entity port %s", config.Interface.Veth1)
		return
	}

	if err = assignIPAddress(config.Interface.Veth2, portInfo.IPAddr, portInfo.Mask); err != nil {
		log.Errorf("assigning ip to container failed: %v", err)
		return
	}
	// De-registering for VRS port updates
	err = vrsConnection.DeregisterForPortUpdates(config.Interface.Veth1)
	if err != nil {
		log.Errorf("Error de-registering for port updates from VRS for entity port %s", config.Interface.Veth1)
		return
	}

	fmt.Printf("%s", portInfo.Gateway)
}

func setupUpLink(name string) error {
	eth, errStr := netlink.LinkByName(name)
	if errStr != nil {
		return fmt.Errorf("Failed to lookup port %s", name)
	}
	err := netlink.LinkSetUp(eth)
	if err != nil {
		return fmt.Errorf("Error setting up veth %s", eth.Attrs().Name)
	}
	return nil
}

func findMAC(name string) (string, error) {
	eth, errStr := netlink.LinkByName(name)
	if errStr != nil {
		return "", fmt.Errorf("Failed to lookup port %s", name)
	}
	return eth.Attrs().HardwareAddr.String(), nil
}

func assignIPAddress(name, ip, mask string) error {
	netmask := net.IPMask(net.ParseIP(mask))
	prefixSize, _ := netmask.Size()
	ipNet := net.IPNet{IP: net.ParseIP(ip), Mask: net.CIDRMask(prefixSize, 32)}

	veth, errStr := netlink.LinkByName(name)
	if errStr != nil {
		return fmt.Errorf("failed to lookup %s: %v", name, errStr)
	}

	if err := netlink.AddrAdd(veth, &netlink.Addr{IPNet: &ipNet}); err != nil {
		return fmt.Errorf("Failed to assign IP %s: %v", &ipNet, err)
	}
	return nil
}

//SetupLogging sets up logging infrastructure
func setupLogging() {
	customFormatter := new(logTextFormatter)
	log.SetFormatter(customFormatter)
	log.SetOutput(&lumberjack.Logger{
		Filename: LOGFILE,
		MaxSize:  1,
		MaxAge:   30,
	})
	log.SetLevel(log.DebugLevel)
}

//Format custom format method used by logrus. Prints in standard nuage log format
func (f *logTextFormatter) Format(entry *log.Entry) ([]byte, error) {
	outputString := fmt.Sprintf("|%v|%s|%s|%s\n", entry.Time, strings.ToUpper(log.Level.String(entry.Level)), path.Base(os.Args[0]), entry.Message)
	for k, v := range entry.Data {
		outputString += fmt.Sprintf("|%s=%s|", k, v)
	}
	return []byte(outputString), nil
}

func cleanUpStaleEntries(config *VRSBConfig, vrsConnection api.VRSConnection) {
	if err := vrsConnection.DestroyEntityByVMName(config.Name); err != nil {
		log.Errorf("destroying entity entry failed with error: %v", err)
	}

	if err := vrsConnection.DestroyPort(config.Interface.Veth1); err != nil {
		log.Errorf("destroying port entry failed with error: %v", err)
	}

	if err := vrsConnection.RemovePortFromAlubr0(config.Interface.Veth1); err != nil {
		log.Errorf("removing port from bridge failed with error: %v", err)
	}

	localVethPair := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{Name: config.Interface.Veth1},
		PeerName:  config.Interface.Veth2,
	}

	err := netlink.LinkDel(localVethPair)
	if err != nil {
		log.Errorf("Deleting veth pair %+v failed with error: %s", localVethPair, err)
	}
}
