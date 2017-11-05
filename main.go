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
	"github.com/go-iam/handler/group"
	"github.com/go-iam/handler/key"
	"github.com/go-iam/handler/policy"
	"github.com/go-iam/handler/user"
	"github.com/go-iam/middleware"
	"github.com/go-iam/mux"
	"github.com/golang/glog"
)

func setupRouter() *mux.Router {
	router := mux.NewRouter()
	router.HookFunc(mux.HookAfterRouter, common.LogHandler).HookFunc(mux.HookAfterRouter, common.SendResponseHandler)

	router.HandleFunc("/", account.CreateAccountHandler).Methods("GET").Queries("Action", "CreateAccount")
	router.HandleFunc("/", account.GetAccountHandler).Methods("GET").Queries("Action", "GetAccount")
	router.HandleFunc("/", account.DeleteAccountHandler).Methods("GET").Queries("Action", "DeleteAccount")
	router.HandleFunc("/", account.UpdateAccountHandler).Methods("GET").Queries("Action", "UpdateAccount")
	router.HandleFunc("/", account.ListAccountHandler).Methods("GET").Queries("Action", "ListAccount")
	router.HandleFunc("/", account.LoginAccountHandler).Methods("GET").Queries("Action", "LoginAccount")

	router.HandleFunc("/", user.CreateUserHandler).Methods("GET").Queries("Action", "CreateIamUser")
	router.HandleFunc("/", user.GetIAMUserHandler).Methods("GET").Queries("Action", "GetIamUser")
	router.HandleFunc("/", user.DeleteIAMUserHandler).Methods("GET").Queries("Action", "DeleteIamUser")
	router.HandleFunc("/", user.UpdateIAMUserHandler).Methods("GET").Queries("Action", "UpdateIamUser")
	router.HandleFunc("/", user.ListIAMUserHandler).Methods("GET").Queries("Action", "ListIamUser")

	router.HandleFunc("/", group.CreateGroupHandler).Methods("GET").Queries("Action", "CreateGroup")
	router.HandleFunc("/", group.GetGroupHandler).Methods("GET").Queries("Action", "GetGroup")
	router.HandleFunc("/", group.DeleteGroupHandler).Methods("GET").Queries("Action", "DeleteGroup")
	router.HandleFunc("/", group.UpdateGroupHandler).Methods("GET").Queries("Action", "UpdateGroup")
	router.HandleFunc("/", group.ListGroupHandler).Methods("GET").Queries("Action", "ListGroup")

	router.HandleFunc("/", group.GroupAddUserHandler).Methods("GET").Queries("Action", "AddUserToGroup")
	router.HandleFunc("/", group.GroupRemoveUserHandler).Methods("GET").Queries("Action", "RemoveUserFromGroup")
	router.HandleFunc("/", group.ListUserGroupHandler).Methods("GET").Queries("Action", "ListGroupsForUser")
	router.HandleFunc("/", group.ListGroupUserHandler).Methods("GET").Queries("Action", "ListUsersForGroup")

	router.HandleFunc("/", policy.CreatePolicyHandler).Methods("GET").Queries("Action", "CreatePolicy")
	router.HandleFunc("/", policy.GetPolicyHandler).Methods("GET").Queries("Action", "GetPolicy")
	router.HandleFunc("/", policy.DeletePolicyHandler).Methods("GET").Queries("Action", "DeletePolicy")
	router.HandleFunc("/", policy.UpdatePolicyHandler).Methods("GET").Queries("Action", "UpdatePolicy")
	router.HandleFunc("/", policy.ListPolicyHandler).Methods("GET").Queries("Action", "ListPolicy")

	router.HandleFunc("/", key.CreateKeyHandler).Methods("GET").Queries("Action", "CreateAccessKey")
	router.HandleFunc("/", key.GetKeyHandler).Methods("GET").Queries("Action", "GetAccessKey")
	router.HandleFunc("/", key.DeleteKeyHandler).Methods("GET").Queries("Action", "DeleteAccessKey")
	router.HandleFunc("/", key.UpdateKeyHandler).Methods("GET").Queries("Action", "UpdateAccessKey")
	router.HandleFunc("/", key.ListKeyHandler).Methods("GET").Queries("Action", "ListAccessKey")

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
