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
	assert.Equal(t, expected.state, actual.state, "game state is not equal")
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
	for j, expectedMessage := range expected.messages {
		actualMessage := actual.messages[j]
		assert.Equal(t, expectedMessage.Text, actualMessage.Text, "bot message is not equal")
		AssertTimeAlmostEqual(t, actualMessage.CreatedAt, expectedMessage.CreatedAt, DELTA, "message createdAt is not within range")
	}
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
