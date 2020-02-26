package internal

import (
	pb "../api"
	"container/list"
	"strings"
)

// Special W-character, beginning of the sequence
var cb = pb.Wchar{
	Id:      &pb.Wid{Ns: "_", Ng: 0},
	Visible: false,
}

// Special W-character, ending of the sequence
var ce = pb.Wchar{
	Id:      &pb.Wid{Ns: "_", Ng: 1},
	Visible: false,
}

type Sequence struct {
	// sequence of W-characters
	stringS *list.List
}

type Site struct {
	// site identifier
	NumSite string
	// Local Logical clock
	Hs   int64
	seq  Sequence
	pool chan *pb.Operation
	queue []*pb.Operation
}

func initialSequence(initial string) Sequence {
	ls := list.New()
	ls.PushBack(&cb)

	for i, c := range initial {
		ls.PushBack(&pb.Wchar{
			Id:      &pb.Wid{Ns: "_", Ng: int64(2 + i)},
			CodePoint: c,
			Visible: true,
		})
	}
	ls.PushBack(&ce)
	return Sequence{ls}
}

func NewSite(id string, initial string) Site {
	site := Site{
		id,
		0,
		initialSequence(initial),
		make(chan *pb.Operation),
		make([]*pb.Operation, 0),
	}
	return site
}

func (s Sequence) pos(wid pb.Wid) int {
	p, found := s.head().findElementById(wid)
	if found != nil {
		return p
	} else {
		panic("TODO illegal input")
	}
}

func (s Sequence) insert(c *pb.Wchar, i int) {
	_, neighbor := s.head().find(func(_ *pb.Wchar) bool {
		if i == 0 {
			return true
		} else {
			i -= 1
			return false
		}
	})
	s.stringS.InsertBefore(c, neighbor.elem)
}

func (s Sequence) subseq(c pb.Wid, d pb.Wid) (int, *SubSeq) {
	_, head := s.head().findElementById(c)
	ret := head.Next()
	length, _ := ret.findElementById(d)
	return length, ret
}

func (s Sequence) contains(c pb.Wid) bool {
	_, found := s.head().findElementById(c)
	return found != nil
}

func (s Sequence) value() string {
	var builder strings.Builder
	for head := s.head(); head != nil; head = head.Next() {
		wchar := head.Val()
		if wchar.Visible {
			builder.WriteRune(wchar.CodePoint)
		}
	}
	return builder.String()
}

func (s Sequence) ithVisible(i int) *pb.Wchar {
	if i == 0 {
		return &cb
	}
	i -= 1
	_, found := s.head().find(func(c *pb.Wchar) bool {
		if c.Visible {
			if i == 0 {
				return true
			} else {
				i -= 1
			}
		}
		return false
	})
	if found != nil {
		return found.Val()
	} else {
		return &ce
	}

}

func ins(c *pb.Wchar) *pb.Operation {
	return &pb.Operation{
		Type: pb.OperationType_INSERT,
		C:    c,
	}
}

func del(c *pb.Wchar) *pb.Operation {
	return &pb.Operation{
		Type: pb.OperationType_DELETE,
		C:    c,
	}
}

func (site *Site) GenerateIns(pos int, alpha rune) *pb.Operation {
	site.Hs += 1
	cp := site.seq.ithVisible(pos)
	cn := site.seq.ithVisible(pos + 1)
	wid := pb.Wid{
		Ns: site.NumSite,
		Ng: site.Hs,
	}
	wchar := &pb.Wchar{
		Id:         &wid,
		CodePoint:  alpha,
		Visible:    true,
		PreviousId: cp.Id,
		NextId:     cn.Id,
	}
	site.IntegrateIns(wchar, *cp.Id, *cn.Id)
	return ins(wchar)
	// TODO broadcast
}

func (site *Site) GenerateDel(pos int) *pb.Operation {
	wchar := site.seq.ithVisible(pos + 1)
	site.IntegrateDel(wchar)
	return del(wchar)
	// TODO broadcast
}

func (site *Site) IntegrateIns(c *pb.Wchar, cp pb.Wid, cn pb.Wid) {
	s := site.seq
	length, ss := s.subseq(cp, cn)
	if length == 0 {
		s.insert(c, s.pos(cn))
	} else {
		var l []pb.Wid
		l = append(l, cp)

		for i := 0; i < length; i += 1 {
			di := ss.Val()
			if s.pos(*di.PreviousId) <= s.pos(cp) && s.pos(cn) <= s.pos(*di.NextId) {
				l = append(l, *di.Id)
			}
			ss = ss.Next()
		}
		l = append(l, cn)

		i := 1
		for ; i < len(l) - 1 && compare(l[i], c.Id) < 0; {
			i += 1
		}
		site.IntegrateIns(c, l[i-1], l[i])
	}
}

func compare(a pb.Wid, b *pb.Wid) int {
	if a.Ns == b.Ns {
		return int(a.Ng - b.Ng)
	} else {
		if a.Ns < b.Ns {
			return -1
		} else if a.Ns > b.Ns {
			return 1
		} else {
			return 0
		}
	}
}

func (site *Site) IntegrateDel(wchar *pb.Wchar) {
	wchar.Visible = false
}

func (site *Site) isExecutable(op *pb.Operation) bool {
	c := op.C
	if op.Type == pb.OperationType_DELETE {
		return site.seq.contains(*c.Id)
	} else {
		cp, cn := c.NextId, c.PreviousId
		return site.seq.contains(*cp) && site.seq.contains(*cn)
	}
}

func (site *Site) Reception(op *pb.Operation) {
	site.pool <- op
}

func (site *Site) integrate(op *pb.Operation) {
	if op.Type == pb.OperationType_DELETE {
		_, c := site.seq.head().findElementById(*op.C.Id)
		site.IntegrateDel(c.Val())
	} else {
		c := op.C
		site.IntegrateIns(c, *c.PreviousId, *c.NextId)
	}
}

func (site *Site) Integrate(op *pb.Operation) int {
	consumed := 0
	site.queue = append(site.queue, op)
	for {
		l := len(site.queue)
		pool := make([]*pb.Operation, 0)
		for _, op := range site.queue {
			if site.isExecutable(op) {
				site.integrate(op)
			} else {
				pool = append(pool, op)
			}
		}
		site.queue = pool
		if len(site.queue) == l {
			// need missing operation (wait operation to arrive)
			break
		} else {
			consumed += l - len(site.queue)
		}
	}
	return consumed
}

func (site *Site) Main() {
	for {
		for op := range site.pool {
			if site.isExecutable(op) {
				site.integrate(op)
			} else {
				site.pool <- op
			}
		}
	}
}

func (site *Site) Value() string {
	return site.seq.value()
}

func (site *Site) Raw() *SubSeq {
	return site.seq.head()
}

type predicate func(c *pb.Wchar) bool

type SubSeq struct {
	elem *list.Element
}

func (s Sequence) head() *SubSeq {
	return &SubSeq{s.stringS.Front()}
}

func (s SubSeq) find(pred predicate) (int, *SubSeq) {
	c := 0
	for elem := &s; elem != nil; elem = elem.Next() {
		if pred(elem.Val()) {
			return c, elem
		} else {
			c += 1
		}
	}
	return -1, nil
}

func (s SubSeq) findElementById(wid pb.Wid) (int, *SubSeq) {
	return s.find(func(c *pb.Wchar) bool {
		return Equal(c.Id, wid)
	})
}

func (s SubSeq) Next() *SubSeq {
	next := s.elem.Next()
	if next != nil {
		return &SubSeq{next}
	} else {
		return nil
	}
}

func (s SubSeq) Val() *pb.Wchar {
	return s.elem.Value.(*pb.Wchar)
}

func Equal(a *pb.Wid, b pb.Wid) bool {
	return a.Ns == b.Ns && a.Ng == b.Ng
}
