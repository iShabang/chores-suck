/*
TODO: Add a login method that creates and returns a jason web token after successfull login. The key will
need to be generated separately from the jwt-go package.
*/
package main

import (
	//"golang.org/x/crypto/bcrypt"
	//"github.com/dgrijalva/jwt-go"
	"fmt"
	"login"
	"net/http"
	"users"
)

type App struct {
	UserHandler  *user.UserHandler
	LoginHandler *login.LoginHandler
}

func (h App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	switch r.URL.Path {
	case "/users":
		h.UserHandler.ServeHTTP(w, r)
	case "/login":
		h.LoginHandler.ServeHTTP(w, r)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func main() {
	var app App
	var userHandler user.UserHandler
	var loginHandler login.LoginHandler
	app.UserHandler = &userHandler
	app.LoginHandler = &loginHandler
	user.Users = []user.User{
		{Name: "Shannon", Password: "password1"},
		{Name: "Bee", Password: "password2"},
	}
	http.ListenAndServe(":8080", app)
}
