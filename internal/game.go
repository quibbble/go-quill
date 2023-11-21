package game

import (
	"github.com/google/uuid"
	"github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/internal/state"
)

type Game struct {
	engine.Engine
	state.State
}

func (g *Game) PlayCard(player, card uuid.UUID, targets ...uuid.UUID) error

func (g *Game) MoveUnit(player, unit, tile uuid.UUID) error

func (g *Game) AttackUnit(player, attackingUnit, defendingUnit uuid.UUID) error

func (g *Game) EndTurn(player uuid.UUID) error
