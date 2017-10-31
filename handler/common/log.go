package common

import (
	"fmt"
	"github.com/go-iam/context"
	"net/http"
	"strings"
	"time"
)

func LogHandler(w http.ResponseWriter, r *http.Request) {
	reqId := context.Get(r, "request_id").(string)
	start := context.Get(r, "request_start").(string)
	startUnix := context.Get(r, "request_start_unix").(int64)
	queries := ""
	Queries := r.URL.Query()
	for k, v := range Queries {
		queries = fmt.Sprintf("%s%s=%s&", queries, k, v[0])
	}
	queries = strings.Trim(queries, "&")
	url := r.Host + r.URL.Path
	if queries != "" {
		url = url + "?" + queries
	}

	if err := context.Get(r, "request_error"); err != nil {
		fmt.Printf("Req:%s, Method:%s, URL:%s, Start:%s, Cost:%d us, Error:%s\n", reqId, r.Method, url, start, (time.Now().UnixNano()-startUnix)/1000.0, err)
	} else {
		fmt.Printf("Req:%s, Method:%s, URL:%s, Start:%s, Cost:%d us\n", reqId, r.Method, url, start, (time.Now().UnixNano()-startUnix)/1000.0)
	}
}
