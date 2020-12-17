package main

import (
	"fmt"
	"log"
	"net/http"

	"chores-suck/users"
)

// ServeHTTP Exported Function for HTTP Server
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	u := users.User{Name: "Shannon"}
	fmt.Print(u, "\n")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
