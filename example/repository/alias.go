package repository

import (
	"database/sql"
	"github.com/bigmate/bsql"
)

type (
	Transactor = bsql.Transactor
)

var (
	ErrNotFound = sql.ErrNoRows
)
