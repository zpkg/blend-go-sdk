package db

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
)

func TestInvocationErr(t *testing.T) {
	assert := assert.New(t)

	inv := &Invocation{err: fmt.Errorf("this is only a test")}
	assert.NotNil(inv.Err())
}

func TestInvocationLabels(t *testing.T) {
	assert := assert.New(t)

	inv := &Invocation{}
	inv = inv.WithLabel("test")
	assert.NotEmpty(inv.Label())
}

func TestInvocationPrepare(t *testing.T) {
	assert := assert.New(t)

	inv := &Invocation{err: fmt.Errorf("test")}
	_, err := inv.Prepare("select 'ok!'")
	assert.NotNil(err)
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
	return Default().InTx(tx).Exec("create table json_test (id serial primary key, name varchar(255), not_null json, nullable json)")
}

func dropJSONTextTable(tx *sql.Tx) error {
	return Default().InTx(tx).Exec("drop table if exists json_test")
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
	assert.Nil(Default().InTx(tx).Create(&obj0))

	var verify0 jsonTest
	assert.Nil(Default().InTx(tx).Get(&verify0, obj0.ID))

	assert.Equal(obj0.ID, verify0.ID)
	assert.Equal(obj0.Name, verify0.Name)
	assert.Equal(obj0.Nullable, verify0.Nullable)
	assert.Equal(obj0.NotNull.Label, verify0.NotNull.Label)

	// try creating partially set object and reading it out
	obj1 := jsonTest{Name: uuid.V4().String(), NotNull: jsonTestChild{Label: uuid.V4().String()}} //note `Nullable` isn't set
	assert.Nil(Default().InTx(tx).Create(&obj1))

	var verify1 jsonTest
	assert.Nil(Default().InTx(tx).Get(&verify1, obj1.ID))

	assert.Equal(obj1.ID, verify1.ID)
	assert.Equal(obj1.Name, verify1.Name)
	assert.Nil(verify1.Nullable)
	assert.Equal(obj1.NotNull.Label, verify1.NotNull.Label)
}
