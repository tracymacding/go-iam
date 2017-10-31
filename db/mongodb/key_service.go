package mongodb

import (
	"github.com/go-iam/db"
)

// TODO
func (ms *mongoService) GetKey(keyId string) (*db.KeyBean, error) {
	return &db.KeyBean{KeyId: "fake_ak_id", CreatorType: 1}, nil
}
