package model

import (
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

const GAME_EXPIRY_DURATION = -4 * time.Hour

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
		return nil, errors.Errorf("cannot get random bot from an empty list")
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

func (game *Game) BotWithPlayerId(playerId string) *Bot {
	for _, bot := range game.bots {
		if bot.player != nil && bot.player.id == playerId {
			return bot
		}
	}
	return nil
}

func (game *Game) getCurrentTurnBotId() string {
	turnIndex := game.currentTurnIndex % int64(len(game.turnOrder))
	return game.turnOrder[turnIndex]
}

func (game *Game) getNextTurnBotId() string {
	turnIndex := (game.currentTurnIndex + 1) % int64(len(game.turnOrder))
	return game.turnOrder[turnIndex]
}

func (game *Game) HasPlayer(playerId string) bool {
	if utilities.IsBlank(playerId) {
		return false
	}
	return game.BotWithPlayerId(playerId) != nil
}

func (game *Game) StateHasBeenHandled() bool {
	return game.stateHandled
}

func (game *Game) IsInStatePlayersJoined() bool {
	return game.state == playersJoined
}

func (game *Game) IsInStateWaitingForBotQuestion() bool {
	return game.state.isQuestion() && game.state.isWaitingForAi()
}

func (game *Game) IsInStateWaitingForBotAnswer() bool {
	return game.state.isAnswer() && game.state.isWaitingForAi()
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

func (game *Game) RecentlyUpdated() bool {
	recent := time.Now().Add(GAME_EXPIRY_DURATION)
	return recent.Before(game.updatedAt)
}

type GameUpdate struct {
	State                   gameState
	CurrentTurnIndex        *int64
	StateHandled            *bool
	LastQuestion            *string
	LastQuestionTargetBotId *string
}

func (game *Game) GetGameUpdateAfterIncomingMessage(sourceBotId string, targetBotId string, text string) (*GameUpdate, error) {
	sourceBot := game.BotWithId(sourceBotId)
	targetBot := game.BotWithId(targetBotId)

	if sourceBot == nil {
		return nil, errors.New("invalid sourceBotId")
	}

	if targetBot == nil {
		return nil, errors.New("invalid targetBotId")
	}

	state := game.state
	expectedSourceBotId, err := game.expectedSourceBotIdForWaitingMessage()
	if err != nil {
		return nil, err
	}

	if expectedSourceBotId != sourceBotId {
		return nil, errors.New("incorrect sourceBotId")
	}

	if state.isWaitingForAi() && !sourceBot.IsAi() {
		return nil, errors.New("expecting AI message but did not receive one")
	}

	if state.isWaitingForHuman() && !sourceBot.IsHuman() {
		return nil, errors.New("expecting Human message but did not receive one")
	}

	update := GameUpdate{}
	var nextBot *Bot

	if state.isQuestion() {
		if sourceBotId == targetBotId {
			return nil, errors.New("questioning message should have different source and target bot")
		}
		nextBot = targetBot
	} else if state.isAnswer() {
		if sourceBotId != targetBotId {
			return nil, errors.New("answering message should have same source and target bot")
		}
		nextBot = game.BotWithId(game.getNextTurnBotId())
	}

	update.State = getNewStateForNextBot(state, nextBot)

	if update.State.isQuestion() {
		nextIndex := game.currentTurnIndex + 1
		update.CurrentTurnIndex = &nextIndex
	} else if update.State.isAnswer() {
		update.LastQuestion = &text
		update.LastQuestionTargetBotId = &(nextBot.id)
	}

	// TODO: Wondering if we really need stateHandled on all the games.
	stateHandled := false
	update.StateHandled = &stateHandled

	return &update, nil
}

func getNewStateForNextBot(currentState gameState, nextBot *Bot) gameState {
	if currentState.isQuestion() {
		if nextBot.IsAi() {
			return waitingForBotAnswer
		} else if nextBot.IsHuman() {
			return waitingForPlayerAnswer
		}
	} else if currentState.isAnswer() {
		if nextBot.IsAi() {
			return waitingForBotQuestion
		} else if nextBot.IsHuman() {
			return waitingForPlayerQuestion
		}
	}
	return undefinedGameState
}

func (game *Game) expectedSourceBotIdForWaitingMessage() (string, error) {
	if !game.state.isWaitingForMessage() {
		return "", errors.New("this game is not waiting for messages currently")
	}

	if game.state.isQuestion() {
		return game.getCurrentTurnBotId(), nil
	}

	if game.state.isAnswer() {
		return game.lastQuestionTargetBotId, nil
	}

	return "", utilities.NewBadError("game in an unexpected state")
}

func (game *Game) GetTargetBotForNextQuestion() (*Bot, error) {
	possibleTargetBotsForNextQuestion := []*Bot{}
	currentTurnBotId := game.getCurrentTurnBotId()
	for _, bot := range game.bots {
		if currentTurnBotId != bot.id {
			possibleTargetBotsForNextQuestion = append(possibleTargetBotsForNextQuestion, bot)
		}
	}

	botsWithLeastNumberOfMessages := getBotsWithLeastNumberOfMessages(possibleTargetBotsForNextQuestion)
	randomBotWithLeastNumberOfMessages, err := getRandomBot(botsWithLeastNumberOfMessages)
	if err != nil {
		return nil, err
	}

	return randomBotWithLeastNumberOfMessages, nil
}

func getBotsWithLeastNumberOfMessages(bots []*Bot) []*Bot {
	if len(bots) == 0 {
		return bots
	}
	leastNumberOfMessages := len(bots[0].messages)
	for _, bot := range bots {
		if len(bot.messages) < leastNumberOfMessages {
			leastNumberOfMessages = len(bot.messages)
		}
	}

	botsWithLeastNumberOfMessages := []*Bot{}
	for _, bot := range bots {
		if len(bot.messages) == leastNumberOfMessages {
			botsWithLeastNumberOfMessages = append(botsWithLeastNumberOfMessages, bot)
		}
	}

	return botsWithLeastNumberOfMessages
}

func (game *Game) getCurrentTurnBot() *Bot {
	return game.BotWithId(game.getCurrentTurnBotId())
}

func (game *Game) GetBotThatGameIsWaitingOn() *Bot {
	if game.state.isQuestion() {
		return game.getCurrentTurnBot()
	}
	if game.state.isAnswer() {
		return game.BotWithId(game.lastQuestionTargetBotId)
	}
	return nil
}
