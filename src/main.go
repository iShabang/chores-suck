package main

import (
	"log"
	"net/http"

	"github.com/gorilla/context"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux)))
}
