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

func (p *PlayerAccessorMockSuccess) CreatePlayer() (*model.Player, error) {
	return model.NewPlayer(model.PlayerOptions{
		Id: p.PlayerId,
	})
}

func (p *PlayerAccessorMockSuccess) UpdatePlayerWithUserIdUsingTransaction(playerId, userId string, transaction DatabaseTransaction) (*model.Player, error) {
	return model.NewPlayer(model.PlayerOptions{
		Id: p.PlayerId,
	})
}

func (p *PlayerAccessorMockSuccess) GetPlayerForUserIfExists(userId string) (*model.Player, error) {
	return model.NewPlayer(model.PlayerOptions{
		Id: p.PlayerId,
	})
}

func (p *PlayerAccessorMockSuccess) CreatePlayerForUser(userId string) (*model.Player, error) {
	return model.NewPlayer(model.PlayerOptions{
		Id: p.PlayerId,
	})
}

type PlayerAccessorMockFailure struct {
}

func (p *PlayerAccessorMockFailure) GetPlayerUsingTransaction(playerId string, transaction DatabaseTransaction) (*model.Player, error) {
	return nil, errors.New("unable to get player")
}

func (p *PlayerAccessorMockFailure) CreatePlayer() (*model.Player, error) {
	return nil, errors.New("unable to create player")
}

func (p *PlayerAccessorMockFailure) UpdatePlayerWithUserIdUsingTransaction(playerId, userId string, transaction DatabaseTransaction) (*model.Player, error) {
	return nil, errors.New("unable to update player")
}

func (p *PlayerAccessorMockFailure) GetPlayerForUserIfExists(userId string) (*model.Player, error) {
	return nil, errors.New("unable to get player")
}

func (p *PlayerAccessorMockFailure) CreatePlayerForUser(userId string) (*model.Player, error) {
	return nil, errors.New("unable to create player")
}

type PlayerAccessorMockConfigurable struct {
	GetPlayerUsingTransactionInternal              func(playerId string, transaction DatabaseTransaction) (*model.Player, error)
	CreatePlayerTransaction                        func() (*model.Player, error)
	UpdatePlayerWithUserIdUsingTransactionInternal func(playerId, userId string, transaction DatabaseTransaction) (*model.Player, error)
	GetPlayerForUserIfExistsInternal               func(userId string) (*model.Player, error)
	CreatePlayerForUserInternal                    func(userId string) (*model.Player, error)
}

func (p *PlayerAccessorMockConfigurable) GetPlayerUsingTransaction(playerId string, transaction DatabaseTransaction) (*model.Player, error) {
	return p.GetPlayerUsingTransactionInternal(playerId, transaction)
}

func (p *PlayerAccessorMockConfigurable) CreatePlayer() (*model.Player, error) {
	return p.CreatePlayerTransaction()
}

func (p *PlayerAccessorMockConfigurable) UpdatePlayerWithUserIdUsingTransaction(playerId, userId string, transaction DatabaseTransaction) (*model.Player, error) {
	return p.UpdatePlayerWithUserIdUsingTransactionInternal(playerId, userId, transaction)
}

func (p *PlayerAccessorMockConfigurable) GetPlayerForUserIfExists(userId string) (*model.Player, error) {
	return p.GetPlayerForUserIfExistsInternal(userId)
}

func (p *PlayerAccessorMockConfigurable) CreatePlayerForUser(userId string) (*model.Player, error) {
	return p.CreatePlayerForUserInternal(userId)
}
