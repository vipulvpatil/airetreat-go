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
					state: started,
					bots: []*Bot{
						{
							player: &Player{
								id: "player_id1",
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
				Bots: []BotView{
					{
						Id:   "",
						Name: "",
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
					state: started,
					bots: []*Bot{
						{
							player: &Player{
								id: "player_id1",
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
				Bots: []BotView{
					{
						Id:   "",
						Name: "",
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
