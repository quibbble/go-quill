package event

import (
	"context"

	"github.com/mitchellh/mapstructure"
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
				for _, poison := range unit.GetTraits(tr.PoisonTrait) {
					var args tr.PoisonArgs
					if err := mapstructure.Decode(poison.GetArgs(), &args); err != nil {
						return errors.Wrap(err)
					}
					event := &Event{
						uuid: state.Gen.New(en.EventUUID),
						typ:  DamageUnitEvent,
						args: &DamageUnitArgs{
							DamageType: dg.MagicDamage,
							Amount:     args.Amount,
							ChooseUnit: parse.Choose{
								Type: ch.UUIDChoice,
								Args: ch.UUIDArgs{
									UUID: unit.GetUUID(),
								},
							},
						},
						affect: DamageUnitAffect,
					}
					if err := engine.Do(context.Background(), event, state); err != nil {
						return errors.Wrap(err)
					}
				}
			}
		}
	}

	// update units movement and cooldown
	events := []*Event{
		{
			uuid: state.Gen.New(en.EventUUID),
			typ:  RefreshMovementEvent,
			args: &RefreshMovementArgs{
				ChooseUnits: parse.Choose{
					Type: ch.CompositeChoice,
					Args: &ch.CompositeArgs{
						SetFunction: ch.SetIntersect,
						Choices: []parse.Choose{
							{
								Type: ch.OwnedUnitsChoice,
								Args: &ch.OwnedUnitsArgs{
									ChoosePlayer: parse.Choose{
										Type: ch.CurrentPlayerChoice,
										Args: &ch.CurrentPlayerArgs{},
									},
								},
							},
							{
								Type: ch.UnitsChoice,
								Args: &ch.UnitsArgs{
									Types: []string{
										cd.CreatureUnit,
									},
								},
							},
						},
					},
				},
			},
			affect: RefreshMovementAffect,
		},
		{
			uuid: state.Gen.New(en.EventUUID),
			typ:  CooldownEvent,
			args: &CooldownArgs{
				ChooseUnits: parse.Choose{
					Type: ch.CompositeChoice,
					Args: &ch.CompositeArgs{
						SetFunction: ch.SetIntersect,
						Choices: []parse.Choose{
							{
								Type: ch.OwnedUnitsChoice,
								Args: &ch.OwnedUnitsArgs{
									ChoosePlayer: parse.Choose{
										Type: ch.CurrentPlayerChoice,
										Args: &ch.CurrentPlayerArgs{},
									},
								},
							},
							{
								Type: ch.UnitsChoice,
								Args: &ch.UnitsArgs{
									Types: []string{
										cd.CreatureUnit,
										cd.StructureUnit,
									},
								},
							},
						},
					},
				},
			},
			affect: CooldownAffect,
		},
	}
	for _, event := range events {
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
		events := []*Event{
			{
				uuid: state.Gen.New(en.EventUUID),
				typ:  DamageUnitsEvent,
				args: &DamageUnitsArgs{
					DamageType: dg.PureDamage,
					Amount:     state.Recycle[player],
					ChooseUnits: parse.Choose{
						Type: ch.CompositeChoice,
						Args: &ch.CompositeArgs{
							SetFunction: ch.SetIntersect,
							Choices: []parse.Choose{
								{
									Type: ch.OwnedUnitsChoice,
									Args: &ch.OwnedUnitsArgs{
										ChoosePlayer: parse.Choose{
											Type: ch.CurrentPlayerChoice,
											Args: &ch.CurrentPlayerArgs{},
										},
									},
								},
								{
									Type: ch.UnitsChoice,
									Args: &ch.UnitsArgs{
										Types: []string{
											cd.BaseUnit,
										},
									},
								},
							},
						},
					},
				},
				affect: DamageUnitsAffect,
			},
			{
				uuid: state.Gen.New(en.EventUUID),
				typ:  RecycleDeckEvent,
				args: &RecycleDeckArgs{
					ChoosePlayer: parse.Choose{
						Type: ch.UUIDChoice,
						Args: &ch.UUIDArgs{
							UUID: player,
						},
					},
				},
				affect: RecycleDeckAffect,
			},
		}
		for _, event := range events {
			if err := engine.Do(context.Background(), event, state); err != nil {
				return errors.Wrap(err)
			}
		}
		size = state.Deck[player].GetSize()
	}

	// refresh mana and draw a card
	events = []*Event{
		{
			uuid: state.Gen.New(en.EventUUID),
			typ:  GainManaEvent,
			args: &GainManaArgs{
				ChoosePlayer: parse.Choose{
					Type: ch.CurrentPlayerChoice,
					Args: &ch.CurrentPlayerArgs{},
				},
				Amount: state.Mana[player].BaseAmount - state.Mana[player].Amount,
			},
			affect: GainManaAffect,
		},
		{
			uuid: state.Gen.New(en.EventUUID),
			typ:  DrawCardEvent,
			args: &DrawCardArgs{
				ChoosePlayer: parse.Choose{
					Type: ch.UUIDChoice,
					Args: &ch.UUIDArgs{
						UUID: player,
					},
				},
			},
			affect: DrawCardAffect,
		},
	}

	for _, event := range events {
		if err := engine.Do(context.Background(), event, state); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
