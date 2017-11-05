package policy

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
)

type ListPolicyEntityApi struct {
	req        *http.Request
	status     int
	err        error
	policy     string
	policyType PolicyType
	policyId   string
	users      []*PolicyUser
	groups     []*PolicyGroup
	account    string
}

func (lpea *ListPolicyEntityApi) Parse() {
	params := util.ParseParameters(lpea.req)
	lpea.policy = params["PolicyName"]
	lpea.policyType = ParsePolicyType(params["PolicyType"])
}

func (lpea *ListPolicyEntityApi) Validate() {
	if lpea.policy == "" {
		lpea.err = MissPolicyNameError
		lpea.status = http.StatusBadRequest
		return
	}
	if !IsValidType(lpea.policyType) {
		lpea.err = InvalidPolicyTypeError
		lpea.status = http.StatusBadRequest
		return
	}
}

func (lpea *ListPolicyEntityApi) Auth() {
	lpea.err = doAuth(lpea.req)
	if lpea.err != nil {
		lpea.status = http.StatusForbidden
	}
}

func (lpea *ListPolicyEntityApi) Response() {
	json := simplejson.New()
	if lpea.err != nil {
		context.Set(lpea.req, "request_error", gerror.NewIAMError(lpea.status, lpea.err))
		json.Set("ErrorMessage", lpea.err.Error())
	} else {
		uJsons := make([]*simplejson.Json, 0)
		for _, u := range lpea.users {
			j := u.Json()
			uJsons = append(uJsons, j)
		}

		gJsons := make([]*simplejson.Json, 0)
		for _, g := range lpea.groups {
			j := g.Json()
			gJsons = append(gJsons, j)
		}
		json.Set("Users", uJsons)
		json.Set("Groups", gJsons)
	}
	json.Set("RequestId", context.Get(lpea.req, "request_id"))
	data, _ := json.Encode()
	context.Set(lpea.req, "response", data)
}

func (lpea *ListPolicyEntityApi) listPolicyEntity() {
	users := make([]*db.PolicyUserBean, 0)
	lpea.err = db.ActiveService().ListPolicyUser(lpea.policyId, &users)
	if lpea.err != nil {
		lpea.status = http.StatusInternalServerError
		return
	}
	for _, u := range users {
		var user db.UserBean
		lpea.err = db.ActiveService().GetIamUserById(u.UserId, &user)
		if lpea.err != nil {
			return
		}
		pu := &PolicyUser{
			userId:     user.UserId.Hex(),
			userName:   user.UserName,
			attachDate: u.AttachDate,
		}
		lpea.users = append(lpea.users, pu)
	}

	groups := make([]*db.PolicyGroupBean, 0)
	lpea.err = db.ActiveService().ListPolicyGroup(lpea.policyId, &groups)
	if lpea.err != nil {
		lpea.status = http.StatusInternalServerError
		return
	}
	for _, g := range groups {
		var group db.GroupBean
		lpea.err = db.ActiveService().GetGroupById(g.GroupId, &group)
		if lpea.err != nil {
			return
		}
		pg := &PolicyGroup{
			groupName:  group.GroupName,
			attachDate: g.AttachDate,
		}
		lpea.groups = append(lpea.groups, pg)
	}
}

func ListPolicyEntityHandler(w http.ResponseWriter, r *http.Request) {
	lpea := ListPolicyEntityApi{req: r, status: http.StatusOK}
	defer lpea.Response()

	if lpea.Auth(); lpea.err != nil {
		return
	}

	if lpea.Parse(); lpea.err != nil {
		return
	}

	if lpea.Validate(); lpea.err != nil {
		return
	}

	policyId, err := GetPolicyId(lpea.account, lpea.policy, lpea.policyType)
	if err != nil {
		lpea.err = err
		return
	}
	lpea.policyId = policyId

	if lpea.listPolicyEntity(); lpea.err != nil {
		return
	}
}
