package state

type Mana struct {
	Amount     int
	BaseAmount int
}

func NewMana() *Mana {
	return &Mana{
		Amount:     0,
		BaseAmount: 0,
	}
}
