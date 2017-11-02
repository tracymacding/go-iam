package user

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
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

type CreateIAMUserApi struct {
	req    *http.Request
	status int
	err    error
	user   User
}

var (
	MissUserNameError = errors.New("user name missing")
	MissPasswordError = errors.New("password missing")
	TooManyUsersError = errors.New("The count of users beyond the current limits")
)

func (caa *CreateIAMUserApi) Validate() {
	if caa.user.userName == "" {
		caa.err = MissUserNameError
		caa.status = http.StatusBadRequest
		return
	}
	if caa.user.password == "" {
		caa.err = MissPasswordError
		caa.status = http.StatusBadRequest
		return
	}
}

func (caa *CreateIAMUserApi) Parse() {
	params := util.ParseParameters(caa.req)
	caa.user.userName = params["UserName"]
	caa.user.displayName = params["DisplayName"]
	caa.user.phone = params["MobilePhone"]
	caa.user.email = params["Email"]
	caa.user.comments = params["Comment"]
	caa.user.password = params["Password"]
}

func (caa *CreateIAMUserApi) Auth() {
	caa.err = doAuth(caa.req)
	if caa.err != nil {
		caa.status = http.StatusForbidden
	}
}

func (caa *CreateIAMUserApi) Response() {
	json := simplejson.New()
	if caa.err == nil {
		userJson := simplejson.New()
		userJson.Set("UserId", caa.user.userId)
		userJson.Set("UserName", caa.user.userName)
		userJson.Set("DisplayName", caa.user.displayName)
		userJson.Set("MobilePhone", caa.user.phone)
		userJson.Set("Email", caa.user.email)
		userJson.Set("Comments", caa.user.comments)
		userJson.Set("CreateDate", caa.user.createDate)
		json.Set("User", userJson)
	} else {
		context.Set(caa.req, "request_error", gerror.NewIAMError(caa.status, caa.err))
		json.Set("ErrorMessage", caa.err.Error())
	}
	json.Set("RequestId", context.Get(caa.req, "request_id"))
	data, _ := json.Encode()
	context.Set(caa.req, "response", data)
}

const (
	MAX_IAM_USER_PER_ACCOUNT = 100
)

func (caa *CreateIAMUserApi) createIamUser() {
	cnt := 0
	cnt, caa.err = db.ActiveService().UserCountOfAccount(caa.user.account)
	if caa.err != nil {
		caa.status = http.StatusInternalServerError
		return
	}

	if cnt >= MAX_IAM_USER_PER_ACCOUNT {
		caa.status = http.StatusConflict
		caa.err = TooManyUsersError
		return
	}

	bean := &db.UserBean{
		UserName:    caa.user.userName,
		DisplayName: caa.user.displayName,
		Phone:       caa.user.phone,
		Email:       caa.user.email,
		Comments:    caa.user.comments,
		Password:    caa.user.password,
		Account:     caa.user.account,
		CreateDate:  time.Now().Format(time.RFC3339),
	}
	bean, caa.err = db.ActiveService().CreateIamUser(bean)
	if caa.err != nil {
		if caa.err == db.UserExistError {
			caa.status = http.StatusConflict
		} else {
			caa.status = http.StatusInternalServerError
		}
		return
	}
	caa.user.userId = bean.UserId.Hex()
	caa.user.createDate = bean.CreateDate
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	caa := CreateIAMUserApi{req: r, status: http.StatusOK}
	defer caa.Response()

	if caa.Auth(); caa.err != nil {
		return
	}

	caa.Parse()

	if caa.Validate(); caa.err != nil {
		return
	}

	if caa.createIamUser(); caa.err != nil {
		return
	}

	// TODO: create system type ak/sk pair for iam user
}
