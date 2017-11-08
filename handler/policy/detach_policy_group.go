package policy

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/group"
	"github.com/go-iam/handler/util"
	"net/http"
)

type GroupDetachPolicyApi struct {
	req        *http.Request
	status     int
	err        error
	policy     string
	policyType PolicyType
	policyId   string
	group      string
	groupId    string
	account    string
}

func (udpa *GroupDetachPolicyApi) Parse() {
	params := util.ParseParameters(udpa.req)
	udpa.policy = params["PolicyName"]
	udpa.policyType = ParsePolicyType(params["PolicyType"])
	udpa.group = params["GroupName"]
}

func (udpa *GroupDetachPolicyApi) Validate() {
	if udpa.policy == "" {
		udpa.err = MissPolicyNameError
		udpa.status = http.StatusBadRequest
		return
	}
	if !IsValidType(udpa.policyType) {
		udpa.err = MissPolicyTypeError
		udpa.status = http.StatusBadRequest
		return
	}
	if udpa.group == "" {
		udpa.err = group.MissGroupNameError
		udpa.status = http.StatusBadRequest
		return
	}
	if ok, err := group.IsGroupNameValid(udpa.group); !ok {
		udpa.err = err
		udpa.status = http.StatusBadRequest
		return
	}
	if ok, err := IsPolicyNameValid(udpa.policy); !ok {
		udpa.err = err
		udpa.status = http.StatusBadRequest
		return
	}
}

func (udpa *GroupDetachPolicyApi) Auth() {
	udpa.err = doAuth(udpa.req)
	if udpa.err != nil {
		udpa.status = http.StatusForbidden
	}
}

func (udpa *GroupDetachPolicyApi) Response() {
	json := simplejson.New()
	if udpa.err != nil {
		context.Set(udpa.req, "request_error", gerror.NewIAMError(udpa.status, udpa.err))
		json.Set("ErrorMessage", udpa.err.Error())
	}
	json.Set("RequestId", context.Get(udpa.req, "request_id"))
	data, _ := json.Encode()
	context.Set(udpa.req, "response", data)
}

func (udpa *GroupDetachPolicyApi) detachPolicyFromGroup() {
	groupId, err := group.GetGroupId(udpa.account, udpa.group)
	if err != nil {
		udpa.err = err
		return
	}
	udpa.groupId = groupId

	policyId, err := GetPolicyId(udpa.account, udpa.policy, udpa.policyType)
	if err != nil {
		udpa.err = err
		return
	}
	udpa.policyId = policyId

	bean := &db.PolicyGroupBean{
		PolicyId: udpa.policyId,
		GroupId:  udpa.groupId,
	}
	udpa.err = db.ActiveService().GroupDetachPolicy(bean)
	if udpa.err != nil {
		if udpa.err == db.PolicyNotAttachedGroupError {
			udpa.status = http.StatusNotFound
		} else {
			udpa.status = http.StatusInternalServerError
		}
		return
	}
}

func GroupDetachPolicyHandler(w http.ResponseWriter, r *http.Request) {
	udpa := GroupDetachPolicyApi{req: r, status: http.StatusOK}
	defer udpa.Response()

	if udpa.Auth(); udpa.err != nil {
		return
	}

	if udpa.Parse(); udpa.err != nil {
		return
	}

	if udpa.Validate(); udpa.err != nil {
		return
	}

	if udpa.detachPolicyFromGroup(); udpa.err != nil {
		return
	}
}
