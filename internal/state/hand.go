package state

import (
	c "github.com/quibbble/go-boardgame/pkg/collection"
	"github.com/quibbble/go-quill/internal/state/card"
)

const (
	InitHandSize = 5
	MaxHandSize  = 10
)

type Hand struct {
	c.Collection[card.Card]
}

func NewHand(card ...*card.Card) *Hand
