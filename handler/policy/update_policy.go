package policy

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
)

type UpdatePolicyApi struct {
	req            *http.Request
	status         int
	err            error
	policy         Policy
	newPolicy      string
	newDocument    string
	newDescription string
}

func (upa *UpdatePolicyApi) Parse() {
	params := util.ParseParameters(upa.req)
	upa.policy.policyName = params["PolicyName"]
	if params["NewPolicyName"] != "" {
		upa.newPolicy = params["NewPolicyName"]
	}
	if params["NewPolicyDocument"] != "" {
		upa.newDocument = params["NewPolicyDocument"]
	}
	if params["NewDescription"] != "" {
		upa.newDescription = params["NewDescription"]
	}
}

func (upa *UpdatePolicyApi) Validate() {
	if upa.policy.policyName == "" {
		upa.err = MissPolicyNameError
		upa.status = http.StatusBadRequest
		return
	}
}

func (upa *UpdatePolicyApi) Auth() {
	upa.err = doAuth(upa.req)
	if upa.err != nil {
		upa.status = http.StatusForbidden
	}
}

func (upa *UpdatePolicyApi) Response() {
	json := simplejson.New()
	if upa.err == nil {
		j := upa.policy.Json()
		json.Set("Policy", j)
	} else {
		json.Set("ErrorMessage", upa.err.Error())
		context.Set(upa.req, "request_error", gerror.NewIAMError(upa.status, upa.err))
	}
	json.Set("RequestId", context.Get(upa.req, "request_id"))
	data, _ := json.Encode()
	context.Set(upa.req, "response", data)
}

func (upa *UpdatePolicyApi) updatePolicy() {
	bean := db.PolicyBean{
		PolicyName:  upa.policy.policyName,
		PolicyType:  upa.policy.policyType,
		Account:     upa.policy.account,
		Document:    upa.policy.document,
		Description: upa.policy.description,
		Version:     upa.policy.version,
		CreateDate:  upa.policy.createDate,
		UpdateDate:  upa.policy.updateDate,
	}
	policy, account := upa.policy.policyName, upa.policy.account
	if upa.newPolicy != "" {
		bean.PolicyName = upa.newPolicy
	}
	upa.err = db.ActiveService().UpdatePolicy(policy, account, &bean)
	if upa.err == db.PolicyNotExistError {
		upa.status = http.StatusNotFound
	} else if upa.err == db.PolicyExistError {
		upa.status = http.StatusConflict
	} else {
		upa.status = http.StatusInternalServerError
	}
}

func UpdatePolicyHandler(w http.ResponseWriter, r *http.Request) {
	upa := UpdatePolicyApi{req: r, status: http.StatusOK}
	defer upa.Response()

	if upa.Auth(); upa.err != nil {
		return
	}

	if upa.Parse(); upa.err != nil {
		return
	}

	if upa.Validate(); upa.err != nil {
		return
	}

	gpa := GetPolicyApi{}
	gpa.policy.policyName = upa.policy.policyName
	gpa.policy.account = upa.policy.account

	if gpa.getPolicy(); gpa.err != nil {
		upa.err = gpa.err
		return
	}

	upa.policy.policyType = gpa.policy.policyType
	upa.policy.document = gpa.policy.document
	upa.policy.description = gpa.policy.description
	upa.policy.version = gpa.policy.version
	upa.policy.createDate = gpa.policy.createDate
	upa.policy.updateDate = gpa.policy.updateDate

	if upa.newDocument != "" {
		upa.policy.document = upa.newDocument
	}
	if upa.newDescription != "" {
		upa.policy.description = upa.newDescription
	}

	if upa.updatePolicy(); upa.err != nil {
		return
	}
}
