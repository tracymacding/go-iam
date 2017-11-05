package mysql

import (
	"github.com/go-iam/db"
)

func (ms *mysqlService) GroupAddUser(bean *db.GroupUserBean) (*db.GroupUserBean, error) {
	return nil, nil
}

func (ms *mysqlService) GroupRemoveUser(bean *db.GroupUserBean) error {
	return nil
}

func (ms *mysqlService) ListUserGroup(userId string, beans *[]*db.GroupUserBean) error {
	return nil
}

func (ms *mysqlService) ListGroupUser(groupId string, beans *[]*db.GroupUserBean) error {
	return nil
}

func (ms *mysqlService) UserJoinedGroupsNum(userId string) (int, error) {
	return 0, nil
}
