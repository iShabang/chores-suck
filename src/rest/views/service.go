package views

import (
	"chores-suck/rest/sessions"
	"chores-suck/types"
	"net/http"
)

// Service provides functionality for generating views
type Service interface {
	BuildDashboard(http.ResponseWriter, *http.Request, string) error
}

// Repository describes the interface necessary for grabbing data necessary for views
type Repository interface {
	GetUserByID(string) (types.User, error)
}

type service struct {
	store *sessions.Store
	repo  Repository
}

func (s *service) BuildDashboard(wr http.ResponseWriter, req *http.Request, uid string) error {
	// Get User
	// Get Memberships
	// Get Group Data
	// Populate template with data
	return nil
}
