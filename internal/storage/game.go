package storage

import (
	"time"

	"github.com/vipulvpatil/airetreat-go/internal/model"
)

type GameAccessor interface {
	CreateGame() (string, error)
	GetGame(gameId string) (*model.Game, error)
	GetGameUsingTransaction(gameId string, transaction DatabaseTransaction) (*model.Game, error)
	GetGames(playerId string) ([]string, error)
	UpdateGameState(gameId string, updateOpts GameUpdateOptions) error
	UpdateGameStateUsingTransaction(gameId string, updateOpts GameUpdateOptions, transaction DatabaseTransaction) error
	GetUnhandledGameIdsForState(gameStateString string) ([]string, error)
	DeleteGame(gameId string) error
	GetOldGames(gameExpiryDuration time.Duration) ([]string, error)
	UpdateGameStateIfEnoughPlayersHaveJoinedUsingTransaction(gameId string, transaction DatabaseTransaction) error
}
