package port

// MetadataKey represent Nuage specific metadata key
type MetadataKey string

// Different keys supported for Nuage Port metadata
const (
	MetadataKeyDomain                 MetadataKey = "nuage-enterprise-domain"
	MetadataKeyNetwork                MetadataKey = "nuage-enterprise-network"
	MetadataKeyNetworkType            MetadataKey = "nuage-enterprise-network-type"
	MetadataKeyZone                   MetadataKey = "nuage-enterprise-zone"
	MetadataKeyStaticIP               MetadataKey = "static-ip"
	MetadataKeyNuageVPort             MetadataKey = "nuage-vport"
	MetadataKeyNuageVPortTag          MetadataKey = "nuage-vport-tag"
	MetadataKeyNuageRedirectionTarget MetadataKey = "nuage-redirection-target"
	MetadataNuagePolicyGroup          MetadataKey = "nuage-policy-group"
	MetadataNuageSubnetAddress        MetadataKey = "nuage-subnet-address"
	MetadataNuageSubnetMask           MetadataKey = "nuage-subnet-mask"
	MetadataNuageSubnetGateway        MetadataKey = "nuage-subnet-gateway"
	MetadataKeyPg                     MetadataKey = "pg"
	MetadataKeyPortBindings           MetadataKey = "nuage-port-mapping"
)
