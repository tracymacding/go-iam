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

type ListUserGroupApi struct {
	req     *http.Request
	status  int
	err     error
	user    string
	userId  string
	groups  []*UserGroup
	account string
}

func (luga *ListUserGroupApi) Parse() {
	params := util.ParseParameters(luga.req)
	luga.user = params["UserName"]
}

func (luga *ListUserGroupApi) Validate() {
	if luga.user == "" {
		luga.err = user.MissUserNameError
		luga.status = http.StatusBadRequest
		return
	}
}

func (luga *ListUserGroupApi) Auth() {
	luga.err = doAuth(luga.req)
	if luga.err != nil {
		luga.status = http.StatusForbidden
	}
}

func (luga *ListUserGroupApi) Response() {
	json := simplejson.New()
	if luga.err != nil {
		context.Set(luga.req, "request_error", gerror.NewIAMError(luga.status, luga.err))
		json.Set("ErrorMessage", luga.err.Error())
	} else {
		jsons := make([]*simplejson.Json, 0)
		for _, g := range luga.groups {
			j := g.Json()
			jsons = append(jsons, j)
		}
		json.Set("Groups", jsons)
		json.Set("UserName", luga.user)
	}
	json.Set("RequestId", context.Get(luga.req, "request_id"))
	data, _ := json.Encode()
	context.Set(luga.req, "response", data)
}

func (luga *ListUserGroupApi) listUserGroup() {
	beans := make([]*db.GroupUserBean, 0)
	luga.err = db.ActiveService().ListUserGroup(luga.userId, &beans)
	if luga.err != nil {
		luga.status = http.StatusInternalServerError
		return
	}
	for _, bean := range beans {
		var group db.GroupBean
		luga.err = db.ActiveService().GetGroupById(bean.GroupId, &group)
		if luga.err != nil {
			return
		}
		ug := &UserGroup{
			groupName:  group.GroupName,
			comments:   group.Comments,
			createDate: group.CreateDate,
			joinDate:   bean.JoinDate,
		}
		luga.groups = append(luga.groups, ug)
	}
}

func ListUserGroupHandler(w http.ResponseWriter, r *http.Request) {
	luga := ListUserGroupApi{req: r, status: http.StatusOK}
	defer luga.Response()

	if luga.Auth(); luga.err != nil {
		return
	}

	if luga.Parse(); luga.err != nil {
		return
	}

	if luga.Validate(); luga.err != nil {
		return
	}

	userId, err := user.GetUserId(luga.account, luga.user)
	if err != nil {
		luga.err = err
		return
	}
	luga.userId = userId

	if luga.listUserGroup(); luga.err != nil {
		return
	}
}
