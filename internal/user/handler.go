package user

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"super_web_app/internal/handlers"
	"super_web_app/pkg/logging"
)

type handler struct {
	logger *logging.Logger
}

func NewHandler(logger *logging.Logger) handlers.Handler {
	return &handler{logger: logger}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET("/users", h.GetUsersList)
}

func (h *handler) GetUsersList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("this is a list of users"))
}
