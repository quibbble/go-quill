package event

import (
	"reflect"

	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	cd "github.com/quibbble/go-quill/internal/state/card"
	tr "github.com/quibbble/go-quill/internal/state/card/trait"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
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
					choose, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), args.Choose.Type, args.Choose.Args)
					if err != nil {
						return errors.Wrap(err)
					}
					owned, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), ch.OwnedChoice, &ch.OwnedArgs{
						Player: tile.Unit.GetPlayer(),
					})
					if err != nil {
						return errors.Wrap(err)
					}
					uuids, err := ch.NewChoices(choose, owned).Retrieve(engine, state)
					if err != nil {
						return errors.Wrap(err)
					}
					after := uuids
					if err := updateUnits(engine, state, before, after, args.Trait); err != nil {
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
					choose, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), args.Choose.Type, args.Choose.Args)
					if err != nil {
						return errors.Wrap(err)
					}
					owned, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), ch.OwnedChoice, &ch.OwnedArgs{
						Player: state.GetOpponent(tile.Unit.GetPlayer()),
					})
					if err != nil {
						return errors.Wrap(err)
					}
					uuids, err := ch.NewChoices(choose, owned).Retrieve(engine, state)
					if err != nil {
						return errors.Wrap(err)
					}
					after := uuids
					if err := updateUnits(engine, state, before, after, args.Trait); err != nil {
						return errors.Wrap(err)
					}
				}
			}
		}
	}
	return nil
}

func updateUnits(engine *en.Engine, state *st.State, before, after []uuid.UUID, trait Trait) error {
	remove := uuid.Diff(before, after)
	for _, u := range remove {
		x, y, err := state.Board.GetUnitXY(u)
		if err != nil {
			return errors.Wrap(err)
		}
		found := false
		unit := state.Board.XYs[x][y].Unit.(*cd.UnitCard)
		for _, t := range unit.GetTraits(trait.Type) {
			if reflect.DeepEqual(t.GetArgs(), trait) {
				event, err := NewEvent(state.Gen.New(st.EventUUID), RemoveTraitFromCard, &RemoveTraitFromCardArgs{
					Trait: t.GetUUID(),
					Choose: Choose{
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
			return errors.Errorf("failed to find '%s' trait for '%s'", trait.Type, u)
		}
	}
	add := uuid.Diff(after, before)
	for _, u := range add {
		event, err := NewEvent(state.Gen.New(st.EventUUID), AddTraitToCard, &AddTraitToCardArgs{
			Trait: trait,
			Choose: Choose{
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
