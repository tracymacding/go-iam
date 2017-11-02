package db

import (
	"errors"
)

var (
	UserExistError       = errors.New("user already exist")
	UserNotExistError    = errors.New("user not exist")
	AccountExistError    = errors.New("account already exist")
	AccountNotExistError = errors.New("account not exist")
	GroupExistError      = errors.New("group already exist")
	GroupNotExistError   = errors.New("group not exist")
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

type GroupService interface {
	CreateGroup(group *GroupBean) (*GroupBean, error)
	GetGroup(account, group string, grp *GroupBean) error
	DeleteGroup(account, group string) error
	UpdateGroup(account, group string, grp *GroupBean) error
	ListGroup(account, marker string, max int, groups *[]*GroupBean) error
	GroupCountOfAccount(accountId string) (int, error)
}

type KeyService interface {
	GetKey(keyId string) (*KeyBean, error)
}

type Service interface {
	AccountService
	UserService
	GroupService
	KeyService
}
