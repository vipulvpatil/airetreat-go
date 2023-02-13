package workers

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/vipulvpatil/airetreat-go/internal/clients/instagram"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
)

const GET_INSTAGRAM_LONG_LIVED_ACCESS_TOKEN = "get_instagram_long_lived_access_token"

var workerStorage storage.StorageAccessor
var instagramClient instagram.InstagramClient

type PoolDependencies struct {
	Namespace       string
	RedisPool       *redis.Pool
	Storage         storage.StorageAccessor
	InstagramClient instagram.InstagramClient
}

func NewPool(deps PoolDependencies) *work.WorkerPool {
	pool := work.NewWorkerPool(jobContext{}, 10, deps.Namespace, deps.RedisPool)

	pool.Job(GET_INSTAGRAM_LONG_LIVED_ACCESS_TOKEN, (*jobContext).getInstagramLongLivedAccessToken)

	// TODO: Not sure if this is the best way to do this. But using Package variables for all dependencies required inside any of the jobs.
	workerStorage = deps.Storage
	instagramClient = deps.InstagramClient
	return pool
}
