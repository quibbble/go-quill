package state

import (
	"github.com/quibbble/go-quill/internal/state/card"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type State struct {
	Turn   int
	Teams  []uuid.UUID
	Winner *uuid.UUID

	Board *Board

	Decks map[uuid.UUID]*Deck
	Hands map[uuid.UUID]*Hand
	Mana  map[uuid.UUID]*Mana
}

func NewState(player1, player2 uuid.UUID, deck1, deck2 []string) (*State, error) {
	board, err := NewBoard(player1, player2)
	if err != nil {
		return nil, err
	}

	d1, err := NewDeck(deck1)
	if err != nil {
		return nil, err
	}
	d2, err := NewDeck(deck2)
	if err != nil {
		return nil, err
	}

	hand1 := make([]*card.Card, 0)
	hand2 := make([]*card.Card, 0)
	for i := 0; i < InitHandSize; i++ {
		card1, err := d1.Draw()
		if err != nil {
			return nil, err
		}
		card2, err := d2.Draw()
		if err != nil {
			return nil, err
		}
		hand1 = append(hand1, card1)
		hand2 = append(hand2, card2)
	}

	return &State{
		Turn:   0,
		Teams:  []uuid.UUID{player1, player2},
		Winner: nil,

		Board: board,
		Decks: map[uuid.UUID]*Deck{player1: d1, player2: d2},
		Hands: map[uuid.UUID]*Hand{player1: NewHand(hand1...), player2: NewHand(hand2...)},
		Mana:  map[uuid.UUID]*Mana{player1: NewMana(), player2: NewMana()},
	}, nil
}

func (s *State) GetTurn() uuid.UUID {
	return s.Teams[s.Turn%len(s.Teams)]
}
