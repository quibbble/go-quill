package event

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/internal/state/damage"
	"github.com/quibbble/go-quill/internal/state/hook/choose"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	AttackUnitEvent = "attack_unit"
)

type AttackUnitArgs struct {
	X, Y int
	ch.Choose
}

func AttackUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	a, ok := args.(AttackUnitArgs)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	choices, err := a.Choose.Retrieve(engine, state, targets...)
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
	attacker := state.Board.XYs[x][y].Unit
	if attacker.Range <= 0 && !attacker.CheckCodex(x, y, a.X, a.Y) {
		return errors.Errorf("unit '%s' cannot attack due to failed codex check", attacker.UUID)
	}
	if state.Board.XYs[a.X][a.Y].Unit == nil {
		return errors.Errorf("unit '%s' cannot attack at (%d,%d) as no unit exists", attacker.UUID, a.X, a.Y)
	}
	defender := state.Board.XYs[a.X][a.Y].Unit
	attackerDamage, defenderDamage, err := damage.Battle(attacker, defender)
	if err != nil {
		return errors.Wrap(err)
	}
	if attackerDamage > 0 {
		event := &Event{
			uuid: uuid.New(st.EventUUID),
			typ:  DamageUnitEvent,
			args: &DamageUnitArgs{
				DamageType: attacker.DamageType,
				Amount:     attackerDamage,
				Choose: &choose.UUIDChoice{
					UUID: defender.UUID,
				},
			},
			affect: DamageUnitAffect,
		}
		if err := engine.Do(event, state); err != nil {
			return errors.Wrap(err)
		}
	}
	if defenderDamage > 0 {
		event := &Event{
			uuid: uuid.New(st.EventUUID),
			typ:  DamageUnitEvent,
			args: &DamageUnitArgs{
				DamageType: defender.DamageType,
				Amount:     defenderDamage,
				Choose: &choose.UUIDChoice{
					UUID: attacker.UUID,
				},
			},
			affect: DamageUnitAffect,
		}
		if err := engine.Do(event, state); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
