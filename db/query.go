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

	stmt *sql.Stmt
	conn *Connection
	inv  *Invocation
	tx   *sql.Tx
	err  error
}

func (q *Query) exec() (rows *sql.Rows, err error) {
	if q.err != nil {
		err = q.err
		return
	}
	rows, err = q.stmt.QueryContext(q.context, q.args...)
	if err != nil {
		q.inv.invalidateCachedStatement()
		err = exception.New(err)
	}
	return
}

// Execute runs a given query, yielding the raw results.
func (q *Query) Execute() (rows *sql.Rows, err error) {
	defer func() { q.finalizer(recover(), err) }()
	rows, err = q.exec()
	return
}

// Any returns if there are any results for the query.
func (q *Query) Any() (hasRows bool, err error) {
	defer func() { err = q.finalizer(recover(), err) }()

	var rows *sql.Rows
	rows, err = q.exec()
	if err != nil {
		return
	}
	defer rows.Close()

	rowsErr := rows.Err()
	if rowsErr != nil {
		err = exception.New(rowsErr)
		return
	}

	hasRows = rows.Next()
	return
}

// None returns if there are no results for the query.
func (q *Query) None() (hasRows bool, err error) {
	defer func() { err = q.finalizer(recover(), err) }()

	var rows *sql.Rows
	rows, err = q.exec()
	if err != nil {
		return
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		err = exception.New(err)
		return
	}

	hasRows = !rows.Next()
	return
}

// Scan writes the results to a given set of local variables.
func (q *Query) Scan(args ...interface{}) (err error) {
	defer func() { err = q.finalizer(recover(), err) }()

	var rows *sql.Rows
	rows, err = q.exec()
	if err != nil {
		return
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		err = exception.New(err)
		return
	}

	if rows.Next() {
		err = rows.Scan(args...)
		if err != nil {
			err = exception.New(err)
			return
		}
	}

	return
}

// Out writes the query result to a single object via. reflection mapping.
func (q *Query) Out(object interface{}) (err error) {
	defer func() { err = q.finalizer(recover(), err) }()

	var rows *sql.Rows
	rows, err = q.exec()
	if err != nil {
		return
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		err = exception.New(err)
		return
	}

	sliceType := reflectType(object)
	if sliceType.Kind() != reflect.Struct {
		err = exception.New("destination object is not a struct")
		return
	}

	columnMeta := getCachedColumnCollectionFromInstance(object)
	if rows.Next() {
		if populatable, ok := object.(Populatable); ok {
			err = populatable.Populate(rows)
		} else {
			err = PopulateByName(object, rows, columnMeta)
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

	var rows *sql.Rows
	rows, err = q.exec()
	if err != nil {
		return err
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		err = exception.New(err)
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
	for rows.Next() {
		newObj := makeNew(sliceInnerType)

		if isPopulatable {
			err = asPopulatable(newObj).Populate(rows)
		} else {
			err = PopulateByName(newObj, rows, meta)
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

	var rows *sql.Rows
	rows, err = q.exec()
	if err != nil {
		return
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		err = exception.New(err)
		return
	}

	for rows.Next() {
		err = consumer(rows)
		if err != nil {
			err = exception.New(err)
			return
		}
	}
	return
}

// First executes the consumer for the first result of a query.
func (q *Query) First(consumer RowsConsumer) (err error) {
	defer func() { err = q.finalizer(recover(), err) }()

	var rows *sql.Rows
	rows, err = q.exec()
	if err != nil {
		return
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		err = exception.New(err)
		return
	}

	if rows.Next() {
		err = consumer(rows)
		if err != nil {
			return
		}
	}
	return
}

// --------------------------------------------------------------------------------
// helpers
// --------------------------------------------------------------------------------

func (q *Query) finalizer(r interface{}, err error) error {
	if r != nil {
		err = exception.Nest(err, exception.New(r))
	}
	// close the statement if it's set using the invocation
	if q.stmt != nil {
		err = q.inv.closeStatement(err, q.stmt)
	}
	// call the invocation finisher
	q.inv.finish(q.statement, nil, err)
	return err
}
