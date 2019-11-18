package login

import (
	//"golang.org/x/crypto/bcrypt"
	//"github.com/dgrijalva/jwt-go"
	//"encoding/json"
	//"fmt"
	"net/http"
	"users"
)

func (h LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.handleGET(w, r)
	case "POST":
		h.handlePOST(w, r)
	default:
		http.Error(w, "Invlalid login command", http.StatusMethodNotAllowed)
	}

}

type LoginHandler struct {
	users []user.User // Are slices always passed by reference?
}

// FYI: This is how you do dependency injection in Go
func NewLogin(u []user.User) *LoginHandler {
	return &LoginHandler{
		users: u,
	}
}

func (h LoginHandler) handleGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Status", "200")
}

func (h LoginHandler) handlePOST(w http.ResponseWriter, r *http.Request) {
}
