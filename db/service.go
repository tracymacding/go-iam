package db

import (
	"errors"
)

var (
	UserExistError              = errors.New("user already exist")
	UserNotExistError           = errors.New("user not exist")
	AccountExistError           = errors.New("account already exist")
	AccountNotExistError        = errors.New("account not exist")
	GroupExistError             = errors.New("group already exist")
	GroupNotExistError          = errors.New("group not exist")
	UserJoinedGroupError        = errors.New("the user has already joined the group")
	UserNotJoinedGroupError     = errors.New("the user not joined the group")
	PolicyExistError            = errors.New("policy already exist")
	PolicyNotExistError         = errors.New("policy not exist")
	PolicyAttachedUserError     = errors.New("the policy has already attached to the user")
	PolicyNotAttachedUserError  = errors.New("the policy not attached to the group")
	PolicyAttachedGroupError    = errors.New("the policy has already attached to the group")
	PolicyNotAttachedGroupError = errors.New("the policy not attached to the user")
	KeyExistError               = errors.New("access key already exist")
	KeyNotExistError            = errors.New("access key not exist")
)

type AccountService interface {
	CreateAccount(account *AccountBean) error
	GetAccount(accountId string, account *AccountBean) error
	DeleteAccount(accountId string) error
	UpdateAccount(accountId string, account *AccountBean) error
	ListAccount(accountType int, accounts *[]*AccountBean) error
}

type UserService interface {
	CreateIamUser(usr *UserBean) error
	GetIamUser(account, user string, usr *UserBean) error
	GetIamUserById(userId string, usr *UserBean) error
	DeleteIamUser(account, user string) error
	UpdateIamUser(usr *UserBean) error
	ListIamUser(account, marker string, max int, usrs *[]*UserBean) error
	UserCountOfAccount(accountId string) (int, error)
}

type GroupService interface {
	CreateGroup(group *GroupBean) error
	GetGroup(account, group string, grp *GroupBean) error
	GetGroupById(groupId string, grp *GroupBean) error
	DeleteGroup(account, group string) error
	UpdateGroup(grp *GroupBean) error
	ListGroup(account, marker string, max int, groups *[]*GroupBean) error
	GroupCountOfAccount(accountId string) (int, error)
}

type GroupUserService interface {
	GroupAddUser(bean *GroupUserBean) (*GroupUserBean, error)
	GroupRemoveUser(bean *GroupUserBean) error
	ListUserGroup(userId string, beans *[]*GroupUserBean) error
	ListGroupUser(group string, beans *[]*GroupUserBean) error
	UserJoinedGroupsNum(userId string) (int, error)
}

type PolicyUserService interface {
	UserAttachPolicy(bean *PolicyUserBean) (*PolicyUserBean, error)
	UserDetachPolicy(bean *PolicyUserBean) error
	ListUserPolicy(userId string, beans *[]*PolicyUserBean) error
	ListPolicyUser(policyId string, beans *[]*PolicyUserBean) error
	UserAttachedPolicyNum(userId string) (int, error)
}

type PolicyGroupService interface {
	GroupAttachPolicy(bean *PolicyGroupBean) (*PolicyGroupBean, error)
	GroupDetachPolicy(bean *PolicyGroupBean) error
	ListGroupPolicy(groupId string, beans *[]*PolicyGroupBean) error
	ListPolicyGroup(policyId string, beans *[]*PolicyGroupBean) error
	GroupAttachedPolicyNum(groupId string) (int, error)
}

type PolicyService interface {
	CreatePolicy(policy *PolicyBean) error
	GetPolicy(account, policy string, bean *PolicyBean) error
	GetPolicyById(policyId string, policy *PolicyBean) error
	DeletePolicy(account, policy string) error
	UpdatePolicy(policyId string, bean *PolicyBean) error
	ListPolicy(account, marker string, ptype, max int, groups *[]*PolicyBean) error
	ListAllPolicy(account, marker string, max int, groups *[]*PolicyBean) error
	PolicyCountOfAccount(accountId string) (int, error)
}

type KeyService interface {
	CreateKey(key *KeyBean) error
	GetKey(id string, key *KeyBean) error
	DeleteKey(accessKeyId string) error
	UpdateKey(id string, key *KeyBean) error
	ListKey(entity string, entitype int, keys *[]*KeyBean) error
	KeyCountOfEntity(entity string, entitype int) (int, error)
}

type Service interface {
	AccountService
	UserService
	GroupService
	GroupUserService
	PolicyService
	PolicyUserService
	PolicyGroupService
	KeyService
}
