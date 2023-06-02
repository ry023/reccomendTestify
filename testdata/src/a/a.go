package a

import (
	"testing"
)

func TestSomething(t *testing.T) {
	t.Errorf("Something Error: %s", "foobar") // want "assert.Failf"
	t.Error("Something Error")                // want "assert.Fail"
	t.Fatalf("Something Error: %s", "foobar") // want "assert.FailNowf"
	t.Fatal("Something Error")                // want "assert.FailNow"
}
