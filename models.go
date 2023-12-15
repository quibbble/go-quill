package go_quill

import (
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

// Action types
const (
	ActionNextTarget = "NextTarget"
	ActionPlayCard   = "PlayCard"
	ActionSackCard   = "SackCard"
	ActionAttackUnit = "AttackUnit"
	ActionMoveUnit   = "MoveUnit"
)

const (
	DecksTag = "Decks"
)

type QuillMoreOptions struct {
	Seed  int64
	Decks [][]string
}

type QuillSnapshotData struct {
	Board  [st.Cols][st.Rows]*st.Tile
	Hand   map[string][]st.ICard
	Deck   map[string]int
	Mana   map[string]*st.Mana
	Sacked map[string]bool
}

type NextTargetActionDetails struct {
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