package game

import (
	"github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type Game struct {
	engine.Engine
	state.State
}

func (g *Game) PlayCard(player, card uuid.UUID, targets ...uuid.UUID) error {
	// TODO get card from game
	// TODO remove mana cost from player mana pool
	// TODO remove card from hand

	// TODO do play card event

	// TODO if item card add to unit
	// TODO if spell card send to discard
	// TODO if unit card place on board

	return nil
}

func (g *Game) MoveUnit(player, unit, tile uuid.UUID) error

func (g *Game) AttackUnit(player, attackingUnit, defendingUnit uuid.UUID) error

func (g *Game) EndTurn(player uuid.UUID) error
