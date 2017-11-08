package policy

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
)

type GetPolicyApi struct {
	req    *http.Request
	status int
	err    error
	policy Policy
}

func (gpa *GetPolicyApi) Parse() {
	params := util.ParseParameters(gpa.req)
	gpa.policy.policyName = params["PolicyName"]
}

func (gpa *GetPolicyApi) Validate() {
	if gpa.policy.policyName == "" {
		gpa.err = MissPolicyNameError
		gpa.status = http.StatusBadRequest
		return
	}
}

func (gpa *GetPolicyApi) Auth() {
	gpa.err = doAuth(gpa.req)
	if gpa.err != nil {
		gpa.status = http.StatusForbidden
	}
}

func (gpa *GetPolicyApi) Response() {
	json := simplejson.New()
	if gpa.err == nil {
		j := gpa.policy.Json()
		json.Set("User", j)
	} else {
		gerr := gerror.NewIAMError(gpa.status, gpa.err)
		context.Set(gpa.req, "request_error", gerr)
		json.Set("ErrorMessage", gpa.err.Error())
	}
	json.Set("RequestId", context.Get(gpa.req, "request_id"))
	data, _ := json.Encode()
	context.Set(gpa.req, "response", data)
}

func (gpa *GetPolicyApi) getPolicy() {
	var bean db.PolicyBean

	gpa.err = db.ActiveService().GetPolicy(gpa.policy.account, gpa.policy.policyName, &bean)
	if gpa.err != nil {
		if gpa.err == db.PolicyNotExistError {
			gpa.status = http.StatusNotFound
		} else {
			gpa.status = http.StatusInternalServerError
		}
		return
	}
	gpa.policy = FromBean(&bean)
}

func GetPolicyHandler(w http.ResponseWriter, r *http.Request) {
	gpa := GetPolicyApi{req: r, status: http.StatusOK}
	defer gpa.Response()

	if gpa.Auth(); gpa.err != nil {
		return
	}

	if gpa.Parse(); gpa.err != nil {
		return
	}

	if gpa.Validate(); gpa.err != nil {
		return
	}

	if gpa.getPolicy(); gpa.err != nil {
		return
	}
}
