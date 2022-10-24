package main

import (
	"context"
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
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "kek")
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "there will be a login page")
}

func takeUserWithId1(w http.ResponseWriter, r *http.Request) {
	var username, password string
	err := conn.QueryRow(context.Background(), "select username, password from users where id=$1", 1).Scan(&username, &password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(w, username+" "+password)
}

func handleRequest() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/login/", loginPage)
	http.HandleFunc("/take_user_with_id_1/", takeUserWithId1)
	http.ListenAndServe(":3000", nil)
}

func main() {
	connectToDatabase()
	fmt.Println("Connected to database")
	handleRequest()

}
