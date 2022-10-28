package main

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"super_web_app/internal/config"
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
	cfg := config.GetConfig()

	handler := user.NewHandler(logger)
	handler.Register(router)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindApi, cfg.Listen.Port))

	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Infof("server is listening port %s:%s", cfg.Listen.BindApi, cfg.Listen.Port)
	logger.Fatal(server.Serve(listener))

}
