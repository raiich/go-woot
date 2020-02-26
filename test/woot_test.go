package test

import (
	"../internal"
	"testing"
)

func TestInsert(t *testing.T) {
	site1 := internal.NewSite("site1", "")
	site2 := internal.NewSite("site2", "")
	op1 := site1.GenerateIns(0, 'a')
	op2 := site2.GenerateIns(0, '1')
	site2.Integrate(op1)
	site1.Integrate(op2)

	if !(site1.Value() == "a1" && site1.Value() == site2.Value()) {
		t.Fatalf("failed %v != %v", site1.Value(), site2.Value())
	}
}

func TestExample(t *testing.T) {
	site1 := internal.NewSite("site1", "ab")
	site2 := internal.NewSite("site2", "ab")
	site3 := internal.NewSite("site3", "ab")

	println(site1.Value())


	op1 := site1.GenerateIns(1, '1')
	op2 := site2.GenerateIns(1, '2')
	op3 := site1.GenerateIns(1, '3')

	site3.Integrate(op1)
	site1.Integrate(op2)
	site3.Integrate(op2)
	site3.Integrate(op3)
	site2.Integrate(op3)
	site2.Integrate(op1)

	if !(site1.Value() == "a312b" && site1.Value() == site2.Value() && site1.Value() == site3.Value()) {
		t.Fatalf("failed %v, %v, %v", site1.Value(), site2.Value(), site3.Value())
	}
}

func TestInsDel1(t *testing.T) {
	site1 := internal.NewSite("site1", "abcde")
	site2 := internal.NewSite("site2", "abcde")
	op1 := site1.GenerateIns(1, '1')
	op2 := site2.GenerateDel(2)
	site1.Integrate(op2)
	site2.Integrate(op1)

	if !(site1.Value() == "a1bde" && site1.Value() == site2.Value()) {
		t.Fatalf("failed %v != %v", site1.Value(), site2.Value())
	}
}

func TestInsIns1(t *testing.T) {
	site1 := internal.NewSite("site1", "ab")
	site2 := internal.NewSite("site2", "ab")
	site3 := internal.NewSite("site3", "ab")

	op := site1.GenerateIns(0, 'a')
	site2.Integrate(op)
	site3.Integrate(op)

	op = site2.GenerateIns(1, 'b')
	site1.Integrate(op)
	site3.Integrate(op)

	op1 := site1.GenerateIns(1, '1')
	op2 := site2.GenerateIns(1, '2')

	site3.Integrate(op1)
	op3 := site3.GenerateIns(1, '3')

	site1.Integrate(op3)
	site1.Integrate(op2)
	site2.Integrate(op3)
	site2.Integrate(op1)
	site3.Integrate(op2)

	if !(site1.Value() == "a31b" && site1.Value() == site2.Value() && site1.Value() == site3.Value()) {
		t.Fatalf("failed %v, %v, %v", site1.Value(), site2.Value(), site3.Value())
	}
}

func TestInsDel(t *testing.T) {
	site1 := internal.NewSite("site1", "abc")
	site2 := internal.NewSite("site2", "abc")
	op1 := site1.GenerateIns(2, '1')
	op2 := site2.GenerateDel(2)

	site1.Integrate(op2)
	site2.Integrate(op1)
	if !(site1.Value() == "ab1" && site1.Value() == site2.Value()) {
		t.Fatalf("failed %v, %v", site1.Value(), site2.Value())
	}
}

func TestDelDel(t *testing.T) {
	site1 := internal.NewSite("site1", "abc")
	site2 := internal.NewSite("site2", "abc")
	op1 := site1.GenerateDel(1)
	op2 := site2.GenerateDel(2)

	site1.Integrate(op2)
	site2.Integrate(op1)
	if !(site1.Value() == "a" && site1.Value() == site2.Value()) {
		t.Fatalf("failed %v, %v", site1.Value(), site2.Value())
	}
}