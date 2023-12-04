package event

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	tr "github.com/quibbble/go-quill/internal/state/card/trait"
	dg "github.com/quibbble/go-quill/internal/state/damage"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	EndTurnEvent = "EndTurn"
)

type EndTurnArgs struct {
}

func EndTurnAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	// check for poison traits on player's units
	for _, col := range state.Board.XYs {
		for _, tile := range col {
			unit := tile.Unit
			if unit != nil && unit.GetPlayer() == state.GetTurn() {
				for _, poison := range unit.GetTraits(tr.PoisonTrait) {
					args := poison.GetArgs().(tr.PoisonArgs)
					event := &Event{
						uuid: state.Gen.New(st.EventUUID),
						typ:  DamageUnitEvent,
						args: &DamageUnitArgs{
							DamageType: dg.MagicDamage,
							Amount:     args.Amount,
							Choose: Choose{
								Type: ch.UUIDChoice,
								Args: ch.UUIDArgs{
									UUID: unit.GetUUID(),
								},
							},
						},
						affect: DamageUnitAffect,
					}
					if err := engine.Do(event, state); err != nil {
						return errors.Wrap(err)
					}
				}
			}
		}
	}

	state.Turn++
	player := state.GetTurn()

	state.Sacked[player] = false

	// if deck is empty then damage bases and recycle deck
	size := state.Deck[player].GetSize()
	for size <= 0 {
		state.Recycle[player]++
		events := []*Event{
			{
				uuid: state.Gen.New(st.EventUUID),
				typ:  DamageUnitsEvent,
				args: &DamageUnitsArgs{
					DamageType: dg.PureDamage,
					Amount:     state.Recycle[player],
					Choose: Choose{
						Type: ch.BasesChoice,
						Args: &ch.BasesArgs{
							Players: []uuid.UUID{player},
						},
					},
				},
				affect: DamageUnitsAffect,
			},
			{
				uuid: state.Gen.New(st.EventUUID),
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

	// refresh mana, refresh units movement, and draw a card
	events := []*Event{
		{
			uuid: state.Gen.New(st.EventUUID),
			typ:  GainManaEvent,
			args: &GainBaseManaArgs{
				Player: player,
				Amount: state.Mana[player].BaseAmount - state.Mana[player].Amount,
			},
			affect: GainManaAffect,
		},
		{
			uuid: state.Gen.New(st.EventUUID),
			typ:  RefreshMovementEvent,
			args: &RefreshMovementArgs{
				Choose: Choose{
					Type: ch.UnitsChoice,
					Args: &ch.UnitsArgs{
						Players: []uuid.UUID{player},
					},
				},
			},
			affect: RefreshMovementAffect,
		},
		{
			uuid: state.Gen.New(st.EventUUID),
			typ:  CooldownEvent,
			args: &CooldownArgs{
				Choose: Choose{
					Type: ch.UnitsChoice,
					Args: &ch.UnitsArgs{
						Players: []uuid.UUID{player},
					},
				},
			},
			affect: CooldownAffect,
		},
		{
			uuid: state.Gen.New(st.EventUUID),
			typ:  DrawCardEvent,
			args: &DrawCardArgs{
				Player: player,
			},
			affect: DrawCardAffect,
		},
	}

	for _, event := range events {
		if err := engine.Do(event, state); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
