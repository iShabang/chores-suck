package login

import (
	//"golang.org/x/crypto/bcrypt"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
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

var (
	SECRET_KEY = []byte("cnBkyv93jqZ1DMWkDxHqCbfb@II*bq8!IUJnf#859VBz&n80$WQ9kIUEn5zOGz5M")
)

type LoginHandler struct {
	users map[string]string // Are slices always passed by reference?
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
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

	// insert JWT logic here
	expireTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: newUser.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SECRET_KEY)
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
