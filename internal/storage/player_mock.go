package storage

import "errors"

type PlayerAccessorMockSuccess struct {
	PlayerId string
}

func (p *PlayerAccessorMockSuccess) CreatePlayer(userId *string) (string, error) {
	return p.PlayerId, nil
}

func (p *PlayerAccessorMockSuccess) UpdatePlayerWithUserIdUsingTransaction(playerId, userId string, transaction DatabaseTransaction) error {
	return nil
}

type PlayerAccessorMockFailure struct {
}

func (p *PlayerAccessorMockFailure) CreatePlayer(userId *string) (string, error) {
	return "", errors.New("unable to create player")
}

func (p *PlayerAccessorMockFailure) UpdatePlayerWithUserIdUsingTransaction(playerId, userId string, transaction DatabaseTransaction) error {
	return errors.New("unable to update player")
}
