package account

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
)

type ListAccountApi struct {
	req      *http.Request
	status   int
	err      error
	saccType string
	accType  AccountType
	accounts []*Account
}

func (lsaa *ListAccountApi) Parse() {
	params := util.ParseParameters(lsaa.req)
	lsaa.saccType = params["AccountType"]
	lsaa.accType = ParseAccountType(lsaa.saccType)
}

func (lsaa *ListAccountApi) Validate() {
	if lsaa.saccType != "" && !IsValidType(lsaa.accType) {
		lsaa.err = InvalidAccountTypeError
		lsaa.status = http.StatusBadRequest
		return
	}
}

func (lsaa *ListAccountApi) Auth() {
	lsaa.err = doAuth(lsaa.req)
	if lsaa.err != nil {
		lsaa.status = http.StatusForbidden
	}
}

func (lsaa *ListAccountApi) Response() {
	json := simplejson.New()
	if lsaa.err == nil {
		jsons := make([]*simplejson.Json, 0)
		for _, account := range lsaa.accounts {
			j := account.Json()
			jsons = append(jsons, j)
		}
		json.Set("Account", jsons)
	} else {
		json.Set("ErrorMessage", lsaa.err.Error())
		gerr := gerror.NewIAMError(lsaa.status, lsaa.err)
		context.Set(lsaa.req, "request_error", gerr)
	}
	json.Set("RequestId", context.Get(lsaa.req, "request_id"))
	data, _ := json.Encode()
	context.Set(lsaa.req, "response", data)
}

func (lsaa *ListAccountApi) listAccount() {
	beans := make([]*db.AccountBean, 0)
	if lsaa.accType == 0 {
		lsaa.err = db.ActiveService().ListAccount(0, &beans)
	} else {
		lsaa.err = db.ActiveService().ListAccount(int(lsaa.accType), &beans)
	}
	if lsaa.err != nil {
		lsaa.status = http.StatusInternalServerError
		return
	}

	for _, bean := range beans {
		account := FromBean(bean)
		lsaa.accounts = append(lsaa.accounts, &account)
	}
}

func ListAccountHandler(w http.ResponseWriter, r *http.Request) {
	lsaa := ListAccountApi{req: r, status: http.StatusOK}
	defer lsaa.Response()

	if lsaa.Auth(); lsaa.err != nil {
		return
	}

	if lsaa.Parse(); lsaa.err != nil {
		return
	}

	if lsaa.Validate(); lsaa.err != nil {
		return
	}

	if lsaa.listAccount(); lsaa.err != nil {
		return
	}
}
