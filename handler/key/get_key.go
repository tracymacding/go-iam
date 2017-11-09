package key

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
)

type GetKeyApi struct {
	req    *http.Request
	status int
	err    error
	key    Key
}

func (gua *GetKeyApi) Parse() {
	params := util.ParseParameters(gua.req)
	gua.key.accessKeyId = params["AccessKeyId"]
}

func (gua *GetKeyApi) Validate() {
	if gua.key.accessKeyId == "" {
		gua.err = MissAccessKeyIdError
		gua.status = http.StatusBadRequest
		return
	}
}

func (gua *GetKeyApi) Auth() {
	gua.err = doAuth(gua.req)
	if gua.err != nil {
		gua.status = http.StatusForbidden
	}
}

func (gua *GetKeyApi) Response() {
	json := simplejson.New()
	if gua.err == nil {
		j := gua.key.Json()
		json.Set("AccessKey", j)
	} else {
		gerr := gerror.NewIAMError(gua.status, gua.err)
		context.Set(gua.req, "request_error", gerr)
		json.Set("ErrorMessage", gua.err.Error())
	}
	json.Set("RequestId", context.Get(gua.req, "request_id"))
	data, _ := json.Encode()
	context.Set(gua.req, "response", data)
}

func (gua *GetKeyApi) getKey() {
	var bean db.KeyBean

	gua.err = db.ActiveService().GetKey(gua.key.accessKeyId, &bean)
	if gua.err != nil {
		if gua.err == db.KeyNotExistError {
			gua.status = http.StatusNotFound
		} else {
			gua.status = http.StatusInternalServerError
		}
		return
	}
	gua.key = FromBean(&bean)
}

func GetKeyHandler(w http.ResponseWriter, r *http.Request) {
	gua := GetKeyApi{req: r, status: http.StatusOK}
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

	if gua.getKey(); gua.err != nil {
		return
	}
}
