package user

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/db"
	"regexp"
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

var (
	UserNameTooLongError    = errors.New("UserName beyond the length limit")
	DisplayNameTooLongError = errors.New("DisplayName beyond the length limit")
	CommentTooLongError     = errors.New("comments beyond the length limit")
	UserNameInvalidError    = errors.New("user name contains invalid char")
	DisplayNameInvalidError = errors.New("display name contains invalid char")
	MobilePhoneInvalidError = errors.New("mobile phone syntax error")
	EmailInvalidError       = errors.New("email syntax error")
)

const (
	MaxUserNameLength    = 64
	MaxDisplayNameLength = 12
	MaxCommentLength     = 128
)

func (user *User) checkUserName() error {
	if len(user.userName) > MaxUserNameLength {
		return UserNameTooLongError
	}
	match, _ := regexp.MatchString("^[a-zA-Z0-9\\.@\\-_]+$", user.userName)
	if !match {
		return UserNameInvalidError
	}
	return nil
}

func (user *User) checkDisplayName() error {
	if user.displayName == "" {
		return nil
	}

	if len(user.displayName) > MaxDisplayNameLength {
		return DisplayNameTooLongError
	}
	match, _ := regexp.MatchString("^[a-zA-Z0-9\\.@\\-\u4e00-\u9fa5]+$", user.displayName)
	if !match {
		return DisplayNameInvalidError
	}
	return nil
}

func (user *User) checkPhone() error {
	reg := `^1([38][0-9]|14[57]|5[^4])\d{8}$`
	rgx := regexp.MustCompile(reg)
	if user.phone != "" && !rgx.MatchString(user.phone) {
		return MobilePhoneInvalidError
	}
	return nil
}

func (user *User) checkEmail() error {
	if user.email == "" {
		return nil
	}

	m, _ := regexp.MatchString("^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+", user.email)
	if !m {
		return EmailInvalidError
	}
	return nil
}

func (user *User) checkComments() error {
	if len(user.comments) > MaxCommentLength {
		return CommentTooLongError
	}
	return nil
}

func (user *User) validate() error {
	err := user.checkUserName()
	if err != nil {
		return err
	}

	err = user.checkDisplayName()
	if err != nil {
		return err
	}

	err = user.checkComments()
	if err != nil {
		return err
	}

	err = user.checkEmail()
	if err != nil {
		return err
	}

	err = user.checkPhone()
	if err != nil {
		return err
	}

	return nil
}

func FromBean(bean *db.UserBean) User {
	user := User{}
	user.userName = bean.UserName
	user.userId = bean.UserId.Hex()
	user.displayName = bean.DisplayName
	user.phone = bean.Phone
	user.email = bean.Email
	user.comments = bean.Comments
	user.password = bean.Password
	user.createDate = bean.CreateDate
	return user
}

func (user *User) ToBean() db.UserBean {
	return db.UserBean{
		UserName:    user.userName,
		DisplayName: user.displayName,
		Phone:       user.phone,
		Email:       user.email,
		Comments:    user.comments,
		Password:    user.password,
		Account:     user.account,
		CreateDate:  user.createDate,
	}
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
