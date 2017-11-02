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
