package tools

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type RegisterHandler struct {
	c *Connection
}

func (h RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.handlePOST(w, r)
	default:
		http.Error(w, "Invlalid login command", http.StatusMethodNotAllowed)
	}
}

func (h RegisterHandler) handlePOST(w http.ResponseWriter, r *http.Request) {
	var nu UserLarge
	json.NewDecoder(r.Body).Decode(&nu)
	if nu.FirstName == "" ||
	   nu.LastName == "" ||
	   nu.Email == "" ||
	   nu.Password == "" ||
	   nu.Username == "" || {

	   }

}
