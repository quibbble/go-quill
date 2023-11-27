package damage

import (
	"github.com/quibbble/go-quill/internal/state/card"
)

const (
	PureDamage     = "Pure"
	PhysicalDamage = "Physical"
	MagicDamage    = "Magic"
)

func Damage(unit *card.UnitCard, amount int, typ string) (int, error)
