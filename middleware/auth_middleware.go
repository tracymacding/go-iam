package middleware

import (
	"fmt"
	"net/http"
)

type AuthMiddleware struct {
}

func (m *AuthMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	fmt.Printf("enter auth middleware\n")
	next(w, req)
}
