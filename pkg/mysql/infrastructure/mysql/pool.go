package mysql

import (
	"context"
	"errors"
	"sync"

	"github.com/veresnikov/golib/pkg/cache/application/memory"
	inframemory "github.com/veresnikov/golib/pkg/cache/infrastructure/memory"
)

var (
	errConnectionNotFound = errors.New("connection not found")
)

type Pool interface {
	TransactionalClientContext(ctx context.Context) (TransactionalClientContext, error)
	ReleaseConnection(ctx context.Context) error
}

func NewConnectionPool(connector Connector) Pool {
	return &pool{
		connector: connector,
		cache:     inframemory.NewMemoryCache(nil),
	}
}

type cachedConnection struct {
	client     TransactionalConnection
	useCounter int
}

type pool struct {
	connector Connector

	mu    sync.Mutex
	cache memory.Cache
}

func (p *pool) TransactionalClientContext(ctx context.Context) (TransactionalClientContext, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	v, err := p.cache.Get(ctx)
	if err == nil {
		cachedConn := v.(cachedConnection)
		cachedConn.useCounter++
		p.cache.Set(ctx, cachedConn, nil)
		return cachedConn.client, nil
	}
	conn, err := p.connector.Connection(ctx)
	if err != nil {
		return nil, err
	}
	p.cache.Set(ctx, cachedConnection{
		client:     conn,
		useCounter: 1,
	}, nil)
	return conn, nil
}

func (p *pool) ReleaseConnection(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	v, err := p.cache.Get(ctx)
	if err != nil {
		return errConnectionNotFound
	}
	cachedConn := v.(cachedConnection)
	cachedConn.useCounter--
	if cachedConn.useCounter == 0 {
		err = cachedConn.client.ReleaseConnection()
		p.cache.Delete(ctx)
		return err
	}
	p.cache.Set(ctx, cachedConn, nil)
	return nil
}
