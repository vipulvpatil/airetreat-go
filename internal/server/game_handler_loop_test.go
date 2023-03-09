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

func Test_GameHandlerLoop(t *testing.T) {
	t.Run("looks for appropriate games and calls job starter, until canceled", func(t *testing.T) {
		jobStarterMock := &workers.JobStarterMockCallCheck{}
		tickerDuration := 10 * time.Millisecond
		gamesAccessorGetUnhandledGameIdsMockCaller := GetUnhandledGameIdsMockCaller{
			MapByInput: map[string]*functionCallInspectableMock{
				"PLAYERS_JOINED": {
					ReturnData:  [][]string{{"game_id1", "game_id2"}},
					ReturnCount: 1,
				},
			},
		}
		gamesAccessorGetOldGamesMockCaller := GetOldGamesMockCaller{
			&functionCallInspectableMock{
				ReturnData:  [][]string{{"old_game_id1"}, {"old_game_id2"}},
				ReturnCount: 2,
			},
		}

		functionsToCheck := []struct {
			name              string
			functionCall      functionCallInspectable
			expectedCallCount int
		}{
			{
				name:              "GetUnhandledGameIds for state PLAYERS_JOINED, %s",
				functionCall:      gamesAccessorGetUnhandledGameIdsMockCaller.MapByInput["PLAYERS_JOINED"],
				expectedCallCount: 4,
			},
			{
				name:              "GetOldGames, %s",
				functionCall:      gamesAccessorGetOldGamesMockCaller,
				expectedCallCount: 4,
			},
		}

		jobStartedCallsToVerify := []struct {
			jobName string
			jobArgs []map[string]any
		}{
			{
				jobName: workers.START_GAME_ONCE_PLAYERS_HAVE_JOINED,
				jobArgs: []map[string]any{
					{"gameId": "game_id1"},
					{"gameId": "game_id2"},
				},
			},
			{
				jobName: workers.DELETE_EXPIRED_GAMES,
				jobArgs: []map[string]any{
					{"gameId": "old_game_id1"},
					{"gameId": "old_game_id2"},
				},
			},
		}

		server, _ := NewServer(
			ServerDependencies{
				Storage: storage.NewStorageAccessorMock(
					storage.WithGameAccessorMock(
						&storage.GameAccessorConfigurableMock{
							GetUnhandledGameIdsForStateInternal: gamesAccessorGetUnhandledGameIdsMockCaller.getUnhandledGameIdsForStateInternal,
							GetOldGamesInternal:                 gamesAccessorGetOldGamesMockCaller.getOldGames,
						},
					),
				),
			},
		)

		var wg sync.WaitGroup
		gameHandlerLoopCtx, cancelGameHandlerLoop := context.WithCancel(context.Background())
		go server.GameHandlerLoop(gameHandlerLoopCtx, tickerDuration, &wg, jobStarterMock)
		time.Sleep(45 * time.Millisecond)

		for _, jobsStarted := range jobStartedCallsToVerify {
			assertJobStarterCalledWithArgsForJob(
				t,
				jobsStarted.jobArgs,
				jobStarterMock,
				jobsStarted.jobName,
			)
		}

		for _, f := range functionsToCheck {
			assertCallCount(t, f.expectedCallCount, f.functionCall, f.name, "loop should run continuously until canceled")
		}
		cancelGameHandlerLoop()
		time.Sleep(45 * time.Millisecond)
		for _, f := range functionsToCheck {
			assertCallCount(t, f.expectedCallCount, f.functionCall, f.name, "function call count should not change once loop is canceled")
		}
	})
}

type functionCallInspectable interface {
	FunctionCalledCount() int
}

type functionCallInspectableMock struct {
	ReturnData  [][]string
	ReturnCount int
	callCount   int
}

func (f *functionCallInspectableMock) FunctionCalledCount() int {
	return f.callCount
}

func assertCallCount(t *testing.T, expectedCallCount int, functionCall functionCallInspectable, msgAndArgs ...any) bool {
	return assert.Equal(t, expectedCallCount, functionCall.FunctionCalledCount(), msgAndArgs...)
}

type GetUnhandledGameIdsMockCaller struct {
	MapByInput map[string]*functionCallInspectableMock
}

func (m *GetUnhandledGameIdsMockCaller) getUnhandledGameIdsForStateInternal(gameStateString string) ([]string, error) {
	f, ok := m.MapByInput[gameStateString]
	if !ok {
		f = &functionCallInspectableMock{}
	}
	f.callCount++
	if f.ReturnCount >= f.callCount {
		m.MapByInput[gameStateString] = f
		return f.ReturnData[f.callCount-1], nil
	}
	m.MapByInput[gameStateString] = f
	return nil, nil
}

type GetOldGamesMockCaller struct {
	*functionCallInspectableMock
}

func (m *GetOldGamesMockCaller) getOldGames(gameExpiryDuration time.Duration) ([]string, error) {
	m.callCount++
	if m.ReturnCount >= m.callCount {
		return m.ReturnData[m.callCount-1], nil
	}
	return nil, nil
}

func assertJobStarterCalledWithArgsForJob(t *testing.T, expectedCalledArgs []map[string]any, jobStarter *workers.JobStarterMockCallCheck, jobName string) bool {
	return assert.EqualValues(
		t,
		expectedCalledArgs,
		jobStarter.CalledArgs[jobName],
		"appropriate jobs should be started from the loop",
	)
}
