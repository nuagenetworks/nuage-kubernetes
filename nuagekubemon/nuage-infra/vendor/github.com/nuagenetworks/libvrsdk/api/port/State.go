package port

import "github.com/nuagenetworks/libvrsdk/ovsdb"

// StateKey represents the keys use to query the state information for a port
type StateKey string

// Keys to query the state of the port
const (
	StateKeyIPAddress  StateKey = ovsdb.NuagePortTableColumnIPAddress
	StateKeySubnetMask StateKey = ovsdb.NuagePortTableColumnSubnetMask
	StateKeyGateway    StateKey = ovsdb.NuagePortTableColumnGateway
	StateKeyVrfID      StateKey = ovsdb.NuagePortTableColumnVRFId
	StateKeyEvpnID     StateKey = ovsdb.NuagePortTableColumnEVPNID
)
