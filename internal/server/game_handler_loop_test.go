package server

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
	"github.com/vipulvpatil/airetreat-go/internal/workers"
)

type GetUnhandledGameIdsMockCaller struct {
	MapByInput map[string]struct {
		ReturnData  [][]string
		CallCount   int
		ReturnCount int
	}
}

func (m *GetUnhandledGameIdsMockCaller) getUnhandledGameIdsForStateInternal(gameStateString string) []string {
	f, ok := m.MapByInput[gameStateString]
	if !ok {
		f = struct {
			ReturnData  [][]string
			CallCount   int
			ReturnCount int
		}{}
	}
	f.CallCount++
	if f.ReturnCount >= f.CallCount {
		m.MapByInput[gameStateString] = f
		return f.ReturnData[f.CallCount-1]
	}
	m.MapByInput[gameStateString] = f
	return nil
}

func Test_GameHandlerLoop(t *testing.T) {
	tests := []struct {
		name                                       string
		jobStarterMock                             *workers.JobStarterMockCallCheck
		gamesAccessorGetUnhandledGameIdsMockCaller GetUnhandledGameIdsMockCaller
		tickerDuration                             time.Duration
	}{
		{
			name:           "looks for appropriate games and calls job starter, until canceled",
			jobStarterMock: &workers.JobStarterMockCallCheck{},
			gamesAccessorGetUnhandledGameIdsMockCaller: GetUnhandledGameIdsMockCaller{
				MapByInput: map[string]struct {
					ReturnData  [][]string
					CallCount   int
					ReturnCount int
				}{
					"PLAYERS_JOINED": {
						ReturnData:  [][]string{{"game_id1", "game_id2"}},
						CallCount:   0,
						ReturnCount: 1,
					},
				},
			},
			tickerDuration: 10 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := NewServer(
				ServerDependencies{
					Storage: storage.NewStorageAccessorMock(
						storage.WithGameAccessorMock(
							&storage.GameAccessorConfigurableMock{
								GetUnhandledGameIdsForStateInternal: tt.gamesAccessorGetUnhandledGameIdsMockCaller.getUnhandledGameIdsForStateInternal,
							},
						),
					),
				},
			)

			var wg sync.WaitGroup
			gameHandlerLoopCtx, cancelGameHandlerLoop := context.WithCancel(context.Background())
			go server.GameHandlerLoop(gameHandlerLoopCtx, tt.tickerDuration, &wg, tt.jobStarterMock)
			time.Sleep(45 * time.Millisecond)
			assert.EqualValues(
				t,
				[]map[string]interface{}{
					{"gameId": "game_id1"},
					{"gameId": "game_id2"},
				},
				tt.jobStarterMock.CalledArgs[workers.START_GAME_ONCE_PLAYERS_HAVE_JOINED],
				"appropriate jobs should be started from the loop",
			)
			assert.Equal(t, 4, tt.gamesAccessorGetUnhandledGameIdsMockCaller.MapByInput["PLAYERS_JOINED"].CallCount, "loop should run continuously until canceled")
			cancelGameHandlerLoop()
			time.Sleep(45 * time.Millisecond)
			assert.Equal(t, 4, tt.gamesAccessorGetUnhandledGameIdsMockCaller.MapByInput["PLAYERS_JOINED"].CallCount, "loop should not run once canceled")
		})
	}
}
