package simulate

import (
	pb "../../api"
	"../../internal"
	"container/heap"
	"fmt"
	"math/rand"
)

func randomEdit(site *internal.Site, slot []rune) *pb.Operation {
	mode := rand.Int31n(100)
	length := len(site.Value())
	if length > 3 && mode%4 == 0 || length > 10 {
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

type Peer struct {
	internal.Site
	slot     []rune
	pq       PriorityQueue
	neighbor *Peer
}

func NewPeers(slots [][]rune) []*Peer {
	peers := make([]*Peer, len(slots))
	for i, slot := range slots {
		peers[i] = &Peer{
			Site: internal.NewSite(fmt.Sprintf("site%v", i), " "),
			slot: slot,
			pq:   make(PriorityQueue, 0),
		}
	}
	for i, p := range peers {
		p.neighbor = peers[(i+1)%len(peers)]
	}
	return peers
}

func (p *Peer) broadcast(op *pb.Operation) {
	for peer := p.neighbor; peer != p; peer = peer.neighbor {
		heap.Push(&peer.pq, &Item{rand.Int(), op})
	}
}

// random: consume some operations
func (p *Peer) consume() {
	item := heap.Pop(&p.pq).(*Item)
	p.Site.Integrate(item.value)
}

// flush: consume all operations
func (p *Peer) flush() {
	for ; len(p.pq) > 0; {
		item := heap.Pop(&p.pq).(*Item)
		p.Site.Integrate(item.value)
	}
}

type predicate func(p *Peer) bool

func filter(peers []*Peer, pred predicate) []*Peer {
	ret := make([]*Peer, 0)
	for _, p := range peers {
		if pred(p) {
			ret = append(ret, p)
		}
	}
	return ret
}

type Info struct {
	Peer  int
	Value string
}

func dump(peers []*Peer, out chan Info) {
	for i, p := range peers {
		out <- Info{i, p.Value()}
	}
}

func Run(peers []*Peer, out chan Info, trial int) {
	for i := 0; i < trial; i++ {
		for _, p := range peers {
			op := randomEdit(&p.Site, p.slot)
			p.broadcast(op)
		}
		for range peers {
			candidates := filter(peers, func(p *Peer) bool {
				return len(p.pq) > 0
			})
			p := rand.Intn(len(candidates))
			candidates[p].consume()
		}
		dump(peers, out)
	}

	for _, p := range peers {
		p.flush()
	}
	dump(peers, out)
	close(out)
}
