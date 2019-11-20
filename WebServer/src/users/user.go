package user

import (
	//"golang.org/x/crypto/bcrypt"
	"auth"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
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
	//w.Header().Set("Content-Type", "application/json")
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	ok, err, claims := auth.AuthToken(c.Value)
	if !ok || err == jwt.ErrSignatureInvalid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))

	//json.NewEncoder(w).Encode(names)
}

func (h UserHandler) handlePOST(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UserHandler handlePost")
	var newUser User
	json.NewDecoder(r.Body).Decode(&newUser)
	Users[newUser.Name] = newUser.Password
	w.Header().Set("Status", "201")
}
