package event

import (
	"context"
	"reflect"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	tr "github.com/quibbble/go-quill/internal/game/state/card/trait"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

func FriendsTraitCheck(engine *en.Engine, state *st.State) error {
	// checks if any units need to add/remove their trait through friend trait
	for _, col := range state.Board.XYs {
		for _, tile := range col {
			if tile.Unit != nil {
				for _, trait := range tile.Unit.GetTraits(tr.FriendsTrait) {
					args := trait.GetArgs().(*tr.FriendsArgs)
					before := args.Current
					if before == nil {
						before = make([]uuid.UUID, 0)
					}
					choose1, err := ch.NewChoose(state.Gen.New(en.ChooseUUID), args.ChooseUnits.Type, args.ChooseUnits.Args)
					if err != nil {
						return errors.Wrap(err)
					}
					choose2, err := ch.NewChoose(state.Gen.New(en.ChooseUUID), ch.OwnedUnitsChoice, &ch.OwnedUnitsArgs{
						ChoosePlayer: parse.Choose{
							Type: ch.UUIDChoice,
							Args: ch.UUIDArgs{
								UUID: tile.Unit.GetPlayer(),
							},
						},
					})
					if err != nil {
						return errors.Wrap(err)
					}
					after, err := ch.NewChooseChain(ch.SetIntersect, choose1, choose2).Retrieve(context.WithValue(context.Background(), en.CardCtx, tile.Unit.GetUUID()), engine, state)
					if err != nil {
						return errors.Wrap(err)
					}
					args.Current = after

					t, err := tr.NewTrait(uuid.Nil, args.Trait.Type, args.Trait.Args)
					if err != nil {
						return errors.Wrap(err)
					}

					if err := updateUnits(engine, state, before, after, t); err != nil {
						return errors.Wrap(err)
					}
				}
			}
		}
	}
	return nil
}

func EnemiesTraitCheck(engine *en.Engine, state *st.State) error {
	// checks if any units need to add/remove their trait through friend trait
	for _, col := range state.Board.XYs {
		for _, tile := range col {
			if tile.Unit != nil {
				for _, trait := range tile.Unit.GetTraits(tr.EnemiesTrait) {
					args := trait.GetArgs().(*tr.EnemiesArgs)
					before := args.Current
					if before == nil {
						before = make([]uuid.UUID, 0)
					}
					choose1, err := ch.NewChoose(state.Gen.New(en.ChooseUUID), args.ChooseUnits.Type, args.ChooseUnits.Args)
					if err != nil {
						return errors.Wrap(err)
					}
					choose2, err := ch.NewChoose(state.Gen.New(en.ChooseUUID), ch.OwnedUnitsChoice, &ch.OwnedUnitsArgs{
						ChoosePlayer: parse.Choose{
							Type: ch.UUIDChoice,
							Args: ch.UUIDArgs{
								UUID: state.GetOpponent(tile.Unit.GetPlayer()),
							},
						},
					})
					if err != nil {
						return errors.Wrap(err)
					}
					after, err := ch.NewChooseChain(ch.SetIntersect, choose1, choose2).Retrieve(context.WithValue(context.Background(), en.CardCtx, tile.Unit.GetUUID()), engine, state)
					if err != nil {
						return errors.Wrap(err)
					}
					args.Current = after

					t, err := tr.NewTrait(uuid.Nil, args.Trait.Type, args.Trait.Args)
					if err != nil {
						return errors.Wrap(err)
					}

					if err := updateUnits(engine, state, before, after, t); err != nil {
						return errors.Wrap(err)
					}
				}
			}
		}
	}
	return nil
}

func updateUnits(engine *en.Engine, state *st.State, before, after []uuid.UUID, trait st.ITrait) error {
	remove := uuid.Diff(before, after)
	for _, u := range remove {
		x, y, err := state.Board.GetUnitXY(u)
		if err != nil {
			return errors.Wrap(err)
		}
		found := false
		unit := state.Board.XYs[x][y].Unit.(*cd.UnitCard)
		for _, t := range unit.GetTraits(trait.GetType()) {
			if reflect.DeepEqual(t.GetArgs(), trait.GetArgs()) {
				event, err := NewEvent(state.Gen.New(en.EventUUID), RemoveTraitFromCard, &RemoveTraitFromCardArgs{
					ChooseTrait: parse.Choose{
						Type: ch.UUIDChoice,
						Args: ch.UUIDArgs{
							UUID: t.GetUUID(),
						},
					},
					ChooseCard: parse.Choose{
						Type: ch.UUIDChoice,
						Args: ch.UUIDArgs{
							UUID: u,
						},
					},
				})
				if err != nil {
					return errors.Wrap(err)
				}
				if err := engine.Do(context.Background(), event, state); err != nil {
					return errors.Wrap(err)
				}
				found = true
				break
			}
		}
		if !found {
			return errors.Errorf("failed to find '%s' trait for '%s'", trait.GetType(), u)
		}
	}
	add := uuid.Diff(after, before)
	for _, u := range add {
		event, err := NewEvent(state.Gen.New(en.EventUUID), AddTraitToCard, &AddTraitToCardArgs{
			Trait: parse.Trait{
				Type: trait.GetType(),
				Args: trait.GetArgs(),
			},
			ChooseCard: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: u,
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
	return nil
}
