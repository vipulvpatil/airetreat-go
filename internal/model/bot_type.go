package model

type botType int64

const (
	undefinedBotType botType = iota
	ai
	human
)

func BotType(str string) botType {
	switch str {
	case "AI":
		return ai
	case "HUMAN":
		return human
	default:
		return undefinedBotType
	}
}

func (b botType) String() string {
	switch b {
	case ai:
		return "AI"
	case human:
		return "HUMAN"
	default:
		return "UNDEFINED"
	}
}

func (b botType) Valid() bool {
	return b.String() != "UNDEFINED"
}

func (b botType) IsAi() bool {
	return b == ai
}

func (b botType) IsHuman() bool {
	return b == human
}
