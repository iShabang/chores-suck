package tools

import (
	"encoding/json"
	db "fmserver/tools/database"
	"fmt"
	"net/http"
)

type GroupName struct {
	String string `json:"name"`
}

type GroupHandler struct {
	conn *db.Connection
	auth *AuthHandler
}

func NewGroupHandler(c *db.Connection, a *AuthHandler) GroupHandler {
	return GroupHandler{
		conn: c,
		auth: a,
	}
}

func (h *GroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.String() {
	case "/group/create":
		h.CreateGroup(w, r)
	case "/group/update":
		h.UpdateGroup(w, r)
	case "/group/find":
		h.FindGroup(w, r)
	case "/group/join":
		h.JoinGroup(w, r)
	}
}

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	id, success := h.auth.AuthorizeRequest(r)
	if !success {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var name GroupName
	json.NewDecoder(r.Body).Decode(&name)

	group := db.Group{
		Name:  name.String,
		Admin: id,
	}
	_, err := h.conn.AddGroup(&group)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *GroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	id, success := h.auth.AuthorizeRequest(r)
	if !success {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var name GroupName
	json.NewDecoder(r.Body).Decode(&name)

}

func (h *GroupHandler) FindGroup(w http.ResponseWriter, r *http.Request) {
	id, success := h.auth.AuthorizeRequest(r)
	if !success {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var name GroupName
	json.NewDecoder(r.Body).Decode(&name)

	group, err := h.conn.FindUserGroup(id, name.String)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *GroupHandler) JoinGroup(w http.ResponseWriter, r *http.Request) {}
