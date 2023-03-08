package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GameState(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput gameState
	}{
		{
			name:           "creates STARTED account type",
			input:          "STARTED",
			expectedOutput: started,
		},
		{
			name:           "creates PLAYERS_JOINED account type",
			input:          "PLAYERS_JOINED",
			expectedOutput: playersJoined,
		},
		{
			name:           "creates WAITING_FOR_BOT_QUESTION account type",
			input:          "WAITING_FOR_BOT_QUESTION",
			expectedOutput: waitingForBotQuestion,
		},
		{
			name:           "creates WAITING_FOR_BOT_ANSWER account type",
			input:          "WAITING_FOR_BOT_ANSWER",
			expectedOutput: waitingForBotAnswer,
		},
		{
			name:           "creates WAITING_FOR_PLAYER_QUESTION account type",
			input:          "WAITING_FOR_PLAYER_QUESTION",
			expectedOutput: waitingForPlayerQuestion,
		},
		{
			name:           "creates WAITING_FOR_PLAYER_ANSWER account type",
			input:          "WAITING_FOR_PLAYER_ANSWER",
			expectedOutput: waitingForPlayerAnswer,
		},
		{
			name:           "creates FINISHED account type",
			input:          "FINISHED",
			expectedOutput: finished,
		},
		{
			name:           "handles unknown account type",
			input:          "unknown",
			expectedOutput: undefinedGameState,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := GameState(tt.input)
			assert.Equal(t, state, tt.expectedOutput)
		})
	}
}

func Test_GameState_String(t *testing.T) {
	tests := []struct {
		name           string
		input          gameState
		expectedOutput string
	}{
		{
			name:           "gets STARTED from started game state",
			input:          started,
			expectedOutput: "STARTED",
		},
		{
			name:           "gets PLAYERS_JOINED from playersJoined game state",
			input:          playersJoined,
			expectedOutput: "PLAYERS_JOINED",
		},
		{
			name:           "gets WAITING_FOR_BOT_QUESTION from waitingForBotQuestion game state",
			input:          waitingForBotQuestion,
			expectedOutput: "WAITING_FOR_BOT_QUESTION",
		},
		{
			name:           "gets WAITING_FOR_BOT_ANSWER from waitingForBotAnswer game state",
			input:          waitingForBotAnswer,
			expectedOutput: "WAITING_FOR_BOT_ANSWER",
		},
		{
			name:           "gets WAITING_FOR_PLAYER_QUESTION from waitingForPlayerQuestion game state",
			input:          waitingForPlayerQuestion,
			expectedOutput: "WAITING_FOR_PLAYER_QUESTION",
		},
		{
			name:           "gets WAITING_FOR_PLAYER_ANSWER from waitingForPlayerAnswer game state",
			input:          waitingForPlayerAnswer,
			expectedOutput: "WAITING_FOR_PLAYER_ANSWER",
		},
		{
			name:           "gets FINISHED from finished game state",
			input:          finished,
			expectedOutput: "FINISHED",
		},
		{
			name:           "gets unknown from undefinedGameState game state",
			input:          undefinedGameState,
			expectedOutput: "UNDEFINED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameStateString := tt.input.String()
			assert.Equal(t, gameStateString, tt.expectedOutput)
		})
	}
}

func Test_GameState_Valid(t *testing.T) {
	t.Run("returns true for a valid account type", func(t *testing.T) {
		assert.True(t, started.Valid())
	})

	t.Run("returns false for a invalid account type", func(t *testing.T) {
		assert.False(t, undefinedGameState.Valid())
	})
}

func Test_GameState_isWaitingForAi(t *testing.T) {
	t.Run("returns true for a waitingForBotQuestion", func(t *testing.T) {
		assert.True(t, waitingForBotQuestion.isWaitingForAi())
	})

	t.Run("returns true for a waitingForBotAnswer", func(t *testing.T) {
		assert.True(t, waitingForBotAnswer.isWaitingForAi())
	})

	t.Run("returns false for other states", func(t *testing.T) {
		assert.False(t, waitingForPlayerQuestion.isWaitingForAi())
	})
}

func Test_GameState_isWaitingForHuman(t *testing.T) {
	t.Run("returns true for a waitingForPlayerQuestion", func(t *testing.T) {
		assert.True(t, waitingForPlayerQuestion.isWaitingForHuman())
	})

	t.Run("returns true for a waitingForPlayerAnswer", func(t *testing.T) {
		assert.True(t, waitingForPlayerAnswer.isWaitingForHuman())
	})

	t.Run("returns false for other states", func(t *testing.T) {
		assert.False(t, waitingForBotQuestion.isWaitingForHuman())
	})
}

func Test_GameState_isQuestion(t *testing.T) {
	t.Run("returns true for a waitingForBotQuestion", func(t *testing.T) {
		assert.True(t, waitingForBotQuestion.isQuestion())
	})

	t.Run("returns true for a waitingForPlayerQuestion", func(t *testing.T) {
		assert.True(t, waitingForPlayerQuestion.isQuestion())
	})

	t.Run("returns false for other states", func(t *testing.T) {
		assert.False(t, waitingForBotAnswer.isQuestion())
	})
}
func Test_GameState_isAnswer(t *testing.T) {
	t.Run("returns true for a waitingForBotAnswer", func(t *testing.T) {
		assert.True(t, waitingForBotAnswer.isAnswer())
	})

	t.Run("returns true for a waitingForPlayerAnswer", func(t *testing.T) {
		assert.True(t, waitingForPlayerAnswer.isAnswer())
	})

	t.Run("returns false for other states", func(t *testing.T) {
		assert.False(t, waitingForBotQuestion.isAnswer())
	})
}

func Test_GameState_isWaitingForMessage(t *testing.T) {
	t.Run("returns true for a waitingForBotQuestion", func(t *testing.T) {
		assert.True(t, waitingForBotQuestion.isWaitingForMessage())
	})

	t.Run("returns true for a waitingForPlayerQuestion", func(t *testing.T) {
		assert.True(t, waitingForPlayerQuestion.isWaitingForMessage())
	})

	t.Run("returns true for a waitingForBotAnswer", func(t *testing.T) {
		assert.True(t, waitingForBotAnswer.isWaitingForMessage())
	})

	t.Run("returns true for a waitingForPlayerAnswer", func(t *testing.T) {
		assert.True(t, waitingForPlayerAnswer.isWaitingForMessage())
	})

	t.Run("returns false for other states", func(t *testing.T) {
		assert.False(t, playersJoined.isWaitingForMessage())
	})
}
