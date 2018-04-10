package assert

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestEmpty(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	a := Empty().WithOutput(buf)
	if a.Output() == nil {
		t.Error("The empty assertion helper should have an output set")
		t.Fail()
	}
	if a.NonFatal().True(false, "this should fail") {
		t.Error("NonFatal true(false) didn't fail.")
	}
	if !a.NonFatal().True(true, "this should pass") {
		t.Error("NonFatal true(true) didn't pass.")
	}

	if len(buf.String()) == 0 {
		t.Error("We should have produced output.")
	}
}

func TestIsZero(t *testing.T) {
	zeroShort := int16(0)
	if !isZero(zeroShort) {
		t.Error("isZero failed")
	}

	notZeroShort := int16(3)
	if isZero(notZeroShort) {
		t.Error("isZero failed")
	}

	zero := 0
	if !isZero(zero) {
		t.Error("isZero failed")
	}
	notZero := 3
	if isZero(notZero) {
		t.Error("isZero failed")
	}

	zeroFloat64 := 0.0
	if !isZero(zeroFloat64) {
		t.Error("isZero failed")
	}
	notZeroFloat64 := 3.14
	if isZero(notZeroFloat64) {
		t.Error("isZero failed")
	}

	zeroFloat32 := float32(0.0)
	if !isZero(zeroFloat32) {
		t.Error("isZero failed")
	}
	notZeroFloat32 := float32(3.14)
	if isZero(notZeroFloat32) {
		t.Error("isZero failed")
	}
}

func TestGetLength(t *testing.T) {
	emptyString := ""
	l := getLength(emptyString)
	if l != 0 {
		t.Errorf("getLength incorrect.")
	}

	notEmptyString := "foo"
	l = getLength(notEmptyString)
	if l != 3 {
		t.Errorf("getLength incorrect.")
	}

	emptyArray := []int{}
	l = getLength(emptyArray)
	if l != 0 {
		t.Errorf("getLength incorrect.")
	}

	notEmptyArray := []int{1, 2, 3}
	l = getLength(notEmptyArray)
	if l != 3 {
		t.Errorf("getLength incorrect.")
	}

	emptyMap := map[string]int{}
	l = getLength(emptyMap)
	if l != 0 {
		t.Errorf("getLength incorrect.")
	}

	notEmptyMap := map[string]int{"foo": 1, "bar": 2, "baz": 3}
	l = getLength(notEmptyMap)
	if l != 3 {
		t.Errorf("getLength incorrect.")
	}
}

type myNestedStruct struct {
	ID   int
	Name string
}

type myTestStruct struct {
	ID          int
	Name        string
	SingleValue float32
	DoubleValue float64
	Timestamp   time.Time
	Struct      myNestedStruct

	IDPtr     *int
	NamePptr  *string
	StructPtr *myNestedStruct

	Slice    []myNestedStruct
	SlicePtr *[]myNestedStruct
}

func createTestStruct() myTestStruct {

	testInt := 1
	testName := "test struct"

	nestedA := myNestedStruct{1, "A"}
	nestedB := myNestedStruct{1, "B"}
	nestedC := myNestedStruct{1, "C"}

	testStruct := myTestStruct{
		ID:          testInt,
		Name:        testName,
		SingleValue: float32(3.14),
		DoubleValue: 6.28,
		Timestamp:   time.Now(),
		Struct:      nestedA,
		IDPtr:       &testInt,
		NamePptr:    &testName,
		StructPtr:   &nestedB,
		Slice:       []myNestedStruct{nestedA, nestedB, nestedC},
	}

	testStruct.SlicePtr = &testStruct.Slice
	return testStruct

}

func TestStructsAreEqual(t *testing.T) {
	testStructA := createTestStruct()
	testStructB := createTestStruct()
	testStructB.Name = "not test struct"

	if didFail, _ := shouldBeEqual(testStructA, testStructA); didFail {
		t.Error("shouldBeEqual Failed.")
		t.FailNow()
	}

	if didFail, _ := shouldBeEqual(testStructA, testStructB); !didFail {
		t.Error("shouldBeEqual Failed.")
		t.FailNow()
	}
}

func TestShouldBeEqual(t *testing.T) {
	byteA := byte('a')
	byteB := byte('b')

	if didFail, _ := shouldBeEqual(byteA, byteA); didFail {
		t.Error("shouldBeEqual Failed.")
		t.FailNow()
	}
	if didFail, _ := shouldBeEqual(byteA, byteB); !didFail {
		t.Error("shouldBeEqual Failed.")
		t.FailNow()
	}

	stringA := "test"
	stringB := "not test"

	if didFail, _ := shouldBeEqual(stringA, stringA); didFail {
		t.Error("shouldBeEqual Equal Failed.")
		t.FailNow()
	}
	if didFail, _ := shouldBeEqual(stringA, stringB); !didFail {
		t.Error("shouldBeEqual Failed.")
		t.FailNow()
	}

	intA := 1
	intB := 2

	if didFail, _ := shouldBeEqual(intA, intA); didFail {
		t.Error("shouldBeEqual Equal Failed.")
		t.FailNow()
	}
	if didFail, _ := shouldBeEqual(intA, intB); !didFail {
		t.Error("shouldBeEqual Failed.")
		t.FailNow()
	}

	float32A := float32(3.14)
	float32B := float32(6.28)

	if didFail, _ := shouldBeEqual(float32A, float32A); didFail {
		t.Error("shouldBeEqual Equal Failed.")
		t.FailNow()
	}
	if didFail, _ := shouldBeEqual(float32A, float32B); !didFail {
		t.Error("shouldBeEqual Failed.")
		t.FailNow()
	}

	floatA := 3.14
	floatB := 6.28

	if didFail, _ := shouldBeEqual(floatA, floatA); didFail {
		t.Error("shouldBeEqual Equal Failed.")
		t.FailNow()
	}
	if didFail, _ := shouldBeEqual(floatA, floatB); !didFail {
		t.Error("shouldBeEqual Failed.")
		t.FailNow()
	}
}

func makesThings(shouldReturnNil bool) *myTestStruct {
	if !shouldReturnNil {
		return &myTestStruct{}
	}
	return nil
}

func TestShouldBeNil(t *testing.T) {
	assertsToNil := makesThings(true)
	assertsToNotNil := makesThings(false)

	didFail, didFailErrMsg := shouldBeNil(assertsToNil)
	if didFail {
		t.Error(didFailErrMsg)
		t.FailNow()
	}

	didFail, didFailErrMsg = shouldBeNil(assertsToNotNil)
	if !didFail {
		t.Error("shouldBeNil returned did_fail as `true` for a not nil object")
		t.FailNow()
	}
}

func TestShouldNotBeNil(t *testing.T) {
	assertsToNil := makesThings(true)
	assertsToNotNil := makesThings(false)

	didFail, didFailErrMsg := shouldNotBeNil(assertsToNotNil)
	if didFail {
		t.Error(didFailErrMsg)
		t.FailNow()
	}

	didFail, didFailErrMsg = shouldNotBeNil(assertsToNil)
	if !didFail {
		t.Error("shouldNotBeNil returned did_fail as `true` for a not nil object")
		t.FailNow()
	}
}

func TestShouldContain(t *testing.T) {
	shouldNotHaveFailed, _ := shouldContain("is a", "this is a test")
	if shouldNotHaveFailed {
		t.Errorf("shouldConatain failed.")
		t.FailNow()
	}

	shouldHaveFailed, _ := shouldContain("beer", "this is a test")
	if !shouldHaveFailed {
		t.Errorf("shouldConatain failed.")
		t.FailNow()
	}
}

type anyTestObj struct {
	ID   int
	Name string
}

func TestAny(t *testing.T) {
	testObjs := []anyTestObj{{1, "Test"}, {2, "Test2"}, {3, "Foo"}}

	didFail, _ := shouldAny(testObjs, func(obj interface{}) bool {
		if typed, didType := obj.(anyTestObj); didType {
			return strings.HasPrefix(typed.Name, "Foo")
		}
		return false
	})
	if didFail {
		t.Errorf("shouldAny failed.")
		t.FailNow()
	}

	didFail, _ = shouldAny(testObjs, func(obj interface{}) bool {
		if typed, didType := obj.(anyTestObj); didType {
			return strings.HasPrefix(typed.Name, "Bar")
		}
		return false
	})
	if !didFail {
		t.Errorf("shouldAny should have failed.")
		t.FailNow()
	}

	didFail, _ = shouldAny(anyTestObj{1, "test"}, func(obj interface{}) bool {
		return true
	})
	if !didFail {
		t.Errorf("shouldAny should have failed on non-slice target.")
		t.FailNow()
	}
}

func TestAll(t *testing.T) {
	testObjs := []anyTestObj{{1, "Test"}, {2, "Test2"}, {3, "Foo"}}

	didFail, _ := shouldAll(testObjs, func(obj interface{}) bool {
		if typed, didType := obj.(anyTestObj); didType {
			return typed.ID > 0
		}
		return false
	})
	if didFail {
		t.Errorf("shouldAll shouldnt have failed.")
		t.FailNow()
	}

	didFail, _ = shouldAll(testObjs, func(obj interface{}) bool {
		if typed, didType := obj.(anyTestObj); didType {
			return strings.HasPrefix(typed.Name, "Test")
		}
		return false
	})
	if !didFail {
		t.Errorf("shouldAll should have failed.")
		t.FailNow()
	}

	didFail, _ = shouldAll(anyTestObj{1, "test"}, func(obj interface{}) bool {
		return true
	})
	if !didFail {
		t.Errorf("shouldAll should have failed on non-slice target.")
		t.FailNow()
	}
}

func TestNone(t *testing.T) {
	testObjs := []anyTestObj{{1, "Test"}, {2, "Test2"}, {3, "Foo"}}

	didFail, _ := shouldNone(testObjs, func(obj interface{}) bool {
		if typed, didType := obj.(anyTestObj); didType {
			return typed.ID > 4
		}
		return false
	})
	if didFail {
		t.Errorf("shouldAll shouldnt have failed.")
		t.FailNow()
	}

	didFail, _ = shouldNone(testObjs, func(obj interface{}) bool {
		if typed, didType := obj.(anyTestObj); didType {
			return typed.ID > 0
		}
		return false
	})
	if !didFail {
		t.Errorf("shouldNone should have failed.")
		t.FailNow()
	}
}

func TestInTimeDelta(t *testing.T) {
	value1 := time.Date(2016, 1, 29, 9, 0, 0, 0, time.UTC)
	value2 := time.Date(2016, 1, 29, 9, 0, 0, 1, time.UTC)
	value3 := time.Date(2016, 1, 29, 8, 0, 0, 0, time.UTC)
	value4 := time.Date(2015, 1, 29, 9, 0, 0, 0, time.UTC)

	didFail, _ := shouldBeInTimeDelta(value1, value2, 1*time.Minute)
	if didFail {
		t.Errorf("shouldBeInTimeDelta shouldnt have failed.")
		t.FailNow()
	}

	didFail, _ = shouldBeInTimeDelta(value1, value3, 1*time.Minute)
	if !didFail {
		t.Errorf("shouldBeInTimeDelta should have failed.")
		t.FailNow()
	}

	didFail, _ = shouldBeInTimeDelta(value1, value4, 1*time.Minute)
	if !didFail {
		t.Errorf("shouldBeInTimeDelta should have failed.")
		t.FailNow()
	}
}
