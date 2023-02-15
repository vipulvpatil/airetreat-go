package model

import (
	"time"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

type Game struct {
	id                      string
	state                   gameState
	currentTurnIndex        int64
	turnOrder               []string
	stateHandled            bool
	stateHandledAt          *time.Time
	currentStateTotalTime   int64
	lastQuestion            string
	lastQuestionTargetBotId string
	createdAt               time.Time
	updatedAt               time.Time
	bots                    []*Bot
	lastQuestionTargetBot   *Bot
}

type GameOptions struct {
	Id                      string
	State                   string
	CurrentTurnIndex        int64
	TurnOrder               []string
	StateHandled            bool
	StateHandledAt          *time.Time
	CurrentStateTotalTime   int64
	LastQuestion            string
	LastQuestionTargetBotId string
	CreatedAt               time.Time
	UpdatedAt               time.Time
	Bots                    []*Bot
	LastQuestionTargetBot   *Bot
}

func NewGame(opts GameOptions) (*Game, error) {
	if utilities.IsBlank(opts.Id) {
		return nil, errors.New("cannot create game with an empty id")
	}

	state := GameState(opts.State)
	if !state.Valid() {
		return nil, errors.New("cannot create game with an invalid state")
	}

	if len(opts.TurnOrder) == 0 {
		return nil, errors.New("cannot create game with empty turn order array")
	}

	if len(opts.Bots) == 0 {
		return nil, errors.New("cannot create game with empty bots array")
	}

	return &Game{
		id:                      opts.Id,
		state:                   state,
		currentTurnIndex:        opts.CurrentTurnIndex,
		turnOrder:               opts.TurnOrder,
		stateHandled:            opts.StateHandled,
		stateHandledAt:          opts.StateHandledAt,
		currentStateTotalTime:   opts.CurrentStateTotalTime,
		lastQuestion:            opts.LastQuestion,
		lastQuestionTargetBotId: opts.LastQuestionTargetBotId,
		createdAt:               opts.CreatedAt,
		updatedAt:               opts.UpdatedAt,
		bots:                    opts.Bots,
		lastQuestionTargetBot:   opts.LastQuestionTargetBot,
	}, nil
}
