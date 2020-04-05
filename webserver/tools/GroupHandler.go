/*
A Group is a group of users that share chores. The chores are added to the
group and then distributed among the users. When a chore is assigned to a
user, it must be associated with a group (via group-id).

Each group has an administrator (just defined by a user id) that controls who is in the group and the chores in it.
Each group will have an array of users (not including the administrator). The administrator will also be assigned chores.

Administrator functions:
Create a group - this is techinically not an "administrator" privelage since anyone can create a group, but it must be included in this handler. The user that creates the group is the default administrator.
Change the group name
Add/remove users from the group
Add/remove chores from the group

*/

package tools

import (
	"encoding/json"
	db "fmserver/tools/database"
	"fmt"
	"net/http"
)

GroupChore struct {
	Name
}

type GroupJson struct {
	Name string `json:"name"`
	Chores 
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
	}
}

/*
The CreateGroup method will be restricted to authorized users only. Any user that creates a group will be the default administrator.
*/
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	id, success := h.auth.AuthorizeRequest(r)
	if !success {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var gj GroupJson
	json.NewDecoder(r.Body).Decode(&gj)

	group := db.Group{
		Name:  gj.Name,
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

/* Administrators will update fields in a group using a form. We will need to
* get all the data in JSON format and commit the changes to the database. */
func (h *GroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	// First make sure this is a valid user
	id, success := h.auth.AuthorizeRequest(r)
	if !success {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// TODO: Grab all fields needed to update
	var name GroupName
	json.NewDecoder(r.Body).Decode(&name)

	// Before commiting these changes, we need to ensure that this user is an
	// administrator for this group. Find the group in the database and compare
	// the administrator id with the current user id

}
