package user

import (
	"github.com/bitly/go-simplejson"
)

type User struct {
	userId      string
	userName    string
	displayName string
	phone       string
	email       string
	comments    string
	password    string
	createDate  string
	account     string
}

func (user *User) Json() *simplejson.Json {
	j := simplejson.New()
	j.Set("UserId", user.userId)
	j.Set("UserName", user.userName)
	j.Set("DisplayName", user.displayName)
	j.Set("MobilePhone", user.phone)
	j.Set("Email", user.email)
	j.Set("Comments", user.comments)
	j.Set("CreateDate", user.createDate)
	return j
}
