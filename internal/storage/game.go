package storage

import "github.com/vipulvpatil/airetreat-go/internal/model"

type GameAccessor interface {
	CreateGame() (string, error)
	JoinGame(gameId, playerId string) error
	GetGame(gameId string) (*model.Game, error)
	GetGames(playerId string) ([]string, error)
}
