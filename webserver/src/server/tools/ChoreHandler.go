package tools

import (
	"encoding/json"
	"net/http"
	"server/tools/database"
)

/********************************************************
TYPES
********************************************************/
type ChoreHandler struct {
	conn *db.Connection
	auth *AuthHandler
}

/********************************************************
INITIALIZER
********************************************************/
func NewChoreHandler(conn *db.Connection, auth *AuthHandler) *ChoreHandler {
	return &ChoreHandler{
		conn: conn,
		auth: auth,
	}
}

/********************************************************
EXPORTED METHODS
********************************************************/
func (h *ChoreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.String() {
	case "/chores/user":
		h.ChoresUser(w, r)
	case "/chores/group":
		h.ChoresGroup(w, r)
	}
}
func (h *ChoreHandler) ChoresUser(w http.ResponseWriter, r *http.Request) {
	id, success := h.auth.AuthorizeRequest(r)

	var err error
	var chores []*db.Chore
	if success {
		chores, err = h.conn.GetUserChores(id)
		success = (err == nil) && (len(chores) > 0)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var cj []byte
	if success {
		cj, err = json.Marshal(chores)
		success = (err == nil)
	}

	if success {
		w.Header().Set("content-type", "application/json")
		w.Write(cj)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *ChoreHandler) ChoresGroup(w http.ResponseWriter, r *http.Request) {
}
