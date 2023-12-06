package trait

const (
	GiftTrait = "Gift"
)

type GiftArgs struct {
	Trait struct {
		Type string
		Args interface{}
	}
}
