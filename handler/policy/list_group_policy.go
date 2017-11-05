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

type ListGroupPolicyApi struct {
	req      *http.Request
	status   int
	err      error
	group    string
	groupId  string
	policies []*GroupPolicy
	account  string
}

func (luga *ListGroupPolicyApi) Parse() {
	params := util.ParseParameters(luga.req)
	luga.group = params["GroupName"]
}

func (luga *ListGroupPolicyApi) Validate() {
	if luga.group == "" {
		luga.err = group.MissGroupNameError
		luga.status = http.StatusBadRequest
		return
	}
}

func (luga *ListGroupPolicyApi) Auth() {
	luga.err = doAuth(luga.req)
	if luga.err != nil {
		luga.status = http.StatusForbidden
	}
}

func (luga *ListGroupPolicyApi) Response() {
	json := simplejson.New()
	if luga.err != nil {
		context.Set(luga.req, "request_error", gerror.NewIAMError(luga.status, luga.err))
		json.Set("ErrorMessage", luga.err.Error())
	} else {
		jsons := make([]*simplejson.Json, 0)
		for _, p := range luga.policies {
			j := p.Json()
			jsons = append(jsons, j)
		}
		json.Set("Policies", jsons)
		json.Set("GroupName", luga.group)
	}
	json.Set("RequestId", context.Get(luga.req, "request_id"))
	data, _ := json.Encode()
	context.Set(luga.req, "response", data)
}

func (luga *ListGroupPolicyApi) listGroupPolicy() {
	beans := make([]*db.PolicyGroupBean, 0)
	luga.err = db.ActiveService().ListGroupPolicy(luga.groupId, &beans)
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
		up := &GroupPolicy{
			policyName:  policy.PolicyName,
			policyType:  policy.PolicyType,
			description: policy.Description,
			version:     policy.Version,
			attachDate:  bean.AttachDate,
		}
		luga.policies = append(luga.policies, up)
	}
}

func ListGroupPolicyHandler(w http.ResponseWriter, r *http.Request) {
	luga := ListGroupPolicyApi{req: r, status: http.StatusOK}
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

	groupId, err := group.GetGroupId(luga.account, luga.group)
	if err != nil {
		luga.err = err
		return
	}
	luga.groupId = groupId

	if luga.listGroupPolicy(); luga.err != nil {
		return
	}
}
