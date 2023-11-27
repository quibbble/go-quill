package state

import (
	"github.com/quibbble/go-quill/internal/state/card"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type State struct {
	Turn   int
	Teams  []uuid.UUID
	Winner *uuid.UUID

	Board *Board

	Deck    map[uuid.UUID]*Deck
	Discard map[uuid.UUID]*Deck
	Trash   map[uuid.UUID]*Deck
	Hand    map[uuid.UUID]*Hand
	Mana    map[uuid.UUID]*Mana
	Recycle map[uuid.UUID]int
}

func NewState(player1, player2 uuid.UUID, deck1, deck2 []string) (*State, error) {
	board, err := NewBoard(player1, player2)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	d1, err := NewDeck(deck1...)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	d2, err := NewDeck(deck2...)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	hand1 := make([]card.ICard, 0)
	hand2 := make([]card.ICard, 0)
	for i := 0; i < InitHandSize; i++ {
		card1, err := d1.Draw()
		if err != nil {
			return nil, errors.Wrap(err)
		}
		card2, err := d2.Draw()
		if err != nil {
			return nil, errors.Wrap(err)
		}
		hand1 = append(hand1, *card1)
		hand2 = append(hand2, *card2)
	}

	return &State{
		Turn:   0,
		Teams:  []uuid.UUID{player1, player2},
		Winner: nil,

		Board:   board,
		Deck:    map[uuid.UUID]*Deck{player1: d1, player2: d2},
		Discard: map[uuid.UUID]*Deck{player1: NewEmptyDeck(), player2: NewEmptyDeck()},
		Trash:   map[uuid.UUID]*Deck{player1: NewEmptyDeck(), player2: NewEmptyDeck()},
		Hand:    map[uuid.UUID]*Hand{player1: NewHand(hand1...), player2: NewHand(hand2...)},
		Mana:    map[uuid.UUID]*Mana{player1: NewMana(), player2: NewMana()},
		Recycle: map[uuid.UUID]int{player1: 0, player2: 0},
	}, nil
}

func (s *State) GetTurn() uuid.UUID {
	return s.Teams[s.Turn%len(s.Teams)]
}

func (s *State) GetOpponent(player uuid.UUID) uuid.UUID {
	if s.Teams[0] == player {
		return s.Teams[1]
	}
	return s.Teams[0]
}

func (s *State) GameOver() bool {
	return s.Winner != nil
}
