package migration

// ColumnNotExists creates a table on the given connection if it does not exist.
func ColumnNotExists(tableName, columnName string) GuardFunc {
	return ColumnNotExistsFromPredicate(PredicateColumnExists, tableName, columnName)
}

// ConstraintNotExists creates a table on the given connection if it does not exist.
func ConstraintNotExists(constraintName string) GuardFunc {
	return ConstraintExistsFromPredicate(PredicateConstraintExists, constraintName)
}

// TableNotExists creates a table on the given connection if it does not exist.
func TableNotExists(tableName string) GuardFunc {
	return TableNotExistsFromPredicate(PredicateTableExists, tableName)
}

// IndexNotExists creates a index on the given connection if it does not exist.
func IndexNotExists(tableName, indexName string) GuardFunc {
	return IndexNotExistsFromPredicate(PredicateIndexExists, tableName, indexName)
}

// RoleNotExists creates a new role if it doesn't exist.
func RoleNotExists(roleName string) GuardFunc {
	return RoleNotExistsFromPredicate(PredicateRoleExists, roleName)
}

// ColumnExists alters an existing column, erroring if it doesn't exist
func ColumnExists(tableName, columnName string) GuardFunc {
	return ColumnExistsFromPredicate(PredicateColumnExists, tableName, columnName)
}

// ConstraintExists alters an existing constraint, erroring if it doesn't exist
func ConstraintExists(constraintName string) GuardFunc {
	return ConstraintExistsFromPredicate(PredicateConstraintExists, constraintName)
}

// TableExists alters an existing table, erroring if it doesn't exist
func TableExists(tableName string) GuardFunc {
	return TableExistsFromPredicate(PredicateTableExists, tableName)
}

// IndexExists alters an existing index, erroring if it doesn't exist
func IndexExists(tableName, indexName string) GuardFunc {
	return IndexExistsFromPredicate(PredicateIndexExists, tableName, indexName)
}

// RoleExists alters an existing role in the db
func RoleExists(roleName string) GuardFunc {
	return RoleExistsFromPredicate(PredicateRoleExists, roleName)
}
