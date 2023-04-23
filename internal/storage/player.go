package storage

import (
	"fmt"

	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

type PlayerAccessor interface {
	CreatePlayer(userId *string) (string, error)
	UpdatePlayer(playerId string, userId string) error
}

func (s *Storage) CreatePlayer(userId *string) (string, error) {
	id := s.IdGenerator.Generate()

	playerOpts := model.PlayerOptions{
		Id:     id,
		UserId: userId,
	}

	_, err := model.NewPlayer(playerOpts)
	if err != nil {
		return "", utilities.WrapBadError(err, "failed to create player")
	}

	result, err := s.db.Exec(
		`INSERT INTO public."players" ("id", "user_id") VALUES ($1, $2)`,
		playerOpts.Id,
		playerOpts.UserId,
	)
	if err != nil {
		return "", err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", utilities.WrapBadError(err, "dbError while inserting player and changing db")
	}

	if rowsAffected != 1 {
		return "", utilities.NewBadError(fmt.Sprintf("Very few or too many rows were affected when inserting player in db. This is highly unexpected. rowsAffected: %d", rowsAffected))
	}

	return playerOpts.Id, nil
}

func (s *Storage) UpdatePlayer(playerId string, userId string) error {
	if utilities.IsBlank(playerId) {
		return utilities.NewBadError("playerId cannot be blank")
	}

	if utilities.IsBlank(userId) {
		return utilities.NewBadError("userId cannot be blank")
	}

	result, err := s.db.Exec(
		`UPDATE public."players" SET "user_id" = $1 WHERE id = $2`,
		userId,
		playerId,
	)
	if err != nil {
		return utilities.WrapBadError(err, "dbError while attempting player update")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utilities.WrapBadError(err, "dbError while updating player and changing db")
	}

	if rowsAffected != 1 {
		return utilities.NewBadError(fmt.Sprintf("Very few or too many rows were affected when updating player in db. This is highly unexpected. rowsAffected: %d", rowsAffected))
	}

	return nil
}
