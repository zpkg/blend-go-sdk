package db

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/blend/go-sdk/exception"
)

// --------------------------------------------------------------------------------
// Query Result
// --------------------------------------------------------------------------------

// Query is the intermediate result of a query.
type Query struct {
	context        context.Context
	statement      string
	statementLabel string
	args           []interface{}

	rows *sql.Rows
	err  error

	stmt *sql.Stmt
	conn *Connection
	inv  *Invocation
	tx   *sql.Tx
}

// Prepare prepares a statement query.
func (q *Query) Prepare() *Query {
	if q.err != nil {
		return q
	}
	q.stmt, q.err = q.inv.Prepare(q.statement)
	return q
}

// Close finishes a query.
func (q *Query) Close() (err error) {
	if q.err != nil {
		err = q.err
		return
	}
	if finishErr := q.inv.finish(q.statement, nil, err, q.stmt, q.rows); finishErr != nil {
		err = exception.Nest(err, finishErr)
	}
	return
}

// Execute runs a given query, yielding the raw results.
func (q *Query) Execute() (rows *sql.Rows, err error) {
	defer func() { err = q.finalizer(recover(), err) }()
	rows, err = q.query()
	return
}

// Any returns if there are any results for the query.
func (q *Query) Any() (hasRows bool, err error) {
	defer func() { err = q.finalizer(recover(), err) }()

	q.rows, q.err = q.query()
	if q.err != nil {
		err = q.err
		return
	}
	if err = q.rows.Err(); err != nil {
		err = exception.New(err)
		return
	}

	hasRows = q.rows.Next()
	return
}

// None returns if there are no results for the query.
func (q *Query) None() (hasRows bool, err error) {
	defer func() { err = q.finalizer(recover(), err) }()

	q.rows, q.err = q.query()
	if q.err != nil {
		err = q.err
		return
	}
	if err = q.rows.Err(); err != nil {
		err = exception.New(err)
		return
	}
	hasRows = !q.rows.Next()
	return
}

// Scan writes the results to a given set of local variables.
func (q *Query) Scan(args ...interface{}) (err error) {
	defer func() { err = q.finalizer(recover(), err) }()

	q.rows, q.err = q.query()
	if q.err != nil {
		err = q.err
		return
	}

	if q.rows.Next() {
		if err = q.rows.Err(); err != nil {
			err = exception.New(err)
			return
		}
		if err = q.rows.Scan(args...); err != nil {
			err = exception.New(err)
			return
		}
	}

	return
}

// Out writes the query result to a single object via. reflection mapping.
func (q *Query) Out(object interface{}) (err error) {
	defer func() { err = q.finalizer(recover(), err) }()

	q.rows, q.err = q.query()
	if q.err != nil {
		err = q.err
		return
	}

	sliceType := reflectType(object)
	if sliceType.Kind() != reflect.Struct {
		err = exception.New("destination object is not a struct")
		return
	}

	columnMeta := getCachedColumnCollectionFromInstance(object)
	if q.rows.Next() {
		if rowsErr := q.rows.Err(); rowsErr != nil {
			err = exception.New(rowsErr)
			return
		}
		if populatable, ok := object.(Populatable); ok {
			err = populatable.Populate(q.rows)
		} else {
			err = PopulateByName(object, q.rows, columnMeta)
		}
		if err != nil {
			return
		}
	}

	return
}

// OutMany writes the query results to a slice of objects.
func (q *Query) OutMany(collection interface{}) (err error) {
	defer func() { err = q.finalizer(recover(), err) }()

	q.rows, q.err = q.query()
	if q.err != nil {
		err = q.err
		return
	}

	sliceType := reflectType(collection)
	if sliceType.Kind() != reflect.Slice {
		err = exception.New(ErrCollectionNotSlice)
		return
	}

	sliceInnerType := reflectSliceType(collection)
	collectionValue := reflectValue(collection)
	v := makeNew(sliceInnerType)
	meta := getCachedColumnCollectionFromType(newColumnCacheKey(sliceInnerType), sliceInnerType)

	isPopulatable := isPopulatable(v)

	didSetRows := false
	for q.rows.Next() {
		if rowsErr := q.rows.Err(); rowsErr != nil {
			err = exception.New(rowsErr)
			return
		}
		newObj := makeNew(sliceInnerType)

		if isPopulatable {
			err = asPopulatable(newObj).Populate(q.rows)
		} else {
			err = PopulateByName(newObj, q.rows, meta)
		}
		if err != nil {
			return
		}

		newObjValue := reflectValue(newObj)
		collectionValue.Set(reflect.Append(collectionValue, newObjValue))
		didSetRows = true
	}

	if !didSetRows {
		collectionValue.Set(reflect.MakeSlice(sliceType, 0, 0))
	}
	return
}

// Each executes the consumer for each result of the query (one to many).
func (q *Query) Each(consumer RowsConsumer) (err error) {
	defer func() { err = q.finalizer(recover(), err) }()

	q.rows, q.err = q.query()
	if q.err != nil {
		err = q.err
		return
	}

	for q.rows.Next() {
		if err = q.rows.Err(); err != nil {
			err = exception.New(err)
			return
		}
		if err = consumer(q.rows); err != nil {
			err = exception.New(err)
			return
		}
	}
	return
}

// First executes the consumer for the first result of a query.
func (q *Query) First(consumer RowsConsumer) (err error) {
	defer func() { err = q.finalizer(recover(), err) }()

	q.rows, q.err = q.query()
	if q.err != nil {
		err = q.err
		return
	}

	if q.rows.Next() {
		if err = q.rows.Err(); err != nil {
			err = exception.New(err)
			return
		}
		if err = consumer(q.rows); err != nil {
			return
		}
	}
	return
}

// ThenOut applies a row of the result set to the given object.
// It automatically advances to the next result set after.
func (q *Query) ThenOut(object interface{}) *Query {
	defer func() { q.recover(recover()) }()

	if q.rows == nil {
		q.rows, q.err = q.query()
	}
	if q.err != nil {
		return q
	}

	sliceType := reflectType(object)
	if sliceType.Kind() != reflect.Struct {
		q.err = exception.New("destination object is not a struct")
		return q
	}

	columnMeta := getCachedColumnCollectionFromInstance(object)
	if q.rows.Next() {
		if populatable, ok := object.(Populatable); ok {
			q.err = populatable.Populate(q.rows)
		} else {
			q.err = PopulateByName(object, q.rows, columnMeta)
		}
		if q.err != nil {
			return q
		}
	}
	q.rows.NextResultSet()
	return q
}

// ThenOutMany reads the results into a collection.
// It automatically advances to the next result set after.
func (q *Query) ThenOutMany(collection interface{}) *Query {
	defer func() { q.recover(recover()) }()

	if q.rows == nil {
		q.rows, q.err = q.query()
	}
	if q.err != nil {
		return q
	}
	sliceType := reflectType(collection)
	if sliceType.Kind() != reflect.Slice {
		q.err = exception.New(ErrCollectionNotSlice)
		return q
	}

	sliceInnerType := reflectSliceType(collection)
	collectionValue := reflectValue(collection)
	v := makeNew(sliceInnerType)
	meta := getCachedColumnCollectionFromType(newColumnCacheKey(sliceInnerType), sliceInnerType)

	isPopulatable := isPopulatable(v)

	didSetRows := false
	for q.rows.Next() {
		newObj := makeNew(sliceInnerType)
		if isPopulatable {
			q.err = asPopulatable(newObj).Populate(q.rows)
		} else {
			q.err = PopulateByName(newObj, q.rows, meta)
		}
		if q.err != nil {
			return q
		}

		newObjValue := reflectValue(newObj)
		collectionValue.Set(reflect.Append(collectionValue, newObjValue))
		didSetRows = true
	}

	if !didSetRows {
		collectionValue.Set(reflect.MakeSlice(sliceType, 0, 0))
	}
	q.rows.NextResultSet()
	return q
}

// ThenEach applies a consumer to each row in the result.
// It automatically advances to the next result set after.
func (q *Query) ThenEach(consumer RowsConsumer) *Query {
	defer func() { q.recover(recover()) }()

	if q.rows == nil {
		q.rows, q.err = q.query()
	}
	if q.err != nil {
		return q
	}

	var err error
	for q.rows.Next() {
		err = consumer(q.rows)
		if err != nil {
			q.err = exception.New(err)
			return q
		}
	}
	q.rows.NextResultSet()
	return q
}

// ThenFirst applies a consumer to each row in the result.
// It automatically advances to the next result set after.
func (q *Query) ThenFirst(consumer RowsConsumer) *Query {
	defer func() { q.recover(recover()) }()

	if q.rows == nil {
		q.rows, q.err = q.query()
	}
	if q.err != nil {
		return q
	}

	if q.rows.Next() {
		q.err = consumer(q.rows)
		if q.err != nil {
			return q
		}
	}
	q.rows.NextResultSet()
	return q
}

// --------------------------------------------------------------------------------
// helpers
// --------------------------------------------------------------------------------

func (q *Query) query() (rows *sql.Rows, err error) {
	if q.err != nil {
		err = q.err
		return
	}

	if q.stmt != nil {
		rows, err = q.stmt.QueryContext(q.context, q.args...)
	} else if q.inv.tx != nil {
		rows, err = q.inv.tx.QueryContext(q.context, q.statement, q.args...)
	} else {
		rows, err = q.inv.conn.connection.QueryContext(q.context, q.statement, q.args...)
	}

	if err != nil {
		if q.stmt != nil {
			err = q.inv.maybeCloseStatement(err, q.stmt)
		}
		err = exception.New(err)
	}
	return
}

func (q *Query) recover(r interface{}) {
	if r != nil {
		q.err = exception.Nest(q.err, exception.New(r))
	}
}

func (q *Query) finalizer(r interface{}, err error) error {
	err = exception.Nest(err, q.inv.finish(q.statement, r, err, q.stmt, q.rows))
	return err
}
