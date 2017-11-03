package policy

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
)

type DeletePolicyApi struct {
	req    *http.Request
	status int
	err    error
	policy Policy
}

func (dpa *DeletePolicyApi) Parse() {
	params := util.ParseParameters(dpa.req)
	dpa.policy.policyName = params["PolicyName"]
}

func (dpa *DeletePolicyApi) Validate() {
	if dpa.policy.policyName == "" {
		dpa.err = MissPolicyNameError
		dpa.status = http.StatusBadRequest
		return
	}
}

func (dpa *DeletePolicyApi) Auth() {
	dpa.err = doAuth(dpa.req)
	if dpa.err != nil {
		dpa.status = http.StatusForbidden
	}
}

func (dpa *DeletePolicyApi) Response() {
	json := simplejson.New()
	if dpa.err != nil {
		context.Set(dpa.req, "request_error", gerror.NewIAMError(dpa.status, dpa.err))
		json.Set("ErrorMessage", dpa.err.Error())
	}
	json.Set("RequestId", context.Get(dpa.req, "request_id"))
	data, _ := json.Encode()
	context.Set(dpa.req, "response", data)
}

func (dpa *DeletePolicyApi) deletePolicy() {
	dpa.err = db.ActiveService().DeletePolicy(dpa.policy.account, dpa.policy.policyName)
	if dpa.err != nil {
		if dpa.err == db.PolicyNotExistError {
			dpa.status = http.StatusNotFound
		} else {
			dpa.status = http.StatusInternalServerError
		}
	}
}

func DeletePolicyHandler(w http.ResponseWriter, r *http.Request) {
	dpa := DeletePolicyApi{req: r, status: http.StatusOK}

	defer dpa.Response()

	if dpa.Auth(); dpa.err != nil {
		return
	}

	if dpa.Parse(); dpa.err != nil {
		return
	}

	if dpa.Validate(); dpa.err != nil {
		return
	}

	if dpa.deletePolicy(); dpa.err != nil {
		return
	}
}
