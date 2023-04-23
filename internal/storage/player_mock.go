package storage

import (
	"errors"

	"github.com/vipulvpatil/airetreat-go/internal/model"
)

type PlayerAccessorMockSuccess struct {
	PlayerId string
}

func (p *PlayerAccessorMockSuccess) GetPlayerUsingTransaction(playerId string, transaction DatabaseTransaction) (*model.Player, error) {
	return model.NewPlayer(model.PlayerOptions{
		Id: p.PlayerId,
	})
}

func (p *PlayerAccessorMockSuccess) CreatePlayer(userId *string) (string, error) {
	return p.PlayerId, nil
}

func (p *PlayerAccessorMockSuccess) UpdatePlayerWithUserIdUsingTransaction(playerId, userId string, transaction DatabaseTransaction) error {
	return nil
}

type PlayerAccessorMockFailure struct {
}

func (p *PlayerAccessorMockFailure) GetPlayerUsingTransaction(playerId string, transaction DatabaseTransaction) (*model.Player, error) {
	return nil, errors.New("unable to get player")
}

func (p *PlayerAccessorMockFailure) CreatePlayer(userId *string) (string, error) {
	return "", errors.New("unable to create player")
}

func (p *PlayerAccessorMockFailure) UpdatePlayerWithUserIdUsingTransaction(playerId, userId string, transaction DatabaseTransaction) error {
	return errors.New("unable to update player")
}
