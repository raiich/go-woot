package simulate

import (
	pb "../../api"
	"../../internal"
	"container/heap"
	"fmt"
	"github.com/golang/protobuf/proto"
	"math"
	"math/rand"
)

func randomEdit(site *internal.Site, slot []rune) *pb.Operation {
	mode := rand.Int31n(100)
	length := len(site.Value())
	if length > 2 && mode%2 == 0 || length > 5 {
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
			Site: internal.NewSite(fmt.Sprintf("site%v", i), "A"),
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
		priority := rand.Int()
		if op.Type == pb.OperationType_DELETE {
			//priority /= 10000
			priority = int(math.Log(float64(priority)))
		}
		op = proto.Clone(op).(*pb.Operation)
		op.C = proto.Clone(op.C).(*pb.Wchar)
		heap.Push(&peer.pq, &Item{priority, op})
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
	Operation string
	Value string
}

func dump(peers []*Peer, op string, out chan Info) {
	for i, p := range peers {
		out <- Info{i, op, p.Value()}
	}
}

func Run(peers []*Peer, out chan Info, trial int) {
	for i := 0; i < trial; i++ {
		for _, p := range peers {
			op := randomEdit(&p.Site, p.slot)
			p.broadcast(op)
		}
		dump(peers, "Edit ", out)

		for range peers {
			candidates := filter(peers, func(p *Peer) bool {
				return len(p.pq) > 0
			})
			p := rand.Intn(len(candidates))
			candidates[p].consume()
		}
		dump(peers, "Merge", out)
	}

	for _, p := range peers {
		p.flush()
	}
	dump(peers, "merge", out)
	close(out)
}
