package engine

type Context string

const (
	CardCtx      Context = "Card"
	HookCardCtx  Context = "HookCard"
	HookEventCtx Context = "HookEvent"
	TargetsCtx   Context = "Targets"
)

// Context
// - hook, condition, choose, event all have a ctx arg
// - ctx contains the following
//   - Card - passed when a card kicks off an event
//   - HookCard - passed when a hook is triggered
//   - TriggerEvent - passed when a hook is triggered
//   - Targets - []uuid.UUID

// Self refers to HookCard when it exists
// Otherwise Self refers to Card
// Error if neither are found

// TriggerEvent is used in conditions
