package mysql

import (
	"github.com/go-iam/db"
)

func (ms *mysqlService) UserAttachPolicy(bean *db.PolicyUserBean) (*db.PolicyUserBean, error) {
	return bean, nil
}

func (ms *mysqlService) UserAttachedPolicyNum(userId string) (int, error) {
	return 0, nil
}

func (ms *mysqlService) UserDetachPolicy(bean *db.PolicyUserBean) error {
	return nil
}

func (ms *mysqlService) ListUserPolicy(userId string, beans *[]*db.PolicyUserBean) error {
	return nil
}

func (ms *mysqlService) ListPolicyUser(policyId string, beans *[]*db.PolicyUserBean) error {
	return nil
}
