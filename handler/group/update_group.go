package group

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type UpdateGroupApi struct {
	req      *http.Request
	status   int
	err      error
	group    Group
	newGroup string
}

func (uga *UpdateGroupApi) Parse() {
	params := util.ParseParameters(uga.req)
	uga.group.groupName = params["GroupName"]
	if params["NewGroupName"] != "" {
		uga.newGroup = params["NewGroupName"]
	}
	if params["NewComments"] != "" {
		uga.group.comments = params["NewComments"]
	}
}

func (uga *UpdateGroupApi) Validate() {
	if uga.group.groupName == "" {
		uga.err = MissGroupNameError
		uga.status = http.StatusBadRequest
		return
	}
	ok, err := IsGroupNameValid(uga.group.groupName)
	if !ok {
		uga.err = err
		uga.status = http.StatusBadRequest
		return
	}
	if uga.newGroup != "" {
		ok, err := IsGroupNameValid(uga.newGroup)
		if !ok {
			uga.err = err
			uga.status = http.StatusBadRequest
			return
		}
	}
	if uga.group.comments != "" {
		ok, err := IsCommentsValid(uga.group.comments)
		if !ok {
			uga.err = err
			uga.status = http.StatusBadRequest
			return
		}
	}
}

func (uga *UpdateGroupApi) Auth() {
	uga.err = doAuth(uga.req)
	if uga.err != nil {
		uga.status = http.StatusForbidden
	}
}

func (uga *UpdateGroupApi) updateGroup() {
	gga := GetGroupApi{}
	gga.group.groupName = uga.group.groupName
	gga.group.account = uga.group.account

	if gga.getGroup(); gga.err != nil {
		uga.err = gga.err
		return
	}

	uga.group.groupId = gga.group.groupId
	if uga.group.comments == "" {
		uga.group.comments = gga.group.comments
	}
	uga.group.createDate = gga.group.createDate

	bean := uga.group.ToBean()
	bean.GroupId = bson.ObjectIdHex(gga.group.groupId)

	if uga.newGroup != "" {
		bean.GroupName = uga.newGroup
	}
	uga.err = db.ActiveService().UpdateGroup(&bean)
	if uga.err == db.GroupNotExistError {
		uga.status = http.StatusNotFound
	} else if uga.err == db.GroupExistError {
		uga.status = http.StatusConflict
	} else {
		uga.status = http.StatusInternalServerError
	}
	uga.group = FromBean(&bean)
}

func (uga *UpdateGroupApi) Response() {
	json := simplejson.New()
	if uga.err == nil {
		j := uga.group.Json()
		json.Set("Group", j)
	} else {
		gerr := gerror.NewIAMError(uga.status, uga.err)
		context.Set(uga.req, "request_error", gerr)
		json.Set("ErrorMessage", uga.err.Error())
	}
	json.Set("RequestId", context.Get(uga.req, "request_id"))
	data, _ := json.Encode()
	context.Set(uga.req, "response", data)
}

func UpdateGroupHandler(w http.ResponseWriter, r *http.Request) {
	uga := UpdateGroupApi{req: r, status: http.StatusOK}
	defer uga.Response()

	if uga.Auth(); uga.err != nil {
		return
	}

	if uga.Parse(); uga.err != nil {
		return
	}

	if uga.Validate(); uga.err != nil {
		return
	}

	if uga.updateGroup(); uga.err != nil {
		return
	}
}
