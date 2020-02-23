package internal

import (
	pb "../api"
)

type Editor struct {
	SiteId     string
	LocalClock int64
}

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
}

func (s *Sequence) Len() int {
	return 0
}

func (s *Sequence) At(position int) *pb.Wchar {
	return nil
}

func (s *Sequence) Pos(c *pb.Wchar) int {
	return 0
}

func (s *Sequence) Insert(c rune, position int) {}

func (s *Sequence) Subseq(c *pb.Wid, d *pb.Wid) *Sequence {
	return nil
}

func (s *Sequence) Contains(c *pb.Wid) bool {
	return false
}

func (s *Sequence) Value() *string {
	return nil
}

func (s *Sequence) IthVisible(i int) bool {
	return false
}

func Ins(c *pb.Wchar) *pb.Operation {
	return &pb.Operation{
		Type: pb.OperationType_INSERT,
		C: c,
	}
}

func Del(c *pb.Wchar) *pb.Operation {
	return &pb.Operation{
		Type: pb.OperationType_DELETE,
		C: c,
	}
}
