package assert

import (
	"reflect"
	"strings"
	"testing"
)

func Equal[V any](t *testing.T, got, expected V) {
	t.Helper()
	if reflect.DeepEqual(got, expected) {
		return
	} else {
		t.Errorf(`assert.Equal(
got: %v
expected: %v
)`, got, expected)
		return
	}
}

func Nil(t *testing.T, value any) {
	t.Helper()

	if value == nil {
		return
	}

	v := reflect.ValueOf(value)
	if !v.IsNil() {
		t.Errorf(`expected nil, got %v`, value)
	}
}

func EqualStringSlices(t *testing.T, got, expected []string) {
	t.Helper()

	if !reflect.DeepEqual(got, expected) {
		t.Errorf(`assert.EqualStringSlices(
got: %v
expected: %v
)`, got, expected)
	}
}

func Contains(t *testing.T, str string, element string) {
	t.Helper()

	if !strings.Contains(str, element) {
		t.Errorf(`assert.Contains(
string: %v
element: %v
)`, str, element)
	}
}

func HasError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Errorf(`expected error, got nil`)
	}
}

func NoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf(`expected no error, got %v`, err)
	}
}
