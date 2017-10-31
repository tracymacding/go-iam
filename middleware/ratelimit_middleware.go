package middleware

import (
	"fmt"
	"net/http"
)

type RateLimitMiddleware struct {
}

func (m *RateLimitMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	fmt.Printf("enter ratelimit middleware\n")
	next(w, req)
}
