package trait

import (
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"

	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	RecodeTrait = "Recode"
)

type RecodeArgs struct {
	Code string
}

func AddRecode(t *Trait, card st.ICard) error {
	a := t.GetArgs().(*RecodeArgs)
	unit, ok := card.(*cd.UnitCard)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	return unit.Recode(a.Code)
}

func RemoveRecode(t *Trait, card st.ICard) error {
	a := t.GetArgs().(*RecodeArgs)
	unit, ok := card.(*cd.UnitCard)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	return unit.Recode(a.Code)
}
