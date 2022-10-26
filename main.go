package main

import (
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"super_web_app/internal/user"
	"time"
)

var conn *pgx.Conn

func main() {
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
	log.Fatalln(server.Serve(listener))

	// conn = web.ConnectToDatabase()

}
