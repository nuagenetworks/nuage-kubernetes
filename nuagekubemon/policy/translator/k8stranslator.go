package translator

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	"github.com/nuagenetworks/nuagepolicyapi/policies"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kapi "k8s.io/kubernetes/pkg/api"
	"strconv"
)

const priorityLabel = "nuage.io/priority"

func CreateNuagePGPolicy(
	pe *api.NetworkPolicyEvent,
	policyGroupMap map[string]api.PgInfo,
	nuageMetadata map[string]string) (*policies.NuagePolicy, error) {

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
	var defaultPolicyElement policies.DefaultPolicyElement

	for _, ingressRule := range k8sNetworkPolicySpec.Ingress {
		for _, from := range ingressRule.From {
			var fromPG api.PgInfo
			if sourceSelector, err := metav1.LabelSelectorAsSelector(from.PodSelector); err == nil {
				if fromPG, ok = policyGroupMap[sourceSelector.String()]; !ok {
					return nil, fmt.Errorf("Policy group missing for %s", sourceSelector.String())
				}
			} else {
				return nil, fmt.Errorf("Source policy group information was not found for %s", sourceSelector.String())
			}
			if len(ingressRule.Ports) == 0 {
				defaultPolicyElement = policies.DefaultPolicyElement{
					Name: fmt.Sprintf("%s-%d", policyName, 0),
					From: policies.EndPoint{Type: policies.PolicyGroup,
						Name: fromPG.PgName},
					To: policies.EndPoint{Type: policies.PolicyGroup,
						Name: targetPG.PgName},
					Action: policies.Allow,
					NetworkParameters: policies.NetworkParameters{
						Protocol: policies.TCP,
						DestinationPortRange: policies.PortRange{StartPort: 0,
							EndPort: 0}},
				}
				glog.Infof("Adding policy event %+v", defaultPolicyElement)
				defaultPolicyElements = append(defaultPolicyElements, defaultPolicyElement)
			} else {
				for idx, targetPort := range ingressRule.Ports {
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

					defaultPolicyElement = policies.DefaultPolicyElement{
						Name: fmt.Sprintf("%s-%d", policyName, idx),
						From: policies.EndPoint{Type: policies.PolicyGroup,
							Name: fromPG.PgName},
						To: policies.EndPoint{Type: policies.PolicyGroup,
							Name: targetPG.PgName},
						Action: policies.Allow,
						NetworkParameters: policies.NetworkParameters{
							Protocol: targetProtocol,
							DestinationPortRange: policies.PortRange{StartPort: port,
								EndPort: port}},
					}

					glog.Infof("Adding policy event %+v", defaultPolicyElement)
					defaultPolicyElements = append(defaultPolicyElements, defaultPolicyElement)
				}
			}
		}
	}

	nuagePolicy.PolicyElements = defaultPolicyElements

	return &nuagePolicy, nil
}
