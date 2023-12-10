package event

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	tr "github.com/quibbble/go-quill/internal/game/state/card/trait"
	dm "github.com/quibbble/go-quill/internal/game/state/damage"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/maths"
)

const (
	AttackUnitEvent = "AttackUnit"
)

type AttackUnitArgs struct {
	ChooseAttacker parse.Choose
	ChooseDefender parse.Choose
}

func AttackUnitAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	a := args.(*AttackUnitArgs)
	attackerChoice, err := ch.GetUnitChoice(ctx, a.ChooseAttacker, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	aX, aY, err := state.Board.GetUnitXY(attackerChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	attacker := state.Board.XYs[aX][aY].Unit.(*cd.UnitCard)

	defenderChoice, err := ch.GetUnitChoice(ctx, a.ChooseDefender, engine, state)
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
		event1, err := NewEvent(state.Gen.New(en.EventUUID), RemoveItemFromUnitEvent, RemoveItemFromUnitArgs{
			ChooseItem: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: item.UUID,
				},
			},
			ChooseUnit: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: defender.UUID,
				},
			},
		})
		if err != nil {
			return errors.Wrap(err)
		}
		event2, err := NewEvent(state.Gen.New(en.EventUUID), AddItemToUnitEvent, AddItemToUnitArgs{
			ChooseItem: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: item.UUID,
				},
			},
			ChooseUnit: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: attacker.UUID,
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

	defenders := []*cd.UnitCard{defender}

	// lobber trait check
	if (len(attacker.GetTraits(tr.LobberTrait)) > 0) && attacker.Range > 0 {
		choose, err := ch.NewChoose(state.Gen.New(en.ChooseUUID), ch.CodexChoice, &ch.CodexArgs{
			Types: []string{cd.CreatureUnit, cd.StructureUnit},
			Codex: attacker.Codex,
			ChooseUnitOrTile: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: defender.UUID,
				},
			},
		})
		if err != nil {
			return errors.Wrap(err)
		}
		uuids, err := choose.Retrieve(context.Background(), engine, state)
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
			var event en.IEvent
			event, err = NewEvent(state.Gen.New(en.EventUUID), DamageUnitEvent, DamageUnitArgs{
				DamageType: attacker.DamageType,
				Amount:     attackerDamage,
				ChooseUnit: parse.Choose{
					Type: ch.UUIDChoice,
					Args: ch.UUIDArgs{
						UUID: defender.UUID,
					},
				},
			})
			if err != nil {
				return errors.Wrap(err)
			}

			// execute trait check
			if len(attacker.GetTraits(tr.ExecuteTrait)) > 0 &&
				defender.Health < defender.GetInit().(*parse.UnitCard).Health {
				event, err = NewEvent(state.Gen.New(en.EventUUID), KillUnitEvent, KillUnitArgs{
					ChooseUnit: parse.Choose{
						Type: ch.UUIDChoice,
						Args: ch.UUIDArgs{
							UUID: defender.UUID,
						},
					},
				})
				if err != nil {
					return errors.Wrap(err)
				}
			}

			if err := engine.Do(context.Background(), event, state); err != nil {
				return errors.Wrap(err)
			}

			// pillage trait check
			if defender.GetID() == "U0001" {
				for _, trait := range attacker.GetTraits(tr.PillageTrait) {
					args := trait.GetArgs().(*tr.PillageArgs)
					for _, h := range args.Hooks {
						hook, err := state.NewHook(state.Gen, attacker.GetUUID(), h)
						if err != nil {
							return errors.Wrap(err)
						}
						engine.Register(hook)
					}
					for _, e := range args.Events {
						event, err := NewEvent(state.Gen.New(en.EventUUID), e.Type, e.Args)
						if err != nil {
							return errors.Wrap(err)
						}
						if err := engine.Do(context.Background(), event, state); err != nil {
							return errors.Wrap(err)
						}
					}
				}
			}

			// gift trait check
			if defender.Type == cd.CreatureUnit {
				for _, trait := range attacker.GetTraits(tr.GiftTrait) {
					args := trait.GetArgs().(*tr.GiftArgs)
					event, err := NewEvent(state.Gen.New(en.EventUUID), AddTraitToCard, &AddTraitToCardArgs{
						Trait: args.Trait,
						ChooseCard: parse.Choose{
							Type: ch.UUIDChoice,
							Args: ch.UUIDArgs{
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
			event, err := NewEvent(state.Gen.New(en.EventUUID), DamageUnitEvent, DamageUnitArgs{
				DamageType: defender.DamageType,
				Amount:     defenderDamage,
				ChooseUnit: parse.Choose{
					Type: ch.UUIDChoice,
					Args: ch.UUIDArgs{
						UUID: attacker.UUID,
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

	// if attacker still on board then reset cooldown
	if _, _, err := state.Board.GetUnitXY(attacker.UUID); err == nil {
		attacker.Cooldown = maths.MaxInt(0, attacker.BaseCooldown)

		// berserk trait check - if defender was killed then allow attacker to attack again
		_, _, err := state.Board.GetUnitXY(defender.UUID)
		if err != nil && len(attacker.GetTraits(tr.BerserkTrait)) > 0 {
			attacker.Cooldown = 0
		}
	}
	return nil
}
