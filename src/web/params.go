package web

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Params struct {
	Writer  http.ResponseWriter
	Request *http.Request
	UserID  uint64
	Query   httprouter.Params
}
