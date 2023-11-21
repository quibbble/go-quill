package cards

type Card interface {
}

type ItemCard struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`

	Cost uint `yaml:"cost"`

	Conditions []Condition `yaml:"conditions"`

	Traits []Traits `yaml:"traits"`
}

type SpellCard struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`

	Cost uint `yaml:"cost"`

	Conditions []Condition `yaml:"conditions"`
	TargetReqs []TargetReq `yaml:"target_reqs"`

	Events []Event `yaml:"events"`
	Hooks  []Hook  `yaml:"hooks"`
}

type UnitCard struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`

	Cost uint `yaml:"cost"`

	Conditions []Condition `yaml:"conditions"`

	Events []Event `yaml:"events"`
	Hooks  []Hook  `yaml:"hooks"`

	DamageType string   `yaml:"damage_type"`
	Attack     uint     `yaml:"attack"`
	Health     uint     `yaml:"health"`
	Cooldown   uint     `yaml:"cooldown"`
	Movement   uint     `yaml:"movement"`
	Codex      uint8    `yaml:"codex"`
	Traits     []Traits `yaml:"traits"`
}

type Details map[string]interface{}

type Condition struct {
	Type    string `yaml:"type"`
	Details `yaml:"details"`
}

type TargetReq struct {
	Type    string `yaml:"type"`
	Details `yaml:"details"`
}

type Hook struct {
	When            string      `yaml:"when"`
	Conditions      []Condition `yaml:"conditions"`
	Event           `yaml:"event"`
	ReuseConditions []Condition `yaml:"reuse_conditions"`
}

type Event struct {
	Type    string `yaml:"type"`
	Details `yaml:"details"`
}

type Traits struct {
	Type    string `yaml:"type"`
	Details `yaml:"details"`
}
