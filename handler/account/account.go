package account

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/db"
)

type Account struct {
	accountId   string
	accountName string
	password    string
	accountType AccountType
	createDate  string
}

func FromBean(bean *db.AccountBean) Account {
	account := Account{}
	account.accountId = bean.AccountId.Hex()
	account.accountName = bean.AccountName
	account.accountType = AccountType(bean.AccountType)
	account.password = bean.Password
	account.createDate = bean.CreateDate
	return account
}

func (acc *Account) ToBean() db.AccountBean {
	bean := db.AccountBean{
		AccountName: acc.accountName,
		Password:    acc.password,
		AccountType: int(acc.accountType),
		CreateDate:  acc.createDate,
	}
	return bean
}

func (acc *Account) Json() *simplejson.Json {
	j := simplejson.New()
	j.Set("AccountId", acc.accountId)
	j.Set("AccountName", acc.accountName)
	j.Set("AccountType", acc.accountType.String())
	j.Set("CreateDate", acc.createDate)
	return j
}
