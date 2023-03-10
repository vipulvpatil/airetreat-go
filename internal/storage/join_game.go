package storage

import (
	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

func (s *Storage) JoinGame(gameId, playerId string) error {
	if utilities.IsBlank(gameId) {
		return errors.New("gameId cannot be blank")
	}

	if utilities.IsBlank(playerId) {
		return errors.New("playerId cannot be blank")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return utilities.WrapBadError(err, "failed to start db transaction")
	}
	defer tx.Rollback()

	game, err := getGame(tx, gameId)
	if err != nil {
		return err
	}

	if !game.HasJustStarted() {
		return errors.Errorf("Cannot join this game: %v", gameId)
	}

	if game.HasPlayer(playerId) {
		return nil
	}

	aiBot, err := game.GetOneRandomAiBot()
	if err != nil {
		return err
	}

	err = connectPlayerToBot(tx, playerId, aiBot.Id())
	if err != nil {
		return err
	}

	err = updateGameStateToPlayersJoined(tx, gameId)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}
