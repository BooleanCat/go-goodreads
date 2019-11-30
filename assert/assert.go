package assert

import (
	"regexp"
	"testing"
)

func Nil(t *testing.T, v interface{}) {
	if v != nil {
		t.Fatalf(`expected "%v" to be nil`, v)
	}
}

func Equal(t *testing.T, v, w interface{}) {
	if v != w {
		t.Fatalf(`expected "%v" to equal "%v"`, v, w)
	}
}

func ErrorMatches(t *testing.T, err error, pattern string) {
	expression := regexp.MustCompile(pattern)

	if err == nil {
		t.Fatal("expected err to have occurred")
	}

	if !expression.MatchString(err.Error()) {
		t.Fatalf(`expected error "%v" to match pattern "%s"`, err, pattern)
	}
}
