package policy

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/user"
	"github.com/go-iam/handler/util"
	"net/http"
)

type UserDetachPolicyApi struct {
	req        *http.Request
	status     int
	err        error
	policy     string
	policyType PolicyType
	policyId   string
	user       string
	userId     string
	account    string
}

func (udpa *UserDetachPolicyApi) Parse() {
	params := util.ParseParameters(udpa.req)
	udpa.policy = params["PolicyName"]
	udpa.policyType = ParsePolicyType(params["PolicyType"])
	udpa.user = params["UserName"]
}

func (udpa *UserDetachPolicyApi) Validate() {
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
	if udpa.user == "" {
		udpa.err = user.MissUserNameError
		udpa.status = http.StatusBadRequest
		return
	}
}

func (udpa *UserDetachPolicyApi) Auth() {
	udpa.err = doAuth(udpa.req)
	if udpa.err != nil {
		udpa.status = http.StatusForbidden
	}
}

func (udpa *UserDetachPolicyApi) Response() {
	json := simplejson.New()
	if udpa.err != nil {
		context.Set(udpa.req, "request_error", gerror.NewIAMError(udpa.status, udpa.err))
		json.Set("ErrorMessage", udpa.err.Error())
	}
	json.Set("RequestId", context.Get(udpa.req, "request_id"))
	data, _ := json.Encode()
	context.Set(udpa.req, "response", data)
}

func (udpa *UserDetachPolicyApi) detachPolicyFromUser() {
	bean := &db.PolicyUserBean{
		PolicyId: udpa.policyId,
		UserId:   udpa.userId,
	}
	udpa.err = db.ActiveService().UserDetachPolicy(bean)
	if udpa.err != nil {
		if udpa.err == db.PolicyNotAttachedUserError {
			udpa.status = http.StatusNotFound
		} else {
			udpa.status = http.StatusInternalServerError
		}
		return
	}
}

func UserDetachPolicyHandler(w http.ResponseWriter, r *http.Request) {
	udpa := UserDetachPolicyApi{req: r, status: http.StatusOK}
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

	userId, err := user.GetUserId(udpa.account, udpa.user)
	if err != nil {
		udpa.err = err
		return
	}
	udpa.userId = userId

	policyId, err := GetPolicyId(udpa.account, udpa.policy, udpa.policyType)
	if err != nil {
		udpa.err = err
		return
	}
	udpa.policyId = policyId

	if udpa.detachPolicyFromUser(); udpa.err != nil {
		return
	}
}
