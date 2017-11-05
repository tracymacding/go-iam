package mysql

import (
	"github.com/go-iam/db"
)

func (ms *mysqlService) CreatePolicy(policy *db.PolicyBean) (*db.PolicyBean, error) {
	return nil, nil
}

func (ms *mysqlService) PolicyCountOfAccount(accountId string) (int, error) {
	return 0, nil
}

func (ms *mysqlService) DeletePolicy(account, policy string) error {
	return nil
}

func (ms *mysqlService) GetPolicy(account, policy string, bean *db.PolicyBean) error {
	return nil
}

func (ms *mysqlService) GetPolicyById(policyId string, bean *db.PolicyBean) error {
	return nil
}

func (ms *mysqlService) UpdatePolicy(account, policy string, bean *db.PolicyBean) error {
	return nil
}

func (ms *mysqlService) ListPolicy(account, marker string, ptype, max int, policys *[]*db.PolicyBean) error {
	return nil
}
