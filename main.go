package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	"os"
	"super_web_app/utils"
	"time"
)

var conn *pgx.Conn

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func newClient(ctx context.Context, username, password, host, port, database string) (*pgx.Conn, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, host, port, database)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_conn, err := pgx.Connect(ctx, dsn)
	return _conn, err
}

func connectToDatabase() {

	err := utils.DoWithTries(func() error {
		_conn, err := newClient(context.Background(),
			"api",
			"8791",
			"db.physphile.ru",
			"5432",
			"api",
		)

		if err != nil {
			return err
		}

		conn = _conn

		return nil
	}, 5, 5*time.Second)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connected to database")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(&w)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<em>This is api for <a href=\"http://physphile.ru\">physphile.ru</a> site</em>")
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	connectToDatabase()
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
	connectToDatabase()
}

func handleRequest() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/login/", loginPage)
	http.HandleFunc("/register/", registerPage)
	http.ListenAndServe(":3000", nil)
}

func main() {
	connectToDatabase()
	handleRequest()
}
