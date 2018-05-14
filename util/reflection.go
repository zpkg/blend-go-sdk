package util

import (
	"encoding/base64"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/blend/go-sdk/exception"
)

var (
	// Reflection is a namespace for reflection utilities.
	Reflection = reflectionUtil{}
)

// Patchable describes an object that can be patched with raw values.
type Patchable interface {
	Patch(values map[string]interface{}) error
}

type reflectionUtil struct{}

// FollowValuePointer derefs a reflectValue until it isn't a pointer, but will preseve it's nilness.
func (ru reflectionUtil) FollowValuePointer(v reflect.Value) interface{} {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil
	}

	val := v
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val.Interface()
}

// FollowType derefs a type until it isn't a pointer or an interface.
func (ru reflectionUtil) FollowType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}
	return t
}

// FollowValue derefs a value until it isn't a pointer or an interface.
func (ru reflectionUtil) FollowValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}

// ReflectValue returns the integral reflect.Value for an object.
func (ru reflectionUtil) ReflectValue(obj interface{}) reflect.Value {
	v := reflect.ValueOf(obj)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}

// ReflectType returns the integral type for an object.
func (ru reflectionUtil) ReflectType(obj interface{}) reflect.Type {
	t := reflect.TypeOf(obj)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t
}

// MakeNew returns a new instance of a reflect.Type.
func (ru reflectionUtil) MakeNew(t reflect.Type) interface{} {
	return reflect.New(t).Interface()
}

// MakeSliceOfType returns a new slice of a given reflect.Type.
func (ru reflectionUtil) MakeSliceOfType(t reflect.Type) interface{} {
	return reflect.New(reflect.SliceOf(t)).Interface()
}

// TypeName returns the string type name for an object's integral type.
func (ru reflectionUtil) TypeName(obj interface{}) string {
	return ru.ReflectType(obj).Name()
}

// GetValueByName returns a value for a given struct field by name.
func (ru reflectionUtil) GetValueByName(target interface{}, fieldName string) interface{} {
	targetValue := ru.ReflectValue(target)
	field := targetValue.FieldByName(fieldName)
	return field.Interface()
}

// GetFieldByNameOrJSONTag returns a value for a given struct field by name or by json tag name.
func (ru reflectionUtil) GetFieldByNameOrJSONTag(targetValue reflect.Type, fieldName string) *reflect.StructField {
	for index := 0; index < targetValue.NumField(); index++ {
		field := targetValue.Field(index)

		if field.Name == fieldName {
			return &field
		}
		tag := field.Tag
		jsonTag := tag.Get("json")
		if String.CaseInsensitiveEquals(jsonTag, fieldName) {
			return &field
		}
	}

	return nil
}

func (ru reflectionUtil) SetValueByName(target interface{}, fieldName string, fieldValue interface{}) error {
	targetValue := ru.ReflectValue(target)
	targetType := ru.ReflectType(target)
	return ru.SetValueByNameFromType(target, targetType, targetValue, fieldName, fieldValue)
}

// SetValueByName sets a value on an object by its field name.
func (ru reflectionUtil) SetValueByNameFromType(obj interface{}, targetType reflect.Type, targetValue reflect.Value, fieldName string, fieldValue interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = exception.Newf("panic setting value by name").WithMessagef("field: %s panic: %v", fieldName, r)
		}
	}()

	relevantField := ru.GetFieldByNameOrJSONTag(targetType, fieldName)
	if relevantField == nil {
		err = exception.New("unknown field").WithMessagef("%s `%s`", targetType.Name(), fieldName)
		return
	}

	field := targetValue.FieldByName(relevantField.Name)
	if !field.CanSet() {
		err = exception.New("cannot set field").WithMessagef("%s `%s`", targetType.Name(), fieldName)
		return
	}

	fieldType := field.Type()
	value := ru.ReflectValue(fieldValue)
	valueType := value.Type()
	if !value.IsValid() {
		err = exception.New("invalid value").WithMessagef("%s `%s`", targetType.Name(), fieldName)
		return
	}

	assigned, assignErr := ru.tryAssignment(fieldType, valueType, field, value)
	if assignErr != nil {
		err = assignErr
		return
	}
	if !assigned {
		err = exception.New("cannot set field").WithMessagef("%s `%s`", targetType.Name(), fieldName)
		return
	}
	return
}

func (ru reflectionUtil) tryAssignment(fieldType, valueType reflect.Type, field, value reflect.Value) (assigned bool, err error) {
	if valueType.AssignableTo(fieldType) {
		field.Set(value)
		assigned = true
		return
	}

	if valueType.ConvertibleTo(fieldType) {
		convertedValue := value.Convert(fieldType)
		if convertedValue.Type().AssignableTo(fieldType) {
			field.Set(convertedValue)
			assigned = true
			return
		}
	}

	if fieldType.Kind() == reflect.Ptr {
		if valueType.AssignableTo(fieldType.Elem()) {
			elem := reflect.New(fieldType.Elem())
			elem.Elem().Set(value)
			field.Set(elem)
			assigned = true
			return
		} else if valueType.ConvertibleTo(fieldType.Elem()) {
			elem := reflect.New(fieldType.Elem())
			elem.Elem().Set(value.Convert(fieldType.Elem()))
			field.Set(elem)
			assigned = true
			return
		}
	}

	return
}

// Patch updates an object based on a map of field names to values.
func (ru reflectionUtil) Patch(obj interface{}, patchValues map[string]interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = exception.Newf("%v", r)
		}
	}()

	if patchable, isPatchable := obj.(Patchable); isPatchable {
		return patchable.Patch(patchValues)
	}

	targetValue := ru.ReflectValue(obj)
	targetType := targetValue.Type()

	for key, value := range patchValues {
		err = ru.SetValueByNameFromType(obj, targetType, targetValue, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// Decompose is a *very* inefficient way to turn an object into a map string => interface.
func (ru reflectionUtil) Decompose(object interface{}) map[string]interface{} {
	var output map[string]interface{}
	JSON.Deserialize(&output, JSON.Serialize(object))
	return output
}

// checks if a value is a zero value or its types default value
func (ru reflectionUtil) IsZero(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && ru.IsZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && ru.IsZero(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}

// IsExported returns if a field is exported given its name and capitalization.
func (ru reflectionUtil) IsExported(fieldName string) bool {
	return fieldName != "" && strings.ToUpper(fieldName)[0] == fieldName[0]
}

// CoalesceFields merges non-zero fields into destination fields marked with the `coalesce:...` struct field tag.
func (ru reflectionUtil) CoalesceFields(object interface{}) {
	objectValue := ru.ReflectValue(object)
	objectType := ru.ReflectType(object)
	if objectType.Kind() == reflect.Struct {
		numberOfFields := objectValue.NumField()
		for index := 0; index < numberOfFields; index++ {
			field := objectType.Field(index)
			fieldValue := objectValue.Field(index)
			// only alter the field if it is exported (uppercase variable name) and is not already a non-zero value
			if ru.IsExported(field.Name) && ru.IsZero(fieldValue) {
				alternateFieldNames := strings.Split(field.Tag.Get("coalesce"), ",")

				// find the first non-zero value in the list of backup values
				for j := 0; j < len(alternateFieldNames); j++ {
					alternateFieldName := alternateFieldNames[j]
					alternateValue := objectValue.FieldByName(alternateFieldName)
					// will panic if trying to set a non-exported value or a zero value, so ignore those
					if ru.IsExported(alternateFieldName) && !ru.IsZero(alternateValue) {
						fieldValue.Set(alternateValue)
						break
					}
				}
			}
			// recurse, in case nested values of this field need to be set as well
			if ru.IsExported(field.Name) && !ru.IsZero(fieldValue) {
				ru.CoalesceFields(fieldValue.Addr().Interface())
			}
		}
	} else if objectType.Kind() == reflect.Array || objectType.Kind() == reflect.Slice {
		arrayLength := objectValue.Len()
		for i := 0; i < arrayLength; i++ {
			ru.CoalesceFields(objectValue.Index(i).Addr().Interface())
		}
	}
}

// PatchStrings options.
const (
	// FieldTagEnv is the struct tag for what environment variable to use to populate a field.
	FieldTagEnv = "env"
	// FieldFlagCSV is a field tag flag (say that 10 times fast).
	FieldFlagCSV = "csv"
	// FieldFlagBase64 is a field tag flag (say that 10 times fast).
	FieldFlagBase64 = "base64"
	// FieldFlagBytes is a field tag flag (say that 10 times fast).
	FieldFlagBytes = "bytes"
)

// MapStringsUnmarshaler is a type that handles unmarshalling a map of strings into itself.
type MapStringsUnmarshaler interface {
	UnmarshalMapStrings(data map[string]string) error
}

// PatchStrings sets an object from a set of strings mapping field names to string values (to be parsed).
func (ru reflectionUtil) PatchStrings(tagName string, data map[string]string, obj interface{}) error {
	// check if the type implements marshaler.
	if typed, isTyped := obj.(MapStringsUnmarshaler); isTyped {
		return typed.UnmarshalMapStrings(data)
	}

	objMeta := ru.ReflectType(obj)
	objValue := ru.ReflectValue(obj)

	typeDuration := reflect.TypeOf(time.Duration(time.Nanosecond))

	var field reflect.StructField
	var fieldType reflect.Type
	var fieldValue reflect.Value
	var tag string
	var err error
	var pieces []string
	var dataField string
	var dataValue string
	var dataFieldValue interface{}
	var hasDataValue bool

	var isCSV bool
	var isBytes bool
	var isBase64 bool

	for x := 0; x < objMeta.NumField(); x++ {
		isCSV = false
		isBytes = false
		isBase64 = false

		field = objMeta.Field(x)
		fieldValue = objValue.FieldByName(field.Name)

		// Treat structs as nested values.
		if field.Type.Kind() == reflect.Struct {
			if err = ru.PatchStrings(tagName, data, objValue.Field(x).Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		tag = field.Tag.Get(tagName)
		if len(tag) > 0 {
			pieces = strings.Split(tag, ",")
			dataField = pieces[0]
			if len(pieces) > 1 {
				for y := 1; y < len(pieces); y++ {
					if pieces[y] == FieldFlagCSV {
						isCSV = true
					} else if pieces[y] == FieldFlagBase64 {
						isBase64 = true
					} else if pieces[y] == FieldFlagBytes {
						isBytes = true
					}
				}
			}

			dataValue, hasDataValue = data[dataField]
			if !hasDataValue {
				continue
			}

			if isCSV {
				dataFieldValue = strings.Split(dataValue, ",")
			} else if isBase64 {
				dataFieldValue, err = base64.StdEncoding.DecodeString(dataValue)
				if err != nil {
					return err
				}
			} else if isBytes {
				dataFieldValue = []byte(dataValue)
			} else {
				// figure out the rootmost type (i.e. deref ****ptr etc.)
				fieldType = ru.FollowType(field.Type)
				switch fieldType {
				case typeDuration:
					dataFieldValue, err = time.ParseDuration(dataValue)
					if err != nil {
						return exception.Wrap(err)
					}
				default:
					switch fieldType.Kind() {
					case reflect.Bool:
						if hasDataValue {
							dataFieldValue = Parse.Bool(dataValue)
						} else {
							continue
						}
					case reflect.Float32:
						dataFieldValue, err = strconv.ParseFloat(dataValue, 32)
						if err != nil {
							return exception.Wrap(err)
						}
					case reflect.Float64:
						dataFieldValue, err = strconv.ParseFloat(dataValue, 64)
						if err != nil {
							return exception.Wrap(err)
						}
					case reflect.Int8:
						dataFieldValue, err = strconv.ParseInt(dataValue, 10, 8)
						if err != nil {
							return exception.Wrap(err)
						}
					case reflect.Int16:
						dataFieldValue, err = strconv.ParseInt(dataValue, 10, 16)
						if err != nil {
							return exception.Wrap(err)
						}
					case reflect.Int32:
						dataFieldValue, err = strconv.ParseInt(dataValue, 10, 32)
						if err != nil {
							return exception.Wrap(err)
						}
					case reflect.Int:
						dataFieldValue, err = strconv.ParseInt(dataValue, 10, 64)
						if err != nil {
							return exception.Wrap(err)
						}
					case reflect.Int64:
						dataFieldValue, err = strconv.ParseInt(dataValue, 10, 64)
						if err != nil {
							return exception.Wrap(err)
						}
					case reflect.Uint8:
						dataFieldValue, err = strconv.ParseUint(dataValue, 10, 8)
						if err != nil {
							return exception.Wrap(err)
						}
					case reflect.Uint16:
						dataFieldValue, err = strconv.ParseUint(dataValue, 10, 8)
						if err != nil {
							return exception.Wrap(err)
						}
					case reflect.Uint32:
						dataFieldValue, err = strconv.ParseUint(dataValue, 10, 32)
						if err != nil {
							return exception.Wrap(err)
						}
					case reflect.Uint64:
						dataFieldValue, err = strconv.ParseUint(dataValue, 10, 64)
						if err != nil {
							return exception.Wrap(err)
						}
					case reflect.Uint, reflect.Uintptr:
						dataFieldValue, err = strconv.ParseUint(dataValue, 10, 64)
						if err != nil {
							return exception.Wrap(err)
						}
					case reflect.String:
						dataFieldValue = dataValue
					default:
						return exception.New("map strings into; unhandled assignment").WithMessagef("type %s", fieldType.String())
					}
				}
			}

			value := ru.ReflectValue(dataFieldValue)
			valueType := value.Type()
			if !value.IsValid() {
				return exception.New("invalid value").WithMessagef("%s `%s`", objMeta.Name(), field.Name)
			}

			assigned, err := ru.tryAssignment(fieldType, valueType, fieldValue, value)
			if err != nil {
				return err
			}
			if !assigned {
				return exception.New("cannot set field").WithMessagef("%s `%s`", objMeta.Name(), field.Name)
			}
		}
	}
	return nil
}
