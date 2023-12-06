package trait

const (
	EnrageTrait = "Enrage"
)

type EnrageArgs struct {
	Event struct {
		Type string
		Args interface{}
	}
}
