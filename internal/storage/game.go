package storage

import "github.com/vipulvpatil/airetreat-go/internal/model"

type GameAccessor interface {
	CreateGame() model.Game
}
