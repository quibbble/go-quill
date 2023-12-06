package trait

const (
	DeathCryTrait = "DeathCry"
)

type DeathCryArgs struct {
	Event struct {
		Type string
		Args interface{}
	}
}
