package group

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
)

type ListGroupUserApi struct {
	req     *http.Request
	status  int
	err     error
	group   string
	groupId string
	users   []*GroupUser
	account string
}

func (lgua *ListGroupUserApi) Parse() {
	params := util.ParseParameters(lgua.req)
	lgua.group = params["GroupName"]
}

func (lgua *ListGroupUserApi) Validate() {
	if lgua.group == "" {
		lgua.err = MissGroupNameError
		lgua.status = http.StatusBadRequest
		return
	}
	if ok, err := IsGroupNameValid(lgua.group); !ok {
		lgua.err = err
		lgua.status = http.StatusBadRequest
		return
	}
}

func (lgua *ListGroupUserApi) Auth() {
	lgua.err = doAuth(lgua.req)
	if lgua.err != nil {
		lgua.status = http.StatusForbidden
	}
}

func (lgua *ListGroupUserApi) Response() {
	json := simplejson.New()
	if lgua.err != nil {
		context.Set(lgua.req, "request_error", gerror.NewIAMError(lgua.status, lgua.err))
		json.Set("ErrorMessage", lgua.err.Error())
	} else {
		jsons := make([]*simplejson.Json, 0)
		for _, u := range lgua.users {
			j := u.Json()
			jsons = append(jsons, j)
		}
		json.Set("Users", jsons)
		json.Set("GroupName", lgua.group)
	}
	json.Set("RequestId", context.Get(lgua.req, "request_id"))
	data, _ := json.Encode()
	context.Set(lgua.req, "response", data)
}

func (lgua *ListGroupUserApi) listGroupUser() {
	beans := make([]*db.GroupUserBean, 0)
	lgua.err = db.ActiveService().ListGroupUser(lgua.groupId, &beans)
	if lgua.err != nil {
		lgua.status = http.StatusInternalServerError
		return
	}
	for _, bean := range beans {
		var user db.UserBean
		lgua.err = db.ActiveService().GetIamUserById(bean.UserId, &user)
		if lgua.err != nil {
			return
		}
		gu := &GroupUser{
			userId:      user.UserId.Hex(),
			userName:    user.UserName,
			displayName: user.DisplayName,
			createDate:  user.CreateDate,
			joinDate:    bean.JoinDate,
		}
		lgua.users = append(lgua.users, gu)
	}
}

func ListGroupUserHandler(w http.ResponseWriter, r *http.Request) {
	lgua := ListGroupUserApi{req: r, status: http.StatusOK}
	defer lgua.Response()

	if lgua.Auth(); lgua.err != nil {
		return
	}

	if lgua.Parse(); lgua.err != nil {
		return
	}

	if lgua.Validate(); lgua.err != nil {
		return
	}

	groupId, err := GetGroupId(lgua.account, lgua.group)
	if err != nil {
		lgua.err = err
		return
	}
	lgua.groupId = groupId

	if lgua.listGroupUser(); lgua.err != nil {
		return
	}
}
