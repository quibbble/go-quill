package event

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	tr "github.com/quibbble/go-quill/internal/game/state/card/trait"
	dg "github.com/quibbble/go-quill/internal/game/state/damage"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	EndTurnEvent = "EndTurn"
)

type EndTurnArgs struct{}

func EndTurnAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	// poison trait check on player's units
	for _, col := range state.Board.XYs {
		for _, tile := range col {
			unit := tile.Unit
			if unit != nil && unit.GetPlayer() == state.GetTurn() {
				for _, trait := range unit.GetTraits(tr.PoisonTrait) {
					args := trait.GetArgs().(*tr.PoisonArgs)
					event, err := NewEvent(state.Gen.New(en.EventUUID), DamageUnitEvent, DamageUnitArgs{
						DamageType: dg.MagicDamage,
						Amount:     args.Amount,
						ChooseUnit: parse.Choose{
							Type: ch.UUIDChoice,
							Args: ch.UUIDArgs{
								UUID: unit.GetUUID(),
							},
						},
					})
					if err != nil {
						return errors.Wrap(err)
					}
					if err := engine.Do(context.Background(), event, state); err != nil {
						return errors.Wrap(err)
					}
				}
			}
		}
	}

	// update units movement and cooldown
	event1, err := NewEvent(state.Gen.New(en.EventUUID), RefreshMovementEvent, RefreshMovementArgs{
		ChooseUnits: parse.Choose{
			Type: ch.CompositeChoice,
			Args: ch.CompositeArgs{
				SetFunction: ch.SetIntersect,
				ChooseChain: []parse.Choose{
					{
						Type: ch.OwnedUnitsChoice,
						Args: ch.OwnedUnitsArgs{
							ChoosePlayer: parse.Choose{
								Type: ch.CurrentPlayerChoice,
								Args: ch.CurrentPlayerArgs{},
							},
						},
					},
					{
						Type: ch.UnitsChoice,
						Args: ch.UnitsArgs{
							Types: []string{
								cd.CreatureUnit,
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return errors.Wrap(err)
	}
	event2, err := NewEvent(state.Gen.New(en.EventUUID), CooldownEvent, CooldownArgs{
		ChooseUnits: parse.Choose{
			Type: ch.CompositeChoice,
			Args: ch.CompositeArgs{
				SetFunction: ch.SetIntersect,
				ChooseChain: []parse.Choose{
					{
						Type: ch.OwnedUnitsChoice,
						Args: ch.OwnedUnitsArgs{
							ChoosePlayer: parse.Choose{
								Type: ch.CurrentPlayerChoice,
								Args: ch.CurrentPlayerArgs{},
							},
						},
					},
					{
						Type: ch.UnitsChoice,
						Args: ch.UnitsArgs{
							Types: []string{
								cd.CreatureUnit,
								cd.StructureUnit,
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return errors.Wrap(err)
	}
	for _, event := range []en.IEvent{event1, event2} {
		if err := engine.Do(context.Background(), event, state); err != nil {
			return errors.Wrap(err)
		}
	}

	state.Turn++
	player := state.GetTurn()

	state.Sacked[player] = false

	// if deck is empty then damage bases and recycle deck
	size := state.Deck[player].GetSize()
	for size <= 0 {
		state.Recycle[player]++
		event1, err := NewEvent(state.Gen.New(en.EventUUID), DamageUnitsEvent, DamageUnitsArgs{
			DamageType: dg.PureDamage,
			Amount:     state.Recycle[player],
			ChooseUnits: parse.Choose{
				Type: ch.CompositeChoice,
				Args: ch.CompositeArgs{
					SetFunction: ch.SetIntersect,
					ChooseChain: []parse.Choose{
						{
							Type: ch.OwnedUnitsChoice,
							Args: ch.OwnedUnitsArgs{
								ChoosePlayer: parse.Choose{
									Type: ch.CurrentPlayerChoice,
									Args: ch.CurrentPlayerArgs{},
								},
							},
						},
						{
							Type: ch.UnitsChoice,
							Args: ch.UnitsArgs{
								Types: []string{
									cd.BaseUnit,
								},
							},
						},
					},
				},
			},
		})
		if err != nil {
			return errors.Wrap(err)
		}
		event2, err := NewEvent(state.Gen.New(en.EventUUID), RecycleDeckEvent, RecycleDeckArgs{
			ChoosePlayer: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: player,
				},
			},
		})
		if err != nil {
			return errors.Wrap(err)
		}
		for _, event := range []en.IEvent{event1, event2} {
			if err := engine.Do(context.Background(), event, state); err != nil {
				return errors.Wrap(err)
			}
		}
		size = state.Deck[player].GetSize()
	}

	// refresh mana and draw a card
	event1, err = NewEvent(state.Gen.New(en.EventUUID), GainManaEvent, GainManaArgs{
		ChoosePlayer: parse.Choose{
			Type: ch.CurrentPlayerChoice,
			Args: ch.CurrentPlayerArgs{},
		},
		Amount: state.Mana[player].BaseAmount - state.Mana[player].Amount,
	})
	if err != nil {
		return errors.Wrap(err)
	}
	event2, err = NewEvent(state.Gen.New(en.EventUUID), DrawCardEvent, DrawCardArgs{
		ChoosePlayer: parse.Choose{
			Type: ch.UUIDChoice,
			Args: ch.UUIDArgs{
				UUID: player,
			},
		},
	})
	if err != nil {
		return errors.Wrap(err)
	}
	for _, event := range []en.IEvent{event1, event2} {
		if err := engine.Do(context.Background(), event, state); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
