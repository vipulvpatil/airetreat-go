package model

import (
	"math/rand"
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

func Test_Game_HasJustStarted(t *testing.T) {
	tests := []struct {
		name           string
		input          *Game
		expectedOutput bool
	}{
		{
			name:           "returns true if game has just started",
			input:          &Game{state: started},
			expectedOutput: true,
		},
		{
			name:           "returns false if game has moved to another state",
			input:          &Game{state: playersJoined},
			expectedOutput: false,
		},
		{
			name:           "returns false if game is nil",
			input:          nil,
			expectedOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.HasJustStarted()
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func Test_Game_GetOneRandomAiBot(t *testing.T) {
	tests := []struct {
		name           string
		input          *Game
		expectedOutput *Bot
		errorExpected  bool
		errorString    string
	}{
		{
			name:           "errors if game is nil",
			input:          nil,
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "attempting to get bots from a nil game",
		},
		{
			name: "errors if game has no ai bots",
			input: &Game{state: started, bots: []*Bot{
				{
					id:        "bot_id1",
					name:      "bot1",
					typeOfBot: human,
				},
			},
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "no AI bots in the game",
		},
		{
			name: "returns an ai bot if game has only 1 ai bot",
			input: &Game{state: started, bots: []*Bot{
				{
					id:        "bot_id1",
					name:      "bot1",
					typeOfBot: human,
				},
				{
					id:        "bot_id2",
					name:      "bot2",
					typeOfBot: ai,
				},
			}},
			expectedOutput: &Bot{
				id:        "bot_id2",
				name:      "bot2",
				typeOfBot: ai,
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "returns a ai random bot if game has multiple ai bots",
			input: &Game{state: started, bots: []*Bot{
				{
					id:        "bot_id1",
					name:      "bot1",
					typeOfBot: ai,
				},
				{
					id:        "bot_id2",
					name:      "bot2",
					typeOfBot: ai,
				},
				{
					id:        "bot_id3",
					name:      "bot3",
					typeOfBot: human,
				},
			}},
			expectedOutput: &Bot{
				id:        "bot_id1",
				name:      "bot1",
				typeOfBot: ai,
			},
			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rand.Seed(0)
			result, err := tt.input.GetOneRandomAiBot()
			if !tt.errorExpected {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, result)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}
