package event

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	dg "github.com/quibbble/go-quill/internal/state/damage"
	"github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	EndTurnEvent = "end_turn"
)

type EndTurnArgs struct {
}

func EndTurnAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	state.Turn++
	player := state.GetTurn()

	// if deck is empty then damage bases and recycle deck
	size := state.Deck[player].GetSize()
	for size <= 0 {
		state.Recycle[player]++
		events := []*Event{
			{
				uuid: uuid.New(st.EventUUID),
				typ:  DamageUnitsEvent,
				args: &DamageUnitsArgs{
					DamageType: dg.DamageTypePure,
					Amount:     state.Recycle[player],
					Choose: &choose.BasesChoice{
						Player: player,
					},
				},
				affect: DamageUnitsAffect,
			},
			{
				uuid: uuid.New(st.EventUUID),
				typ:  RecycleDeckEvent,
				args: &RecycleDeckArgs{
					Player: player,
				},
				affect: RecycleDeckAffect,
			},
		}
		for _, event := range events {
			if err := engine.Do(event, state); err != nil {
				return errors.Wrap(err)
			}
		}
		size = state.Deck[player].GetSize()
	}

	// draw a card
	if err := engine.Do(&Event{
		uuid: uuid.New(st.EventUUID),
		typ:  DrawCardEvent,
		args: &DrawCardArgs{
			Player: player,
		},
		affect: DrawCardAffect,
	}, state); err != nil {
		return errors.Wrap(err)
	}
	return nil
}
