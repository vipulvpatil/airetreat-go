package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewGame(t *testing.T) {
	bot, _ := NewBot(BotOptions{Id: "botId1"})
	tests := []struct {
		name           string
		input          GameOptions
		expectedOutput *Game
		errorExpected  bool
		errorString    string
	}{
		{
			name:           "id is empty",
			input:          GameOptions{},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create game with an empty id",
		},
		{
			name: "state is invalid",
			input: GameOptions{
				Id: "123",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create game with an invalid state",
		},
		{
			name: "invalid turn order",
			input: GameOptions{
				Id:    "123",
				State: "STARTED",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create game with empty turn order array",
		},
		{
			name: "invalid bots",
			input: GameOptions{
				Id:        "123",
				State:     "STARTED",
				TurnOrder: []string{"b", "p1", "b", "p2"},
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create game with empty bots array",
		},
		{
			name: "Game gets created successfully",
			input: GameOptions{
				Id:        "123",
				State:     "STARTED",
				TurnOrder: []string{"b", "p1", "b", "p2"},
				Bots:      []*Bot{bot},
			},
			expectedOutput: &Game{
				id:        "123",
				state:     started,
				turnOrder: []string{"b", "p1", "b", "p2"},
				bots:      []*Bot{bot},
			},
			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewGame(tt.input)
			if tt.errorExpected {
				assert.EqualError(t, err, tt.errorString)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}
