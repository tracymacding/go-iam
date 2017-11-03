package key

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/context"
	"github.com/go-iam/db"
	"github.com/go-iam/gerror"
	"github.com/go-iam/handler/util"
	"net/http"
)

type ListKeyApi struct {
	req      *http.Request
	status   int
	err      error
	keys     []*Key
	entity   string
	entitype KeyOwnerType
}

func (lka *ListKeyApi) Parse() {
	params := util.ParseParameters(lka.req)
	lka.entity = params["UserName"]
	if lka.entity != "" {
		lka.entitype = IAMUser
	} else {
		lka.entitype = Account
	}
}

func (lka *ListKeyApi) Validate() {
}

func (lka *ListKeyApi) Auth() {
	lka.err = doAuth(lka.req)
	if lka.err != nil {
		lka.status = http.StatusForbidden
	}
}

func (lka *ListKeyApi) Response() {
	json := simplejson.New()
	if lka.err == nil {
		jsons := make([]*simplejson.Json, 0)
		for _, key := range lka.keys {
			j := key.Json()
			jsons = append(jsons, j)
		}
		json.Set("AccessKeys", jsons)
	} else {
		json.Set("ErrorMessage", lka.err.Error())
		context.Set(lka.req, "request_error", gerror.NewIAMError(lka.status, lka.err))
	}
	json.Set("RequestId", context.Get(lka.req, "request_id"))
	data, _ := json.Encode()
	context.Set(lka.req, "response", data)
}

func (lka *ListKeyApi) listKey() {
	beans := make([]*db.KeyBean, 0)

	lka.err = db.ActiveService().ListKey(
		lka.entity,
		int(lka.entitype),
		&beans)
	if lka.err != nil {
		lka.status = http.StatusInternalServerError
		return
	}

	for _, bean := range beans {
		key := &Key{
			accessKeyId:     bean.AccessKeyId.Hex(),
			accessKeySecret: bean.AccessKeySecret,
			status:          KeyStatus(bean.Status),
			createDate:      bean.CreateDate,
		}
		lka.keys = append(lka.keys, key)
	}
}

func ListKeyHandler(w http.ResponseWriter, r *http.Request) {
	lka := ListKeyApi{req: r, status: http.StatusOK}
	defer lka.Response()

	if lka.Auth(); lka.err != nil {
		return
	}

	if lka.Parse(); lka.err != nil {
		return
	}

	if lka.Validate(); lka.err != nil {
		return
	}

	if lka.listKey(); lka.err != nil {
		return
	}
}
