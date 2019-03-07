package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"sync"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/uuid"
)

// TestConnectionSanityCheck tests if we can connect to the db, a.k.a., if the underlying driver works.
func TestConnectionUseBeforeOpen(t *testing.T) {
	assert := assert.New(t)

	conn, err := NewFromEnv()
	assert.Nil(err)

	tx, err := conn.Begin()
	assert.NotNil(err)
	assert.True(exception.Is(ErrConnectionClosed, err))
	assert.Nil(tx)
}

// TestConnectionSanityCheck tests if we can connect to the db, a.k.a., if the underlying driver works.
func TestConnectionSanityCheck(t *testing.T) {
	assert := assert.New(t)

	conn, err := NewFromEnv()
	assert.Nil(err)
	str := conn.Config().CreateDSN()
	_, err = sql.Open("postgres", str)
	assert.Nil(err)
}

func TestPrepare(t *testing.T) {
	a := assert.New(t)
	tx, err := Default().Begin()
	a.Nil(err)
	defer tx.Rollback()

	err = createTable(tx)
	a.Nil(err)
}

func TestQuery(t *testing.T) {
	a := assert.New(t)
	tx, err := Default().Begin()
	a.Nil(err)
	defer tx.Rollback()

	err = seedObjects(100, tx)
	a.Nil(err)

	objs := []benchObj{}
	err = Default().QueryInTx("select * from bench_object", tx).OutMany(&objs)
	a.Nil(err)
	a.NotEmpty(objs)

	var all []benchObj
	err = Default().GetAllInTx(&all, tx)
	a.Nil(err)
	a.Equal(len(objs), len(all))

	obj := benchObj{}
	err = Default().QueryInTx("select * from bench_object limit 1", tx).Out(&obj)
	a.Nil(err)
	a.NotEqual(obj.ID, 0)

	var id int
	err = Default().QueryInTx("select id from bench_object limit 1", tx).Scan(&id)
	a.Nil(err)
	a.NotEqual(id, 0)
}

func TestConnectionStatementCacheExecute(t *testing.T) {
	a := assert.New(t)

	conn, err := NewFromEnv()
	a.Nil(err)
	a.Nil(conn.Open())
	defer conn.Close()
	conn.PlanCache().WithEnabled(true)

	a.Nil(conn.Exec("select 'ok!'"))
	a.Nil(conn.Exec("select 'ok!'"))
	a.False(conn.PlanCache().HasStatement("select 'ok!'"))

	a.Nil(conn.ExecWithCachedPlan("select 'ok!'", "ping"))
	a.Nil(conn.ExecWithCachedPlan("select 'ok!'", "ping"))
	a.True(conn.PlanCache().HasStatement("ping"))
}

func TestConnectionStatementCacheQuery(t *testing.T) {
	a := assert.New(t)

	conn, err := NewFromEnv()
	a.Nil(err)
	a.Nil(conn.Open())
	defer conn.Close()

	conn.PlanCache().WithEnabled(true)

	var ok string
	a.Nil(conn.Invoke(context.TODO()).WithCachedPlan("status").Query("select 'ok!'").Scan(&ok))
	a.Equal("ok!", ok)

	a.Nil(conn.Invoke(context.TODO()).WithCachedPlan("status").Query("select 'ok!'").Scan(&ok))
	a.Equal("ok!", ok)

	a.True(conn.PlanCache().HasStatement("status"))
}

func TestCRUDMethods(t *testing.T) {
	a := assert.New(t)
	tx, err := Default().Begin()
	a.Nil(err)
	defer tx.Rollback()

	seedErr := seedObjects(100, tx)
	a.Nil(seedErr)

	objs := []benchObj{}
	queryErr := Default().QueryInTx("select * from bench_object", tx).OutMany(&objs)

	a.Nil(queryErr)
	a.NotEmpty(objs)

	all := []benchObj{}
	allErr := Default().GetAllInTx(&all, tx)
	a.Nil(allErr)
	a.Equal(len(objs), len(all))

	sampleObj := all[0]

	getTest := benchObj{}
	getTestErr := Default().GetInTx(&getTest, tx, sampleObj.ID)
	a.Nil(getTestErr)
	a.Equal(sampleObj.ID, getTest.ID)
	a.NotEmpty(getTest.UUID)

	exists, existsErr := Default().ExistsInTx(&getTest, tx)
	a.Nil(existsErr)
	a.True(exists)

	getTest.Name = "not_a_test_object"

	updateErr := Default().UpdateInTx(&getTest, tx)
	a.Nil(updateErr)

	verify := benchObj{}
	verifyErr := Default().GetInTx(&verify, tx, getTest.ID)
	a.Nil(verifyErr)
	a.Equal(getTest.Name, verify.Name)

	deleteErr := Default().DeleteInTx(&verify, tx)
	a.Nil(deleteErr)

	delVerify := benchObj{}
	delVerifyErr := Default().GetInTx(&delVerify, tx, getTest.ID)
	a.Nil(delVerifyErr)
}

func TestCRUDMethodsCached(t *testing.T) {
	a := assert.New(t)

	conn, err := NewFromEnv()
	conn.PlanCache().WithEnabled(true)
	a.Nil(err)
	a.Nil(conn.Open())
	defer conn.Close()

	tx, err := conn.Begin()
	a.Nil(err)
	defer tx.Rollback()

	err = seedObjects(100, tx)
	a.Nil(err)

	objs := []benchObj{}
	a.Nil(Default().QueryInTx("select * from bench_object", tx).OutMany(&objs))
	a.NotEmpty(objs)

	all := []benchObj{}
	a.Nil(Default().GetAllInTx(&all, tx))
	a.Equal(len(objs), len(all))

	sampleObj := all[0]

	getTest := benchObj{}
	a.Nil(Default().GetInTx(&getTest, tx, sampleObj.ID))
	a.Equal(sampleObj.ID, getTest.ID)

	exists, existsErr := Default().ExistsInTx(&getTest, tx)
	a.Nil(existsErr)
	a.True(exists)

	getTest.Name = "not_a_test_object"

	updateErr := Default().UpdateInTx(&getTest, tx)
	a.Nil(updateErr)

	verify := benchObj{}
	verifyErr := Default().GetInTx(&verify, tx, getTest.ID)
	a.Nil(verifyErr)
	a.Equal(getTest.Name, verify.Name)

	deleteErr := Default().DeleteInTx(&verify, tx)
	a.Nil(deleteErr)

	delVerify := benchObj{}
	delVerifyErr := Default().GetInTx(&delVerify, tx, getTest.ID)
	a.Nil(delVerifyErr)
}

func TestConnectionOpen(t *testing.T) {
	a := assert.New(t)

	conn, err := NewFromEnv()
	a.Nil(err)
	a.Nil(conn.Open())
	defer conn.Close()

	a.NotNil(conn.bufferPool)
	a.NotNil(conn.connection)
	a.NotNil(conn.planCache)
}

func TestExec(t *testing.T) {
	a := assert.New(t)
	tx, err := Default().Begin()
	a.Nil(err)
	defer tx.Rollback()

	err = Default().ExecInTx("select 'ok!'", tx)
	a.Nil(err)
}

func TestConnectionCreate(t *testing.T) {
	assert := assert.New(t)
	tx, err := Default().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	err = createTable(tx)
	assert.Nil(err)

	obj := &benchObj{
		Name:      fmt.Sprintf("test_object_0"),
		UUID:      uuid.V4().String(),
		Timestamp: time.Now().UTC(),
		Amount:    1000.0 + (5.0 * float32(0)),
		Pending:   true,
		Category:  fmt.Sprintf("category_%d", 0),
	}
	err = Default().CreateInTx(obj, tx)
	assert.Nil(err)
}

func TestConnectionCreateParallel(t *testing.T) {
	assert := assert.New(t)

	err := createTable(nil)
	assert.Nil(err)
	defer dropTableIfExists(nil)

	wg := sync.WaitGroup{}
	wg.Add(5)
	for x := 0; x < 5; x++ {
		go func() {
			defer wg.Done()
			obj := &benchObj{
				Name:      fmt.Sprintf("test_object_0"),
				UUID:      uuid.V4().String(),
				Timestamp: time.Now().UTC(),
				Amount:    1000.0 + (5.0 * float32(0)),
				Pending:   true,
				Category:  fmt.Sprintf("category_%d", 0),
			}
			innerErr := Default().CreateInTx(obj, nil)
			assert.Nil(innerErr)
		}()
	}
	wg.Wait()
}

func TestConnectionUpsert(t *testing.T) {
	assert := assert.New(t)
	tx, err := Default().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	err = createUpserObjectTable(tx)
	assert.Nil(err)

	obj := &upsertObj{
		UUID:      uuid.V4().String(),
		Timestamp: time.Now().UTC(),
		Category:  uuid.V4().String(),
	}
	err = Default().UpsertInTx(obj, tx)
	assert.Nil(err)

	var verify upsertObj
	err = Default().GetInTx(&verify, tx, obj.UUID)
	assert.Nil(err)
	assert.Equal(obj.Category, verify.Category)

	obj.Category = "test"

	err = Default().UpsertInTx(obj, tx)
	assert.Nil(err)

	err = Default().GetInTx(&verify, tx, obj.UUID)
	assert.Nil(err)
	assert.Equal(obj.Category, verify.Category)
}

func TestConnectionUpsertWithSerial(t *testing.T) {
	assert := assert.New(t)
	tx, err := Default().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	err = createTable(tx)
	assert.Nil(err)

	obj := &benchObj{
		Name:      "test_object_0",
		UUID:      uuid.V4().String(),
		Timestamp: time.Now().UTC(),
		Amount:    1005.0,
		Pending:   true,
		Category:  "category_0",
	}
	err = Default().UpsertInTx(obj, tx)
	assert.Nil(err, fmt.Sprintf("%+v", err))
	assert.NotZero(obj.ID)

	var verify benchObj
	err = Default().GetInTx(&verify, tx, obj.ID)
	assert.Nil(err)
	assert.Equal(obj.Category, verify.Category)

	obj.Category = "test"

	err = Default().UpsertInTx(obj, tx)
	assert.Nil(err)
	assert.NotZero(obj.ID)

	err = Default().GetInTx(&verify, tx, obj.ID)
	assert.Nil(err)
	assert.Equal(obj.Category, verify.Category)
}

func TestConnectionCreateMany(t *testing.T) {
	assert := assert.New(t)
	tx, err := Default().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	err = createTable(tx)
	assert.Nil(err)

	var objects []DatabaseMapped
	for x := 0; x < 10; x++ {
		objects = append(objects, benchObj{
			Name:      fmt.Sprintf("test_object_%d", x),
			UUID:      uuid.V4().String(),
			Timestamp: time.Now().UTC(),
			Amount:    1005.0,
			Pending:   true,
			Category:  fmt.Sprintf("category_%d", x),
		})
	}

	err = Default().CreateManyInTx(objects, tx)
	assert.Nil(err)

	var verify []benchObj
	err = Default().QueryInTx(`select * from bench_object`, tx).OutMany(&verify)
	assert.Nil(err)
	assert.NotEmpty(verify)
}

func TestConnectionTruncate(t *testing.T) {
	assert := assert.New(t)
	tx, err := Default().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	err = createTable(tx)
	assert.Nil(err)

	var objects []DatabaseMapped
	for x := 0; x < 10; x++ {
		objects = append(objects, benchObj{
			Name:      fmt.Sprintf("test_object_%d", x),
			UUID:      uuid.V4().String(),
			Timestamp: time.Now().UTC(),
			Amount:    1005.0,
			Pending:   true,
			Category:  fmt.Sprintf("category_%d", x),
		})
	}

	err = Default().CreateManyInTx(objects, tx)
	assert.Nil(err)

	var count int
	err = Default().QueryInTx(`select count(*) from bench_object`, tx).Scan(&count)
	assert.Nil(err)
	assert.NotZero(count)

	err = Default().TruncateInTx(benchObj{}, tx)
	assert.Nil(err)

	err = Default().QueryInTx(`select count(*) from bench_object`, tx).Scan(&count)
	assert.Nil(err)
	assert.Zero(count)
}

func TestConnectionCreateIfNotExists(t *testing.T) {
	assert := assert.New(t)
	tx, err := Default().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	err = createUpserObjectTable(tx)
	assert.Nil(err)

	obj := &upsertObj{
		UUID:      uuid.V4().String(),
		Timestamp: time.Now().UTC(),
		Category:  uuid.V4().String(),
	}
	err = Default().CreateIfNotExistsInTx(obj, tx)
	assert.Nil(err)

	var verify upsertObj
	err = Default().GetInTx(&verify, tx, obj.UUID)
	assert.Nil(err)
	assert.Equal(obj.Category, verify.Category)

	oldCategory := obj.Category
	obj.Category = "test"

	err = Default().CreateIfNotExistsInTx(obj, tx)
	assert.Nil(err)

	err = Default().GetInTx(&verify, tx, obj.UUID)
	assert.Nil(err)
	assert.Equal(oldCategory, verify.Category)
}

func TestConnectionInvalidatesBadCachedStatements(t *testing.T) {
	assert := assert.New(t)

	conn, err := NewFromEnv()
	assert.Nil(err)
	assert.Nil(conn.Open())
	defer conn.Close()

	conn.PlanCache().WithEnabled(true)

	createTableStatement := `CREATE TABLE state_invalidation (id int not null, name varchar(64))`
	insertStatement := `INSERT INTO state_invalidation (id, name) VALUES ($1, $2)`
	alterTableStatement := `ALTER TABLE state_invalidation ALTER COLUMN id TYPE bigint;`
	dropTableStatement := `DROP TABLE state_invalidation`
	queryStatement := `SELECT * from state_invalidation`

	defer func() {
		err = conn.Exec(dropTableStatement)
		assert.Nil(err)
	}()

	err = conn.Exec(createTableStatement)
	assert.Nil(err)

	err = conn.Exec(insertStatement, 1, "Foo")
	assert.Nil(err)

	err = conn.Exec(insertStatement, 2, "Bar")
	assert.Nil(err)

	_, err = conn.Query(queryStatement).Any()
	assert.Nil(err)

	err = conn.Exec(alterTableStatement)
	assert.Nil(err)

	// normally this would result in a busted cached query plan.
	// we need to invalidate the cache and make this work.
	_, err = conn.Query(queryStatement).Any()
	assert.Nil(err)

	_, err = conn.Query(queryStatement).Any()
	assert.Nil(err)
}

// TestConnectionConfigSetsDatabase tests if we set the .database property on open.
func TestConnectionConfigSetsDatabase(t *testing.T) {
	assert := assert.New(t)
	conn, err := NewFromEnv()
	assert.Nil(err)
	assert.Nil(conn.Open())
	defer conn.Close()
	assert.NotEmpty(conn.Config().GetDatabase())
}
