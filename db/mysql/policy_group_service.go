package mysql

import (
	"github.com/go-iam/db"
)

func (ms *mysqlService) GroupAttachPolicy(bean *db.PolicyGroupBean) (*db.PolicyGroupBean, error) {
	return bean, nil
}

func (ms *mysqlService) GroupAttachedPolicyNum(groupId string) (int, error) {
	return 0, nil
}

func (ms *mysqlService) GroupDetachPolicy(bean *db.PolicyGroupBean) error {
	return nil
}

func (ms *mysqlService) ListGroupPolicy(groupId string, beans *[]*db.PolicyGroupBean) error {
	return nil
}

func (ms *mysqlService) ListPolicyGroup(policyId string, beans *[]*db.PolicyGroupBean) error {
	return nil
}
