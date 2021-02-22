package main

import (
	groupcore "chores-suck/core/groups"
	"chores-suck/core/storage/postgres"
	"chores-suck/core/users"
	"chores-suck/web"
	"chores-suck/web/auth"
	"chores-suck/web/groups"
	"chores-suck/web/sessions"
	"chores-suck/web/views"
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
	auth := auth.NewService(users, store)
	views := views.NewService(store, repo)
	groups := groups.NewService(groupCore, users)
	handler := web.Handler(web.NewServices(auth, views, groups))
	log.Fatal(http.ListenAndServe(":8080", context.ClearHandler(handler)))
}
