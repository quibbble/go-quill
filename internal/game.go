package game

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	ev "github.com/quibbble/go-quill/internal/state/hook/event"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

var (
	ErrWrongTurn = func(player uuid.UUID) error { return errors.Errorf("'%s' cannot play on other player's turn", player) }
)

type Game struct {
	*en.Engine
	*st.State
}

func NewGame(player1, player2 uuid.UUID, deck1, deck2 []string) (*Game, error) {
	// TODO add hook that ends game on 2 base death
	engine := en.NewEngine()
	state, err := st.NewState(player1, player2, deck1, deck2)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return &Game{
		Engine: engine,
		State:  state,
	}, nil
}

func (g *Game) PlayCard(player, card uuid.UUID, targets ...uuid.UUID) error {
	if player != g.State.GetTurn() {
		return ErrWrongTurn(player)
	}
	event, err := ev.NewEvent(ev.PlayCardEvent, ev.PlayCardArgs{
		Player: player,
		Choose: &ch.UUIDChoice{
			UUID: card,
		},
	})
	if err != nil {
		return errors.Wrap(err)
	}
	if err := g.Engine.Do(event, g.State, targets...); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (g *Game) MoveUnit(player, unit uuid.UUID, x, y int) error {
	if player != g.State.GetTurn() {
		return ErrWrongTurn(player)
	}
	event, err := ev.NewEvent(ev.MoveUnitEvent, ev.MoveUnitArgs{
		X: x,
		Y: y,
		Choose: &ch.UUIDChoice{
			UUID: unit,
		},
	})
	if err != nil {
		return errors.Wrap(err)
	}
	if err := g.Engine.Do(event, g.State); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (g *Game) AttackUnit(player, unit uuid.UUID, x, y int) error {
	if player != g.State.GetTurn() {
		return ErrWrongTurn(player)
	}
	event, err := ev.NewEvent(ev.AttackUnitEvent, ev.AttackUnitArgs{
		X: x,
		Y: y,
		Choose: &ch.UUIDChoice{
			UUID: unit,
		},
	})
	if err != nil {
		return errors.Wrap(err)
	}
	if err := g.Engine.Do(event, g.State); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (g *Game) SackCard(player, card uuid.UUID, option string) error {
	if player != g.State.GetTurn() {
		return ErrWrongTurn(player)
	}
	event, err := ev.NewEvent(ev.SackCardEvent, ev.SackCardArgs{
		Player: player,
		Choose: &ch.UUIDChoice{
			UUID: card,
		},
		Option: option,
	})
	if err != nil {
		return errors.Wrap(err)
	}
	if err := g.Engine.Do(event, g.State); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (g *Game) EndTurn(player uuid.UUID) error {
	if player != g.State.GetTurn() {
		return ErrWrongTurn(player)
	}
	event, err := ev.NewEvent(ev.EndTurnEvent, ev.EndTurnArgs{})
	if err != nil {
		return errors.Wrap(err)
	}
	if err := g.Engine.Do(event, g.State); err != nil {
		return errors.Wrap(err)
	}
	return nil
}
