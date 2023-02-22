package model

type gameViewState int64

const (
	undefinedGameViewState gameViewState = iota
	waitingForPlayersToJoin
	waitingOnBotToAskAQuestion
	waitingOnBotToAnswer
	waitingOnYouToAskAQuestion
	waitingOnYouToAnswer
	youLost
	youWon
	timeUp
)

func GameViewState(str string) gameViewState {
	switch str {
	case "WAITING_FOR_PLAYERS_TO_JOIN":
		return waitingForPlayersToJoin
	case "WAITING_ON_BOT_TO_ASK_A_QUESTION":
		return waitingOnBotToAskAQuestion
	case "WAITING_ON_BOT_TO_ANSWER":
		return waitingOnBotToAnswer
	case "WAITING_ON_YOU_TO_ASK_A_QUESTION":
		return waitingOnYouToAskAQuestion
	case "WAITING_ON_YOU_TO_ANSWER":
		return waitingOnYouToAnswer
	case "YOU_LOST":
		return youLost
	case "YOU_WON":
		return youWon
	case "TIME_UP":
		return timeUp
	default:
		return undefinedGameViewState
	}
}

func (s gameViewState) String() string {
	switch s {
	case waitingForPlayersToJoin:
		return "WAITING_FOR_PLAYERS_TO_JOIN"
	case waitingOnBotToAskAQuestion:
		return "WAITING_ON_BOT_TO_ASK_A_QUESTION"
	case waitingOnBotToAnswer:
		return "WAITING_ON_BOT_TO_ANSWER"
	case waitingOnYouToAskAQuestion:
		return "WAITING_ON_YOU_TO_ASK_A_QUESTION"
	case waitingOnYouToAnswer:
		return "WAITING_ON_YOU_TO_ANSWER"
	case youLost:
		return "YOU_LOST"
	case youWon:
		return "YOU_WON"
	case timeUp:
		return "TIME_UP"
	default:
		return "UNDEFINED"
	}
}

func (s gameViewState) Valid() bool {
	return s.String() != "UNDEFINED"
}
