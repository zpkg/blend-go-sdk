/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

// ReturnAll returns all the failing validations as an
// array of validation errors.
func ReturnAll(validators ...Validator) error {
	var output []error
	var err error
	for _, validator := range validators {
		if err = validator(); err != nil {
			output = append(output, err)
		}
	}
	if len(output) > 0 {
		return ValidationErrors(output)
	}
	return nil
}
