package model

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewGame(t *testing.T) {
	bot, err := NewBot(BotOptions{
		Id:        "bot_id1",
		Name:      "Bot name",
		TypeOfBot: "AI",
	})
	assert.NoError(t, err)
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
				TurnOrder: []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create game with empty bots array",
		},
		{
			name: "invalid last question target bot",
			input: GameOptions{
				Id:                      "123",
				State:                   "STARTED",
				TurnOrder:               []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
				LastQuestionTargetBotId: "bot_id2",
				Bots:                    []*Bot{bot},
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "cannot create game with incorrect last question target bot id",
		},
		{
			name: "Game gets created successfully",
			input: GameOptions{
				Id:                      "123",
				State:                   "STARTED",
				TurnOrder:               []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
				LastQuestion:            "Question",
				LastQuestionTargetBotId: "bot_id1",
				Bots:                    []*Bot{bot},
			},
			expectedOutput: &Game{
				id:                      "123",
				state:                   started,
				turnOrder:               []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
				lastQuestion:            "Question",
				lastQuestionTargetBotId: "bot_id1",
				bots:                    []*Bot{bot},
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

func Test_HasJustStarted(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.HasJustStarted()
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func Test_GetOneRandomAiBot(t *testing.T) {
	tests := []struct {
		name           string
		input          *Game
		expectedOutput *Bot
		errorExpected  bool
		errorString    string
	}{
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

func Test_BotWithId(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			game  *Game
			botId string
		}
		expectedOutput *Bot
	}{
		{
			name: "returns bot if present in game",
			input: struct {
				game  *Game
				botId string
			}{
				game: &Game{
					state: started,
					bots: []*Bot{
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
							player: &Player{
								id: "player_id1",
							},
						},
					},
				},
				botId: "bot_id1",
			},
			expectedOutput: &Bot{
				id:        "bot_id1",
				name:      "bot1",
				typeOfBot: 1,
			},
		},
		{
			name: "returns nil if game does not have the specified bot",
			input: struct {
				game  *Game
				botId string
			}{
				game: &Game{
					state: started,
					bots: []*Bot{
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
					},
				},
				botId: "bot_id4",
			},
			expectedOutput: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.game.BotWithId(tt.input.botId)
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}
func Test_HasPlayer(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			game     *Game
			playerId string
		}
		expectedOutput bool
	}{
		{
			name: "returns true if game has player",
			input: struct {
				game     *Game
				playerId string
			}{
				game: &Game{
					state: started,
					bots: []*Bot{
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
							player: &Player{
								id: "player_id1",
							},
						},
					},
				},
				playerId: "player_id1",
			},
			expectedOutput: true,
		},
		{
			name: "returns false if game does not have player",
			input: struct {
				game     *Game
				playerId string
			}{
				game: &Game{
					state: started,
					bots: []*Bot{
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
					},
				},
				playerId: "player_id1",
			},
			expectedOutput: false,
		},
		{
			name: "returns false if playerId is blank",
			input: struct {
				game     *Game
				playerId string
			}{
				game: &Game{
					state: started,
					bots: []*Bot{
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
					},
				},
				playerId: "",
			},
			expectedOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.game.HasPlayer(tt.input.playerId)
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func Test_BotWithPlayerId(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			game     *Game
			playerId string
		}
		expectedOutput *Bot
	}{
		{
			name: "returns the corresponding bot connected to the player",
			input: struct {
				game     *Game
				playerId string
			}{
				game: &Game{
					state: started,
					bots: []*Bot{
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
							player: &Player{
								id: "player_id1",
							},
						},
					},
				},
				playerId: "player_id1",
			},
			expectedOutput: &Bot{
				id:        "bot_id3",
				name:      "bot3",
				typeOfBot: human,
				player: &Player{
					id: "player_id1",
				},
			},
		},
		{
			name: "returns nil if game does not have player",
			input: struct {
				game     *Game
				playerId string
			}{
				game: &Game{
					state: started,
					bots: []*Bot{
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
					},
				},
				playerId: "player_id1",
			},
			expectedOutput: nil,
		},
		{
			name: "returns nil if playerId is blank",
			input: struct {
				game     *Game
				playerId string
			}{
				game: &Game{
					state: started,
					bots: []*Bot{
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
					},
				},
				playerId: "",
			},
			expectedOutput: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.game.BotWithPlayerId(tt.input.playerId)
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func Test_GetGameUpdateAfterIncomingMessage(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			game        *Game
			sourceBotId string
			targetBotId string
			text        string
		}
		expectedOutput *GameUpdate
		errorExpected  bool
		errorString    string
	}{
		{
			name: "returns the game update for the incoming message given a game in waitingForAiQuestion",
			input: struct {
				game        *Game
				sourceBotId string
				targetBotId string
				text        string
			}{
				game: &Game{
					state:            waitingForAiQuestion,
					currentTurnIndex: 1,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3"},
					bots: []*Bot{
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
							player: &Player{
								id: "player_id1",
							},
						},
					},
				},
				sourceBotId: "bot_id2",
				targetBotId: "bot_id3",
				text:        "What is the answer?",
			},
			expectedOutput: func() *GameUpdate {
				stateHandled := false
				lastQuestion := "What is the answer?"
				lastQuestionTargetBotId := "bot_id3"

				return &GameUpdate{
					State:                   GameState("WAITING_FOR_HUMAN_ANSWER"),
					CurrentTurnIndex:        nil,
					StateHandled:            &stateHandled,
					LastQuestion:            &lastQuestion,
					LastQuestionTargetBotId: &lastQuestionTargetBotId,
				}
			}(),
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "returns the game update for the incoming message given a game in waitingForAiAnswer",
			input: struct {
				game        *Game
				sourceBotId string
				targetBotId string
				text        string
			}{
				game: &Game{
					state:            waitingForAiAnswer,
					currentTurnIndex: 1,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3"},
					bots: []*Bot{
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
							player: &Player{
								id: "player_id1",
							},
						},
					},
					lastQuestion:            "Is this a question?",
					lastQuestionTargetBotId: "bot_id2",
				},
				sourceBotId: "bot_id2",
				targetBotId: "bot_id2",
				text:        "This is the answer",
			},
			expectedOutput: func() *GameUpdate {
				currentTurnIndex := int64(2)
				stateHandled := false

				return &GameUpdate{
					State:                   GameState("WAITING_FOR_HUMAN_QUESTION"),
					CurrentTurnIndex:        &currentTurnIndex,
					StateHandled:            &stateHandled,
					LastQuestion:            nil,
					LastQuestionTargetBotId: nil,
				}
			}(),
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "returns the game update for the incoming message given a game in waitingForHumanQuestion",
			input: struct {
				game        *Game
				sourceBotId string
				targetBotId string
				text        string
			}{
				game: &Game{
					state:            waitingForHumanQuestion,
					currentTurnIndex: 2,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3"},
					bots: []*Bot{
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
							player: &Player{
								id: "player_id1",
							},
						},
					},
				},
				sourceBotId: "bot_id3",
				targetBotId: "bot_id1",
				text:        "What is the next question?",
			},
			expectedOutput: func() *GameUpdate {
				stateHandled := false
				lastQuestion := "What is the next question?"
				lastQuestionTargetBotId := "bot_id1"

				return &GameUpdate{
					State:                   GameState("WAITING_FOR_AI_ANSWER"),
					CurrentTurnIndex:        nil,
					StateHandled:            &stateHandled,
					LastQuestion:            &lastQuestion,
					LastQuestionTargetBotId: &lastQuestionTargetBotId,
				}
			}(),
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "returns the game update for the incoming message given a game in waitingForHumanAnswer",
			input: struct {
				game        *Game
				sourceBotId string
				targetBotId string
				text        string
			}{
				game: &Game{
					state:            waitingForHumanAnswer,
					currentTurnIndex: 1,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3"},
					bots: []*Bot{
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
							player: &Player{
								id: "player_id1",
							},
						},
					},
					lastQuestion:            "Is this a question?",
					lastQuestionTargetBotId: "bot_id3",
				},
				sourceBotId: "bot_id3",
				targetBotId: "bot_id3",
				text:        "This is the answer",
			},
			expectedOutput: func() *GameUpdate {
				currentTurnIndex := int64(2)
				stateHandled := false

				return &GameUpdate{
					State:                   GameState("WAITING_FOR_HUMAN_QUESTION"),
					CurrentTurnIndex:        &currentTurnIndex,
					StateHandled:            &stateHandled,
					LastQuestion:            nil,
					LastQuestionTargetBotId: nil,
				}
			}(),
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "returns the game update for the incoming message given a game in waitingForHumanAnswer and the next turn is AI",
			input: struct {
				game        *Game
				sourceBotId string
				targetBotId string
				text        string
			}{
				game: &Game{
					state:            waitingForHumanAnswer,
					currentTurnIndex: 0,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3"},
					bots: []*Bot{
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
							player: &Player{
								id: "player_id1",
							},
						},
					},
					lastQuestion:            "Is this a question?",
					lastQuestionTargetBotId: "bot_id3",
				},
				sourceBotId: "bot_id3",
				targetBotId: "bot_id3",
				text:        "This is the answer",
			},
			expectedOutput: func() *GameUpdate {
				currentTurnIndex := int64(1)
				stateHandled := false

				return &GameUpdate{
					State:                   GameState("WAITING_FOR_AI_QUESTION"),
					CurrentTurnIndex:        &currentTurnIndex,
					StateHandled:            &stateHandled,
					LastQuestion:            nil,
					LastQuestionTargetBotId: nil,
				}
			}(),
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "errors if sourceBotId is not in game",
			input: struct {
				game        *Game
				sourceBotId string
				targetBotId string
				text        string
			}{
				game: &Game{
					state:            waitingForHumanAnswer,
					currentTurnIndex: 1,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3"},
					bots: []*Bot{
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
							player: &Player{
								id: "player_id1",
							},
						},
					},
					lastQuestion:            "Is this a question?",
					lastQuestionTargetBotId: "bot_id3",
				},
				sourceBotId: "bot_id5",
				targetBotId: "bot_id3",
				text:        "This is the answer",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "invalid sourceBotId",
		},
		{
			name: "errors if sourceBotId is blank",
			input: struct {
				game        *Game
				sourceBotId string
				targetBotId string
				text        string
			}{
				game: &Game{
					state:            waitingForHumanAnswer,
					currentTurnIndex: 1,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3"},
					bots: []*Bot{
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
							player: &Player{
								id: "player_id1",
							},
						},
					},
					lastQuestion:            "Is this a question?",
					lastQuestionTargetBotId: "bot_id3",
				},
				sourceBotId: "",
				targetBotId: "bot_id3",
				text:        "This is the answer",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "invalid sourceBotId",
		},
		{
			name: "errors if targetBotId is not in game",
			input: struct {
				game        *Game
				sourceBotId string
				targetBotId string
				text        string
			}{
				game: &Game{
					state:            waitingForHumanAnswer,
					currentTurnIndex: 1,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3"},
					bots: []*Bot{
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
							player: &Player{
								id: "player_id1",
							},
						},
					},
					lastQuestion:            "Is this a question?",
					lastQuestionTargetBotId: "bot_id3",
				},
				sourceBotId: "bot_id3",
				targetBotId: "bot_id5",
				text:        "This is the answer",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "invalid targetBotId",
		},
		{
			name: "errors if targetBotId is blank",
			input: struct {
				game        *Game
				sourceBotId string
				targetBotId string
				text        string
			}{
				game: &Game{
					state:            waitingForHumanAnswer,
					currentTurnIndex: 1,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3"},
					bots: []*Bot{
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
							player: &Player{
								id: "player_id1",
							},
						},
					},
					lastQuestion:            "Is this a question?",
					lastQuestionTargetBotId: "bot_id3",
				},
				sourceBotId: "bot_id3",
				targetBotId: "",
				text:        "This is the answer",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "invalid targetBotId",
		},
		{
			name: "errors if game is not waiting for messages",
			input: struct {
				game        *Game
				sourceBotId string
				targetBotId string
				text        string
			}{
				game: &Game{
					state:            playersJoined,
					currentTurnIndex: 1,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3"},
					bots: []*Bot{
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
							player: &Player{
								id: "player_id1",
							},
						},
					},
					lastQuestion:            "Is this a question?",
					lastQuestionTargetBotId: "bot_id3",
				},
				sourceBotId: "bot_id3",
				targetBotId: "bot_id3",
				text:        "This is the answer",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "this game is not waiting for messages currently",
		},
		{
			name: "errors if game source bot asking question does not match expected source bot",
			input: struct {
				game        *Game
				sourceBotId string
				targetBotId string
				text        string
			}{
				game: &Game{
					state:            waitingForAiQuestion,
					currentTurnIndex: 1,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3"},
					bots: []*Bot{
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
							player: &Player{
								id: "player_id1",
							},
						},
					},
					lastQuestion:            "Is this a question?",
					lastQuestionTargetBotId: "bot_id3",
				},
				sourceBotId: "bot_id1",
				targetBotId: "bot_id3",
				text:        "what?",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "incorrect sourceBotId",
		},
		{
			name: "errors if game source bot answering question does not match expected source bot",
			input: struct {
				game        *Game
				sourceBotId string
				targetBotId string
				text        string
			}{
				game: &Game{
					state:            waitingForHumanAnswer,
					currentTurnIndex: 1,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3"},
					bots: []*Bot{
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
							player: &Player{
								id: "player_id1",
							},
						},
					},
					lastQuestion:            "Is this a question?",
					lastQuestionTargetBotId: "bot_id3",
				},
				sourceBotId: "bot_id1",
				targetBotId: "bot_id3",
				text:        "answer",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "incorrect sourceBotId",
		},
		{
			name: "errors if game is in an unexpected state where game expects an AI bot but the bot is HUMAN",
			input: struct {
				game        *Game
				sourceBotId string
				targetBotId string
				text        string
			}{
				game: &Game{
					state:            waitingForAiAnswer,
					currentTurnIndex: 1,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3"},
					bots: []*Bot{
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
							player: &Player{
								id: "player_id1",
							},
						},
					},
					lastQuestion:            "Is this a question?",
					lastQuestionTargetBotId: "bot_id3",
				},
				sourceBotId: "bot_id3",
				targetBotId: "bot_id3",
				text:        "answer",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "expecting AI message but did not receive one",
		},
		{
			name: "errors if game is in an unexpected state where game expects a HUMAN bot but the bot is AI",
			input: struct {
				game        *Game
				sourceBotId string
				targetBotId string
				text        string
			}{
				game: &Game{
					state:            waitingForHumanQuestion,
					currentTurnIndex: 1,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3"},
					bots: []*Bot{
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
							player: &Player{
								id: "player_id1",
							},
						},
					},
					lastQuestion:            "Is this a question?",
					lastQuestionTargetBotId: "bot_id3",
				},
				sourceBotId: "bot_id2",
				targetBotId: "bot_id2",
				text:        "Another question?",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "expecting Human message but did not receive one",
		},

		{
			name: "errors if answer has different source and target bot",
			input: struct {
				game        *Game
				sourceBotId string
				targetBotId string
				text        string
			}{
				game: &Game{
					state:            waitingForHumanAnswer,
					currentTurnIndex: 1,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3"},
					bots: []*Bot{
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
							player: &Player{
								id: "player_id1",
							},
						},
					},
					lastQuestion:            "Is this a question?",
					lastQuestionTargetBotId: "bot_id3",
				},
				sourceBotId: "bot_id3",
				targetBotId: "bot_id2",
				text:        "answer",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "answering message should have same source and target bot",
		},
		{
			name: "errors if question has same source and target bot",
			input: struct {
				game        *Game
				sourceBotId string
				targetBotId string
				text        string
			}{
				game: &Game{
					state:            waitingForAiQuestion,
					currentTurnIndex: 1,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3"},
					bots: []*Bot{
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
							player: &Player{
								id: "player_id1",
							},
						},
					},
					lastQuestion:            "Is this a question?",
					lastQuestionTargetBotId: "bot_id3",
				},
				sourceBotId: "bot_id2",
				targetBotId: "bot_id2",
				text:        "Another question?",
			},
			expectedOutput: nil,
			errorExpected:  true,
			errorString:    "questioning message should have different source and target bot",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.input.game.GetGameUpdateAfterIncomingMessage(tt.input.sourceBotId, tt.input.targetBotId, tt.input.text)
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

func Test_StateHasBeenHandled(t *testing.T) {
	tests := []struct {
		name           string
		input          *Game
		expectedOutput bool
	}{
		{
			name: "returns true",
			input: &Game{
				state:        started,
				stateHandled: true,
			},
			expectedOutput: true,
		},
		{
			name: "returns false",
			input: &Game{
				state:        started,
				stateHandled: false,
			},
			expectedOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.StateHasBeenHandled()
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func Test_IsInStatePlayersJoined(t *testing.T) {
	tests := []struct {
		name           string
		input          *Game
		expectedOutput bool
	}{
		{
			name: "returns true",
			input: &Game{
				state:        playersJoined,
				stateHandled: true,
			},
			expectedOutput: true,
		},
		{
			name: "returns false",
			input: &Game{
				state:        started,
				stateHandled: false,
			},
			expectedOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.IsInStatePlayersJoined()
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func Test_IsInStateWaitingForAiQuestion(t *testing.T) {
	tests := []struct {
		name           string
		input          *Game
		expectedOutput bool
	}{
		{
			name: "returns true",
			input: &Game{
				state:        waitingForAiQuestion,
				stateHandled: true,
			},
			expectedOutput: true,
		},
		{
			name: "returns false",
			input: &Game{
				state:        started,
				stateHandled: false,
			},
			expectedOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.IsInStateWaitingForAiQuestion()
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func Test_IsInStateWaitingForAiAnswer(t *testing.T) {
	tests := []struct {
		name           string
		input          *Game
		expectedOutput bool
	}{
		{
			name: "returns true",
			input: &Game{
				state:        waitingForAiAnswer,
				stateHandled: true,
			},
			expectedOutput: true,
		},
		{
			name: "returns false",
			input: &Game{
				state:        started,
				stateHandled: false,
			},
			expectedOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.IsInStateWaitingForAiAnswer()
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func Test_RandomizedTurnOrder(t *testing.T) {
	tests := []struct {
		name           string
		input          *Game
		expectedOutput []string
	}{
		{
			name: "returns randomized turn order",
			input: &Game{
				state: playersJoined,
				bots: []*Bot{
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
					{
						id:        "bot_id4",
						name:      "bot4",
						typeOfBot: ai,
					},
					{
						id:        "bot_id5",
						name:      "bot5",
						typeOfBot: ai,
					},
				},
			},
			expectedOutput: []string{"bot_id3", "bot_id4", "bot_id2", "bot_id1", "bot_id5"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rand.Seed(0)
			result := tt.input.RandomizedTurnOrder()
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func Test_RecentlyUpdated(t *testing.T) {
	timeNow := time.Now()
	timeOld := time.Now().Add(-5 * time.Hour)
	tests := []struct {
		name           string
		input          *Game
		expectedOutput bool
	}{
		{
			name: "returns true",
			input: &Game{
				state:        waitingForAiAnswer,
				stateHandled: true,
				updatedAt:    timeNow,
			},
			expectedOutput: true,
		},
		{
			name: "returns false",
			input: &Game{
				state:        started,
				stateHandled: false,
				updatedAt:    timeOld,
			},
			expectedOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.RecentlyUpdated()
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func Test_GetTargetBotForNextQuestion(t *testing.T) {
	tests := []struct {
		name           string
		input          *Game
		expectedOutput string
		errorExpected  bool
		errorString    string
	}{
		{
			name: "errors if game has no bots",
			input: &Game{
				state:            started,
				turnOrder:        []string{"bot_id1"},
				currentTurnIndex: 0,
			},
			expectedOutput: "",
			errorExpected:  true,
			errorString:    "cannot get target bot from an empty list",
		},
		{
			name: "errors if only one bot",
			input: &Game{
				state:            started,
				turnOrder:        []string{"bot_id1"},
				currentTurnIndex: 0,
				bots: []*Bot{
					{
						id:        "bot_id1",
						name:      "bot1",
						typeOfBot: ai,
					},
				},
			},
			expectedOutput: "",
			errorExpected:  true,
			errorString:    "cannot get target bot from an empty list",
		},
		{
			name: "returns random bot with least messages",
			input: &Game{
				state:            started,
				turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
				currentTurnIndex: 2,
				bots: []*Bot{
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
						typeOfBot: ai,
					},
					{
						id:        "bot_id4",
						name:      "bot4",
						typeOfBot: ai,
					},
					{
						id:        "bot_id5",
						name:      "bot5",
						typeOfBot: ai,
					},
				},
				messages: []*Message{
					{SourceBotId: "bot_id2", TargetBotId: "bot_id1", Text: "question 1", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id1", TargetBotId: "bot_id1", Text: "answer 1", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id2", TargetBotId: "bot_id1", Text: "question 2", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id1", TargetBotId: "bot_id1", Text: "answer 2", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id2", TargetBotId: "bot_id1", Text: "question 3", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id1", TargetBotId: "bot_id1", Text: "answer 3", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id1", TargetBotId: "bot_id2", Text: "question 1", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id2", TargetBotId: "bot_id2", Text: "answer 1", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id1", TargetBotId: "bot_id2", Text: "question 2", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id2", TargetBotId: "bot_id2", Text: "answer 2", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id1", TargetBotId: "bot_id3", Text: "question 1", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id3", TargetBotId: "bot_id3", Text: "answer 1", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id1", TargetBotId: "bot_id3", Text: "question 2", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id3", TargetBotId: "bot_id3", Text: "answer 2", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id1", TargetBotId: "bot_id3", Text: "question 3", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id3", TargetBotId: "bot_id3", Text: "answer 3", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id1", TargetBotId: "bot_id3", Text: "question 4", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id3", TargetBotId: "bot_id3", Text: "answer 4", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id1", TargetBotId: "bot_id3", Text: "question 5", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id3", TargetBotId: "bot_id3", Text: "answer 5", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id3", TargetBotId: "bot_id4", Text: "question 1", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id4", TargetBotId: "bot_id4", Text: "answer 1", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id3", TargetBotId: "bot_id4", Text: "question 2", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id4", TargetBotId: "bot_id4", Text: "answer 2", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id4", TargetBotId: "bot_id5", Text: "question 1", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "answer 1", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id4", TargetBotId: "bot_id5", Text: "question 2", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "answer 2", CreatedAt: time.Now(), MessageType: "answer"},
				},
			},
			expectedOutput: "bot_id4",
			errorExpected:  false,
			errorString:    "",
		},
		{
			name: "returns random bot with least messages excluding current turn bot",
			input: &Game{
				state:            started,
				turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
				currentTurnIndex: 1,
				bots: []*Bot{
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
						typeOfBot: ai,
					},
					{
						id:        "bot_id4",
						name:      "bot4",
						typeOfBot: ai,
					},
					{
						id:        "bot_id5",
						name:      "bot5",
						typeOfBot: ai,
					},
				},
				messages: []*Message{
					{SourceBotId: "bot_id2", TargetBotId: "bot_id1", Text: "question 1", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id1", TargetBotId: "bot_id1", Text: "answer 1", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id2", TargetBotId: "bot_id1", Text: "question 2", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id1", TargetBotId: "bot_id1", Text: "answer 2", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id1", TargetBotId: "bot_id2", Text: "question 1", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id2", TargetBotId: "bot_id2", Text: "answer 1", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id2", TargetBotId: "bot_id3", Text: "question 1", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id3", TargetBotId: "bot_id3", Text: "answer 1", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id2", TargetBotId: "bot_id3", Text: "question 2", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id3", TargetBotId: "bot_id3", Text: "answer 2", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id2", TargetBotId: "bot_id4", Text: "question 1", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id4", TargetBotId: "bot_id4", Text: "answer 1", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id2", TargetBotId: "bot_id4", Text: "question 2", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id4", TargetBotId: "bot_id4", Text: "answer 2", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id2", TargetBotId: "bot_id5", Text: "question 1", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "answer 1", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id2", TargetBotId: "bot_id5", Text: "question 2", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "answer 2", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id2", TargetBotId: "bot_id5", Text: "question 3", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "answer 3", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id2", TargetBotId: "bot_id5", Text: "question 4", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "answer 4", CreatedAt: time.Now(), MessageType: "answer"},
					{SourceBotId: "bot_id2", TargetBotId: "bot_id5", Text: "question 5", CreatedAt: time.Now(), MessageType: "question"},
					{SourceBotId: "bot_id5", TargetBotId: "bot_id5", Text: "answer 5", CreatedAt: time.Now(), MessageType: "answer"},
				},
			},
			expectedOutput: "bot_id3",
			errorExpected:  false,
			errorString:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rand.Seed(0)
			result, err := tt.input.GetTargetBotIdForNextQuestion()
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

func Test_GetBotThatGameIsWaitingOn(t *testing.T) {
	tests := []struct {
		name           string
		input          *Game
		expectedOutput *Bot
	}{
		{
			name: "returns bot with current turn when waiting for bot to ask question",
			input: &Game{
				state:            waitingForAiQuestion,
				currentTurnIndex: 1,
				turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3"},
				bots: []*Bot{
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
						player: &Player{
							id: "player_id1",
						},
					},
				},
			},
			expectedOutput: &Bot{
				id:        "bot_id2",
				name:      "bot2",
				typeOfBot: ai,
			},
		},
		{
			name: "returns bot with lastQuestionTargetBotId when waiting for bot to answer question",
			input: &Game{
				state:                   waitingForAiAnswer,
				currentTurnIndex:        1,
				turnOrder:               []string{"bot_id1", "bot_id2", "bot_id3"},
				lastQuestionTargetBotId: "bot_id1",
				bots: []*Bot{
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
						player: &Player{
							id: "player_id1",
						},
					},
				},
			},
			expectedOutput: &Bot{
				id:        "bot_id1",
				name:      "bot1",
				typeOfBot: ai,
			},
		},
		{
			name: "returns bot with current turn when waiting for human to ask question",
			input: &Game{
				state:            waitingForHumanQuestion,
				currentTurnIndex: 2,
				turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3"},
				bots: []*Bot{
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
						player: &Player{
							id: "player_id1",
						},
					},
				},
			},
			expectedOutput: &Bot{
				id:        "bot_id3",
				name:      "bot3",
				typeOfBot: human,
				player: &Player{
					id: "player_id1",
				},
			},
		},
		{
			name: "returns bot with lastQuestionTargetBotId when waiting for huma to answer question",
			input: &Game{
				state:                   waitingForHumanAnswer,
				currentTurnIndex:        1,
				turnOrder:               []string{"bot_id1", "bot_id2", "bot_id3"},
				lastQuestionTargetBotId: "bot_id3",
				bots: []*Bot{
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
						player: &Player{
							id: "player_id1",
						},
					},
				},
			},
			expectedOutput: &Bot{
				id:        "bot_id3",
				name:      "bot3",
				typeOfBot: human,
				player: &Player{
					id: "player_id1",
				},
			},
		},
		{
			name: "returns nil if game in unexpected state",
			input: &Game{
				state:                   playersJoined,
				currentTurnIndex:        1,
				turnOrder:               []string{"bot_id1", "bot_id2", "bot_id3"},
				lastQuestionTargetBotId: "bot_id3",
				bots: []*Bot{
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
						player: &Player{
							id: "player_id1",
						},
					},
				},
			},
			expectedOutput: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.GetBotThatGameIsWaitingOn()
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}
