package policy

type PolicyType int

const (
	PolicyError  PolicyType = 0
	PolicySystem PolicyType = 1
	PolicyCustom PolicyType = 2
)

func (p PolicyType) String() string {
	switch p {
	case PolicySystem:
		return "System"
	case PolicyCustom:
		return "Custom"
	default:
		panic("unknown policy type")
	}
}

func ParsePolicyType(stype string) PolicyType {
	if stype == "System" {
		return PolicySystem
	}
	if stype == "Custom" {
		return PolicyCustom
	}
	return PolicyError
}

func IsValidType(policyType PolicyType) bool {
	return policyType == PolicySystem || policyType == PolicyCustom
}
