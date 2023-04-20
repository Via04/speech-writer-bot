package utils

import (
	"testing"
)

func TestGetNameNoExt(t *testing.T) {
	name, err := GetNameNoExt("test.txt")
	if err != nil {
		t.Fatalf("failed function: %s", err.Error())
	}
	if name != "test" {
		t.Fatalf("want test, got %s", name)
	}
}
