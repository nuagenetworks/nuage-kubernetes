package client

import (
	"github.com/golang/glog"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/config"
	"net/http"
	"net/url"
	"sync"
)

type PodList struct {
	list           map[string]*SubnetNode   // namespace/podName -> specific subnet
	namespaces     map[string]NamespaceData // namespace name -> data
	editLock       sync.RWMutex
	newSubnetQueue chan config.NamespaceUpdateRequest // send a subnet on this channel to request that another subnet be created after it in the list
}

type podListJson struct {
	SubnetName string
}

func NewPodList(namespaces map[string]NamespaceData, updateChan chan config.NamespaceUpdateRequest) *PodList {
	pods := PodList{}
	pods.list = make(map[string]*SubnetNode)
	pods.namespaces = namespaces
	pods.newSubnetQueue = updateChan
	return &pods
}

func (pods PodList) Get(urlVars map[string]string, values url.Values,
	header http.Header) (int, interface{}, http.Header) {
	namespace, exists := urlVars["namespace"]
	if !exists {
		return http.StatusNotFound, nil, nil
	}
	podName, exists := urlVars["podName"]
	if !exists {
		// In the future, maybe return a list of all pods, but that should be
		// unnecessary at the moment. Assume any GET requests without a
		// specific pod name to be erroneous.
		return http.StatusNotFound, nil, nil
	}
	pods.editLock.RLock()
	defer pods.editLock.RUnlock()
	if item, exists := pods.list[namespace+"/"+podName]; exists {
		return http.StatusOK, podListJson{SubnetName: item.SubnetName}, nil
	}
	return http.StatusNotFound, nil, nil
}

func (pods PodList) Post(urlVars map[string]string, values url.Values,
	header http.Header, bodyJson map[string]interface{}) (int, interface{},
	http.Header) {
	namespace, exists := urlVars["namespace"]
	if !exists {
		glog.Error("Namespace info missing for the POST request")
		return http.StatusNotFound, nil, nil
	}
	podName, exists := bodyJson["podName"]
	if !exists {
		glog.Error("Podname missing from the JSON data")
		return http.StatusBadRequest, nil, nil
	}
	podNameString, isString := podName.(string)
	if !isString || podNameString == "" {
		glog.Error("Invalid pod name")
		return http.StatusBadRequest, nil, nil
	}

	desiredZone, zoneSpecified:= bodyJson["desiredZone"]
	if zoneSpecified {
		desiredZoneStr, isString := desiredZone.(string)
		if !isString || desiredZoneStr == "" {
			glog.Error("Invalid zone name")
			return http.StatusBadRequest, nil, nil
		}
		glog.Info("Specified zone: %s", desiredZoneStr)

		desiredSubnet, subnetSpecified := bodyJson["desiredSubnet"]
		if !subnetSpecified {
			glog.Error("Invalid request: Subnet absent")
			return http.StatusBadRequest, nil, nil;
		}

		desiredSubnetStr, isString := desiredSubnet.(string)
		if !isString || desiredSubnetStr == "" {
			glog.Error("Invalid subnet name")
			return http.StatusBadRequest, nil, nil;
		}

		_, nuageMonManagedZone := pods.namespaces[desiredZoneStr]
		if !nuageMonManagedZone {
			return http.StatusOK, podListJson{SubnetName: desiredSubnetStr}, nil
		} else {
			glog.Error("Invalid zone parameter: Zone controlled by Nuage Monitor")
			return http.StatusBadRequest, nil, nil
		}
	}

	// lock for reading. Holding RLock() doesn't block other RLock() calls, but
	// does block Lock() calls, which should be used for writing.  We can't
	// defer RUnlock() here, because we need to RUnlock() then Lock() before
	// writing, and a double RUnlock() causes a panic
	pods.editLock.RLock()
	if podData, exists := pods.list[namespace+"/"+podNameString]; exists {
		pods.editLock.RUnlock()
		return http.StatusConflict, podListJson{SubnetName: podData.SubnetName}, nil
	}
	pods.editLock.RUnlock()
	// After this point, objects received from pods.namespaces may be edited, so
	// RLock() is no longer sufficient.  Lock() must be used.
	pods.editLock.Lock()
	defer pods.editLock.Unlock()
	nsData, exists := pods.namespaces[namespace]
	if !exists {
		glog.Warningf("Attempted to create a pod in namespace %q, but no"+
			" namespace was found.", namespace)
		// TODO: handle case where kubernetes creates the namespace and pod
		// before the namespace's create event is handled by the vsd client
		return http.StatusNotFound, nil, nil
	}
	nsSubnetsHead := nsData.Subnets
	if nsSubnetsHead == nil {
		glog.Warningf(
			"Namespace %q was found, but didn't contain any subnets",
			namespace)
		return http.StatusInternalServerError, nil, nil
	}
	for currentNode := nsSubnetsHead; currentNode != nil; currentNode = currentNode.Next {
		// total available IPs, minus broadcast (e.g. a.b.c.255/24), the network
		// ID (e.g. a.b.c.0/24), and a space for a gateway (e.g. a.b.c.1 or
		// equivalent, usually)
		maxIPs := 1<<(uint(32-currentNode.Subnet.CIDRMask)) - 3
		if currentNode.ActiveIPs < maxIPs {
			pods.list[namespace+"/"+podNameString] = currentNode
			currentNode.ActiveIPs++
			if currentNode.Next == nil && float64(currentNode.ActiveIPs)/float64(maxIPs) > 0.8 {
				// If this is the last node (.next is nil), then all other
				// subnets are full. If more than 80% of the final subnet's IPs
				// are allocated, create another subnet.
				if !nsData.NeedsNewSubnet {
					glog.Infof("Asking for new subnet in namespace %q", namespace)
					nsData.NeedsNewSubnet = true
					pods.namespaces[namespace] = nsData
					pods.newSubnetQueue <- config.NamespaceUpdateRequest{
						NamespaceID: namespace,
						Event:       config.AddSubnet,
					}
				}
			}
			return http.StatusOK, podListJson{SubnetName: currentNode.SubnetName}, nil
		}
	}
	// All subnets were full. Return an internal error for now?
	// TODO: force create a new subnet
	glog.Warningf("All subnets in namespace %q are full", namespace)
	return http.StatusInternalServerError, nil, nil
}

func (pods PodList) Delete(urlVars map[string]string, values url.Values,
	header http.Header) (int, interface{}, http.Header) {
	namespace, exists := urlVars["namespace"]
	if !exists {
		return http.StatusNotFound, nil, nil
	}
	podName, exists := urlVars["podName"]
	if !exists {
		return http.StatusNotFound, nil, nil
	}
	pods.editLock.Lock()
	defer pods.editLock.Unlock()
	if subnetNode, exists := pods.list[namespace+"/"+podName]; exists {
		subnetNode.ActiveIPs--
		//TODO: Check if the subnet is no longer necessary. If so, delete it.
	}
	delete(pods.list, namespace+"/"+podName)
	return http.StatusOK, nil, nil
}
