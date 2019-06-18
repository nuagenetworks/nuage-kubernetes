package policy

import (
	"github.com/golang/glog"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/subnet"
	networkingV1 "k8s.io/api/networking/v1"
)

//NWMacroMap map of network macros on VSD per zone
type NWMacroMap map[string]int

//NWMacroExceptMap map of network macros for except CIDR on VSD per zone
type NWMacroExceptMap map[string]int

func (rm *ResourceManager) createNetworkMacros(ipBlock *networkingV1.IPBlock, pe *api.NetworkPolicyEvent) error {

	if ipBlock == nil {
		return nil
	}

	allowCidr := ipBlock.CIDR
	exceptList := ipBlock.Except

	if err := rm.checkAndCreateNM(allowCidr); err != nil {
		return err
	}

	for _, exceptCidr := range exceptList {
		if err := rm.checkAndCreateNM(exceptCidr); err != nil {
			return err
		}
	}

	return nil
}

func (rm *ResourceManager) checkAndCreateNM(cidr string) error {
	if _, ok := rm.networkMacroMap[cidr]; !ok {
		nm, err := rm.createNetworkMacroObject(cidr)
		if err != nil {
			glog.Errorf("creating network macro object failed: %v", err)
			return err
		}
		if _, err := rm.callBacks.AddNetworkMacro(nm); err != nil {
			glog.Errorf("adding network macro to VSD failed: %v", err)
			return err
		}
	} else {
		rm.networkMacroMap[cidr]++
	}
	return nil
}

func (rm *ResourceManager) createNetworkMacroObject(cidr string) (*api.VsdNetworkMacro, error) {
	macroCidr, err := subnet.IPv4SubnetFromString(cidr)
	if err != nil {
		glog.Errorf("Failed converting cidr string to cidr object: %v\n", err)
		return nil, err
	}

	return &api.VsdNetworkMacro{
		Name:       cidr,
		IPType:     "IPV4",
		Address:    macroCidr.Address.String(),
		Netmask:    macroCidr.Netmask().String(),
		ExternalID: rm.externalID,
	}, nil
}
