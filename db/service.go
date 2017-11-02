package db

import (
	"errors"
)

var (
	UserExistError       = errors.New("user already exist")
	UserNotExistError    = errors.New("user not exist")
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
	GetIamUser(account, user string, usr *UserBean) error
	DeleteIamUser(account, user string) error
	UpdateIamUser(account, user string, usr *UserBean) error
	ListIamUser(account, marker string, max int, usrs *[]*UserBean) error
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
