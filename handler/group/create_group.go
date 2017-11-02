package group

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
	"time"
)

type CreateGroupApi struct {
	req    *http.Request
	status int
	err    error
	group  Group
}

var (
	MissGroupNameError = errors.New("user name missing")
	TooManyGroupsError = errors.New("The count of group beyond the current limits")
)

func (cga *CreateGroupApi) Parse() {
	params := util.ParseParameters(cga.req)
	cga.group.groupName = params["GroupName"]
	cga.group.comments = params["Comments"]
}

func (cga *CreateGroupApi) Validate() {
	if cga.group.groupName == "" {
		cga.err = MissGroupNameError
		cga.status = http.StatusBadRequest
		return
	}
}

func (cga *CreateGroupApi) Auth() {
	cga.err = doAuth(cga.req)
	if cga.err != nil {
		cga.status = http.StatusForbidden
	}
}

func (cga *CreateGroupApi) Response() {
	json := simplejson.New()
	if cga.err == nil {
		j := cga.group.Json()
		json.Set("Group", j)
	} else {
		context.Set(cga.req, "request_error", gerror.NewIAMError(cga.status, cga.err))
		json.Set("ErrorMessage", cga.err.Error())
	}
	json.Set("RequestId", context.Get(cga.req, "request_id"))
	data, _ := json.Encode()
	context.Set(cga.req, "response", data)
}

const (
	MAX_GROUP_PER_ACCOUNT = 100
)

func (cga *CreateGroupApi) createGroup() {
	cnt := 0
	cnt, cga.err = db.ActiveService().GroupCountOfAccount(cga.group.account)
	if cga.err != nil {
		cga.status = http.StatusInternalServerError
		return
	}

	if cnt >= MAX_GROUP_PER_ACCOUNT {
		cga.status = http.StatusConflict
		cga.err = TooManyGroupsError
		return
	}

	bean := &db.GroupBean{
		GroupName:  cga.group.groupName,
		Comments:   cga.group.comments,
		Account:    cga.group.account,
		CreateDate: time.Now().Format(time.RFC3339),
	}
	bean, cga.err = db.ActiveService().CreateGroup(bean)
	if cga.err != nil {
		if cga.err == db.GroupExistError {
			cga.status = http.StatusConflict
		} else {
			cga.status = http.StatusInternalServerError
		}
		return
	}
	cga.group.groupId = bean.GroupId.Hex()
	cga.group.createDate = bean.CreateDate
}

func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	cga := CreateGroupApi{req: r, status: http.StatusOK}
	defer cga.Response()

	if cga.Auth(); cga.err != nil {
		return
	}

	if cga.Parse(); cga.err != nil {
		return
	}

	if cga.Validate(); cga.err != nil {
		return
	}

	if cga.createGroup(); cga.err != nil {
		return
	}
}
