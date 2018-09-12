package db

import (
	"context"
	"database/sql"
	"sync"

	"github.com/blend/go-sdk/exception"
)

// NewStatementCache returns a new `StatementCache`.
func NewStatementCache() *StatementCache {
	return &StatementCache{
		enabled: true,
		cache:   make(map[string]*sql.Stmt),
	}
}

// StatementCache is a cache of prepared statements.
type StatementCache struct {
	sync.Mutex
	dbc     *sql.DB
	enabled bool
	cache   map[string]*sql.Stmt
}

// WithConnection sets the statement cache connection.
func (sc *StatementCache) WithConnection(conn *sql.DB) *StatementCache {
	sc.dbc = conn
	return sc
}

// Connection returns the underlying connection.
func (sc *StatementCache) Connection() *sql.DB {
	return sc.dbc
}

// WithEnabled sets if the cache is enabled.
func (sc *StatementCache) WithEnabled(enabled bool) *StatementCache {
	sc.enabled = enabled
	return sc
}

// Enabled returns if the statement cache is enabled.
func (sc *StatementCache) Enabled() bool {
	if sc == nil {
		return false
	}
	return sc.enabled
}

// Close implements io.Closer.
func (sc *StatementCache) Close() error {
	sc.Lock()
	defer sc.Unlock()

	var err error
	for _, stmt := range sc.cache {
		err = stmt.Close()
		if err != nil {
			return err
		}
	}
	sc.cache = make(map[string]*sql.Stmt)
	return err
}

// HasStatement returns if the cache contains a statement.
func (sc *StatementCache) HasStatement(statementID string) bool {
	sc.Lock()
	defer sc.Unlock()
	_, hasStmt := sc.cache[statementID]
	return hasStmt
}

// InvalidateStatement removes a statement from the cache.
func (sc *StatementCache) InvalidateStatement(statementID string) error {
	sc.Lock()
	defer sc.Unlock()

	if statement, hasStatement := sc.cache[statementID]; hasStatement {
		delete(sc.cache, statementID)
		if statement != nil {
			return exception.New(statement.Close())
		}
	}
	return nil
}

// PrepareContext returns a cached expression for a statement, or creates and caches a new one.
func (sc *StatementCache) PrepareContext(context context.Context, statementID, statement string, tx *sql.Tx) (*sql.Stmt, error) {
	if tx != nil {
		return tx.PrepareContext(context, statement)
	}

	if !sc.enabled {
		return sc.dbc.PrepareContext(context, statement)
	}

	sc.Lock()
	defer sc.Unlock()

	if stmt, hasStmt := sc.cache[statementID]; hasStmt {
		return stmt, nil
	}

	stmt, err := sc.dbc.PrepareContext(context, statement)
	if err != nil {
		return nil, err
	}

	sc.cache[statementID] = stmt
	return stmt, nil
}
