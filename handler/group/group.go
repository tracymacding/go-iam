package group

import (
	"github.com/bitly/go-simplejson"
)

type Group struct {
	groupId    string
	groupName  string
	comments   string
	createDate string
	account    string
}

func (group *Group) Json() *simplejson.Json {
	j := simplejson.New()
	j.Set("GroupId", group.groupId)
	j.Set("GroupName", group.groupName)
	j.Set("Comments", group.comments)
	j.Set("CreateDate", group.createDate)
	return j
}

type UserGroup struct {
	groupName  string
	comments   string
	createDate string
	joinDate   string
}

func (ug *UserGroup) Json() *simplejson.Json {
	j := simplejson.New()
	j.Set("GroupName", ug.groupName)
	j.Set("Comments", ug.comments)
	j.Set("CreateDate", ug.createDate)
	j.Set("JoinDate", ug.joinDate)
	return j
}

type GroupUser struct {
	userId      string
	userName    string
	displayName string
	createDate  string
	joinDate    string
}

func (gu *GroupUser) Json() *simplejson.Json {
	j := simplejson.New()
	j.Set("UserId", gu.userId)
	j.Set("UserName", gu.userName)
	j.Set("DisplayName", gu.displayName)
	j.Set("CreateDate", gu.createDate)
	j.Set("JoinDate", gu.joinDate)
	return j
}
