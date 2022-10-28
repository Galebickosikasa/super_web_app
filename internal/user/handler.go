package user

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"super_web_app/internal/handlers"
	"super_web_app/pkg/logging"
)

type handler struct {
	logger   *logging.Logger
	database *pgx.Conn
}

func NewHandler(logger *logging.Logger, database *pgx.Conn) handlers.Handler {
	return &handler{logger: logger, database: database}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET("/users", h.GetUsersList)
	router.GET("/users/:uuid", h.GetUserById)
}

func (h *handler) GetUsersList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	query, err := h.database.Query(context.Background(), "select username, password from users")
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(400)
	} else {
		var username, password string
		var res string
		for true {
			if query.Next() {
				err := query.Scan(&username, &password)
				if err != nil {
					continue
				}
				res += fmt.Sprintf("username: %s, password: %s\n", username, password)
			} else {
				break
			}
		}
		w.WriteHeader(200)
		w.Write([]byte(res))
	}
}

func (h *handler) GetUserById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	uuid := params.ByName("uuid")
	var username, password string
	err := h.database.QueryRow(context.Background(), "select username, password from users where id = $1", uuid).Scan(&username, &password)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(400)
	} else {
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("username: %s, password: %s", username, password)))
	}
}
