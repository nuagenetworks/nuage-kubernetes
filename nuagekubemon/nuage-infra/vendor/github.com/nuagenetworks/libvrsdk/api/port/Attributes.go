package port

import "github.com/nuagenetworks/libvrsdk/api/entity"

// Attributes represents the attributes of the Nuage VRS port
type Attributes struct {
	MAC      string
	Platform entity.Domain
	Bridge   string
}
