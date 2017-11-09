package policy

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/db"
	"regexp"
)

type Policy struct {
	policyId    string
	policyName  string
	policyType  PolicyType
	document    string
	description string
	version     string
	createDate  string
	updateDate  string
	account     string
}

func (p *Policy) ToBean() db.PolicyBean {
	return db.PolicyBean{
		PolicyName:  p.policyName,
		PolicyType:  int(p.policyType),
		Document:    p.document,
		Description: p.description,
		Version:     p.version,
		CreateDate:  p.createDate,
		UpdateDate:  p.updateDate,
		Account:     p.account,
	}
}

func FromBean(bean *db.PolicyBean) Policy {
	return Policy{
		policyId:    bean.PolicyId.Hex(),
		policyName:  bean.PolicyName,
		policyType:  PolicyType(bean.PolicyType),
		document:    bean.Document,
		description: bean.Description,
		version:     bean.Version,
		createDate:  bean.CreateDate,
		updateDate:  bean.UpdateDate,
		account:     bean.Account,
	}
}

func (policy *Policy) Json() *simplejson.Json {
	j := simplejson.New()
	j.Set("PolicyId", policy.policyId)
	j.Set("PolicyName", policy.policyName)
	j.Set("PolicyType", policy.policyType.String())
	j.Set("Description", policy.description)
	j.Set("DefaultVersion", policy.version)
	j.Set("CreateDate", policy.createDate)
	j.Set("UpdateDate", policy.updateDate)
	return j
}

var (
	PolicyNameTooLongError     = errors.New("policy name beyond the length limit")
	DescriptionTooLongError    = errors.New("description beyond the length limit")
	PolicyDocumentTooLongError = errors.New("policy document beyond the length limit")
	PolicyNameInvalidError     = errors.New("policy name contains invalid char")
)

const (
	MaxPolicyNameLength     = 128
	MaxDescriptionLength    = 1024
	MaxPolicyDocumentLength = 128
)

func IsPolicyNameValid(policyName string) (bool, error) {
	if len(policyName) > MaxPolicyNameLength {
		return false, PolicyNameTooLongError
	}

	reg := `^[a-zA-Z0-9-]+$`
	rgx := regexp.MustCompile(reg)
	if !rgx.MatchString(policyName) {
		return false, PolicyNameInvalidError
	}

	return true, nil
}

func IsPolicyDocumentValid(document string) (bool, error) {
	if len(document) > MaxPolicyDocumentLength {
		return false, PolicyDocumentTooLongError
	}

	// TODO
	return true, nil
}

func IsDescriptionValid(description string) (bool, error) {
	if description == "" {
		return true, nil
	}

	if len(description) > MaxDescriptionLength {
		return false, DescriptionTooLongError
	}

	return true, nil
}

func (policy *Policy) validate() error {
	ok, err := IsPolicyNameValid(policy.policyName)
	if !ok {
		return err
	}

	ok, err = IsDescriptionValid(policy.description)
	if !ok {
		return err
	}

	ok, err = IsPolicyDocumentValid(policy.document)
	if !ok {
		return err
	}

	return nil
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
