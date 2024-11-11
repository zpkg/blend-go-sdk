/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/zpkg/blend-go-sdk/db"
	"github.com/zpkg/blend-go-sdk/db/migration"
	"github.com/zpkg/blend-go-sdk/logger"
	"github.com/zpkg/blend-go-sdk/stringutil"
)

func createTable(tableName string, log logger.Log, conn *db.Connection) error {
	log.Infof("creating %s", tableName)
	action := migration.NewStep(
		migration.TableNotExists(tableName),
		migration.Statements(
			fmt.Sprintf("CREATE TABLE %s (id serial not null primary key, name varchar(255))", tableName)),
	)

	suite := migration.New(migration.OptGroups(migration.NewGroup(migration.OptGroupActions(action))))
	suite.Log = log
	return suite.Apply(context.TODO(), conn)
}

func dropTable(tableName string, log logger.Log, conn *db.Connection) error {
	log.Infof("dropping %s", tableName)
	action := migration.NewStep(
		migration.TableExists(tableName),
		migration.Statements(
			fmt.Sprintf("DROP TABLE %s", tableName),
		),
	)

	suite := migration.New(migration.OptGroups(migration.NewGroup(migration.OptGroupActions(action))))
	suite.Log = log
	return suite.Apply(context.TODO(), conn)
}

func maybeFatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func reportStats(log logger.Log, conn *db.Connection) {
	ticker := time.Tick(500 * time.Millisecond)
	for range ticker {
		log.Infof("[%d] connections currently open", conn.Connection.Stats().OpenConnections)
		log.Infof("[%v] wait duration", conn.Connection.Stats().WaitDuration)
	}
}

func main() {
	log := logger.All(logger.OptDisabled(db.QueryFlag))
	conn, err := db.New(db.OptConfigFromEnv())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if err := conn.Open(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	go reportStats(log, conn)

	tableName := strings.ToLower(stringutil.Random(stringutil.Letters, 12))

	maybeFatal(createTable(tableName, log, conn))
	defer func() { maybeFatal(dropTable(tableName, log, conn)) }()

	for x := 0; x < 1<<12; x++ {
		ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
		err := db.IgnoreExecResult(conn.Invoke(db.OptContext(ctx)).Exec(fmt.Sprintf("INSERT INTO %s VALUES ($1)", tableName), strconv.Itoa(x)))
		maybeFatal(err)
		cancel()
	}

	wg := sync.WaitGroup{}
	wg.Add(4)
	for routine := 0; routine < 4; routine++ {
		go func() {
			defer wg.Done()
			for x := 0; x < 1<<10; x++ {
				ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
				if _, err := conn.Invoke(db.OptContext(ctx)).Query(fmt.Sprintf("select * from %s", tableName)).Any(); err != nil {
					maybeFatal(err)
				}
				cancel()
			}
		}()
	}
	wg.Wait()

	log.DrainContext(context.TODO())
	log.Infof("OK")
}
