package db

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
)

const (
	//DBNilError is a common error
	DBNilError = "connection is nil"
)

const (
	runeComma   = rune(',')
	runeNewline = rune('\n')
	runeTab     = rune('\t')
	runeSpace   = rune(' ')
)

// --------------------------------------------------------------------------------
// Connection
// --------------------------------------------------------------------------------

// New returns a new Connection.
// It will use very bare bones defaults for the config.
func New() *Connection {
	return &Connection{
		config:    &Config{},
		planCache: NewPlanCache(),
	}
}

// NewFromConfig returns a new connection from a config.
func NewFromConfig(cfg *Config) (*Connection, error) {
	dsn := cfg.CreateDSN()
	parsed, err := NewConfigFromDSN(dsn)
	if err != nil {
		return nil, err
	}
	return New().WithConfig(parsed), nil
}

// MustNewFromConfig returns a new connection from a config
// and panics if there is an error.
func MustNewFromConfig(cfg *Config) *Connection {
	conn, err := NewFromConfig(cfg)
	if err != nil {
		panic(err)
	}
	return conn
}

// MustNewFromEnv creates a new db connection from environment variables.
// It will panic if there is an error.
func MustNewFromEnv() *Connection {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		panic(err)
	}
	return MustNewFromConfig(cfg)
}

// NewFromEnv will returns a new connection from a config
// set from environment variables.
func NewFromEnv() (*Connection, error) {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return NewFromConfig(cfg)
}

// Open opens a connection, testing an error and returning it if not nil, and if nil, opening the connection.
// It's designed ot be used in conjunction with a constructor, i.e.
//    conn, err := db.Open(db.NewFromConfig(cfg))
func Open(conn *Connection, err error) (*Connection, error) {
	if err != nil {
		return nil, err
	}
	if err = conn.Open(); err != nil {
		return nil, err
	}
	return conn, nil
}

// Connection is the basic wrapper for connection parameters and saves a reference to the created sql.Connection.
type Connection struct {
	sync.Mutex
	tracer               Tracer
	statementInterceptor StatementInterceptor

	connection *sql.DB
	config     *Config
	bufferPool *BufferPool
	log        logger.FullReceiver
	planCache  *PlanCache
}

// WithConfig sets the config.
func (dbc *Connection) WithConfig(cfg *Config) *Connection {
	dbc.config = cfg
	return dbc
}

// Config returns the config.
func (dbc *Connection) Config() *Config {
	return dbc.config
}

// WithTracer sets the connection tracer and returns a reference.
func (dbc *Connection) WithTracer(tracer Tracer) *Connection {
	dbc.tracer = tracer
	return dbc
}

// Tracer returns the tracer.
func (dbc *Connection) Tracer() Tracer {
	return dbc.tracer
}

// WithStatementInterceptor sets the connection statement interceptor.
func (dbc *Connection) WithStatementInterceptor(interceptor StatementInterceptor) *Connection {
	dbc.statementInterceptor = interceptor
	return dbc
}

// StatementInterceptor returns the statement interceptor.
func (dbc *Connection) StatementInterceptor() StatementInterceptor {
	return dbc.statementInterceptor
}

// Connection returns the underlying driver connection.
func (dbc *Connection) Connection() *sql.DB {
	return dbc.connection
}

// Close implements a closer.
func (dbc *Connection) Close() error {
	if dbc.planCache != nil {
		if err := dbc.planCache.Close(); err != nil {
			return err
		}
	}
	return dbc.connection.Close()
}

// WithLogger sets the connection's diagnostic agent.
func (dbc *Connection) WithLogger(log logger.FullReceiver) *Connection {
	dbc.log = log
	return dbc
}

// Logger returns the diagnostics agent.
func (dbc *Connection) Logger() logger.FullReceiver {
	return dbc.log
}

// PlanCache returns the statement cache.
func (dbc *Connection) PlanCache() *PlanCache {
	return dbc.planCache
}

// Open returns a connection object, either a cached connection object or creating a new one in the process.
func (dbc *Connection) Open() error {
	dbc.Lock()
	defer dbc.Unlock()

	// bail if we've already opened the connection.
	if dbc.connection != nil {
		return Error(ErrConnectionAlreadyOpen)
	}
	if dbc.config == nil {
		return Error(ErrConfigUnset)
	}
	if dbc.bufferPool == nil {
		dbc.bufferPool = NewBufferPool(dbc.config.GetBufferPoolSize())
	}
	if dbc.planCache == nil {
		dbc.planCache = NewPlanCache()
	}

	dsn := dbc.config.CreateDSN()
	namedValues, err := ParseURL(dsn)
	if err != nil {
		return err
	}

	// open the connection
	dbConn, err := sql.Open(dbc.config.GetEngine(), namedValues)
	if err != nil {
		return Error(err)
	}

	dbc.planCache.WithConnection(dbConn)
	dbc.planCache.WithEnabled(!dbc.config.GetPlanCacheDisabled())
	dbc.connection = dbConn
	dbc.connection.SetConnMaxLifetime(dbc.config.GetMaxLifetime())
	dbc.connection.SetMaxIdleConns(dbc.config.GetIdleConnections())
	dbc.connection.SetMaxOpenConns(dbc.config.GetMaxConnections())
	return nil
}

// Begin starts a new transaction.
func (dbc *Connection) Begin(opts ...*sql.TxOptions) (*sql.Tx, error) {
	return dbc.BeginContext(context.Background(), opts...)
}

// BeginContext starts a new transaction in a givent context.
func (dbc *Connection) BeginContext(context context.Context, opts ...*sql.TxOptions) (*sql.Tx, error) {
	if dbc.connection == nil {
		return nil, exception.New(ErrConnectionClosed)
	}
	if len(opts) > 0 {
		tx, err := dbc.connection.BeginTx(context, opts[0])
		return tx, Error(err)
	}
	tx, err := dbc.connection.BeginTx(context, nil)
	return tx, Error(err)
}

// PrepareContext prepares a statement potentially returning a cached version of the statement.
func (dbc *Connection) PrepareContext(context context.Context, cachedPlanKey, statement string, txs ...*sql.Tx) (stmt *sql.Stmt, err error) {
	if dbc.tracer != nil {
		tf := dbc.tracer.Prepare(context, dbc, statement)
		if tf != nil {
			defer func() { tf.Finish(err) }()
		}
	}
	if tx := Tx(txs...); tx != nil {
		stmt, err = tx.PrepareContext(context, statement)
		return
	}
	if dbc.planCache != nil && dbc.planCache.Enabled() && cachedPlanKey != "" {
		stmt, err = dbc.planCache.PrepareContext(context, cachedPlanKey, statement)
		return
	}
	stmt, err = dbc.connection.PrepareContext(context, statement)
	return
}

// --------------------------------------------------------------------------------
// Invocation
// --------------------------------------------------------------------------------

// Invoke returns a new invocation.
func (dbc *Connection) Invoke(context context.Context, txs ...*sql.Tx) *Invocation {
	return &Invocation{
		context:              context,
		tracer:               dbc.tracer,
		statementInterceptor: dbc.statementInterceptor,
		conn:                 dbc,
		startTime:            time.Now().UTC(),
		tx:                   OptionalTx(txs...),
	}
}

// Background returns an empty context.Context.
func (dbc *Connection) Background() context.Context {
	return context.Background()
}

// Ping checks the db connection.
func (dbc *Connection) Ping() error {
	return Error(dbc.connection.Ping())
}

// PingContext checks the db connection.
func (dbc *Connection) PingContext(context context.Context) (err error) {
	if dbc.tracer != nil {
		tf := dbc.tracer.Ping(context, dbc)
		if tf != nil {
			defer func() {
				tf.Finish(err)
			}()
		}
	}

	err = Error(dbc.connection.PingContext(context))
	return
}
