package internal

// WChar is W-character c
type WChar struct {
	// Id is the id of c
	Id Wid
	// Visible is {true | false}, if the character c is visible
	Visible bool
	// Alpha is alphabetical value of the effective character of c
	Alpha rune
	// PreviousId is the id of the previous W-character of c
	PreviousId Wid
	// NextId is the id of the next W-character of c
	NextId Wid
}

// Wid is the id of W-character
type Wid struct {
	// The identifier of a site (a peer)
	Ns string
	// The local clock of the W-character is generated on a site
	Ng int64
}

type Editor struct {
	SiteId     string
	LocalClock int64
}

// Special W-character, beginning of the sequence
var Cb = WChar{
	Id:      Wid{"_", 0},
	Visible: false,
}

// Special W-character, ending of the sequence
var Ce = WChar{
	Id:      Wid{"_", 1},
	Visible: false,
}

type Sequence struct {
}

func (s *Sequence) len() int {
	return 0
}

func (s *Sequence) at(position int) *WChar {
	return nil
}

func (s *Sequence) pos(c WChar) int {
	return 0
}

func (s *Sequence) Insert(c rune, position int) {}

func (s *Sequence) Subseq(c Wid, d Wid) *Sequence {
	return nil
}

func (s *Sequence) Contains(c Wid) bool {
	return false
}

func (s *Sequence) Value() *string {
	return nil
}

func (s *Sequence) IthVisible(i int) bool {
	return false
}
