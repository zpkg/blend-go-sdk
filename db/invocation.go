package db

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/blend/go-sdk/exception"
)

const (
	connectionErrorMessage = "invocation context; db connection is nil"
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

// Start returns the invocation start time.
func (i *Invocation) Start() time.Time {
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
	if i.conn.StatementCache().Enabled() && len(i.statementLabel) > 0 {
		return i.conn.PrepareCachedContext(i.Context(), i.statementLabel, statement, i.tx)
	}
	return i.conn.PrepareContext(i.Context(), statement, i.tx)
}

// Exec executes a sql statement with a given set of arguments.
func (i *Invocation) Exec(statement string, args ...interface{}) (err error) {
	err = i.validate()
	if err != nil {
		return
	}

	i.start(statement)
	defer func() { err = i.finish(statement, recover(), err) }()

	stmt, stmtErr := i.Prepare(statement)
	if stmtErr != nil {
		err = exception.New(stmtErr)
		return
	}

	defer func() { err = i.closeStatement(err, stmt) }()

	if _, execErr := stmt.Exec(args...); execErr != nil {
		err = exception.New(execErr)
		if err != nil {
			i.invalidateCachedStatement()
		}
		return
	}

	return
}

// Query returns a new query object for a given sql query and arguments.
func (i *Invocation) Query(statement string, args ...interface{}) *Query {
	stmt, err := i.Prepare(statement)
	i.start(statement)
	return &Query{
		stmt:           stmt,
		err:            err,
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
	err = i.validate()
	if err != nil {
		return
	}

	if len(ids) == 0 {
		return exception.New(ErrInvalidIDs)
	}

	var queryBody string
	meta := getCachedColumnCollectionFromInstance(object)
	standardCols := meta.NotReadOnly()
	tableName := TableName(object)
	if len(i.statementLabel) == 0 {
		i.statementLabel = fmt.Sprintf("%s_get", tableName)
	}

	defer func() { err = i.finish(queryBody, recover(), err) }()

	columnNames := standardCols.ColumnNames()
	pks := standardCols.PrimaryKeys()
	if pks.Len() == 0 {
		err = exception.New(ErrNoPrimaryKey)
		return
	}

	queryBodyBuffer := i.conn.bufferPool.Get()
	defer i.conn.bufferPool.Put(queryBodyBuffer)

	queryBodyBuffer.WriteString("SELECT ")
	for i, name := range columnNames {
		queryBodyBuffer.WriteString(name)
		if i < (len(columnNames) - 1) {
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

	queryBody = queryBodyBuffer.String()

	stmt, stmtErr := i.Prepare(queryBody)
	if stmtErr != nil {
		err = exception.New(stmtErr)
		return
	}
	defer i.closeStatement(err, stmt)

	i.start(queryBody)
	rows, queryErr := stmt.QueryContext(i.Context(), ids...)

	if queryErr != nil {
		err = exception.New(queryErr)
		i.invalidateCachedStatement()
		return
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			err = exception.Nest(err, closeErr)
		}
	}()

	var popErr error
	if rows.Next() {
		if isPopulatable(object) {
			popErr = asPopulatable(object).Populate(rows)
		} else {
			popErr = PopulateInOrder(object, rows, standardCols)
		}

		if popErr != nil {
			err = exception.New(popErr)
			return
		}
	}

	err = exception.New(rows.Err())
	return
}

// GetAll returns all rows of an object mapped table wrapped in a transaction.
func (i *Invocation) GetAll(collection interface{}) (err error) {
	err = i.validate()
	if err != nil {
		return
	}

	var queryBody string
	defer func() { err = i.finish(queryBody, recover(), err) }()

	collectionValue := reflectValue(collection)
	t := reflectSliceType(collection)
	tableName := TableNameByType(t)

	if len(i.statementLabel) == 0 {
		i.statementLabel = fmt.Sprintf("%s_get_all", tableName)
	}

	meta := getCachedColumnCollectionFromType(tableName, t).NotReadOnly()

	columnNames := meta.ColumnNames()

	queryBodyBuffer := i.conn.bufferPool.Get()
	defer i.conn.bufferPool.Put(queryBodyBuffer)

	queryBodyBuffer.WriteString("SELECT ")
	for i, name := range columnNames {
		queryBodyBuffer.WriteString(name)
		if i < (len(columnNames) - 1) {
			queryBodyBuffer.WriteRune(runeComma)
		}
	}
	queryBodyBuffer.WriteString(" FROM ")
	queryBodyBuffer.WriteString(tableName)

	queryBody = queryBodyBuffer.String()
	stmt, stmtErr := i.Prepare(queryBody)
	if stmtErr != nil {
		err = exception.New(stmtErr)
		i.invalidateCachedStatement()
		return
	}

	defer func() { err = i.closeStatement(err, stmt) }()

	i.start(queryBody)

	rows, queryErr := stmt.QueryContext(i.Context())
	if queryErr != nil {
		err = exception.New(queryErr)
		return
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			err = exception.Nest(err, closeErr)
		}
	}()

	v, err := makeNewDatabaseMapped(t)
	if err != nil {
		return
	}
	isPopulatable := isPopulatable(v)

	var popErr error
	for rows.Next() {
		newObj, _ := makeNewDatabaseMapped(t)

		if isPopulatable {
			popErr = asPopulatable(newObj).Populate(rows)
		} else {
			popErr = PopulateInOrder(newObj, rows, meta)
			if popErr != nil {
				err = exception.New(popErr)
				return
			}
		}
		newObjValue := reflectValue(newObj)
		collectionValue.Set(reflect.Append(collectionValue, newObjValue))
	}

	err = exception.New(rows.Err())
	return
}

// Create writes an object to the database within a transaction.
func (i *Invocation) Create(object DatabaseMapped) (err error) {
	err = i.validate()
	if err != nil {
		return
	}

	var queryBody string
	defer func() { err = i.finish(queryBody, recover(), err) }()

	cols := getCachedColumnCollectionFromInstance(object)
	writeCols := cols.NotReadOnly().NotAutos()

	autos := cols.Autos()
	tableName := TableName(object)

	if len(i.statementLabel) == 0 {
		i.statementLabel = fmt.Sprintf("%s_create", tableName)
	}

	colNames := writeCols.ColumnNames()
	colValues := writeCols.ColumnValues(object)

	queryBodyBuffer := i.conn.bufferPool.Get()
	defer i.conn.bufferPool.Put(queryBodyBuffer)

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

	if autos.Len() > 0 {
		queryBodyBuffer.WriteString(" RETURNING ")
		queryBodyBuffer.WriteString(autos.ColumnNamesCSV())
	}

	queryBody = queryBodyBuffer.String()
	stmt, stmtErr := i.Prepare(queryBody)
	if stmtErr != nil {
		err = exception.New(stmtErr)
		return
	}
	defer func() { err = i.closeStatement(err, stmt) }()

	i.start(queryBody)

	var execErr error
	if autos.Len() == 0 {
		if i.context != nil {
			_, execErr = stmt.ExecContext(i.context, colValues...)
		} else {
			_, execErr = stmt.Exec(colValues...)
		}

		if execErr != nil {
			err = exception.New(execErr)
			i.invalidateCachedStatement()
			return
		}
	} else {
		autoValues := make([]interface{}, autos.Len())
		for i, autoCol := range autos.Columns() {
			autoValues[i] = reflect.New(reflect.PtrTo(autoCol.FieldType)).Interface()
		}

		if i.context != nil {
			execErr = stmt.QueryRowContext(i.context, colValues...).Scan(autoValues...)
		} else {
			execErr = stmt.QueryRow(colValues...).Scan(autoValues...)
		}

		if execErr != nil {
			err = exception.New(execErr)
			i.invalidateCachedStatement()
			return
		}

		for index := 0; index < len(autoValues); index++ {
			setErr := autos.Columns()[index].SetValue(object, autoValues[index])
			if setErr != nil {
				err = exception.New(setErr)
				return
			}
		}
	}

	return nil
}

// CreateIfNotExists writes an object to the database if it does not already exist within a transaction.
func (i *Invocation) CreateIfNotExists(object DatabaseMapped) (err error) {
	err = i.validate()
	if err != nil {
		return
	}

	var queryBody string
	defer func() { err = i.finish(queryBody, recover(), err) }()

	cols := getCachedColumnCollectionFromInstance(object)
	writeCols := cols.NotReadOnly().NotAutos()

	//NOTE: we're only using one.
	autos := cols.Autos()
	pks := cols.PrimaryKeys()
	tableName := TableName(object)

	if len(i.statementLabel) == 0 {
		i.statementLabel = fmt.Sprintf("%s_create_if_not_exists", tableName)
	}

	colNames := writeCols.ColumnNames()
	colValues := writeCols.ColumnValues(object)

	queryBodyBuffer := i.conn.bufferPool.Get()
	defer i.conn.bufferPool.Put(queryBodyBuffer)

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
	stmt, stmtErr := i.Prepare(queryBody)
	if stmtErr != nil {
		err = exception.New(stmtErr)
		return
	}
	defer func() { err = i.closeStatement(err, stmt) }()

	i.start(queryBody)

	var execErr error
	if autos.Len() == 0 {
		if i.context != nil {
			_, execErr = stmt.ExecContext(i.context, colValues...)
		} else {
			_, execErr = stmt.Exec(colValues...)
		}
		if execErr != nil {
			err = exception.New(execErr)
			i.invalidateCachedStatement()
			return
		}
	} else {
		autoValues := make([]interface{}, autos.Len())
		for i, autoCol := range autos.Columns() {
			autoValues[i] = reflect.New(reflect.PtrTo(autoCol.FieldType)).Interface()
		}

		if i.context != nil {
			execErr = stmt.QueryRowContext(i.context, colValues...).Scan(autoValues...)
		} else {
			execErr = stmt.QueryRow(colValues...).Scan(autoValues...)
		}

		if execErr != nil {
			err = exception.New(execErr).WithMessagef("query: %s", queryBody)
			return
		}

		for index := 0; index < len(autoValues); index++ {
			setErr := autos.Columns()[index].SetValue(object, autoValues[index])
			if setErr != nil {
				err = exception.New(setErr)
				return
			}
		}
	}

	return nil
}

// CreateMany writes many an objects to the database within a transaction.
func (i *Invocation) CreateMany(objects interface{}) (err error) {
	err = i.validate()
	if err != nil {
		return
	}

	var queryBody string
	defer func() { err = i.finish(queryBody, recover(), err) }()

	sliceValue := reflectValue(objects)
	if sliceValue.Len() == 0 {
		return nil
	}

	sliceType := reflectSliceType(objects)
	tableName := TableNameByType(sliceType)

	cols := getCachedColumnCollectionFromType(tableName, sliceType)
	writeCols := cols.NotReadOnly().NotAutos()

	//NOTE: we're only using one.
	//serials := cols.Serials()
	colNames := writeCols.ColumnNames()

	queryBodyBuffer := i.conn.bufferPool.Get()
	defer i.conn.bufferPool.Put(queryBodyBuffer)

	queryBodyBuffer.WriteString("INSERT INTO ")
	queryBodyBuffer.WriteString(tableName)
	queryBodyBuffer.WriteString(" (")
	for i, name := range colNames {
		queryBodyBuffer.WriteString(name)
		if i < len(colNames)-1 {
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
	stmt, stmtErr := i.Prepare(queryBody)
	if stmtErr != nil {
		err = exception.New(stmtErr)
		return
	}
	defer func() { err = i.closeStatement(err, stmt) }()
	i.start(queryBody)

	var colValues []interface{}
	for row := 0; row < sliceValue.Len(); row++ {
		colValues = append(colValues, writeCols.ColumnValues(sliceValue.Index(row).Interface())...)
	}

	var execErr error
	if i.context != nil {
		_, execErr = stmt.ExecContext(i.context, colValues...)
	} else {
		_, execErr = stmt.Exec(colValues...)
	}
	if execErr != nil {
		err = exception.New(execErr)
		i.invalidateCachedStatement()
		return
	}

	return nil
}

// Update updates an object wrapped in a transaction.
func (i *Invocation) Update(object DatabaseMapped) (err error) {
	err = i.validate()
	if err != nil {
		return
	}

	var queryBody string
	defer func() { err = i.finish(queryBody, recover(), err) }()

	tableName := TableName(object)
	if len(i.statementLabel) == 0 {
		i.statementLabel = fmt.Sprintf("%s_update", tableName)
	}

	cols := getCachedColumnCollectionFromInstance(object)
	writeCols := cols.WriteColumns()
	pks := cols.PrimaryKeys()
	updateCols := cols.UpdateColumns()
	updateValues := updateCols.ColumnValues(object)
	numColumns := writeCols.Len()

	queryBodyBuffer := i.conn.bufferPool.Get()
	defer i.conn.bufferPool.Put(queryBodyBuffer)

	queryBodyBuffer.WriteString("UPDATE ")
	queryBodyBuffer.WriteString(tableName)
	queryBodyBuffer.WriteString(" SET ")

	var writeColIndex int
	var col Column
	for ; writeColIndex < writeCols.Len(); writeColIndex++ {
		col = writeCols.columns[writeColIndex]
		queryBodyBuffer.WriteString(col.ColumnName)
		queryBodyBuffer.WriteString(" = $" + strconv.Itoa(writeColIndex+1))
		if writeColIndex != numColumns-1 {
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
	stmt, stmtErr := i.Prepare(queryBody)
	if stmtErr != nil {
		err = exception.New(stmtErr)
		return
	}

	defer func() { err = i.closeStatement(err, stmt) }()

	i.start(queryBody)

	var execErr error
	if i.context != nil {
		_, execErr = stmt.ExecContext(i.context, updateValues...)
	} else {
		_, execErr = stmt.Exec(updateValues...)
	}
	if execErr != nil {
		err = exception.New(execErr)
		i.invalidateCachedStatement()
		return
	}

	return
}

// Upsert inserts the object if it doesn't exist already (as defined by its primary keys) or updates it wrapped in a transaction.
func (i *Invocation) Upsert(object DatabaseMapped) (err error) {
	err = i.validate()
	if err != nil {
		return
	}

	var queryBody string
	defer func() { err = i.finish(queryBody, recover(), err) }()

	cols := getCachedColumnCollectionFromInstance(object)
	writeCols := cols.NotReadOnly().NotAutos()

	conflictUpdateCols := cols.NotReadOnly().NotAutos().NotPrimaryKeys()

	serials := cols.Autos()
	pks := cols.PrimaryKeys()
	tableName := TableName(object)

	if len(i.statementLabel) == 0 {
		i.statementLabel = fmt.Sprintf("%s_upsert", tableName)
	}

	colNames := writeCols.ColumnNames()
	colValues := writeCols.ColumnValues(object)

	queryBodyBuffer := i.conn.bufferPool.Get()
	defer i.conn.bufferPool.Put(queryBodyBuffer)

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

	var serial = serials.FirstOrDefault()
	if serials.Len() != 0 {
		queryBodyBuffer.WriteString(" RETURNING ")
		queryBodyBuffer.WriteString(serial.ColumnName)
	}

	queryBody = queryBodyBuffer.String()

	stmt, stmtErr := i.Prepare(queryBody)
	if stmtErr != nil {
		err = exception.New(stmtErr)
		i.invalidateCachedStatement()
		return
	}
	defer func() { err = i.closeStatement(err, stmt) }()

	i.start(queryBody)

	var execErr error
	if serials.Len() != 0 {
		var id interface{}
		if i.context != nil {
			execErr = stmt.QueryRowContext(i.context, colValues...).Scan(&id)
		} else {
			execErr = stmt.QueryRow(colValues...).Scan(&id)
		}
		if execErr != nil {
			err = exception.New(execErr)
			i.invalidateCachedStatement()
			return
		}
		setErr := serial.SetValue(object, id)
		if setErr != nil {
			err = exception.New(setErr)
			return
		}
	} else {
		if i.context != nil {
			_, execErr = stmt.ExecContext(i.context, colValues...)
		} else {
			_, execErr = stmt.Exec(colValues...)
		}
		if execErr != nil {
			err = exception.New(execErr).WithMessagef("query: %s", queryBody)
			return
		}
	}

	return nil
}

// Exists returns a bool if a given object exists (utilizing the primary key columns if they exist) wrapped in a transaction.
func (i *Invocation) Exists(object DatabaseMapped) (exists bool, err error) {
	err = i.validate()
	if err != nil {
		return
	}

	var queryBody string
	defer func() { err = i.finish(queryBody, recover(), err) }()

	tableName := TableName(object)
	if len(i.statementLabel) == 0 {
		i.statementLabel = fmt.Sprintf("%s_exists", tableName)
	}
	cols := getCachedColumnCollectionFromInstance(object)
	pks := cols.PrimaryKeys()

	if pks.Len() == 0 {
		err = exception.New("No primary key on object.")
		return
	}

	queryBodyBuffer := i.conn.bufferPool.Get()
	defer i.conn.bufferPool.Put(queryBodyBuffer)

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

	queryBody = queryBodyBuffer.String()
	stmt, stmtErr := i.Prepare(queryBody)
	if stmtErr != nil {
		err = exception.New(stmtErr)
		return
	}

	defer func() { err = i.closeStatement(err, stmt) }()
	i.start(queryBody)

	pkValues := pks.ColumnValues(object)
	var rows *sql.Rows
	var queryErr error
	if i.context != nil {
		rows, queryErr = stmt.QueryContext(i.context, pkValues...)
	} else {
		rows, queryErr = stmt.Query(pkValues...)
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			err = exception.Nest(err, closeErr)
		}
	}()

	if queryErr != nil {
		exists = false
		err = exception.New(queryErr)
		i.invalidateCachedStatement()
		return
	}

	exists = rows.Next()
	return
}

// Delete deletes an object from the database wrapped in a transaction.
func (i *Invocation) Delete(object DatabaseMapped) (err error) {
	err = i.validate()
	if err != nil {
		return
	}

	var queryBody string
	defer func() { err = i.finish(queryBody, recover(), err) }()

	tableName := TableName(object)

	if len(i.statementLabel) == 0 {
		i.statementLabel = fmt.Sprintf("%s_delete", tableName)
	}

	cols := getCachedColumnCollectionFromInstance(object)
	pks := cols.PrimaryKeys()

	if len(pks.Columns()) == 0 {
		err = exception.New("No primary key on object.")
		return
	}

	queryBodyBuffer := i.conn.bufferPool.Get()
	defer i.conn.bufferPool.Put(queryBodyBuffer)

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

	queryBody = queryBodyBuffer.String()
	stmt, stmtErr := i.Prepare(queryBody)
	if stmtErr != nil {
		err = exception.New(stmtErr)
		return
	}
	defer func() { err = i.closeStatement(err, stmt) }()
	i.start(queryBody)

	pkValues := pks.ColumnValues(object)

	var execErr error
	if i.context != nil {
		_, execErr = stmt.ExecContext(i.context, pkValues...)
	} else {
		_, execErr = stmt.Exec(pkValues...)
	}
	if execErr != nil {
		err = exception.New(execErr)
		i.invalidateCachedStatement()
	}
	return
}

// Truncate completely empties a table in a single command.
func (i *Invocation) Truncate(object DatabaseMapped) (err error) {
	err = i.validate()
	if err != nil {
		return
	}

	var queryBody string
	defer func() { err = i.finish(queryBody, recover(), err) }()

	tableName := TableName(object)

	if len(i.statementLabel) == 0 {
		i.statementLabel = fmt.Sprintf("%s_truncate", tableName)
	}

	queryBodyBuffer := i.conn.bufferPool.Get()
	defer i.conn.bufferPool.Put(queryBodyBuffer)

	queryBodyBuffer.WriteString("TRUNCATE ")
	queryBodyBuffer.WriteString(tableName)

	queryBody = queryBodyBuffer.String()
	stmt, stmtErr := i.Prepare(queryBody)
	if stmtErr != nil {
		err = exception.New(stmtErr)
		return
	}
	defer func() { err = i.closeStatement(err, stmt) }()

	i.start(queryBody)

	var execErr error
	if i.context != nil {
		_, execErr = stmt.ExecContext(i.context)
	} else {
		_, execErr = stmt.Exec()
	}

	if execErr != nil {
		err = exception.New(execErr)
		i.invalidateCachedStatement()
	}
	return
}

// --------------------------------------------------------------------------------
// helpers
// --------------------------------------------------------------------------------

// validate the invocation is ready
func (i *Invocation) validate() error {
	if i.conn == nil {
		return exception.New(connectionErrorMessage)
	}
	return nil
}

func (i *Invocation) invalidateCachedStatement() error {
	if i.conn.statementCache.Enabled() && len(i.statementLabel) > 0 && i.tx == nil {
		return i.conn.statementCache.InvalidateStatement(i.statementLabel)
	}
	return nil
}

func (i *Invocation) closeStatement(err error, stmt *sql.Stmt) error {
	if i.conn.StatementCache().Enabled() && len(i.statementLabel) > 0 && i.tx == nil {
		return err
	}

	return exception.Nest(err, stmt.Close())
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
	if i.traceFinisher != nil {
		i.traceFinisher.Finish(err)
	}
	i.conn.finish(i.context, statement, i.statementLabel, since(i.startTime), err)
	return err
}
