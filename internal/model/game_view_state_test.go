package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GameViewState(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput gameViewState
	}{
		{
			name:           "creates WAITING_FOR_PLAYERS_TO_JOIN account type",
			input:          "WAITING_FOR_PLAYERS_TO_JOIN",
			expectedOutput: waitingForPlayersToJoin,
		},
		{
			name:           "creates WAITING_ON_BOT_TO_ASK_A_QUESTION account type",
			input:          "WAITING_ON_BOT_TO_ASK_A_QUESTION",
			expectedOutput: waitingOnBotToAskAQuestion,
		},
		{
			name:           "creates WAITING_ON_BOT_TO_ANSWER account type",
			input:          "WAITING_ON_BOT_TO_ANSWER",
			expectedOutput: waitingOnBotToAnswer,
		},
		{
			name:           "creates WAITING_ON_YOU_TO_ASK_A_QUESTION account type",
			input:          "WAITING_ON_YOU_TO_ASK_A_QUESTION",
			expectedOutput: waitingOnYouToAskAQuestion,
		},
		{
			name:           "creates WAITING_ON_YOU_TO_ANSWER account type",
			input:          "WAITING_ON_YOU_TO_ANSWER",
			expectedOutput: waitingOnYouToAnswer,
		},
		{
			name:           "creates YOU_LOST account type",
			input:          "YOU_LOST",
			expectedOutput: youLost,
		},
		{
			name:           "creates YOU_WON account type",
			input:          "YOU_WON",
			expectedOutput: youWon,
		},
		{
			name:           "creates TIME_UP account type",
			input:          "TIME_UP",
			expectedOutput: timeUp,
		},
		{
			name:           "handles unknown account type",
			input:          "unknown",
			expectedOutput: undefinedGameViewState,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := GameViewState(tt.input)
			assert.Equal(t, state, tt.expectedOutput)
		})
	}
}

func Test_GameViewState_String(t *testing.T) {
	tests := []struct {
		name           string
		input          gameViewState
		expectedOutput string
	}{
		{
			name:           "gets WAITING_FOR_PLAYERS_TO_JOIN from started game state",
			input:          waitingForPlayersToJoin,
			expectedOutput: "WAITING_FOR_PLAYERS_TO_JOIN",
		},
		{
			name:           "gets WAITING_ON_BOT_TO_ASK_A_QUESTION from playersJoined game state",
			input:          waitingOnBotToAskAQuestion,
			expectedOutput: "WAITING_ON_BOT_TO_ASK_A_QUESTION",
		},
		{
			name:           "gets WAITING_ON_BOT_TO_ANSWER from waitingForBotQuestion game state",
			input:          waitingOnBotToAnswer,
			expectedOutput: "WAITING_ON_BOT_TO_ANSWER",
		},
		{
			name:           "gets WAITING_ON_YOU_TO_ASK_A_QUESTION from waitingForBotAnswer game state",
			input:          waitingOnYouToAskAQuestion,
			expectedOutput: "WAITING_ON_YOU_TO_ASK_A_QUESTION",
		},
		{
			name:           "gets WAITING_ON_YOU_TO_ANSWER from waitingForPlayerQuestion game state",
			input:          waitingOnYouToAnswer,
			expectedOutput: "WAITING_ON_YOU_TO_ANSWER",
		},
		{
			name:           "gets YOU_LOST from waitingForPlayerAnswer game state",
			input:          youLost,
			expectedOutput: "YOU_LOST",
		},
		{
			name:           "gets YOU_WON from finished game state",
			input:          youWon,
			expectedOutput: "YOU_WON",
		},
		{
			name:           "gets TIME_UP from finished game state",
			input:          timeUp,
			expectedOutput: "TIME_UP",
		},
		{
			name:           "gets unknown from undefinedGameViewState game state",
			input:          undefinedGameViewState,
			expectedOutput: "UNDEFINED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameViewStateString := tt.input.String()
			assert.Equal(t, gameViewStateString, tt.expectedOutput)
		})
	}
}

func Test_GameViewState_Valid(t *testing.T) {
	t.Run("returns true for a valid account type", func(t *testing.T) {
		assert.True(t, started.Valid())
	})

	t.Run("returns false for a invalid account type", func(t *testing.T) {
		assert.False(t, undefinedGameViewState.Valid())
	})
}
