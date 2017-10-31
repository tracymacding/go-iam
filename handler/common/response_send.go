package common

import (
	"encoding/base64"
	"github.com/go-iam/context"
	"net/http"
)

const (
	base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"
)

func base64Encode(src []byte) string {
	return base64.NewEncoding(base64Table).EncodeToString(src)
}

func SendResponseHandler(w http.ResponseWriter, r *http.Request) {
	resp := context.Get(r, "response").([]byte)
	if resp != nil {
		w.Header().Set("Server", "GO-IAM")
		w.Write(resp)
	}
}
