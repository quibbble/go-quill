package event

import (
	"github.com/mitchellh/mapstructure"
	"github.com/quibbble/go-quill/cards"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	cd "github.com/quibbble/go-quill/internal/state/card"
	tr "github.com/quibbble/go-quill/internal/state/card/trait"
	dm "github.com/quibbble/go-quill/internal/state/damage"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	AttackUnitEvent = "AttackUnit"
)

type AttackUnitArgs struct {
	X, Y   int
	Choose Choose
}

func AttackUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a AttackUnitArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	choose, err := ch.NewChoice(a.Choose.Type, a.Choose.Args)
	if err != nil {
		return errors.Wrap(err)
	}
	choices, err := choose.Retrieve(engine, state, targets...)
	if err != nil {
		return errors.Wrap(err)
	}
	if len(choices) != 1 {
		return errors.ErrInvalidSliceLength
	}
	if choices[0].Type() != st.UnitUUID {
		return st.ErrInvalidUUIDType(choices[0], st.UnitUUID)
	}
	x, y, err := state.Board.GetUnitXY(choices[0])
	if err != nil {
		return errors.Wrap(err)
	}
	attacker := state.Board.XYs[x][y].Unit.(*cd.UnitCard)
	if attacker.Cooldown != 0 {
		return errors.Errorf("unit '%s' cannot attack due to cooldown", attacker.UUID)
	}
	if attacker.Range <= 0 && !attacker.CheckCodex(x, y, a.X, a.Y) {
		return errors.Errorf("unit '%s' cannot attack due to failed codex check", attacker.UUID)
	}
	if attacker.Range > 0 && !attacker.CheckRange(x, y, a.X, a.Y) {
		return errors.Errorf("unit '%s' cannot attack due to failed range check", attacker.UUID)
	}
	if state.Board.XYs[a.X][a.Y].Unit == nil {
		return errors.Errorf("unit '%s' cannot attack at (%d,%d) as no unit exists", attacker.UUID, a.X, a.Y)
	}
	defender := state.Board.XYs[a.X][a.Y].Unit.(*cd.UnitCard)

	// thief trait check
	if len(attacker.GetTraits(tr.ThiefTrait)) > 0 && len(defender.Items) > 0 {
		item := defender.Items[state.Rand.Intn(len(defender.Items))]
		// set thief player as owner and add to that players hand to allow adding to thief
		item.Player = attacker.Player
		state.Hand[attacker.Player].Add(item)
		events := []*Event{
			{
				uuid: uuid.New(st.EventUUID),
				typ:  RemoveItemFromUnitEvent,
				args: &RemoveItemFromUnitArgs{
					ChooseItem: Choose{
						Type: ch.UUIDChoice,
						Args: &ch.UUIDArgs{
							UUID: item.UUID,
						},
					},
					ChooseUnit: Choose{
						Type: ch.UUIDChoice,
						Args: &ch.UUIDArgs{
							UUID: defender.UUID,
						},
					},
				},
				affect: RemoveItemFromUnitAffect,
			},
			{
				uuid: uuid.New(st.EventUUID),
				typ:  AddItemToUnitEvent,
				args: &AddItemToUnitArgs{
					ChooseItem: Choose{
						Type: ch.UUIDChoice,
						Args: &ch.UUIDArgs{
							UUID: item.UUID,
						},
					},
					ChooseUnit: Choose{
						Type: ch.UUIDChoice,
						Args: &ch.UUIDArgs{
							UUID: attacker.UUID,
						},
					},
				},
				affect: AddItemToUnitAffect,
			},
		}
		for _, event := range events {
			if err := engine.Do(event, state); err != nil {
				return errors.Wrap(err)
			}
		}
		return nil
	}

	attackerDamage, defenderDamage, err := dm.Battle(state, attacker, defender)
	if err != nil {
		return errors.Wrap(err)
	}
	if attackerDamage > 0 {
		var event *Event
		event = &Event{
			uuid: uuid.New(st.EventUUID),
			typ:  DamageUnitEvent,
			args: &DamageUnitArgs{
				DamageType: attacker.DamageType,
				Amount:     attackerDamage,
				Choose: Choose{
					Type: ch.UUIDChoice,
					Args: &ch.UUIDArgs{
						UUID: defender.UUID,
					},
				},
			},
			affect: DamageUnitAffect,
		}

		// execute trait check
		if len(attacker.GetTraits(tr.ExecuteTrait)) > 0 &&
			defender.Health < defender.GetInit().(*cards.UnitCard).Health {
			event = &Event{
				uuid: uuid.New(st.EventUUID),
				typ:  KillUnitEvent,
				args: &KillUnitArgs{
					Choose: Choose{
						Type: ch.UUIDChoice,
						Args: &ch.UUIDArgs{
							UUID: defender.UUID,
						},
					},
				},
				affect: KillUnitAffect,
			}
		}

		if err := engine.Do(event, state); err != nil {
			return errors.Wrap(err)
		}

		// pillage trait check
		if defender.GetInit().(*cards.Card).ID == "U0001" {
			pillages := attacker.GetTraits(tr.PillageTrait)
			for _, pillage := range pillages {
				args := pillage.GetArgs().(*tr.PillageArgs)
				event, err := NewEvent(args.Event.Type, args.Event.Args)
				if err != nil {
					return errors.Wrap(err)
				}
				if err := engine.Do(event, state); err != nil {
					return errors.Wrap(err)
				}
			}
		}

		// gift trait check
		if defender.Type == cd.CreatureUnit {
			gifts := attacker.GetTraits(tr.GiftTrait)
			for _, gift := range gifts {
				args := gift.GetArgs().(*tr.GiftArgs)
				event, err := NewEvent(AddTraitToCard, &AddTraitToCardArgs{
					Trait: args.Trait,
					Choose: Choose{
						Type: ch.UUIDChoice,
						Args: &ch.UUIDArgs{
							UUID: defender.UUID,
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
		}
	}

	if defenderDamage > 0 && attacker.Range <= 0 {
		event := &Event{
			uuid: uuid.New(st.EventUUID),
			typ:  DamageUnitEvent,
			args: &DamageUnitArgs{
				DamageType: defender.DamageType,
				Amount:     defenderDamage,
				Choose: Choose{
					Type: ch.UUIDChoice,
					Args: &ch.UUIDArgs{
						UUID: attacker.UUID,
					},
				},
			},
			affect: DamageUnitAffect,
		}
		if err := engine.Do(event, state); err != nil {
			return errors.Wrap(err)
		}
	}

	// if attacker still on board then reset cooldown
	if _, _, err := state.Board.GetUnitXY(attacker.UUID); err == nil {
		attacker.Cooldown = attacker.GetInit().(*cards.UnitCard).Cooldown

		// Berserk trait check - if defender was killed then allow attacker to attack again
		_, _, err := state.Board.GetUnitXY(defender.UUID)
		if err != nil && len(attacker.GetTraits(tr.BerserkTrait)) > 0 {
			attacker.Cooldown = 0
		}
	}
	return nil
}
