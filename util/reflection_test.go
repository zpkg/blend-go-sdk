package util

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

type subType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type testType struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	NotTagged string
	Tagged    string    `json:"is_tagged"`
	SubTypes  []subType `json:"children"`
}

func TestDecompose(t *testing.T) {
	assert := assert.New(t)

	myObj := testType{}
	myObj.ID = 123
	myObj.Name = "Test Object"
	myObj.NotTagged = "Not Tagged"
	myObj.Tagged = "Is Tagged"
	myObj.SubTypes = append([]subType{}, subType{1, "One"})
	myObj.SubTypes = append(myObj.SubTypes, subType{2, "Two"})
	myObj.SubTypes = append(myObj.SubTypes, subType{3, "Three"})
	myObj.SubTypes = append(myObj.SubTypes, subType{4, "Four"})

	decomposed := Reflection.Decompose(myObj)

	_, hasKey := decomposed["id"]
	assert.True(hasKey)

	_, hasKey = decomposed["name"]
	assert.True(hasKey)

	_, hasKey = decomposed["NotTagged"]
	assert.True(hasKey)

	_, hasKey = decomposed["is_tagged"]
	assert.True(hasKey)

	_, hasKey = decomposed["children"]
	assert.True(hasKey)
}

type TestType2 struct {
	SomeVal    string `coalesce:"SomeVal2"`
	SomeVal2   string
	OtherVal   string `coalesce:"OtherVal2,OtherVal3"`
	OtherVal2  string
	OtherVal3  string
	StructVal  subType `coalesce:"StructVal2"`
	StructVal2 subType
}

func TestCoalesceFieldsNoChange(t *testing.T) {
	assert := assert.New(t)

	testVal := TestType2{
		SomeVal:    "Foo1",
		SomeVal2:   "Foo2",
		OtherVal:   "Foo3",
		OtherVal2:  "Foo4",
		OtherVal3:  "Foo5",
		StructVal:  subType{1, "Name"},
		StructVal2: subType{2, "Name2"},
	}

	Reflection.CoalesceFields(&testVal)
	assert.Equal("Foo1", testVal.SomeVal)
	assert.Equal("Foo3", testVal.OtherVal)
	assert.Equal(subType{1, "Name"}, testVal.StructVal)
	// omit values to check that values coalesce
}

func TestCoalesceFieldsSimple(t *testing.T) {
	assert := assert.New(t)
	testVal2 := TestType2{
		SomeVal2:   "Foo2",
		OtherVal2:  "Foo4",
		OtherVal3:  "Foo5",
		StructVal2: subType{2, "Name2"},
	}
	Reflection.CoalesceFields(&testVal2)
	assert.Equal("Foo2", testVal2.SomeVal)
	assert.Equal("Foo4", testVal2.OtherVal)
	assert.Equal(subType{2, "Name2"}, testVal2.StructVal)
}

func TestCoalesceFieldsStructAndMultiples(t *testing.T) {
	assert := assert.New(t)
	// omit values to check that values coalesce
	testVal3 := TestType2{
		SomeVal:   "Foo1",
		OtherVal3: "Foo5",
		StructVal: subType{1, "Name"},
	}
	Reflection.CoalesceFields(&testVal3)
	assert.Equal("Foo5", testVal3.OtherVal)
}

type testType3 struct {
	Sub  subType2 `coalesce:"Sub2"`
	Sub2 subType2
}

type subType2 struct {
	Val1 string `coalesce:"Val2"`
	Val2 string
}

func TestCoalesceFieldsNested(t *testing.T) {
	assert := assert.New(t)
	t1 := testType3{
		Sub: subType2{"", "foo"},
	}

	Reflection.CoalesceFields(&t1)
	assert.Equal("foo", t1.Sub.Val1)

	t2 := testType3{
		Sub2: subType2{"", "foo2"},
	}

	Reflection.CoalesceFields(&t2)
	assert.Equal("foo2", t2.Sub.Val1)
}

type testType4 struct {
	Subs []subType2
}

func TestCoalesceFieldsArray(t *testing.T) {
	assert := assert.New(t)
	t1 := testType4{[]subType2{{"", "foo"}, subType2{"foo2", ""}}}
	Reflection.CoalesceFields(&t1)
	assert.Equal("foo", t1.Subs[0].Val1)
	assert.Equal("foo2", t1.Subs[1].Val1)
}

func TestPatchObject(t *testing.T) {
	assert := assert.New(t)

	myObj := testType{}
	myObj.ID = 123
	myObj.Name = "Test Object"
	myObj.NotTagged = "Not Tagged"
	myObj.Tagged = "Is Tagged"
	myObj.SubTypes = append([]subType{}, subType{1, "One"})
	myObj.SubTypes = append(myObj.SubTypes, subType{2, "Two"})
	myObj.SubTypes = append(myObj.SubTypes, subType{3, "Three"})
	myObj.SubTypes = append(myObj.SubTypes, subType{4, "Four"})

	patchData := make(map[string]interface{})
	patchData["is_tagged"] = "Is Not Tagged"

	err := Reflection.Patch(&myObj, patchData)
	assert.Nil(err)
	assert.Equal("Is Not Tagged", myObj.Tagged)
}

type TestObject struct {
	ID   int
	Name string
}

func testCachedObject(obj interface{}) func() interface{} {
	return func() interface{} {
		return obj
	}
}

func TestReflectTypeInterface(t *testing.T) {
	assert := assert.New(t)

	proto := testCachedObject(TestObject{ID: 1, Name: "Test"})

	assert.NotNil(proto())

	objType := Reflection.ReflectType(proto())
	assert.NotNil(objType)
}

func TestReflectValueInterface(t *testing.T) {
	assert := assert.New(t)

	proto := testCachedObject(&TestObject{ID: 1, Name: "Test"})

	assert.NotNil(proto())

	objValue := Reflection.ReflectValue(proto())
	assert.NotNil(objValue)
	assert.True(objValue.CanSet())
}

type mapStringsTest struct {
	Bool     bool          `secret:"bool"`
	Float32  float32       `secret:"float32"`
	Float64  float64       `secret:"float64"`
	Int8     int8          `secret:"int8"`
	Int16    int16         `secret:"int16"`
	Int32    int32         `secret:"int32"`
	Int64    int64         `secret:"int64"`
	Uint8    uint8         `secret:"uint8"`
	Uint16   uint16        `secret:"uint16"`
	Uint32   uint32        `secret:"uint32"`
	Uint64   uint32        `secret:"uint64"`
	String   string        `secret:"string"`
	Duration time.Duration `secret:"duration"`

	CSV    []string `secret:"csvField,csv"`
	Base64 []byte   `secret:"base64Field,bytes"`
	Bytes  []byte   `secret:"bytesField,bytes"`
}

func TestMapStringsInto(t *testing.T) {
	assert := assert.New(t)

	var mule mapStringsTest

	// ----
	// bool
	// ----

	boolValid := map[string]string{
		"bool": "true",
	}
	boolInvalid := map[string]string{
		"bool": "nottrue",
	}
	assert.Nil(Reflection.MapStringsInto("secret", boolValid, &mule))
	assert.Equal(true, mule.Bool)
	assert.Nil(Reflection.MapStringsInto("secret", boolInvalid, &mule))
	assert.Equal(false, mule.Bool)

	// -------
	// float32
	// -------

	float32Valid := map[string]string{
		"float32": "3.14",
	}
	float32Invalid := map[string]string{
		"float32": "random",
	}
	assert.Nil(Reflection.MapStringsInto("secret", float32Valid, &mule))
	assert.Equal(3.14, mule.Float32)
	assert.NotNil(Reflection.MapStringsInto("secret", float32Invalid, &mule))

	// -------
	// float64
	// -------

	float64Valid := map[string]string{
		"float64": "6.28",
	}
	float64Invalid := map[string]string{
		"float64": "random",
	}
	assert.Nil(Reflection.MapStringsInto("secret", float64Valid, &mule))
	assert.Equal(6.28, mule.Float64)
	assert.NotNil(Reflection.MapStringsInto("secret", float64Invalid, &mule))
}
