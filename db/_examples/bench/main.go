package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/db/migration"
	"github.com/blend/go-sdk/db/migration/pg"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stringutil"
)

type benchObject struct {
	ID   int    `db:"id,pk,auto"`
	Name string `db:"name"`
}

func createTable(tableName string, log logger.Log, conn *db.Connection) error {
	log.SyncInfof("creating %s", tableName)
	return migration.New(
		migration.Group(
			migration.Step(
				pg.TableNotExists(tableName),
				migration.Statements(
					fmt.Sprintf("CREATE TABLE %s (id serial not null primary key, name varchar(255))", tableName),
				),
			),
		),
	).WithLogger(log).Apply(conn)
}

func dropTable(tableName string, log logger.Log, conn *db.Connection) error {
	log.SyncInfof("dropping %s", tableName)
	return migration.New(
		migration.Group(
			migration.Step(
				pg.TableExists(tableName),
				migration.Statements(
					fmt.Sprintf("DROP TABLE %s", tableName),
				),
			),
		),
	).WithLogger(log).Apply(conn)
}

func maybeFatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func reportStats(log logger.Log, conn *db.Connection) {
	ticker := time.Tick(500 * time.Millisecond)
	for {
		select {
		case <-ticker:
			log.SyncInfof("[%d] connections currently open", conn.Connection().Stats().OpenConnections)
			log.SyncInfof("[%v] wait duration", conn.Connection().Stats().WaitDuration)
		}
	}
}

func main() {
	log := logger.All().WithDisabled(logger.Query)
	conn := db.MustNewFromEnv()

	if err := conn.Open(); err != nil {
		log.SyncFatalExit(err)
	}

	go reportStats(log, conn)

	tableName := strings.ToLower(stringutil.Random(stringutil.Letters, 12))

	maybeFatal(createTable(tableName, log, conn))
	defer func() { maybeFatal(dropTable(tableName, log, conn)) }()

	for x := 0; x < 1<<12; x++ {
		ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
		maybeFatal(conn.Invoke(ctx).Exec(fmt.Sprintf("INSERT INTO %s VALUES ($1)", tableName), strconv.Itoa(x)))
		cancel()
	}

	wg := sync.WaitGroup{}
	wg.Add(4)
	for routine := 0; routine < 4; routine++ {
		go func() {
			defer wg.Done()
			for x := 0; x < 1<<10; x++ {
				ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
				if _, err := conn.Invoke(ctx).Query(fmt.Sprintf("select * from %s", tableName)).Any(); err != nil {
					maybeFatal(err)
				}
				cancel()
			}
		}()
	}
	wg.Wait()

	maybeFatal(log.Drain())
	log.SyncInfof("OK")
}
