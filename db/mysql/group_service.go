package mysql

import (
	"github.com/go-iam/db"
)

func (ms *mysqlService) CreateGroup(group *db.GroupBean) (*db.GroupBean, error) {
	return nil, nil
}

func (ms *mysqlService) GroupCountOfAccount(accountId string) (int, error) {
	return 0, nil
}

func (ms *mysqlService) DeleteGroup(account, group string) error {
	return nil
}

func (ms *mysqlService) GetGroup(account, group string, grp *db.GroupBean) error {
	return nil
}

func (ms *mysqlService) UpdateGroup(account, group string, grp *db.GroupBean) error {
	return nil
}

func (ms *mysqlService) ListGroup(account, marker string, max int, groups *[]*db.GroupBean) error {
	return nil
}
