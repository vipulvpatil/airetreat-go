package workers

import (
	"github.com/gocraft/work"
	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

type jobContext struct{}

func (j *jobContext) startGameOncePlayersHaveJoined(job *work.Job) error {
	gameId := job.ArgString("gameId")

	if utilities.IsBlank(gameId) {
		return errors.New("gameId is required")
	}
	game, err := workerStorage.GetGame(gameId)
	if err != nil {
		return err
	}

	if game.StateHasBeenHandled() {
		return errors.Errorf("game has already been handled: %s", gameId)
	}

	if !game.IsInStatePlayersJoined() {
		return errors.Errorf("game should be in PlayersJoined state: %s", gameId)
	}

	randomizedTurnOrder := game.RandomizedTurnOrder()

	firstTurnBot := game.BotWithId(randomizedTurnOrder[0])

	var newGameState string

	if firstTurnBot.IsAi() {
		newGameState = "WAITING_FOR_BOT_QUESTION"
	} else if firstTurnBot.IsHuman() {
		newGameState = "WAITING_FOR_PLAYER_QUESTION"
	}

	startTurnIndex := int64(0)

	updateOpts := storage.GameUpdateOptions{
		State:            &newGameState,
		CurrentTurnIndex: &startTurnIndex,
		TurnOrder:        randomizedTurnOrder,
	}

	return workerStorage.UpdateGameState(gameId, updateOpts)
}

func (j *jobContext) askQuestionOnBehalfOfBot(job *work.Job) error {
	// WAITING_FOR_BOT_QUESTION
	// verify state
	// identify question target.
	// get question from AI
	// update state. and
	// set handled, if question target is HUMAN
	// set unhandled, if question target is AI
	return nil
}

func (j *jobContext) answerQuestionOnBehalfOfBot(job *work.Job) error {
	// WAITING_FOR_BOT_ANSWER
	// verify state
	// get answer from AI
	// update state. and
	// set handled, if next turn is HUMAN
	// set unhandled, if next turn is AI
	return nil
}

func (j *jobContext) deleteExpiredGames(job *work.Job) error {
	gameId := job.ArgString("gameId")

	if utilities.IsBlank(gameId) {
		return errors.New("gameId is required")
	}
	game, err := workerStorage.GetGame(gameId)
	if err != nil {
		return err
	}

	if game.RecentlyUpdated() {
		return nil
	}

	return workerStorage.DeleteGame(gameId)
}
