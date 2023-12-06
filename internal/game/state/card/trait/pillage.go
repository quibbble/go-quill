package trait

const (
	PillageTrait = "Pillage"
)

type PillageArgs struct {
	Event struct {
		Type string
		Args interface{}
	}
}
