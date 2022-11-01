package main

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"super_web_app/internal/config"
	"super_web_app/internal/user"
	"super_web_app/pkg/logging"
	"super_web_app/pkg/postgresql"
	"time"
)

var conn *pgx.Conn

func main() {
	logger := logging.GetLogger()
	logger.Info("start")

	cfg := config.GetConfig()
	conn = postgresql.ConnectToDatabase(postgresql.Config{
		Username: cfg.Database.Username,
		Password: cfg.Database.Password,
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Database: cfg.Database.Database,
	})

	router := httprouter.New()

	handler := user.NewHandler(logger, conn)
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
