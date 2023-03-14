package connection

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Client interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)

	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row

	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

type ClientContext interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row

	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Transaction interface {
	Client
	Commit() error
	Rollback() error
}

type TransactionContext interface {
	ClientContext
	Commit() error
	Rollback() error
}

type TransactionalClient interface {
	Client
	BeginTransaction() (Transaction, error)
}

type TransactionalClientContext interface {
	ClientContext
	BeginTransactionContext(ctx context.Context, opts *sql.TxOptions) (TransactionContext, error)
}

type TransactionalConnection interface {
	TransactionalClientContext
	ReleaseConnection() error
}

type transactionalClient struct {
	*sqlx.DB
}

func (c *transactionalClient) BeginTransaction() (Transaction, error) {
	return c.Beginx()
}

func (c *transactionalClient) BeginTransactionContext(ctx context.Context, opts *sql.TxOptions) (TransactionContext, error) {
	return c.BeginTxx(ctx, opts)
}

type transactionalConnection struct {
	*sqlx.Conn
}

func (c *transactionalConnection) BeginTransactionContext(ctx context.Context, opts *sql.TxOptions) (TransactionContext, error) {
	return c.BeginTxx(ctx, opts)
}

func (c *transactionalConnection) ReleaseConnection() error {
	return c.Conn.Close()
}
