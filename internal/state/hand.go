package state

import (
	c "github.com/quibbble/go-boardgame/pkg/collection"
	"github.com/quibbble/go-quill/internal/state/card"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	InitHandSize = 5
	MaxHandSize  = 10
)

type Hand struct {
	c.Collection[card.ICard]
}

func NewHand(card ...card.ICard) *Hand

func (h *Hand) GetCard(card uuid.UUID) (card.ICard, error) {
	for _, it := range h.GetItems() {
		if it.GetUUID() == card {
			return it, nil
		}
	}
	return nil, ErrNotFound(card)
}

func (h *Hand) RemoveCard(card uuid.UUID) error {
	for i, it := range h.GetItems() {
		if it.GetUUID() == card {
			if err := h.Remove(i); err != nil {
				return errors.Wrap(err)
			}
			return nil
		}
	}
	return ErrNotFound(card)
}
