package group

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/user"
	"github.com/go-iam/handler/util"
	"net/http"
	"time"
)

type GroupAddUserApi struct {
	req     *http.Request
	status  int
	err     error
	group   string
	groupId string
	user    string
	userId  string
	account string
}

func (gaua *GroupAddUserApi) Parse() {
	params := util.ParseParameters(gaua.req)
	gaua.group = params["GroupName"]
	gaua.user = params["UserName"]
}

func (gaua *GroupAddUserApi) Validate() {
	if gaua.group == "" {
		gaua.err = MissGroupNameError
		gaua.status = http.StatusBadRequest
		return
	}
	if gaua.user == "" {
		gaua.err = user.MissUserNameError
		gaua.status = http.StatusBadRequest
		return
	}
	if ok, err := user.IsUserNameValid(gaua.user); !ok {
		gaua.err = err
		gaua.status = http.StatusBadRequest
		return
	}
	if ok, err := IsGroupNameValid(gaua.group); !ok {
		gaua.err = err
		gaua.status = http.StatusBadRequest
		return
	}
}

func (gaua *GroupAddUserApi) Auth() {
	gaua.err = doAuth(gaua.req)
	if gaua.err != nil {
		gaua.status = http.StatusForbidden
	}
}

func (gaua *GroupAddUserApi) Response() {
	json := simplejson.New()
	if gaua.err != nil {
		gerr := gerror.NewIAMError(gaua.status, gaua.err)
		context.Set(gaua.req, "request_error", gerr)
		json.Set("ErrorMessage", gaua.err.Error())
	}
	json.Set("RequestId", context.Get(gaua.req, "request_id"))
	data, _ := json.Encode()
	context.Set(gaua.req, "response", data)
}

const (
	MAX_GROUP_PER_USER_ATTACHED = 100
)

var (
	TooManyGroupUserJoinedError = errors.New("The count of groups the target user joined beyond the current limits")
)

func (gaua *GroupAddUserApi) groupAddUser() {
	userId, err := user.GetUserId(gaua.account, gaua.user)
	if err != nil {
		gaua.err = err
		return
	}
	gaua.userId = userId

	groupId, err := GetGroupId(gaua.account, gaua.group)
	if err != nil {
		gaua.err = err
		return
	}
	gaua.groupId = groupId

	cnt := 0
	cnt, gaua.err = db.ActiveService().UserJoinedGroupsNum(gaua.userId)
	if gaua.err != nil {
		gaua.status = http.StatusInternalServerError
		return
	}

	if cnt >= MAX_GROUP_PER_USER_ATTACHED {
		gaua.status = http.StatusConflict
		gaua.err = TooManyGroupUserJoinedError
		return
	}

	bean := &db.GroupUserBean{
		GroupId:  gaua.groupId,
		UserId:   gaua.userId,
		JoinDate: time.Now().Format(time.RFC3339),
	}
	bean, gaua.err = db.ActiveService().GroupAddUser(bean)
	if gaua.err != nil {
		if gaua.err == db.UserJoinedGroupError {
			gaua.status = http.StatusConflict
		} else {
			gaua.status = http.StatusInternalServerError
		}
		return
	}
}

func GroupAddUserHandler(w http.ResponseWriter, r *http.Request) {
	gaua := GroupAddUserApi{req: r, status: http.StatusOK}
	defer gaua.Response()

	if gaua.Auth(); gaua.err != nil {
		return
	}

	if gaua.Parse(); gaua.err != nil {
		return
	}

	if gaua.Validate(); gaua.err != nil {
		return
	}

	if gaua.groupAddUser(); gaua.err != nil {
		return
	}
}
