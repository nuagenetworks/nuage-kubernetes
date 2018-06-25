package entity

// MetadataKey identifies metadata keys for a Nuage entity
type MetadataKey string

// Different metadata keys supported for Nuage entity
const (
	MetadataKeyEnterprise      MetadataKey = "enterprise"
	MetadataKeyUser            MetadataKey = "user"
	MetadataKeyOrchestrationID MetadataKey = "nuage-orchestrationID"
	MetadataKeySiteID          MetadataKey = "nuage-siteID"
	MetadataKeyDeleteMode      MetadataKey = "nuage-delete-mode"
	MetadataKeyDeleteExpiry    MetadataKey = "nuage-delete-expiry"
	MetadataKeyExtension       MetadataKey = "nuage-extension"
)
