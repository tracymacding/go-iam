package account

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"net/http"
)

type UpdateAccountApi struct {
	req     *http.Request
	status  int
	err     error
	account Account
}

func parseParameters(r *http.Request) map[string]string {
	params := make(map[string]string, 0)
	vals := r.URL.Query()
	for k, _ := range vals {
		params[k] = vals.Get(k)
	}
	return params
}

func (uaa *UpdateAccountApi) Parse() {
	params := parseParameters(uaa.req)
	uaa.account.accountId = params["AccountId"]
	if params["NewAccountName"] != "" {
		uaa.account.accountName = params["NewAccountName"]
	}
	if params["NewPassword"] != "" {
		uaa.account.password = params["NewPassword"]
	}
}

func (uaa *UpdateAccountApi) Validate() {
	if uaa.account.accountId == "" {
		uaa.err = MissAccountIdError
		uaa.status = http.StatusBadRequest
		return
	}
}

func (uaa *UpdateAccountApi) Auth() {
	uaa.err = doAuth(uaa.req)
	if uaa.err != nil {
		uaa.status = http.StatusForbidden
	}
}

func (uaa *UpdateAccountApi) updateAccount() {
	bean := db.AccountBean{
		AccountName: uaa.account.accountName,
		Password:    uaa.account.password,
		AccountType: int(uaa.account.accountType),
		CreateDate:  uaa.account.createDate,
	}
	uaa.err = db.ActiveService().UpdateAccount(uaa.account.accountId, &bean)
	if uaa.err == db.AccountNotExistError {
		uaa.status = http.StatusNotFound
	} else if uaa.err == db.AccountExistError {
		uaa.status = http.StatusConflict
	} else {
		uaa.status = http.StatusInternalServerError
	}
}

func (uaa *UpdateAccountApi) Response() {
	json := simplejson.New()
	if uaa.err == nil {
		accJson := simplejson.New()
		accJson.Set("AccountId", uaa.account.accountId)
		accJson.Set("AccountName", uaa.account.accountName)
		json.Set("Account", accJson)
	} else {
		json.Set("ErrorMessage", uaa.err.Error())
		context.Set(uaa.req, "request_error", gerror.NewIAMError(uaa.status, uaa.err))
	}
	json.Set("RequestId", context.Get(uaa.req, "request_id"))
	data, _ := json.Encode()
	context.Set(uaa.req, "response", data)
}

func UpdateAccountHandler(w http.ResponseWriter, r *http.Request) {
	uaa := UpdateAccountApi{req: r, status: http.StatusOK}

	defer uaa.Response()

	if uaa.Auth(); uaa.err != nil {
		return
	}

	uaa.Parse()

	if uaa.Validate(); uaa.err != nil {
		return
	}

	gaa := GetAccountApi{}
	gaa.account.accountId = uaa.account.accountId

	if gaa.getAccount(); gaa.err != nil {
		uaa.err = gaa.err
		return
	}

	uaa.account.accountType = gaa.account.accountType
	uaa.account.createDate = gaa.account.createDate
	if uaa.account.accountName == "" {
		uaa.account.accountName = gaa.account.accountName
	}
	if uaa.account.password == "" {
		uaa.account.password = gaa.account.password
	}

	uaa.updateAccount()
}
