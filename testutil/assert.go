package testutil

import (
	"reflect"
	"strings"
	"testing"
)

// AssertNoError fails the test if err is not nil.
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// AssertError fails the test if err is nil.
func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// AssertErrorContains fails the test if err is nil or doesn't contain substr.
func AssertErrorContains(t *testing.T, err error, substr string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), substr) {
		t.Errorf("expected error to contain %q, got %q", substr, err.Error())
	}
}

// AssertEqual fails the test if got != want.
func AssertEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// AssertNotEqual fails the test if got == want.
func AssertNotEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want not equal", got)
	}
}

// AssertContains fails the test if s does not contain substr.
func AssertContains(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Errorf("expected %q to contain %q", s, substr)
	}
}

// AssertNotContains fails the test if s contains substr.
func AssertNotContains(t *testing.T, s, substr string) {
	t.Helper()
	if strings.Contains(s, substr) {
		t.Errorf("expected %q to not contain %q", s, substr)
	}
}

// AssertTrue fails the test if condition is false.
func AssertTrue(t *testing.T, condition bool, msg string) {
	t.Helper()
	if !condition {
		t.Errorf("expected true: %s", msg)
	}
}

// AssertFalse fails the test if condition is true.
func AssertFalse(t *testing.T, condition bool, msg string) {
	t.Helper()
	if condition {
		t.Errorf("expected false: %s", msg)
	}
}

// AssertNil fails the test if value is not nil.
func AssertNil(t *testing.T, value interface{}) {
	t.Helper()
	if value != nil && !reflect.ValueOf(value).IsNil() {
		t.Errorf("expected nil, got %v", value)
	}
}

// AssertNotNil fails the test if value is nil.
func AssertNotNil(t *testing.T, value interface{}) {
	t.Helper()
	if value == nil || reflect.ValueOf(value).IsNil() {
		t.Error("expected non-nil value")
	}
}

// AssertLen fails the test if the length of v is not expected.
func AssertLen(t *testing.T, v interface{}, expected int) {
	t.Helper()
	rv := reflect.ValueOf(v)
	if rv.Len() != expected {
		t.Errorf("expected length %d, got %d", expected, rv.Len())
	}
}

// AssertEmpty fails the test if v is not empty.
func AssertEmpty(t *testing.T, v interface{}) {
	t.Helper()
	rv := reflect.ValueOf(v)
	if rv.Len() != 0 {
		t.Errorf("expected empty, got length %d", rv.Len())
	}
}

// AssertNotEmpty fails the test if v is empty.
func AssertNotEmpty(t *testing.T, v interface{}) {
	t.Helper()
	rv := reflect.ValueOf(v)
	if rv.Len() == 0 {
		t.Error("expected non-empty")
	}
}
