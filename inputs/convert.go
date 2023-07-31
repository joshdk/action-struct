// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package inputs

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func convert(kind, value string) (any, error) {
	switch kind {
	case "bool":
		return strconv.ParseBool(value)

	case "float32":
		value, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return nil, err
		}

		return float32(value), nil

	case "float64":
		return strconv.ParseFloat(value, 64)

	case "int":
		return strconv.Atoi(value)

	case "string":
		return value, nil

	case "time.Duration":
		return time.ParseDuration(value)

	case "time.Time":
		return time.Parse(time.RFC3339, value)

	case "[]string":
		return split(value), nil

	case "[]uint8":
		// This is not "[]byte" since byte is just a type alias for uint8.
		return []byte(value), nil

	default:
		return nil, fmt.Errorf("unsupported type %s", kind)
	}
}

// split the given string on newlines and commas, trim, and return the
// non-blank fields.
func split(value string) (fields []string) {
	for _, field := range strings.FieldsFunc(value, func(r rune) bool {
		return r == '\n' || r == ','
	}) {
		field = strings.TrimSpace(field)
		if field != "" {
			fields = append(fields, field)
		}
	}

	return
}
