package pg

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/uuid"
)

func TestConnectionCreate(t *testing.T) {
	assert := assert.New(t)
	tx, err := db.Default().Begin()
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
	err = db.Default().Invoke(OptTx(tx)).Insert(obj)
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
			innerErr := db.Default().CreateInTx(obj, nil)
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
	err = db.Default().UpsertInTx(obj, tx)
	assert.Nil(err)

	var verify upsertObj
	err = db.Default().GetInTx(&verify, tx, obj.UUID)
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
	tx, err := db.Default().Begin()
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
	err = db.Default().UpsertInTx(obj, tx)
	assert.Nil(err, fmt.Sprintf("%+v", err))
	assert.NotZero(obj.ID)

	var verify benchObj
	err = db.Default().GetInTx(&verify, tx, obj.ID)
	assert.Nil(err)
	assert.Equal(obj.Category, verify.Category)

	obj.Category = "test"

	err = db.Default().UpsertInTx(obj, tx)
	assert.Nil(err)
	assert.NotZero(obj.ID)

	err = db.Default().GetInTx(&verify, tx, obj.ID)
	assert.Nil(err)
	assert.Equal(obj.Category, verify.Category)
}

func TestConnectionCreateMany(t *testing.T) {
	assert := assert.New(t)
	tx, err := db.Default().Begin()
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
