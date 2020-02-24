package test

import (
	"../internal"
	"testing"
)

func TestWoot(t *testing.T) {
	site1 := internal.Site{NumSite: "site1", Hs: 0}
	site2 := internal.Site{NumSite: "site1", Hs: 0}
	op1 := site1.GenerateIns(0, 'a')
	op2 := site2.GenerateIns(0, '1')
	site2.Integrate(op1)
	site1.Integrate(op2)

	//site1.GenerateDel()

	t.Fatalf("failed %#v", nil)
}

func TestExample(t *testing.T) {
	t.Fatalf("failed %#v", nil)
}

func TestExample2(t *testing.T) {
	t.Fatal("failed")
}
