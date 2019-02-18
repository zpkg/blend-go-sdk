package secrets

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/reflectutil"
)

// Constants
const (
	StructTag = "secret"
)

// DecomposeJSON decomposes an object into json fields marked with the `secret` struct tag.
// Top level fields will get their own keys.
// Nested objects are serialized as json.
func DecomposeJSON(obj interface{}) (map[string]string, error) {
	output := make(map[string]string)

	objType := reflectutil.Type(obj)
	objValue := reflectutil.Value(obj)

	var field reflect.StructField
	var fieldValue reflect.Value
	var tagValue string
	var outputKey string
	var tagPieces []string
	for x := 0; x < objType.NumField(); x++ {
		field = objType.Field(x)
		if !reflectutil.IsExported(field.Name) {
			continue
		}

		fieldValue = objValue.FieldByName(field.Name)
		tagValue = field.Tag.Get(StructTag)
		if tagValue == "" || tagValue == "-" {
			continue
		}
		tagPieces = strings.Split(tagValue, ",")
		outputKey = tagPieces[0]
		contents, err := json.Marshal(fieldValue.Interface())
		if err != nil {
			return nil, exception.New(err)
		}
		output[outputKey] = string(contents)
	}

	return output, nil
}

// RestoreJSON restores an object from a given data bag as JSON.
func RestoreJSON(data map[string]string, obj interface{}) error {
	objType := reflectutil.Type(obj)
	objValue := reflectutil.Value(obj)

	fieldLookup := make(map[string]string)

	var field reflect.StructField
	var fieldValue reflect.Value
	var tagValue string
	var outputKey string
	var tagPieces []string
	for x := 0; x < objType.NumField(); x++ {
		field = objType.Field(x)
		if !reflectutil.IsExported(field.Name) {
			continue
		}

		fieldValue = objValue.FieldByName(field.Name)
		tagValue = field.Tag.Get(StructTag)
		if tagValue == "" || tagValue == "-" {
			continue
		}
		tagPieces = strings.Split(tagValue, ",")
		outputKey = tagPieces[0]
		fieldLookup[outputKey] = field.Name
	}

	var fieldName string
	var ok bool
	for key, value := range data {
		// figure out which field matches the key ...
		if fieldName, ok = fieldLookup[key]; !ok {
			continue
		}
		fieldValue = objValue.FieldByName(fieldName)
		if err := json.Unmarshal([]byte(value), fieldValue.Addr().Interface()); err != nil {
			return exception.New(err)
		}
	}
	return nil
}
