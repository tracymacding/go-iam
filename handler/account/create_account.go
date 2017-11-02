package account

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
	"time"
)

type ParamValidater interface {
	Validate()
}

type ParamParser interface {
	Parse()
}

type Auther interface {
	Auth()
}

type Responser interface {
	Response()
}

type CreateAccountApi struct {
	req     *http.Request
	status  int
	err     error
	account Account
}

var (
	MissAccountNameError    = errors.New("account name missing")
	MissPasswordError       = errors.New("password missing")
	MissAccountTypeError    = errors.New("account type missing")
	InvalidAccountTypeError = errors.New("invalid account type")
)

func (caa *CreateAccountApi) Validate() {
	if caa.account.accountName == "" {
		caa.err = MissAccountNameError
		caa.status = http.StatusBadRequest
		return
	}
	if caa.account.password == "" {
		caa.err = MissPasswordError
		caa.status = http.StatusBadRequest
		return
	}
	if caa.account.accountType == 0 {
		caa.err = MissAccountTypeError
		caa.status = http.StatusBadRequest
		return
	}
	if !IsValidType(caa.account.accountType) {
		caa.err = InvalidAccountTypeError
		caa.status = http.StatusBadRequest
		return
	}
}

func (caa *CreateAccountApi) Parse() {
	params := util.ParseParameters(caa.req)
	caa.account.accountName = params["AccountName"]
	caa.account.password = params["Password"]
	caa.account.accountType = ParseAccountType(params["AccountType"])
}

func (caa *CreateAccountApi) Auth() {
	caa.err = doAuth(caa.req)
	if caa.err != nil {
		caa.status = http.StatusForbidden
	}
}

func (caa *CreateAccountApi) Response() {
	json := simplejson.New()
	if caa.err == nil {
		j := caa.account.Json()
		json.Set("Account", j)
	} else {
		context.Set(caa.req, "request_error", gerror.NewIAMError(caa.status, caa.err))
		json.Set("ErrorMessage", caa.err.Error())
	}
	json.Set("RequestId", context.Get(caa.req, "request_id"))
	data, _ := json.Encode()
	context.Set(caa.req, "response", data)
}

func (caa *CreateAccountApi) createAccount() {
	bean := &db.AccountBean{
		AccountName: caa.account.accountName,
		Password:    caa.account.password,
		AccountType: int(caa.account.accountType),
		CreateDate:  time.Now().Format(time.RFC3339),
	}
	bean, caa.err = db.ActiveService().CreateAccount(bean)
	if caa.err != nil {
		if caa.err == db.AccountExistError {
			caa.status = http.StatusConflict
		} else {
			caa.status = http.StatusInternalServerError
		}
	} else {
		caa.account.accountId = bean.AccountId.Hex()
		caa.account.createDate = bean.CreateDate
	}
}

func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	caa := CreateAccountApi{req: r, status: http.StatusOK}

	defer caa.Response()

	if caa.Auth(); caa.err != nil {
		return
	}

	caa.Parse()

	if caa.Validate(); caa.err != nil {
		return
	}

	if caa.createAccount(); caa.err != nil {
		return
	}

	// TODO: create system type ak/sk pair for account
}
