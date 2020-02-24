package internal

import (
	pb "../api"
	"container/list"
	"strings"
)

// Special W-character, beginning of the sequence
var Cb = pb.Wchar{
	Id:      &pb.Wid{Ns: "_", Ng: 0},
	Visible: false,
}

// Special W-character, ending of the sequence
var Ce = pb.Wchar{
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
}

func InitialSequence() Sequence {
	ls := list.New()
	ls.PushBack(Cb)
	ls.PushBack(Ce)
	return Sequence{ls}
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
	s.stringS.InsertAfter(c, neighbor)
}

func (s Sequence) subseq(c pb.Wid, d pb.Wid) (int, *list.Element) {
	_, head := s.head().findElementById(c)
	length, _ := SubSeq{head}.findElementById(d)
	return length, head
}

func (s Sequence) contains(c pb.Wid) bool {
	_, found := s.head().findElementById(c)
	return found != nil
}

func (s Sequence) value() string {
	var builder strings.Builder
	for head := s.stringS.Front(); head != nil ; head = head.Next() {
		wchar := head.Value.(*pb.Wchar)
		if wchar.Visible {
			builder.WriteRune(wchar.CodePoint)
		}
	}
	return builder.String()
}

func (s Sequence) ithVisible(i int) *pb.Wchar {
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
	return found.Value.(*pb.Wchar)
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
	return broadcast(ins(wchar))
}

func (site *Site) GenerateDel(pos int) *pb.Operation {
	wchar := site.seq.ithVisible(pos)
	site.IntegrateDel(wchar)
	return broadcast(del(wchar))
}

func (site *Site) IntegrateIns(c *pb.Wchar, cp pb.Wid, cn pb.Wid) {
	s := site.seq
	length, ss := s.subseq(cp, cn)
	if length == 0 {
		s.insert(c, s.pos(cn))
	} else {
		var l []pb.Wid
		l = append(l, cp)
		current := ss
		for i := 0; i < length; i += 1 {
			di := current.Value.(*pb.Wchar)
			if s.pos(*di.PreviousId) <= s.pos(cp) && s.pos(cn) <= s.pos(*di.NextId) {
				l = append(l, *di.Id)
			}
			current = current.Next()
		}

		var prev, next pb.Wid
		for i, li := range l {
			if compare(li, c.Id) < 0 {
				prev = l[i-1]
				next = li
				break
			}
		}
		site.IntegrateIns(c, prev, next)
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

func (site *Site)IntegrateDel(wchar *pb.Wchar) {
	wchar.Visible = false
}

func broadcast(op *pb.Operation) *pb.Operation {
	return op
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

func (site *Site) Integrate(op *pb.Operation) {
	if op.Type == pb.OperationType_DELETE {
		_, c := site.seq.head().findElementById(*op.C.Id)
		site.IntegrateDel(c.Value.(*pb.Wchar))
	} else {
		c := op.C
		site.IntegrateIns(c, *c.PreviousId, *c.NextId)
	}
}

func (site *Site) Main() {
	for {
		for op := range site.pool {
			if site.isExecutable(op) {
				site.Integrate(op)
			} else {
				site.pool <- op
			}
		}
	}
}

type predicate func(c *pb.Wchar) bool

type SubSeq struct {
	head *list.Element
}

func (s Sequence) head() SubSeq{
	return SubSeq{s.stringS.Front()}
}

func (s SubSeq) find(pred predicate) (int, *list.Element) {
	c := 0
	for elem := s.head; elem != nil; elem = elem.Next() {
		if pred(elem.Value.(*pb.Wchar)) {
			return c, elem
		} else {
			c += 1
		}
	}
	return -1, nil
}

func (s SubSeq) findElementById(wid pb.Wid) (int, *list.Element) {
	return s.find(func(c *pb.Wchar) bool {
		return equal(c.Id, wid)
	})
}

func equal(a *pb.Wid, b pb.Wid) bool {
	return a.Ns == b.Ns && a.Ng == a.Ng
}
