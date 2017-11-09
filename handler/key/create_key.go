package key

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/user"
	"github.com/go-iam/handler/util"
	"net/http"
	"time"
)

type CreateKeyApi struct {
	req    *http.Request
	status int
	err    error
	key    Key
}

var (
	TooManyKeysError = errors.New("The count of keys beyond the current limits")
)

func (cka *CreateKeyApi) Validate() {
	if cka.key.owner != "" {
		if ok, err := user.IsUserNameValid(cka.key.owner); !ok {
			cka.err = err
			cka.status = http.StatusBadRequest
			return
		}
	}
}

func (cka *CreateKeyApi) Parse() {
	params := util.ParseParameters(cka.req)
	cka.key.owner = params["UserName"]
	if cka.key.owner != "" {
		cka.key.ownerType = IAMUser
	}
}

func (cka *CreateKeyApi) Auth() {
	cka.err = doAuth(cka.req)
	if cka.err != nil {
		cka.status = http.StatusForbidden
	}
}

func (cka *CreateKeyApi) Response() {
	json := simplejson.New()
	if cka.err == nil {
		j := cka.key.Json()
		json.Set("AccessKey", j)
	} else {
		gerr := gerror.NewIAMError(cka.status, cka.err)
		context.Set(cka.req, "request_error", gerr)
		json.Set("ErrorMessage", cka.err.Error())
	}
	json.Set("RequestId", context.Get(cka.req, "request_id"))
	data, _ := json.Encode()
	context.Set(cka.req, "response", data)
}

const (
	MAX_KEY_PER_ENTITY = 2
)

func (cka *CreateKeyApi) createKey() {
	cnt := 0
	cnt, cka.err = db.ActiveService().KeyCountOfEntity(cka.key.owner, int(cka.key.ownerType))
	if cka.err != nil {
		cka.status = http.StatusInternalServerError
		return
	}

	if cnt >= MAX_KEY_PER_ENTITY {
		cka.status = http.StatusConflict
		cka.err = TooManyKeysError
		return
	}

	bean := cka.key.ToBean()
	bean.CreateDate = time.Now().Format(time.RFC3339)
	bean.Status = int(Active)
	cka.err = db.ActiveService().CreateKey(&bean)
	if cka.err != nil {
		cka.status = http.StatusInternalServerError
		return
	}
	cka.key = FromBean(&bean)
}

func CreateKeyHandler(w http.ResponseWriter, r *http.Request) {
	cka := CreateKeyApi{req: r, status: http.StatusOK}
	defer cka.Response()

	if cka.Auth(); cka.err != nil {
		return
	}

	if cka.Parse(); cka.err != nil {
		return
	}

	if cka.Validate(); cka.err != nil {
		return
	}

	if cka.createKey(); cka.err != nil {
		return
	}
}
