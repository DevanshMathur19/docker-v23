package strslice // import "github.com/DevanshMathur19/docker-v23/api/types/strslice"

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestStrSliceMarshalJSON(t *testing.T) {
	for _, testcase := range []struct {
		input    StrSlice
		expected string
	}{
		// MADNESS(stevvooe): No clue why nil would be "" but empty would be
		// "null". Had to make a change here that may affect compatibility.
		{input: nil, expected: "null"},
		{StrSlice{}, "[]"},
		{StrSlice{"/bin/sh", "-c", "echo"}, `["/bin/sh","-c","echo"]`},
	} {
		data, err := json.Marshal(testcase.input)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != testcase.expected {
			t.Fatalf("%#v: expected %v, got %v", testcase.input, testcase.expected, string(data))
		}
	}
}

func TestStrSliceUnmarshalJSON(t *testing.T) {
	parts := map[string][]string{
		"":                        {"default", "values"},
		"[]":                      {},
		`["/bin/sh","-c","echo"]`: {"/bin/sh", "-c", "echo"},
	}
	for input, expected := range parts {
		strs := StrSlice{"default", "values"}
		if err := strs.UnmarshalJSON([]byte(input)); err != nil {
			t.Fatal(err)
		}

		actualParts := []string(strs)
		if !reflect.DeepEqual(actualParts, expected) {
			t.Fatalf("%#v: expected %v, got %v", input, expected, actualParts)
		}
	}
}

func TestStrSliceUnmarshalString(t *testing.T) {
	var e StrSlice
	echo, err := json.Marshal("echo")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(echo, &e); err != nil {
		t.Fatal(err)
	}

	if len(e) != 1 {
		t.Fatalf("expected 1 element after unmarshal: %q", e)
	}

	if e[0] != "echo" {
		t.Fatalf("expected `echo`, got: %q", e[0])
	}
}

func TestStrSliceUnmarshalSlice(t *testing.T) {
	var e StrSlice
	echo, err := json.Marshal([]string{"echo"})
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(echo, &e); err != nil {
		t.Fatal(err)
	}

	if len(e) != 1 {
		t.Fatalf("expected 1 element after unmarshal: %q", e)
	}

	if e[0] != "echo" {
		t.Fatalf("expected `echo`, got: %q", e[0])
	}
}
