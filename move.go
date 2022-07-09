package chess

import (
	"encoding/json"
	"errors"
)

// A MoveTag represents a notable consequence of a move.
type MoveTag uint16

const (
	// KingSideCastle indicates that the move is a king side castle.
	KingSideCastle MoveTag = 1 << iota
	// QueenSideCastle indicates that the move is a queen side castle.
	QueenSideCastle
	// Capture indicates that the move captures a piece.
	Capture
	// EnPassant indicates that the move captures via en passant.
	EnPassant
	// Check indicates that the move puts the opposing player in check.
	Check
	// inCheck indicates that the move puts the moving player in check and
	// is therefore invalid.
	inCheck
)

// A Move is the movement of a piece from one square to another.
type Move struct {
	s1    Square
	s2    Square
	promo PieceType
	tags  MoveTag
}

func NewMove(s1 Square, s2 Square, promo PieceType, tags MoveTag) Move {
	return Move{
		s1: s1,
		s2: s2,
		promo: promo,
		tags: tags,
	}
}

// String returns a string useful for debugging.  String doesn't return
// algebraic notation.
func (m *Move) String() string {
	return m.s1.String() + m.s2.String() + m.promo.String()
}

// S1 returns the origin square of the move.
func (m *Move) S1() Square {
	return m.s1
}

// S2 returns the destination square of the move.
func (m *Move) S2() Square {
	return m.s2
}

// Promo returns promotion piece type of the move.
func (m *Move) Promo() PieceType {
	return m.promo
}

// HasTag returns true if the move contains the MoveTag given.
func (m *Move) HasTag(tag MoveTag) bool {
	return (tag & m.tags) > 0
}

func (m *Move) addTag(tag MoveTag) {
	m.tags = m.tags | tag
}

func (m Move) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(m.s1.String() + m.s2.String() + m.promo.String()))
}

func (m *Move) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if len(s) < 4 {
		return errors.New("chess: unable to unmarshal move: incorrect data length")
	}
	var ok bool
	if m.s1, ok = strToSquareMap[s[0:2]]; !ok {
		return errors.New("chess: unable to unmarshal move: invalid src square")
	}
	if m.s2, ok = strToSquareMap[s[2:4]]; !ok {
		return errors.New("chess: unable to unmarshal move: invalid dst square")
	}
	if len(s) > 4 {
		if m.promo, ok = strToPieceTypeMap[s[4:5]]; !ok {
			return errors.New("chess: unable to unmarshal move: invalid promo piece type")
		}
	} else {
		m.promo = NoPieceType
	}
	m.tags = MoveTag(0)
	return nil
}

type moveSlice []*Move

func (a moveSlice) find(m *Move) *Move {
	if m == nil {
		return nil
	}
	for _, move := range a {
		if move.String() == m.String() {
			return move
		}
	}
	return nil
}
