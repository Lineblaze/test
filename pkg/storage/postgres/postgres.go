package postgres

import (
	"avito_test/internal/config"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"log"
	"time"
)

const (
	_defaultConnectionAttempts = 10
	_defaultConnectionTimeout  = time.Second
	_maxConnections            = int32(800)
	_minConnections            = int32(50)
	_maxConnectionLifeTime     = time.Second * 300
	_maxIdleLifeTime           = time.Second * 15
)

type Postgres interface {
	Stats() *pgxpool.Stat
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	TxRunner
}

type Pool struct {
	db *pgxpool.Pool
}

func InitPsqlDB(c *config.Config) (Postgres, error) { // nolint: ireturn
	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.DBName,
		c.Postgres.SSLMode)
	connectionUrl += fmt.Sprintf(" pool_max_conns=%d pool_min_conns=%d pool_max_conn_lifetime=%v pool_max_conn_idle_time=%v",
		_maxConnections, _minConnections, _maxConnectionLifeTime, _maxIdleLifeTime)

	connectionAttempts := _defaultConnectionAttempts
	var result *pgxpool.Pool
	var err error

	for connectionAttempts > 0 {
		result, err = pgxpool.New(context.Background(), connectionUrl)
		if err == nil {
			break
		}

		log.Printf("ATTEMPT %d TO CONNECT TO POSTGRES BY URL %s FAILED: %s\n", connectionAttempts, connectionUrl, err.Error())
		connectionAttempts--
		time.Sleep(_defaultConnectionTimeout)
	}

	if result == nil {
		log.Printf("POSTGRES CONNECTION(%s) ERROR: %s\n", connectionUrl, err.Error())
		return nil, errors.WithMessage(err, "failed to initialize PostgreSQL connection")
	}

	return &Pool{db: result}, nil
}

func (p *Pool) Stats() *pgxpool.Stat {
	return p.db.Stat()
}

func (p *Pool) Begin(ctx context.Context) (pgx.Tx, error) { // nolint: ireturn
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to begin transaction")
	}
	return tx, nil
}

func (p *Pool) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) { // nolint: ireturn
	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.WithMessage(err, "query failed")
	}
	return rows, nil
}

func (p *Pool) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return errors.WithMessage(err, "failed to execute query")
	}
	if err = pgxscan.ScanOne(dest, rows); err != nil {
		return errors.WithMessage(err, "failed to scan one record")
	}
	return nil
}

func (p *Pool) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return errors.WithMessage(err, "failed to execute query")
	}
	if err = pgxscan.ScanAll(dest, rows); err != nil {
		return errors.WithMessage(err, "failed to scan multiple records")
	}
	return nil
}

func (p *Pool) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	tag, err := p.db.Exec(ctx, query, args...)
	if err != nil {
		return tag, errors.WithMessage(err, "execution failed")
	}
	return tag, nil
}

func (p *Pool) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row { // nolint: ireturn
	return p.db.QueryRow(ctx, query, args...)
}
