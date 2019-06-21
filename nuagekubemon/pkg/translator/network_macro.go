package translator

import (
	"github.com/golang/glog"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	xlateApi "github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/apis/translate"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/subnet"
	networkingV1 "k8s.io/api/networking/v1"
)

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
	if _, ok := rm.vsdObjsMap.NWMacroMap[cidr]; !ok {
		nm, err := rm.createNetworkMacroObject(cidr)
		if err != nil {
			glog.Errorf("creating network macro object failed: %v", err)
			return err
		}
		id, err := rm.callBacks.AddNetworkMacro(nm)
		if err != nil {
			glog.Errorf("adding network macro to VSD failed: %v", err)
			return err
		}
		rm.vsdObjsMap.NWMacroMap[cidr] = xlateApi.NWMacroInfo{
			Name: nm.Name,
			ID:   id,
			CIDR: cidr,
		}
	} else {
		nwMacroInfo := rm.vsdObjsMap.NWMacroMap[cidr]
		nwMacroInfo.RefCount = nwMacroInfo.RefCount + 1
		rm.vsdObjsMap.NWMacroMap[cidr] = nwMacroInfo

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
