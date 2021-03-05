package web

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type authBasicHandle func(http.ResponseWriter, *http.Request, uint64)
type authParamHandle func(http.ResponseWriter, *http.Request, httprouter.Params, uint64)

// Services holds references to services that handlers utilize to carry out requests
type Services struct {
	auth   AuthService
	views  ViewService
	groups GroupService
	users  UserService
}

// NewServices creates a new Services object
func NewServices(a AuthService, v ViewService, g GroupService, u UserService) *Services {
	return &Services{
		auth:   a,
		views:  v,
		groups: g,
		users:  u,
	}
}

//Handler creates and returns a new http.Handler with the request handlers and functions pre-registered/routed
func Handler(s *Services) http.Handler {
	ro := httprouter.New()
	ro.GET("/groups/:groupID", s.authorizeParam(s.views.EditGroupForm))
	ro.POST("/groups/:groupID", s.authorizeParam(s.groups.GroupAccess(s.groups.EditGroup)))
	ro.POST("/groups/:groupID/remove/user/:userID", s.authorizeParam(s.groups.GroupAccess(s.groups.DeleteMember)))
	ro.POST("/groups/:groupID/add/user", s.authorizeParam(s.groups.GroupAccess(s.groups.AddMember)))
	ro.GET("/groups/:groupID/add/role", s.authorizeParam(s.groups.GroupAccess(s.views.NewRoleForm)))
	ro.POST("/groups/:groupID/add/role", s.authorizeParam(s.groups.GroupAccess(s.groups.AddRole)))
	ro.GET("/groups/:groupID/update/role/:roleID", s.authorizeParam(s.groups.GroupAccess(s.views.UpdateRoleForm)))
	ro.HandlerFunc("GET", "/login", s.views.LoginForm)
	ro.HandlerFunc("POST", "/login", s.auth.Login)
	ro.HandlerFunc("GET", "/logout", s.auth.Logout)
	ro.HandlerFunc("POST", "/new/user", s.users.CreateUser)
	ro.HandlerFunc("GET", "/dashboard", s.authorize(s.views.Dashboard))
	ro.HandlerFunc("GET", "/new/user", s.views.RegisterForm)
	ro.HandlerFunc("GET", "/new/group", s.authorize(s.views.NewGroupForm))
	ro.HandlerFunc("POST", "/new/group", s.authorize(s.groups.CreateGroup))
	ro.HandlerFunc("GET", "/", s.views.Index)
	return ro
}

/////////////////////////////////////////////////////////////////
// Middleware methods
/////////////////////////////////////////////////////////////////
func (s *Services) authorize(handler func(wr http.ResponseWriter, req *http.Request, uid uint64)) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		// TODO: Save the requested url in a cookie that can be redirected to after logging in successfully
		uid, err := s.auth.Authorize(wr, req)
		if err != nil {
			log.Print(err)
			http.Redirect(wr, req, "/login", 302)
			return
		}

		handler(wr, req, uid)
	}
}

func (s *Services) authorizeParam(handler func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, uid uint64)) httprouter.Handle {
	return func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		uid, err := s.auth.Authorize(wr, req)
		if err != nil {
			log.Print(err)
			http.Redirect(wr, req, "/login", 302)
			return
		}

		handler(wr, req, ps, uid)
	}
}
