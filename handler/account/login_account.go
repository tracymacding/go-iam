package account

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
)

type LoginAccountApi struct {
	req     *http.Request
	status  int
	err     error
	account Account
}

func (laa *LoginAccountApi) Parse() {
	params := util.ParseParameters(laa.req)
	laa.account.accountId = params["AccountId"]
	laa.account.accountName = params["AccountName"]
	laa.account.password = params["Password"]
}

func (laa *LoginAccountApi) Validate() {
	if laa.account.accountId == "" {
		laa.err = MissAccountIdError
		laa.status = http.StatusBadRequest
		return
	}
	if laa.account.accountName == "" {
		laa.err = MissAccountNameError
		laa.status = http.StatusBadRequest
		return
	}
	if laa.account.password == "" {
		laa.err = MissPasswordError
		laa.status = http.StatusBadRequest
		return
	}
}

func (laa *LoginAccountApi) Auth() {
	laa.err = doAuth(laa.req)
	if laa.err != nil {
		laa.status = http.StatusForbidden
	}
}

func (laa *LoginAccountApi) Response() {
	json := simplejson.New()
	if laa.err != nil {
		json.Set("ErrorMessage", laa.err.Error())
		context.Set(laa.req, "request_error", gerror.NewIAMError(laa.status, laa.err))
	} else {
		j := laa.account.Json()
		json.Set("Account", j)
		// TODO: fill in system ak/sk pair
	}
	json.Set("RequestId", context.Get(laa.req, "request_id"))
	data, _ := json.Encode()
	context.Set(laa.req, "response", data)
}

func LoginAccountHandler(w http.ResponseWriter, r *http.Request) {
	laa := LoginAccountApi{req: r, status: http.StatusOK}

	defer laa.Response()

	if laa.Auth(); laa.err != nil {
		return
	}

	laa.Parse()

	if laa.Validate(); laa.err != nil {
		return
	}

	gaa := GetAccountApi{}
	gaa.account.accountId = laa.account.accountId

	if gaa.getAccount(); gaa.err != nil {
		laa.err = gaa.err
		return
	}

	if laa.account.accountName != gaa.account.accountName {
		laa.status = http.StatusUnauthorized
		laa.err = errors.New("Invalid UserName")
		return
	}
	if laa.account.password != gaa.account.password {
		laa.status = http.StatusUnauthorized
		laa.err = errors.New("Invalid Password")
		return
	}

	// TODO: list ak/sk pair
}
