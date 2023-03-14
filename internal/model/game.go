package model

import (
	"fmt"
	"math/rand"
	"sort"
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

	if !utilities.IsBlank(opts.LastQuestionTargetBotId) {
		targetBotFound := false
		for _, bot := range opts.Bots {
			if bot.id == opts.LastQuestionTargetBotId {
				targetBotFound = true
			}
		}
		if !(targetBotFound) {
			return nil, errors.New("cannot create game with incorrect last question target bot id")
		}
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

func (game *Game) IsInStateWaitingForAiQuestion() bool {
	return game.state.isWaitingForAQuestion() && game.state.isWaitingOnAi()
}

func (game *Game) IsInStateWaitingForAiAnswer() bool {
	return game.state.isWaitingForAnAnswer() && game.state.isWaitingOnAi()
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

	if state.isWaitingOnAi() && !sourceBot.IsAi() {
		return nil, errors.New("expecting AI message but did not receive one")
	}

	if state.isWaitingOnHuman() && !sourceBot.IsHuman() {
		return nil, errors.New("expecting Human message but did not receive one")
	}

	update := GameUpdate{}
	var nextBot *Bot

	if state.isWaitingForAQuestion() {
		if sourceBotId == targetBotId {
			return nil, errors.New("questioning message should have different source and target bot")
		}
		nextBot = targetBot
	} else if state.isWaitingForAnAnswer() {
		if sourceBotId != targetBotId {
			return nil, errors.New("answering message should have same source and target bot")
		}
		nextBot = game.BotWithId(game.getNextTurnBotId())
	}

	update.State = getNewStateForNextBot(state, nextBot)

	if update.State.isWaitingForAQuestion() {
		nextIndex := game.currentTurnIndex + 1
		update.CurrentTurnIndex = &nextIndex
	} else if update.State.isWaitingForAnAnswer() {
		update.LastQuestion = &text
		update.LastQuestionTargetBotId = &(nextBot.id)
	}

	// TODO: Wondering if we really need stateHandled on all the games.
	stateHandled := false
	update.StateHandled = &stateHandled

	return &update, nil
}

func getNewStateForNextBot(currentState gameState, nextBot *Bot) gameState {
	if currentState.isWaitingForAQuestion() {
		if nextBot.IsAi() {
			return waitingForAiAnswer
		} else if nextBot.IsHuman() {
			return waitingForHumanAnswer
		}
	} else if currentState.isWaitingForAnAnswer() {
		if nextBot.IsAi() {
			return waitingForAiQuestion
		} else if nextBot.IsHuman() {
			return waitingForHumanQuestion
		}
	}
	return undefinedGameState
}

func (game *Game) expectedSourceBotIdForWaitingMessage() (string, error) {
	if !game.state.isWaitingForMessage() {
		return "", errors.New("this game is not waiting for messages currently")
	}

	if game.state.isWaitingForAQuestion() {
		return game.getCurrentTurnBotId(), nil
	}

	if game.state.isWaitingForAnAnswer() {
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
	if game.state.isWaitingForAQuestion() {
		return game.getCurrentTurnBot()
	}
	if game.state.isWaitingForAnAnswer() {
		return game.BotWithId(game.lastQuestionTargetBotId)
	}
	return nil
}

func (game *Game) GetCoversation() []string {
	messages := []Message{}
	for _, bot := range game.bots {
		// TODO: this is brittle. and based on the assumption that the messages are always question followed by answer.
		for j, message := range bot.messages {
			if j%2 == 0 {
				message.Text = fmt.Sprintf("Question: %s", message.Text)
			} else {
				message.Text = fmt.Sprintf("%s: %s", bot.name, message.Text)
			}
			messages = append(messages, message)
		}
	}
	sort.Sort(byCreatedAt(messages))
	conversation := []string{}
	for _, message := range messages {
		conversation = append(conversation, message.Text)
	}
	return conversation
}

type byCreatedAt []Message

func (m byCreatedAt) Len() int {
	return len(m)
}
func (m byCreatedAt) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
func (m byCreatedAt) Less(i, j int) bool {
	return m[i].CreatedAt.Before(m[j].CreatedAt)
}
