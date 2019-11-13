package main

import (
	//"golang.org/x/crypto/bcrypt"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

type App struct {
	UserHandler *UserHandler
}

type User struct {
	Name     string
	Password string
}

type UserHandler struct{}

var (
	users []User
)

func (h App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	if r.URL.Path == "/users" {
		fmt.Println("got users")
		h.UserHandler.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Not Found", http.StatusNotFound)
}

func (h UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UserHandler ServeHTTP")
	fmt.Println(r.Method)
	switch r.Method {
	case "GET":
		h.handleGet(w, r)
	case "POST":
		h.handlePOST(w, r)
	default:
		http.Error(w, "Invlalid user command", http.StatusMethodNotAllowed)
	}

}

func (h UserHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UserHandler handleGet")
	w.Header().Set("Content-Type", "application/json")
	names := make([]string, len(users))
	for i, v := range users {
		names[i] = v.Name
	}
	json.NewEncoder(w).Encode(names)
}

func (h UserHandler) handlePOST(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UserHandler handlePost")
	var newUser User
	json.NewDecoder(r.Body).Decode(&newUser)
	users = append(users, newUser)
	w.Header().Set("Status", "201")
}

func main() {
	var app App
	var userHandler UserHandler
	app.UserHandler = &userHandler
	users = []User{
		{"Shannon", "password1"},
		{"Bee", "password2"},
	}
	http.ListenAndServe(":8080", app)
}
