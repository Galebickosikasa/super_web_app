package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"net/http"
	"super_web_app/utils"
	"super_web_app/web"
)

var conn *pgx.Conn

func homePage(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(&w)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<em>This is api for <a href=\"http://physphile.ru\">physphile.ru</a> site</em>")
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("go to login page")
	conn = web.ConnectToDatabase()
	utils.SetCORSHeaders(&w)

	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	type User struct {
		Username string
		Password string
	}
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var password string

	err = conn.QueryRow(context.Background(), "SELECT password FROM users WHERE username = $1", user.Username).Scan(&password)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "User does not exist", http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	if password != user.Password {
		http.Error(w, "Wrong password", http.StatusBadRequest)
		return
	}

	authToken := utils.GetToken()
	conn.Query(context.Background(), "UPDATE users SET auth_token = $1 WHERE username = $2", authToken, user.Username)
	type Response struct {
		AuthToken string
	}
	res, _ := json.Marshal(Response{authToken})

	fmt.Fprintf(w, "%s", string(res))
}

func registerPage(w http.ResponseWriter, r *http.Request) {
	conn = web.ConnectToDatabase()
}

func handleRequest() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/login/", loginPage)
	http.HandleFunc("/register/", registerPage)
	http.ListenAndServe(":3000", nil)
}

func main() {
	conn = web.ConnectToDatabase()
	handleRequest()
}
