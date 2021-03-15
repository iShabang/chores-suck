package web

import (
	"chores-suck/core"
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
	roles  RoleService
	chores ChoreService
}

// NewServices creates a new Services object
func NewServices(a AuthService, v ViewService, g GroupService, u UserService, r RoleService, c ChoreService) *Services {
	return &Services{
		auth:   a,
		views:  v,
		groups: g,
		users:  u,
		roles:  r,
		chores: c,
	}
}

//Handler creates and returns a new http.Handler with the request handlers and functions pre-registered/routed
func Handler(s *Services) http.Handler {
	ro := httprouter.New()
	ro.GET("/groups/update/:groupID", s.groupMW(s.views.EditGroupForm))
	ro.GET("/roles/create/:groupID", s.groupMW(s.views.NewRoleForm))
	ro.GET("/roles/update/:roleID", s.roleMW(s.views.UpdateRoleForm))
	ro.GET("/chores/create/:groupID", s.groupMW(s.views.NewChoreForm))
	ro.POST("/groups/update/:groupID", s.groupMW(s.groups.UpdateGroup))
	ro.POST("/roles/create/:groupID", s.groupMW(s.groups.AddRole))
	ro.POST("/roles/update/:roleID", s.roleMW(s.roles.Update))
	ro.POST("/chores/create/:groupID", s.groupMW(s.chores.Create))
	ro.HandlerFunc("GET", "/", s.views.Index)
	ro.HandlerFunc("GET", "/login", s.views.LoginForm)
	ro.HandlerFunc("GET", "/logout", s.auth.Logout)
	ro.HandlerFunc("GET", "/dashboard", s.authorize(s.views.Dashboard))
	ro.HandlerFunc("GET", "/register", s.views.RegisterForm)
	ro.HandlerFunc("GET", "/groups/create", s.authorize(s.views.NewGroupForm))
	ro.HandlerFunc("POST", "/login", s.auth.Login)
	ro.HandlerFunc("POST", "/register", s.users.CreateUser)
	ro.HandlerFunc("POST", "/groups/create", s.authorize(s.groups.CreateGroup))
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

func (s *Services) groupMW(handler func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, user *core.User, group *core.Group)) httprouter.Handle {
	return s.authorizeParam(s.groups.GroupAccess(handler))
}

func (s *Services) roleMW(handler func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, user *core.User, role *core.Role)) httprouter.Handle {
	return s.authorizeParam(s.roles.RoleMW(handler))
}
