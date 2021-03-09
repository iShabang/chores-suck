package main

import (
	"chores-suck/core"
	"chores-suck/core/storage/postgres"
	"chores-suck/web"
	"chores-suck/web/sessions"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/context"
)

func main() {
	repo := postgres.NewStorage()
	userCore := core.NewUserService(repo)
	groupCore := core.NewGroupService(repo)
	roleCore := core.NewRoleService(repo)

	store := sessions.NewStore(repo, []byte(os.Getenv("SESSION_KEY")))
	auth := web.NewAuthService(userCore, store)
	views := web.NewViewService(store, userCore, auth, groupCore)
	users := web.NewUserService(userCore, views)
	groups := web.NewGroupService(groupCore, userCore, views)
	roles := web.NewRoleService(groupCore, roleCore, userCore, views)
	handler := web.Handler(web.NewServices(auth, views, groups, users, roles))
	log.Fatal(http.ListenAndServe(":8080", context.ClearHandler(handler)))
}
