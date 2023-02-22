package storage

import (
	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/model"
)

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

type GameJoinerMockSuccess struct {
	GameAccessor
}

func (g *GameJoinerMockSuccess) JoinGame(string, string) error {
	return nil
}

type GameJoinerMockFailure struct {
	GameAccessor
}

func (g *GameJoinerMockFailure) JoinGame(string, string) error {
	return errors.New("unable to join game")
}

type GameGetterMockSuccess struct {
	GameAccessor
	Game *model.Game
}

func (g *GameGetterMockSuccess) GetGame(string) (*model.Game, error) {
	return g.Game, nil
}

type GameGetterMockFailure struct {
	GameAccessor
}

func (g *GameGetterMockFailure) GetGame(string) (*model.Game, error) {
	return nil, errors.New("unable to get game")
}
