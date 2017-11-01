package account

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"net/http"
)

type ListAccountApi struct {
	req      *http.Request
	status   int
	err      error
	accType  AccountType
	accounts []*Account
}

func (lsaa *ListAccountApi) Parse() {
	params := parseParameters(lsaa.req)
	lsaa.accType = ParseAccountType(params["AccountType"])
}

func (lsaa *ListAccountApi) Validate() {
	if !IsValidType(lsaa.accType) {
		lsaa.err = MissAccountIdError
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
			accountJson := simplejson.New()
			accountJson.Set("AccountId", account.accountId)
			accountJson.Set("AccountName", account.accountName)
			accountJson.Set("AccountType", account.accountType.String())
			accountJson.Set("CreateDate", account.createDate)
			jsons = append(jsons, accountJson)
		}
		json.Set("Account", jsons)
	} else {
		json.Set("ErrorMessage", lsaa.err.Error())
		context.Set(lsaa.req, "request_error", gerror.NewIAMError(lsaa.status, lsaa.err))
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
		account := &Account{
			accountId:   bean.AccountId.Hex(),
			accountName: bean.AccountName,
			accountType: AccountType(bean.AccountType),
			createDate:  bean.CreateDate,
		}
		lsaa.accounts = append(lsaa.accounts, account)
	}
}

func ListAccountHandler(w http.ResponseWriter, r *http.Request) {
	lsaa := ListAccountApi{req: r, status: http.StatusOK}

	defer lsaa.Response()

	if lsaa.Auth(); lsaa.err != nil {
		return
	}

	lsaa.Parse()

	if lsaa.Validate(); lsaa.err != nil {
		return
	}

	if lsaa.listAccount(); lsaa.err != nil {
		return
	}
}
