package util

import (
	"net/http"
)

func ParseParameters(r *http.Request) map[string]string {
	params := make(map[string]string, 0)
	vals := r.URL.Query()
	for k, _ := range vals {
		params[k] = vals.Get(k)
	}
	return params
}
