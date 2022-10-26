package web

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"os"
	"super_web_app/internal/utils"
	"time"
)

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

func ConnectToDatabase() (conn *pgx.Conn) {
	fmt.Println("Trying to connect to database")

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
	return
}
