package mysql

import (
	"github.com/go-iam/db"
)

func (ms *mysqlService) CreateIamUser(usr *db.UserBean) error {
	return nil
}

func (ms *mysqlService) UserCountOfAccount(accountId string) (int, error) {
	return 0, nil
}

func (ms *mysqlService) GetIamUser(account, user string, usr *db.UserBean) error {
	return nil
}

func (ms *mysqlService) GetIamUserById(userId string, usr *db.UserBean) error {
	return nil
}

func (ms *mysqlService) DeleteIamUser(account, user string) error {
	return nil
}

func (ms *mysqlService) UpdateIamUser(account, user string, usr *db.UserBean) error {
	return nil
}

func (ms *mysqlService) ListIamUser(account, marker string, max int, usrs *[]*db.UserBean) error {
	return nil
}
