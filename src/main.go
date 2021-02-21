package main

import (
	"chores-suck/core/storage/postgres"
	"chores-suck/core/users"
	"chores-suck/rest"
	"chores-suck/rest/auth"
	"chores-suck/rest/sessions"
	"chores-suck/rest/views"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/context"
)

func main() {
	repo := postgres.NewStorage()
	store := sessions.NewStore(repo, []byte(os.Getenv("SESSION_KEY")))
	users := users.NewService(repo)

	auth := auth.NewService(users, store)
	views := views.NewService(store, repo)
	handler := rest.Handler(rest.NewServices(auth, views))
	log.Fatal(http.ListenAndServe(":8080", context.ClearHandler(handler)))
}
