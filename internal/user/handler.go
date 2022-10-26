package user

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"super_web_app/internal/handlers"
)

type handler struct {
}

func NewHandler() handlers.Handler {
	return &handler{}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET("/users", h.GetUsersList)
}

func (h *handler) GetUsersList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("this is a list of users"))
}
