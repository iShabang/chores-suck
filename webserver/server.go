package main

import (
	"fmserver/tools"
	"fmserver/tools/database"
	"fmt"
	"log"
	"net/http"
)

type App struct {
	LoginHandler    *tools.LoginHandler
	RegisterHandler *tools.RegisterHandler
	AuthHandler     *tools.AuthHandler
	ChoreHandler    *tools.ChoreHandler
	FileHandler     *tools.FileHandler
}

func (h App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Print("serving http\n")
	fmt.Println(r.URL.Path)
	switch r.URL.Path {
	case "/login":
		h.LoginHandler.Login(w, r)
	case "/logout":
		h.LoginHandler.Logout(w, r)
	case "/register":
		h.RegisterHandler.ServeHTTP(w, r)
	default:
		h.FileHandler.ServeFile(w, r)
	}
}

func main() {
	conn := db.NewConnection()
	err := conn.Connect("mongodb://127.0.0.1:27017")
	if err != nil {
		log.Print(err)
	}

	var app App
	app.LoginHandler = tools.NewLogin(&conn)
	app.RegisterHandler = tools.NewRegister(&conn)
	app.AuthHandler = tools.NewAuthHandler(&conn)
	app.ChoreHandler = tools.NewChoreHandler(&conn, app.AuthHandler)
	app.FileHandler = tools.NewFileHandler(&conn, app.AuthHandler, "files/")
	log.Fatal(http.ListenAndServe(":8080", app))
}
