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
		j := gua.user.Json()
		json.Set("User", j)
	} else {
		gerr := gerror.NewIAMError(gua.status, gua.err)
		context.Set(gua.req, "request_error", gerr)
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
	gua.user = FromBean(&bean)
}

func GetIAMUserHandler(w http.ResponseWriter, r *http.Request) {
	gua := GetUserApi{req: r, status: http.StatusOK}

	defer gua.Response()

	if gua.Auth(); gua.err != nil {
		return
	}

	if gua.Parse(); gua.err != nil {
		return
	}

	if gua.Validate(); gua.err != nil {
		return
	}

	if gua.getUser(); gua.err != nil {
		return
	}
}
