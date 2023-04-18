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
				HelpCount: 3,
			},
			expectedOutput: &Bot{
				id:        "123",
				name:      "some name",
				typeOfBot: ai,
				helpCount: 3,
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
				HelpCount:       3,
			},
			expectedOutput: &Bot{
				id:        "123",
				name:      "some name",
				typeOfBot: human,
				player:    &player,
				helpCount: 3,
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
				helpCount: 0,
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

func Test_Bot_Id(t *testing.T) {
	tests := []struct {
		name           string
		input          *Bot
		expectedOutput string
	}{
		{
			name:           "returns Id successfully",
			input:          &Bot{id: "id1"},
			expectedOutput: "id1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Id()
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func Test_Bot_IsAi(t *testing.T) {
	tests := []struct {
		name           string
		input          *Bot
		expectedOutput bool
	}{
		{
			name:           "returns true",
			input:          &Bot{id: "id1", typeOfBot: ai},
			expectedOutput: true,
		},
		{
			name:           "returns false",
			input:          &Bot{id: "id1", typeOfBot: human},
			expectedOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.IsAi()
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func Test_Bot_IsHuman(t *testing.T) {
	tests := []struct {
		name           string
		input          *Bot
		expectedOutput bool
	}{
		{
			name:           "returns true",
			input:          &Bot{id: "id1", typeOfBot: human},
			expectedOutput: true,
		},
		{
			name:           "returns false",
			input:          &Bot{id: "id1", typeOfBot: ai},
			expectedOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.IsHuman()
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func Test_Bot_CanGetHelp(t *testing.T) {
	tests := []struct {
		name           string
		input          *Bot
		expectedOutput bool
	}{
		{
			name:           "returns true",
			input:          &Bot{id: "id1", typeOfBot: human, helpCount: 2},
			expectedOutput: true,
		},
		{
			name:           "returns false",
			input:          &Bot{id: "id1", typeOfBot: human, helpCount: 0},
			expectedOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.CanGetHelp()
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func Test_Bot_ConectPlayer(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			bot    *Bot
			player *Player
		}
		expectedOutput *Bot
		errorExpected  bool
		errorString    string
	}{
		{
			name: "conects player to bot",
			input: struct {
				bot    *Bot
				player *Player
			}{
				bot:    &Bot{id: "b1", typeOfBot: ai},
				player: &Player{id: "p1"},
			},
			expectedOutput: &Bot{
				id:        "b1",
				typeOfBot: human,
				player:    &Player{id: "p1"},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "errors when conecting empty player to a bot",
			input: struct {
				bot    *Bot
				player *Player
			}{
				bot:    &Bot{id: "b1", typeOfBot: ai},
				player: nil,
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "Cannot connect an empty player",
		},
		{
			name: "errors when conecting a new player to a bot already connected to a different player",
			input: struct {
				bot    *Bot
				player *Player
			}{
				bot:    &Bot{id: "b1", typeOfBot: human, player: &Player{id: "p2"}},
				player: &Player{id: "p1"},
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "Cannot replace the connected player",
		},
		{
			name: "errors when conecting a player to a non ai bot",
			input: struct {
				bot    *Bot
				player *Player
			}{
				bot:    &Bot{id: "b1", typeOfBot: undefinedBotType},
				player: &Player{id: "p1"},
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "Can only conect to bot that is currently ai",
		},
		{
			name: "does not error when the same player is connected to the bot",
			input: struct {
				bot    *Bot
				player *Player
			}{
				bot:    &Bot{id: "b1", typeOfBot: human, player: &Player{id: "p1"}},
				player: &Player{id: "p1"},
			},
			expectedOutput: &Bot{
				id:        "b1",
				typeOfBot: human,
				player:    &Player{id: "p1"},
			},
			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.bot.ConnectPlayer(tt.input.player)
			if tt.errorExpected {
				assert.EqualError(t, err, tt.errorString)
			} else {
				assert.EqualValues(t, tt.expectedOutput.id, tt.input.bot.id)
				assert.EqualValues(t, tt.expectedOutput.name, tt.input.bot.name)
				assert.EqualValues(t, tt.expectedOutput.typeOfBot, tt.input.bot.typeOfBot)
				assert.EqualValues(t, tt.expectedOutput.player, tt.input.bot.player)
			}
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
			expectedOutput: []string{"R4-D4", "HAL 99", "Davide", "C-21PO", "EVE-a-L"},
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
