package game

import (
	"math/rand"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	tr "github.com/quibbble/go-quill/internal/game/state/card/trait"
	hk "github.com/quibbble/go-quill/internal/game/state/hook"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	cn "github.com/quibbble/go-quill/internal/game/state/hook/condition"
	ev "github.com/quibbble/go-quill/internal/game/state/hook/event"
	tg "github.com/quibbble/go-quill/internal/game/state/target"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

var (
	ErrWrongTurn = func(player uuid.UUID) error { return errors.Errorf("'%s' cannot play on other player's turn", player) }
)

var build = func(gen *uuid.Gen, id string, player uuid.UUID) (st.ICard, error) {
	if len(id) == 0 {
		return nil, parse.ErrInvalidCardID
	}
	builders := &cd.Builders{
		BuildCondition: cn.NewCondition,
		BuildEvent:     ev.NewEvent,
		BuildHook:      hk.NewHook,
		BuildTargetReq: tg.NewTargetReq,
		BuildTrait:     tr.NewTrait,
		Gen:            gen,
	}
	return cd.NewCard(builders, id, player)
}

type Game struct {
	*en.Engine
	*st.State
	*uuid.Gen
}

func NewGame(seed int64, player1, player2 uuid.UUID, deck1, deck2 map[string]int) (*Game, error) {
	gen := uuid.NewGen(rand.New(rand.NewSource(seed)))
	b := func(id string, player uuid.UUID) (st.ICard, error) {
		return build(gen, id, player)
	}

	d1 := make([]string, 0)
	d2 := make([]string, 0)

	for id, count := range deck1 {
		for i := 0; i < count; i++ {
			d1 = append(d1, id)
		}
	}
	for id, count := range deck2 {
		for i := 0; i < count; i++ {
			d2 = append(d2, id)
		}
	}

	engine := en.NewEngine()
	state, err := st.NewState(seed, b, player1, player2, d1, d2)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return &Game{
		Engine: engine,
		State:  state,
		Gen:    gen,
	}, nil
}

func (g *Game) PlayCard(player, card uuid.UUID, targets ...uuid.UUID) error {
	if player != g.State.GetTurn() {
		return ErrWrongTurn(player)
	}
	event, err := ev.NewEvent(g.Gen.New(st.EventUUID), ev.PlayCardEvent, ev.PlayCardArgs{
		Player: player,
		Choose: ch.RawChoose{
			Type: ch.UUIDChoice,
			Args: &ch.UUIDArgs{
				UUID: card,
			},
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
	event, err := ev.NewEvent(g.Gen.New(st.EventUUID), ev.MoveUnitEvent, ev.MoveUnitArgs{
		X: x,
		Y: y,
		Choose: ch.RawChoose{
			Type: ch.UUIDChoice,
			Args: &ch.UUIDArgs{
				UUID: unit,
			},
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
	event, err := ev.NewEvent(g.Gen.New(st.EventUUID), ev.AttackUnitEvent, ev.AttackUnitArgs{
		X: x,
		Y: y,
		Choose: ch.RawChoose{
			Type: ch.UUIDChoice,
			Args: &ch.UUIDArgs{
				UUID: unit,
			},
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
	event, err := ev.NewEvent(g.Gen.New(st.EventUUID), ev.SackCardEvent, ev.SackCardArgs{
		Player: player,
		Choose: ch.RawChoose{
			Type: ch.UUIDChoice,
			Args: &ch.UUIDArgs{
				UUID: card,
			},
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
	event, err := ev.NewEvent(g.Gen.New(st.EventUUID), ev.EndTurnEvent, ev.EndTurnArgs{})
	if err != nil {
		return errors.Wrap(err)
	}
	if err := g.Engine.Do(event, g.State); err != nil {
		return errors.Wrap(err)
	}
	return nil
}
