package login

import (
	//"golang.org/x/crypto/bcrypt"
	"auth"
	"encoding/json"
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
	users map[string]string // Are slices always passed by reference?
}

// FYI: This is how you do dependency injection in Go
func NewLogin(u map[string]string) *LoginHandler {
	return &LoginHandler{
		users: u,
	}
}

func (h LoginHandler) handleGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Status", "200")
}

func (h LoginHandler) handlePOST(w http.ResponseWriter, r *http.Request) {
	var newUser user.User
	json.NewDecoder(r.Body).Decode(&newUser)
	pass, ok := h.users[newUser.Name]
	if !ok || pass != newUser.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokenString, err, expireTime := auth.GenToken(newUser.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expireTime,
	})
}
