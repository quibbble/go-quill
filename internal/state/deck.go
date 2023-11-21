package state

import (
	c "github.com/quibbble/go-boardgame/pkg/collection"
	"github.com/quibbble/go-quill/internal/state/card"
)

type Deck struct {
	c.Collection[card.Card]
}

func NewDeck([]string) (*Deck, error)
