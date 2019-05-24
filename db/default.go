package db

var (
	defaultConnection *Connection
)

// SetDefault sets an alias created with `CreateDbAlias` as default. This lets you refer to it later via. `Default()`
//
//	db.CreateDbAlias("main", spiffy.NewDbConnection("localhost", "test_db", "", ""))
//	db.SetDefault("main")
//	execErr := db.Default().Execute("select 'ok!')
//
// This will then let you refer to the alias via. `Default()`
func SetDefault(conn *Connection) {
	defaultConnection = conn
}

// Default returns a reference to the DbConnection set as default.
//
//	db.Default().Exec("select 'ok!")
//
func Default() *Connection {
	return defaultConnection
}

// OpenDefault sets the default connection and opens it.
func OpenDefault(conn *Connection) error {
	err := conn.Open()
	if err != nil {
		return err
	}
	SetDefault(conn)
	return nil
}
