package group

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/db"
	"regexp"
)

type Group struct {
	groupId    string
	groupName  string
	comments   string
	createDate string
	account    string
}

func FromBean(bean *db.GroupBean) Group {
	return Group{
		groupId:    bean.GroupId.Hex(),
		groupName:  bean.GroupName,
		comments:   bean.Comments,
		createDate: bean.CreateDate,
	}
}

func (group *Group) ToBean() db.GroupBean {
	return db.GroupBean{
		GroupName:  group.groupName,
		Comments:   group.comments,
		CreateDate: group.createDate,
	}
}

func (group *Group) Json() *simplejson.Json {
	j := simplejson.New()
	j.Set("GroupId", group.groupId)
	j.Set("GroupName", group.groupName)
	j.Set("Comments", group.comments)
	j.Set("CreateDate", group.createDate)
	return j
}

var (
	GroupNameTooLongError = errors.New("group name beyond the length limit")
	CommentTooLongError   = errors.New("comments beyond the length limit")
	GroupNameInvalidError = errors.New("group name contains invalid char")
)

const (
	MaxGroupNameLength = 64
	MaxCommentLength   = 128
)

func IsGroupNameValid(groupName string) (bool, error) {
	if len(groupName) > MaxGroupNameLength {
		return false, GroupNameTooLongError
	}

	reg := `^[a-zA-Z0-9\\.@\\-_]+$`
	rgx := regexp.MustCompile(reg)
	if !rgx.MatchString(groupName) {
		return false, GroupNameInvalidError
	}

	return true, nil
}

func IsCommentsValid(comments string) (bool, error) {
	if len(comments) > MaxCommentLength {
		return false, CommentTooLongError
	}

	return true, nil
}

func (group *Group) validate() error {
	ok, err := IsGroupNameValid(group.groupName)
	if !ok {
		return err
	}

	ok, err = IsCommentsValid(group.comments)
	if !ok {
		return err
	}

	return nil
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
