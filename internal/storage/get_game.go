package storage

import "github.com/vipulvpatil/airetreat-go/internal/model"

func (s *Storage) GetGame(gameId string) (*model.Game, error) {
	return getGame(s.db, gameId)
}
