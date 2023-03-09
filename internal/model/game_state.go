package model

type gameState int64

const (
	undefinedGameState gameState = iota
	started
	playersJoined
	waitingForAiQuestion
	waitingForAiAnswer
	waitingForHumanQuestion
	waitingForHumanAnswer
	finished
)

func GameState(str string) gameState {
	switch str {
	case "STARTED":
		return started
	case "PLAYERS_JOINED":
		return playersJoined
	case "WAITING_FOR_AI_QUESTION":
		return waitingForAiQuestion
	case "WAITING_FOR_AI_ANSWER":
		return waitingForAiAnswer
	case "WAITING_FOR_HUMAN_QUESTION":
		return waitingForHumanQuestion
	case "WAITING_FOR_HUMAN_ANSWER":
		return waitingForHumanAnswer
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
	case waitingForAiQuestion:
		return "WAITING_FOR_AI_QUESTION"
	case waitingForAiAnswer:
		return "WAITING_FOR_AI_ANSWER"
	case waitingForHumanQuestion:
		return "WAITING_FOR_HUMAN_QUESTION"
	case waitingForHumanAnswer:
		return "WAITING_FOR_HUMAN_ANSWER"
	case finished:
		return "FINISHED"
	default:
		return "UNDEFINED"
	}
}

func (s gameState) Valid() bool {
	return s.String() != "UNDEFINED"
}

func (s gameState) isWaitingForAi() bool {
	return s == waitingForAiQuestion || s == waitingForAiAnswer
}

func (s gameState) isWaitingForHuman() bool {
	return s == waitingForHumanQuestion || s == waitingForHumanAnswer
}

func (s gameState) isQuestion() bool {
	return s == waitingForAiQuestion || s == waitingForHumanQuestion
}

func (s gameState) isAnswer() bool {
	return s == waitingForAiAnswer || s == waitingForHumanAnswer
}

func (s gameState) isWaitingForMessage() bool {
	return s == waitingForAiQuestion || s == waitingForAiAnswer || s == waitingForHumanQuestion || s == waitingForHumanAnswer
}
