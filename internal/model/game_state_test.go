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
			name:           "creates WAITING_FOR_AI_QUESTION account type",
			input:          "WAITING_FOR_AI_QUESTION",
			expectedOutput: waitingForAiQuestion,
		},
		{
			name:           "creates WAITING_FOR_AI_ANSWER account type",
			input:          "WAITING_FOR_AI_ANSWER",
			expectedOutput: waitingForAiAnswer,
		},
		{
			name:           "creates WAITING_FOR_HUMAN_QUESTION account type",
			input:          "WAITING_FOR_HUMAN_QUESTION",
			expectedOutput: waitingForHumanQuestion,
		},
		{
			name:           "creates WAITING_FOR_HUMAN_ANSWER account type",
			input:          "WAITING_FOR_HUMAN_ANSWER",
			expectedOutput: waitingForHumanAnswer,
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
			name:           "gets WAITING_FOR_AI_QUESTION from waitingForAiQuestion game state",
			input:          waitingForAiQuestion,
			expectedOutput: "WAITING_FOR_AI_QUESTION",
		},
		{
			name:           "gets WAITING_FOR_AI_ANSWER from waitingForAiAnswer game state",
			input:          waitingForAiAnswer,
			expectedOutput: "WAITING_FOR_AI_ANSWER",
		},
		{
			name:           "gets WAITING_FOR_HUMAN_QUESTION from waitingForHumanQuestion game state",
			input:          waitingForHumanQuestion,
			expectedOutput: "WAITING_FOR_HUMAN_QUESTION",
		},
		{
			name:           "gets WAITING_FOR_HUMAN_ANSWER from waitingForHumanAnswer game state",
			input:          waitingForHumanAnswer,
			expectedOutput: "WAITING_FOR_HUMAN_ANSWER",
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

func Test_GameState_isWaitingOnAi(t *testing.T) {
	t.Run("returns true for a waitingForAiQuestion", func(t *testing.T) {
		assert.True(t, waitingForAiQuestion.isWaitingOnAi())
	})

	t.Run("returns true for a waitingForAiAnswer", func(t *testing.T) {
		assert.True(t, waitingForAiAnswer.isWaitingOnAi())
	})

	t.Run("returns false for other states", func(t *testing.T) {
		assert.False(t, waitingForHumanQuestion.isWaitingOnAi())
	})
}

func Test_GameState_isWaitingOnHuman(t *testing.T) {
	t.Run("returns true for a waitingForHumanQuestion", func(t *testing.T) {
		assert.True(t, waitingForHumanQuestion.isWaitingOnHuman())
	})

	t.Run("returns true for a waitingForHumanAnswer", func(t *testing.T) {
		assert.True(t, waitingForHumanAnswer.isWaitingOnHuman())
	})

	t.Run("returns false for other states", func(t *testing.T) {
		assert.False(t, waitingForAiQuestion.isWaitingOnHuman())
	})
}

func Test_GameState_isWaitingForAQuestion(t *testing.T) {
	t.Run("returns true for a waitingForAiQuestion", func(t *testing.T) {
		assert.True(t, waitingForAiQuestion.isWaitingForAQuestion())
	})

	t.Run("returns true for a waitingForHumanQuestion", func(t *testing.T) {
		assert.True(t, waitingForHumanQuestion.isWaitingForAQuestion())
	})

	t.Run("returns false for other states", func(t *testing.T) {
		assert.False(t, waitingForAiAnswer.isWaitingForAQuestion())
	})
}
func Test_GameState_isWaitingForAnAnswer(t *testing.T) {
	t.Run("returns true for a waitingForAiAnswer", func(t *testing.T) {
		assert.True(t, waitingForAiAnswer.isWaitingForAnAnswer())
	})

	t.Run("returns true for a waitingForHumanAnswer", func(t *testing.T) {
		assert.True(t, waitingForHumanAnswer.isWaitingForAnAnswer())
	})

	t.Run("returns false for other states", func(t *testing.T) {
		assert.False(t, waitingForAiQuestion.isWaitingForAnAnswer())
	})
}

func Test_GameState_isWaitingForMessage(t *testing.T) {
	t.Run("returns true for a waitingForAiQuestion", func(t *testing.T) {
		assert.True(t, waitingForAiQuestion.isWaitingForMessage())
	})

	t.Run("returns true for a waitingForHumanQuestion", func(t *testing.T) {
		assert.True(t, waitingForHumanQuestion.isWaitingForMessage())
	})

	t.Run("returns true for a waitingForAiAnswer", func(t *testing.T) {
		assert.True(t, waitingForAiAnswer.isWaitingForMessage())
	})

	t.Run("returns true for a waitingForHumanAnswer", func(t *testing.T) {
		assert.True(t, waitingForHumanAnswer.isWaitingForMessage())
	})

	t.Run("returns false for other states", func(t *testing.T) {
		assert.False(t, playersJoined.isWaitingForMessage())
	})
}
