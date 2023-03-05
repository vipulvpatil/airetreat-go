package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gocraft/work"
	"github.com/vipulvpatil/airetreat-go/internal/workers"
)

func (s *AiRetreatGoService) GameHandlerLoop(ctx context.Context, tickerDuration time.Duration, wg *sync.WaitGroup, jobStarter workers.JobStarter) {
	wg.Add(1)
	defer wg.Done()

	ticker := time.NewTicker(tickerDuration)
	for {
		select {
		case <-ticker.C:
			gameIds := s.storage.GetUnhandledGameIdsForState("PLAYERS_JOINED")
			for _, gameId := range gameIds {
				_, err := jobStarter.EnqueueUnique(workers.START_GAME_ONCE_PLAYERS_HAVE_JOINED, work.Q{"gameId": gameId})
				if err != nil {
					fmt.Println(err)
				}
			}

		case <-ctx.Done():
			return
		}
	}
}
