package policies

type PolicyType string

const (
	DEFAULT PolicyType = "default"
)

type VERSION string

const (
	V1Alpha VERSION = "v1-alpha"
)

type NuagePolicy struct {
	Version    VERSION
	Type       PolicyType
	Enterprise string
	Domain     string

	Name           string
	ID             string
	Priority       int
	PolicyElements interface{} `yaml:"policy-elements"`
}
