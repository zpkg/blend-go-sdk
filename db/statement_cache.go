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
		cache:   sync.Map{},
	}
}

// StatementCache is a cache of prepared statements.
type StatementCache struct {
	conn    *sql.DB
	enabled bool
	cache   sync.Map
}

// WithConnection sets the statement cache connection.
func (sc *StatementCache) WithConnection(conn *sql.DB) *StatementCache {
	sc.conn = conn
	return sc
}

// Connection returns the underlying connection.
func (sc *StatementCache) Connection() *sql.DB {
	return sc.conn
}

// WithEnabled sets if the cache is enabled.
func (sc *StatementCache) WithEnabled(enabled bool) *StatementCache {
	sc.enabled = enabled
	return sc
}

// Enabled returns if the statement cache is enabled.
func (sc *StatementCache) Enabled() bool {
	return sc.enabled
}

// Close implements io.Closer.
func (sc *StatementCache) Close() (err error) {
	sc.cache.Range(func(k, v interface{}) bool {
		err = v.(*sql.Stmt).Close()
		return err == nil
	})
	return
}

// HasStatement returns if the cache contains a statement.
func (sc *StatementCache) HasStatement(statementID string) bool {
	_, hasStmt := sc.cache.Load(statementID)
	return hasStmt
}

// InvalidateStatement removes a statement from the cache.
func (sc *StatementCache) InvalidateStatement(statementID string) (err error) {
	stmt, ok := sc.cache.Load(statementID)
	if !ok {
		return
	}
	sc.cache.Delete(statementID)
	return stmt.(*sql.Stmt).Close()
}

// PrepareContext returns a cached expression for a statement, or creates and caches a new one.
func (sc *StatementCache) PrepareContext(context context.Context, statementID, statement string) (*sql.Stmt, error) {
	if len(statementID) == 0 {
		return nil, exception.New(ErrStatementLabelUnset)
	}

	if stmt, hasStmt := sc.cache.Load(statementID); hasStmt {
		return stmt.(*sql.Stmt), nil
	}

	stmt, err := sc.conn.PrepareContext(context, statement)
	if err != nil {
		return nil, err
	}

	sc.cache.Store(statementID, stmt)
	return stmt, nil
}
