package key

import (
	// "github.com/go-iam/db"
	// "github.com/go-iam/mux"
	// "github.com/go-iam/security"
	"net/http"
)

func doAuth(r *http.Request) error {
	// keyInfo := db.KeyBean{}
	// err := db.ActiveService().GetKey(mux.Vars(r)["AccessKeyId"], &keyInfo)
	// if err != nil {
	// 	return err
	// }

	// owner, err := GetKeyOwner(keyInfo.AccessKeyId.Hex(), keyInfo.Entitype)
	// if err != nil {
	// 	return err
	// }

	// resource := "ccs:iam:*:" + owner + ":user/*"
	// return security.DoAuth(r.Method, "CreateIamUser", resource, mux.Vars(r))
	return nil
}
