package tools

import (
	"net/http"
	"server/tools/database"
)

/********************************************************
TYPES
********************************************************/

type AuthHandler struct {
	conn *db.Connection
}

/********************************************************
INITIALIZER
********************************************************/

func NewAuthHandler(c *db.Connection) *AuthHandler {
	return &AuthHandler{
		conn: c,
	}
}

/********************************************************
HTTP
********************************************************/
func (h AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
