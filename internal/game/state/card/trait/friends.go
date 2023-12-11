package trait

import (
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	FriendsTrait = "Friends"
)

type FriendsArgs struct {
	ChooseUnits parse.Choose
	Trait       parse.Trait

	// DO NOT SET IN YAML - SET BY ENGINE
	// current units that have the trait applied
	Current []uuid.UUID
}
