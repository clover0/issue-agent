package assert

import (
	"reflect"
	"testing"
)

func Equal[V comparable](t *testing.T, got, expected V) {
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
