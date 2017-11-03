package policy

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
	"strconv"
)

type ListPolicyApi struct {
	req      *http.Request
	status   int
	err      error
	policies []*Policy
	marker   string
	maxItems int
	ptype    PolicyType
	account  string
}

func (lpa *ListPolicyApi) Parse() {
	params := util.ParseParameters(lpa.req)
	ptype := params["PolicyType"]
	lpa.marker = params["Marker"]
	items := params["MaxItems"]

	lpa.ptype = ParsePolicyType(ptype)
	if items == "" {
		lpa.maxItems = 100
	}
	lpa.maxItems, lpa.err = strconv.Atoi(items)
}

var (
	InvalidMaxItemsError   = errors.New("MaxItems parameter out of range")
	InvalidPolicyTypeError = errors.New("Invalid policy type")
)

func (lpa *ListPolicyApi) Validate() {
	if lpa.maxItems < 1 || lpa.maxItems > 1000 {
		lpa.err = InvalidMaxItemsError
		lpa.status = http.StatusBadRequest
		return
	}
	if !IsValidType(lpa.ptype) {
		lpa.err = InvalidPolicyTypeError
		lpa.status = http.StatusBadRequest
		return
	}
}

func (lpa *ListPolicyApi) Auth() {
	lpa.err = doAuth(lpa.req)
	if lpa.err != nil {
		lpa.status = http.StatusForbidden
	}
}

func (lpa *ListPolicyApi) Response() {
	json := simplejson.New()
	if lpa.err == nil {
		jsons := make([]*simplejson.Json, 0)
		for _, policy := range lpa.policies {
			j := policy.Json()
			jsons = append(jsons, j)
		}
		json.Set("Policies", jsons)
	} else {
		json.Set("ErrorMessage", lpa.err.Error())
		context.Set(lpa.req, "request_error", gerror.NewIAMError(lpa.status, lpa.err))
	}
	json.Set("RequestId", context.Get(lpa.req, "request_id"))
	data, _ := json.Encode()
	context.Set(lpa.req, "response", data)
}

func (lpa *ListPolicyApi) listPolicy() {
	beans := make([]*db.PolicyBean, 0)

	lpa.err = db.ActiveService().ListPolicy(
		lpa.account,
		lpa.marker,
		int(lpa.ptype),
		lpa.maxItems,
		&beans)
	if lpa.err != nil {
		lpa.status = http.StatusInternalServerError
		return
	}

	for _, bean := range beans {
		policy := &Policy{
			policyId:    bean.PolicyId.Hex(),
			policyName:  bean.PolicyName,
			policyType:  bean.PolicyType,
			document:    bean.Document,
			description: bean.Description,
			version:     bean.Version,
			createDate:  bean.CreateDate,
			updateDate:  bean.UpdateDate,
		}
		lpa.policies = append(lpa.policies, policy)
	}
}

func ListPolicyHandler(w http.ResponseWriter, r *http.Request) {
	lpa := ListPolicyApi{req: r, status: http.StatusOK}

	defer lpa.Response()

	if lpa.Auth(); lpa.err != nil {
		return
	}

	if lpa.Parse(); lpa.err != nil {
		return
	}

	if lpa.Validate(); lpa.err != nil {
		return
	}

	if lpa.listPolicy(); lpa.err != nil {
		return
	}
}
