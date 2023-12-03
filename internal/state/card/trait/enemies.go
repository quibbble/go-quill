package trait

const (
	EnemiesTrait = "Enemies"
)

type EnemiesArgs struct {
	Location string
	Trait    struct {
		Type string
		Args interface{}
	}
}
