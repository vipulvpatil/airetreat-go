package server

import (
	"context"
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
			s.beginGames(jobStarter)
			s.askQuestionsUsingAi(jobStarter)
			s.answerQuestionsUsingAi(jobStarter)
			s.deleteExpiredGames(jobStarter)
		case <-ctx.Done():
			return
		}
	}
}

func (s *AiRetreatGoService) beginGames(jobStarter workers.JobStarter) {
	gameIds, err := s.storage.GetUnhandledGameIdsForState("PLAYERS_JOINED")
	if err != nil {
		s.logger.LogMessageln(err)
		return
	}
	for _, gameId := range gameIds {
		_, err := jobStarter.EnqueueUnique(workers.START_GAME_ONCE_PLAYERS_HAVE_JOINED, work.Q{"gameId": gameId})
		if err != nil {
			s.logger.LogMessageln(err)
		}
	}
}

func (s *AiRetreatGoService) askQuestionsUsingAi(jobStarter workers.JobStarter) {
	gameIds, err := s.storage.GetUnhandledGameIdsForState("WAITING_FOR_AI_QUESTION")
	if err != nil {
		s.logger.LogMessageln(err)
		return
	}
	for _, gameId := range gameIds {
		_, err := jobStarter.EnqueueUnique(workers.ASK_QUESTION_ON_BEHALF_OF_BOT, work.Q{"gameId": gameId})
		if err != nil {
			s.logger.LogMessageln(err)
		}
	}
}

func (s *AiRetreatGoService) answerQuestionsUsingAi(jobStarter workers.JobStarter) {
	gameIds, err := s.storage.GetUnhandledGameIdsForState("WAITING_FOR_AI_ANSWER")
	if err != nil {
		s.logger.LogMessageln(err)
		return
	}
	for _, gameId := range gameIds {
		_, err := jobStarter.EnqueueUnique(workers.ANSWER_QUESTION_ON_BEHALF_OF_BOT, work.Q{"gameId": gameId})
		if err != nil {
			s.logger.LogMessageln(err)
		}
	}
}

func (s *AiRetreatGoService) deleteExpiredGames(jobStarter workers.JobStarter) {
	gameIds, err := s.storage.GetOldGames(-1 * time.Hour)
	if err != nil {
		s.logger.LogMessageln(err)
		return
	}
	for _, gameId := range gameIds {
		_, err := jobStarter.EnqueueUnique(workers.DELETE_EXPIRED_GAMES, work.Q{"gameId": gameId})
		if err != nil {
			s.logger.LogMessageln(err)
		}
	}
}
