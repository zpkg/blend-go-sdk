package migration

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/stringutil"
)

// GuardFunc is a control for migration steps.
type GuardFunc func(context.Context, *db.Connection, *sql.Tx, Action) error

// --------------------------------------------------------------------------------
// Guards
// --------------------------------------------------------------------------------

// Guard returns a function that determines if a step in a group should run.
func Guard(description string, predicate func(c *db.Connection, tx *sql.Tx) (bool, error)) GuardFunc {
	return func(ctx context.Context, c *db.Connection, tx *sql.Tx, step Action) error {
		proceed, err := predicate(c, tx)
		if err != nil {
			if suite := GetContextSuite(ctx); suite != nil {
				return suite.Error(WithLabel(ctx, description), err)
			}
			return err
		}

		if !proceed {
			if suite := GetContextSuite(ctx); suite != nil {
				suite.Skipf(ctx, description)
			}
			return nil
		}

		err = step(ctx, c, tx)
		if err != nil {
			if suite := GetContextSuite(ctx); suite != nil {
				return suite.Error(WithLabel(ctx, description), err)
			}
			return err
		}
		if suite := GetContextSuite(ctx); suite != nil {
			suite.Applyf(ctx, description)
		}
		return nil
	}
}

// Always always runs a step.
func Always() GuardFunc {
	return Guard("always run", func(_ *db.Connection, _ *sql.Tx) (bool, error) { return true, nil })
}

// IfExists only runs the statement if the given item exists.
func IfExists(statement string) GuardFunc {
	return Guard("if exists run", func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return Exists(c, tx, statement)
	})
}

// IfNotExists only runs the statement if the given item doesn't exist.
func IfNotExists(statement string) GuardFunc {
	return Guard("if not exists run", func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return NotExists(c, tx, statement)
	})
}

// GuardPredicate wraps a Predicate in a GuardFunc
func GuardPredicate(description string, p Predicate, arg1 string) GuardFunc {
	return Guard(description, func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return p(c, tx, arg1)
	})
}

// GuardNotPredicate inverts a Predicate, and wraps that in a GuardFunc
func GuardNotPredicate(description string, p Predicate, arg1 string) GuardFunc {
	return Guard(description, func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return Not(p(c, tx, arg1))
	})
}

// GuardPredicate2 wraps a Predicate2 in a GuardFunc
func GuardPredicate2(description string, p Predicate2, arg1, arg2 string) GuardFunc {
	return Guard(description, func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return p(c, tx, arg1, arg2)
	})
}

// GuardNotPredicate2 inverts a Predicate2, and wraps that in a GuardFunc
func GuardNotPredicate2(description string, p Predicate2, arg1, arg2 string) GuardFunc {
	return Guard(description, func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return Not(p(c, tx, arg1, arg2))
	})
}

// GuardPredicate3 wraps a Predicate3 in a GuardFunc
func GuardPredicate3(description string, p Predicate3, arg1, arg2, arg3 string) GuardFunc {
	return Guard(description, func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return p(c, tx, arg1, arg2, arg3)
	})
}

// GuardNotPredicate3 inverts a Predicate3, and wraps that in a GuardFunc
func GuardNotPredicate3(description string, p Predicate3, arg1, arg2, arg3 string) GuardFunc {
	return Guard(description, func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return Not(p(c, tx, arg1, arg2, arg3))
	})
}


// Predicate is a function that evaluates based on a string param.
type Predicate func(*db.Connection, *sql.Tx, string) (bool, error)

// Predicate2 is a function that evaluates based on two string params.
type Predicate2 func(*db.Connection, *sql.Tx, string, string) (bool, error)

// Predicate3 is a function that evaluates based on three string params.
type Predicate3 func(*db.Connection, *sql.Tx, string, string, string) (bool, error)

// Not inverts the output of a predicate.
func Not(proceed bool, err error) (bool, error) {
	return !proceed, err
}

// --------------------------------------------------------------------------------
// Guard Helpers
// --------------------------------------------------------------------------------

// Exists returns if a statement has results.
func Exists(c *db.Connection, tx *sql.Tx, selectStatement string) (bool, error) {
	if !stringutil.HasPrefixCaseless(selectStatement, "select") {
		return false, fmt.Errorf("statement must be a `SELECT`")
	}
	return c.Invoke(db.OptTx(tx)).Query(selectStatement).Any()
}

// NotExists returns if a statement doesnt have results.
func NotExists(c *db.Connection, tx *sql.Tx, selectStatement string) (bool, error) {
	if !stringutil.HasPrefixCaseless(selectStatement, "select") {
		return false, fmt.Errorf("statement must be a `SELECT`")
	}
	return c.Invoke(db.OptTx(tx)).Query(selectStatement).None()
}
