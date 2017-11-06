package user

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

type ListUserApi struct {
	req      *http.Request
	status   int
	err      error
	users    []*User
	marker   string
	maxItems int
	account  string
}

func (lua *ListUserApi) Parse() {
	params := util.ParseParameters(lua.req)
	lua.marker = params["Marker"]
	items := params["MaxItems"]

	if items == "" {
		lua.maxItems = 100
	} else {
		lua.maxItems, lua.err = strconv.Atoi(items)
	}
}

var (
	InvalidMaxItemsError = errors.New("MaxItems parameter out of range")
)

func (lua *ListUserApi) Validate() {
	if lua.maxItems < 1 || lua.maxItems > 1000 {
		lua.err = InvalidMaxItemsError
		lua.status = http.StatusBadRequest
		return
	}
}

func (lua *ListUserApi) Auth() {
	lua.err = doAuth(lua.req)
	if lua.err != nil {
		lua.status = http.StatusForbidden
	}
}

func (lua *ListUserApi) Response() {
	json := simplejson.New()
	if lua.err == nil {
		jsons := make([]*simplejson.Json, 0)
		for _, user := range lua.users {
			j := user.Json()
			jsons = append(jsons, j)
		}
		json.Set("Users", jsons)
	} else {
		json.Set("ErrorMessage", lua.err.Error())
		gerr := gerror.NewIAMError(lua.status, lua.err)
		context.Set(lua.req, "request_error", gerr)
	}
	json.Set("RequestId", context.Get(lua.req, "request_id"))
	data, _ := json.Encode()
	context.Set(lua.req, "response", data)
}

func (lua *ListUserApi) listUser() {
	beans := make([]*db.UserBean, 0)

	lua.err = db.ActiveService().ListIamUser(
		lua.account,
		lua.marker,
		lua.maxItems,
		&beans)
	if lua.err != nil {
		lua.status = http.StatusInternalServerError
		return
	}

	for _, bean := range beans {
		usr := FromBean(bean)
		lua.users = append(lua.users, &usr)
	}
}

func ListIAMUserHandler(w http.ResponseWriter, r *http.Request) {
	lua := ListUserApi{req: r, status: http.StatusOK}
	defer lua.Response()

	if lua.Auth(); lua.err != nil {
		return
	}

	if lua.Parse(); lua.err != nil {
		return
	}

	if lua.Validate(); lua.err != nil {
		return
	}

	if lua.listUser(); lua.err != nil {
		return
	}
}
