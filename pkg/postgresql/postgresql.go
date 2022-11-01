package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"os"
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

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func newClient(ctx context.Context, cfg Config) (*pgx.Conn, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_conn, err := pgx.Connect(ctx, dsn)
	return _conn, err
}

func ConnectToDatabase(cfg Config) (conn *pgx.Conn) {
	logger := logging.GetLogger()
	logger.Info("Trying to connect to database")

	err := utils.DoWithTries(func() error {
		_conn, err := newClient(context.Background(), cfg)

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
