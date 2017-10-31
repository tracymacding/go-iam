package mysql

import (
	"github.com/go-iam/db"
)

func (ms *mysqlService) CreateIamUser(usr *db.UserBean) (*db.UserBean, error) {
	return nil, nil
}

func (ms *mysqlService) UserCountOfAccount(accountId string) (int, error) {
	return 0, nil
}
