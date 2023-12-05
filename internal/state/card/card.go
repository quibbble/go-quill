package card

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	CostStat      = "Cost"
	AttackStat    = "Attack"
	HealthStat    = "Health"
	MovementStat  = "Movement"
	CooldownState = "Cooldown"
	RangeState    = "Range"
)

type Card struct {
	init *parse.Card

	UUID uuid.UUID

	Player uuid.UUID

	Cost int

	// Conditions required to play the card
	Conditions en.Conditions
	// Target requirements to play the card
	TargetReqs []en.ITargetReq

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
	en.BuildTargetReq
	BuildTrait
	*uuid.Gen
}

func NewCard(builders *Builders, card *parse.Card, player uuid.UUID) (*Card, error) {
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

	conditions, err := buildConditions(card.Conditions)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	targetReqs := make([]en.ITargetReq, 0)
	for _, t := range card.TargetReqs {
		targetReq, err := builders.BuildTargetReq(builders.Gen.New(st.TargetReqUUID), t.Type, t.Args)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		targetReqs = append(targetReqs, targetReq)
	}

	hooks := make([]en.IHook, 0)
	for _, h := range card.Hooks {
		hookConditions, err := buildConditions(h.Conditions)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		hookEvent, err := builders.BuildEvent(builders.Gen.New(st.EventUUID), h.Event.Type, h.Event.Args)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		hookReuseConditions, err := buildConditions(h.ReuseConditions)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		hook, err := builders.BuildHook(builders.Gen.New(st.HookUUID), h.When, h.Type, hookConditions, hookEvent, hookReuseConditions)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		hooks = append(hooks, hook)
	}

	events := make([]en.IEvent, 0)
	for _, e := range card.Events {
		event, err := builders.BuildEvent(builders.Gen.New(st.EventUUID), e.Type, e.Args)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		events = append(events, event)
	}

	traits := make([]st.ITrait, 0)
	for _, t := range card.Traits {
		trait, err := builders.BuildTrait(builders.Gen.New(st.TraitUUID), t.Type, t.Args)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		traits = append(traits, trait)
	}

	return &Card{
		init:       card,
		UUID:       builders.Gen.New(rune(card.ID[0])),
		Player:     player,
		Cost:       card.Cost,
		Conditions: conditions,
		TargetReqs: targetReqs,
		Hooks:      hooks,
		Events:     events,
		Traits:     traits,
	}, nil
}

func (c *Card) GetID() string {
	return c.init.ID
}

func (c *Card) GetUUID() uuid.UUID {
	return c.UUID
}

func (c *Card) GetPlayer() uuid.UUID {
	return c.Player
}

func (c *Card) GetInit() parse.ICard {
	return c.init
}

func (c *Card) Playable(engine en.IEngine, state en.IState) (bool, error) {
	return c.Conditions.Pass(engine, state)
}

func (c *Card) ValidTargets(engine en.IEngine, state en.IState, targets ...uuid.UUID) (bool, error) {
	if len(targets) != len(c.TargetReqs) {
		return false, nil
	}
	pass := true
	for i, req := range c.TargetReqs {
		p, err := req.Validate(engine, state, targets[i], targets[:i]...)
		if err != nil {
			return false, errors.Wrap(err)
		}
		pass = p && pass
	}
	return pass, nil
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

func (c *Card) AddTrait(engine en.IEngine, trait st.ITrait) error {
	c.Traits = append(c.Traits, trait)
	return trait.Add(engine, c)
}

func (c *Card) RemoveTrait(engine en.IEngine, trait uuid.UUID) error {
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
	return tr.Remove(engine, c)
}
