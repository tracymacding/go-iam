package account

import (
	"github.com/bitly/go-simplejson"
)

type Account struct {
	accountId   string
	accountName string
	password    string
	accountType AccountType
	createDate  string
}

func (acc *Account) Json() *simplejson.Json {
	j := simplejson.New()
	j.Set("AccountId", acc.accountId)
	j.Set("AccountName", acc.accountName)
	j.Set("AccountType", acc.accountType.String())
	j.Set("CreateDate", acc.createDate)
	return j
}
