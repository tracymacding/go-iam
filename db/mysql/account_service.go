package mysql

import (
	"github.com/go-iam/db"
)

func (ms *mysqlService) CreateAccount(account *db.AccountBean) error {
	return nil
}

func (ms *mysqlService) GetAccount(accountId string, account *db.AccountBean) error {
	return nil
}

func (ms *mysqlService) DeleteAccount(accountId string) error {
	return nil
}

func (ms *mysqlService) UpdateAccount(accountId string, account *db.AccountBean) error {
	return nil
}

func (ms *mysqlService) ListAccount(accountType int, accounts *[]*db.AccountBean) error {
	return nil
}
