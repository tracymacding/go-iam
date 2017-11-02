package group

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
)

type GetGroupApi struct {
	req    *http.Request
	status int
	err    error
	group  Group
}

func (gua *GetGroupApi) Parse() {
	params := util.ParseParameters(gua.req)
	gua.group.groupName = params["GroupName"]
}

func (gua *GetGroupApi) Validate() {
	if gua.group.groupName == "" {
		gua.err = MissGroupNameError
		gua.status = http.StatusBadRequest
		return
	}
}

func (gua *GetGroupApi) Auth() {
	gua.err = doAuth(gua.req)
	if gua.err != nil {
		gua.status = http.StatusForbidden
	}
}

func (gua *GetGroupApi) Response() {
	json := simplejson.New()
	if gua.err == nil {
		j := gua.group.Json()
		json.Set("User", j)
	} else {
		context.Set(gua.req, "request_error", gerror.NewIAMError(gua.status, gua.err))
		json.Set("ErrorMessage", gua.err.Error())
	}
	json.Set("RequestId", context.Get(gua.req, "request_id"))
	data, _ := json.Encode()
	context.Set(gua.req, "response", data)
}

func (gua *GetGroupApi) getGroup() {
	var bean db.GroupBean

	gua.err = db.ActiveService().GetGroup(gua.group.account, gua.group.groupName, &bean)
	if gua.err != nil {
		if gua.err == db.GroupNotExistError {
			gua.status = http.StatusNotFound
		} else {
			gua.status = http.StatusInternalServerError
		}
		return
	}

	gua.group.groupId = bean.GroupId.Hex()
	gua.group.comments = bean.Comments
	gua.group.createDate = bean.CreateDate
}

func GetGroupHandler(w http.ResponseWriter, r *http.Request) {
	gua := GetGroupApi{req: r, status: http.StatusOK}

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

	if gua.getGroup(); gua.err != nil {
		return
	}
}
