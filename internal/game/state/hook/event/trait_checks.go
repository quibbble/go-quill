package event

import (
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
					args := trait.GetArgs().(tr.FriendsArgs)
					before := args.Current
					if before == nil {
						before = make([]uuid.UUID, 0)
					}
					choose, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), args.ChooseUnits.Type, args.ChooseUnits.Args)
					if err != nil {
						return errors.Wrap(err)
					}
					owned, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), ch.OwnedUnitsChoice, &ch.OwnedUnitsArgs{
						ChoosePlayer: parse.Choose{
							Type: ch.UUIDChoice,
							Args: &ch.UUIDArgs{
								UUID: tile.Unit.GetPlayer(),
							},
						},
					})
					if err != nil {
						return errors.Wrap(err)
					}
					after, err := ch.NewChoices(choose, owned).Retrieve(engine, state, tile.Unit.GetUUID())
					if err != nil {
						return errors.Wrap(err)
					}
					args.Current = after
					trait.SetArgs(args)

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
					args := trait.GetArgs().(tr.EnemiesArgs)
					before := args.Current
					if before == nil {
						before = make([]uuid.UUID, 0)
					}
					choose, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), args.ChooseUnits.Type, args.ChooseUnits.Args)
					if err != nil {
						return errors.Wrap(err)
					}
					owned, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), ch.OwnedUnitsChoice, &ch.OwnedUnitsArgs{
						ChoosePlayer: parse.Choose{
							Type: ch.UUIDChoice,
							Args: &ch.UUIDArgs{
								UUID: state.GetOpponent(tile.Unit.GetPlayer()),
							},
						},
					})
					if err != nil {
						return errors.Wrap(err)
					}
					after, err := ch.NewChoices(choose, owned).Retrieve(engine, state, tile.Unit.GetUUID())
					if err != nil {
						return errors.Wrap(err)
					}
					args.Current = after
					trait.SetArgs(args)

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
				event, err := NewEvent(state.Gen.New(st.EventUUID), RemoveTraitFromCard, &RemoveTraitFromCardArgs{
					Trait: t.GetUUID(),
					ChooseCard: parse.Choose{
						Type: ch.UUIDChoice,
						Args: &ch.UUIDArgs{
							UUID: u,
						},
					},
				})
				if err != nil {
					return errors.Wrap(err)
				}
				if err := engine.Do(event, state); err != nil {
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
		event, err := NewEvent(state.Gen.New(st.EventUUID), AddTraitToCard, &AddTraitToCardArgs{
			Trait: parse.Trait{
				Type: trait.GetType(),
				Args: trait.GetArgs(),
			},
			ChooseCard: parse.Choose{
				Type: ch.UUIDChoice,
				Args: &ch.UUIDArgs{
					UUID: u,
				},
			},
		})
		if err != nil {
			return errors.Wrap(err)
		}
		if err := engine.Do(event, state); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
