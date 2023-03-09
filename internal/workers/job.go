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
	gameId := job.ArgString("gameId")

	if utilities.IsBlank(gameId) {
		return errors.New("gameId is required")
	}

	tx, err := workerStorage.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	game, err := workerStorage.GetGameUsingTransaction(gameId, tx)
	if err != nil {
		return err
	}

	if game.StateHasBeenHandled() {
		return errors.Errorf("game has already been handled: %s", gameId)
	}

	if !game.IsInStateWaitingForBotQuestion() {
		return errors.Errorf("game should be in WaitingForBotQuestion state: %s", gameId)
	}

	sourceBot := game.GetBotThatGameIsWaitingOn()
	targetBot, err := game.GetTargetBotForNextQuestion()
	if err != nil {
		return err
	}

	// TODO: Get question from OpenAiClient
	question := "Some question from AI"
	gameUpdate, err := game.GetGameUpdateAfterIncomingMessage(sourceBot.Id(), targetBot.Id(), question)
	if err != nil {
		return err
	}

	newGameState := gameUpdate.State.String()
	updateOptions := storage.GameUpdateOptions{
		State:                   &newGameState,
		CurrentTurnIndex:        gameUpdate.CurrentTurnIndex,
		StateHandled:            gameUpdate.StateHandled,
		LastQuestion:            gameUpdate.LastQuestion,
		LastQuestionTargetBotId: gameUpdate.LastQuestionTargetBotId,
	}

	err = workerStorage.UpdateGameStateUsingTransaction(gameId, updateOptions, tx)
	if err != nil {
		return err
	}

	err = workerStorage.CreateMessageUsingTransaction(targetBot.Id(), question, tx)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (j *jobContext) answerQuestionOnBehalfOfBot(job *work.Job) error {
	gameId := job.ArgString("gameId")

	if utilities.IsBlank(gameId) {
		return errors.New("gameId is required")
	}

	tx, err := workerStorage.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	game, err := workerStorage.GetGameUsingTransaction(gameId, tx)
	if err != nil {
		return err
	}

	if game.StateHasBeenHandled() {
		return errors.Errorf("game has already been handled: %s", gameId)
	}

	if !game.IsInStateWaitingForBotAnswer() {
		return errors.Errorf("game should be in WaitingForBotAnswer state: %s", gameId)
	}

	sourceBot := game.GetBotThatGameIsWaitingOn()

	// TODO: Get question from OpenAiClient
	answer := "Some answer from AI"
	gameUpdate, err := game.GetGameUpdateAfterIncomingMessage(sourceBot.Id(), sourceBot.Id(), answer)
	if err != nil {
		return err
	}

	newGameState := gameUpdate.State.String()
	updateOptions := storage.GameUpdateOptions{
		State:                   &newGameState,
		CurrentTurnIndex:        gameUpdate.CurrentTurnIndex,
		StateHandled:            gameUpdate.StateHandled,
		LastQuestion:            gameUpdate.LastQuestion,
		LastQuestionTargetBotId: gameUpdate.LastQuestionTargetBotId,
	}

	err = workerStorage.UpdateGameStateUsingTransaction(gameId, updateOptions, tx)
	if err != nil {
		return err
	}

	err = workerStorage.CreateMessageUsingTransaction(sourceBot.Id(), answer, tx)
	if err != nil {
		return err
	}

	tx.Commit()
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
