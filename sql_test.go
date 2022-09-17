package bsql

import (
	"bytes"
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"reflect"
	"testing"
)

func TestTransactor_Commit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectCommit()

	buf := &bytes.Buffer{}
	DBx := sqlx.NewDb(db, "sqlmock")
	trx, f := New(DBx, logger{buf})
	ctx := context.Background()

	ctx, err = trx.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin: %v", err)
	}

	repo := f.FromContext(ctx)

	if _, ok := repo.(*sqlx.Tx); !ok {
		t.Fatalf("repo is not type of *sqlx.Tx")
	}

	if err = trx.Commit(ctx); err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	if len(buf.Bytes()) > 0 {
		t.Fatalf("buffer is expected to be empty: %s", buf.Bytes())
	}
}

func TestTransactor_Rollback(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectRollback()

	buf := &bytes.Buffer{}
	DBx := sqlx.NewDb(db, "sqlmock")
	trx, f := New(DBx, logger{buf})

	ctx, err := trx.Begin(context.Background())
	if err != nil {
		t.Fatalf("failed to begin: %v", err)
	}

	repo := f.FromContext(ctx)

	if _, ok := repo.(*sqlx.Tx); !ok {
		t.Fatalf("repo is not type of *sqlx.Tx")
	}

	if err = trx.Rollback(ctx); err != nil {
		t.Fatalf("failed to rollback: %v", err)
	}

	if len(buf.Bytes()) > 0 {
		t.Fatalf("buffer is expected to be empty: %s", buf.Bytes())
	}
}

func TestTransactor_Begin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}

	err = errors.New("connection is closed")

	mock.ExpectBegin().WillReturnError(err)

	buf := &bytes.Buffer{}
	DBx := sqlx.NewDb(db, "sqlmock")
	trx, _ := New(DBx, logger{buf})

	_, begErr := trx.Begin(context.Background())

	if !reflect.DeepEqual(err, begErr) {
		t.Fatalf("expected error be equal")
	}

	if len(buf.Bytes()) == 0 {
		t.Fatalf("buf is empty")
	}
}

func TestFactory_FromContext(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}

	DBx := sqlx.NewDb(db, "sqlmock")
	_, f := New(DBx, nil)

	repo := f.FromContext(context.Background())

	if !reflect.DeepEqual(repo, DBx) {
		t.Fatalf("expected DBx and repo be same, DBx: %T, repo: %T", DBx, repo)
	}
}

func TestNopCommitter_Commit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}

	mock.ExpectBegin()

	DBx := sqlx.NewDb(db, "sqlmock")
	trx, f := New(DBx, nil)

	ctx, err := trx.Begin(context.Background())
	if err != nil {
		t.Fatalf("failed to begin: %v", err)
	}

	ctx, err = trx.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin: %v", err)
	}

	if _, ok := f.FromContext(ctx).(dummyTx); !ok {
		t.Fatal("repo is not type of dummyTx")
	}

	ctx, err = trx.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin: %v", err)
	}

	if _, ok := f.FromContext(ctx).(dummyTx); !ok {
		t.Fatal("repo is not type of dummyTx")
	}

	if err = trx.Commit(ctx); err != nil {
		t.Fatalf("failed to commit: %v", err)
	}
}
