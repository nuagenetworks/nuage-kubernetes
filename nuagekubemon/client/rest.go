package client

import (
	"github.com/golang/glog"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type PodList struct {
	list             map[string]*SubnetNode   // namespace/podName -> specific subnet
	namespaces       map[string]NamespaceData // namespace name -> data
	editLock         *sync.RWMutex
	podChannel       chan *api.PodEvent
	getPgsFunc       GetPgsFunc
	autoScaleSubnets string
}

type podListJson struct {
	SubnetName   string   `json:"subnetName"`
	PolicyGroups []string `json:"policyGroups"`
}

type restErrorJson struct {
	Error string `json:"error"`
}

type GetPgsFunc func(string, string) (*[]string, error)

func NewPodList(namespaces map[string]NamespaceData, podChannel chan *api.PodEvent, getPgsFunc GetPgsFunc, autoScaleSubnets string) *PodList {
	pods := PodList{}
	pods.list = make(map[string]*SubnetNode)
	pods.namespaces = namespaces
	pods.podChannel = podChannel
	pods.editLock = &sync.RWMutex{}
	pods.getPgsFunc = getPgsFunc
	pods.autoScaleSubnets = strings.ToLower(autoScaleSubnets)
	return &pods
}

func (pods *PodList) Get(urlVars map[string]string, values url.Values,
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

func (pods *PodList) Post(urlVars map[string]string, values url.Values,
	header http.Header, bodyJson map[string]interface{}) (int, interface{},
	http.Header) {
	namespace, exists := urlVars["namespace"]
	if !exists {
		errText := "Namespace info missing for the POST request"
		glog.Error(errText)
		return http.StatusNotFound, restErrorJson{Error: errText}, nil
	}
	podName, exists := bodyJson["podName"]
	if !exists {
		errText := "Podname missing from the JSON data"
		glog.Error(errText)
		return http.StatusBadRequest, restErrorJson{Error: errText}, nil
	}
	podNameString, isString := podName.(string)
	if !isString || podNameString == "" {
		errText := "Invalid pod name"
		glog.Error(errText)
		return http.StatusBadRequest, restErrorJson{Error: errText}, nil
	}

	action, exists := bodyJson["action"]
	if !exists {
		action = "added"
	}

	pgList, err := pods.getPgsFunc(podNameString, namespace)
	if err != nil {
		glog.Error("Couldn't get the policy groups matching this pod")
	}

	desiredZone, zoneSpecified := bodyJson["desiredZone"]
	if zoneSpecified {
		//operations work flow, we dont need to do anything here
		if action == "delete" {
			return http.StatusOK, podListJson{}, nil
		}
		desiredZoneStr, isString := desiredZone.(string)
		if !isString || desiredZoneStr == "" {
			errText := "Invalid zone name"
			glog.Error(errText)
			return http.StatusBadRequest, restErrorJson{Error: errText}, nil
		}
		glog.Info("Specified zone: %s", desiredZoneStr)

		desiredSubnet, subnetSpecified := bodyJson["desiredSubnet"]
		if !subnetSpecified {
			errText := "Invalid request: Subnet absent"
			glog.Error(errText)
			return http.StatusBadRequest, restErrorJson{Error: errText}, nil
		}

		desiredSubnetStr, isString := desiredSubnet.(string)
		if !isString || desiredSubnetStr == "" {
			errText := "Invalid subnet name"
			glog.Error(errText)
			return http.StatusBadRequest, restErrorJson{Error: errText}, nil
		}

		_, nuageMonManagedZone := pods.namespaces[desiredZoneStr]
		if !nuageMonManagedZone {
			return http.StatusOK, podListJson{SubnetName: desiredSubnetStr, PolicyGroups: *pgList}, nil
		} else {
			errText := "Invalid zone parameter: Zone controlled by Nuage Monitor"
			glog.Error(errText)
			return http.StatusBadRequest, restErrorJson{Error: errText}, nil
		}
	}

	//if subnet scaling is not enabled, follow the normal monitor behavior
	if pods.autoScaleSubnets != "1" {
		if action == "delete" {
			return http.StatusOK, podListJson{}, nil
		}
		return http.StatusOK, podListJson{SubnetName: namespace + "-0", PolicyGroups: *pgList}, nil
	}

	podRespChan := make(chan *api.PodEventResp)
	event := &api.PodEvent{
		Type:      api.Added,
		Name:      podNameString,
		Namespace: namespace,
		RespChan:  podRespChan,
	}

	if action == "delete" {
		event.Type = api.Deleted
	}

	pods.podChannel <- event
	resp := <-event.RespChan
	if resp.Error != nil {
		glog.Errorf("Allocating/Deallocating pod to subnet failed: %v", resp.Error)
		return http.StatusInternalServerError, restErrorJson{Error: resp.Error.Error()}, nil
	}
	if action == "delete" {
		return http.StatusOK, podListJson{}, nil
	}
	return http.StatusOK, podListJson{SubnetName: resp.Data.(string), PolicyGroups: *pgList}, nil
}

func (pods *PodList) Delete(urlVars map[string]string, values url.Values,
	header http.Header) (int, interface{}, http.Header) {
	namespace, exists := urlVars["namespace"]
	if !exists {
		return http.StatusNotFound, nil, nil
	}
	podName, exists := urlVars["podName"]
	if !exists {
		return http.StatusNotFound, nil, nil
	}
	glog.Infof("Deleting %s/%s", namespace, podName)
	pods.editLock.Lock()
	defer pods.editLock.Unlock()
	if subnetNode, exists := pods.list[namespace+"/"+podName]; exists {
		subnetNode.ActiveIPs--
		//TODO: Check if the subnet is no longer necessary. If so, delete it.
		delete(pods.list, namespace+"/"+podName)
	}
	glog.Infof("Successfully deleted %s/%s", namespace, podName)
	return http.StatusOK, nil, nil
}
