package policy

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
	"time"
)

type CreatePolicyApi struct {
	req    *http.Request
	status int
	err    error
	policy Policy
}

var (
	MissPolicyNameError     = errors.New("policy name missing")
	MissPolicyDocumentError = errors.New("policy document missing")
	TooManyPolicysError     = errors.New("The count of policy beyond the current limits")
)

func (cpa *CreatePolicyApi) Parse() {
	params := util.ParseParameters(cpa.req)
	cpa.policy.policyName = params["PolicyName"]
	cpa.policy.document = params["PolicyDocument"]
	cpa.policy.description = params["PolicyDescription"]
}

func (cpa *CreatePolicyApi) Validate() {
	if cpa.policy.policyName == "" {
		cpa.err = MissPolicyNameError
		cpa.status = http.StatusBadRequest
		return
	}
	if cpa.policy.document == "" {
		cpa.err = MissPolicyDocumentError
		cpa.status = http.StatusBadRequest
		return
	}

	// TODO: check document synax
}

func (cpa *CreatePolicyApi) Auth() {
	cpa.err = doAuth(cpa.req)
	if cpa.err != nil {
		cpa.status = http.StatusForbidden
	}
}

func (cpa *CreatePolicyApi) Response() {
	json := simplejson.New()
	if cpa.err == nil {
		j := cpa.policy.Json()
		json.Set("Policy", j)
	} else {
		context.Set(cpa.req, "request_error", gerror.NewIAMError(cpa.status, cpa.err))
		json.Set("ErrorMessage", cpa.err.Error())
	}
	json.Set("RequestId", context.Get(cpa.req, "request_id"))
	data, _ := json.Encode()
	context.Set(cpa.req, "response", data)
}

const (
	MAX_POLICY_PER_ACCOUNT = 100
)

func (cpa *CreatePolicyApi) createPolicy() {
	cnt := 0
	cnt, cpa.err = db.ActiveService().PolicyCountOfAccount(cpa.policy.account)
	if cpa.err != nil {
		cpa.status = http.StatusInternalServerError
		return
	}

	if cnt >= MAX_POLICY_PER_ACCOUNT {
		cpa.status = http.StatusConflict
		cpa.err = TooManyPolicysError
		return
	}

	now := time.Now().Format(time.RFC3339)
	bean := &db.PolicyBean{
		PolicyName:  cpa.policy.policyName,
		PolicyType:  cpa.policy.policyName,
		Document:    cpa.policy.document,
		Description: cpa.policy.description,
		Version:     cpa.policy.version,
		Account:     cpa.policy.account,
		CreateDate:  now,
		UpdateDate:  now,
	}
	bean, cpa.err = db.ActiveService().CreatePolicy(bean)
	if cpa.err != nil {
		if cpa.err == db.PolicyExistError {
			cpa.status = http.StatusConflict
		} else {
			cpa.status = http.StatusInternalServerError
		}
		return
	}
	cpa.policy.policyId = bean.PolicyId.Hex()
	cpa.policy.createDate = bean.CreateDate
	cpa.policy.updateDate = bean.UpdateDate
}

func CreatePolicyHandler(w http.ResponseWriter, r *http.Request) {
	cpa := CreatePolicyApi{req: r, status: http.StatusOK}
	defer cpa.Response()

	if cpa.Auth(); cpa.err != nil {
		return
	}

	if cpa.Parse(); cpa.err != nil {
		return
	}

	if cpa.Validate(); cpa.err != nil {
		return
	}

	if cpa.createPolicy(); cpa.err != nil {
		return
	}
}
