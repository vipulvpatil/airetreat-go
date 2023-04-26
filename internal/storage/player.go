package storage

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

type PlayerAccessor interface {
	CreatePlayer() (*model.Player, error)
	GetPlayerUsingTransaction(playerId string, transaction DatabaseTransaction) (*model.Player, error)
	UpdatePlayerWithUserIdUsingTransaction(playerId, userId string, transaction DatabaseTransaction) (*model.Player, error)
	GetPlayerForUser(userId string) (*model.Player, error)
	CreatePlayerForUser(userId string) (*model.Player, error)
}

func (s *Storage) GetPlayerUsingTransaction(playerId string, transaction DatabaseTransaction) (*model.Player, error) {
	if utilities.IsBlank(playerId) {
		return nil, errors.New("playerId cannot be blank")
	}

	var nullableUserId sql.NullString

	row := transaction.QueryRow(`SELECT user_id FROM public."players" WHERE id = $1 FOR UPDATE`, playerId)
	err := row.Scan(&nullableUserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Errorf("getting player for %s: no such player", playerId)
		}
		return nil, errors.Errorf("getting player for %s: %v", playerId, err)
	}

	var userId *string
	if nullableUserId.Valid {
		userId = &nullableUserId.String
	}

	return model.NewPlayer(model.PlayerOptions{
		Id:     playerId,
		UserId: userId,
	})
}

func (s *Storage) CreatePlayer() (*model.Player, error) {
	id := s.IdGenerator.Generate()

	playerOpts := model.PlayerOptions{
		Id: id,
	}

	player, err := model.NewPlayer(playerOpts)
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to create player")
	}

	result, err := s.db.Exec(
		`INSERT INTO public."players" ("id") VALUES ($1)`,
		playerOpts.Id,
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, utilities.WrapBadError(err, "dbError while inserting player and changing db")
	}

	if rowsAffected != 1 {
		return nil, utilities.NewBadError(fmt.Sprintf("Very few or too many rows were affected when inserting player in db. This is highly unexpected. rowsAffected: %d", rowsAffected))
	}

	return player, nil
}

func (s *Storage) UpdatePlayerWithUserIdUsingTransaction(playerId, userId string, transaction DatabaseTransaction) (*model.Player, error) {
	if utilities.IsBlank(playerId) {
		return nil, errors.New("playerId cannot be blank")
	}

	if utilities.IsBlank(userId) {
		return nil, errors.New("userId cannot be blank")
	}

	result, err := transaction.Exec(
		`UPDATE public."players" SET "user_id" = $1 WHERE id = $2`,
		userId,
		playerId,
	)
	if err != nil {
		return nil, utilities.WrapBadError(err, "dbError while attempting player update")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, utilities.WrapBadError(err, "dbError while updating player and changing db")
	}

	if rowsAffected != 1 {
		return nil, utilities.NewBadError(fmt.Sprintf("Very few or too many rows were affected when updating player in db. This is highly unexpected. rowsAffected: %d", rowsAffected))
	}

	return model.NewPlayer(model.PlayerOptions{Id: playerId, UserId: &userId})
}

func (s *Storage) GetPlayerForUser(userId string) (*model.Player, error) {
	return nil, nil
}

func (s *Storage) CreatePlayerForUser(userId string) (*model.Player, error) {
	return nil, nil
}
