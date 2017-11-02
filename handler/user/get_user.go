package user

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
)

type GetUserApi struct {
	req    *http.Request
	status int
	err    error
	user   User
}

func (gua *GetUserApi) Parse() {
	params := util.ParseParameters(gua.req)
	gua.user.userName = params["UserName"]
}

func (gua *GetUserApi) Validate() {
	if gua.user.userName == "" {
		gua.err = MissUserNameError
		gua.status = http.StatusBadRequest
		return
	}
}

func (gua *GetUserApi) Auth() {
	gua.err = doAuth(gua.req)
	if gua.err != nil {
		gua.status = http.StatusForbidden
	}
}

func (gua *GetUserApi) Response() {
	json := simplejson.New()
	if gua.err == nil {
		userJson := simplejson.New()
		userJson.Set("UserId", gua.user.userId)
		userJson.Set("UserName", gua.user.userName)
		userJson.Set("DisplayName", gua.user.displayName)
		userJson.Set("MobilePhone", gua.user.phone)
		userJson.Set("Email", gua.user.email)
		userJson.Set("Comments", gua.user.comments)
		userJson.Set("CreateDate", gua.user.createDate)
		json.Set("User", userJson)
	} else {
		context.Set(gua.req, "request_error", gerror.NewIAMError(gua.status, gua.err))
		json.Set("ErrorMessage", gua.err.Error())
	}
	json.Set("RequestId", context.Get(gua.req, "request_id"))
	data, _ := json.Encode()
	context.Set(gua.req, "response", data)
}

func (gua *GetUserApi) getUser() {
	var bean db.UserBean

	gua.err = db.ActiveService().GetIamUser(gua.user.account, gua.user.userName, &bean)
	if gua.err != nil {
		if gua.err == db.UserNotExistError {
			gua.status = http.StatusNotFound
		} else {
			gua.status = http.StatusInternalServerError
		}
		return
	}

	gua.user.userId = bean.UserId.Hex()
	gua.user.displayName = bean.DisplayName
	gua.user.phone = bean.Phone
	gua.user.email = bean.Email
	gua.user.comments = bean.Comments
	gua.user.createDate = bean.CreateDate
	gua.user.password = bean.Password
}

func GetIAMUserHandler(w http.ResponseWriter, r *http.Request) {
	gua := GetUserApi{req: r, status: http.StatusOK}

	defer gua.Response()

	if gua.Auth(); gua.err != nil {
		return
	}

	gua.Parse()

	if gua.Validate(); gua.err != nil {
		return
	}

	if gua.getUser(); gua.err != nil {
		return
	}
}
