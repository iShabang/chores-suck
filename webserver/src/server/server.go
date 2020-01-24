package main

import (
	//"golang.org/x/crypto/bcrypt"
	"fmt"
	"log"
	"net/http"
	"tools"
)

type App struct {
	LoginHandler    *tools.LoginHandler
	RegisterHandler *tools.RegisterHandler
	fileDir         string
}

func (h App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	switch r.URL.Path {
	case "/login":
		h.LoginHandler.ServeHTTP(w, r)
	case "/register":
		h.RegisterHandler.ServeHTTP(w, r)
	case "/":
		http.ServeFile(w, r, h.fileDir+"index.html")
	default:
		http.ServeFile(w, r, h.fileDir+r.URL.Path[1:])
	}
}

func main() {
	conn := tools.NewConnection()
	err := conn.Connect("mongodb://127.0.0.1:27017")
	if err != nil {
		log.Print(err)
	}

	var app App
	loginHandler := tools.NewLogin(&conn)
	regHandler := tools.NewRegister(&conn)
	app.LoginHandler = loginHandler
	app.RegisterHandler = regHandler
	app.fileDir = "./files/"
	log.Fatal(http.ListenAndServe(":8080", app))
}
