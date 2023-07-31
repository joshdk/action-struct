// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package inputs provides functionality for unmarshalling values from the
// GitHub Actions environment into a tagged struct.
package inputs

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

// validType performs a sanity check that the given value is appropriate for
// use in an unmarshal function. Specifically, it checks that the value is a
// non-nil pointer to a non-nil struct.
func validType(v any) error {
	reflectValue := reflect.ValueOf(v)
	reflectType := reflect.TypeOf(v)

	switch {
	case reflectType == nil:
		// Given value was nil.
		return fmt.Errorf("nil")
	case reflectType.Kind() != reflect.Pointer:
		// Given value was not a pointer.
		return fmt.Errorf("non-pointer %v", reflectType)
	case reflectValue.IsNil():
		// Given value was a nil pointer.
		return fmt.Errorf("nil %v", reflectType)
	case reflectType.Elem().Kind() != reflect.Struct:
		// Given value was not a pointer to a struct.
		return fmt.Errorf("non-struct %v", reflectType.Elem().Kind())
	default:
		return nil
	}
}

// Unmarshal updates fields in the given struct with corresponding GitHub
// Actions input values. If a field has an "input" tag, the associated named
// input value will be obtained, and converted into the correct type for the
// underlying field. Struct fields that have no tag, or an "input" tag with a
// name of "-" are ignored.
func Unmarshal(action *githubactions.Action, v any) error {
	if err := validType(v); err != nil {
		return err
	}

	reflectValue := reflect.ValueOf(v).Elem()
	reflectType := reflectValue.Type()

	// Loop over every struct field and attempt to update it with a matching
	// input value.
	for index := 0; index < reflectValue.NumField(); index++ {
		field := reflectValue.Field(index)

		// Ignore fields that cannot be addressed or are private.
		if !field.CanSet() {
			continue
		}

		// Ignore fields that don't have a struct tag.
		tag, ok := reflectType.Field(index).Tag.Lookup("input")
		if !ok {
			continue
		}

		// Determine if the input was flagged as required.
		inputName, required := strings.CutSuffix(tag, ",required")

		// Ignore fields that want to be ignored.
		if inputName == "-" {
			continue
		}

		// Obtain a value for this field using the tagged input name.
		inputValue := action.GetInput(inputName)
		if inputValue == "" {
			if required {
				return fmt.Errorf("no value for required input %s", inputName)
			}
			continue
		}

		// Parse the string input value into a concrete type.
		value, err := convert(field.Type().String(), inputValue)
		if err != nil {
			return fmt.Errorf("field %s: %w", reflectType.Field(index).Name, err)
		}

		// Finally, update the struct field with the new value.
		field.Set(reflect.ValueOf(value))
	}

	return nil
}
