package tools

import (
	"net/http"
	"server/tools/database"
)

/********************************************************
TYPES
********************************************************/
type FileHandler struct {
	conn    *db.Connection
	auth    *AuthHandler
	fileDir string
}

/********************************************************
INITIALIZER
********************************************************/
func NewFileHandler(c *db.Connection, a *AuthHandler, f string) *FileHandler {
	return &FileHandler{
		conn:    c,
		auth:    a,
		fileDir: f,
	}
}

/********************************************************
EXPORTED METHODS
********************************************************/
func (h *FileHandler) ServeFile(w http.ResponseWriter, r *http.Request) {
	// Check for publicly accessible files
	public := r.URL.Path == "/login.html" || r.URL.Path == "/register.html"
	if public {
		http.ServeFile(w, r, h.fileDir+r.URL.Path[1:])
	}
	_, success := h.auth.AuthorizeRequest(r)
	if !success {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		path := r.URL.Path
		if path == "/" {
			path = h.fileDir + "index.html"
		} else {
			path = h.fileDir + path[1:]
		}
		http.ServeFile(w, r, path)
	}
}
