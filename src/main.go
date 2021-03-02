package main

import (
	groupcore "chores-suck/core/groups"
	"chores-suck/core/storage/postgres"
	"chores-suck/core/users"
	"chores-suck/web"
	"chores-suck/web/sessions"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/context"
)

func main() {
	repo := postgres.NewStorage()
	users := users.NewService(repo)
	groupCore := groupcore.NewService(repo)

	store := sessions.NewStore(repo, []byte(os.Getenv("SESSION_KEY")))
	auth := web.NewAuthService(users, store)
	views := web.NewViewService(store, users)
	groups := web.NewGroupService(groupCore, users)
	handler := web.Handler(web.NewServices(auth, views, groups, users))
	log.Fatal(http.ListenAndServe(":8080", context.ClearHandler(handler)))
}
