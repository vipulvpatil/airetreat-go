package storage

import (
	"time"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/model"
)

type GameCreatorMockSuccess struct {
	GameAccessor
	GameId string
}

func (g *GameCreatorMockSuccess) CreateGame(public bool) (string, error) {
	return g.GameId, nil
}

type GameCreatorMockFailure struct {
	GameAccessor
}

func (g *GameCreatorMockFailure) CreateGame(public bool) (string, error) {
	return "", errors.New("unable to create game")
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

type GamesGetterMockSuccess struct {
	GameAccessor
	GameIds []string
}

func (g *GamesGetterMockSuccess) GetGames(string) ([]string, error) {
	return g.GameIds, nil
}

type GamesGetterMockFailure struct {
	GameAccessor
}

func (g *GamesGetterMockFailure) GetGames(string) ([]string, error) {
	return nil, errors.New("unable to get games")
}

type GameIdsGetterMockNil struct {
	GameAccessor
}

func (g *GameIdsGetterMockNil) GetUnhandledGameIdsForState(gameStateString string) ([]string, error) {
	return nil, nil
}

type GameIdsGetterMockEmpty struct {
	GameAccessor
}

func (g *GameIdsGetterMockEmpty) GetUnhandledGameIdsForState(gameStateString string) ([]string, error) {
	return []string{}, nil
}

type GameAccessorConfigurableMock struct {
	CreateGameInternal                                               func() (string, error)
	GetGameInternal                                                  func(gameId string) (*model.Game, error)
	GetGameUsingTransactionInternal                                  func(gameId string, transaction DatabaseTransaction) (*model.Game, error)
	GetGamesInternal                                                 func(playerId string) ([]string, error)
	UpdateGameStateInternal                                          func(gameId string, updateOpts GameUpdateOptions) error
	UpdateGameStateUsingTransactionInternal                          func(gameId string, updateOpts GameUpdateOptions, transaction DatabaseTransaction) error
	GetUnhandledGameIdsForStateInternal                              func(gameStateString string) ([]string, error)
	DeleteGameInternal                                               func(gameId string) error
	GetOldGamesInternal                                              func(gameExpiryDuration time.Duration) ([]string, error)
	UpdateGameStateIfEnoughPlayersHaveJoinedUsingTransactionInternal func(gameId string, transaction DatabaseTransaction) error
	GetAutoJoinableGamesInternal                                     func() ([]string, error)
}

func (g *GameAccessorConfigurableMock) CreateGame(public bool) (string, error) {
	return g.CreateGameInternal()
}
func (g *GameAccessorConfigurableMock) GetGame(gameId string) (*model.Game, error) {
	return g.GetGameInternal(gameId)
}
func (g *GameAccessorConfigurableMock) GetGameUsingTransaction(gameId string, transaction DatabaseTransaction) (*model.Game, error) {
	return g.GetGameUsingTransactionInternal(gameId, transaction)
}
func (g *GameAccessorConfigurableMock) GetGames(playerId string) ([]string, error) {
	return g.GetGamesInternal(playerId)
}
func (g *GameAccessorConfigurableMock) UpdateGameState(gameId string, updateOpts GameUpdateOptions) error {
	return g.UpdateGameStateInternal(gameId, updateOpts)
}
func (g *GameAccessorConfigurableMock) UpdateGameStateUsingTransaction(gameId string, updateOpts GameUpdateOptions, transaction DatabaseTransaction) error {
	return g.UpdateGameStateUsingTransactionInternal(gameId, updateOpts, transaction)
}
func (g *GameAccessorConfigurableMock) GetUnhandledGameIdsForState(gameStateString string) ([]string, error) {
	return g.GetUnhandledGameIdsForStateInternal(gameStateString)
}
func (g *GameAccessorConfigurableMock) DeleteGame(gameId string) error {
	return g.DeleteGameInternal(gameId)
}
func (g *GameAccessorConfigurableMock) GetOldGames(gameExpiryDuration time.Duration) ([]string, error) {
	return g.GetOldGamesInternal(gameExpiryDuration)
}
func (g *GameAccessorConfigurableMock) UpdateGameStateIfEnoughPlayersHaveJoinedUsingTransaction(gameId string, transaction DatabaseTransaction) error {
	return g.UpdateGameStateIfEnoughPlayersHaveJoinedUsingTransactionInternal(gameId, transaction)
}
func (g *GameAccessorConfigurableMock) GetAutoJoinableGames() ([]string, error) {
	return g.GetAutoJoinableGamesInternal()
}
