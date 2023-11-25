package state

type Mana struct {
	Amount  int
	Refresh int
}

func NewMana() *Mana
