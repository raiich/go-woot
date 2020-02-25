package test

import (
	"../internal"
	"testing"
)

func TestInsert(t *testing.T) {
	site1 := internal.NewSite("site1")
	site2 := internal.NewSite("site2")
	op1 := site1.GenerateIns(0, 'a')
	op2 := site2.GenerateIns(0, '1')
	site2.Integrate(op1)
	site1.Integrate(op2)

	if site1.Value() != site2.Value() {
		t.Fatalf("failed %v != %v", site1.Value(), site2.Value())
	}
}
