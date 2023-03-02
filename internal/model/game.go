package model

import (
	"math/rand"
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
	stateTotalTime          int64
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
	StateTotalTime          int64
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
		stateTotalTime:          opts.StateTotalTime,
		lastQuestion:            opts.LastQuestion,
		lastQuestionTargetBotId: opts.LastQuestionTargetBotId,
		createdAt:               opts.CreatedAt,
		updatedAt:               opts.UpdatedAt,
		bots:                    opts.Bots,
		lastQuestionTargetBot:   opts.LastQuestionTargetBot,
	}, nil
}

func (game *Game) HasJustStarted() bool {
	return game.state == started
}

func (game *Game) GetOneRandomAiBot() (*Bot, error) {
	aiBots := []*Bot{}
	for _, bot := range game.bots {
		if bot.IsAi() {
			aiBots = append(aiBots, bot)
		}
	}

	if len(aiBots) <= 0 {
		return nil, errors.New("no AI bots in the game")
	}

	return getRandomBot(aiBots)
}

func getRandomBot(bots []*Bot) (*Bot, error) {
	if len(bots) == 0 {
		return nil, errors.Errorf("Cannot get random bot from an empty list")
	}

	rand.Shuffle(len(bots), func(i, j int) {
		bots[i], bots[j] = bots[j], bots[i]
	})

	return bots[0], nil
}

func (game *Game) BotWithId(botId string) *Bot {
	for _, bot := range game.bots {
		if bot.id == botId {
			return bot
		}
	}
	return nil
}

func (game *Game) botWithPlayerId(playerId string) *Bot {
	for _, bot := range game.bots {
		if bot.player != nil && bot.player.id == playerId {
			return bot
		}
	}
	return nil
}

func (game *Game) getTargetBot() *Bot {
	switch game.state {
	case waitingForBotQuestion, waitingForPlayerQuestion:
		return game.BotWithId(game.getCurrentTurnBotId())
	case waitingForBotAnswer, waitingForPlayerAnswer:
		return game.BotWithId(game.lastQuestionTargetBotId)
	default:
		return nil
	}
}

func (game *Game) getCurrentTurnBotId() string {
	turnIndex := game.currentTurnIndex % int64(len(game.turnOrder))
	return game.turnOrder[turnIndex]
}

func (game *Game) HasPlayer(playerId string) bool {
	if utilities.IsBlank(playerId) {
		return false
	}
	return game.botWithPlayerId(playerId) != nil
}

func (game *Game) StateHasBeenHandled() bool {
	return game.stateHandled
}

func (game *Game) IsInStatePlayersJoined() bool {
	return game.state == playersJoined
}

func (game *Game) RandomizedTurnOrder() []string {
	botIds := []string{}
	for _, bot := range game.bots {
		botIds = append(botIds, bot.id)
	}

	rand.Shuffle(len(botIds), func(i, j int) {
		botIds[i], botIds[j] = botIds[j], botIds[i]
	})
	return botIds
}
