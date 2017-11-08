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

type ListUserPolicyApi struct {
	req      *http.Request
	status   int
	err      error
	user     string
	userId   string
	policies []*UserPolicy
	account  string
}

func (luga *ListUserPolicyApi) Parse() {
	params := util.ParseParameters(luga.req)
	luga.user = params["UserName"]
}

func (luga *ListUserPolicyApi) Validate() {
	if luga.user == "" {
		luga.err = user.MissUserNameError
		luga.status = http.StatusBadRequest
		return
	}
	if ok, err := user.IsUserNameValid(luga.user); !ok {
		luga.err = err
		luga.status = http.StatusBadRequest
		return
	}
}

func (luga *ListUserPolicyApi) Auth() {
	luga.err = doAuth(luga.req)
	if luga.err != nil {
		luga.status = http.StatusForbidden
	}
}

func (luga *ListUserPolicyApi) Response() {
	json := simplejson.New()
	if luga.err != nil {
		gerr := gerror.NewIAMError(luga.status, luga.err)
		context.Set(luga.req, "request_error", gerr)
		json.Set("ErrorMessage", luga.err.Error())
	} else {
		jsons := make([]*simplejson.Json, 0)
		for _, p := range luga.policies {
			j := p.Json()
			jsons = append(jsons, j)
		}
		json.Set("Policies", jsons)
		json.Set("UserName", luga.user)
	}
	json.Set("RequestId", context.Get(luga.req, "request_id"))
	data, _ := json.Encode()
	context.Set(luga.req, "response", data)
}

func (luga *ListUserPolicyApi) listUserPolicy() {
	userId, err := user.GetUserId(luga.account, luga.user)
	if err != nil {
		luga.err = err
		return
	}
	luga.userId = userId

	beans := make([]*db.PolicyUserBean, 0)
	luga.err = db.ActiveService().ListUserPolicy(luga.userId, &beans)
	if luga.err != nil {
		luga.status = http.StatusInternalServerError
		return
	}
	for _, bean := range beans {
		var policy db.PolicyBean
		luga.err = db.ActiveService().GetPolicyById(bean.PolicyId, &policy)
		if luga.err != nil {
			return
		}
		up := &UserPolicy{
			policyName:  policy.PolicyName,
			policyType:  PolicyType(policy.PolicyType).String(),
			description: policy.Description,
			version:     policy.Version,
			attachDate:  bean.AttachDate,
		}
		luga.policies = append(luga.policies, up)
	}
}

func ListUserPolicyHandler(w http.ResponseWriter, r *http.Request) {
	luga := ListUserPolicyApi{req: r, status: http.StatusOK}
	defer luga.Response()

	if luga.Auth(); luga.err != nil {
		return
	}

	if luga.Parse(); luga.err != nil {
		return
	}

	if luga.Validate(); luga.err != nil {
		return
	}

	if luga.listUserPolicy(); luga.err != nil {
		return
	}
}
