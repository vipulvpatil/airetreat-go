package model

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewBot(t *testing.T) {
	player := Player{id: "p_id"}
	tests := []struct {
		name           string
		input          BotOptions
		expectedOutput *Bot
		errorExpected  bool
		errorString    string
	}{
		{
			name:           "id is empty",
			input:          BotOptions{},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create bot with an empty id",
		},
		{
			name:           "name is empty",
			input:          BotOptions{Id: "1"},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create bot with an empty name",
		},
		{
			name:           "botType is invalid",
			input:          BotOptions{Id: "1", Name: "botname", TypeOfBot: "CHAT"},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create bot with an invalid botType",
		},
		{
			name: "errors when botType is not human but a connected Player is provided",
			input: BotOptions{
				Id:              "123",
				Name:            "some name",
				TypeOfBot:       "AI",
				ConnectedPlayer: &player,
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create a bot of non-human type with a connected Player",
		},
		{
			name: "errors when botType is human but a connected Player is not provided",
			input: BotOptions{
				Id:        "123",
				Name:      "some name",
				TypeOfBot: "HUMAN",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create a bot of human type without a connected Player",
		},
		{
			name: "Bot gets created successfully with provided botType",
			input: BotOptions{
				Id:        "123",
				Name:      "some name",
				TypeOfBot: "AI",
			},
			expectedOutput: &Bot{
				id:        "123",
				name:      "some name",
				typeOfBot: ai,
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "Bot gets created successfully with human botType",
			input: BotOptions{
				Id:              "123",
				Name:            "some name",
				TypeOfBot:       "HUMAN",
				ConnectedPlayer: &player,
			},
			expectedOutput: &Bot{
				id:        "123",
				name:      "some name",
				typeOfBot: human,
				player:    &player,
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "Bot gets created successfully with default botType",
			input: BotOptions{
				Id:   "123",
				Name: "some name",
			},
			expectedOutput: &Bot{
				id:        "123",
				name:      "some name",
				typeOfBot: ai,
			},
			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewBot(tt.input)
			if tt.errorExpected {
				assert.EqualError(t, err, tt.errorString)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func Test_RandomBotNames(t *testing.T) {
	tests := []struct {
		name           string
		input          int64
		expectedOutput []string
	}{
		{
			name:           "random names are generated",
			input:          10,
			expectedOutput: []string{"The Hivey-five", "Gart", "RoboCupp", "T-3PO", "Avis"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rand.Seed(tt.input)
			result := RandomBotNames()
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}
