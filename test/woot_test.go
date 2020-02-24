package test

import (
	"../internal"
	"testing"
)

func TestWoot(t *testing.T) {
	site1 := internal.NewSite("site1")
	site2 := internal.NewSite("site2")
	op1 := site1.GenerateIns(0, 'a')
	op2 := site2.GenerateIns(0, '1')
	site2.Integrate(op1)
	site1.Integrate(op2)

	println(site1.Value())
	println(site2.Value())


	//site1.GenerateDel()

	t.Fatalf("failed %#v", nil)
}

func TestExample(t *testing.T) {
	t.Fatalf("failed %#v", nil)
}

func TestExample2(t *testing.T) {
	t.Fatal("failed")
}
