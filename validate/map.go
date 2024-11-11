/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

import (
	"reflect"

	"github.com/zpkg/blend-go-sdk/ex"
	"github.com/zpkg/blend-go-sdk/reflectutil"
)

// Errors
const (
	ErrInstanceNotMap ex.Class = "validated reference is not a map"
	ErrMapKeys        ex.Class = "map should have keys"
)

// Map returns validators for a map type reference.
func Map(instance interface{}) MapValidators {
	return MapValidators{instance}
}

// MapValidators is a set of validators for maps.
type MapValidators struct {
	Value interface{}
}

// Keys validates a map contains a given set of keys.
func (mv MapValidators) Keys(keys ...interface{}) Validator {
	return func() error {
		value := reflectutil.Value(mv.Value)
		if value.Kind() != reflect.Map {
			return ErrInstanceNotMap
		}

		for _, key := range keys {
			mapValue := value.MapIndex(reflect.ValueOf(key))
			if !mapValue.IsValid() {
				return Errorf(ErrMapKeys, mv.Value, "missing key: %v", key)
			}
		}
		return nil
	}
}
