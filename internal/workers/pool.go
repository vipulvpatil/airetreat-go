package workers

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
)

const START_GAME_ONCE_PLAYERS_HAVE_JOINED = "start_game_once_players_have_joined"

var workerStorage storage.StorageAccessor

type PoolDependencies struct {
	Namespace string
	RedisPool *redis.Pool
	Storage   storage.StorageAccessor
}

func NewPool(deps PoolDependencies) *work.WorkerPool {
	pool := work.NewWorkerPool(jobContext{}, 10, deps.Namespace, deps.RedisPool)

	pool.Job(START_GAME_ONCE_PLAYERS_HAVE_JOINED, (*jobContext).startGameOncePlayersHaveJoined)

	// TODO: Not sure if this is the best way to do this. But using Package variables for all dependencies required inside any of the jobs.
	workerStorage = deps.Storage
	return pool
}
