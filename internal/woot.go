package internal

import (
	pb "../api"
	"container/list"
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
	String *list.List
}

type Site struct {
	// site identifier
	NumSite string
	// Local Logical clock
	Hs   int64
	Seq  *Sequence
	pool chan *pb.Operation
}

func InitialSequence() *Sequence {
	ls := list.New()
	ls.PushBack(Cb)
	ls.PushBack(Ce)
	return &Sequence{ls}
}

func (s *Sequence) pos(c *pb.Wid) int {
	current := s.String.Front()
	for p := 0; current != nil; p += 1 {
		if current.Value.(*pb.Wchar).Id == c {
			return p
		}
		current = current.Next()
	}
	panic("TODO illegal input")
}

func (s *Sequence) insert(c *pb.Wchar, position int) {
	current := s.String.Front()
	for i := 0; i < position; i += 1 {
		current = current.Next()
	}
	s.String.InsertAfter(c, current)
}

func (s *Sequence) subseq(c *pb.Wid, d *pb.Wid) (int, *list.Element) {
	head := s.findElement(c)
	length := 0
	for current := head; current != nil; current = current.Next() {
		length += 1
		a := current.Value.(*pb.Wchar)
		if equal(a.Id, d) {
			return length, head
		}
	}
	panic("TODO illegal input")
}

func (s *Sequence) contains(c *pb.Wid) bool {
	return s.findElement(c) != nil
}

func (s *Sequence) find(c *pb.Wid) *pb.Wchar {
	if found := s.findElement(c); found != nil {
		return found.Value.(*pb.Wchar)
	}
	return nil
}

func (s *Sequence) findElement(c *pb.Wid) *list.Element {
	for elem := s.String.Front(); elem != nil; elem = elem.Next() {
		if equal(c, elem.Value.(*pb.Wchar).Id) {
			return elem
		}
	}
	return nil
}

func equal(a *pb.Wid, b *pb.Wid) bool {
	return a.Ns == b.Ns && a.Ng == a.Ng
}

func (s *Sequence) value() *string {
	return nil
}

func (s *Sequence) ithVisible(i int) *pb.Wchar {
	return nil
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
	cp := site.Seq.ithVisible(pos)
	cn := site.Seq.ithVisible(pos + 1)
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
	site.Seq.IntegrateIns(wchar, cp.Id, cn.Id)
	return broadcast(ins(wchar))
}

func (site *Site) GenerateDel(pos int) *pb.Operation {
	wchar := site.Seq.ithVisible(pos)
	IntegrateDel(wchar)
	return broadcast(del(wchar))
}

func (s *Sequence) IntegrateIns(c *pb.Wchar, cp *pb.Wid, cn *pb.Wid) {
	length, ss := s.subseq(cp, cn)
	if length == 0 {
		s.insert(c, s.pos(cn))
	} else {
		var l []*pb.Wid
		l = append(l, cp)
		current := ss
		for i := 0; i < length; i += 1 {
			di := current.Value.(*pb.Wchar)
			if s.pos(di.PreviousId) <= s.pos(cp) && s.pos(cn) <= s.pos(di.NextId) {
				l = append(l, di.Id)
			}
			current = current.Next()
		}

		var prev, next *pb.Wid
		for i, li := range l {
			if compare(li, c.Id) < 0 {
				prev = l[i-1]
				next = li
				break
			}
		}
		s.IntegrateIns(c, prev, next)
	}
}

func compare(a *pb.Wid, b *pb.Wid) int {
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

func IntegrateDel(wchar *pb.Wchar) {
	wchar.Visible = false
}

func broadcast(op *pb.Operation) *pb.Operation {
	return op
}

func (site *Site) isExecutable(op *pb.Operation) bool {
	c := op.C
	if op.Type == pb.OperationType_DELETE {
		return site.Seq.contains(c.Id)
	} else {
		cp, cn := c.NextId, c.PreviousId
		return site.Seq.contains(cp) && site.Seq.contains(cn)
	}
}

func (site *Site) Reception(op *pb.Operation) {
	site.pool <- op
}

func (site *Site) Main() {
	for {
		for op := range site.pool {
			if site.isExecutable(op) {
				if op.Type == pb.OperationType_DELETE {
					c := site.Seq.find(op.C.Id)
					IntegrateDel(c)
				} else {
					c := op.C
					site.Seq.IntegrateIns(c, c.PreviousId, c.NextId)
				}
			}
		}
	}
}
