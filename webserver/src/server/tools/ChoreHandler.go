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
	if !success {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var err error
	var chores []*db.Chore
	chores, err = h.conn.GetUserChores(id)
	if !((err == nil) && (len(chores) > 0)) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var cj []byte
	cj, err = json.Marshal(chores)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(cj)
}

func (h *ChoreHandler) ChoresGroup(w http.ResponseWriter, r *http.Request) {
}
