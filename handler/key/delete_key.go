package key

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
)

type DeleteKeyApi struct {
	req    *http.Request
	status int
	err    error
	key    Key
}

func (dka *DeleteKeyApi) Parse() {
	params := util.ParseParameters(dka.req)
	dka.key.accessKeyId = params["UserAccessKeyId"]
}

var (
	MissAccessKeyIdError = errors.New("The user access key does not exist")
)

func (dka *DeleteKeyApi) Validate() {
	if dka.key.accessKeyId == "" {
		dka.err = MissAccessKeyIdError
		dka.status = http.StatusBadRequest
		return
	}
}

func (dka *DeleteKeyApi) Auth() {
	dka.err = doAuth(dka.req)
	if dka.err != nil {
		dka.status = http.StatusForbidden
	}
}

func (dka *DeleteKeyApi) Response() {
	json := simplejson.New()
	if dka.err != nil {
		context.Set(dka.req, "request_error", gerror.NewIAMError(dka.status, dka.err))
		json.Set("ErrorMessage", dka.err.Error())
	}
	json.Set("RequestId", context.Get(dka.req, "request_id"))
	data, _ := json.Encode()
	context.Set(dka.req, "response", data)
}

func (dka *DeleteKeyApi) deleteKey() {
	dka.err = db.ActiveService().DeleteKey(dka.key.accessKeyId)
	if dka.err != nil {
		if dka.err == db.KeyNotExistError {
			dka.status = http.StatusNotFound
		} else {
			dka.status = http.StatusInternalServerError
		}
	}
}

func DeleteKeyHandler(w http.ResponseWriter, r *http.Request) {
	dka := DeleteKeyApi{req: r, status: http.StatusOK}

	defer dka.Response()

	if dka.Auth(); dka.err != nil {
		return
	}

	if dka.Parse(); dka.err != nil {
		return
	}

	if dka.Validate(); dka.err != nil {
		return
	}

	if dka.deleteKey(); dka.err != nil {
		return
	}
}
