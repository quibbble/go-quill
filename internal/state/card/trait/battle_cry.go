package trait

const (
	BattleCryTrait = "BattleCry"
)

type BattleCryArgs struct {
	Event struct {
		Type string
		Args interface{}
	}
}
