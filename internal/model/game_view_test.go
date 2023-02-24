package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GameViewForPlayer(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			playerId string
			game     *Game
		}
		output        *GameView
		errorExpected bool
		errorString   string
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
				},
			},
			output: &GameView{
				State:          waitingForPlayersToJoin,
				DisplayMessage: "Please wait as players join in",
				StateTotalTime: 60,
				LastQuestion:   "no question",
				MyBotId:        "bot_id1",
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
			errorExpected: false,
			errorString:   "",
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
				},
			},
			output: &GameView{
				State:          waitingForPlayersToJoin,
				DisplayMessage: "Please wait as players join in",
				StateTotalTime: 60,
				LastQuestion:   "no question",
				MyBotId:        "bot_id1",
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
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "successfully returns a game state when game is waiting for Bot to ask question",
			input: struct {
				playerId string
				game     *Game
			}{
				playerId: "player_id1",
				game: &Game{
					state:            waitingForBotQuestion,
					turnOrder:        []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
					currentTurnIndex: 1,
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
				},
			},
			output: &GameView{
				State:          waitingOnBotToAskAQuestion,
				DisplayMessage: "Please wait as someone is asking a question",
				StateTotalTime: 60,
				LastQuestion:   "no question",
				MyBotId:        "bot_id1",
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
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "successfully returns a game state when game is waiting for Bot to answer a question",
			input: struct {
				playerId string
				game     *Game
			}{
				playerId: "player_id1",
				game: &Game{
					state:                   waitingForBotAnswer,
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
				},
			},
			output: &GameView{
				State:          waitingOnBotToAnswer,
				DisplayMessage: "Please wait as bot2 is answering the question",
				StateTotalTime: 60,
				LastQuestion:   "no question",
				MyBotId:        "bot_id1",
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
			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameView := tt.input.game.GameViewForPlayer(tt.input.playerId)
			assert.Equal(t, gameView, tt.output)
		})
	}
}
