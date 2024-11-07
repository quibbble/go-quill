package go_quill

import (
	"fmt"
	"strconv"
	"strings"

	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgn"
	"github.com/quibbble/go-quill/internal/game"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/parse"
)

const key = "Quill"

type Builder struct{}

func (b *Builder) Create(options *bg.BoardGameOptions) (bg.BoardGame, error) {
	return NewQuill(options)
}

func (b *Builder) CreateWithBGN(options *bg.BoardGameOptions) (bg.BoardGameWithBGN, error) {
	return NewQuill(options)
}

func (b *Builder) Load(game *bgn.Game) (bg.BoardGameWithBGN, error) {
	if game.Tags[bgn.GameTag] != key {
		return nil, loadFailure(fmt.Errorf("game tag does not match game key"))
	}
	teamsStr, ok := game.Tags[bgn.TeamsTag]
	if !ok {
		return nil, loadFailure(fmt.Errorf("missing teams tag"))
	}
	teams := strings.Split(teamsStr, ", ")
	seedStr, ok := game.Tags[bgn.SeedTag]
	if !ok {
		return nil, loadFailure(fmt.Errorf("missing seed tag"))
	}
	seed, err := strconv.Atoi(seedStr)
	if err != nil {
		return nil, loadFailure(err)
	}
	decksStr, ok := game.Tags[DecksTag]
	if !ok {
		return nil, loadFailure(fmt.Errorf("missing decks tag"))
	}
	decksList := strings.Split(decksStr, " : ")
	if len(decksList) != 2 {
		return nil, loadFailure(fmt.Errorf("must have two decks"))
	}
	decks := [][]string{strings.Split(decksList[0], ", "), strings.Split(decksList[1], ", ")}
	g, err := b.CreateWithBGN(&bg.BoardGameOptions{
		Teams: teams,
		MoreOptions: QuillMoreOptions{
			Seed:  int64(seed),
			Decks: decks,
		},
	})
	if err != nil {
		return nil, err
	}
	for _, action := range game.Actions {
		if action.TeamIndex >= len(teams) {
			return nil, loadFailure(fmt.Errorf("team index %d out of range", action.TeamIndex))
		}
		team := teams[action.TeamIndex]
		actionType := notationToAction[string(action.ActionKey)]
		if actionType == "" {
			return nil, loadFailure(fmt.Errorf("invalid action key %s", string(action.ActionKey)))
		}
		var details interface{}
		switch actionType {
		case ActionPlayCard:
			result, err := decodePlayCard(action.Details)
			if err != nil {
				return nil, err
			}
			details = result
		case ActionSackCard:
			result, err := decodeSackCard(action.Details)
			if err != nil {
				return nil, err
			}
			details = result
		case ActionAttackUnit:
			result, err := decodeAttackUnit(action.Details)
			if err != nil {
				return nil, err
			}
			details = result
		case ActionMoveUnit:
			result, err := decodeMoveUnit(action.Details)
			if err != nil {
				return nil, err
			}
			details = result
		case ActionEndTurn:
			result, err := decodeEndTurn(action.Details)
			if err != nil {
				return nil, err
			}
			details = result
		case bg.ActionSetWinners:
			result, err := bg.DecodeSetWinnersActionDetailsBGN(action.Details, teams)
			if err != nil {
				return nil, err
			}
			details = result
		}
		if err := g.Do(&bg.BoardGameAction{
			Team:        team,
			ActionType:  actionType,
			MoreDetails: details,
		}); err != nil {
			return nil, err
		}
	}
	return g, nil
}

func (b *Builder) Info() *bg.BoardGameInfo {
	ids, err := parse.AllCards()
	if err != nil {
		return nil
	}
	cards := make([]st.ICard, 0)
	for _, id := range ids {
		card, err := game.NewDummyCard(id)
		if err != nil {
			if err == parse.ErrNotEnabled {
				continue
			}
			return nil
		}
		cards = append(cards, card)
	}
	return &bg.BoardGameInfo{
		GameKey:  b.Key(),
		MinTeams: minTeams,
		MaxTeams: maxTeams,
		MoreInfo: &QuillMoreInfo{
			Cards: cards,
		},
	}
}

func (b *Builder) Key() string {
	return key
}
