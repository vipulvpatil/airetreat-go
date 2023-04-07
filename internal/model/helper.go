package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const DELTA = 5 * time.Second

// These function helps in testing Game and Bot Objects
func AssertEqualGame(t *testing.T, expected, actual *Game) {
	assert.Equal(t, expected.id, actual.id, "game id is not equal")
	assert.Equal(t, expected.state, actual.state, "gam state is not equal")
	assert.Equal(t, expected.currentTurnIndex, actual.currentTurnIndex, "game currentTurnIndex is not equal")
	assert.Equal(t, expected.turnOrder, actual.turnOrder, "game turnOrder is not equal")
	assert.Equal(t, expected.stateHandled, actual.stateHandled, "game stateHandled is not equal")
	assert.Equal(t, expected.stateTotalTime, actual.stateTotalTime, "game stateTotalTime is not equal")
	assert.Equal(t, expected.lastQuestion, actual.lastQuestion, "game lastQuestion is not equal")
	assert.Equal(t, expected.lastQuestionTargetBotId, actual.lastQuestionTargetBotId, "game lastQuestionTargetBotId is not equal")

	// Since we cannot mock postgres time operations. We just check that the updated times are near expected times.
	AssertTimeAlmostEqual(t, actual.createdAt, expected.createdAt, DELTA, "game createdAt is not within range")
	AssertTimeAlmostEqual(t, actual.updatedAt, expected.updatedAt, DELTA, "game updatedAt is not within range")
	if expected.stateHandledAt != nil {
		AssertTimeAlmostEqual(t, *actual.stateHandledAt, *expected.stateHandledAt, DELTA, "game stateHandledAt is not within range")
	}

	for i, expectedBot := range expected.bots {
		actualBot := actual.bots[i]
		AssertEqualBot(t, expectedBot, actualBot)
	}
}

func AssertEqualBot(t *testing.T, expected, actual *Bot) {
	assert.Equal(t, expected.id, actual.id, "bot id is not equal")
	assert.Equal(t, expected.name, actual.name, "bot name is not equal")
	assert.Equal(t, expected.typeOfBot, actual.typeOfBot, "bot type is not equal")
	assert.Equal(t, expected.player, actual.player, "bot player is not equal")
	// TODO: correctly verify messages
	// for j, expectedMessage := range expected.messages {
	// 	actualMessage := actual.messages[j]
	// 	assert.Equal(t, expectedMessage.Text, actualMessage.Text, "bot message is not equal")
	// 	AssertTimeAlmostEqual(t, actualMessage.CreatedAt, expectedMessage.CreatedAt, DELTA, "message createdAt is not within range")
	// }
}

func AssertTimeAlmostEqual(t *testing.T, actual, expected time.Time, delta time.Duration, msgAndArgs ...interface{}) bool {
	return assert.WithinRange(
		t,
		actual,
		expected.Add(-1*delta),
		expected.Add(delta),
		msgAndArgs,
	)
}

func AssertEqualGameView(t *testing.T, expected, actual *GameView) {

	assert.Equal(t, expected.State, actual.State, "gameView State is not equal")
	assert.Equal(t, expected.DisplayMessage, actual.DisplayMessage, "gameView DisplayMessage is not equal")
	assert.Equal(t, expected.StateTotalTime, actual.StateTotalTime, "gameView StateTotalTime is not equal")
	assert.Equal(t, expected.LastQuestion, actual.LastQuestion, "gameView LastQuestion is not equal")
	assert.Equal(t, expected.MyBotId, actual.MyBotId, "gameView MyBotId is not equal")
	assert.Equal(t, expected.Bots, actual.Bots, "gameView Bots is not equal")

	// Since we cannot mock postgres time operations. We just check that the updated times are near expected times.
	if expected.StateStartedAt != nil {
		AssertTimeAlmostEqual(t, *actual.StateStartedAt, *expected.StateStartedAt, DELTA, "gameView StateStartedAt is not within range")
	}

	for i, expectedDetailMessage := range expected.DetailedMessages {
		actualDetailedMessage := actual.DetailedMessages[i]
		AssertEqualDetailedMessage(t, expectedDetailMessage, actualDetailedMessage)
	}
}

func AssertEqualDetailedMessage(t *testing.T, expected, actual DetailedMessage) {
	assert.Equal(t, expected.Text, actual.Text, "detailedMessage Text is not equal")
	assert.Equal(t, expected.SourceBotId, actual.SourceBotId, "detailedMessage SourceBotId is not equal")
	assert.Equal(t, expected.SourceBotName, actual.SourceBotName, "detailedMessage SourceBotName is not equal")
	assert.Equal(t, expected.TargetBotId, actual.TargetBotId, "detailedMessage TargetBotId is not equal")
	assert.Equal(t, expected.TargetBotName, actual.TargetBotName, "detailedMessage TargetBotName is not equal")
	assert.Equal(t, expected.MessageType, actual.MessageType, "detailedMessage MessageType is not equal")
	AssertTimeAlmostEqual(t, actual.CreatedAt, expected.CreatedAt, DELTA, "detailedMessage CreatedAt is not within range")
}
