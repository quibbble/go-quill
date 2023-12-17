package go_quill

import (
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

// Action types
const (
	ActionNextTargets = "NextTargets"
	ActionPlayCard    = "PlayCard"
	ActionSackCard    = "SackCard"
	ActionAttackUnit  = "AttackUnit"
	ActionMoveUnit    = "MoveUnit"
	ActionEndTurn     = "EndTurn"
)

const (
	DecksTag = "Decks"
)

type QuillMoreOptions struct {
	Seed  int64
	Decks [][]string
}

type QuillSnapshotData struct {
	Board      [st.Cols][st.Rows]*st.Tile
	PlayRange  map[string][]int
	UUIDToTeam map[uuid.UUID]string
	Hand       map[string][]st.ICard
	Deck       map[string]int
	Mana       map[string]*st.Mana
	Sacked     map[string]bool
}

type NextTargetsActionDetails struct {
	Targets []uuid.UUID
}

type PlayCardActionDetails struct {
	Card    uuid.UUID
	Targets []uuid.UUID
}

type SackCardActionDetails struct {
	Card   uuid.UUID
	Option string
}

type AttackUnitActionDetails struct {
	Attacker, Defender uuid.UUID
}

type MoveUnitActionDetails struct {
	Unit, Tile uuid.UUID
}

type EndTurnActionDetails struct{}
