package card

type ItemCard struct {
	Card

	// Traits applied to a unit when an item is held
	HeldTraits []Trait
}
