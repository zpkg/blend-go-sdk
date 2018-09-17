package db

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
)

// Invocation is a specific operation against a context.
type Invocation struct {
	statementLabel string

	conn          *Connection
	context       context.Context
	cancel        func()
	tracer        Tracer
	traceFinisher TraceFinisher
	startTime     time.Time
	tx            *sql.Tx
}

// StartTime returns the invocation start time.
func (i *Invocation) StartTime() time.Time {
	return i.startTime
}

// WithContext sets the context and returns a reference to the invocation.
func (i *Invocation) WithContext(context context.Context) *Invocation {
	i.context = context
	return i
}

// Context returns the underlying context.
func (i *Invocation) Context() context.Context {
	if i.context == nil {
		return context.Background()
	}
	return i.context
}

// WithCancel sets an optional cancel callback.
func (i *Invocation) WithCancel(cancel func()) *Invocation {
	i.cancel = cancel
	return i
}

// Cancel returns the optional cancel callback.
func (i *Invocation) Cancel() func() {
	return i.cancel
}

// WithLabel instructs the query generator to get or create a cached prepared statement.
func (i *Invocation) WithLabel(label string) *Invocation {
	i.statementLabel = label
	return i
}

// Label returns the statement / plan cache label for the context.
func (i *Invocation) Label() string {
	return i.statementLabel
}

// WithTx sets the tx
func (i *Invocation) WithTx(tx *sql.Tx) *Invocation {
	i.tx = tx
	return i
}

// Tx returns the underlying transaction.
func (i *Invocation) Tx() *sql.Tx {
	return i.tx
}

// Prepare returns a cached or newly prepared statment plan for a given sql statement.
func (i *Invocation) Prepare(statement string) (*sql.Stmt, error) {
	return i.conn.PrepareContext(i.Context(), i.statementLabel, statement, i.tx)
}

// Exec executes a sql statement with a given set of arguments.
func (i *Invocation) Exec(statement string, args ...interface{}) (err error) {
	var stmt *sql.Stmt
	i.start(statement)
	defer func() { err = i.finish(statement, recover(), err) }()

	stmt, err = i.Prepare(statement)
	if err != nil {
		err = exception.New(err)
		return
	}
	defer func() { err = i.closeStatement(stmt, err) }()

	if _, err = stmt.ExecContext(i.Context(), args...); err != nil {
		err = exception.New(err)
		return
	}
	return
}

// Query returns a new query object for a given sql query and arguments.
func (i *Invocation) Query(statement string, args ...interface{}) *Query {
	i.start(statement)
	return &Query{
		context:        i.Context(),
		statement:      statement,
		statementLabel: i.statementLabel,
		args:           args,
		conn:           i.conn,
		inv:            i,
		tx:             i.tx,
	}
}

// Get returns a given object based on a group of primary key ids within a transaction.
func (i *Invocation) Get(object DatabaseMapped, ids ...interface{}) (err error) {
	if len(ids) == 0 {
		err = exception.New(ErrInvalidIDs)
		return
	}

	var queryBody string
	var stmt *sql.Stmt
	var cols *ColumnCollection
	defer func() { err = i.finish(queryBody, recover(), err) }()

	if i.statementLabel, queryBody, cols, err = i.generateGet(object); err != nil {
		err = exception.New(err)
		return
	}

	if stmt, err = i.Prepare(queryBody); err != nil {
		err = exception.New(err)
		return
	}
	defer func() { err = i.closeStatement(stmt, err) }()
	i.start(queryBody)

	row := stmt.QueryRowContext(i.Context(), ids...)
	var populateErr error
	if typed, ok := object.(Populatable); ok {
		populateErr = typed.Populate(row)
	} else {
		populateErr = PopulateInOrder(object, row, cols)
	}
	if populateErr != nil && !exception.Is(populateErr, sql.ErrNoRows) {
		err = exception.New(populateErr)
		return
	}

	return
}

// GetAll returns all rows of an object mapped table wrapped in a transaction.
func (i *Invocation) GetAll(collection interface{}) (err error) {
	var queryBody string
	var stmt *sql.Stmt
	var rows *sql.Rows
	var cols *ColumnCollection
	var collectionType reflect.Type
	defer func() { err = i.finish(queryBody, recover(), err) }()

	i.statementLabel, queryBody, cols, collectionType = i.generateGetAll(collection)

	if stmt, err = i.Prepare(queryBody); err != nil {
		err = exception.New(err)
		return
	}
	defer func() { err = i.closeStatement(stmt, err) }()
	i.start(queryBody)

	if rows, err = stmt.QueryContext(i.Context()); err != nil {
		err = exception.New(err)
		return
	}
	defer func() { err = exception.Nest(err, rows.Close()) }()

	collectionValue := reflectValue(collection)
	for rows.Next() {
		var obj interface{}
		if obj, err = makeNewDatabaseMapped(collectionType); err != nil {
			err = exception.New(err)
			return
		}

		if typed, ok := obj.(Populatable); ok {
			err = typed.Populate(rows)
		} else {
			err = PopulateInOrder(obj, rows, cols)
		}
		if err != nil {
			err = exception.New(err)
			return
		}

		objValue := reflectValue(obj)
		collectionValue.Set(reflect.Append(collectionValue, objValue))
	}
	return
}

// Create writes an object to the database within a transaction.
func (i *Invocation) Create(object DatabaseMapped) (err error) {
	var queryBody string
	var stmt *sql.Stmt
	var writeCols, autos *ColumnCollection
	defer func() { err = i.finish(queryBody, recover(), err) }()

	i.statementLabel, queryBody, writeCols, autos = i.generateCreate(object)

	if stmt, err = i.Prepare(queryBody); err != nil {
		err = exception.New(err)
		return
	}
	defer func() { err = i.closeStatement(stmt, err) }()

	i.start(queryBody)

	if autos.Len() == 0 {
		if _, err = stmt.ExecContext(i.Context(), writeCols.ColumnValues(object)...); err != nil {
			err = exception.New(err)
			return
		}
		return
	}

	autoValues := i.autoValues(autos)
	if err = stmt.QueryRowContext(i.Context(), writeCols.ColumnValues(object)...).Scan(autoValues...); err != nil {
		err = exception.New(err)
		return
	}
	if err = i.setAutos(object, autos, autoValues); err != nil {
		err = exception.New(err)
		return
	}

	return
}

// CreateIfNotExists writes an object to the database if it does not already exist within a transaction.
func (i *Invocation) CreateIfNotExists(object DatabaseMapped) (err error) {
	var queryBody string
	var stmt *sql.Stmt
	var autos, writeCols *ColumnCollection
	defer func() { err = i.finish(queryBody, recover(), err) }()

	i.statementLabel, queryBody, autos, writeCols = i.generateCreateIfNotExists(object)

	if stmt, err = i.Prepare(queryBody); err != nil {
		err = exception.New(err)
		return
	}
	defer func() { err = i.closeStatement(stmt, err) }()
	i.start(queryBody)

	if autos.Len() == 0 {
		if _, err = stmt.ExecContext(i.context, writeCols.ColumnValues(object)...); err != nil {
			err = exception.New(err)
		}
		return
	}

	autoValues := i.autoValues(autos)
	if err = stmt.QueryRowContext(i.Context(), writeCols.ColumnValues(object)...).Scan(autoValues...); err != nil {
		err = exception.New(err)
		return
	}
	if err = i.setAutos(object, autos, autoValues); err != nil {
		err = exception.New(err)
		return
	}

	return
}

// CreateMany writes many objects to the database in a single insert.
// Important; this will not use cached statements ever because the generated query
// is different for each cardinality of objects.
func (i *Invocation) CreateMany(objects interface{}) (err error) {
	var queryBody string
	var writeCols *ColumnCollection
	var sliceValue reflect.Value
	defer func() { err = i.finish(queryBody, recover(), err) }()

	queryBody, writeCols, sliceValue = i.generateCreateMany(objects)

	var colValues []interface{}
	for row := 0; row < sliceValue.Len(); row++ {
		colValues = append(colValues, writeCols.ColumnValues(sliceValue.Index(row).Interface())...)
	}

	if i.tx != nil {
		_, err = i.tx.ExecContext(i.Context(), queryBody, colValues...)
	} else {
		_, err = i.conn.connection.ExecContext(i.Context(), queryBody, colValues...)
	}
	if err != nil {
		err = exception.New(err)
		return
	}
	return
}

// Update updates an object wrapped in a transaction.
func (i *Invocation) Update(object DatabaseMapped) (err error) {
	var queryBody string
	var stmt *sql.Stmt
	var pks, writeCols *ColumnCollection
	defer func() { err = i.finish(queryBody, recover(), err) }()

	i.statementLabel, queryBody, pks, writeCols = i.generateUpdate(object)

	if stmt, err = i.Prepare(queryBody); err != nil {
		err = exception.New(err)
		return
	}
	defer func() { err = i.closeStatement(stmt, err) }()
	i.start(queryBody)

	if _, err = stmt.ExecContext(i.Context(), append(writeCols.ColumnValues(object), pks.ColumnValues(object)...)...); err != nil {
		err = exception.New(err)
		return
	}
	return
}

// Upsert inserts the object if it doesn't exist already (as defined by its primary keys) or updates it wrapped in a transaction.
func (i *Invocation) Upsert(object DatabaseMapped) (err error) {
	var queryBody string
	var autos, writeCols *ColumnCollection
	var stmt *sql.Stmt
	//defer func() { err = i.finish(queryBody, recover(), err) }()

	i.statementLabel, queryBody, autos, writeCols = i.generateUpsert(object)

	if stmt, err = i.Prepare(queryBody); err != nil {
		err = exception.New(err)
		return
	}
	defer func() { err = i.closeStatement(stmt, err) }()
	i.start(queryBody)

	if autos.Len() == 0 {
		if _, err = stmt.ExecContext(i.Context(), writeCols.ColumnValues(object)...); err != nil {
			err = exception.New(err)
			return
		}
		return
	}

	autoValues := i.autoValues(autos)
	if err = stmt.QueryRowContext(i.Context(), writeCols.ColumnValues(object)...).Scan(autoValues...); err != nil {
		err = exception.New(err)
		return
	}
	if err = i.setAutos(object, autos, autoValues); err != nil {
		err = exception.New(err)
		return
	}

	return
}

// Exists returns a bool if a given object exists (utilizing the primary key columns if they exist) wrapped in a transaction.
func (i *Invocation) Exists(object DatabaseMapped) (exists bool, err error) {
	var queryBody string
	var pks *ColumnCollection
	var stmt *sql.Stmt
	defer func() { err = i.finish(queryBody, recover(), err) }()

	if i.statementLabel, queryBody, pks, err = i.generateExists(object); err != nil {
		err = exception.New(err)
		return
	}
	if stmt, err = i.Prepare(queryBody); err != nil {
		err = exception.New(err)
		return
	}
	defer func() { err = i.closeStatement(stmt, err) }()
	i.start(queryBody)

	var value int
	if err = stmt.QueryRowContext(i.Context(), pks.ColumnValues(object)...).Scan(&value); err != nil {
		err = exception.New(err)
		return
	}

	exists = value != 0
	return
}

// Delete deletes an object from the database wrapped in a transaction.
func (i *Invocation) Delete(object DatabaseMapped) (err error) {
	var queryBody string
	var stmt *sql.Stmt
	var pks *ColumnCollection
	defer func() { err = i.finish(queryBody, recover(), err) }()

	if i.statementLabel, queryBody, pks, err = i.generateDelete(object); err != nil {
		return
	}

	if stmt, err = i.Prepare(queryBody); err != nil {
		err = exception.New(err)
		return
	}
	defer func() { err = i.closeStatement(stmt, err) }()
	i.start(queryBody)

	if _, err = stmt.ExecContext(i.Context(), pks.ColumnValues(object)...); err != nil {
		err = exception.New(err)
		return
	}
	return
}

// Truncate completely empties a table in a single command.
func (i *Invocation) Truncate(object DatabaseMapped) (err error) {
	var queryBody string
	var stmt *sql.Stmt
	defer func() { err = i.finish(queryBody, recover(), err) }()

	i.statementLabel, queryBody = i.generateTruncate(object)

	if stmt, err = i.Prepare(queryBody); err != nil {
		err = exception.New(err)
		return
	}
	defer func() { err = i.closeStatement(stmt, err) }()
	i.start(queryBody)

	if _, err = stmt.ExecContext(i.Context()); err != nil {
		err = exception.New(err)
		return
	}
	return
}

// --------------------------------------------------------------------------------
// query body generators
// --------------------------------------------------------------------------------

func (i *Invocation) generateGet(object DatabaseMapped) (statementLabel, queryBody string, cols *ColumnCollection, err error) {
	tableName := TableName(object)

	cols = getCachedColumnCollectionFromInstance(object)
	pks := cols.PrimaryKeys()
	if pks.Len() == 0 {
		err = exception.New(ErrNoPrimaryKey)
		return
	}

	queryBodyBuffer := i.conn.bufferPool.Get()

	queryBodyBuffer.WriteString("SELECT ")
	for i, name := range cols.ColumnNames() {
		queryBodyBuffer.WriteString(name)
		if i < (cols.Len() - 1) {
			queryBodyBuffer.WriteRune(runeComma)
		}
	}

	queryBodyBuffer.WriteString(" FROM ")
	queryBodyBuffer.WriteString(tableName)
	queryBodyBuffer.WriteString(" WHERE ")

	for i, pk := range pks.Columns() {
		queryBodyBuffer.WriteString(pk.ColumnName)
		queryBodyBuffer.WriteString(" = ")
		queryBodyBuffer.WriteString("$" + strconv.Itoa(i+1))

		if i < (pks.Len() - 1) {
			queryBodyBuffer.WriteString(" AND ")
		}
	}

	statementLabel = tableName + "_get"
	queryBody = queryBodyBuffer.String()
	i.conn.bufferPool.Put(queryBodyBuffer)
	return
}

func (i *Invocation) generateGetAll(collection interface{}) (statementLabel, queryBody string, cols *ColumnCollection, collectionType reflect.Type) {
	collectionType = reflectSliceType(collection)
	tableName := TableNameByType(collectionType)

	cols = getCachedColumnCollectionFromType(tableName, reflectSliceType(collection)).NotReadOnly()

	queryBodyBuffer := i.conn.bufferPool.Get()
	queryBodyBuffer.WriteString("SELECT ")
	for i, name := range cols.ColumnNames() {
		queryBodyBuffer.WriteString(name)
		if i < (cols.Len() - 1) {
			queryBodyBuffer.WriteRune(runeComma)
		}
	}
	queryBodyBuffer.WriteString(" FROM ")
	queryBodyBuffer.WriteString(tableName)

	queryBody = queryBodyBuffer.String()
	statementLabel = tableName + "_get_all"
	i.conn.bufferPool.Put(queryBodyBuffer)
	return
}

func (i *Invocation) generateCreate(object DatabaseMapped) (statementLabel, queryBody string, writeCols, autos *ColumnCollection) {
	tableName := TableName(object)

	cols := getCachedColumnCollectionFromInstance(object)
	writeCols = cols.WriteColumns()
	autos = cols.Autos()

	queryBodyBuffer := i.conn.bufferPool.Get()

	queryBodyBuffer.WriteString("INSERT INTO ")
	queryBodyBuffer.WriteString(tableName)
	queryBodyBuffer.WriteString(" (")
	for i, name := range writeCols.ColumnNames() {
		queryBodyBuffer.WriteString(name)
		if i < (writeCols.Len() - 1) {
			queryBodyBuffer.WriteRune(runeComma)
		}
	}
	queryBodyBuffer.WriteString(") VALUES (")
	for x := 0; x < writeCols.Len(); x++ {
		queryBodyBuffer.WriteString("$" + strconv.Itoa(x+1))
		if x < (writeCols.Len() - 1) {
			queryBodyBuffer.WriteRune(runeComma)
		}
	}
	queryBodyBuffer.WriteString(")")

	if autos.Len() > 0 {
		queryBodyBuffer.WriteString(" RETURNING ")
		queryBodyBuffer.WriteString(autos.ColumnNamesCSV())
	}

	queryBody = queryBodyBuffer.String()
	statementLabel = tableName + "_create"
	i.conn.bufferPool.Put(queryBodyBuffer)
	return
}

func (i *Invocation) generateCreateIfNotExists(object DatabaseMapped) (statementLabel, queryBody string, autos, writeCols *ColumnCollection) {
	cols := getCachedColumnCollectionFromInstance(object)

	writeCols = cols.WriteColumns()
	autos = cols.Autos()

	pks := cols.PrimaryKeys()
	tableName := TableName(object)

	queryBodyBuffer := i.conn.bufferPool.Get()

	queryBodyBuffer.WriteString("INSERT INTO ")
	queryBodyBuffer.WriteString(tableName)
	queryBodyBuffer.WriteString(" (")
	for i, name := range writeCols.ColumnNames() {
		queryBodyBuffer.WriteString(name)
		if i < (writeCols.Len() - 1) {
			queryBodyBuffer.WriteRune(runeComma)
		}
	}
	queryBodyBuffer.WriteString(") VALUES (")
	for x := 0; x < writeCols.Len(); x++ {
		queryBodyBuffer.WriteString("$" + strconv.Itoa(x+1))
		if x < (writeCols.Len() - 1) {
			queryBodyBuffer.WriteRune(runeComma)
		}
	}
	queryBodyBuffer.WriteString(")")

	if pks.Len() > 0 {
		queryBodyBuffer.WriteString(" ON CONFLICT (")
		pkColumnNames := pks.ColumnNames()
		for i, name := range pkColumnNames {
			queryBodyBuffer.WriteString(name)
			if i < len(pkColumnNames)-1 {
				queryBodyBuffer.WriteRune(runeComma)
			}
		}
		queryBodyBuffer.WriteString(") DO NOTHING")
	}

	if autos.Len() > 0 {
		queryBodyBuffer.WriteString(" RETURNING ")
		queryBodyBuffer.WriteString(autos.ColumnNamesCSV())
	}

	queryBody = queryBodyBuffer.String()
	statementLabel = tableName + "_create_if_not_exists"
	i.conn.bufferPool.Put(queryBodyBuffer)
	return
}

func (i *Invocation) generateCreateMany(objects interface{}) (queryBody string, writeCols *ColumnCollection, sliceValue reflect.Value) {
	sliceValue = reflectValue(objects)
	sliceType := reflectSliceType(objects)
	tableName := TableNameByType(sliceType)

	cols := getCachedColumnCollectionFromType(tableName, sliceType)
	writeCols = cols.WriteColumns()

	queryBodyBuffer := i.conn.bufferPool.Get()

	queryBodyBuffer.WriteString("INSERT INTO ")
	queryBodyBuffer.WriteString(tableName)
	queryBodyBuffer.WriteString(" (")
	for i, name := range writeCols.ColumnNames() {
		queryBodyBuffer.WriteString(name)
		if i < (writeCols.Len() - 1) {
			queryBodyBuffer.WriteRune(runeComma)
		}
	}

	queryBodyBuffer.WriteString(") VALUES ")

	metaIndex := 1
	for x := 0; x < sliceValue.Len(); x++ {
		queryBodyBuffer.WriteString("(")
		for y := 0; y < writeCols.Len(); y++ {
			queryBodyBuffer.WriteString(fmt.Sprintf("$%d", metaIndex))
			metaIndex = metaIndex + 1
			if y < writeCols.Len()-1 {
				queryBodyBuffer.WriteRune(runeComma)
			}
		}
		queryBodyBuffer.WriteString(")")
		if x < sliceValue.Len()-1 {
			queryBodyBuffer.WriteRune(runeComma)
		}
	}

	queryBody = queryBodyBuffer.String()
	i.conn.bufferPool.Put(queryBodyBuffer)
	return
}

func (i *Invocation) generateUpdate(object DatabaseMapped) (statementLabel, queryBody string, pks, writeCols *ColumnCollection) {
	tableName := TableName(object)

	cols := getCachedColumnCollectionFromInstance(object)

	pks = cols.PrimaryKeys()
	writeCols = cols.WriteColumns()

	queryBodyBuffer := i.conn.bufferPool.Get()

	queryBodyBuffer.WriteString("UPDATE ")
	queryBodyBuffer.WriteString(tableName)
	queryBodyBuffer.WriteString(" SET ")

	var writeColIndex int
	var col Column
	for ; writeColIndex < writeCols.Len(); writeColIndex++ {
		col = writeCols.columns[writeColIndex]
		queryBodyBuffer.WriteString(col.ColumnName)
		queryBodyBuffer.WriteString(" = $" + strconv.Itoa(writeColIndex+1))
		if writeColIndex != (writeCols.Len() - 1) {
			queryBodyBuffer.WriteRune(runeComma)
		}
	}

	queryBodyBuffer.WriteString(" WHERE ")
	for i, pk := range pks.Columns() {
		queryBodyBuffer.WriteString(pk.ColumnName)
		queryBodyBuffer.WriteString(" = ")
		queryBodyBuffer.WriteString("$" + strconv.Itoa(i+(writeColIndex+1)))

		if i < (pks.Len() - 1) {
			queryBodyBuffer.WriteString(" AND ")
		}
	}

	queryBody = queryBodyBuffer.String()
	statementLabel = tableName + "_update"
	i.conn.bufferPool.Put(queryBodyBuffer)
	return
}

func (i *Invocation) generateUpsert(object DatabaseMapped) (statementLabel, queryBody string, autos, writeCols *ColumnCollection) {
	tableName := TableName(object)
	cols := getCachedColumnCollectionFromInstance(object)
	conflictUpdateCols := cols.NotReadOnly().NotAutos().NotPrimaryKeys()

	writeCols = cols.NotReadOnly().NotAutos()
	autos = cols.Autos()
	pks := cols.PrimaryKeys()

	colNames := writeCols.ColumnNames()

	queryBodyBuffer := i.conn.bufferPool.Get()

	queryBodyBuffer.WriteString("INSERT INTO ")
	queryBodyBuffer.WriteString(tableName)
	queryBodyBuffer.WriteString(" (")
	for i, name := range colNames {
		queryBodyBuffer.WriteString(name)
		if i < len(colNames)-1 {
			queryBodyBuffer.WriteRune(runeComma)
		}
	}
	queryBodyBuffer.WriteString(") VALUES (")

	for x := 0; x < writeCols.Len(); x++ {
		queryBodyBuffer.WriteString("$" + strconv.Itoa(x+1))
		if x < (writeCols.Len() - 1) {
			queryBodyBuffer.WriteRune(runeComma)
		}
	}

	queryBodyBuffer.WriteString(")")

	if pks.Len() > 0 {
		tokenMap := map[string]string{}
		for i, col := range writeCols.Columns() {
			tokenMap[col.ColumnName] = "$" + strconv.Itoa(i+1)
		}

		queryBodyBuffer.WriteString(" ON CONFLICT (")
		pkColumnNames := pks.ColumnNames()
		for i, name := range pkColumnNames {
			queryBodyBuffer.WriteString(name)
			if i < len(pkColumnNames)-1 {
				queryBodyBuffer.WriteRune(runeComma)
			}
		}
		queryBodyBuffer.WriteString(") DO UPDATE SET ")

		conflictCols := conflictUpdateCols.Columns()
		for i, col := range conflictCols {
			queryBodyBuffer.WriteString(col.ColumnName + " = " + tokenMap[col.ColumnName])
			if i < (len(conflictCols) - 1) {
				queryBodyBuffer.WriteRune(runeComma)
			}
		}
	}
	if autos.Len() > 0 {
		queryBodyBuffer.WriteString(" RETURNING ")
		queryBodyBuffer.WriteString(autos.ColumnNamesCSV())
	}

	queryBody = queryBodyBuffer.String()
	statementLabel = tableName + "_upsert"
	i.conn.bufferPool.Put(queryBodyBuffer)
	return
}

func (i *Invocation) generateExists(object DatabaseMapped) (statementLabel, queryBody string, pks *ColumnCollection, err error) {
	tableName := TableName(object)
	pks = getCachedColumnCollectionFromInstance(object).PrimaryKeys()
	if pks.Len() == 0 {
		err = exception.New(ErrNoPrimaryKey)
		return
	}
	queryBodyBuffer := i.conn.bufferPool.Get()
	queryBodyBuffer.WriteString("SELECT 1 FROM ")
	queryBodyBuffer.WriteString(tableName)
	queryBodyBuffer.WriteString(" WHERE ")
	for i, pk := range pks.Columns() {
		queryBodyBuffer.WriteString(pk.ColumnName)
		queryBodyBuffer.WriteString(" = ")
		queryBodyBuffer.WriteString("$" + strconv.Itoa(i+1))

		if i < (pks.Len() - 1) {
			queryBodyBuffer.WriteString(" AND ")
		}
	}
	statementLabel = tableName + "_exists"
	queryBody = queryBodyBuffer.String()
	i.conn.bufferPool.Put(queryBodyBuffer)
	return
}

func (i *Invocation) generateDelete(object DatabaseMapped) (statementLabel, queryBody string, pks *ColumnCollection, err error) {
	tableName := TableName(object)
	pks = getCachedColumnCollectionFromInstance(object).PrimaryKeys()
	if len(pks.Columns()) == 0 {
		err = exception.New(ErrNoPrimaryKey)
		return
	}
	queryBodyBuffer := i.conn.bufferPool.Get()
	queryBodyBuffer.WriteString("DELETE FROM ")
	queryBodyBuffer.WriteString(tableName)
	queryBodyBuffer.WriteString(" WHERE ")
	for i, pk := range pks.Columns() {
		queryBodyBuffer.WriteString(pk.ColumnName)
		queryBodyBuffer.WriteString(" = ")
		queryBodyBuffer.WriteString("$" + strconv.Itoa(i+1))

		if i < (pks.Len() - 1) {
			queryBodyBuffer.WriteString(" AND ")
		}
	}
	statementLabel = tableName + "_delete"
	queryBody = queryBodyBuffer.String()
	i.conn.bufferPool.Put(queryBodyBuffer)
	return
}

func (i *Invocation) generateTruncate(object DatabaseMapped) (statmentLabel, queryBody string) {
	tableName := TableName(object)

	queryBodyBuffer := i.conn.bufferPool.Get()
	queryBodyBuffer.WriteString("TRUNCATE ")
	queryBodyBuffer.WriteString(tableName)

	queryBody = queryBodyBuffer.String()
	statmentLabel = tableName + "_truncate"
	i.conn.bufferPool.Put(queryBodyBuffer)
	return
}

// --------------------------------------------------------------------------------
// helpers
// --------------------------------------------------------------------------------

func (i *Invocation) autoValues(autos *ColumnCollection) []interface{} {
	autoValues := make([]interface{}, autos.Len())
	for i, autoCol := range autos.Columns() {
		autoValues[i] = reflect.New(reflect.PtrTo(autoCol.FieldType)).Interface()
	}
	return autoValues
}

func (i *Invocation) setAutos(object DatabaseMapped, autos *ColumnCollection, autoValues []interface{}) (err error) {
	for index := 0; index < len(autoValues); index++ {
		err = autos.Columns()[index].SetValue(object, autoValues[index])
		if err != nil {
			err = exception.New(err)
			return
		}
	}
	return
}

func (i *Invocation) closeStatement(stmt *sql.Stmt, err error) error {
	if stmt == nil {
		return err
	}
	if i.tx != nil || i.conn.statementCache == nil || !i.conn.statementCache.Enabled() || i.statementLabel == "" {
		return exception.Nest(err, stmt.Close())
	}
	return err
}

func (i *Invocation) start(statement string) {
	if i.tracer != nil {
		i.traceFinisher = i.tracer.Query(i.context, i.conn, i, statement)
	}
}

func (i *Invocation) finish(statement string, r interface{}, err error) error {
	if i.cancel != nil {
		i.cancel()
	}
	if r != nil {
		err = exception.Nest(err, exception.New(r))
	}
	if i.conn.log != nil {
		i.conn.log.Trigger(
			logger.NewQueryEvent(statement, time.Now().UTC().Sub(i.startTime)).
				WithUsername(i.conn.config.GetUsername()).
				WithDatabase(i.conn.config.GetDatabase()).
				WithQueryLabel(i.statementLabel).
				WithEngine(i.conn.config.GetEngine()).
				WithErr(err),
		)
	}
	if i.traceFinisher != nil {
		i.traceFinisher.Finish(err)
	}
	if err != nil {
		err = exception.New(err)
	}
	return err
}
