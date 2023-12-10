package event

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	EndGameEvent = "EndGame"
)

type EndGameArgs struct {
	ChooseWinner parse.Choose
}

func EndGameAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	a := args.(*EndGameArgs)
	player, err := ch.GetPlayerChoice(ctx, a.ChooseWinner, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	state.Winner = &player
	return nil
}
