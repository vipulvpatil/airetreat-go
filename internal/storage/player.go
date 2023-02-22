package storage

import (
	"fmt"

	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

type PlayerCreator interface {
	CreatePlayer() (string, error)
}

func (s *Storage) CreatePlayer() (string, error) {
	id := s.IdGenerator.Generate()

	playerOpts := model.PlayerOptions{
		Id: id,
	}

	_, err := model.NewPlayer(playerOpts)
	if err != nil {
		return "", utilities.WrapBadError(err, "failed to create player")
	}

	result, err := s.db.Exec(
		`INSERT INTO public."players" ("id") VALUES ($1)`,
		playerOpts.Id,
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
