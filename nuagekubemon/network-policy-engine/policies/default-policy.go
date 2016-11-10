package policies

type DefaultPolicyElement struct {
	Name              string
	From              EndPoint `yaml:"from"`
	To                EndPoint `yaml:"to"`
	Action            ActionType
	NetworkParameters NetworkParameters `yaml:"network-parameters"`
}
