package damage

import (
	"github.com/quibbble/go-quill/internal/state/card"
)

const (
	DamageTypePure     = "pure"
	DamageTypePhysical = "physical"
	DamageTypeMagic    = "magic"
)

func Damage(unit *card.UnitCard, amount int, typ string) (int, error)
