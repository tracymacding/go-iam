package user

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"net/http"
	"time"
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

// TODO
func (usr *User) validate() (bool, error) {
	return true, nil
}

func parseUserFromRequest(r *http.Request, usr *User) {
	vals := r.URL.Query()
	v, ok := vals["UserName"]
	if ok && len(v) > 0 {
		usr.userName = v[0]
	}
	v, ok = vals["DisplayName"]
	if ok && len(v) > 0 {
		usr.displayName = v[0]
	}
	v, ok = vals["MobilePhone"]
	if ok && len(v) > 0 {
		usr.phone = v[0]
	}
	v, ok = vals["Email"]
	if ok && len(v) > 0 {
		usr.email = v[0]
	}
	v, ok = vals["Comments"]
	if ok && len(v) > 0 {
		usr.comments = v[0]
	}
	v, ok = vals["Password"]
	if ok && len(v) > 0 {
		usr.comments = v[0]
	}
	// TODO:
	usr.account = "test"
}

func createIamUser(usr *User) error {
	userBean := db.UserBean{
		UserName:    usr.userName,
		DisplayName: usr.displayName,
		Phone:       usr.phone,
		Email:       usr.email,
		Comments:    usr.comments,
		Password:    usr.password,
		Account:     usr.account,
		CreateDate:  time.Now().Format(time.RFC3339),
	}
	bean, err := db.ActiveService().CreateIamUser(&userBean)
	if err != nil {
		return err
	}
	usr.userId = bean.UserId.Hex()
	usr.createDate = userBean.CreateDate
	return nil
}

func CreateUserResponse(r *http.Request, usr *User, err error) []byte {
	json := simplejson.New()
	json.Set("RequestId", context.Get(r, "request_id"))
	if err != nil {
		json.Set("ErrorMessage", err.Error())
	} else {
		userJson := simplejson.New()
		userJson.Set("UserId", usr.userId)
		userJson.Set("UserName", usr.userName)
		userJson.Set("DisplayName", usr.displayName)
		userJson.Set("MobilePhone", usr.phone)
		userJson.Set("Email", usr.email)
		userJson.Set("Comments", usr.comments)
		userJson.Set("CreateDate", usr.createDate)
		json.Set("User", userJson)
	}
	data, _ := json.Encode()
	return data
}

const (
	MAX_IAM_USER_PER_ACCOUNT = 100
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var usr User
	var err error
	status := http.StatusOK

	defer func() {
		if err != nil {
			context.Set(r, "request_error", gerror.NewIAMError(status, err))
		}
		resp := CreateUserResponse(r, &usr, err)
		context.Set(r, "response", resp)
	}()

	if err = doAuth(r); err != nil {
		status = http.StatusForbidden
		return
	}

	parseUserFromRequest(r, &usr)

	ok := true
	if ok, err = usr.validate(); !ok {
		status = http.StatusBadRequest
		return
	}

	cnt := 0
	cnt, err = db.ActiveService().UserCountOfAccount(usr.account)
	if err != nil {
		status = http.StatusInternalServerError
		return
	}

	if cnt >= MAX_IAM_USER_PER_ACCOUNT {
		status = http.StatusConflict
		err = fmt.Errorf("The count of users beyond the current limits")
		return
	}

	if err = createIamUser(&usr); err != nil {
		if err == db.UserExistError {
			status = http.StatusConflict
		} else {
			status = http.StatusInternalServerError
		}
		return
	}

	// TODO: create system type ak/sk pair for iam user
}
