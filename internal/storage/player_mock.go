package storage

import (
	"errors"

	"github.com/vipulvpatil/airetreat-go/internal/model"
)

type PlayerAccessorMockSuccess struct {
	PlayerId string
	UserId   *string
}

func (p *PlayerAccessorMockSuccess) GetPlayer(playerId string) (*model.Player, error) {
	return model.NewPlayer(model.PlayerOptions{
		Id:     p.PlayerId,
		UserId: p.UserId,
	})
}

func (p *PlayerAccessorMockSuccess) GetPlayerUsingTransaction(playerId string, transaction DatabaseTransaction) (*model.Player, error) {
	return model.NewPlayer(model.PlayerOptions{
		Id:     p.PlayerId,
		UserId: p.UserId,
	})
}

func (p *PlayerAccessorMockSuccess) CreatePlayer() (*model.Player, error) {
	return model.NewPlayer(model.PlayerOptions{
		Id:     p.PlayerId,
		UserId: p.UserId,
	})
}

func (p *PlayerAccessorMockSuccess) UpdatePlayerWithUserIdUsingTransaction(playerId, userId string, transaction DatabaseTransaction) (*model.Player, error) {
	return model.NewPlayer(model.PlayerOptions{
		Id:     p.PlayerId,
		UserId: p.UserId,
	})
}

func (p *PlayerAccessorMockSuccess) GetPlayerForUserOrNil(userId string) (*model.Player, error) {
	return model.NewPlayer(model.PlayerOptions{
		Id:     p.PlayerId,
		UserId: p.UserId,
	})
}

func (p *PlayerAccessorMockSuccess) CreatePlayerForUser(userId string) (*model.Player, error) {
	return model.NewPlayer(model.PlayerOptions{
		Id:     p.PlayerId,
		UserId: p.UserId,
	})
}

type PlayerAccessorMockFailure struct {
}

func (p *PlayerAccessorMockFailure) GetPlayer(playerId string) (*model.Player, error) {
	return nil, errors.New("unable to get player")
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

func (p *PlayerAccessorMockFailure) GetPlayerForUserOrNil(userId string) (*model.Player, error) {
	return nil, errors.New("unable to get player")
}

func (p *PlayerAccessorMockFailure) CreatePlayerForUser(userId string) (*model.Player, error) {
	return nil, errors.New("unable to create player")
}

type PlayerAccessorMockConfigurable struct {
	GetPlayerInternal                              func(playerId string) (*model.Player, error)
	GetPlayerUsingTransactionInternal              func(playerId string, transaction DatabaseTransaction) (*model.Player, error)
	CreatePlayerInternal                           func() (*model.Player, error)
	UpdatePlayerWithUserIdUsingTransactionInternal func(playerId, userId string, transaction DatabaseTransaction) (*model.Player, error)
	GetPlayerForUserOrNilInternal                  func(userId string) (*model.Player, error)
	CreatePlayerForUserInternal                    func(userId string) (*model.Player, error)
}

func (p *PlayerAccessorMockConfigurable) GetPlayer(playerId string) (*model.Player, error) {
	return p.GetPlayerInternal(playerId)
}

func (p *PlayerAccessorMockConfigurable) GetPlayerUsingTransaction(playerId string, transaction DatabaseTransaction) (*model.Player, error) {
	return p.GetPlayerUsingTransactionInternal(playerId, transaction)
}

func (p *PlayerAccessorMockConfigurable) CreatePlayer() (*model.Player, error) {
	return p.CreatePlayerInternal()
}

func (p *PlayerAccessorMockConfigurable) UpdatePlayerWithUserIdUsingTransaction(playerId, userId string, transaction DatabaseTransaction) (*model.Player, error) {
	return p.UpdatePlayerWithUserIdUsingTransactionInternal(playerId, userId, transaction)
}

func (p *PlayerAccessorMockConfigurable) GetPlayerForUserOrNil(userId string) (*model.Player, error) {
	return p.GetPlayerForUserOrNilInternal(userId)
}

func (p *PlayerAccessorMockConfigurable) CreatePlayerForUser(userId string) (*model.Player, error) {
	return p.CreatePlayerForUserInternal(userId)
}
