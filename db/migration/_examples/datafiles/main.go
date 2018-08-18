package main

import (
	"database/sql"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/db/migration"
	"github.com/blend/go-sdk/logger"
)

func main() {
	log := logger.All()
	conn := db.NewFromEnv()
	if err := conn.Open(); err != nil {
		logger.FatalExit(err)
	}

	err := migration.New(
		migration.NewGroup(
			migration.NewStep(
				migration.TableExists("test_vocab"),
				migration.Statements(
					"DROP TABLE test_vocab",
				),
			),
			migration.NewStep(
				migration.TableNotExists("test_vocab"),
				migration.Statements(
					"CREATE TABLE test_vocab (id serial not null, word varchar(32) not null);",
					"ALTER TABLE test_vocab ADD CONSTRAINT pk_test_vocab_id PRIMARY KEY(id);",
				),
			),
			migration.ReadDataFile("data.sql"),
			migration.NewStep(
				migration.Guard("test custom step", func(c *db.Connection, tx *sql.Tx) (bool, error) {
					return c.QueryInTx("select 1 from test_vocab where word = $1", tx, "foo").None()
				}),
				migration.Actions(func(c *db.Connection, tx *sql.Tx) error {
					return c.ExecInTx("insert into test_vocab (word) values ($1)", tx, "foo")
				}),
			),
			migration.NewStep(
				migration.TableExists("test_vocab"),
				migration.Statements(
					"DROP TABLE test_vocab",
				),
			),
		),
	).WithLogger(log).Apply(conn)

	if err != nil {
		log.SyncFatalExit(err)
	}
}
