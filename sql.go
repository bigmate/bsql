package bsql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/jmoiron/sqlx"
	"os"
)

type txKey struct{}

type Database interface {
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Transactor interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Factory interface {
	FromContext(ctx context.Context) Database
}

type transactor struct {
	db  *sqlx.DB
	log Logger
}

func New(db *sqlx.DB, log Logger) (Transactor, Factory) {
	if log == nil {
		log = logger{os.Stderr}
	}

	trx := &transactor{
		db:  db,
		log: log,
	}

	return trx, trx
}

func (t *transactor) Begin(ctx context.Context) (context.Context, error) {
	switch tx := ctx.Value(txKey{}).(type) {
	case dummyTx:
		return ctx, nil
	case *sqlx.Tx:
		return context.WithValue(ctx, txKey{}, dummyTx{tx}), nil
	}

	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		t.log.Errorf("transactor: failed to begin transaction: %v", err)
		return ctx, err
	}

	return context.WithValue(ctx, txKey{}, tx), nil
}

func (t *transactor) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(txKey{}).(driver.Tx)
	if !ok {
		return nil
	}

	if err := tx.Commit(); err != nil {
		t.log.Errorf("transactor: failed to commit transaction: %v", err)
		return err
	}

	return nil
}

func (t *transactor) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(txKey{}).(driver.Tx)
	if !ok {
		return nil
	}

	if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		t.log.Errorf("transactor: failed to rollback transaction: %v", err)
		return err
	}

	return nil
}

func (t *transactor) FromContext(ctx context.Context) Database {
	if tx, ok := ctx.Value(txKey{}).(Database); ok {
		return tx
	}

	return t.db
}

type dummyTx struct {
	*sqlx.Tx
}

func (dummyTx) Commit() error {
	return nil
}

func (dummyTx) Rollback() error {
	return nil
}
