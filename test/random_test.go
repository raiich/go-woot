package test

import (
	pb "../api"
	"../internal"
	"container/heap"
	"math/rand"
	"testing"
)

func randomEdit(site *internal.Site, slot []rune) *pb.Operation {
	mode := rand.Int31n(100)
	length := len(site.Value())
	if length > 5 && mode%4 == 0 || length > 50 {
		pos := rand.Intn(len(site.Value()))
		return site.GenerateDel(pos)
	} else {
		c := slot[rand.Intn(len(slot))]
		pos := rand.Intn(len(site.Value()))
		return site.GenerateIns(pos, c)
	}
}

type Item struct {
	priority int
	value    *pb.Operation
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

func TestRandomOperation(t *testing.T) {
	sites := []internal.Site{
		internal.NewSite("site1", " "),
		internal.NewSite("site2", " "),
		internal.NewSite("site3", " "),
	}

	pq := []PriorityQueue{
		make(PriorityQueue, 0),
		make(PriorityQueue, 0),
		make(PriorityQueue, 0),
	}

	for i := 0; i < 38; i += 1 {
		ops := []*pb.Operation{
			randomEdit(&sites[0], []rune("abcdefghij")),
			randomEdit(&sites[1], []rune("1234567890")),
			randomEdit(&sites[2], []rune("^*-/_=+?:.")),
		}

		// random: consume some operations
		for i := 0; i < 3; i++ {
			heap.Push(&pq[i], &Item{rand.Int(), ops[(i+1)%3]})
			heap.Push(&pq[i], &Item{rand.Int(), ops[(i+2)%3]})
			for c := rand.Intn(len(pq[i])); c >= 0; c -= 1 {
				item := heap.Pop(&pq[i]).(*Item)
				sites[i].Integrate(item.value)
			}
		}
	}

	// flush: consume all operations
	for i := 0; i < 3; i++ {
		for c := len(pq[i]) - 1; c >= 0; c -= 1 {
			item := heap.Pop(&pq[i]).(*Item)
			sites[i].Integrate(item.value)
		}
	}

	a, b, c := sites[0].Value(), sites[1].Value(), sites[2].Value()
	if a != b || b != c {
		t.Fatalf("failed %v %v %v", a, b, c)
	}
}
