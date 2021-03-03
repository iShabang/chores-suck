package web

import (
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"chores-suck/core"
	"chores-suck/web/messages"
)

type authBasicHandle func(http.ResponseWriter, *http.Request, uint64)
type authParamHandle func(http.ResponseWriter, *http.Request, httprouter.Params, uint64)

// Services holds references to services that handlers utilize to carry out requests
type Services struct {
	auth   AuthService
	views  ViewService
	groups GroupService
	users  core.UserService
}

// NewServices creates a new Services object
func NewServices(a AuthService, v ViewService, g GroupService, u core.UserService) *Services {
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
	ro.GET("/editgroup/:id", s.authorizeParam(s.editGroupGet))
	ro.POST("/editgroup/:id", s.authorizeParam(s.groupAccess(s.editGroupPost)))
	ro.HandlerFunc("GET", "/login", s.loginGet)
	ro.HandlerFunc("POST", "/login", s.loginPost)
	ro.HandlerFunc("GET", "/logout", s.logout)
	ro.HandlerFunc("POST", "/createuser", s.createUser)
	ro.HandlerFunc("GET", "/dashboard", s.authorize(s.dashboard))
	ro.HandlerFunc("GET", "/register", s.register)
	ro.HandlerFunc("GET", "/creategroup", s.authorize(s.createGroupGet))
	ro.HandlerFunc("POST", "/creategroup", s.authorize(s.createGroupPost))
	ro.HandlerFunc("GET", "/", s.index)
	return ro
}

func (s *Services) index(wr http.ResponseWriter, req *http.Request) {
	err := s.views.Index(wr, req)
	if err != nil {
		handleError(err, wr)
	}
}

func (s *Services) dashboard(wr http.ResponseWriter, req *http.Request, uid uint64) {
	err := s.views.BuildDashboard(wr, req, uid)
	if err != nil {
		handleError(err, wr)
	}
}

func (s *Services) logout(wr http.ResponseWriter, req *http.Request) {
	err := s.auth.Logout(wr, req)
	if err != nil {
		handleError(err, wr)
	}
	http.Redirect(wr, req, "/", 302)
}

func (s *Services) loginPost(wr http.ResponseWriter, req *http.Request) {
	err := s.auth.Login(wr, req)
	if err != nil {
		// TODO: Resend login form with error message
		handleError(err, wr)
		return
	}
	// TODO: Check cookie for a redirect url
	http.Redirect(wr, req, "/dashboard", 302)
}

func (s *Services) loginGet(wr http.ResponseWriter, req *http.Request) {
	_, err := s.auth.Authorize(wr, req)
	if err == nil {
		http.Redirect(wr, req, "/dashboard", 302)
	} else {
		err = s.views.LoginForm(wr, req)
		if err != nil {
			handleError(err, wr)
		}
	}
}

func (s *Services) register(wr http.ResponseWriter, req *http.Request) {
	msg := messages.RegisterMessage{}
	err := s.views.RegisterForm(wr, req, &msg)
	if err != nil {
		handleError(err, wr)
	}
}

func (s *Services) createUser(wr http.ResponseWriter, req *http.Request) {
	msg := messages.RegisterMessage{}
	err := s.auth.Create(wr, req, &msg)
	if err != nil {
		if err == ErrInvalidInput {
			s.views.RegisterForm(wr, req, &msg)
		} else {
			handleError(err, wr)
		}
		return
	}
	http.Redirect(wr, req, "/login", 302)
}

func (s *Services) createGroupGet(wr http.ResponseWriter, req *http.Request, uid uint64) {
	e := s.views.NewGroupForm(wr, req, uid, &messages.CreateGroup{})
	if e != nil {
		log.Printf("handler.createGroupGet: %s", e.Error())
		handleError(e, wr)
	}
}

func (s *Services) createGroupPost(wr http.ResponseWriter, req *http.Request, uid uint64) {
	msg := messages.CreateGroup{}
	e := s.groups.CreateGroup(wr, req, uid, &msg)
	if e != nil {
		log.Printf("handler.createGroup: %s", e.Error())
		msg.General = "There was an error creating a new group"
		s.views.NewGroupForm(wr, req, uid, &msg)
		return
	}

	// TODO: if successfull, redirect to the group page
	http.Redirect(wr, req, "/dashboard", 302)
}

func (s *Services) editGroupGet(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, uid uint64) {
	groupID, e := strconv.ParseUint(ps.ByName("id"), 10, 64)
	if e != nil {
		http.Error(wr, e.Error(), http.StatusInternalServerError)
	}
	group := core.Group{ID: groupID}
	e = s.groups.GetGroup(&group)

	edit, e := s.groups.CanEdit(&group, uid)
	if e != nil {
		handleError(e, wr)
		return
	} else if !edit {
		http.Error(wr, "You do not have permission to edit this group.", http.StatusUnauthorized)
		return
	}

	user := core.User{ID: uid}
	e = s.users.GetUserByID(&user)
	if e != nil {
		http.Error(wr, e.Error(), http.StatusInternalServerError)
		return
	}
	e = s.views.EditGroupForm(wr, req, &group, &user)
	if e != nil {
		handleError(e, wr)
	}
}

func (s *Services) editGroupPost(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, uid uint64, group *core.Group) {
	user := core.User{ID: uid}
	e := s.users.GetUserByID(&user)
	if e != nil {
		http.Error(wr, e.Error(), http.StatusInternalServerError)
		return
	}
	editGroup := false
	for _, r := range group.Roles {
		for _, u := range r.Users {
			if u.ID == user.ID && r.Can(core.EditGroup) {
				editGroup = true
				break
			}
		}
	}
	if !editGroup {
		http.Error(wr, "You do not have permission to edit this group", http.StatusUnauthorized)
		return
	}
	groupName := req.PostFormValue("groupname")
	// TODO: validate input
	group.Name = groupName
	e = s.groups.UpdateGroup(group)
	if e != nil {
		log.Print(e.Error())
	}
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

func (s *Services) groupAccess(handler func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, uid uint64, group *core.Group)) authParamHandle {
	return func(wr http.ResponseWriter, req *http.Request, ps httprouter.Params, uid uint64) {
		groupID, e := strconv.ParseUint(ps.ByName("id"), 10, 64)
		if e != nil {
			http.Error(wr, e.Error(), http.StatusInternalServerError)
		}
		group := core.Group{ID: groupID}
		e = s.groups.GetGroup(&group)

		edit, e := s.groups.CanEdit(&group, uid)
		if e != nil {
			handleError(e, wr)
			return
		} else if !edit {
			http.Error(wr, "You do not have permission to edit this group.", http.StatusUnauthorized)
			return
		}

		handler(wr, req, ps, uid, &group)
	}
}

/////////////////////////////////////////////////////////////////
// Helper methods
/////////////////////////////////////////////////////////////////

func handleError(err error, wr http.ResponseWriter) {
	if err != nil {
		switch e := err.(type) {
		case HttpError:
			http.Error(wr, e.Error(), e.Status())
		default:
			http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
