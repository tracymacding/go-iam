package handler

import (
	"fmt"
	"galaxy-s3-gateway/context"
	"net/http"
)

func WrapIAMErrorResponseForRequest(status int, r *http.Request, code, resource string) S3Responser {
	message, ok := ErrorMessage(code)
	if !ok {
		panic(fmt.Sprintf("invalid error code: %s", code))
	}

	return NewIAMErrorResponse(
		status,
		code,
		message,
		resource,
		context.Get(r, "req_id").(string),
	)
}
