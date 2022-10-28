package web

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"os"
	"super_web_app/internal/config"
	"super_web_app/internal/utils"
	"super_web_app/pkg/logging"
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
	logger := logging.GetLogger()
	logger.Info("Trying to connect to database")

	cfg := config.GetConfig()

	err := utils.DoWithTries(func() error {
		_conn, err := newClient(context.Background(),
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Database,
		)

		if err != nil {
			return err
		}

		conn = _conn

		return nil
	}, 5, 5*time.Second)

	if err != nil {
		logger.Errorf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	logger.Info("Connected to database")
	return
}
