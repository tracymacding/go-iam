package db

import (
	"errors"
)

var (
	UserExistError       = errors.New("user already exist")
	AccountExistError    = errors.New("account already exist")
	AccountNotExistError = errors.New("account not exist")
)

type AccountService interface {
	CreateAccount(account *AccountBean) (*AccountBean, error)
	GetAccount(accountId string, account *AccountBean) error
	DeleteAccount(accountId string) error
	UpdateAccount(accountId string, account *AccountBean) error
	ListAccount(accountType int, accounts *[]*AccountBean) error
}

type UserService interface {
	CreateIamUser(usr *UserBean) (*UserBean, error)
	UserCountOfAccount(accountId string) (int, error)
}

type KeyService interface {
	GetKey(keyId string) (*KeyBean, error)
}

type Service interface {
	AccountService
	UserService
	KeyService
}
