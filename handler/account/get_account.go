package account

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
)

type GetAccountApi struct {
	req     *http.Request
	status  int
	err     error
	account Account
}

var (
	MissAccountIdError = errors.New("account id missing")
)

func (gaa *GetAccountApi) Parse() {
	params := util.ParseParameters(gaa.req)
	gaa.account.accountId = params["AccountId"]
}

func (gaa *GetAccountApi) Validate() {
	if gaa.account.accountId == "" {
		gaa.err = MissAccountIdError
		gaa.status = http.StatusBadRequest
		return
	}
}

func (gaa *GetAccountApi) Auth() {
	gaa.err = doAuth(gaa.req)
	if gaa.err != nil {
		gaa.status = http.StatusForbidden
	}
}

func (gaa *GetAccountApi) Response() {
	json := simplejson.New()
	if gaa.err == nil {
		j := gaa.account.Json()
		json.Set("Account", j)
	} else {
		gerr := gerror.NewIAMError(gaa.status, gaa.err)
		context.Set(gaa.req, "request_error", gerr)
		json.Set("ErrorMessage", gaa.err.Error())
	}
	json.Set("RequestId", context.Get(gaa.req, "request_id"))
	data, _ := json.Encode()
	context.Set(gaa.req, "response", data)
}

func (gaa *GetAccountApi) getAccount() {
	var bean db.AccountBean
	gaa.err = db.ActiveService().GetAccount(gaa.account.accountId, &bean)
	if gaa.err != nil {
		if gaa.err == db.AccountNotExistError {
			gaa.status = http.StatusNotFound
		} else {
			gaa.status = http.StatusInternalServerError
		}
		return
	}
	gaa.account = FromBean(&bean)
}

func GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	gaa := GetAccountApi{req: r, status: http.StatusOK}
	defer gaa.Response()

	if gaa.Auth(); gaa.err != nil {
		return
	}

	if gaa.Parse(); gaa.err != nil {
		return
	}

	if gaa.Validate(); gaa.err != nil {
		return
	}

	if gaa.getAccount(); gaa.err != nil {
		return
	}
}
