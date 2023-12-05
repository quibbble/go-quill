package parse

type ICard interface{}

type Card struct {
	ID          string
	Name        string
	Description string

	Cost int

	Conditions []Condition
	TargetReqs []TargetReq

	Hooks  []Hook
	Events []Event

	Traits []Traits
}

type ItemCard struct {
	Card `yaml:",inline" mapstructure:",squash"`

	HeldTraits []Traits
}

type SpellCard struct {
	Card `yaml:",inline" mapstructure:",squash"`
}

type UnitCard struct {
	Card `yaml:",inline" mapstructure:",squash"`

	Type       string
	DamageType string
	Attack     int
	Health     int
	Cooldown   int
	Movement   int
	Range      int
	Codex      string
}

type Args map[string]interface{}

type Condition struct {
	Type string
	Not  bool
	Args Args
}

type TargetReq struct {
	Type string
	Args Args
}

type Hook struct {
	When            string
	Type            string
	Conditions      []Condition
	Event           Event
	ReuseConditions []Condition
}

type Event struct {
	Type string
	Args Args
}

type Traits struct {
	Type string
	Args Args
}
