package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_GameViewForPlayer(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			playerId string
			game     *Game
		}
		output *GameView
	}{
		{
			name: "successfully returns a game state when game has just started with the player in it",
			input: struct {
				playerId string
				game     *Game
			}{
				playerId: "player_id1",
				game: &Game{
					state:            started,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
					currentTurnIndex: 1,
					stateTotalTime:   60,
					lastQuestion:     "last question",
					bots: []*Bot{
						{
							id:   "bot_id1",
							name: "bot1",
							player: &Player{
								id: "player_id1",
							},
						},
						{
							id:   "bot_id2",
							name: "bot2",
						},
						{
							id:   "bot_id3",
							name: "bot3",
						},
						{
							id:   "bot_id4",
							name: "bot4",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
						},
						{
							id:   "bot_id5",
							name: "bot5",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
						},
					},
				},
			},
			output: &GameView{
				State:          waitingForPlayersToJoin,
				DisplayMessage: "Please wait as players join in",
				StateTotalTime: 60,
				LastQuestion:   "last question",
				MyBotId:        "bot_id1",
				Conversation: []ConversationElement{
					{
						IsQuestion: true,
						BotId:      "bot_id4",
						Text:       "A question",
					},
					{
						IsQuestion: false,
						BotId:      "bot_id4",
						Text:       "My answer",
					},
					{
						IsQuestion: true,
						BotId:      "bot_id5",
						Text:       "A question",
					},
					{
						IsQuestion: false,
						BotId:      "bot_id5",
						Text:       "My answer",
					},
				},
				Bots: []BotView{
					{
						Id:   "bot_id1",
						Name: "bot1",
					},
					{
						Id:   "bot_id2",
						Name: "bot2",
					},
					{
						Id:   "bot_id3",
						Name: "bot3",
					},
					{
						Id:   "bot_id4",
						Name: "bot4",
						Messages: []string{
							"A question",
							"My answer",
						},
					},
					{
						Id:   "bot_id5",
						Name: "bot5",
						Messages: []string{
							"A question",
							"My answer",
						},
					},
				},
			},
		},
		{
			name: "successfully returns a game state when game has all players joined with the player in it",
			input: struct {
				playerId string
				game     *Game
			}{
				playerId: "player_id1",
				game: &Game{
					state:            playersJoined,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
					currentTurnIndex: 1,
					stateTotalTime:   60,
					lastQuestion:     "last question",
					bots: []*Bot{
						{
							id:   "bot_id1",
							name: "bot1",
							player: &Player{
								id: "player_id1",
							},
						},
						{
							id:   "bot_id2",
							name: "bot2",
						},
						{
							id:   "bot_id3",
							name: "bot3",
						},
						{
							id:   "bot_id4",
							name: "bot4",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
						},
						{
							id:   "bot_id5",
							name: "bot5",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
							player: &Player{
								id: "player_id2",
							},
						},
					},
				},
			},
			output: &GameView{
				State:          waitingForPlayersToJoin,
				DisplayMessage: "Please wait as players join in",
				StateTotalTime: 60,
				LastQuestion:   "last question",
				MyBotId:        "bot_id1",
				Conversation: []ConversationElement{
					{
						IsQuestion: true,
						BotId:      "bot_id4",
						Text:       "A question",
					},
					{
						IsQuestion: false,
						BotId:      "bot_id4",
						Text:       "My answer",
					},
					{
						IsQuestion: true,
						BotId:      "bot_id5",
						Text:       "A question",
					},
					{
						IsQuestion: false,
						BotId:      "bot_id5",
						Text:       "My answer",
					},
				},
				Bots: []BotView{
					{
						Id:   "bot_id1",
						Name: "bot1",
					},
					{
						Id:   "bot_id2",
						Name: "bot2",
					},
					{
						Id:   "bot_id3",
						Name: "bot3",
					},
					{
						Id:   "bot_id4",
						Name: "bot4",
						Messages: []string{
							"A question",
							"My answer",
						},
					},
					{
						Id:   "bot_id5",
						Name: "bot5",
						Messages: []string{
							"A question",
							"My answer",
						},
					},
				},
			},
		},
		{
			name: "successfully returns a game state when game is waiting for Bot to ask question",
			input: struct {
				playerId string
				game     *Game
			}{
				playerId: "player_id1",
				game: &Game{
					state:            waitingForAiQuestion,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
					currentTurnIndex: 1,
					stateTotalTime:   60,
					lastQuestion:     "last question",
					bots: []*Bot{
						{
							id:   "bot_id1",
							name: "bot1",
							player: &Player{
								id: "player_id1",
							},
						},
						{
							id:   "bot_id2",
							name: "bot2",
						},
						{
							id:   "bot_id3",
							name: "bot3",
						},
						{
							id:   "bot_id4",
							name: "bot4",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
						},
						{
							id:   "bot_id5",
							name: "bot5",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
							player: &Player{
								id: "player_id2",
							},
						},
					},
				},
			},
			output: &GameView{
				State:          waitingOnBotToAskAQuestion,
				DisplayMessage: "Please wait as someone is asking a question",
				StateTotalTime: 60,
				LastQuestion:   "last question",
				MyBotId:        "bot_id1",
				Conversation: []ConversationElement{
					{
						IsQuestion: true,
						BotId:      "bot_id4",
						Text:       "A question",
					},
					{
						IsQuestion: false,
						BotId:      "bot_id4",
						Text:       "My answer",
					},
					{
						IsQuestion: true,
						BotId:      "bot_id5",
						Text:       "A question",
					},
					{
						IsQuestion: false,
						BotId:      "bot_id5",
						Text:       "My answer",
					},
				},
				Bots: []BotView{
					{
						Id:   "bot_id1",
						Name: "bot1",
					},
					{
						Id:   "bot_id2",
						Name: "bot2",
					},
					{
						Id:   "bot_id3",
						Name: "bot3",
					},
					{
						Id:   "bot_id4",
						Name: "bot4",
						Messages: []string{
							"A question",
							"My answer",
						},
					},
					{
						Id:   "bot_id5",
						Name: "bot5",
						Messages: []string{
							"A question",
							"My answer",
						},
					},
				},
			},
		},
		{
			name: "successfully returns a game state when game is waiting for Bot to answer a question",
			input: struct {
				playerId string
				game     *Game
			}{
				playerId: "player_id1",
				game: &Game{
					state:                   waitingForAiAnswer,
					turnOrder:               []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
					currentTurnIndex:        4,
					stateTotalTime:          60,
					lastQuestion:            "last question",
					lastQuestionTargetBotId: "bot_id2",
					bots: []*Bot{
						{
							id:   "bot_id1",
							name: "bot1",
							player: &Player{
								id: "player_id1",
							},
						},
						{
							id:   "bot_id2",
							name: "bot2",
						},
						{
							id:   "bot_id3",
							name: "bot3",
						},
						{
							id:   "bot_id4",
							name: "bot4",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
						},
						{
							id:   "bot_id5",
							name: "bot5",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
							player: &Player{
								id: "player_id2",
							},
						},
					},
				},
			},
			output: &GameView{
				State:          waitingOnBotToAnswer,
				DisplayMessage: "Please wait as bot2 is answering the question",
				StateTotalTime: 60,
				LastQuestion:   "last question",
				MyBotId:        "bot_id1",
				Conversation: []ConversationElement{
					{
						IsQuestion: true,
						BotId:      "bot_id4",
						Text:       "A question",
					},
					{
						IsQuestion: false,
						BotId:      "bot_id4",
						Text:       "My answer",
					},
					{
						IsQuestion: true,
						BotId:      "bot_id5",
						Text:       "A question",
					},
					{
						IsQuestion: false,
						BotId:      "bot_id5",
						Text:       "My answer",
					},
				},
				Bots: []BotView{
					{
						Id:   "bot_id1",
						Name: "bot1",
					},
					{
						Id:   "bot_id2",
						Name: "bot2",
					},
					{
						Id:   "bot_id3",
						Name: "bot3",
					},
					{
						Id:   "bot_id4",
						Name: "bot4",
						Messages: []string{
							"A question",
							"My answer",
						},
					},
					{
						Id:   "bot_id5",
						Name: "bot5",
						Messages: []string{
							"A question",
							"My answer",
						},
					},
				},
			},
		},
		{
			name: "successfully returns a game state when game is waiting for requesting player to ask a question",
			input: struct {
				playerId string
				game     *Game
			}{
				playerId: "player_id1",
				game: &Game{
					state:                   waitingForHumanQuestion,
					turnOrder:               []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
					currentTurnIndex:        10,
					stateTotalTime:          60,
					lastQuestion:            "last question",
					lastQuestionTargetBotId: "bot_id2",
					bots: []*Bot{
						{
							id:   "bot_id1",
							name: "bot1",
							player: &Player{
								id: "player_id1",
							},
						},
						{
							id:   "bot_id2",
							name: "bot2",
						},
						{
							id:   "bot_id3",
							name: "bot3",
						},
						{
							id:   "bot_id4",
							name: "bot4",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
						},
						{
							id:   "bot_id5",
							name: "bot5",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
							player: &Player{
								id: "player_id2",
							},
						},
					},
				},
			},
			output: &GameView{
				State:          waitingOnYouToAskAQuestion,
				DisplayMessage: "Please pick a bot and ask a question. OR Click suggest for help!",
				StateTotalTime: 60,
				LastQuestion:   "last question",
				MyBotId:        "bot_id1",
				Conversation: []ConversationElement{
					{
						IsQuestion: true,
						BotId:      "bot_id4",
						Text:       "A question",
					},
					{
						IsQuestion: false,
						BotId:      "bot_id4",
						Text:       "My answer",
					},
					{
						IsQuestion: true,
						BotId:      "bot_id5",
						Text:       "A question",
					},
					{
						IsQuestion: false,
						BotId:      "bot_id5",
						Text:       "My answer",
					},
				},
				Bots: []BotView{
					{
						Id:   "bot_id1",
						Name: "bot1",
					},
					{
						Id:   "bot_id2",
						Name: "bot2",
					},
					{
						Id:   "bot_id3",
						Name: "bot3",
					},
					{
						Id:   "bot_id4",
						Name: "bot4",
						Messages: []string{
							"A question",
							"My answer",
						},
					},
					{
						Id:   "bot_id5",
						Name: "bot5",
						Messages: []string{
							"A question",
							"My answer",
						},
					},
				},
			},
		},
		{
			name: "successfully returns a game state when game is waiting for requesting player to answer a question",
			input: struct {
				playerId string
				game     *Game
			}{
				playerId: "player_id1",
				game: &Game{
					state:                   waitingForHumanAnswer,
					turnOrder:               []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
					currentTurnIndex:        3,
					stateTotalTime:          60,
					lastQuestion:            "last question",
					lastQuestionTargetBotId: "bot_id1",
					bots: []*Bot{
						{
							id:   "bot_id1",
							name: "bot1",
							player: &Player{
								id: "player_id1",
							},
						},
						{
							id:   "bot_id2",
							name: "bot2",
						},
						{
							id:   "bot_id3",
							name: "bot3",
						},
						{
							id:   "bot_id4",
							name: "bot4",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
						},
						{
							id:   "bot_id5",
							name: "bot5",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
							player: &Player{
								id: "player_id2",
							},
						},
					},
				},
			},
			output: &GameView{
				State:          waitingOnYouToAnswer,
				DisplayMessage: "Please answer the question. OR Click suggest for help!",
				StateTotalTime: 60,
				LastQuestion:   "last question",
				MyBotId:        "bot_id1",
				Conversation: []ConversationElement{
					{
						IsQuestion: true,
						BotId:      "bot_id4",
						Text:       "A question",
					},
					{
						IsQuestion: false,
						BotId:      "bot_id4",
						Text:       "My answer",
					},
					{
						IsQuestion: true,
						BotId:      "bot_id5",
						Text:       "A question",
					},
					{
						IsQuestion: false,
						BotId:      "bot_id5",
						Text:       "My answer",
					},
				},
				Bots: []BotView{
					{
						Id:   "bot_id1",
						Name: "bot1",
					},
					{
						Id:   "bot_id2",
						Name: "bot2",
					},
					{
						Id:   "bot_id3",
						Name: "bot3",
					},
					{
						Id:   "bot_id4",
						Name: "bot4",
						Messages: []string{
							"A question",
							"My answer",
						},
					},
					{
						Id:   "bot_id5",
						Name: "bot5",
						Messages: []string{
							"A question",
							"My answer",
						},
					},
				},
			},
		},
		{
			name: "successfully returns a game state when game is waiting for other player to ask a question",
			input: struct {
				playerId string
				game     *Game
			}{
				playerId: "player_id2",
				game: &Game{
					state:                   waitingForHumanQuestion,
					turnOrder:               []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
					currentTurnIndex:        10,
					stateTotalTime:          60,
					lastQuestion:            "last question",
					lastQuestionTargetBotId: "bot_id2",
					bots: []*Bot{
						{
							id:   "bot_id1",
							name: "bot1",
							player: &Player{
								id: "player_id1",
							},
						},
						{
							id:   "bot_id2",
							name: "bot2",
						},
						{
							id:   "bot_id3",
							name: "bot3",
						},
						{
							id:   "bot_id4",
							name: "bot4",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
						},
						{
							id:   "bot_id5",
							name: "bot5",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
							player: &Player{
								id: "player_id2",
							},
						},
					},
				},
			},
			output: &GameView{
				State:          waitingOnBotToAskAQuestion,
				DisplayMessage: "Please wait as someone is asking a question",
				StateTotalTime: 60,
				LastQuestion:   "last question",
				MyBotId:        "bot_id5",
				Conversation: []ConversationElement{
					{
						IsQuestion: true,
						BotId:      "bot_id4",
						Text:       "A question",
					},
					{
						IsQuestion: false,
						BotId:      "bot_id4",
						Text:       "My answer",
					},
					{
						IsQuestion: true,
						BotId:      "bot_id5",
						Text:       "A question",
					},
					{
						IsQuestion: false,
						BotId:      "bot_id5",
						Text:       "My answer",
					},
				},
				Bots: []BotView{
					{
						Id:   "bot_id1",
						Name: "bot1",
					},
					{
						Id:   "bot_id2",
						Name: "bot2",
					},
					{
						Id:   "bot_id3",
						Name: "bot3",
					},
					{
						Id:   "bot_id4",
						Name: "bot4",
						Messages: []string{
							"A question",
							"My answer",
						},
					},
					{
						Id:   "bot_id5",
						Name: "bot5",
						Messages: []string{
							"A question",
							"My answer",
						},
					},
				},
			},
		},
		{
			name: "successfully returns a game state when game is waiting for other player to answer a question",
			input: struct {
				playerId string
				game     *Game
			}{
				playerId: "player_id2",
				game: &Game{
					state:                   waitingForHumanAnswer,
					turnOrder:               []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
					currentTurnIndex:        3,
					stateTotalTime:          60,
					lastQuestion:            "last question",
					lastQuestionTargetBotId: "bot_id1",
					bots: []*Bot{
						{
							id:   "bot_id1",
							name: "bot1",
							player: &Player{
								id: "player_id1",
							},
						},
						{
							id:   "bot_id2",
							name: "bot2",
						},
						{
							id:   "bot_id3",
							name: "bot3",
						},
						{
							id:   "bot_id4",
							name: "bot4",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
						},
						{
							id:   "bot_id5",
							name: "bot5",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
							player: &Player{
								id: "player_id2",
							},
						},
					},
				},
			},
			output: &GameView{
				State:          waitingOnBotToAnswer,
				DisplayMessage: "Please wait as bot1 is answering the question",
				StateTotalTime: 60,
				LastQuestion:   "last question",
				MyBotId:        "bot_id5",
				Conversation: []ConversationElement{
					{
						IsQuestion: true,
						BotId:      "bot_id4",
						Text:       "A question",
					},
					{
						IsQuestion: false,
						BotId:      "bot_id4",
						Text:       "My answer",
					},
					{
						IsQuestion: true,
						BotId:      "bot_id5",
						Text:       "A question",
					},
					{
						IsQuestion: false,
						BotId:      "bot_id5",
						Text:       "My answer",
					},
				},
				Bots: []BotView{
					{
						Id:   "bot_id1",
						Name: "bot1",
					},
					{
						Id:   "bot_id2",
						Name: "bot2",
					},
					{
						Id:   "bot_id3",
						Name: "bot3",
					},
					{
						Id:   "bot_id4",
						Name: "bot4",
						Messages: []string{
							"A question",
							"My answer",
						},
					},
					{
						Id:   "bot_id5",
						Name: "bot5",
						Messages: []string{
							"A question",
							"My answer",
						},
					},
				},
			},
		},
		{
			name: "returns nil if playerId is blank",
			input: struct {
				playerId string
				game     *Game
			}{
				playerId: "",
				game: &Game{
					state:                   waitingForAiAnswer,
					turnOrder:               []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
					currentTurnIndex:        4,
					stateTotalTime:          60,
					lastQuestionTargetBotId: "bot_id2",
					bots: []*Bot{
						{
							id:   "bot_id1",
							name: "bot1",
							player: &Player{
								id: "player_id1",
							},
						},
						{
							id:   "bot_id2",
							name: "bot2",
						},
						{
							id:   "bot_id3",
							name: "bot3",
						},
						{
							id:   "bot_id4",
							name: "bot4",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
						},
						{
							id:   "bot_id5",
							name: "bot5",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
							player: &Player{
								id: "player_id2",
							},
						},
					},
				},
			},
			output: nil,
		},
		{
			name: "returns nil if playerId is not in the game",
			input: struct {
				playerId string
				game     *Game
			}{
				playerId: "player_id3",
				game: &Game{
					state:                   waitingForAiAnswer,
					turnOrder:               []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
					currentTurnIndex:        4,
					lastQuestionTargetBotId: "bot_id2",
					bots: []*Bot{
						{
							id:   "bot_id1",
							name: "bot1",
							player: &Player{
								id: "player_id1",
							},
						},
						{
							id:   "bot_id2",
							name: "bot2",
						},
						{
							id:   "bot_id3",
							name: "bot3",
						},
						{
							id:   "bot_id4",
							name: "bot4",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
						},
						{
							id:   "bot_id5",
							name: "bot5",
							messages: []Message{
								{Text: "A question", CreatedAt: time.Now()},
								{Text: "My answer", CreatedAt: time.Now()},
							},
							player: &Player{
								id: "player_id2",
							},
						},
					},
				},
			},
			output: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameView := tt.input.game.GameViewForPlayer(tt.input.playerId)
			assert.Equal(t, gameView, tt.output)
		})
	}
}
