package internal

import (
	"testing"
)

func TestExample(t *testing.T) {
	t.Fatalf("failed %#v", nil)
}

func TestExample2(t *testing.T) {
	t.Fatal("failed")
}
