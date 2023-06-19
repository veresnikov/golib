package connection

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cenkalti/backoff/v4"
	// nolint
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	driverName = "mysql"
)

type DSN struct {
	User     string
	Password string
	Host     string
	Database string
}

func (dsn *DSN) String() string {
	return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=true", dsn.User, dsn.Password, dsn.Host, dsn.Database)
}

type Connector interface {
	Client() Client
	ClientContext() ClientContext
	TransactionalClient() TransactionalClient
	TransactionalClientContext() TransactionalClientContext
	Connection(ctx context.Context) (TransactionalConnection, error)
	Migrator(fileSystem http.FileSystem) Migrator

	Open(dsn DSN, maxConnections int) error
	Close() error
}

func NewConnector() Connector {
	return &connector{}
}

type connector struct {
	db *sqlx.DB
}

func (c *connector) Migrator(fileSystem http.FileSystem) Migrator {
	return NewMigrator(c.db, fileSystem)
}

func (c *connector) Connection(ctx context.Context) (TransactionalConnection, error) {
	conn, err := c.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	return &transactionalConnection{conn}, nil
}

func (c *connector) Client() Client {
	return c.db
}

func (c *connector) ClientContext() ClientContext {
	return c.db
}

func (c *connector) TransactionalClient() TransactionalClient {
	return &transactionalClient{c.db}
}

func (c *connector) TransactionalClientContext() TransactionalClientContext {
	return &transactionalClient{c.db}
}

func (c *connector) Open(dsn DSN, maxConnections int) error {
	db, err := sqlx.Open(driverName, dsn.String())
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(maxConnections)
	err = backoff.Retry(func() error {
		return db.Ping()
	}, backoff.NewExponentialBackOff())
	if err != nil {
		closeErr := db.Close()
		if closeErr != nil {
			err = errors.Wrap(err, closeErr.Error())
		}
		return err
	}
	c.db = db
	return nil
}

func (c *connector) Close() error {
	return c.db.Close()
}
