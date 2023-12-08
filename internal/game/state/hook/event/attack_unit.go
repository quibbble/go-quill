package event

import (
	"context"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	tr "github.com/quibbble/go-quill/internal/game/state/card/trait"
	dm "github.com/quibbble/go-quill/internal/game/state/damage"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	AttackUnitEvent = "AttackUnit"
)

type AttackUnitArgs struct {
	ChooseAttacker parse.Choose
	ChooseDefender parse.Choose
}

func AttackUnitAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	var a AttackUnitArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	attackerChoice, err := GetUnitChoice(ctx, a.ChooseAttacker, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	aX, aY, err := state.Board.GetUnitXY(attackerChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	attacker := state.Board.XYs[aX][aY].Unit.(*cd.UnitCard)

	defenderChoice, err := GetUnitChoice(ctx, a.ChooseDefender, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	dX, dY, err := state.Board.GetUnitXY(defenderChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	defender := state.Board.XYs[dX][dY].Unit.(*cd.UnitCard)

	if attacker.Cooldown != 0 {
		return errors.Errorf("unit '%s' cannot attack due to cooldown", attacker.UUID)
	}
	if attacker.Range <= 0 && !attacker.CheckCodex(aX, aY, dX, dY) {
		return errors.Errorf("unit '%s' cannot attack due to failed codex check", attacker.UUID)
	}
	if attacker.Range > 0 && !attacker.CheckRange(aX, aY, dX, dY) {
		return errors.Errorf("unit '%s' cannot attack due to failed range check", attacker.UUID)
	}

	// thief trait check
	if len(attacker.GetTraits(tr.ThiefTrait)) > 0 && len(defender.Items) > 0 {
		item := defender.Items[state.Rand.Intn(len(defender.Items))]
		// set thief player as owner and add to that players hand to allow adding to thief
		item.Player = attacker.Player
		state.Hand[attacker.Player].Add(item)
		events := []*Event{
			{
				uuid: state.Gen.New(st.EventUUID),
				typ:  RemoveItemFromUnitEvent,
				args: &RemoveItemFromUnitArgs{
					ChooseItem: parse.Choose{
						Type: ch.UUIDChoice,
						Args: &ch.UUIDArgs{
							UUID: item.UUID,
						},
					},
					ChooseUnit: parse.Choose{
						Type: ch.UUIDChoice,
						Args: &ch.UUIDArgs{
							UUID: defender.UUID,
						},
					},
				},
				affect: RemoveItemFromUnitAffect,
			},
			{
				uuid: state.Gen.New(st.EventUUID),
				typ:  AddItemToUnitEvent,
				args: &AddItemToUnitArgs{
					ChooseItem: parse.Choose{
						Type: ch.UUIDChoice,
						Args: &ch.UUIDArgs{
							UUID: item.UUID,
						},
					},
					ChooseUnit: parse.Choose{
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
			if err := engine.Do(context.Background(), event, state); err != nil {
				return errors.Wrap(err)
			}
		}
		return nil
	}

	defenders := []*cd.UnitCard{defender}

	// lobber trait check
	if (len(attacker.GetTraits(tr.LobberTrait)) > 0) && attacker.Range > 0 {
		choose, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), ch.CodexChoice, &ch.CodexArgs{
			Types: []string{cd.CreatureUnit, cd.StructureUnit},
			Codex: attacker.Codex,
		})
		if err != nil {
			return errors.Wrap(err)
		}
		ctx := context.WithValue(context.Background(), en.TargetsCtx, []uuid.UUID{defender.UUID})
		uuids, err := choose.Retrieve(ctx, engine, state)
		if err != nil {
			return errors.Wrap(err)
		}
		for _, uuid := range uuids {
			x, y, err := state.Board.GetUnitXY(uuid)
			if err != nil {
				return errors.Wrap(err)
			}
			defenders = append(defenders, state.Board.XYs[x][y].Unit.(*cd.UnitCard))
		}
	}

	for _, defender := range defenders {
		attackerDamage, defenderDamage, err := dm.Battle(state, attacker, defender)
		if err != nil {
			return errors.Wrap(err)
		}
		if attackerDamage > 0 {
			var event *Event
			event = &Event{
				uuid: state.Gen.New(st.EventUUID),
				typ:  DamageUnitEvent,
				args: &DamageUnitArgs{
					DamageType: attacker.DamageType,
					Amount:     attackerDamage,
					ChooseUnit: parse.Choose{
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
				defender.Health < defender.GetInit().(*parse.UnitCard).Health {
				event = &Event{
					uuid: state.Gen.New(st.EventUUID),
					typ:  KillUnitEvent,
					args: &KillUnitArgs{
						ChooseUnit: parse.Choose{
							Type: ch.UUIDChoice,
							Args: &ch.UUIDArgs{
								UUID: defender.UUID,
							},
						},
					},
					affect: KillUnitAffect,
				}
			}

			if err := engine.Do(context.Background(), event, state); err != nil {
				return errors.Wrap(err)
			}

			// pillage trait check
			if defender.GetID() == "U0001" {
				pillages := attacker.GetTraits(tr.PillageTrait)
				for _, pillage := range pillages {
					args := pillage.GetArgs().(tr.PillageArgs)
					event, err := NewEvent(state.Gen.New(st.EventUUID), args.Event.Type, args.Event.Args)
					if err != nil {
						return errors.Wrap(err)
					}
					if err := engine.Do(context.Background(), event, state); err != nil {
						return errors.Wrap(err)
					}
				}
			}

			// gift trait check
			if defender.Type == cd.CreatureUnit {
				for _, gift := range attacker.GetTraits(tr.GiftTrait) {
					args := gift.GetArgs().(tr.GiftArgs)
					event, err := NewEvent(state.Gen.New(st.EventUUID), AddTraitToCard, &AddTraitToCardArgs{
						Trait: args.Trait,
						ChooseCard: parse.Choose{
							Type: ch.UUIDChoice,
							Args: &ch.UUIDArgs{
								UUID: defender.UUID,
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

		if defenderDamage > 0 && attacker.Range <= 0 {
			event := &Event{
				uuid: state.Gen.New(st.EventUUID),
				typ:  DamageUnitEvent,
				args: &DamageUnitArgs{
					DamageType: defender.DamageType,
					Amount:     defenderDamage,
					ChooseUnit: parse.Choose{
						Type: ch.UUIDChoice,
						Args: &ch.UUIDArgs{
							UUID: attacker.UUID,
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

	// if attacker still on board then reset cooldown
	if _, _, err := state.Board.GetUnitXY(attacker.UUID); err == nil {
		attacker.Cooldown = attacker.GetInit().(*parse.UnitCard).Cooldown

		// berserk trait check - if defender was killed then allow attacker to attack again
		_, _, err := state.Board.GetUnitXY(defender.UUID)
		if err != nil && len(attacker.GetTraits(tr.BerserkTrait)) > 0 {
			attacker.Cooldown = 0
		}
	}
	return nil
}
