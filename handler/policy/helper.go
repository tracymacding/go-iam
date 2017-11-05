package policy

func GetPolicyId(account, policy string, ptype PolicyType) (string, error) {
	gpa := GetPolicyApi{}
	gpa.policy.policyName = policy
	gpa.policy.account = account
	if gpa.getPolicy(); gpa.err != nil {
		return "", gpa.err
	}
	return gpa.policy.policyId, nil
}
