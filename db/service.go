package db

import (
	"errors"
)

var (
	UserExistError = errors.New("user already exist")
)

type UserService interface {
	CreateIamUser(usr *UserBean) (*UserBean, error)
	UserCountOfAccount(accountId string) (int, error)
}

type KeyService interface {
	GetKey(keyId string) (*KeyBean, error)
}

type Service interface {
	UserService
	KeyService
}
