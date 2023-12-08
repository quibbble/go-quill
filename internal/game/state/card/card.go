package card

import (
	"context"
	"slices"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	CostStat     = "Cost"
	AttackStat   = "Attack"
	HealthStat   = "Health"
	MovementStat = "Movement"
	CooldownStat = "Cooldown"
	RangeState   = "Range"
)

type Card struct {
	Init parse.ICard

	UUID uuid.UUID

	Player uuid.UUID

	Cost int

	// Conditions required to play the card
	Conditions en.Conditions
	// Target requirements to play the card
	Targets []en.IChoose

	// Hooks registered on unit play
	Hooks []en.IHook
	// Events applied on unit play
	Events []en.IEvent

	// Traits that modify this card
	Traits []st.ITrait
}

type BuildTrait func(uuid uuid.UUID, typ string, args interface{}) (st.ITrait, error)

type Builders struct {
	en.BuildCondition
	en.BuildEvent
	en.BuildHook
	en.BuildChoose
	BuildTrait
	*uuid.Gen
}

func NewCard(builders *Builders, id string, player uuid.UUID) (st.ICard, error) {

	card, err := parse.ParseCard(id)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	var crd *parse.Card
	item, itemOk := card.(*parse.ItemCard)
	if itemOk {
		crd = &item.Card
	}
	spell, spellOk := card.(*parse.SpellCard)
	if spellOk {
		crd = &spell.Card
	}
	unit, unitOk := card.(*parse.UnitCard)
	if unitOk {
		crd = &unit.Card
	}

	if crd == nil {
		return nil, errors.ErrNilInterface
	}

	uuid := builders.Gen.New(rune(crd.ID[0]))

	buildConditions := func(cnds []parse.Condition) ([]en.ICondition, error) {
		conditions := make([]en.ICondition, 0)
		for _, c := range cnds {
			condition, err := builders.BuildCondition(builders.Gen.New(st.ConditionUUID), c.Type, c.Not, c.Args)
			if err != nil {
				return nil, errors.Wrap(err)
			}
			conditions = append(conditions, condition)
		}
		return conditions, nil
	}

	buildEvents := func(evts []parse.Event) ([]en.IEvent, error) {
		events := make([]en.IEvent, 0)
		for _, e := range evts {
			event, err := builders.BuildEvent(builders.Gen.New(st.EventUUID), e.Type, e.Args)
			if err != nil {
				return nil, errors.Wrap(err)
			}
			events = append(events, event)
		}
		return events, nil
	}

	conditions, err := buildConditions(crd.Conditions)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	targets := make([]en.IChoose, 0)
	for _, t := range crd.Targets {
		target, err := builders.BuildChoose(builders.Gen.New(st.ChooseUUID), t.Type, t.Args)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		targets = append(targets, target)
	}

	hooks := make([]en.IHook, 0)
	for _, h := range crd.Hooks {
		hookConditions, err := buildConditions(h.Conditions)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		hookEvents, err := buildEvents(h.Events)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		hookReuseConditions, err := buildConditions(h.ReuseConditions)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		hook, err := builders.BuildHook(builders.Gen.New(st.HookUUID), uuid, h.When, h.Type, hookConditions, hookEvents, hookReuseConditions)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		hooks = append(hooks, hook)
	}

	events, err := buildEvents(crd.Events)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	traits := make([]st.ITrait, 0)
	for _, t := range crd.Traits {
		trait, err := builders.BuildTrait(builders.Gen.New(st.TraitUUID), t.Type, t.Args)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		traits = append(traits, trait)
	}

	core := &Card{
		Init:       card,
		UUID:       uuid,
		Player:     player,
		Cost:       crd.Cost,
		Conditions: conditions,
		Targets:    targets,
		Hooks:      hooks,
		Events:     events,
		Traits:     traits,
	}

	if itemOk {
		heldTraits := make([]st.ITrait, 0)
		for _, trait := range item.HeldTraits {
			trait, err := builders.BuildTrait(builders.Gen.New(st.TraitUUID), trait.Type, trait.Args)
			if err != nil {
				return nil, errors.Wrap(err)
			}
			heldTraits = append(heldTraits, trait)
		}
		return &ItemCard{
			Card:       core,
			HeldTraits: heldTraits,
		}, nil
	} else if spellOk {
		return &SpellCard{
			Card: core,
		}, nil
	} else if unitOk {
		return &UnitCard{
			Card:       core,
			Type:       unit.Type,
			DamageType: unit.DamageType,
			Attack:     unit.Attack,
			Health:     unit.Health,
			Cooldown:   unit.Cooldown,
			Range:      unit.Range,
			Movement:   unit.Movement,
			Codex:      unit.Codex,
			Items:      make([]*ItemCard, 0),
		}, nil
	}
	return nil, errors.ErrNilInterface
}

func (c *Card) GetID() string {
	return c.Init.GetID()
}

func (c *Card) GetUUID() uuid.UUID {
	return c.UUID
}

func (c *Card) GetPlayer() uuid.UUID {
	return c.Player
}

func (c *Card) GetCost() int {
	return c.Cost
}

func (c *Card) SetCost(cost int) {
	c.Cost = cost
}

func (c *Card) GetInit() parse.ICard {
	return c.Init
}

func (c *Card) GetEvents() []en.IEvent {
	return c.Events
}

func (c *Card) GetHooks() []en.IHook {
	return c.Hooks
}

func (c *Card) Playable(engine en.IEngine, state en.IState) (bool, error) {
	return c.Conditions.Pass(context.Background(), engine, state)
}

func (c *Card) NextTargets(ctx context.Context, engine en.IEngine, state en.IState) ([]uuid.UUID, error) {
	targets := ctx.Value(en.TargetsCtx).([]uuid.UUID)
	if len(targets) > len(c.Targets) {
		return nil, errors.ErrIndexOutOfBounds
	}
	last := -1
	for i, target := range targets {
		choices, err := c.Targets[i].Retrieve(ctx, engine, state)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		if !slices.Contains(choices, target) {
			return nil, errors.Errorf("'%s' not a valid target", target)
		}
		last = i
	}
	// all targets are valid and card may be played with target list
	if last+1 == len(c.Targets) {
		return []uuid.UUID{}, nil
	}
	// get the next set of valid targets in the target chain
	return c.Targets[last+1].Retrieve(ctx, engine, state)
}

func (c *Card) GetTraits(typ string) []st.ITrait {
	traits := make([]st.ITrait, 0)
	for _, trait := range c.Traits {
		if trait.GetType() == typ {
			traits = append(traits, trait)
		}
	}
	return traits
}

func (c *Card) addTrait(engine en.IEngine, trait st.ITrait, card st.ICard) error {
	c.Traits = append(c.Traits, trait)
	return trait.Add(engine, card)
}

func (c *Card) removeTrait(engine en.IEngine, trait uuid.UUID, card st.ICard) error {
	idx := -1
	var tr st.ITrait
	for i, t := range c.Traits {
		if t.GetUUID() == trait {
			idx = i
			tr = t
		}
	}
	if idx < 0 {
		return errors.Errorf("failed to find card trait")
	}
	c.Traits = append(c.Traits[:idx], c.Traits[idx+1:]...)
	return tr.Remove(engine, card)
}
