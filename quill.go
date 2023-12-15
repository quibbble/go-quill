package go_quill

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/mitchellh/mapstructure"
	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
	"github.com/quibbble/go-boardgame/pkg/bgn"
	"github.com/quibbble/go-quill/internal/game"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	minTeams = 2
	maxTeams = 2
)

type Quill struct {
	state   *game.Game
	actions []*bg.BoardGameAction

	teams      []string
	teamToUUID map[string]uuid.UUID
	uuidToTeam map[uuid.UUID]string
	targets    []uuid.UUID
	options    *QuillMoreOptions
}

func NewQuill(options *bg.BoardGameOptions) (*Quill, error) {
	if len(options.Teams) < minTeams {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("at least %d teams required to create a game of %s", minTeams, key),
			Status: bgerr.StatusTooFewTeams,
		}
	} else if len(options.Teams) > maxTeams {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("at most %d teams allowed to create a game of %s", maxTeams, key),
			Status: bgerr.StatusTooManyTeams,
		}
	} else if duplicates(options.Teams) {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("duplicate teams found"),
			Status: bgerr.StatusInvalidOption,
		}
	}
	var details QuillMoreOptions
	if err := mapstructure.Decode(options.MoreOptions, &details); err != nil {
		return nil, &bgerr.Error{
			Err:    err,
			Status: bgerr.StatusInvalidOption,
		}
	}
	teamToUUID := make(map[string]uuid.UUID)
	uuidToTeam := make(map[uuid.UUID]string)
	gen := uuid.NewGen(rand.New(rand.NewSource(details.Seed)))
	for _, team := range options.Teams {
		playerUUID := gen.New('P')
		teamToUUID[team] = playerUUID
		uuidToTeam[playerUUID] = team
	}
	state, err := game.NewGame(details.Seed, teamToUUID[options.Teams[0]], teamToUUID[options.Teams[1]], details.Decks[0], details.Decks[1])
	if err != nil {
		return nil, &bgerr.Error{
			Err:    err,
			Status: bgerr.StatusInvalidOption,
		}
	}
	targets, err := state.GetNextTargets(state.GetTurn())
	if err != nil {
		return nil, &bgerr.Error{
			Err:    err,
			Status: bgerr.StatusInvalidAction,
		}
	}
	return &Quill{
		state:      state,
		actions:    make([]*bg.BoardGameAction, 0),
		targets:    targets,
		teams:      options.Teams,
		teamToUUID: teamToUUID,
		uuidToTeam: uuidToTeam,
		options:    &details,
	}, nil
}

func (q *Quill) Do(action *bg.BoardGameAction) error {
	if q.state.Winner != nil {
		return &bgerr.Error{
			Err:    fmt.Errorf("game already over"),
			Status: bgerr.StatusGameOver,
		}
	}
	team, ok := q.teamToUUID[action.Team]
	if !ok {
		return &bgerr.Error{
			Err:    fmt.Errorf("team not found"),
			Status: bgerr.StatusInvalidActionDetails,
		}
	}
	switch action.ActionType {
	case ActionNextTarget:
		var details NextTargetActionDetails
		if err := mapstructure.Decode(action.MoreDetails, &details); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		targets, err := q.state.GetNextTargets(team, details.Targets...)
		if err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidAction,
			}
		}
		q.targets = targets
	case ActionPlayCard:
		var details PlayCardActionDetails
		if err := mapstructure.Decode(action.MoreDetails, &details); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		if err := q.state.PlayCard(team, details.Card, details.Targets...); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidAction,
			}
		}
		q.actions = append(q.actions, action)
		targets, err := q.state.GetNextTargets(team)
		if err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidAction,
			}
		}
		q.targets = targets
	case ActionSackCard:
		var details SackCardActionDetails
		if err := mapstructure.Decode(action.MoreDetails, &details); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		if err := q.state.SackCard(team, details.Card, details.Option); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidAction,
			}
		}
		q.actions = append(q.actions, action)
		targets, err := q.state.GetNextTargets(team)
		if err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidAction,
			}
		}
		q.targets = targets
	case ActionAttackUnit:
		var details AttackUnitActionDetails
		if err := mapstructure.Decode(action.MoreDetails, &details); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		if err := q.state.AttackUnit(team, details.Attacker, details.Defender); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidAction,
			}
		}
		q.actions = append(q.actions, action)
		targets, err := q.state.GetNextTargets(team)
		if err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidAction,
			}
		}
		q.targets = targets
	case ActionMoveUnit:
		var details MoveUnitActionDetails
		if err := mapstructure.Decode(action.MoreDetails, &details); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		if err := q.state.MoveUnit(team, details.Unit, details.Tile); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidAction,
			}
		}
		q.actions = append(q.actions, action)
		targets, err := q.state.GetNextTargets(team)
		if err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidAction,
			}
		}
		q.targets = targets
	case bg.ActionSetWinners:
		var details bg.SetWinnersActionDetails
		if err := mapstructure.Decode(action.MoreDetails, &details); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		if len(details.Winners) != 1 {
			return &bgerr.Error{
				Err:    fmt.Errorf("there can only be one winner"),
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		winner, ok := q.teamToUUID[details.Winners[0]]
		if !ok {
			return &bgerr.Error{
				Err:    fmt.Errorf("team not found"),
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		q.state.Winner = &winner
		q.actions = append(q.actions, action)
		q.targets = make([]uuid.UUID, 0)
	default:
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot process action type %s", action.ActionType),
			Status: bgerr.StatusUnknownActionType,
		}
	}
	return nil
}

func (q *Quill) GetSnapshot(team ...string) (*bg.BoardGameSnapshot, error) {
	if len(team) > 1 {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("get snapshot requires zero or one team"),
			Status: bgerr.StatusTooManyTeams,
		}
	}
	winners := make([]string, 0)
	if q.state.Winner != nil {
		winners = append(winners, q.uuidToTeam[*q.state.Winner])
	}
	hand := make(map[string][]st.ICard)
	for id, h := range q.state.Hand {
		cards := h.GetItems()
		if len(team) == 1 && q.uuidToTeam[id] != team[0] {
			empty := make([]st.ICard, 0)
			for i := 0; i < len(cards); i++ {
				empty = append(empty, cd.NewEmptyCard(id))
			}
			cards = empty
		}
		hand[q.uuidToTeam[id]] = cards
	}
	deck := make(map[string]int)
	for id, d := range q.state.Deck {
		deck[q.uuidToTeam[id]] = d.GetSize()
	}
	mana := make(map[string]*st.Mana)
	for id, m := range q.state.Mana {
		mana[q.uuidToTeam[id]] = m
	}
	sacked := make(map[string]bool)
	for id, s := range q.state.Sacked {
		sacked[q.uuidToTeam[id]] = s
	}
	targets := make([]uuid.UUID, 0)
	if len(team) == 1 && q.state.GetTurn() == q.teamToUUID[team[0]] {
		targets = q.targets
	}
	return &bg.BoardGameSnapshot{
		Turn:    q.uuidToTeam[q.state.GetTurn()],
		Teams:   q.teams,
		Winners: winners,
		MoreData: QuillSnapshotData{
			Board:  q.state.Board.XYs,
			Hand:   hand,
			Deck:   deck,
			Mana:   mana,
			Sacked: sacked,
		},
		Targets: targets,
		Actions: q.actions,
	}, nil
}

func (q *Quill) GetBGN() *bgn.Game {
	tags := map[string]string{
		bgn.GameTag:  key,
		bgn.TeamsTag: strings.Join(q.teams, ", "),
		bgn.SeedTag:  fmt.Sprintf("%d", q.options.Seed),
		DecksTag: strings.Join([]string{
			strings.Join(q.options.Decks[0], ", "),
			strings.Join(q.options.Decks[1], ", "),
		}, " : "),
	}
	actions := make([]bgn.Action, 0)
	for _, action := range q.actions {
		bgnAction := bgn.Action{
			TeamIndex: indexOf(q.teams, action.Team),
			ActionKey: rune(actionToNotation[action.ActionType][0]),
		}
		switch action.ActionType {
		case ActionPlayCard:
			var details PlayCardActionDetails
			_ = mapstructure.Decode(action.MoreDetails, &details)
			bgnAction.Details = details.encodeBGN()
		case ActionSackCard:
			var details SackCardActionDetails
			_ = mapstructure.Decode(action.MoreDetails, &details)
			bgnAction.Details = details.encodeBGN()
		case ActionAttackUnit:
			var details AttackUnitActionDetails
			_ = mapstructure.Decode(action.MoreDetails, &details)
			bgnAction.Details = details.encodeBGN()
		case ActionMoveUnit:
			var details MoveUnitActionDetails
			_ = mapstructure.Decode(action.MoreDetails, &details)
			bgnAction.Details = details.encodeBGN()
		case bg.ActionSetWinners:
			var details bg.SetWinnersActionDetails
			_ = mapstructure.Decode(action.MoreDetails, &details)
			bgnAction.Details, _ = details.EncodeBGN(q.teams)
		}
		actions = append(actions, bgnAction)
	}
	return &bgn.Game{
		Tags:    tags,
		Actions: actions,
	}
}
