package test

import (
	"../examples/simulate"
	"testing"
)

func discard(ch chan simulate.Info) {
	for range ch {
	}
}

func TestRandomOperation(t *testing.T) {
	trial := 3
	var slots = [][]rune{
		[]rune("abcdefghij"),
		[]rune("1234567890"),
		[]rune("^*-/_=+?:."),
	}
	peers := simulate.NewPeers(slots)

	out := make(chan simulate.Info)
	go discard(out)
	simulate.Run(peers, out, trial)

	a, b, c := peers[0].Value(), peers[1].Value(), peers[2].Value()
	if a != b || b != c {
		t.Fatalf("failed %v %v %v", a, b, c)
	}
}
