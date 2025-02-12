package postgres

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

type Tx struct {
	db pgx.Tx
}

type TxReq func(tx Tx) error

type TxRunner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

func ExecTx(ctx context.Context, runner TxRunner, req TxReq) error {
	pgxTx, err := runner.Begin(ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to begin transaction")
	}

	tx := Tx{
		db: pgxTx,
	}

	if err = req(tx); err != nil {
		_ = tx.db.Rollback(ctx)
		return errors.WithMessage(err, "transaction execution failed")
	}

	if err = tx.db.Commit(ctx); err != nil {
		return errors.WithMessage(err, "failed to commit transaction")
	}

	return nil
}

func (p Tx) Begin(ctx context.Context) (pgx.Tx, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to begin transaction")
	}
	return tx, nil
}

func (p Tx) Commit(ctx context.Context) error {
	if err := p.db.Commit(ctx); err != nil {
		return errors.WithMessage(err, "failed to commit transaction")
	}
	return nil
}

func (p Tx) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.WithMessagef(err, "query failed: %s", query)
	}
	return rows, nil
}

func (p Tx) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return errors.WithMessagef(err, "query failed: %s", query)
	}
	defer rows.Close()
	if err = pgxscan.ScanOne(dest, rows); err != nil {
		return errors.WithMessage(err, "failed to scan one record")
	}
	return nil
}

func (p Tx) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	tag, err := p.db.Exec(ctx, sql, arguments...)
	if err != nil {
		return tag, errors.WithMessagef(err, "execution failed: %s", sql)
	}
	return tag, nil
}

func (p Tx) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return errors.WithMessagef(err, "query failed: %s", query)
	}
	defer rows.Close()
	if err = pgxscan.ScanAll(dest, rows); err != nil {
		return errors.WithMessage(err, "failed to scan multiple records")
	}
	return nil
}
