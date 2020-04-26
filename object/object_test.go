package object_test

import (
	"testing"

	"github.com/andy9775/interpreter/object"
)

func TestStringHashKey(t *testing.T) {
	hello1 := &object.String{Value: "Hello World"}
	hello2 := &object.String{Value: "Hello World"}
	diff1 := &object.String{Value: "My name is andy"}
	diff2 := &object.String{Value: "My name is andy"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with the same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with the same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash key")
	}
}