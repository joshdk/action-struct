// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package inputs

import (
	"bytes"
	"testing"
	"time"

	"github.com/sethvargo/go-githubactions"
)

func TestValidType(t *testing.T) {
	// A nil destination is invalid.
	if err := validType(nil); err == nil {
		t.Errorf("expected an error")
	}

	// A non-pointer destination is invalid.
	if err := validType(struct{}{}); err == nil {
		t.Errorf("expected an error")
	}

	// A pointer to a nil destination is invalid.
	var nilptr *struct{}
	if err := validType(nilptr); err == nil {
		t.Errorf("expected an error")
	}

	// A non-struct destination is invalid.
	var str string
	if err := validType(&str); err == nil {
		t.Errorf("expected an error")
	}

	// A non-nil pointer to a non-nil struct is valid.
	if err := validType(&struct{}{}); err != nil {
		t.Errorf("expected no error bot got %v", err)
	}
}

func TestUnmarshal(t *testing.T) {
	env := map[string]string{
		"INPUT_BOOL":     "true",
		"INPUT_FLOAT":    "3.14",
		"INPUT_DOUBLE":   "3.14159",
		"INPUT_INT":      "9001",
		"INPUT_STRING":   "foo",
		"INPUT_DURATION": "1m9s",
		"INPUT_TIME":     "2006-01-02T15:04:05Z",
		"INPUT_LIST":     "foo,bar,baz",
		"INPUT_RAW":      `{"foo": "bar"}`,
	}

	actions := githubactions.New(githubactions.WithGetenv(func(key string) string {
		return env[key]
	}))

	type target struct {
		Bool     bool          `input:"bool"`
		Float32  float32       `input:"float"`
		Float64  float64       `input:"double"`
		Int      int           `input:"int"`
		String   string        `input:"string"`
		Duration time.Duration `input:"duration"`
		Time     time.Time     `input:"time"`
		List     []string      `input:"list"`
		Bytes    []byte        `input:"raw"`
	}

	var actual target
	if err := Unmarshal(actions, &actual); err != nil {
		t.Fatal(err)
	}

	switch {
	case actual.Bool != true:
		t.Error("bool field improperly unmarshalled")
	case actual.Float32 == 0:
		t.Error("float32 field improperly unmarshalled")
	case actual.Float64 == 0:
		t.Error("float64 field improperly unmarshalled")
	case actual.Int != 9001:
		t.Error("int field improperly unmarshalled")
	case actual.String != "foo":
		t.Error("string field improperly unmarshalled")
	case actual.Duration != time.Minute+time.Second*9:
		t.Error("time.Duration field improperly unmarshalled")
	case actual.Time != time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC):
		t.Error("time.Time field improperly unmarshalled")
	case len(actual.List) != 3 || actual.List[0] != "foo" || actual.List[1] != "bar" || actual.List[2] != "baz":
		t.Error("[]string field improperly unmarshalled")
	case bytes.Compare(actual.Bytes, []byte(`{"foo": "bar"}`)) != 0:
		t.Error("[]byte field improperly unmarshalled")
	}
}
