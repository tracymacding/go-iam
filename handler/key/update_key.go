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

type UpdateKeyApi struct {
	req       *http.Request
	status    int
	err       error
	key       Key
	newStatus KeyStatus
}

var (
	MissKeyStatusError    = errors.New("missing key status")
	InvalidKeyStatusError = errors.New("invalid key status")
)

func (uka *UpdateKeyApi) Parse() {
	params := util.ParseParameters(uka.req)
	uka.key.accessKeyId = params["UserAccessKeyId"]
	status := params["Status"]
	if status == "" {
		uka.err = MissKeyStatusError
		return
	}
	uka.newStatus = ParseKeyStatus(status)
}

func (uka *UpdateKeyApi) Validate() {
	if uka.key.accessKeyId == "" {
		uka.err = MissAccessKeyIdError
		uka.status = http.StatusBadRequest
		return
	}
	if !uka.newStatus.IsValid() {
		uka.err = InvalidKeyStatusError
		uka.status = http.StatusBadRequest
		return
	}
}

func (uka *UpdateKeyApi) Auth() {
	uka.err = doAuth(uka.req)
	if uka.err != nil {
		uka.status = http.StatusForbidden
	}
}

func (uka *UpdateKeyApi) updateKey() {
	gka := GetKeyApi{}
	gka.key.accessKeyId = uka.key.accessKeyId

	if gka.getKey(); gka.err != nil {
		uka.err = gka.err
		if gka.err == db.KeyNotExistError {
			uka.status = http.StatusNotFound
		} else {
			uka.status = http.StatusInternalServerError
		}
		return
	}

	// key status not changed
	if uka.newStatus == KeyStatus(gka.key.status) {
		return
	}

	uka.key.accessKeyId = gka.key.accessKeyId
	uka.key.accessKeySecret = gka.key.accessKeySecret
	uka.key.createDate = gka.key.createDate
	uka.key.owner = gka.key.owner
	uka.key.ownerType = gka.key.ownerType
	uka.key.status = uka.newStatus

	bean := uka.key.ToBean()
	uka.err = db.ActiveService().UpdateKey(uka.key.accessKeyId, &bean)
	if uka.err == db.KeyNotExistError {
		uka.status = http.StatusNotFound
	} else {
		uka.status = http.StatusInternalServerError
	}
}

func (uka *UpdateKeyApi) Response() {
	json := simplejson.New()
	if uka.err == nil {
		j := uka.key.Json()
		json.Set("Key", j)
	} else {
		gerr := gerror.NewIAMError(uka.status, uka.err)
		context.Set(uka.req, "request_error", gerr)
		json.Set("ErrorMessage", uka.err.Error())
	}
	json.Set("RequestId", context.Get(uka.req, "request_id"))
	data, _ := json.Encode()
	context.Set(uka.req, "response", data)
}

func UpdateKeyHandler(w http.ResponseWriter, r *http.Request) {
	uka := UpdateKeyApi{req: r, status: http.StatusOK}
	defer uka.Response()

	if uka.Auth(); uka.err != nil {
		return
	}

	if uka.Parse(); uka.err != nil {
		return
	}

	if uka.Validate(); uka.err != nil {
		return
	}

	if uka.updateKey(); uka.err != nil {
		return
	}
}
