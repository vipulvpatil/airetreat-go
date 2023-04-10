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
						},
						{
							id:   "bot_id5",
							name: "bot5",
						},
					},
					messages: []*Message{
						{SourceBotId: "bot_id2", TargetBotId: "bot_id4", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id4", TargetBotId: "bot_id4", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
						{SourceBotId: "bot_id2", TargetBotId: "bot_id5", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
					},
				},
			},
			output: &GameView{
				State:          waitingForPlayersToJoin,
				DisplayMessage: "Waiting for players to join in",
				StateTotalTime: 60,
				LastQuestion:   "last question",
				MyBotId:        "bot_id1",
				DetailedMessages: []DetailedMessage{
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id4",
						SourceBotName: "bot4",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
					},
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id5",
						SourceBotName: "bot5",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
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
					},
					{
						Id:   "bot_id5",
						Name: "bot5",
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
						},
						{
							id:   "bot_id5",
							name: "bot5",
							player: &Player{
								id: "player_id2",
							},
						},
					},
					messages: []*Message{
						{SourceBotId: "bot_id2", TargetBotId: "bot_id4", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id4", TargetBotId: "bot_id4", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
						{SourceBotId: "bot_id2", TargetBotId: "bot_id5", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
					},
				},
			},
			output: &GameView{
				State:          waitingForPlayersToJoin,
				DisplayMessage: "Waiting for players to join in",
				StateTotalTime: 60,
				LastQuestion:   "last question",
				MyBotId:        "bot_id1",
				DetailedMessages: []DetailedMessage{
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id4",
						SourceBotName: "bot4",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
					},
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id5",
						SourceBotName: "bot5",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
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
					},
					{
						Id:   "bot_id5",
						Name: "bot5",
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
						},
						{
							id:   "bot_id5",
							name: "bot5",
							player: &Player{
								id: "player_id2",
							},
						},
					},
					messages: []*Message{
						{SourceBotId: "bot_id2", TargetBotId: "bot_id4", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id4", TargetBotId: "bot_id4", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
						{SourceBotId: "bot_id2", TargetBotId: "bot_id5", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
					},
				},
			},
			output: &GameView{
				State:          waitingOnBotToAskAQuestion,
				DisplayMessage: "Someone is asking a question",
				StateTotalTime: 60,
				LastQuestion:   "last question",
				MyBotId:        "bot_id1",
				DetailedMessages: []DetailedMessage{
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id4",
						SourceBotName: "bot4",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
					},
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id5",
						SourceBotName: "bot5",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
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
					},
					{
						Id:   "bot_id5",
						Name: "bot5",
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
						},
						{
							id:   "bot_id5",
							name: "bot5",
							player: &Player{
								id: "player_id2",
							},
						},
					},
					messages: []*Message{
						{SourceBotId: "bot_id2", TargetBotId: "bot_id4", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id4", TargetBotId: "bot_id4", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
						{SourceBotId: "bot_id2", TargetBotId: "bot_id5", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
					},
				},
			},
			output: &GameView{
				State:          waitingOnBotToAnswer,
				DisplayMessage: "bot2 is answering the question",
				StateTotalTime: 60,
				LastQuestion:   "last question",
				MyBotId:        "bot_id1",
				DetailedMessages: []DetailedMessage{
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id4",
						SourceBotName: "bot4",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
					},
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id5",
						SourceBotName: "bot5",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
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
					},
					{
						Id:   "bot_id5",
						Name: "bot5",
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
						},
						{
							id:   "bot_id5",
							name: "bot5",
							player: &Player{
								id: "player_id2",
							},
						},
					},
					messages: []*Message{
						{SourceBotId: "bot_id2", TargetBotId: "bot_id4", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id4", TargetBotId: "bot_id4", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
						{SourceBotId: "bot_id2", TargetBotId: "bot_id5", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
					},
				},
			},
			output: &GameView{
				State:          waitingOnYouToAskAQuestion,
				DisplayMessage: "Ask a question. OR Click suggest for help!",
				StateTotalTime: 60,
				LastQuestion:   "last question",
				MyBotId:        "bot_id1",
				DetailedMessages: []DetailedMessage{
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id4",
						SourceBotName: "bot4",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
					},
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id5",
						SourceBotName: "bot5",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
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
					},
					{
						Id:   "bot_id5",
						Name: "bot5",
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
						},
						{
							id:   "bot_id5",
							name: "bot5",
							player: &Player{
								id: "player_id2",
							},
						},
					},
					messages: []*Message{
						{SourceBotId: "bot_id2", TargetBotId: "bot_id4", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id4", TargetBotId: "bot_id4", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
						{SourceBotId: "bot_id2", TargetBotId: "bot_id5", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
					},
				},
			},
			output: &GameView{
				State:          waitingOnYouToAnswer,
				DisplayMessage: "Answer the question. OR Click suggest for help!",
				StateTotalTime: 60,
				LastQuestion:   "last question",
				MyBotId:        "bot_id1",
				DetailedMessages: []DetailedMessage{
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id4",
						SourceBotName: "bot4",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
					},
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id5",
						SourceBotName: "bot5",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
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
					},
					{
						Id:   "bot_id5",
						Name: "bot5",
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
						},
						{
							id:   "bot_id5",
							name: "bot5",
							player: &Player{
								id: "player_id2",
							},
						},
					},
					messages: []*Message{
						{SourceBotId: "bot_id2", TargetBotId: "bot_id4", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id4", TargetBotId: "bot_id4", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
						{SourceBotId: "bot_id2", TargetBotId: "bot_id5", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
					},
				},
			},
			output: &GameView{
				State:          waitingOnBotToAskAQuestion,
				DisplayMessage: "Someone is asking a question",
				StateTotalTime: 60,
				LastQuestion:   "last question",
				MyBotId:        "bot_id5",
				DetailedMessages: []DetailedMessage{
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id4",
						SourceBotName: "bot4",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
					},
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id5",
						SourceBotName: "bot5",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
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
					},
					{
						Id:   "bot_id5",
						Name: "bot5",
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
						},
						{
							id:   "bot_id5",
							name: "bot5",
							player: &Player{
								id: "player_id2",
							},
						},
					},
					messages: []*Message{
						{SourceBotId: "bot_id2", TargetBotId: "bot_id4", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id4", TargetBotId: "bot_id4", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
						{SourceBotId: "bot_id2", TargetBotId: "bot_id5", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
					},
				},
			},
			output: &GameView{
				State:          waitingOnBotToAnswer,
				DisplayMessage: "bot1 is answering the question",
				StateTotalTime: 60,
				LastQuestion:   "last question",
				MyBotId:        "bot_id5",
				DetailedMessages: []DetailedMessage{
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id4",
						SourceBotName: "bot4",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
					},
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id5",
						SourceBotName: "bot5",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
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
					},
					{
						Id:   "bot_id5",
						Name: "bot5",
					},
				},
			},
		},
		{
			name: "successfully returns a game state when game is finished with requesting player winning",
			input: struct {
				playerId string
				game     *Game
			}{
				playerId: "player_id2",
				game: &Game{
					state:                   finished,
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
						},
						{
							id:   "bot_id5",
							name: "bot5",
							player: &Player{
								id: "player_id2",
							},
						},
					},
					messages: []*Message{
						{SourceBotId: "bot_id2", TargetBotId: "bot_id4", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id4", TargetBotId: "bot_id4", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
						{SourceBotId: "bot_id2", TargetBotId: "bot_id5", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
					},
					result:       "bot5 won",
					winningBotId: "bot_id5",
				},
			},
			output: &GameView{
				State:          gameViewState(youWon),
				DisplayMessage: "bot5 won",
				StateTotalTime: 60,
				LastQuestion:   "last question",
				MyBotId:        "bot_id5",
				DetailedMessages: []DetailedMessage{
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id4",
						SourceBotName: "bot4",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
					},
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id5",
						SourceBotName: "bot5",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
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
					},
					{
						Id:   "bot_id5",
						Name: "bot5",
					},
				},
			},
		},
		{
			name: "successfully returns a game state when game is finished with requesting player losing",
			input: struct {
				playerId string
				game     *Game
			}{
				playerId: "player_id2",
				game: &Game{
					state:                   finished,
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
						},
						{
							id:   "bot_id5",
							name: "bot5",
							player: &Player{
								id: "player_id2",
							},
						},
					},
					messages: []*Message{
						{SourceBotId: "bot_id2", TargetBotId: "bot_id4", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id4", TargetBotId: "bot_id4", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
						{SourceBotId: "bot_id2", TargetBotId: "bot_id5", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
					},
					result:       "bot3 won",
					winningBotId: "bot_id3",
				},
			},
			output: &GameView{
				State:          gameViewState(youLost),
				DisplayMessage: "bot3 won",
				StateTotalTime: 60,
				LastQuestion:   "last question",
				MyBotId:        "bot_id5",
				DetailedMessages: []DetailedMessage{
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id4",
						SourceBotName: "bot4",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
					},
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id5",
						SourceBotName: "bot5",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
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
					},
					{
						Id:   "bot_id5",
						Name: "bot5",
					},
				},
			},
		},
		{
			name: "successfully returns a game state when game is finished with time out",
			input: struct {
				playerId string
				game     *Game
			}{
				playerId: "player_id2",
				game: &Game{
					state:                   finished,
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
						},
						{
							id:   "bot_id5",
							name: "bot5",
							player: &Player{
								id: "player_id2",
							},
						},
					},
					messages: []*Message{
						{SourceBotId: "bot_id2", TargetBotId: "bot_id4", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id4", TargetBotId: "bot_id4", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
						{SourceBotId: "bot_id2", TargetBotId: "bot_id5", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
					},
					result:       "time ran out",
					winningBotId: "",
				},
			},
			output: &GameView{
				State:          gameViewState(timeUp),
				DisplayMessage: "time ran out",
				StateTotalTime: 60,
				LastQuestion:   "last question",
				MyBotId:        "bot_id5",
				DetailedMessages: []DetailedMessage{
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id4",
						SourceBotName: "bot4",
						TargetBotId:   "bot_id4",
						TargetBotName: "bot4",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
					},
					{
						SourceBotId:   "bot_id2",
						SourceBotName: "bot2",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "A question",
						CreatedAt:     time.Now(),
						MessageType:   "question",
					},
					{
						SourceBotId:   "bot_id5",
						SourceBotName: "bot5",
						TargetBotId:   "bot_id5",
						TargetBotName: "bot5",
						Text:          "My answer",
						CreatedAt:     time.Now(),
						MessageType:   "answer",
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
					},
					{
						Id:   "bot_id5",
						Name: "bot5",
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
						},
						{
							id:   "bot_id5",
							name: "bot5",
							player: &Player{
								id: "player_id2",
							},
						},
					},
					messages: []*Message{
						{SourceBotId: "bot_id2", TargetBotId: "bot_id4", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id4", TargetBotId: "bot_id4", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
						{SourceBotId: "bot_id2", TargetBotId: "bot_id5", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
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
						},
						{
							id:   "bot_id5",
							name: "bot5",
							player: &Player{
								id: "player_id2",
							},
						},
					},
					messages: []*Message{
						{SourceBotId: "bot_id2", TargetBotId: "bot_id4", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id4", TargetBotId: "bot_id4", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
						{SourceBotId: "bot_id2", TargetBotId: "bot_id5", Text: "A question", CreatedAt: time.Now(), MessageType: "question"},
						{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "My answer", CreatedAt: time.Now(), MessageType: "answer"},
					},
				},
			},
			output: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameView := tt.input.game.GameViewForPlayer(tt.input.playerId)
			if tt.output == nil {
				assert.Nil(t, gameView)
			} else {
				AssertEqualGameView(t, gameView, tt.output)
			}
		})
	}
}
