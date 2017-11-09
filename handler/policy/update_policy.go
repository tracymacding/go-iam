package policy

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
	"time"
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
	if ok, err := IsPolicyNameValid(upa.policy.policyName); !ok {
		upa.err = err
		upa.status = http.StatusBadRequest
		return
	}
	if upa.newPolicy != "" {
		if ok, err := IsPolicyNameValid(upa.newPolicy); !ok {
			upa.err = err
			upa.status = http.StatusBadRequest
			return
		}
	}
	if upa.newDocument != "" {
		if ok, err := IsPolicyDocumentValid(upa.newDocument); !ok {
			upa.err = err
			upa.status = http.StatusBadRequest
			return
		}
	}
	if upa.newDescription != "" {
		if ok, err := IsDescriptionValid(upa.newDescription); !ok {
			upa.err = err
			upa.status = http.StatusBadRequest
			return
		}
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
		gerr := gerror.NewIAMError(upa.status, upa.err)
		context.Set(upa.req, "request_error", gerr)
		json.Set("ErrorMessage", upa.err.Error())
	}
	json.Set("RequestId", context.Get(upa.req, "request_id"))
	data, _ := json.Encode()
	context.Set(upa.req, "response", data)
}

func (upa *UpdatePolicyApi) updatePolicy() {
	gpa := GetPolicyApi{}
	gpa.policy.policyName = upa.policy.policyName
	gpa.policy.account = upa.policy.account

	if gpa.getPolicy(); gpa.err != nil {
		upa.err = gpa.err
		upa.status = gpa.status
		return
	}

	upa.policy = gpa.policy
	if upa.newPolicy != "" {
		upa.policy.policyName = upa.newPolicy
	}
	if upa.newDocument != "" {
		upa.policy.document = upa.newDocument
	}
	if upa.newDescription != "" {
		upa.policy.description = upa.newDescription
	}
	upa.policy.updateDate = time.Now().Format(time.RFC3339)

	bean := upa.policy.ToBean()
	upa.err = db.ActiveService().UpdatePolicy(upa.policy.policyId, &bean)
	if upa.err == db.PolicyNotExistError {
		upa.status = http.StatusNotFound
	} else if upa.err == db.PolicyExistError {
		upa.status = http.StatusConflict
	} else {
		upa.status = http.StatusInternalServerError
	}
	upa.policy = FromBean(&bean)
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

	if upa.updatePolicy(); upa.err != nil {
		return
	}
}
