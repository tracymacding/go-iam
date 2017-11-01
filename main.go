package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"

	"github.com/go-iam/db"
	_ "github.com/go-iam/db/mongodb"
	_ "github.com/go-iam/db/mysql"
	"github.com/go-iam/handler/account"
	"github.com/go-iam/handler/common"
	"github.com/go-iam/handler/user"
	"github.com/go-iam/middleware"
	"github.com/go-iam/mux"
	"github.com/golang/glog"
)

func setupRouter() *mux.Router {
	router := mux.NewRouter()
	router.HookFunc(mux.HookAfterRouter, common.LogHandler).HookFunc(mux.HookAfterRouter, common.SendResponseHandler)

	router.HandleFunc("/", user.CreateUserHandler).Methods("GET").Queries("Action", "CreateIamUser")
	router.HandleFunc("/", account.CreateAccountHandler).Methods("GET").Queries("Action", "CreateAccount")
	router.HandleFunc("/", account.GetAccountHandler).Methods("GET").Queries("Action", "GetAccount")
	router.HandleFunc("/", account.DeleteAccountHandler).Methods("GET").Queries("Action", "DeleteAccount")
	router.HandleFunc("/", account.UpdateAccountHandler).Methods("GET").Queries("Action", "UpdateAccount")
	router.HandleFunc("/", account.ListAccountHandler).Methods("GET").Queries("Action", "ListAccount")
	router.HandleFunc("/", account.LoginAccountHandler).Methods("GET").Queries("Action", "LoginAccount")
	return router
}

func startServe(router *mux.Router) {
	router.Use(&middleware.GenerateRequestIdMiddleware{})
	go func() {
		http.ListenAndServe(":3030", nil)
	}()
	http.ListenAndServe(*listenAddr, router)
}

func main() {
	flag.Parse()

	// 使用mongodb作为后端存储,集群地址为192.168.100.100
	err := db.Open("mongodb", *mongoDBAddr)
	// 使用mysql作为后端存储,集群地址为192.168.100.100
	// err := db.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", *dbUser, *dbPassword, *dbAddr, *dbName))
	if err != nil {
		glog.Fatalf("init database error: %s", err)
	}

	startServe(setupRouter())
}
