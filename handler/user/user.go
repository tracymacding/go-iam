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
	UserNameTooLongError    = errors.New("user name beyond the length limit")
	DisplayNameTooLongError = errors.New("display name beyond the length limit")
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

func IsUserNameValid(userName string) (bool, error) {
	if len(userName) > MaxUserNameLength {
		return false, UserNameTooLongError
	}
	reg := `^[a-zA-Z0-9\\.@\\-_]+$`
	rgx := regexp.MustCompile(reg)
	if !rgx.MatchString(userName) {
		return false, UserNameInvalidError
	}
	return true, nil
}

func IsDisplayNameValid(displayName string) (bool, error) {
	if displayName == "" {
		return true, nil
	}

	if len(displayName) > MaxDisplayNameLength {
		return true, DisplayNameTooLongError
	}

	// TODO: fix me: regexp bug
	// reg := `^[a-zA-Z0-9\\.@\\-\u4e00-\u9fa5]+$`
	// rgx := regexp.MustCompile(reg)
	// if !rgx.MatchString(displayName) {
	// 	return false, DisplayNameInvalidError
	// }

	return true, nil
}

func IsPhoneValid(phone string) (bool, error) {
	reg := `^1([38][0-9]|14[57]|5[^4])\d{8}$`
	rgx := regexp.MustCompile(reg)
	if phone != "" && !rgx.MatchString(phone) {
		return false, MobilePhoneInvalidError
	}
	return true, nil
}

func IsEmailValid(email string) (bool, error) {
	reg := `^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+`
	rgx := regexp.MustCompile(reg)
	if email != "" && !rgx.MatchString(email) {
		return false, EmailInvalidError
	}

	return true, nil
}

func IsCommentsValid(comments string) (bool, error) {
	if len(comments) > MaxCommentLength {
		return false, CommentTooLongError
	}
	return true, nil
}

func (user *User) validate() error {
	ok, err := IsUserNameValid(user.userName)
	if !ok {
		return err
	}

	ok, err = IsDisplayNameValid(user.displayName)
	if !ok {
		return err
	}

	ok, err = IsCommentsValid(user.comments)
	if !ok {
		return err
	}

	ok, err = IsEmailValid(user.email)
	if !ok {
		return err
	}

	ok, err = IsPhoneValid(user.phone)
	if !ok {
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
