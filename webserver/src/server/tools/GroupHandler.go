package tools

import db "server/tools/database"

type GroupHandler struct {
	conn *db.Connection
	auth *AuthHandler
}

func NewGroupHandler(c *db.Connection, a *AuthHandler) GroupHandler {
	return GroupHandler{
		conn: c,
		auth, a,
	}
}

func (h *GroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.String() {
	case "/group/create":
		h.CreateGroup(w,r)
	case "/group/update":
		h.UpdateGroup(w,r)
	case "/group/find":
		h.FindGroup(w,r)
	case "/group/join":
		h.JoinGroup(w,r)
	}
}

func 
