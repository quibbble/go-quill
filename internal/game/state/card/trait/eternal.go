package trait

import "github.com/quibbble/go-quill/parse"

const (
	EternalTrait = "Eternal"
)

type EternalArgs struct {
	Conditions []parse.Condition
	ChooseUnit parse.Choose
}
