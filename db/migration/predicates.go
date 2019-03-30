package migration

import (
	"database/sql"
	"strings"

	"github.com/blend/go-sdk/db"
)

// PredicateTableExists returns if a table exists on the given connection.
func PredicateTableExists(c *db.Connection, tx *sql.Tx, tableName string) (bool, error) {
	return c.Invoke(db.OptTx(tx)).Query(`SELECT 1 FROM pg_catalog.pg_tables WHERE tablename = $1`, strings.ToLower(tableName)).Any()
}

// PredicateColumnExists returns if a column exists on a table on the given connection.
func PredicateColumnExists(c *db.Connection, tx *sql.Tx, tableName, columnName string) (bool, error) {
	return c.Invoke(db.OptTx(tx)).Query(`SELECT 1 FROM information_schema.columns i WHERE i.table_name = $1 and i.column_name = $2`, strings.ToLower(tableName), strings.ToLower(columnName)).Any()
}

// PredicateConstraintExists returns if a constraint exists on a table on the given connection.
func PredicateConstraintExists(c *db.Connection, tx *sql.Tx, constraintName string) (bool, error) {
	return c.Invoke(db.OptTx(tx)).Query(`SELECT 1 FROM pg_constraint WHERE conname = $1`, strings.ToLower(constraintName)).Any()
}

// PredicateIndexExists returns if a index exists on a table on the given connection.
func PredicateIndexExists(c *db.Connection, tx *sql.Tx, tableName, indexName string) (bool, error) {
	return c.Invoke(db.OptTx(tx)).Query(`SELECT 1 FROM pg_catalog.pg_index ix join pg_catalog.pg_class t on t.oid = ix.indrelid join pg_catalog.pg_class i on i.oid = ix.indexrelid WHERE t.relname = $1 and i.relname = $2 and t.relkind = 'r'`, tx, strings.ToLower(tableName), strings.ToLower(indexName)).Any()
}

// PredicateRoleExists returns if a role exists or not.
func PredicateRoleExists(c *db.Connection, tx *sql.Tx, roleName string) (bool, error) {
	return c.Invoke(db.OptTx(tx)).Query(`SELECT 1 FROM pg_roles WHERE rolname ilike $1`, roleName).Any()
}
