package translator

import (
	"fmt"
	"strconv"

	"github.com/golang/glog"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	xlateApi "github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/apis/translate"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/policyapi/policies"
	kapi "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const priorityLabel = "nuage.io/priority"

//CreateNuagePGPolicy translates NetworkPolicyEvent to Nuage VSP Constructs
func (rm *ResourceManager) CreateNuagePGPolicy(pe *api.NetworkPolicyEvent) (*policies.NuagePolicy, error) {

	priority, err := getPriorityFromLabels(pe)
	if err != nil {
		return nil, fmt.Errorf("Getting priority from labels failed %v", err)
	}

	nuagePolicy := policies.NewNuagePolicy(
		rm.vsdMetaData["enterprise"],
		rm.vsdMetaData["domain"],
		pe.Name,
		pe.Name,
		priority,
	)

	policyElements, err := rm.convertNetworkPolicySpec(pe)
	if err != nil {
		return nil, fmt.Errorf("converting network policy spec to policy elements failed %v", err)
	}

	glog.Infof("Policy elements created %+v ", policyElements)
	nuagePolicy.PolicyElements = policyElements
	return &nuagePolicy, nil
}

func (rm *ResourceManager) convertNetworkPolicySpec(pe *api.NetworkPolicyEvent) ([]policies.DefaultPolicyElement, error) {
	glog.Infof("Translating(%s) network policy spec %+v", pe.Name, pe.Policy)

	policyElements := []policies.DefaultPolicyElement{}
	k8sNetworkPolicySpec := &pe.Policy
	targetSelector, err := metav1.LabelSelectorAsSelector(&k8sNetworkPolicySpec.PodSelector)
	if err != nil {
		return policyElements, fmt.Errorf("Cannot get label selector as selector")
	}

	targetPG, ok := rm.vsdObjsMap.PGMap[pe.Namespace][targetSelector.String()]
	if !ok {
		return policyElements, fmt.Errorf("Target Pod policy group information missing %+v", targetSelector.String())
	}

	for _, ingressRule := range k8sNetworkPolicySpec.Ingress {
		tmpPolicyElements, err := rm.convertPeerPolicyElements(ingressRule.From,
			ingressRule.Ports,
			pe.Namespace,
			targetPG.PgName,
			pe.Name,
			true)
		if err != nil {
			glog.Errorf("converting k8s ingress policy to nuage policy failed: %v", err)
			return policyElements, err
		}

		policyElements = append(policyElements, tmpPolicyElements...)
	}

	for _, egressRule := range k8sNetworkPolicySpec.Egress {
		tmpPolicyElements, err := rm.convertPeerPolicyElements(egressRule.To,
			egressRule.Ports,
			pe.Namespace,
			targetPG.PgName,
			pe.Name,
			false)
		if err != nil {
			glog.Errorf("converting k8s egress policy to nuage policy failed: %v", err)
			return policyElements, err
		}

		policyElements = append(policyElements, tmpPolicyElements...)
	}
	return policyElements, nil
}

func createPolicyElements(ports []networkingV1.NetworkPolicyPort,
	p *xlateApi.PolicyData) []policies.DefaultPolicyElement {

	policyElement := policies.DefaultPolicyElement{}
	policyElements := []policies.DefaultPolicyElement{}

	glog.Infof("Creating Nuage policy objects as per specified Kubernetes ingress/egress policies")

	if p.Action == policies.Deny {
		for _, proto := range []kapi.Protocol{kapi.ProtocolTCP, kapi.ProtocolUDP} {
			protocol := policies.TCP
			if proto == kapi.ProtocolUDP {
				protocol = policies.UDP
			}
			policyElement := policies.NewPolicyElement(
				fmt.Sprintf("%s-%s-deny", p.Name, proto),
				p.SourceType, p.SourceName,
				p.TargetType, p.TargetName,
				p.Action,
				protocol,
				1, 65535,
			)
			glog.Infof("Adding policy event %+v", policyElement)
			policyElements = append(policyElements, policyElement)
		}
	} else {
		for idx, targetPort := range ports {
			port := targetPort.Port.IntValue()
			protocol := policies.TCP
			if *targetPort.Protocol == kapi.ProtocolUDP {
				protocol = policies.UDP
			}

			policyElement = policies.NewPolicyElement(
				fmt.Sprintf("%s-%d", p.Name, idx),
				p.SourceType, p.SourceName,
				p.TargetType, p.TargetName,
				p.Action,
				protocol,
				port, port)

			glog.Infof("Adding policy event %+v", policyElement)
			policyElements = append(policyElements, policyElement)
		}
	}

	return policyElements
}

func (rm *ResourceManager) convertPeerPolicyElements(peers []networkingV1.NetworkPolicyPeer,
	ports []networkingV1.NetworkPolicyPort,
	policyNamespace string,
	targetPgName string,
	policyName string,
	ingress bool) ([]policies.DefaultPolicyElement, error) {

	var policyElements []policies.DefaultPolicyElement

	for _, peer := range peers {
		var tmpPolicyElements []policies.DefaultPolicyElement
		//Both namespace and pod selectors are specified
		if peer.NamespaceSelector != nil && peer.PodSelector != nil {
			nsSelectorLabel, err := metav1.LabelSelectorAsSelector(peer.NamespaceSelector)
			if err != nil {
				glog.Errorf("Extracting namespace label failed %v", err)
				return nil, err
			}

			namespaces, _ := rm.vsdObjsMap.NSLabelsMap[nsSelectorLabel.String()]
			for _, namespace := range namespaces {
				tmpPolicyElements, err = rm.convertPodSelector(peer,
					ports,
					namespace,
					policyName,
					targetPgName,
					ingress)
				if err != nil {
					glog.Errorf("for namespace %s converting pod selector failed %v", namespace, err)
					return policyElements, err
				}
				policyElements = append(policyElements, tmpPolicyElements...)

			}
			continue
		}

		//only namespace selector is specified
		if peer.NamespaceSelector != nil {
			tmpPolicyElements, err := rm.convertNSSelector(peer,
				ports,
				policyName,
				targetPgName,
				ingress)
			if err != nil {
				glog.Errorf("converting namespace selector failed %v", err)
				return policyElements, err
			}
			policyElements = append(policyElements, tmpPolicyElements...)
			continue
		}

		//only pod selector is specifed
		if peer.PodSelector != nil {
			tmpPolicyElements, err := rm.convertPodSelector(peer,
				ports,
				policyNamespace,
				policyName,
				targetPgName,
				ingress)
			if err != nil {
				glog.Errorf("for namespace %s converting pod selector failed %v", policyNamespace, err)
				return policyElements, err
			}
			policyElements = append(policyElements, tmpPolicyElements...)
			continue
		}

		// only IP Block is specifed
		if peer.IPBlock != nil {
			tmpPolicyElements, err := rm.convertIPBlock(peer,
				ports,
				policyName,
				targetPgName,
				ingress)
			if err != nil {
				glog.Errorf("converting ip block %v failed %v", peer.IPBlock, err)
				return policyElements, err
			}
			policyElements = append(policyElements, tmpPolicyElements...)
			continue
		}
	}

	glog.Infof("Default policy elements created: %v", policyElements)

	return policyElements, nil
}

func getPriorityFromLabels(pe *api.NetworkPolicyEvent) (int, error) {
	priorityStr, ok := pe.Labels[priorityLabel]
	if !ok {
		return 0, fmt.Errorf("Priority missing for the network policy labels")
	}

	priority, err := strconv.Atoi(priorityStr)
	if err != nil {
		return 0, fmt.Errorf("Invalid priority value %s in the network policy labels", priorityStr)
	}

	return priority, nil
}

func (rm *ResourceManager) convertPodSelector(peer networkingV1.NetworkPolicyPeer,
	ports []networkingV1.NetworkPolicyPort,
	namespace string,
	policyName string,
	targetPgName string,
	ingress bool) ([]policies.DefaultPolicyElement, error) {

	var ok bool
	var err error
	var pgInfo *xlateApi.PgInfo
	var sourceSelector labels.Selector
	policyElements := []policies.DefaultPolicyElement{}

	if sourceSelector, err = metav1.LabelSelectorAsSelector(peer.PodSelector); err != nil {
		return policyElements, fmt.Errorf("converting label selector to selector failed for %s", peer.PodSelector.String())
	}
	if pgInfo, ok = rm.vsdObjsMap.PGMap[namespace][sourceSelector.String()]; !ok {
		return policyElements, fmt.Errorf("Policy group missing for %s", sourceSelector.String())
	}

	policyElements = createPolicyElementsUtil(ports,
		&xlateApi.PolicyData{
			Name:       policyName,
			SourceType: policies.PolicyGroup,
			SourceName: pgInfo.PgName,
			TargetType: policies.PolicyGroup,
			TargetName: targetPgName,
			Action:     policies.Allow,
		},
		ingress)

	return policyElements, nil
}

func (rm *ResourceManager) convertNSSelector(peer networkingV1.NetworkPolicyPeer,
	ports []networkingV1.NetworkPolicyPort,
	policyName string,
	targetPgName string,
	ingress bool) ([]policies.DefaultPolicyElement, error) {

	policyElements := []policies.DefaultPolicyElement{}
	//for each of the namespace create a new policy element
	nsSelectorLabel, err := metav1.LabelSelectorAsSelector(peer.NamespaceSelector)
	if err != nil {
		glog.Errorf("Extracting namespace label failed %v", err)
		return policyElements, err
	}

	namespaces, _ := rm.vsdObjsMap.NSLabelsMap[nsSelectorLabel.String()]
	for _, namespace := range namespaces {
		policyElements = append(policyElements,
			createPolicyElementsUtil(ports,
				&xlateApi.PolicyData{
					Name:       policyName,
					SourceType: policies.Zone,
					SourceName: namespace,
					TargetType: policies.PolicyGroup,
					TargetName: targetPgName,
					Action:     policies.Allow,
				},
				ingress,
			)...)
	}

	return policyElements, nil
}

func (rm *ResourceManager) convertIPBlock(peer networkingV1.NetworkPolicyPeer,
	ports []networkingV1.NetworkPolicyPort,
	policyName string,
	targetPgName string,
	ingress bool) ([]policies.DefaultPolicyElement, error) {

	cidrInfo, _ := rm.vsdObjsMap.NWMacroMap[peer.IPBlock.CIDR]
	policyData := &xlateApi.PolicyData{
		Name:       policyName,
		SourceType: policies.NetworkMacro,
		SourceName: cidrInfo.Name,
		TargetType: policies.PolicyGroup,
		TargetName: targetPgName,
		Action:     policies.Allow,
	}

	policyElements := createPolicyElementsUtil(ports, policyData, ingress)

	policyData.Action = policies.Deny
	for _, exceptCIDR := range peer.IPBlock.Except {
		cidrInfo, _ := rm.vsdObjsMap.NWMacroMap[exceptCIDR]
		policyData.SourceName = cidrInfo.Name
		policyElements = append(policyElements, createPolicyElementsUtil(ports, policyData, ingress)...)
	}
	return policyElements, nil
}

func createPolicyElementsUtil(ports []networkingV1.NetworkPolicyPort, policyData *xlateApi.PolicyData,
	ingress bool) []policies.DefaultPolicyElement {

	tmpPolicyElements := []policies.DefaultPolicyElement{}
	if ingress {
		tmpPolicyElements = createPolicyElements(ports, policyData)
	} else {
		policyData.SourceType, policyData.TargetType = policyData.TargetType, policyData.SourceType
		policyData.SourceName, policyData.TargetName = policyData.TargetName, policyData.SourceName
		tmpPolicyElements = createPolicyElements(ports, policyData)
	}
	return tmpPolicyElements
}
