package policy

import (
	"github.com/bitly/go-simplejson"
)

type Policy struct {
	policyId    string
	policyName  string
	policyType  string
	document    string
	description string
	version     string
	createDate  string
	updateDate  string
	account     string
}

func (policy *Policy) Json() *simplejson.Json {
	j := simplejson.New()
	j.Set("PolicyId", policy.policyId)
	j.Set("PolicyName", policy.policyName)
	j.Set("PolicyType", policy.policyType)
	j.Set("Description", policy.description)
	j.Set("DefaultVersion", policy.version)
	j.Set("CreateDate", policy.createDate)
	j.Set("UpdateDate", policy.createDate)
	return j
}
