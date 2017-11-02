package user

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
)

type UpdateUserApi struct {
	req     *http.Request
	status  int
	err     error
	user    User
	newUser string
}

func (uua *UpdateUserApi) Parse() {
	params := util.ParseParameters(uua.req)
	uua.user.userName = params["UserName"]
	if params["NewUserName"] != "" {
		uua.newUser = params["NewUserName"]
	}
	if params["NewMobilePhone"] != "" {
		uua.user.phone = params["NewMobilePhone"]
	}
	if params["NewEmail"] != "" {
		uua.user.email = params["NewEmail"]
	}
	if params["NewDisplayName"] != "" {
		uua.user.displayName = params["NewDisplayName"]
	}
	if params["NewComments"] != "" {
		uua.user.comments = params["NewComments"]
	}
}

func (uua *UpdateUserApi) Validate() {
	if uua.user.userName == "" {
		uua.err = MissUserNameError
		uua.status = http.StatusBadRequest
		return
	}
}

func (uua *UpdateUserApi) Auth() {
	uua.err = doAuth(uua.req)
	if uua.err != nil {
		uua.status = http.StatusForbidden
	}
}

func (uua *UpdateUserApi) updateUser() {
	bean := db.UserBean{
		UserName:    uua.user.userName,
		DisplayName: uua.user.displayName,
		Phone:       uua.user.phone,
		Email:       uua.user.email,
		Comments:    uua.user.comments,
		Password:    uua.user.password,
		CreateDate:  uua.user.createDate,
	}
	if uua.newUser != "" {
		bean.UserName = uua.newUser
	}
	user, account := uua.user.userName, uua.user.account
	uua.err = db.ActiveService().UpdateIamUser(user, account, &bean)
	if uua.err == db.AccountNotExistError {
		uua.status = http.StatusNotFound
	} else if uua.err == db.AccountExistError {
		uua.status = http.StatusConflict
	} else {
		uua.status = http.StatusInternalServerError
	}
}

func (uua *UpdateUserApi) Response() {
	json := simplejson.New()
	if uua.err == nil {
		userJson := simplejson.New()
		userJson.Set("UserId", uua.user.userId)
		userJson.Set("UserName", uua.user.userName)
		userJson.Set("DisplayName", uua.user.displayName)
		userJson.Set("MobilePhone", uua.user.phone)
		userJson.Set("Email", uua.user.email)
		userJson.Set("Comments", uua.user.comments)
		userJson.Set("CreateDate", uua.user.createDate)
		json.Set("User", userJson)
	} else {
		json.Set("ErrorMessage", uua.err.Error())
		context.Set(uua.req, "request_error", gerror.NewIAMError(uua.status, uua.err))
	}
	json.Set("RequestId", context.Get(uua.req, "request_id"))
	data, _ := json.Encode()
	context.Set(uua.req, "response", data)
}

func UpdateIAMUserHandler(w http.ResponseWriter, r *http.Request) {
	uua := UpdateUserApi{req: r, status: http.StatusOK}
	defer uua.Response()

	if uua.Auth(); uua.err != nil {
		return
	}

	uua.Parse()

	if uua.Validate(); uua.err != nil {
		return
	}

	gua := GetUserApi{}
	gua.user.userName = uua.user.userName
	gua.user.account = uua.user.account

	if gua.getUser(); gua.err != nil {
		uua.err = gua.err
		return
	}

	if uua.user.displayName == "" {
		uua.user.displayName = gua.user.displayName
	}
	if uua.user.phone == "" {
		uua.user.phone = gua.user.phone
	}
	if uua.user.email == "" {
		uua.user.email = gua.user.email
	}
	if uua.user.comments == "" {
		uua.user.comments = gua.user.comments
	}
	if uua.user.password == "" {
		uua.user.password = gua.user.password
	}

	uua.updateUser()
}