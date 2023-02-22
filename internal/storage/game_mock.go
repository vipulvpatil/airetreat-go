package storage

import "github.com/pkg/errors"

type GameCreatorMockSuccess struct {
	GameAccessor
	GameId string
}

func (g *GameCreatorMockSuccess) CreateGame() (string, error) {
	return g.GameId, nil
}

type GameCreatorMockFailure struct {
	GameAccessor
}

func (g *GameCreatorMockFailure) CreateGame() (string, error) {
	return "", errors.New("unable to create game")
}
