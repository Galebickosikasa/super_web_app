package main

import (
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"super_web_app/internal/user"
	"super_web_app/internal/web"
	"super_web_app/pkg/logging"
	"time"
)

var conn *pgx.Conn

func main() {
	logger := logging.GetLogger()
	logger.Info("start")
	conn = web.ConnectToDatabase()

	router := httprouter.New()
	handler := user.NewHandler()
	handler.Register(router)

	listener, err := net.Listen("tcp", ":3000")

	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Fatal(server.Serve(listener))

}
