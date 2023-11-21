package state

type Mana struct {
	Amount int
	Limit  int
}

func NewMana() *Mana
