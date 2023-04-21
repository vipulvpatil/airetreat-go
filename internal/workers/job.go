package workers

import (
	"math/rand"
	"time"

	"github.com/gocraft/work"
	"github.com/pkg/errors"
	aibot "github.com/vipulvpatil/airetreat-go/internal/services/ai-bot"
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
		newGameState = "WAITING_FOR_AI_QUESTION"
	} else if firstTurnBot.IsHuman() {
		newGameState = "WAITING_FOR_HUMAN_QUESTION"
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

	if !game.IsInStateWaitingForAiQuestion() {
		return errors.Errorf("game should be in WaitingForAiQuestion state: %s", gameId)
	}

	sourceBot := game.GetBotThatGameIsWaitingOn()
	targetBotId, err := game.GetTargetBotIdForNextQuestion()
	if err != nil {
		return err
	}
	aiBot := aibot.NewAiQuestionGenerator(
		aibot.AiBotOptions{
			BotId:        sourceBot.Id(),
			Game:         game,
			OpenAiClient: openAiClient,
		},
	)
	question := aiBot.GetNextQuestion()

	// Wait a random amount of time.
	time.Sleep(time.Duration(minDelayAfterAIResponse+rand.Intn(maxDelayAfterAIResponse-minDelayAfterAIResponse)) * time.Second)

	gameUpdate, err := game.GetGameUpdateAfterIncomingMessage(sourceBot.Id(), targetBotId, question)
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

	err = workerStorage.CreateMessageUsingTransaction(sourceBot.Id(), targetBotId, question, "question", tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
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

	if !game.IsInStateWaitingForAiAnswer() {
		return errors.Errorf("game should be in WaitingForAiAnswer state: %s", gameId)
	}

	sourceBot := game.GetBotThatGameIsWaitingOn()

	aiBot := aibot.NewAiAnswerGenerator(
		aibot.AiBotOptions{
			BotId:        sourceBot.Id(),
			Game:         game,
			OpenAiClient: openAiClient,
		},
	)
	answer := aiBot.GetNextAnswer()

	// Wait a random amount of time.
	time.Sleep(time.Duration(minDelayAfterAIResponse+rand.Intn(maxDelayAfterAIResponse-minDelayAfterAIResponse)) * time.Second)

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

	err = workerStorage.CreateMessageUsingTransaction(sourceBot.Id(), sourceBot.Id(), answer, "answer", tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
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
