package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BotType(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput botType
	}{
		{
			name:           "creates AI account type",
			input:          "AI",
			expectedOutput: ai,
		},
		{
			name:           "creates HUMAN account type",
			input:          "HUMAN",
			expectedOutput: human,
		},
		{
			name:           "handles unknown account type",
			input:          "unknown",
			expectedOutput: undefinedBotType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := BotType(tt.input)
			assert.Equal(t, state, tt.expectedOutput)
		})
	}
}

func Test_BotType_String(t *testing.T) {
	tests := []struct {
		name           string
		input          botType
		expectedOutput string
	}{
		{
			name:           "gets AI from ai game state",
			input:          ai,
			expectedOutput: "AI",
		},
		{
			name:           "gets HUMAN from human game state",
			input:          human,
			expectedOutput: "HUMAN",
		},
		{
			name:           "gets unknown from undefinedBotType game state",
			input:          undefinedBotType,
			expectedOutput: "UNDEFINED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			botTypeString := tt.input.String()
			assert.Equal(t, botTypeString, tt.expectedOutput)
		})
	}
}

func Test_BotType_Valid(t *testing.T) {
	t.Run("returns true for a valid account type", func(t *testing.T) {
		assert.True(t, ai.Valid())
	})

	t.Run("returns false for a invalid account type", func(t *testing.T) {
		assert.False(t, undefinedBotType.Valid())
	})
}
