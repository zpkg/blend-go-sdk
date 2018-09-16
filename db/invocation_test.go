package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
)

func TestInvocationLabels(t *testing.T) {
	assert := assert.New(t)

	inv := &Invocation{}
	inv = inv.WithLabel("test")
	assert.NotEmpty(inv.Label())
}

type jsonTestChild struct {
	Label string `json:"label"`
}

type jsonTest struct {
	ID   int    `db:"id,pk,auto"`
	Name string `db:"name"`

	NotNull  jsonTestChild `db:"not_null,json"`
	Nullable []string      `db:"nullable,json"`
}

func (jt jsonTest) TableName() string {
	return "json_test"
}

func createJSONTestTable(tx *sql.Tx) error {
	return Default().Invoke(context.Background(), tx).Exec("create table json_test (id serial primary key, name varchar(255), not_null json, nullable json)")
}

func dropJSONTextTable(tx *sql.Tx) error {
	return Default().Invoke(context.Background(), tx).Exec("drop table if exists json_test")
}

func TestInvocationJSONNulls(t *testing.T) {
	assert := assert.New(t)
	tx, err := Default().Begin()
	assert.Nil(err)
	defer tx.Rollback()
	defer dropJSONTextTable(tx)

	assert.Nil(createJSONTestTable(tx))

	// try creating fully set object and reading it out
	obj0 := jsonTest{Name: uuid.V4().String(), NotNull: jsonTestChild{Label: uuid.V4().String()}, Nullable: []string{uuid.V4().String()}}
	assert.Nil(Default().Invoke(context.Background(), tx).Create(&obj0))

	var verify0 jsonTest
	assert.Nil(Default().Invoke(context.Background(), tx).Get(&verify0, obj0.ID))

	assert.Equal(obj0.ID, verify0.ID)
	assert.Equal(obj0.Name, verify0.Name)
	assert.Equal(obj0.Nullable, verify0.Nullable)
	assert.Equal(obj0.NotNull.Label, verify0.NotNull.Label)

	// try creating partially set object and reading it out
	obj1 := jsonTest{Name: uuid.V4().String(), NotNull: jsonTestChild{Label: uuid.V4().String()}} //note `Nullable` isn't set

	columns := Columns(obj1)
	values := columns.ColumnValues(obj1)
	assert.Len(values, 4)
	assert.Nil(values[3], "we shouldn't emit a literal 'null' here")
	assert.NotEqual("null", values[3], "we shouldn't emit a literal 'null' here")

	assert.Nil(Default().Invoke(context.Background(), tx).Create(&obj1))

	var verify1 jsonTest
	assert.Nil(Default().Invoke(context.Background(), tx).Get(&verify1, obj1.ID))

	assert.Equal(obj1.ID, verify1.ID)
	assert.Equal(obj1.Name, verify1.Name)
	assert.Nil(verify1.Nullable)
	assert.Equal(obj1.NotNull.Label, verify1.NotNull.Label)

	any, err := Default().Invoke(context.Background(), tx).Query("select 1 from json_test where id = $1 and nullable is null", obj1.ID).Any()
	assert.Nil(err)
	assert.True(any, "we should have written a sql null, not a literal string 'null'")

	// set it to literal 'null' to test this is backward compatible
	assert.Nil(Default().Invoke(context.Background(), tx).Exec("update json_test set nullable = 'null' where id = $1", obj1.ID))

	var verify2 jsonTest
	assert.Nil(Default().Invoke(context.Background(), tx).Get(&verify2, obj1.ID))
	assert.Equal(obj1.ID, verify2.ID)
	assert.Equal(obj1.Name, verify2.Name)
	assert.Nil(verify2.Nullable, "even if we set it to literal 'null' it should come out golang nil")
	assert.Equal(obj1.NotNull.Label, verify2.NotNull.Label)
}

type uniqueObj struct {
	ID   int    `db:"id,pk"`
	Name string `db:"name"`
}

// TableName returns the mapped table name.
func (uo uniqueObj) TableName() string {
	return "unique_obj"
}

func TestInvocationCreateRepeatInTx(t *testing.T) {
	assert := assert.New(t)
	tx, err := Default().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	assert.Nil(Default().Invoke(context.Background(), tx).Exec("CREATE TABLE IF NOT EXISTS unique_obj (id int not null primary key, name varchar)"))
	assert.Nil(Default().Invoke(context.Background(), tx).Create(&uniqueObj{ID: 1, Name: "one"}))
	var verify uniqueObj
	assert.Nil(Default().Invoke(context.Background(), tx).Get(&verify, 1))
	assert.Equal("one", verify.Name)
	assert.NotNil(Default().Invoke(context.Background(), tx).Create(&uniqueObj{ID: 1, Name: "two"}))
}

func TestInvocationExecError(t *testing.T) {
	assert := assert.New(t)

	conn := MustNewFromEnv()
	conn.StatementCache().WithEnabled(false)
	assert.Nil(conn.Open())
	assert.NotNil(conn.Invoke(context.Background()).Exec("not a select"))
	conn.StatementCache().WithEnabled(true)
	assert.NotNil(conn.Invoke(context.Background()).Exec("not a select"))
	assert.NotNil(conn.Invoke(context.Background()).WithLabel("exec_error_test").Exec("not a select"))
}

type modelTableNameError struct {
	ID string `db:"id,pk"`
}

func (mtne modelTableNameError) TableName() string {
	return uuid.V4().String()
}

func TestInvocationGetError(t *testing.T) {
	assert := assert.New(t)

	var getError modelTableNameError
	conn := MustNewFromEnv()
	conn.StatementCache().WithEnabled(false)
	assert.Nil(conn.Open())
	assert.NotNil(conn.Invoke(context.Background()).Get(&getError, uuid.V4().String()))
	conn.StatementCache().WithEnabled(true)
	assert.NotNil(conn.Invoke(context.Background()).Get(&getError, uuid.V4().String()))
	assert.NotNil(conn.Invoke(context.Background()).WithLabel("get_error_test").Get(&getError, uuid.V4().String()))
}

func TestInvocationGetAllError(t *testing.T) {
	assert := assert.New(t)

	var mustError []modelTableNameError
	conn := MustNewFromEnv()
	conn.StatementCache().WithEnabled(false)
	assert.Nil(conn.Open())
	assert.NotNil(conn.Invoke(context.Background()).GetAll(&mustError))
	conn.StatementCache().WithEnabled(true)
	assert.NotNil(conn.Invoke(context.Background()).GetAll(&mustError))
	assert.NotNil(conn.Invoke(context.Background()).WithLabel("get_all_error_test").GetAll(&mustError))
}

func TestInvocationCreateError(t *testing.T) {
	assert := assert.New(t)

	var mustError modelTableNameError
	conn := MustNewFromEnv()
	conn.StatementCache().WithEnabled(false)
	assert.Nil(conn.Open())
	assert.NotNil(conn.Invoke(context.Background()).Create(&mustError))
	conn.StatementCache().WithEnabled(true)
	assert.NotNil(conn.Invoke(context.Background()).Create(&mustError))
	assert.NotNil(conn.Invoke(context.Background()).WithLabel("create_error_test").Create(&mustError))
}

func TestInvocationCreateIfNotExistsError(t *testing.T) {
	assert := assert.New(t)

	var mustError modelTableNameError
	conn := MustNewFromEnv()
	conn.StatementCache().WithEnabled(false)
	assert.Nil(conn.Open())
	assert.NotNil(conn.Invoke(context.Background()).CreateIfNotExists(&mustError))
	conn.StatementCache().WithEnabled(true)
	assert.NotNil(conn.Invoke(context.Background()).CreateIfNotExists(&mustError))
	assert.NotNil(conn.Invoke(context.Background()).WithLabel("cne_error_test").CreateIfNotExists(&mustError))
}

func TestInvocationUpdateError(t *testing.T) {
	assert := assert.New(t)

	var mustError modelTableNameError
	conn := MustNewFromEnv()
	conn.StatementCache().WithEnabled(false)
	assert.Nil(conn.Open())
	assert.NotNil(conn.Invoke(context.Background()).Update(&mustError))
	conn.StatementCache().WithEnabled(true)
	assert.NotNil(conn.Invoke(context.Background()).Update(&mustError))
	assert.NotNil(conn.Invoke(context.Background()).WithLabel("update_error_test").Update(&mustError))
}

func TestInvocationUpsertError(t *testing.T) {
	assert := assert.New(t)

	var mustError modelTableNameError
	conn := MustNewFromEnv()
	conn.StatementCache().WithEnabled(false)
	assert.Nil(conn.Open())
	assert.NotNil(conn.Invoke(context.Background()).Upsert(&mustError))
	conn.StatementCache().WithEnabled(true)
	assert.NotNil(conn.Invoke(context.Background()).Upsert(&mustError))
	assert.NotNil(conn.Invoke(context.Background()).WithLabel("upsert_error_test").Upsert(&mustError))
}

func boolErr(_ bool, err error) error {
	return err
}

func TestInvocationExistsError(t *testing.T) {
	assert := assert.New(t)

	var mustError modelTableNameError
	conn := MustNewFromEnv()
	conn.StatementCache().WithEnabled(false)
	assert.Nil(conn.Open())
	assert.NotNil(boolErr(conn.Invoke(context.Background()).Exists(mustError)))
	conn.StatementCache().WithEnabled(true)
	assert.NotNil(boolErr(conn.Invoke(context.Background()).Exists(mustError)))
	assert.NotNil(boolErr(conn.Invoke(context.Background()).WithLabel("exists_error_test").Exists(mustError)))
}

func TestInvocationCreateManyError(t *testing.T) {
	assert := assert.New(t)

	mustError := []modelTableNameError{
		{uuid.V4().String()},
		{uuid.V4().String()},
	}
	conn := MustNewFromEnv()
	conn.StatementCache().WithEnabled(false)
	assert.Nil(conn.Open())
	assert.NotNil(conn.Invoke(context.Background()).CreateMany(mustError))
	conn.StatementCache().WithEnabled(true)
	assert.NotNil(conn.Invoke(context.Background()).CreateMany(mustError))
	assert.NotNil(conn.Invoke(context.Background()).WithLabel("cm_error_test").CreateMany(mustError))
}

func TestInvocationDeleteError(t *testing.T) {
	assert := assert.New(t)

	var mustError modelTableNameError
	conn := MustNewFromEnv()
	conn.StatementCache().WithEnabled(false)
	assert.Nil(conn.Open())
	assert.NotNil(conn.Invoke(context.Background()).Delete(&mustError))
	conn.StatementCache().WithEnabled(true)
	assert.NotNil(conn.Invoke(context.Background()).Delete(&mustError))
	assert.NotNil(conn.Invoke(context.Background()).WithLabel("delete_error_test").Delete(&mustError))
}

func TestTruncateError(t *testing.T) {
	assert := assert.New(t)

	var mustError modelTableNameError
	conn := MustNewFromEnv()
	conn.StatementCache().WithEnabled(false)
	assert.Nil(conn.Open())
	assert.NotNil(conn.Invoke(context.Background()).Truncate(&mustError))
	conn.StatementCache().WithEnabled(true)
	assert.NotNil(conn.Invoke(context.Background()).Truncate(&mustError))
	assert.NotNil(conn.Invoke(context.Background()).WithLabel("truncate_error_test").Truncate(&mustError))
}

type uuidTest struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

func (ut uuidTest) TableName() string {
	return "uuid_test"
}

func TestInvocationUUIDs(t *testing.T) {
	assert := assert.New(t)
	tx, err := Default().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	assert.Nil(Default().Invoke(context.Background(), tx).Exec("CREATE TABLE IF NOT EXISTS uuid_test (id uuid not null, name varchar(255) not null)"))

	assert.Nil(Default().Invoke(context.Background(), tx).Create(&uuidTest{ID: uuid.V4(), Name: "foo"}))
	assert.Nil(Default().Invoke(context.Background(), tx).Create(&uuidTest{ID: uuid.V4(), Name: "foo2"}))

	var objs []uuidTest
	assert.Nil(Default().Invoke(context.Background(), tx).GetAll(&objs))

	assert.Len(objs, 2)
}

type EmbeddedTestMeta struct {
	ID           uuid.UUID `db:"id,pk"`
	TimestampUTC time.Time `db:"timestamp_utc"`
}

type embeddedTest struct {
	EmbeddedTestMeta `db:",inline"`
	Name             string `db:"name"`
}

func (et embeddedTest) TableName() string {
	return "embedded_test"
}

func TestInlineMeta(t *testing.T) {
	assert := assert.New(t)
	tx, err := Default().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	test := &embeddedTest{EmbeddedTestMeta: EmbeddedTestMeta{ID: uuid.V4(), TimestampUTC: time.Now().UTC()}, Name: "foo"}
	cols := Columns(test)
	assert.NotEmpty(cols.PrimaryKeys().Columns())
	assert.Equal("id", cols.Columns()[0].ColumnName)
	assert.Equal("timestamp_utc", cols.Columns()[1].ColumnName)
	assert.Equal("name", cols.Columns()[2].ColumnName)

	values := cols.NotReadOnly().NotAutos().ColumnValues(test)
	assert.Len(values, 3)
	assert.Equal(test.ID, values[0])
	assert.False(values[1].(time.Time).IsZero())
	assert.Equal("foo", values[2])

	id0 := uuid.V4()
	id1 := uuid.V4()
	assert.Nil(Default().Invoke(context.Background(), tx).Exec("CREATE TABLE IF NOT EXISTS embedded_test (id uuid not null primary key, timestamp_utc timestamp not null, name varchar(255) not null)"))
	assert.Nil(Default().Invoke(context.Background(), tx).Create(&embeddedTest{EmbeddedTestMeta: EmbeddedTestMeta{ID: id0, TimestampUTC: time.Now().UTC()}, Name: "foo"}))
	assert.Nil(Default().Invoke(context.Background(), tx).Create(&embeddedTest{EmbeddedTestMeta: EmbeddedTestMeta{ID: id1, TimestampUTC: time.Now().UTC()}, Name: "foo2"}))

	var objs []embeddedTest
	assert.Nil(Default().Invoke(context.Background(), tx).GetAll(&objs))

	assert.Len(objs, 2)
	assert.Any(objs, func(v interface{}) bool {
		return v.(embeddedTest).ID.Equals(id0)
	})
	assert.Any(objs, func(v interface{}) bool {
		return v.(embeddedTest).ID.Equals(id1)
	})
	assert.Any(objs, func(v interface{}) bool {
		return v.(embeddedTest).Name == "foo"
	})
	assert.Any(objs, func(v interface{}) bool {
		return v.(embeddedTest).Name == "foo2"
	})
	assert.All(objs, func(v interface{}) bool {
		return !v.(embeddedTest).TimestampUTC.IsZero()
	})
}
