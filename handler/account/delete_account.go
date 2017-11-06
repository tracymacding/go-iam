package account

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
)

type DeleteAccountApi struct {
	req     *http.Request
	status  int
	err     error
	account Account
}

func (daa *DeleteAccountApi) Parse() {
	params := util.ParseParameters(daa.req)
	daa.account.accountId = params["AccountId"]
}

func (daa *DeleteAccountApi) Validate() {
	if daa.account.accountId == "" {
		daa.err = MissAccountIdError
		daa.status = http.StatusBadRequest
		return
	}
}

func (daa *DeleteAccountApi) Auth() {
	daa.err = doAuth(daa.req)
	if daa.err != nil {
		daa.status = http.StatusForbidden
	}
}

func (daa *DeleteAccountApi) Response() {
	json := simplejson.New()
	if daa.err != nil {
		gerr := gerror.NewIAMError(daa.status, daa.err)
		context.Set(daa.req, "request_error", gerr)
		json.Set("ErrorMessage", daa.err.Error())
	}
	json.Set("RequestId", context.Get(daa.req, "request_id"))
	data, _ := json.Encode()
	context.Set(daa.req, "response", data)
}

func (daa *DeleteAccountApi) deleteAccount() {
	daa.err = db.ActiveService().DeleteAccount(daa.account.accountId)
	if daa.err != nil {
		if daa.err == db.AccountNotExistError {
			daa.status = http.StatusNotFound
		} else {
			daa.status = http.StatusInternalServerError
		}
	}
}

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	daa := DeleteAccountApi{req: r, status: http.StatusOK}
	defer daa.Response()

	if daa.Auth(); daa.err != nil {
		return
	}

	if daa.Parse(); daa.err != nil {
		return
	}

	if daa.Validate(); daa.err != nil {
		return
	}

	if daa.deleteAccount(); daa.err != nil {
		return
	}
}
