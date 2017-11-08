package policy

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/user"
	"github.com/go-iam/handler/util"
	"net/http"
	"time"
)

type UserAttachPolicyApi struct {
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

func (uapa *UserAttachPolicyApi) Parse() {
	params := util.ParseParameters(uapa.req)
	uapa.policy = params["PolicyName"]
	uapa.policyType = ParsePolicyType(params["PolicyType"])
	uapa.user = params["UserName"]
}

var (
	MissPolicyTypeError = errors.New("Missing parameter PolicyType")
)

func (uapa *UserAttachPolicyApi) Validate() {
	if uapa.policy == "" {
		uapa.err = MissPolicyNameError
		uapa.status = http.StatusBadRequest
		return
	}
	if !IsValidType(uapa.policyType) {
		uapa.err = MissPolicyTypeError
		uapa.status = http.StatusBadRequest
		return
	}
	if uapa.user == "" {
		uapa.err = user.MissUserNameError
		uapa.status = http.StatusBadRequest
		return
	}
	if ok, err := user.IsUserNameValid(uapa.user); !ok {
		uapa.err = err
		uapa.status = http.StatusBadRequest
		return
	}
	if ok, err := IsPolicyNameValid(uapa.policy); !ok {
		uapa.err = err
		uapa.status = http.StatusBadRequest
		return
	}
}

func (uapa *UserAttachPolicyApi) Auth() {
	uapa.err = doAuth(uapa.req)
	if uapa.err != nil {
		uapa.status = http.StatusForbidden
	}
}

func (uapa *UserAttachPolicyApi) Response() {
	json := simplejson.New()
	if uapa.err != nil {
		gerr := gerror.NewIAMError(uapa.status, uapa.err)
		context.Set(uapa.req, "request_error", gerr)
		json.Set("ErrorMessage", uapa.err.Error())
	}
	json.Set("RequestId", context.Get(uapa.req, "request_id"))
	data, _ := json.Encode()
	context.Set(uapa.req, "response", data)
}

const (
	MAX_POLICY_PER_USER_ATTACHED = 100
)

var (
	TooManyPolicyUserAttachedError = errors.New("The policy count of the user attached beyond the current limits")
)

func (uapa *UserAttachPolicyApi) attachPolicyToUser() {
	userId, err := user.GetUserId(uapa.account, uapa.user)
	if err != nil {
		uapa.err = err
		return
	}
	uapa.userId = userId

	policyId, err := GetPolicyId(uapa.account, uapa.policy, uapa.policyType)
	if err != nil {
		uapa.err = err
		return
	}
	uapa.policyId = policyId

	cnt := 0
	cnt, uapa.err = db.ActiveService().UserAttachedPolicyNum(uapa.userId)
	if uapa.err != nil {
		uapa.status = http.StatusInternalServerError
		return
	}

	if cnt >= MAX_POLICY_PER_USER_ATTACHED {
		uapa.status = http.StatusConflict
		uapa.err = TooManyPolicyUserAttachedError
		return
	}

	bean := &db.PolicyUserBean{
		PolicyId:   uapa.policyId,
		UserId:     uapa.userId,
		AttachDate: time.Now().Format(time.RFC3339),
	}
	bean, uapa.err = db.ActiveService().UserAttachPolicy(bean)
	if uapa.err != nil {
		if uapa.err == db.PolicyAttachedUserError {
			uapa.status = http.StatusConflict
		} else {
			uapa.status = http.StatusInternalServerError
		}
		return
	}
}

func UserAttachPolicyHandler(w http.ResponseWriter, r *http.Request) {
	uapa := UserAttachPolicyApi{req: r, status: http.StatusOK}
	defer uapa.Response()

	if uapa.Auth(); uapa.err != nil {
		return
	}

	if uapa.Parse(); uapa.err != nil {
		return
	}

	if uapa.Validate(); uapa.err != nil {
		return
	}

	if uapa.attachPolicyToUser(); uapa.err != nil {
		return
	}
}
