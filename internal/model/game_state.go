package model

type gameState int64

const (
	undefinedGameState gameState = iota
	started
	playersJoined
	waitingForBotQuestion
	waitingForBotAnswer
	waitingForPlayerQuestion
	waitingForPlayerAnswer
	finished
)

func GameState(str string) gameState {
	switch str {
	case "STARTED":
		return started
	case "PLAYERS_JOINED":
		return playersJoined
	case "WAITING_FOR_BOT_QUESTION":
		return waitingForBotQuestion
	case "WAITING_FOR_BOT_ANSWER":
		return waitingForBotAnswer
	case "WAITING_FOR_PLAYER_QUESTION":
		return waitingForPlayerQuestion
	case "WAITING_FOR_PLAYER_ANSWER":
		return waitingForPlayerAnswer
	case "FINISHED":
		return finished
	default:
		return undefinedGameState
	}
}

func (s gameState) String() string {
	switch s {
	case started:
		return "STARTED"
	case playersJoined:
		return "PLAYERS_JOINED"
	case waitingForBotQuestion:
		return "WAITING_FOR_BOT_QUESTION"
	case waitingForBotAnswer:
		return "WAITING_FOR_BOT_ANSWER"
	case waitingForPlayerQuestion:
		return "WAITING_FOR_PLAYER_QUESTION"
	case waitingForPlayerAnswer:
		return "WAITING_FOR_PLAYER_ANSWER"
	case finished:
		return "FINISHED"
	default:
		return "UNDEFINED"
	}
}

func (s gameState) Valid() bool {
	return s.String() != "UNDEFINED"
}
