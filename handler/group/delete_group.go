package group

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
)

type DeleteGroupApi struct {
	req    *http.Request
	status int
	err    error
	group  Group
}

func (dga *DeleteGroupApi) Parse() {
	params := util.ParseParameters(dga.req)
	dga.group.groupName = params["GroupName"]
}

func (dga *DeleteGroupApi) Validate() {
	if dga.group.groupName == "" {
		dga.err = MissGroupNameError
		dga.status = http.StatusBadRequest
		return
	}
	if ok, err := IsGroupNameValid(dga.group.groupName); !ok {
		dga.err = err
		dga.status = http.StatusBadRequest
		return
	}
}

func (dga *DeleteGroupApi) Auth() {
	dga.err = doAuth(dga.req)
	if dga.err != nil {
		dga.status = http.StatusForbidden
	}
}

func (dga *DeleteGroupApi) Response() {
	json := simplejson.New()
	if dga.err != nil {
		gerr := gerror.NewIAMError(dga.status, dga.err)
		context.Set(dga.req, "request_error", gerr)
		json.Set("ErrorMessage", dga.err.Error())
	}
	json.Set("RequestId", context.Get(dga.req, "request_id"))
	data, _ := json.Encode()
	context.Set(dga.req, "response", data)
}

func (dga *DeleteGroupApi) deleteGroup() {
	dga.err = db.ActiveService().DeleteGroup(dga.group.account, dga.group.groupName)
	if dga.err != nil {
		if dga.err == db.GroupNotExistError {
			dga.status = http.StatusNotFound
		} else {
			dga.status = http.StatusInternalServerError
		}
	}
}

func DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	dga := DeleteGroupApi{req: r, status: http.StatusOK}
	defer dga.Response()

	if dga.Auth(); dga.err != nil {
		return
	}

	if dga.Parse(); dga.err != nil {
		return
	}

	if dga.Validate(); dga.err != nil {
		return
	}

	if dga.deleteGroup(); dga.err != nil {
		return
	}
}
