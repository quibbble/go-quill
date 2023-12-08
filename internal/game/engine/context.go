package engine

type Context string

const (
	CardCtx      Context = "Card"      // passed when a playing a card kicks off some event(s)
	TargetsCtx   Context = "Targets"   // passed when a playing a card that has a list of required targets
	HookCardCtx  Context = "HookCard"  // passed when a hook is triggered and represents the card that registered the hook
	HookEventCtx Context = "HookEvent" // passed when a hook is triggered and represents the event that triggered the hook
)
