package workers

import (
	"fmt"

	"github.com/gocraft/work"
	"github.com/pkg/errors"
)

type jobContext struct{}

func (j *jobContext) startGameOncePlayersHaveJoined(job *work.Job) error {
	gameId := job.ArgString("gameId")

	game, err := workerStorage.GetGame(gameId)
	if err != nil {
		return err
	}

	if game.StateHasBeenHandled() {
		return nil
	}

	if !game.IsInStatePlayersJoined() {
		return errors.Errorf("Game should be in PlayersJoined state: %s", gameId)
	}

	randomizedTurnOrder := game.RandomizedTurnOrder()

	firstTurnBot := game.BotWithId(randomizedTurnOrder[0])

	var newGameState string

	if firstTurnBot.IsAi() {
		newGameState = "WAITING_FOR_BOT_QUESTION"
	} else if firstTurnBot.IsHuman() {
		newGameState = "WAITING_FOR_PLAYER_QUESTION"
	}

	fmt.Println(newGameState)
	// save to db.

	// Temp.
	// process job for 30 seconds with sleep
	return nil
}
