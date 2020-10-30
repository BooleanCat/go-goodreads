package assert

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// Nil succeeds when the value is nil.
func Nil(t *testing.T, v interface{}) {
	if v != nil {
		t.Fatalf(`expected "%v" to be nil`, v)
	}
}

// Equal succeeds when the two values are equal.
func Equal(t *testing.T, v, w interface{}) {
	if !reflect.DeepEqual(v, w) {
		t.Fatalf(`expected "%v" to equal "%v"`, v, w)
	}
}

// EndsWith succeeds when the input string ends with the given suffix.
func EndsWith(t *testing.T, s, u string) {
	if !strings.HasSuffix(s, u) {
		t.Fatalf(`expected "%s" to end with "%s"`, s, u)
	}
}

// ErrorMatches succeeds when the error string matches the given pattern.
func ErrorMatches(t *testing.T, err error, pattern string) {
	if err == nil {
		t.Fatal("expected err to have occurred")
	}

	if !regexp.MustCompile(pattern).MatchString(err.Error()) {
		t.Fatalf(`expected error "%v" to match pattern "%s"`, err, pattern)
	}
}

// DoesNotContainSubstring succeeds when the input string doesn't contain the given substring.
func DoesNotContainSubstring(t *testing.T, x, y string) {
	if strings.Contains(x, y) {
		t.Fatalf(`expected "%s" not to contain "%s"`, x, y)
	}
}

// True succeeds when the input is true.
func True(t *testing.T, v interface{}) {
	if !reflect.DeepEqual(v, true) {
		t.Fatal("expected true")
	}
}
