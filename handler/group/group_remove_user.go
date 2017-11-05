package group

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/user"
	"github.com/go-iam/handler/util"
	"net/http"
)

type GroupRemoveUserApi struct {
	req     *http.Request
	status  int
	err     error
	group   string
	groupId string
	user    string
	userId  string
	account string
}

func (grua *GroupRemoveUserApi) Parse() {
	params := util.ParseParameters(grua.req)
	grua.group = params["GroupName"]
	grua.user = params["UserName"]
}

func (grua *GroupRemoveUserApi) Validate() {
	if grua.group == "" {
		grua.err = MissGroupNameError
		grua.status = http.StatusBadRequest
		return
	}
	if grua.user == "" {
		grua.err = user.MissUserNameError
		grua.status = http.StatusBadRequest
		return
	}
}

func (grua *GroupRemoveUserApi) Auth() {
	grua.err = doAuth(grua.req)
	if grua.err != nil {
		grua.status = http.StatusForbidden
	}
}

func (grua *GroupRemoveUserApi) Response() {
	json := simplejson.New()
	if grua.err != nil {
		context.Set(grua.req, "request_error", gerror.NewIAMError(grua.status, grua.err))
		json.Set("ErrorMessage", grua.err.Error())
	}
	json.Set("RequestId", context.Get(grua.req, "request_id"))
	data, _ := json.Encode()
	context.Set(grua.req, "response", data)
}

func (grua *GroupRemoveUserApi) groupRemoveUser() {
	bean := &db.GroupUserBean{
		GroupId: grua.groupId,
		UserId:  grua.userId,
	}
	grua.err = db.ActiveService().GroupRemoveUser(bean)
	if grua.err != nil {
		if grua.err == db.UserNotJoinedGroupError {
			grua.status = http.StatusNotFound
		} else {
			grua.status = http.StatusInternalServerError
		}
		return
	}
}

func GroupRemoveUserHandler(w http.ResponseWriter, r *http.Request) {
	grua := GroupRemoveUserApi{req: r, status: http.StatusOK}
	defer grua.Response()

	if grua.Auth(); grua.err != nil {
		return
	}

	if grua.Parse(); grua.err != nil {
		return
	}

	if grua.Validate(); grua.err != nil {
		return
	}

	userId, err := user.GetUserId(grua.account, grua.user)
	if err != nil {
		grua.err = err
		return
	}
	grua.userId = userId

	groupId, err := GetGroupId(grua.account, grua.group)
	if err != nil {
		grua.err = err
		return
	}
	grua.groupId = groupId

	if grua.groupRemoveUser(); grua.err != nil {
		return
	}
}
