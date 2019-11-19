package user

import (
	//"golang.org/x/crypto/bcrypt"
	//"github.com/dgrijalva/jwt-go"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	Users = map[string]string{}
)

type User struct {
	Name     string
	Password string
}

type UserHandler struct{}

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
	var names []string
	for i, _ := range Users {
		names = append(names, i)
	}
	json.NewEncoder(w).Encode(names)
}

func (h UserHandler) handlePOST(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UserHandler handlePost")
	var newUser User
	json.NewDecoder(r.Body).Decode(&newUser)
	Users[newUser.Name] = newUser.Password
	w.Header().Set("Status", "201")
}
