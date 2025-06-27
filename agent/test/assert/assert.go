package assert

import (
	"reflect"
	"strings"
	"testing"
)

func Equal[V any](t *testing.T, got, want V) {
	t.Helper()
	if reflect.DeepEqual(got, want) {
		return
	} else {
		t.Errorf(`assert.Equal(
got: %v
want: %v
)`, got, want)
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
		t.Errorf(`want nil, got %v`, value)
	}
}

func EqualStringSlices(t *testing.T, got, want []string) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf(`assert.EqualStringSlices(
got: %v
want: %v
)`, got, want)
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
		t.Errorf(`want error, but got nil`)
	}
}

func NoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf(`want no error, but got %v`, err)
	}
}
