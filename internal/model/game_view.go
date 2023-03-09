package model

import (
	"fmt"
	"time"

	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

type GameView struct {
	State          gameViewState
	DisplayMessage string
	StateStartedAt *time.Time
	StateTotalTime int64
	LastQuestion   string
	MyBotId        string
	Bots           []BotView
}

func (g *Game) GameViewForPlayer(playerId string) *GameView {
	if utilities.IsBlank(playerId) {
		return nil
	}

	myBot := g.BotWithPlayerId(playerId)
	if myBot == nil {
		return nil
	}

	myBotId := myBot.id
	bots := prepareBotViews(g.bots)

	state, displayMessage := convertGameStateToGameViewStateWithMessage(g, myBotId)

	return &GameView{
		State:          state,
		DisplayMessage: displayMessage,
		StateStartedAt: g.stateHandledAt,
		StateTotalTime: g.stateTotalTime,
		LastQuestion:   "no question",
		MyBotId:        myBotId,
		Bots:           bots,
	}
}

func prepareBotViews(bots []*Bot) []BotView {
	botViews := []BotView{}
	for _, bot := range bots {
		botViews = append(botViews, BotView{
			Id:       bot.id,
			Name:     bot.name,
			Messages: bot.messages,
		})
	}

	return botViews
}

func convertGameStateToGameViewStateWithMessage(g *Game, myBotId string) (gameViewState, string) {
	targetBot := g.getTargetBot()
	switch g.state {
	case started, playersJoined:
		return waitingForPlayersToJoin, "Please wait as players join in"
	case waitingForBotQuestion:
		return waitingOnBotToAskAQuestion, "Please wait as someone is asking a question"
	case waitingForBotAnswer:
		return waitingOnBotToAnswer,
			fmt.Sprintf("Please wait as %s is answering the question", targetBot.name)
	case waitingForPlayerQuestion:
		if g.getCurrentTurnBotId() == myBotId {
			return waitingOnYouToAskAQuestion, "Please type a question. OR Click suggest for help!"
		} else {
			return waitingOnBotToAskAQuestion, "Please wait as someone is asking a question"
		}
	case waitingForPlayerAnswer:
		if g.lastQuestionTargetBotId == myBotId {
			return waitingOnYouToAnswer, "Please answer the question. OR Click suggest for help!"
		} else {
			return waitingOnBotToAnswer,
				fmt.Sprintf("Please wait as %s is answering the question", targetBot.name)
		}
	case finished:
		return timeUp, "Time ran out"
	default:
		return undefinedGameViewState, "This is not supposed to happen. What did happen?"
	}
}
