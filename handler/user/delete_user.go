package user

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
)

type DeleteUserApi struct {
	req    *http.Request
	status int
	err    error
	user   User
}

func (dua *DeleteUserApi) Parse() {
	params := util.ParseParameters(dua.req)
	dua.user.userName = params["UserName"]
}

func (dua *DeleteUserApi) Validate() {
	if dua.user.userName == "" {
		dua.err = MissUserNameError
		dua.status = http.StatusBadRequest
		return
	}
}

func (dua *DeleteUserApi) Auth() {
	dua.err = doAuth(dua.req)
	if dua.err != nil {
		dua.status = http.StatusForbidden
	}
}

func (dua *DeleteUserApi) Response() {
	json := simplejson.New()
	if dua.err != nil {
		context.Set(dua.req, "request_error", gerror.NewIAMError(dua.status, dua.err))
		json.Set("ErrorMessage", dua.err.Error())
	}
	json.Set("RequestId", context.Get(dua.req, "request_id"))
	data, _ := json.Encode()
	context.Set(dua.req, "response", data)
}

func (dua *DeleteUserApi) deleteUser() {
	dua.err = db.ActiveService().DeleteIamUser(dua.user.account, dua.user.userName)
	if dua.err != nil {
		if dua.err == db.AccountNotExistError {
			dua.status = http.StatusNotFound
		} else {
			dua.status = http.StatusInternalServerError
		}
	}
}

func DeleteIAMUserHandler(w http.ResponseWriter, r *http.Request) {
	dua := DeleteUserApi{req: r, status: http.StatusOK}

	defer dua.Response()

	if dua.Auth(); dua.err != nil {
		return
	}

	dua.Parse()

	if dua.Validate(); dua.err != nil {
		return
	}

	if dua.deleteUser(); dua.err != nil {
		return
	}
}
