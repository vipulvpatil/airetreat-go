package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// These function helps in testing Game Objects
func AssertEqualGame(t *testing.T, expected, actual *Game) {
	assert.Equal(t, expected.id, actual.id, "game id is not equal")
	assert.Equal(t, expected.state, actual.state, "game state is not equal")
	assert.Equal(t, expected.currentTurnIndex, actual.currentTurnIndex, "game currentTurnIndex is not equal")
	assert.Equal(t, expected.turnOrder, actual.turnOrder, "game turnOrder is not equal")
	assert.Equal(t, expected.stateHandled, actual.stateHandled, "game stateHandled is not equal")
	assert.Equal(t, expected.stateTotalTime, actual.stateTotalTime, "game stateTotalTime is not equal")
	assert.Equal(t, expected.lastQuestion, actual.lastQuestion, "game lastQuestion is not equal")
	assert.Equal(t, expected.lastQuestionTargetBotId, actual.lastQuestionTargetBotId, "game lastQuestionTargetBotId is not equal")
	assert.Equal(t, expected.id, actual.id, "game id is not equal")

	for i, expectedBot := range expected.bots {
		actualBot := actual.bots[i]
		assert.Equal(t, expectedBot, actualBot, "bots did not match")
	}

	// Since we cannot mock postgres time operations. We just check that the updated times are near expected times.
	delta := 5 * time.Second
	AssertTimeAlmostEqual(t, actual.createdAt, expected.createdAt, delta, "game createdAt is not within range")
	AssertTimeAlmostEqual(t, actual.updatedAt, expected.updatedAt, delta, "game updatedAt is not within range")
	if expected.stateHandledAt != nil {
		AssertTimeAlmostEqual(t, *actual.stateHandledAt, *expected.stateHandledAt, delta, "game stateHandledAt is not within range")
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
