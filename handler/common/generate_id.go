package common

import (
	"fmt"
	"net/http"
)

func GenerateRequestIdHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("enter generate request id handler\n")
	// w.Write([]byte("generate request id handler!\n"))
}
