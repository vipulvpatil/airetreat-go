package model

import (
	"time"

	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

type GameView struct {
	State          string
	DisplayMessage string
	StateStartedAt *time.Time
	StateTotalTime int64
	LastQuestion   string
	MyBotId        string
	Bots           []BotView
}

func (g *Game) ForPlayer(playerId string) *GameView {
	if g == nil {
		return nil
	}

	if utilities.IsBlank(playerId) {
		return nil
	}

	myBotId, bots := prepareBotsAndReturnMyBotId(g.bots, playerId)
	if myBotId == nil {
		return nil
	}

	return &GameView{
		State:          "Game STATE",
		DisplayMessage: "Game is in this state",
		StateStartedAt: g.stateHandledAt,
		StateTotalTime: 60,
		LastQuestion:   "no question",
		MyBotId:        *myBotId,
		Bots:           bots,
	}
}

func prepareBotsAndReturnMyBotId(bots []*Bot, playerId string) (*string, []BotView) {
	var myBotId string
	botViews := []BotView{}
	for _, bot := range bots {
		botViews = append(botViews, BotView{
			Id:      bot.id,
			Name:    bot.name,
			Message: bot.messages,
		})
		if bot.player.id == playerId {
			myBotId = bot.id
		}
	}

	return &myBotId, botViews
}
