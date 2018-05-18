package translator

import (
	"fmt"
	"strconv"

	"github.com/golang/glog"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/policyapi/policies"
	kapi "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const priorityLabel = "nuage.io/priority"

func CreateNuagePGPolicy(
	pe *api.NetworkPolicyEvent,
	policyGroupMap map[string]api.PgInfo,
	nuageMetadata map[string]string,
	namespaceLabelsMap map[string][]string,
	ipBlockCidrMap map[string]int) (*policies.NuagePolicy, error) {

	k8sNetworkPolicySpec := &pe.Policy
	policyName := pe.Name
	policyLabels := pe.Labels

	if k8sNetworkPolicySpec == nil || policyGroupMap == nil || nuageMetadata == nil {
		return nil, fmt.Errorf("Invalid arguments")
	}

	glog.Infof("Policy Name %+v ", policyName)
	glog.Infof("Translating network policy spec %+v with labels %+v",
		k8sNetworkPolicySpec, policyLabels)
	glog.Infof("Policy Group Map %+v ", policyGroupMap)

	var ok bool
	var enterprise string
	if enterprise, ok = nuageMetadata["enterpriseName"]; !ok {
		return nil, fmt.Errorf("Enterprise missing from metadata")
	}

	var domain string
	if domain, ok = nuageMetadata["domainName"]; !ok {
		return nil, fmt.Errorf("Domain missing from metadata")
	}

	var targetPG api.PgInfo
	if targetSelector, err := metav1.LabelSelectorAsSelector(&k8sNetworkPolicySpec.PodSelector); err == nil {
		if targetPG, ok = policyGroupMap[targetSelector.String()]; !ok {
			return nil, fmt.Errorf("Target Pod policy group information missing %+v", targetSelector.String())
		}
	} else {
		return nil, fmt.Errorf("Cannot get label selector as selector")
	}

	var priorityStr string
	if priorityStr, ok = pe.Labels[priorityLabel]; !ok {
		return nil, fmt.Errorf("Priority missing for the network policy labels")
	}

	var priority int
	var err error
	if priority, err = strconv.Atoi(priorityStr); err != nil {
		return nil, fmt.Errorf("Invalid priority value %s in the network policy labels", priorityStr)
	}

	nuagePolicy := policies.NuagePolicy{
		Version:    policies.V1Alpha,
		Type:       policies.Default,
		Enterprise: enterprise,
		Domain:     domain,
		Name:       policyName,
		ID:         policyName,
		Priority:   priority,
	}

	var defaultPolicyElements []policies.DefaultPolicyElement

	for _, ingressRule := range k8sNetworkPolicySpec.Ingress {
		tmpPolicyElements, err := convertPeerPolicyElements(ingressRule.From, ingressRule.Ports,
			namespaceLabelsMap, ipBlockCidrMap, targetPG.PgName, policyName, true, policyGroupMap)
		if err != nil {
			glog.Errorf("converting k8s ingress policy to nuage policy failed: %v", err)
			return nil, err
		}

		defaultPolicyElements = append(defaultPolicyElements, tmpPolicyElements...)
	}

	for _, egressRule := range k8sNetworkPolicySpec.Egress {
		tmpPolicyElements, err := convertPeerPolicyElements(egressRule.To, egressRule.Ports,
			namespaceLabelsMap, ipBlockCidrMap, targetPG.PgName, policyName, false, policyGroupMap)
		if err != nil {
			glog.Errorf("converting k8s egress policy to nuage policy failed: %v", err)
			return nil, err
		}

		defaultPolicyElements = append(defaultPolicyElements, tmpPolicyElements...)
	}

	nuagePolicy.PolicyElements = defaultPolicyElements

	return &nuagePolicy, nil
}

func createPolicyElements(ports []networkingV1.NetworkPolicyPort, policyName string,
	sourceType policies.EndPointType, sourceName string,
	targetType policies.EndPointType, targetName string) ([]policies.DefaultPolicyElement, error) {

	var policyElements []policies.DefaultPolicyElement
	var policyElement policies.DefaultPolicyElement
	if len(ports) == 0 {
		protocol := kapi.ProtocolTCP
		port := intstr.FromInt(0)
		//if nothing is specified, this is the default for policy
		ports = append(ports, networkingV1.NetworkPolicyPort{
			Protocol: &protocol,
			Port:     &port,
		})
	}
	for idx, targetPort := range ports {
		if targetPort.Port == nil {
			return nil, fmt.Errorf("Received nil value for port number for non-nil ports section")
		}
		port := targetPort.Port.IntValue()
		targetProtocol := policies.TCP
		if *targetPort.Protocol == kapi.ProtocolUDP {
			targetProtocol = policies.UDP
		} else if *targetPort.Protocol == kapi.ProtocolTCP {
			targetProtocol = policies.TCP
		}

		policyElement = policies.DefaultPolicyElement{
			Name: fmt.Sprintf("%s-%d", policyName, idx),
			From: policies.EndPoint{Type: sourceType,
				Name: sourceName},
			To: policies.EndPoint{Type: targetType,
				Name: targetName},
			Action: policies.Allow,
			NetworkParameters: policies.NetworkParameters{
				Protocol: targetProtocol,
				DestinationPortRange: policies.PortRange{StartPort: port,
					EndPort: port}},
		}

		glog.Infof("Adding policy event %+v", policyElement)
		policyElements = append(policyElements, policyElement)
	}
	return policyElements, nil
}

func convertPeerPolicyElements(peers []networkingV1.NetworkPolicyPeer, ports []networkingV1.NetworkPolicyPort,
	namespaceLabelsMap map[string][]string, ipBlockCidrMap map[string]int, targetPgName string, policyName string, ingress bool,
	policyGroupMap map[string]api.PgInfo) ([]policies.DefaultPolicyElement, error) {
	var defaultPolicyElements []policies.DefaultPolicyElement
	for _, peer := range peers {
		var ok bool
		var err error
		var pgInfo api.PgInfo
		var sourceSelector labels.Selector
		var tmpPolicyElements []policies.DefaultPolicyElement
		if peer.NamespaceSelector != nil {
			//for each of the namespace create a new policy element
			namespaces, _ := namespaceLabelsMap[peer.NamespaceSelector.String()]
			for _, namespace := range namespaces {
				if ingress {
					tmpPolicyElements, err = createPolicyElements(ports, policyName,
						policies.Zone, namespace, policies.PolicyGroup, targetPgName)
				} else {
					tmpPolicyElements, err = createPolicyElements(ports, policyName,
						policies.PolicyGroup, targetPgName, policies.Zone, namespace)
				}
				if err != nil {
					glog.Errorf("creating namespace policy elements failed: %v", err)
					return nil, err
				}
				defaultPolicyElements = append(defaultPolicyElements, tmpPolicyElements...)
			}
			continue
		}

		if peer.IPBlock != nil {
			//for each of the ip cidr create a new policy element
			for nwMacroName, _ := range ipBlockCidrMap {
				if ingress {
					tmpPolicyElements, err = createPolicyElements(ports, policyName,
						policies.NetworkMacro, nwMacroName, policies.PolicyGroup, targetPgName)
				} else {
					tmpPolicyElements, err = createPolicyElements(ports, policyName,
						policies.PolicyGroup, targetPgName, policies.NetworkMacro, nwMacroName)
				}
				if err != nil {
					glog.Errorf("creating ip block cidr policy elements failed: %v", err)
					return nil, err
				}
				defaultPolicyElements = append(defaultPolicyElements, tmpPolicyElements...)
			}
			continue
		}

		if sourceSelector, err = metav1.LabelSelectorAsSelector(peer.PodSelector); err != nil {
			return nil, fmt.Errorf("converting label selector to selector failed for %s", peer.PodSelector.String())
		}
		if pgInfo, ok = policyGroupMap[sourceSelector.String()]; !ok {
			return nil, fmt.Errorf("Policy group missing for %s", sourceSelector.String())
		}

		if ingress {
			tmpPolicyElements, err = createPolicyElements(ports, policyName,
				policies.PolicyGroup, pgInfo.PgName, policies.PolicyGroup, targetPgName)
		} else {
			tmpPolicyElements, err = createPolicyElements(ports, policyName,
				policies.PolicyGroup, targetPgName, policies.PolicyGroup, pgInfo.PgName)
		}
		if err != nil {
			glog.Errorf("creating podselector policy elements failed: %v", err)
			return nil, err
		}
		defaultPolicyElements = append(defaultPolicyElements, tmpPolicyElements...)
	}
	return defaultPolicyElements, nil
}
