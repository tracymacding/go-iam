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

type UserPolicy struct {
	policyName  string
	policyType  string
	description string
	version     string
	attachDate  string
}

func (up *UserPolicy) Json() *simplejson.Json {
	j := simplejson.New()
	j.Set("PolicyName", up.policyName)
	j.Set("PolicyType", up.policyType)
	j.Set("Description", up.description)
	j.Set("DefaultVersion", up.version)
	return j
}

type GroupPolicy struct {
	policyName  string
	policyType  string
	description string
	version     string
	attachDate  string
}

func (gp *GroupPolicy) Json() *simplejson.Json {
	j := simplejson.New()
	j.Set("PolicyName", gp.policyName)
	j.Set("PolicyType", gp.policyType)
	j.Set("Description", gp.description)
	j.Set("DefaultVersion", gp.version)
	return j
}

type PolicyUser struct {
	userId     string
	userName   string
	attachDate string
}

func (pu *PolicyUser) Json() *simplejson.Json {
	j := simplejson.New()
	j.Set("UserId", pu.userId)
	j.Set("UserName", pu.userName)
	j.Set("AttachDate", pu.attachDate)
	return j
}

type PolicyGroup struct {
	groupName  string
	attachDate string
}

func (pg *PolicyGroup) Json() *simplejson.Json {
	j := simplejson.New()
	j.Set("GroupName", pg.groupName)
	j.Set("AttachDate", pg.attachDate)
	return j
}
