package main

import (
	"chores-suck/rest"
	"chores-suck/rest/auth"
	"chores-suck/rest/sessions"
	"chores-suck/rest/views"
	"chores-suck/storage/postgres"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/context"
)

func main() {
	repo := postgres.NewStorage()
	store := sessions.NewStore(repo, []byte(os.Getenv("SESSION_KEY")))

	auth := auth.NewService(repo, store)
	views := views.NewService(store, repo)
	handler := rest.Handler(rest.NewServices(auth, views))
	log.Fatal(http.ListenAndServe(":8080", context.ClearHandler(handler)))
}
