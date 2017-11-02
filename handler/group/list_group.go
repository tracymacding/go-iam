package group

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
	"strconv"
)

type ListGroupApi struct {
	req      *http.Request
	status   int
	err      error
	groups   []*Group
	marker   string
	maxItems int
	account  string
}

func (lga *ListGroupApi) Parse() {
	params := util.ParseParameters(lga.req)
	lga.marker = params["Marker"]
	items := params["MaxItems"]

	if items == "" {
		lga.maxItems = 100
	}
	lga.maxItems, lga.err = strconv.Atoi(items)
}

var (
	InvalidMaxItemsError = errors.New("MaxItems parameter out of range")
)

func (lga *ListGroupApi) Validate() {
	if lga.maxItems < 1 || lga.maxItems > 1000 {
		lga.err = InvalidMaxItemsError
		lga.status = http.StatusBadRequest
		return
	}
}

func (lga *ListGroupApi) Auth() {
	lga.err = doAuth(lga.req)
	if lga.err != nil {
		lga.status = http.StatusForbidden
	}
}

func (lga *ListGroupApi) Response() {
	json := simplejson.New()
	if lga.err == nil {
		jsons := make([]*simplejson.Json, 0)
		for _, group := range lga.groups {
			j := group.Json()
			jsons = append(jsons, j)
		}
		json.Set("Groups", jsons)
	} else {
		json.Set("ErrorMessage", lga.err.Error())
		context.Set(lga.req, "request_error", gerror.NewIAMError(lga.status, lga.err))
	}
	json.Set("RequestId", context.Get(lga.req, "request_id"))
	data, _ := json.Encode()
	context.Set(lga.req, "response", data)
}

func (lga *ListGroupApi) listGroup() {
	beans := make([]*db.GroupBean, 0)

	lga.err = db.ActiveService().ListGroup(
		lga.account,
		lga.marker,
		lga.maxItems,
		&beans)
	if lga.err != nil {
		lga.status = http.StatusInternalServerError
		return
	}

	for _, bean := range beans {
		group := &Group{
			groupId:    bean.GroupId.Hex(),
			groupName:  bean.GroupName,
			comments:   bean.Comments,
			createDate: bean.CreateDate,
		}
		lga.groups = append(lga.groups, group)
	}
}

func ListGroupHandler(w http.ResponseWriter, r *http.Request) {
	lga := ListGroupApi{req: r, status: http.StatusOK}

	defer lga.Response()

	if lga.Auth(); lga.err != nil {
		return
	}

	if lga.Parse(); lga.err != nil {
		return
	}

	if lga.Validate(); lga.err != nil {
		return
	}

	if lga.listGroup(); lga.err != nil {
		return
	}
}
