package damage

import (
	"github.com/quibbble/go-quill/internal/state/card"
)

func Damage(unit *card.UnitCard, amount int, typ string) (int, error)
