package state

import (
	"math/rand"

	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type State struct {
	Seed int64
	Rand *rand.Rand

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

	BuildCard
}

type BuildCard func(id string, player uuid.UUID) (ICard, error)

func NewState(seed int64, build BuildCard, player1, player2 uuid.UUID, deck1, deck2 []string) (*State, error) {
	board, err := NewBoard(build, player1, player2)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	d1, err := NewDeck(seed, build, player1, deck1...)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	d2, err := NewDeck(seed, build, player2, deck2...)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	hand1 := make([]ICard, 0)
	hand2 := make([]ICard, 0)
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
		Seed: seed,
		Rand: rand.New(rand.NewSource(seed)),

		Turn:   0,
		Teams:  []uuid.UUID{player1, player2},
		Winner: nil,

		Board:   board,
		Deck:    map[uuid.UUID]*Deck{player1: d1, player2: d2},
		Discard: map[uuid.UUID]*Deck{player1: NewEmptyDeck(seed), player2: NewEmptyDeck(seed)},
		Trash:   map[uuid.UUID]*Deck{player1: NewEmptyDeck(seed), player2: NewEmptyDeck(seed)},
		Hand:    map[uuid.UUID]*Hand{player1: NewHand(seed, hand1...), player2: NewHand(seed, hand2...)},
		Mana:    map[uuid.UUID]*Mana{player1: NewMana(), player2: NewMana()},
		Recycle: map[uuid.UUID]int{player1: 0, player2: 0},

		BuildCard: build,
	}, nil
}

func (s *State) GetTurn() uuid.UUID {
	return s.Teams[s.Turn%len(s.Teams)]
}

func (s *State) GetCard(uuid uuid.UUID) ICard {
	for _, hand := range s.Hand {
		for _, card := range hand.GetItems() {
			if card.GetUUID() == uuid {
				return card
			}
		}
	}
	for _, tile := range s.Board.UUIDs {
		if tile.Unit != nil && tile.Unit.GetUUID() == uuid {
			return tile.Unit
		}
	}
	for _, deck := range s.Deck {
		for _, card := range deck.GetItems() {
			if card.GetUUID() == uuid {
				return card
			}
		}
	}
	for _, discard := range s.Discard {
		for _, card := range discard.GetItems() {
			if card.GetUUID() == uuid {
				return card
			}
		}
	}
	for _, trash := range s.Trash {
		for _, card := range trash.GetItems() {
			if card.GetUUID() == uuid {
				return card
			}
		}
	}
	return nil
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
