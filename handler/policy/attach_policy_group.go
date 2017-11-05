package policy

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/group"
	"github.com/go-iam/handler/util"
	"net/http"
	"time"
)

type GroupAttachPolicyApi struct {
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

func (gapa *GroupAttachPolicyApi) Parse() {
	params := util.ParseParameters(gapa.req)
	gapa.policy = params["PolicyName"]
	gapa.policyType = ParsePolicyType(params["PolicyType"])
	gapa.group = params["GroupName"]
}

func (gapa *GroupAttachPolicyApi) Validate() {
	if gapa.policy == "" {
		gapa.err = MissPolicyNameError
		gapa.status = http.StatusBadRequest
		return
	}
	if !IsValidType(gapa.policyType) {
		gapa.err = MissPolicyTypeError
		gapa.status = http.StatusBadRequest
		return
	}
	if gapa.group == "" {
		gapa.err = group.MissGroupNameError
		gapa.status = http.StatusBadRequest
		return
	}
}

func (gapa *GroupAttachPolicyApi) Auth() {
	gapa.err = doAuth(gapa.req)
	if gapa.err != nil {
		gapa.status = http.StatusForbidden
	}
}

func (gapa *GroupAttachPolicyApi) Response() {
	json := simplejson.New()
	if gapa.err != nil {
		context.Set(gapa.req, "request_error", gerror.NewIAMError(gapa.status, gapa.err))
		json.Set("ErrorMessage", gapa.err.Error())
	}
	json.Set("RequestId", context.Get(gapa.req, "request_id"))
	data, _ := json.Encode()
	context.Set(gapa.req, "response", data)
}

const (
	MAX_POLICY_PER_GROUP_ATTACHED = 100
)

var (
	TooManyPolicyGroupAttachedError = errors.New("The policy count of the group attached beyond the current limits")
)

func (gapa *GroupAttachPolicyApi) attachPolicyToGroup() {
	cnt := 0
	cnt, gapa.err = db.ActiveService().GroupAttachedPolicyNum(gapa.groupId)
	if gapa.err != nil {
		gapa.status = http.StatusInternalServerError
		return
	}

	if cnt >= MAX_POLICY_PER_GROUP_ATTACHED {
		gapa.status = http.StatusConflict
		gapa.err = TooManyPolicyGroupAttachedError
		return
	}

	bean := &db.PolicyGroupBean{
		PolicyId:   gapa.policyId,
		GroupId:    gapa.groupId,
		AttachDate: time.Now().Format(time.RFC3339),
	}
	bean, gapa.err = db.ActiveService().GroupAttachPolicy(bean)
	if gapa.err != nil {
		if gapa.err == db.PolicyAttachedGroupError {
			gapa.status = http.StatusConflict
		} else {
			gapa.status = http.StatusInternalServerError
		}
		return
	}
}

func GroupAttachPolicyHandler(w http.ResponseWriter, r *http.Request) {
	gapa := GroupAttachPolicyApi{req: r, status: http.StatusOK}
	defer gapa.Response()

	if gapa.Auth(); gapa.err != nil {
		return
	}

	if gapa.Parse(); gapa.err != nil {
		return
	}

	if gapa.Validate(); gapa.err != nil {
		return
	}

	groupId, err := group.GetGroupId(gapa.account, gapa.group)
	if err != nil {
		gapa.err = err
		return
	}
	gapa.groupId = groupId

	policyId, err := GetPolicyId(gapa.account, gapa.policy, gapa.policyType)
	if err != nil {
		gapa.err = err
		return
	}
	gapa.policyId = policyId

	if gapa.attachPolicyToGroup(); gapa.err != nil {
		return
	}
}
