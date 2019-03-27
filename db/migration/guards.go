package pg

import (
	"github.com/blend/go-sdk/db/migration"
)

// ColumnNotExists creates a table on the given connection if it does not exist.
func ColumnNotExists(tableName, columnName string) migration.GuardFunc {
	return migration.ColumnNotExists(PredicateColumnExists, tableName, columnName)
}

// ConstraintNotExists creates a table on the given connection if it does not exist.
func ConstraintNotExists(constraintName string) migration.GuardFunc {
	return migration.ConstraintExists(PredicateConstraintExists, constraintName)
}

// TableNotExists creates a table on the given connection if it does not exist.
func TableNotExists(tableName string) migration.GuardFunc {
	return migration.TableNotExists(PredicateTableExists, tableName)
}

// IndexNotExists creates a index on the given connection if it does not exist.
func IndexNotExists(tableName, indexName string) migration.GuardFunc {
	return migration.IndexNotExists(PredicateIndexExists, tableName, indexName)
}

// RoleNotExists creates a new role if it doesn't exist.
func RoleNotExists(roleName string) migration.GuardFunc {
	return migration.RoleNotExists(PredicateRoleExists, roleName)
}

// ColumnExists alters an existing column, erroring if it doesn't exist
func ColumnExists(tableName, columnName string) migration.GuardFunc {
	return migration.ColumnExists(PredicateColumnExists, tableName, columnName)
}

// ConstraintExists alters an existing constraint, erroring if it doesn't exist
func ConstraintExists(constraintName string) migration.GuardFunc {
	return migration.ConstraintExists(PredicateConstraintExists, constraintName)
}

// TableExists alters an existing table, erroring if it doesn't exist
func TableExists(tableName string) migration.GuardFunc {
	return migration.TableExists(PredicateTableExists, tableName)
}

// IndexExists alters an existing index, erroring if it doesn't exist
func IndexExists(tableName, indexName string) migration.GuardFunc {
	return migration.IndexExists(PredicateIndexExists, tableName, indexName)
}

// RoleExists alters an existing role in the db
func RoleExists(roleName string) migration.GuardFunc {
	return migration.RoleExists(PredicateRoleExists, roleName)
}
