package mysql

import (
	"github.com/go-iam/db"
)

// TODO
func (ms *mysqlService) GetKey(keyId string, key *db.KeyBean) error {
	return nil
}

func (ms *mysqlService) CreateKey(key *db.KeyBean) (*db.KeyBean, error) {
	return nil, nil
}

func (ms *mysqlService) DeleteKey(key string) error {
	return nil
}

func (ms *mysqlService) UpdateKey(id string, key *db.KeyBean) error {
	return nil
}

func (ms *mysqlService) ListKey(entity string, entitype int, keys *[]*db.KeyBean) error {
	return nil
}

func (ms *mysqlService) KeyCountOfEntity(entity string, entitype int) (int, error) {
	return 0, nil
}
